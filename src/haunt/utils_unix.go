//go:build !windows
// +build !windows

package haunt

import "strings"

// isHidden check if a file or a path is hidden
func isHidden(path string) bool {
	arr := strings.Split(path[len(Wdir()):], "/")
	for _, elm := range arr {
		if strings.HasPrefix(elm, ".") {
			return true
		}
	}
	return false
}
