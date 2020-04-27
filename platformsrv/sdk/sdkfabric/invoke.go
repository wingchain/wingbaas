
package sdkfabric

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

//invoke chaincode

func (setup *FabricSetup)ExecuteCC(ccId string,fcn string, paras []string, peers []string) ([]byte,error) {
	args := packArgs(paras)
	response, err := setup.chClient.Execute(channel.Request{ChaincodeID: ccId, Fcn: fcn, Args: args},channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(peers...))
	if err != nil {
		logger.Errorf("ExecuteCC failed,err=%s",err)
		return nil,fmt.Errorf("ExecuteCC failed,err=%s",err)
	}
	res, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("ExecuteCC marshal response failed,err=%s",err)
		return nil,fmt.Errorf("ExecuteCC marshal response failed,err=%s",err)
	}
	logger.Debug("ExecuteCC result=")
	logger.Debug(string(res))
	return res,nil 
}
