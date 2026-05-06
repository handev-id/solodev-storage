package utils

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ResolveUploadKey resolves object key from query/form inputs or generates fallback key.
func ResolveUploadKey(c *fiber.Ctx, originalFilename string) (string, error) {
	key := NormalizeKey(c.Query("key"))
	if key == "" {
		key = NormalizeKey(c.FormValue("key"))
	}

	folder := NormalizeFolder(c.Query("folder"))
	if folder == "" {
		folder = NormalizeFolder(c.FormValue("folder"))
	}

	if strings.Contains(folder, "..") {
		return "", errors.New("invalid folder")
	}

	if key == "" {
		key = generateDefaultKey(originalFilename)
	}

	if strings.Contains(key, "..") {
		return "", errors.New("invalid key")
	}

	if folder != "" {
		key = path.Join(folder, key)
		key = NormalizeKey(key)
	}

	if key == "" || strings.Contains(key, "..") {
		return "", errors.New("invalid key")
	}

	return key, nil
}

// NormalizeKey normalizes object key path format.
func NormalizeKey(key string) string {
	clean := strings.TrimSpace(key)
	clean = strings.ReplaceAll(clean, "\\", "/")
	clean = strings.TrimPrefix(clean, "/")
	return clean
}

// NormalizeFolder normalizes folder path format for object prefix.
func NormalizeFolder(folder string) string {
	clean := NormalizeKey(folder)
	clean = strings.TrimSuffix(clean, "/")
	return clean
}

// SafeFilename extracts filename for download header.
func SafeFilename(key string) string {
	name := path.Base(key)
	if name == "" || name == "." || name == "/" {
		return "file"
	}
	return name
}

func generateDefaultKey(originalFilename string) string {
	name := filepath.Base(strings.TrimSpace(originalFilename))
	if name == "" || name == "." || name == "/" {
		name = "file.bin"
	}
	name = strings.ReplaceAll(name, " ", "-")

	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), name)
}
