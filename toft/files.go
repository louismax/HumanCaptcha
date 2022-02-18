package toft

import "os"

//GetRootDirectory 获取项目根目录
func GetRootDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

//PathExists 检查路径是否存在
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
