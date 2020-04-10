
package utils

import (
	"os"
	"fmt"
	"io"
	"sync"
	"time"
	"strings"
	"io/ioutil"
	"math/rand"
	"github.com/wingbaas/platformsrv/logger"
)

var fileLocker sync.Mutex //file locker

func LoadFile(fileName string) ([]byte, error) {
	fileLocker.Lock()
	bytes, err := ioutil.ReadFile(fileName) //read file
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

//copy file to another dir
func CopyFile(dstName string, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	defer src.Close()
    if err != nil {
        return 0,fmt.Errorf("CopyFile open src file err, %v", err)
    }
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	defer dst.Close()
    if err != nil {
        return 0,fmt.Errorf("CopyFile open dest file err, %v", err)
    }
    return io.Copy(dst, src)
}

//get all file from the dir,include the sub dir
func GetAllFiles(dirPth string) (files []string, err error) {
    var dirs []string
    dir, err := ioutil.ReadDir(dirPth)
    if err != nil {
        return nil, err
    }
    PthSep := string(os.PathSeparator)
    for _, fi := range dir {
        if fi.IsDir() { // is dier
            dirs = append(dirs,dirPth+PthSep+fi.Name())
            GetAllFiles(dirPth + PthSep + fi.Name())
        } else {
			files = append(files, dirPth+PthSep+fi.Name())
        }
    }
    // traverse the sub dir
    for _, table := range dirs {
        temp, _ := GetAllFiles(table)
        for _, temp1 := range temp {
            files = append(files, temp1)
        }
    }
    return files, nil
}

func CopyDir(src string, dest string)(bool,error) {
	var newFile []string
    xfiles, _ := GetAllFiles(src)
    for _, file := range xfiles {
		destNew := strings.Replace(file,src,dest, -1)
		newFile = append(newFile,destNew)
		pos := strings.LastIndex(destNew,"/")
		tmpPath := substring(destNew,0,pos)
		bl,_ := PathExists(tmpPath)
		if !bl {
			bl,_ := CreateDir(tmpPath)
			if !bl {
				logger.Errorf("CopyDir create dir err, %s", tmpPath)
				return false,fmt.Errorf("CopyDir create dir err, %s", tmpPath)
			}
		}  
		_,err := CopyFile(destNew,file)
		if err!=nil {
			logger.Errorf("CopyDir file err, dest=%s   file=%s", destNew,file)
			return false,fmt.Errorf("CopyDir file err, dest=%s   file=%s", destNew,file)
		}
	}
	return true,nil
}

//get the sub string 
func substring(source string,start int, end int) string {
    var r = []rune(source)
    length := len(r)
    if start < 0 || end > length || start > end {
        return ""
    }
    if start == 0 && end == length {
        return source
    }
    return string(r[start:end])
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

func GetCaPrivateKey(chainId string,orgDomain string)(string,error) {
	priKey := ""
	certDir := BAAS_CFG.BlockNetCfgBasePath + chainId + "/crypto-config/peerOrganizations/" + orgDomain + "/ca/"
	files, err := ioutil.ReadDir(certDir)
	if err != nil {
		logger.Errorf("GetCaPrivateKey: err=%v",err)
		return priKey,fmt.Errorf("GetCaPrivateKey: err=%v",err)
	}
	for _, f := range files {
    	if strings.Contains(f.Name(),"_sk") {
			priKey = f.Name()
			return priKey,nil
    	}
	}
	if priKey == "" {
		logger.Errorf("GetCaPrivateKey: not find private key")
		return priKey,fmt.Errorf("GetCaPrivateKey: not find private key")
	}
	return priKey,nil 
}

