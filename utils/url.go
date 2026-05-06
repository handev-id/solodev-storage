package utils

import "strings"

// BuildPublicURL builds a public object URL from base URL and key.
func BuildPublicURL(baseURL, key string) string {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	k := strings.TrimLeft(strings.TrimSpace(key), "/")
	return base + "/" + k
}
