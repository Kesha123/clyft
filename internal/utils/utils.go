package utils

import (
	"os"
	"path/filepath"
)

const CLYFT_USER_STORAGE_PATH = ".local/share/clyft"
const CLYFT_ROOT_STORAGE_PATH = "/var/lib/clyft"

func GetClyftStoragePath() (string, error) {
	if user := os.Getuid(); user != 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, CLYFT_USER_STORAGE_PATH), nil
	}
	return CLYFT_ROOT_STORAGE_PATH, nil
}
