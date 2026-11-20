//go:build linux

package linuxdesk

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

const dirFlags = unix.O_RDONLY | unix.O_DIRECTORY | unix.O_CLOEXEC | unix.O_NOFOLLOW | unix.O_NONBLOCK

type scanLimits struct {
	depth, entries int
	time           time.Duration
}
type fileRow struct {
	name  string
	bytes uint64
}
type scanResult struct {
	count, skipped int
	bytes          uint64
	limited        bool
	top            []fileRow
}

// Every directory is opened relative to the previous validated descriptor;
// O_NOFOLLOW applies to each individual component (including the base).
func openLocalDir(path string) (int, error) {
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) || path == "/" {
		return -1, errors.New("invalid scan root")
	}
	fd, err := unix.Open("/", dirFlags, 0)
	if err != nil {
		return -1, err
	}
	for _, p := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		next, e := unix.Openat(fd, p, dirFlags, 0)
		_ = unix.Close(fd)
		if e != nil {
			return -1, e
		}
		fd = next
	}
	var stat unix.Statfs_t
	if e := unix.Fstatfs(fd, &stat); e != nil {
		_ = unix.Close(fd)
		return -1, e
	}
	switch uint32(stat.Type) {
	case unix.NFS_SUPER_MAGIC, unix.SMB2_SUPER_MAGIC, unix.FUSE_SUPER_MAGIC, 0xff534d42, 0x5346414f: // CIFS and legacy SMB
		_ = unix.Close(fd)
		return -1, errors.New("scan root is not on a known local filesystem")
	}
	return fd, nil
}

func scan(ctx context.Context, path string, uid int, temp bool, limits scanLimits) (scanResult, error) {
	fd, err := openLocalDir(path)
	if err != nil {
		return scanResult{}, err
	}
	base := os.NewFile(uintptr(fd), path)
	defer base.Close()

	var root unix.Stat_t
	if err := unix.Fstat(fd, &root); err != nil {
		return scanResult{}, err
	}
	var out scanResult
	deadline := time.Now().Add(limits.time)
	seen := 0

	budget := func() bool {
		if seen >= limits.entries || ctx.Err() != nil || !time.Now().Before(deadline) {
			out.limited = true
			return false
		}
		return true
	}
	record := func(name string, size uint64) {
		out.count++
		if temp {
			if ^uint64(0)-out.bytes < size {
				out.limited = true
				return
			}
			out.bytes += size
			return
		}
		if size < 50*1024*1024 {
			return
		}
		if len(out.top) >= 15 && size < out.top[14].bytes {
			return
		}
		out.top = append(out.top, fileRow{name, size})
		sort.Slice(out.top, func(i, j int) bool {
			if out.top[i].bytes == out.top[j].bytes {
				return out.top[i].name < out.top[j].name
			}
			return out.top[i].bytes > out.top[j].bytes
		})
		if len(out.top) > 15 {
			out.top = out.top[:15]
		}
	}
	var walk func(*os.File, string, int)
	walk = func(dir *os.File, rel string, depth int) {
		dfd := int(dir.Fd())
		for budget() {
			batch, e := dir.ReadDir(64)
			for _, entry := range batch {
				if !budget() {
					return
				}
				seen++
				name := entry.Name()
				displayPath := filepath.Join(rel, name)
				var stat unix.Stat_t
				if e := unix.Fstatat(dfd, name, &stat, unix.AT_SYMLINK_NOFOLLOW); e != nil || stat.Dev != root.Dev {
					out.skipped++
					continue
				}
				switch stat.Mode & unix.S_IFMT {
				case unix.S_IFDIR:
					if temp && int(stat.Uid) != uid {
						continue
					}
					if depth >= limits.depth {
						out.limited = true
						out.skipped++
						continue
					}
					next, e := unix.Openat(dfd, name, dirFlags, 0)
					if e != nil {
						out.skipped++
						continue
					}
					var opened unix.Stat_t
					if e := unix.Fstat(next, &opened); e != nil || opened.Dev != root.Dev || opened.Ino != stat.Ino || (temp && int(opened.Uid) != uid) {
						_ = unix.Close(next)
						out.skipped++
						continue
					}
					child := os.NewFile(uintptr(next), displayPath)
					walk(child, displayPath, depth+1)
					_ = child.Close()
				case unix.S_IFREG:
					if stat.Size >= 0 && (!temp || int(stat.Uid) == uid) {
						record(displayPath, uint64(stat.Size))
					}
				case unix.S_IFLNK:
					out.skipped++
				}
			}
			if e == io.EOF {
				return
			}
			if e != nil {
				out.skipped++
				return
			}
		}
	}
	walk(base, "", 0)
	return out, nil
}
