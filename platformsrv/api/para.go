
package api

import(
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
)

type FabricDeployCCPara struct {
	BlockChainId string `json:"BlockChainId"`
	OrgName string `json:"OrgName"`
	ChannelId string `json:"ChannelId"`
	ChainCodeId string `json:"ChainCodeId"`
	ChainCodeVersion string `json:"ChainCodeVersion"`
	InitArgs []string `json:"InitArgs"`
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
