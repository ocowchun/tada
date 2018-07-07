package utils

// copy from https://github.com/senorprogrammer/wtf
import (
	"errors"
	"os/user"
	"path/filepath"
)

func Home() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	if currentUser.HomeDir == "" {
		return "", errors.New("cannot find user-specific home dir")
	}

	return currentUser.HomeDir, nil
}

func ExpandHomeDir(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := Home()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}
