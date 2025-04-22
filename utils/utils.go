package utils

import "strings"

func SplitBySupportedPath(fullpath string, supportedPaths []string) (prefix, verb, suffix string, matched bool) {
	for _, path := range supportedPaths {
		pattern := "/" + path + "/"
		if idx := strings.Index(fullpath, pattern); idx != -1 {
			prefix = fullpath[:idx]
			verb = path
			suffix = fullpath[idx+len(pattern):]
			return prefix, suffix, verb, true
		}
	}
	return "", "", "", false // No match found
}
