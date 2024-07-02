package ioutil

import "os"

func PathExists(filePath string) (bool, error) {
	_, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
