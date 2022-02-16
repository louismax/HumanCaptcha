package toft

import "os"

func GetRootDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
