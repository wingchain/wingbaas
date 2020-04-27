
package sdkfabric

import (
	"fmt"
	"encoding/json"
	"sync"
	"io/ioutil"
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

