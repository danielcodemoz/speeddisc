//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type StartupEntry struct {
	Name       string
	Command    string
	Where      string
	File       string
	Hive       string
	Key        string
	CanDisable bool
}

func currentSID() string {
	tok, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return ""
	}
	defer tok.Close()
	tu, err := tok.GetTokenUser()
	if err != nil || tu.User.Sid == nil {
		return ""
	}
	return tu.User.Sid.String()
}

func profileName(sid string) string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList\`+sid, registry.QUERY_VALUE)
	if err != nil {
		return sid
	}
	defer k.Close()
	p, _, err := k.GetStringValue("ProfileImagePath")
	if err != nil || p == "" {
		return sid
	}
	return filepath.Base(p)
}

func skipSID(sid string) bool {
	if sid == "" || strings.EqualFold(sid, ".DEFAULT") {
		return true
	}
	if strings.HasSuffix(strings.ToLower(sid), "_classes") {
		return true
	}
	switch strings.ToUpper(sid) {
	case "S-1-5-18", "S-1-5-19", "S-1-5-20":
		return true
	}
	return false
}

func openHive(hive, key string, access uint32) (registry.Key, error) {
	switch hive {
	case "HKLM":
		return registry.OpenKey(registry.LOCAL_MACHINE, key, access)
	case "HKCU":
		return registry.OpenKey(registry.CURRENT_USER, key, access)
	case "HKU":
		return registry.OpenKey(registry.USERS, key, access)
	default:
		return 0, fmt.Errorf("unknown hive")
	}
}

func listRunKey(hive, key, where string, out *[]StartupEntry, warns *[]string) {
	k, err := openHive(hive, key, registry.READ)
	if err != nil {
		if !errors.Is(err, registry.ErrNotExist) {
			*warns = append(*warns, where+": "+err.Error())
		}
		return
	}
	defer k.Close()
	names, err := k.ReadValueNames(0)
	if err != nil {
		*warns = append(*warns, where+": "+err.Error())
		return
	}
	sort.Strings(names)
	for _, name := range names {
		val, _, err := k.GetStringValue(name)
		e := StartupEntry{
			Name:       name,
			Command:    val,
			Where:      where,
			Hive:       hive,
			Key:        key,
			CanDisable: err == nil,
		}
		if err != nil {
			e.Command = "(not a text value)"
			e.CanDisable = false
		}
		*out = append(*out, e)
	}
}

func listStartupFolders() ([]StartupEntry, []string) {
	var dirs []string
	dirs = append(dirs, filepath.Join(programDataOrDefault(), `Microsoft\Windows\Start Menu\Programs\Startup`))
	if ad := os.Getenv("APPDATA"); ad != "" {
		dirs = append(dirs, filepath.Join(ad, `Microsoft\Windows\Start Menu\Programs\Startup`))
	}
	drive := os.Getenv("SystemDrive")
	if drive == "" {
		drive = `C:`
	}
	if des, err := os.ReadDir(drive + `\Users`); err == nil {
		for _, d := range des {
			if !d.IsDir() {
				continue
			}
			switch strings.ToLower(d.Name()) {
			case "public", "default", "default user", "all users":
				continue
			}
			dirs = append(dirs, filepath.Join(drive+`\Users`, d.Name(), `AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`))
		}
	}
	seenDir := map[string]bool{}
	var entries []StartupEntry
	var warns []string
	for _, dir := range dirs {
		n := normPath(dir)
		if seenDir[n] || !IsStartupDir(dir) {
			continue
		}
		seenDir[n] = true
		items, err := os.ReadDir(dir)
		if err != nil {
			if !os.IsNotExist(err) {
				warns = append(warns, dir+": "+err.Error())
			}
			continue
		}
		for _, it := range items {
			if it.IsDir() || skipStartupName(it.Name()) {
				continue
			}
			if it.Type()&os.ModeSymlink != 0 {
				continue
			}
			full := filepath.Join(dir, it.Name())
			entries = append(entries, StartupEntry{
				Name:       it.Name(),
				Command:    full,
				Where:      dir,
				File:       full,
				CanDisable: true,
			})
		}
	}
	return entries, warns
}

func listStartup() ([]StartupEntry, []string) {
	var entries []StartupEntry
	var warns []string
	const run = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	const run32 = `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Run`
	listRunKey("HKLM", run, `HKLM\`+run, &entries, &warns)
	listRunKey("HKLM", run32, `HKLM\`+run32, &entries, &warns)
	listRunKey("HKCU", run, `HKCU\`+run, &entries, &warns)
	listRunKey("HKCU", run32, `HKCU\`+run32, &entries, &warns)

	me := currentSID()
	uk, err := registry.OpenKey(registry.USERS, "", registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		warns = append(warns, "HKU: "+err.Error())
	} else {
		names, err := uk.ReadSubKeyNames(-1)
		uk.Close()
		if err != nil {
			warns = append(warns, "HKU: "+err.Error())
		} else {
			sort.Strings(names)
			for _, sid := range names {
				if skipSID(sid) || (me != "" && strings.EqualFold(sid, me)) {
					continue
				}
				key := sid + `\Software\Microsoft\Windows\CurrentVersion\Run`
				who := profileName(sid)
				where := fmt.Sprintf(`HKU\%s (%s)\Software\Microsoft\Windows\CurrentVersion\Run`, sid, who)
				listRunKey("HKU", key, where, &entries, &warns)
			}
		}
	}
	files, fw := listStartupFolders()
	entries = append(entries, files...)
	warns = append(warns, fw...)
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Where == entries[j].Where {
			return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
		}
		return entries[i].Where < entries[j].Where
	})
	return entries, warns
}

func disableRunValue(e StartupEntry) error {
	if !e.CanDisable || e.File != "" {
		return ErrNotText
	}
	k, err := openHive(e.Hive, e.Key, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	cur, _, err := k.GetStringValue(e.Name)
	if err != nil {
		return ErrNotText
	}
	if cur != e.Command {
		return ErrValueChanged
	}
	return k.DeleteValue(e.Name)
}
