package util

import "os"

func CheckDir(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}
	Logger.Warnf("Path %s not exists, trying to mkdir", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return true, err
	}
	return true, err
}
