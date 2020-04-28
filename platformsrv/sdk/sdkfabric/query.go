
package sdkfabric

import (
	"fmt"
	//"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	// "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	// pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

//query chaincode
func (setup *FabricSetup)QueryCC(ccId string, fcn string, paras []string, peer string) ([]byte,error) {
	args := packArgs(paras)
	response, err := setup.ChClient.Query(channel.Request{ChaincodeID: ccId, Fcn: fcn, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryCC failed,err=%s",err)
		return nil,fmt.Errorf("QueryCC failed,err=%s",err)
	}
	logger.Debug("QueryCC result=")
	logger.Debug(string(response.Payload)) 
	return response.Payload,nil 
}

/*

// query block by block number
func (setup *FabricSetup)QueryBlockByNumber(blockID uint64, peer string) ([]byte,error) {
	block, err := setup.chClient.QueryBlock(blockID, ledger.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryBlockByNumber failed, err=%s", err)
		return nil,fmt.Errorf("QueryBlockByNumber failed, err=%s", err)
	}
	if block.Data == nil {
		logger.Errorf("QueryBlockByNumber block data is null")
		return nil,fmt.Errorf("QueryBlockByNumber block data is null")
	}
	return processBlock(block),nil
}

//get transaction by transaction ID
func QueryTransactionByID(txid string, peer string) ([]byte,error) {
	processedTransaction, err := setup.chClient.QueryTransaction(fab.TransactionID(txid), ledger.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryTransactionByID failed, err=%s", err)
		return nil,fmt.Errorf("QueryTransactionByID failed, err=%s", err)
	}
	transaction := processTransaction(processedTransaction.GetTransactionEnvelope())
	transaction.ValidationCode = uint8(processedTransaction.GetValidationCode())
	transaction.ValidationCodeName = pb.TxValidationCode_name[int32(transaction.ValidationCode)]
	transactionJSON, _ := json.Marshal(transaction)
	transactionJSONString, _ := Prettyprint(transactionJSON)
	return transactionJSONString,nil
}

//Query blockchain info
func QueryChainInfo(client *ledger.Client, endpoint string) []byte {
	blockchainInfo, _ := client.QueryInfo(ledger.WithTargetEndpoints(endpoint))
	type chainInfo struct {
		Height            uint64
		CurrentBlockHash  string
		PreviousBlockHash string
	}
	bci := chainInfo{}
	bci.Height = blockchainInfo.BCI.Height
	bci.CurrentBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.CurrentBlockHash)
	bci.PreviousBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.PreviousBlockHash)
	res, _ := json.Marshal(bci)
	return res
}

// QueryInstalledChaincodes query installed chaincode - must call with Admin context
func QueryInstalledChaincodes(client *resmgmt.Client, endpoint string) []byte {
	chaincodeQueryRes, err := client.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryInstalledChaincodes: %s", err)
	}
	res, _ := json.Marshal(chaincodeQueryRes)
	return res
}

// QueryInstantiatedChaincodes query instantiated chaincode - must call with Admin context
func QueryInstantiatedChaincodes(client *resmgmt.Client, channelName string, endpoint string) []byte {
	chaincodeQueryRes, err := client.QueryInstantiatedChaincodes(channelName, resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryInstantiatedChaincodes: %s", err)
	}
	res, _ := json.Marshal(chaincodeQueryRes)
	return res
}

// QueryChannels query channels - must call with Admin context
func QueryChannels(client *resmgmt.Client, endpoint string) []byte {
	channelQueryRes, err := client.QueryChannels(resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryChannels: %s", err)
	}
	res, _ := json.Marshal(channelQueryRes)
	return res
}
*/