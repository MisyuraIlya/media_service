package system

import "os"

func CheckExistFile(path string) bool {
	_, err := os.Stat(path)

	if(os.IsNotExist(err)){
		return false
	}

	return err == nil
}