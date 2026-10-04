//go:build windows

package main

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type memoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var procGlobalMemoryStatusEx = windows.NewLazySystemDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")

func physicalMemory() (total, avail uint64, err error) {
	var m memoryStatusEx
	m.dwLength = uint32(unsafe.Sizeof(m))
	r, _, e := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&m)))
	if r == 0 {
		if e != nil {
			return 0, 0, e
		}
		return 0, 0, ErrMemory
	}
	return m.ullTotalPhys, m.ullAvailPhys, nil
}

func readOS() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	product, _, err := k.GetStringValue("ProductName")
	if err != nil {
		return "", err
	}
	display, _, _ := k.GetStringValue("DisplayVersion")
	build, _, _ := k.GetStringValue("CurrentBuild")
	ubr, _, ubrErr := k.GetIntegerValue("UBR")
	arch := os.Getenv("PROCESSOR_ARCHITECTURE")
	var b strings.Builder
	b.WriteString(strings.TrimSpace(product))
	if display != "" {
		b.WriteString(" ")
		b.WriteString(display)
	}
	if build != "" {
		b.WriteString(" (")
		b.WriteString(build)
		if ubrErr == nil {
			fmt.Fprintf(&b, ".%d", ubr)
		}
		b.WriteString(")")
	}
	if arch != "" {
		b.WriteString(" ")
		b.WriteString(arch)
	}
	return b.String(), nil
}

func readCPU() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	name, _, err := k.GetStringValue("ProcessorNameString")
	if err != nil {
		return "", err
	}
	return strings.Join(strings.Fields(name), " "), nil
}

type diskSpace struct {
	Letter string
	Free   uint64
	Total  uint64
}

func fixedDisks() ([]diskSpace, error) {
	mask, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, err
	}
	var out []diskSpace
	for i := 0; i < 26; i++ {
		if mask&(1<<uint(i)) == 0 {
			continue
		}
		root := string(rune('A'+i)) + `:\`
		ptr, err := windows.UTF16PtrFromString(root)
		if err != nil {
			continue
		}
		if windows.GetDriveType(ptr) != windows.DRIVE_FIXED {
			continue
		}
		var free, total, totalFree uint64
		if err := windows.GetDiskFreeSpaceEx(ptr, &free, &total, &totalFree); err != nil {
			continue
		}
		out = append(out, diskSpace{Letter: string(rune('A'+i)) + ":", Free: free, Total: total})
	}
	return out, nil
}

func estimateTemp() (int64, error) {
	var total int64
	var saw bool
	var first error
	paths := append(append([]string{}, userTempTargets()...), windowsTemp())
	seen := map[string]bool{}
	for _, p := range paths {
		n := normPath(p)
		if n == "" || seen[n] {
			continue
		}
		seen[n] = true
		if err := AllowClean(CleanTemp, p); err != nil {
			continue
		}
		nbyte, err := measureTree(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			if first == nil {
				first = err
			}
			continue
		}
		saw = true
		total += nbyte
	}
	if !saw && first != nil {
		return 0, first
	}
	return total, nil
}

func collectInfo(rep *Report) {
	if s, err := readOS(); err != nil {
		rep.addFail(msg{"fail_info", []any{tr("name_os"), errDetail(err)}})
	} else {
		rep.addFound(msg{"found_os", []any{s}})
	}
	if s, err := readCPU(); err != nil {
		rep.addFail(msg{"fail_info", []any{tr("name_cpu"), errDetail(err)}})
	} else if s != "" {
		rep.addFound(msg{"found_cpu", []any{s}})
	}
	if total, avail, err := physicalMemory(); err != nil {
		rep.addFail(msg{"fail_info", []any{tr("name_ram"), errDetail(err)}})
	} else {
		rep.addFound(msg{"found_ram", []any{size(avail), size(total)}})
	}
	disks, err := fixedDisks()
	if err != nil {
		rep.addFail(msg{"fail_info", []any{tr("name_disks"), errDetail(err)}})
	} else if len(disks) == 0 {
		rep.addFail(msg{"fail_info", []any{tr("name_disks"), errDetail(ErrNoDisk)}})
	} else {
		for _, d := range disks {
			rep.addFound(msg{"found_disk", []any{d.Letter, size(d.Free), size(d.Total)}})
		}
	}
	entries, warns := listStartup()
	rep.addFound(msg{"found_startup_count", []any{len(entries)}})
	if len(entries) > 0 {
		limit := 12
		if len(entries) < limit {
			limit = len(entries)
		}
		names := make([]string, 0, limit)
		for i := 0; i < limit; i++ {
			names = append(names, entries[i].Name)
		}
		sample := strings.Join(names, ", ")
		if len(entries) > 12 {
			sample += fmt.Sprintf(" (+%d)", len(entries)-12)
		}
		rep.addFound(msg{"found_startup_sample", []any{sample}})
	}
	for _, w := range warns {
		rep.addFail(msg{"fail_info", []any{tr("name_startup"), clip(w, 180)}})
	}
	if n, err := estimateTemp(); err != nil {
		rep.addFail(msg{"fail_info", []any{tr("name_temp"), errDetail(err)}})
		rep.addFound(msg{"found_temp", []any{size(n)}})
	} else {
		rep.addFound(msg{"found_temp", []any{size(n)}})
	}
}
