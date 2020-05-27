
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
	"github.com/hyperledger/fabric/protos/peer"
	proto "github.com/golang/protobuf/proto"
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
		time.Sleep(5*(time.Second))
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
		time.Sleep(5*time.Second)
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
		time.Sleep(5*time.Second)
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

func getInvokeResult(fileName string,flag string) (interface{},error) {
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
		time.Sleep(5*time.Second)
		bl,_ = utils.PathExists(fileName)
	}
	if !bl {
		logger.Errorf("getInvokeResult:Read file Failed")
		return "",fmt.Errorf("getInvokeResult:Read file Failed")
	}
	for i := 0; i < len(text); i++ {
		lineStr := text[i]
		//logger.Debug("getResultByFlag line str=",lineStr)
		if strings.Contains(lineStr,flag) {
			as := strings.Split(lineStr,flag) 
			if len(as) > 1 {
				bStr := "\n \327j\270\266>\204\327\014\001\225\314m\026\242&\377\321\315\363\027(\336E\r\277\201S\206m?\365Z\022\231\001\n\203\001\0227\n\n_lifecycle\022)\n'\n!namespaces/fields/cctest/Sequence\022\002\010\003\022H\n\006cctest\022>\n\026\n\020\000\364\217\277\277initialized\022\002\010\004\n\007\n\001a\022\002\010\004\n\007\n\001b\022\002\010\004\032\010\n\001a\032\003199\032\010\n\001b\032\003201\032\003\010\310\001\"\014\022\006cctest\032\00210"
				//st := &peer.ProposalResponse{}
				st := &peer.ProposalResponsePayload{}
				err := proto.Unmarshal([]byte(bStr),st)
				if err != nil {
					logger.Errorf("getInvokeResult:unmarshal to response err,result str=")  
					logger.Debug(as[1])
					return "",fmt.Errorf("getInvokeResult:unmarshal to response err")
				}
				logger.Debug("proposalhash=") 
				logger.Debugf("%s",st.ProposalHash) 
				chaincodeAction := &peer.ChaincodeAction{} 
				err = proto.Unmarshal(st.Extension,chaincodeAction)
				if err != nil {
					logger.Errorf("getInvokeResult:unmarshal to chaincodeAction err")  
					return "",fmt.Errorf("getInvokeResult:unmarshal to chaincodeAction err")
				}
				logger.Debug("chaincodeAction=") 
				logger.Debug(chaincodeAction) 
				return st,nil
			}else {
				logger.Errorf("getInvokeResult:result array length <2")  
				return "",fmt.Errorf("getInvokeResult:result array length <2")
			}
		} 
	}
	logger.Errorf("getInvokeResult: output file not find status: ",flag)  
	return "",fmt.Errorf("getInvokeResult: output file not find status: ",flag) 
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
		time.Sleep(5*time.Second)
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
		time.Sleep(5*time.Second)
		bl,_ = utils.PathExists(fileName)
	}
	if !bl {
		logger.Errorf("getResultObj:Read file Failed")
		return nil,fmt.Errorf("getResultObj:Read file Failed")
	}
	logger.Errorf("getResultObj:result not json obj")
	return nil,fmt.Errorf("getResultObj:result not json obj")
}






