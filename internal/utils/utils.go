package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	CLYFT_USER_STORAGE_PATH = ".local/share/clyft"
	CLYFT_ROOT_STORAGE_PATH = "/var/lib/clyft"
	CLYFT_AUTH_PATH         = "containers/auth.json"
)

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

type RegistryAuth struct {
	Auth  string `json:"auth"`
	Email string `json:"email,omitempty"`
}

type ConfigAuths struct {
	Auths map[string]RegistryAuth `json:"auths"`
}

func GetRegistryAuth(registry string) (*RegistryAuth, error) {
	user := os.Geteuid()
	var authJsonPath string
	if user != 0 {
		authJsonPath = filepath.Join(fmt.Sprintf("/run/user/%d/", user), CLYFT_AUTH_PATH)
	} else {
		authJsonPath = filepath.Join("/run/", CLYFT_AUTH_PATH)
		if _, err := os.Stat(authJsonPath); err != nil {
			authJsonPath = filepath.Join(fmt.Sprintf("/run/user/%d/", 0), CLYFT_AUTH_PATH)
		}
		if _, err := os.Stat(authJsonPath); err != nil {
			return nil, fmt.Errorf("couldn't find root auth.json")
		}
	}
	authJsonFile, err := os.Open(authJsonPath)
	if err != nil {
		return nil, fmt.Errorf("error openning %s: %w", authJsonPath, err)
	}
	defer authJsonFile.Close()
	var config ConfigAuths
	decoder := json.NewDecoder(authJsonFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}
	if registryAuth, ok := config.Auths[registry]; ok {
		return &registryAuth, nil
	}
	return nil, fmt.Errorf("registry credentials not found for %q", registry)
}

func GetRepoURI(tag string) string {
	if strings.LastIndex(tag, "/") == -1 {
		return ""
	}
	if index := strings.LastIndex(tag, ":"); index == -1 {
		return tag
	} else {
		return tag[:index]
	}
}
