package allutil

import (
	"io/ioutil"
	"os"
)

func FileIsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func String2File(str string, fileName string) error {
	tmpFile := fileName + ".tmp"
	err := String2FileNotSafe(str, tmpFile)
	if err != nil {
		return err
	}

	bakFile := fileName + ".bak"
	prevContent := File2String(fileName)
	if prevContent != "" {
		err = String2FileNotSafe(prevContent, bakFile)
		if err != nil {
			return err
		}
	}

	_ = os.Remove(fileName)
	err = os.Rename(tmpFile, fileName)
	if err != nil {
		return err
	}

	return nil
}

func String2FileNotSafe(str string, fileName string) error {
	return ioutil.WriteFile(fileName, []byte(str), os.ModePerm)
}

func File2String(fileName string) string {
	if FileIsExist(fileName) {
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			return ""
		}

		return string(b)
	} else {
		return ""
	}
}
