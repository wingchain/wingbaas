
package api

import(
	"fmt"
	"sync"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
)

const (
	LOOP_COUNT int = 30
)
var chainOpLocker map[string]sync.Mutex //chain operation locker
// var pkgLocker sync.Mutex //pkg file locker
// var deployLocker sync.Mutex //install file locker
// var flagLocker sync.Mutex //approve,commit,call,file locker
// var queryLocker sync.Mutex //query file locker

type FabricDeployCCPara struct {
	BlockChainId string `json:"BlockChainId"`
	OrgName string `json:"OrgName"`
	ChannelId string `json:"ChannelId"`
	ChainCodeId string `json:"ChainCodeId"`
	ChainCodeVersion string `json:"ChainCodeVersion"`
	InitArgs []string `json:"InitArgs"`
}

type FabricCCCallPara struct {
	BlockChainId string `json:"BlockChainId"`
	OrgName string `json:"OrgName"`
	ChannelId string `json:"ChannelId"`
	ChainCodeId string `json:"ChainCodeId"`
	Args []string `json:"Args"`
}

type FabricCCQueryPara struct {
	BlockChainId string `json:"BlockChainId"`
	OrgName string `json:"OrgName"`
	ChannelId string `json:"ChannelId"`
	ChainCodeId string `json:"ChainCodeId"`
	Args []string `json:"Args"`
}

type FabricInstantialCCPara struct {
	BlockChainId string `json:"BlockChainId"`
	OrgName string `json:"OrgName"`
	ChannelId string `json:"ChannelId"`
}

func getChainOpLocker(chainId string)sync.Mutex {
	if chainOpLocker == nil {
		chainOpLocker = make(map[string]sync.Mutex)
	}
	v,ok := chainOpLocker[chainId]
	if ok {
		return v
	}else{
		var locker sync.Mutex
		chainOpLocker[chainId] = locker
	}
	return chainOpLocker[chainId]
}

func generateArgJsonStr(args []string) (string,error) {
	type ArgSt struct{
		Function string `json:"function"`
		Args []string `json:"Args"`
	}
	if len(args)<1 {
		logger.Errorf("generateArgJsonStr: args need at least one")
		return "",fmt.Errorf("generateArgJsonStr: args need at least one")
	}
	var obj ArgSt
	obj.Function = args[0]
	if len(args)>1 {
		obj.Args = append(obj.Args,args[1:]...)
	}
	bytes,err := json.Marshal(obj)
	if err != nil {
		logger.Errorf("generateArgJsonStr: marshal args obj err")
		return "",fmt.Errorf("generateArgJsonStr: marshal args obj err")
	}
	return string(bytes),nil
} 

func getResultPkg(fileName string) (string,error) {
	// pkgLocker.Lock()
	// defer pkgLocker.Unlock()
	var text []string
	bl,_ := utils.PathExists(fileName)
	for i:=0;i<LOOP_COUNT;i++ {
		if bl {
			_,text = sdkfabric.ReadFileLine2(fileName) 
			if len(text) < 1 {
				return "",nil
			}else{
				logger.Errorf("getResultPkg: err msg =",text)
				return "",fmt.Errorf("getResultPkg: pkg file Failed")
			}
		}
		time.Sleep(2*time.Second)
		bl,_ = utils.PathExists(fileName)
	}
	logger.Errorf("getResultPkg: get result file time out")
	return "",fmt.Errorf("getResultPkg: get result file time out") 
}

func getResultByIdentifier(fileName string,flag string) (string,error) {
	// deployLocker.Lock()
	// defer deployLocker.Unlock() 
	var text []string
	bl,_ := utils.PathExists(fileName)
	for i:=0;i<LOOP_COUNT;i++ {
		if bl {
			_,text = sdkfabric.ReadFileLine2(fileName)	
			if len(text) > 0 {
				break
			}
		}
		time.Sleep(2*time.Second)
		bl,_ = utils.PathExists(fileName)
	}
	if !bl {
		logger.Errorf("getResultByIdentifier:Read file Failed")
		return "",fmt.Errorf("getResultByIdentifier:Read file Failed")
	}
	for i := 0; i < len(text); i++ {
		lineStr := text[i] 
		//logger.Debug("getResultByIdentifier line str=",lineStr)
		if strings.Contains(lineStr,"identifier:") {
			if strings.Contains(lineStr,flag) {
				//logger.Debug("getResultByIdentifier find flag line str=",lineStr)
				var ret string
				as := strings.Split(lineStr,flag)
				//logger.Debug("getResultByIdentifier line str array=",as)
				if len(as) == 1 {
					ret = as[0]
				}else if len(as) > 1 {
					ret = as[1]
				}
				if len([]rune(ret)) >= 3 {
					ret = strings.Replace(ret,"\n","",-1)
					//logger.Debug("getResultByIdentifier,ret ident=",ret) 
					return ret,nil
				}
			}
		}
	}
	logger.Errorf("getResultByIdentifier: output file not find identifier: ",flag)
	return "",fmt.Errorf("getResultByIdentifier: output file not find identifier: ",flag) 
}

func getResultByFlag(fileName string,flag string) (string,error) {
	// flagLocker.Lock()
	// defer flagLocker.Unlock()
	var text []string
	bl,_ := utils.PathExists(fileName)
	for i:=0;i<LOOP_COUNT;i++ {
		if bl {
			_,text = sdkfabric.ReadFileLine2(fileName)	
			if len(text) > 0 {
				break
			}
		}
		time.Sleep(2*time.Second)
		bl,_ = utils.PathExists(fileName)
	}
	if !bl {
		logger.Errorf("getResultByFlag:Read file Failed")
		return "",fmt.Errorf("getResultByFlag:Read file Failed")
	}
	for i := 0; i < len(text); i++ {
		lineStr := text[i]
		//logger.Debug("getResultByFlag line str=",lineStr)
		if strings.Contains(lineStr,flag) {
			return "",nil
		}
	}
	logger.Errorf("getResultByFlag: output file not find status: ",flag)  
	return "",fmt.Errorf("getResultByFlag: output file not find status: ",flag) 
}

func getResult(fileName string) (interface{},error) {
	// queryLocker.Lock()
	// defer queryLocker.Unlock()
	var text []string
	bl,_ := utils.PathExists(fileName)
	for i:=0;i<LOOP_COUNT;i++ {
		if bl {
			_,text = sdkfabric.ReadFileLine2(fileName)	
			if len(text) > 0 {
				break
			}
		}
		time.Sleep(2*time.Second)
		bl,_ = utils.PathExists(fileName)
	}
	if !bl {
		logger.Errorf("getResult:Read file Failed")
		return nil,fmt.Errorf("getResult:Read file Failed")
	}
	if len(text) == 1 {
		return text[0],nil
	}
	return text,nil
}

func getResultInstatialObj(fileName string) (interface{},error) {
	// queryLocker.Lock()
	// defer queryLocker.Unlock()
	type OutJson struct {
		ChaincodeDefinitions []struct {
			Collections struct {
			} `json:"collections"`
			EndorsementPlugin string `json:"endorsement_plugin"`
			InitRequired bool `json:"init_required"`
			Name string `json:"name"`
			Sequence int `json:"sequence"`
			ValidationParameter string `json:"validation_parameter"`
			ValidationPlugin string `json:"validation_plugin"`
			Version string `json:"version"`
		} `json:"chaincode_definitions"`
	}
	type InstatialCCSt struct {
		Name string `json:"name"`
		Version string `json:"version"`
		Path string `json:"path,omitempty"`
		Input string `json:"input,omitempty"`
		Escc string `json:"escc,omitempty"`
		Vscc string `json:"vscc,omitempty"`
		Sequence int `json:"sequence,omitempty"`
	}
	type RS struct{
		Chaincodes []InstatialCCSt `json:"chaincodes"`
	}
	var rs RS
	bl,_ := utils.PathExists(fileName)
	for i:=0;i<LOOP_COUNT;i++ {
		if bl {
			bytes, err := ioutil.ReadFile(fileName) //read file
			if err == nil && len(bytes) > 1{
				var obj OutJson
				err = json.Unmarshal(bytes,&obj)
				if err == nil {
					for _,it := range obj.ChaincodeDefinitions {
						var r InstatialCCSt
						r.Name = it.Name
						r.Version = it.Version
						r.Sequence = it.Sequence
						rs.Chaincodes = append(rs.Chaincodes,r)
					}
					return rs,nil
				}
			}
		}
		time.Sleep(2*time.Second)
		bl,_ = utils.PathExists(fileName)
	}
	if !bl {
		logger.Errorf("getResultObj:Read file Failed")
		return nil,fmt.Errorf("getResultObj:Read file Failed")
	}
	logger.Errorf("getResultObj:result not json obj")
	return nil,fmt.Errorf("getResultObj:result not json obj")
}






