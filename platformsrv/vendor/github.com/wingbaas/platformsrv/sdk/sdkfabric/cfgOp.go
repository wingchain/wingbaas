
package sdkfabric

import (
	"fmt"
	"io"
	"os"
	"sync"
	"bufio"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric/public" 
)

var fileLocker sync.Mutex //file locker

func LoadChainCfg(chainId string) (public.DeployPara, error) {
	fileName := utils.BAAS_CFG.BlockNetCfgBasePath + chainId + ".json"
	var obj public.DeployPara
	fileLocker.Lock()
	bytes, err := ioutil.ReadFile(fileName) //read file
	fileLocker.Unlock()
	if err != nil {
		logger.Errorf("LoadChainCfg: read cfg file error,%s", err) 
		return obj, fmt.Errorf("%v", err) 
	}
	err = json.Unmarshal(bytes,&obj)
	if err != nil {
		logger.Errorf("LoadChainCfg: unmarshal obj error,%s", err)
		return obj, fmt.Errorf("%s", err) 
	}
	return obj, nil
}

func ReadFileLine(filePath string)(bool, []string ){
	var content []string
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		logger.Errorf("Open File Failed, %v", err)
		return false,[]string{""}
	}	
	reader := bufio.NewReader(file)
    for {
    	var line string
        line, err := reader.ReadString('\n')
        if err != nil && io.EOF != err {
        	logger.Errorf("Read File Line Failed, %v", err)
        	return false,[]string{""}
        } else if ( io.EOF == err ) {
        	break
        }
        if ( "\n" != line ){
        	content = append( content ,line)
        }     
    }
	return true,content
}

func ReWriteFileLine(filePath string,content []string)bool {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
    	logger.Errorf("Open File Failed, %v", err)
        return false
    }	
	defer file.Close()
	for i := 0;i<len(content); i++ {
		_,err = file.WriteString(content[i])
		if err!=nil {
			logger.Errorf("Write File Content Failed, %v", err)
			return false
		}
	}
	return true
}

func AddHosts(hostName string, ip string) bool {
	fileLocker.Lock()
	bl, text := ReadFileLine("/etc/hosts")
	fileLocker.Unlock()
	if !bl {
		logger.Errorf("AddHosts:Read file Failed")
		return false
	}
	exits := false
	newHost := ip + " " + hostName + "\n"
	for i := 0; i < len(text); i++ {
		hosts := string([]byte(text[i])[strings.LastIndex(text[i], " ")+1:])
		if strings.Contains(hosts,hostName) {
			exits = true
			break
		}
	}
	if !exits {
		text = append(text,newHost)
	}
	bl = ReWriteFileLine("/etc/hosts", text)
	if !bl {
		logger.Errorf("AddHosts:Write file Failed")
	}
	return bl
}

