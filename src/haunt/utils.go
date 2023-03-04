package haunt

import (
	"log"
	"os"
	"strings"
)

// Split each arguments in multiple fields
func split(args, fields []string) []string {
	for _, arg := range fields {
		arr := strings.Fields(arg)
		args = append(args, arr...)
	}
	return args
}

// Get file extensions
func ext(path string) string {
	var ext string
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			ext = path[i:]
			if index := strings.LastIndex(ext, "."); index > 0 {
				ext = ext[index:]
			}
		}
	}
	if ext != "" {
		return ext[1:]
	}
	return ""
}

// Replace if isn't empty and create a new array
func replace(a []string, b string) []string {
	if len(b) > 0 {
		return strings.Fields(b)
	}
	return a
}

// Wdir return current working directory
func Wdir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	return dir
}
