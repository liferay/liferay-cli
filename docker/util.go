//go:build linux || darwin

package docker

import (
	"os"
	"strconv"
	"syscall"
)

func GetOsPathGid(path string) *string {
	info, _ := os.Stat(path)
	val := ""

	if info != nil {
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			val = strconv.FormatUint(uint64(stat.Gid), 10)
		}
	}

	return &val
}
