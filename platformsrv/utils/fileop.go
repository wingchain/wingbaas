
package utils

import (
	"os"
	"fmt"
	io "io/ioutil"
	"sync"
	"math/rand"
	"time"
	"github.com/wingbaas/platformsrv/logger"
)

var fileLocker sync.Mutex //file locker

func LoadFile(fileName string) ([]byte, error) {
	fileLocker.Lock()
	bytes, err := io.ReadFile(fileName) //read file
	fileLocker.Unlock()
	if err != nil {
		logger.Errorf("LoadFile: read file error,%v", err)
		return nil, fmt.Errorf("%v", err) 
	}
	return bytes, nil
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

func CreateDir(dirPath string) (bool,error){
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		logger.Errorf("CreateDir: MkdirAll error,%v", err)
		return false,fmt.Errorf("create dir failed,dir=%v, %v", dirPath, err)
	}
	return true,nil 
}

func DirCheck(dirPath string ) error {
	exsist,_ := PathExists(dirPath)
	if exsist {
		return nil
	} else {
		_,err := CreateDir(dirPath)
		if err != nil {
			logger.Errorf("DirCheck: CreateDir error,%v", err)
			return fmt.Errorf("%v", err)
		} else {
			return nil
		}
	} 
	return nil
}

func WriteFile(filePath string,content string) error {
	fileLocker.Lock()
	defer fileLocker.Unlock()
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		logger.Errorf("write file create failed, %v", err)
		return fmt.Errorf("%v", err)
	}
	_,err = file.WriteString(content)
	if err != nil {
		logger.Errorf("write file content failed, %v", err)
		return fmt.Errorf("%v", err)
	}
	return nil
}

func GenerateRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result) 
}