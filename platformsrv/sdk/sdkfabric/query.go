
package sdkfabric

import (
	"fmt"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ChainInfo struct {
	Height            uint64
	CurrentBlockHash  string
	PreviousBlockHash string
}

//query chaincode
func (setup *FabricSetup) QueryCC(ccId string, fcn string, paras []string, peer string) ([]byte,error) {
	args := packArgs(paras)
	response, err := setup.ChClient.Query(channel.Request{ChaincodeID: ccId, Fcn: fcn, Args: args},channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryCC failed,err=%s",err)
		return nil,fmt.Errorf("QueryCC failed,err=%s",err)
	}
	// logger.Debug("QueryCC result=")
	// logger.Debug(string(response.Payload)) 
	return response.Payload,nil  
}

// query block by block number
func (setup *FabricSetup) QueryBlockById(userName string,blockId uint64, peer string) (interface{},error) {
	channelContext := setup.Sdk.ChannelContext(setup.ChannelId,fabsdk.WithUser(userName),fabsdk.WithOrg(setup.OrgName))
	client,err := ledger.New(channelContext)
	if err != nil {
		logger.Errorf("QueryBlockByNumber: Failed to create new ledger client: %s", err)
		return nil,fmt.Errorf("QueryBlockByNumber: Failed to create new ledger client")
	}
	block, err := client.QueryBlock(blockId,ledger.WithTargetEndpoints(peer))
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
func (setup *FabricSetup) QueryTransactionByID(userName string,txId string, peer string) (interface{},error) {
	channelContext := setup.Sdk.ChannelContext(setup.ChannelId,fabsdk.WithUser(userName),fabsdk.WithOrg(setup.OrgName))
	client,err := ledger.New(channelContext)
	if err != nil {
		logger.Errorf("QueryTransactionByID: Failed to create new ledger client: %s", err)
		return nil,fmt.Errorf("QueryTransactionByID: Failed to create new ledger client")
	}
	processedTransaction, err := client.QueryTransaction(fab.TransactionID(txId), ledger.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryTransactionByID failed, err=%s", err)
		return nil,fmt.Errorf("QueryTransactionByID failed, err=%s", err)
	}
	transaction := processTransaction2(processedTransaction.GetTransactionEnvelope())
	transaction.ValidationCode = uint8(processedTransaction.GetValidationCode()) 
	transaction.ValidationCodeName = pb.TxValidationCode_name[int32(transaction.ValidationCode)]
	return transaction,nil 
}

//Query blockchain info
func (setup *FabricSetup) QueryChainInfo(usrName string,peer string) (interface{},error) { 
	channelContext := setup.Sdk.ChannelContext(setup.ChannelId,fabsdk.WithUser(usrName),fabsdk.WithOrg(setup.OrgName))
	client,err := ledger.New(channelContext) 
	if err != nil {
		logger.Errorf("QueryChainInfo: Failed to create new ledger client: %s", err)
		return nil,fmt.Errorf("QueryChainInfo: Failed to create new ledger client")
	}
	blockchainInfo,err := client.QueryInfo(ledger.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryChainInfo: Failed to query info: %s", err)
		return nil,fmt.Errorf("QueryChainInfo: Failed to query info")
	}
	bci := ChainInfo{}
	bci.Height = blockchainInfo.BCI.Height
	bci.CurrentBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.CurrentBlockHash)
	bci.PreviousBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.PreviousBlockHash)
	return bci,nil
}

// QueryInstalledChaincodes query installed chaincode
func (setup *FabricSetup) QueryInstalledChaincodes(peer string) (interface{},error) {
	chaincodeQueryRes,err := setup.netAdmin.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryInstalledChaincodes: Failed %s", err)
		return nil,fmt.Errorf("QueryInstalledChaincodes:failed")
	}
	return chaincodeQueryRes,nil
}

// QueryInstantiatedChaincodes query instantiated chaincode
func  (setup *FabricSetup) QueryInstantiatedChaincodes(peer string) (interface{},error) {
	chaincodeQueryRes, err := setup.netAdmin.QueryInstantiatedChaincodes(setup.ChannelId, resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryInstantiatedChaincodes: Failed %s", err)
		return nil,fmt.Errorf("QueryInstantiatedChaincodes: Failed")
	}
	return chaincodeQueryRes,nil
}

// QueryChannels query channels
func (setup *FabricSetup) QueryChannels(peer string) (interface{},error) {
	channelQueryRes,err := setup.netAdmin.QueryChannels(resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		logger.Errorf("QueryChannels: Failed %s", err)
		return nil,fmt.Errorf("QueryChannels failed")
	} 
	return channelQueryRes,nil
}
