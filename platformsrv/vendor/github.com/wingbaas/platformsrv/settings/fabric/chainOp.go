
package fabric

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
)

const (
	CHANNEL_DEFAULT_USER 	string = "chuser"
	USER_DEFAULT_SECRET 	string = "mysecret"
	USER_DEFAULT_TYPE 		string = "user"
)

func OrgInvokeChainCode(chainId string,orgName string,channelId string,ChainCodeID string,args []string) (interface{},error) {
	if len(args)<2 {
		logger.Errorf("OrgInovkeChainCode: len args<2 error")
		return nil,fmt.Errorf("OrgInovkeChainCode: len args<2 error")
	}
	fcn := args[0]
	paras := args[1:]
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgInovkeChainCode:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	for _,org := range obj.DeployNetCfg.PeerOrgs {
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer := p.Hostname + "." + org.Domain
				fSetup.Peers = append(fSetup.Peers,peer)
			}
			break
		}
	}
	err = fSetup.Initialize()
	if err != nil {
		return nil,fmt.Errorf("OrgInovkeChainCode:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 
	
	err = fSetup.SetupCA(orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgInovkeChainCode: set ca failed: org=%s\n", orgName,)
	}

	_,err = fSetup.GetRegisteredUser(CHANNEL_DEFAULT_USER,orgName,USER_DEFAULT_SECRET,USER_DEFAULT_TYPE)
	if err != nil {
		return nil,fmt.Errorf("OrgInovkeChainCode: GetRegisteredUser failed: org=%s\n", orgName)
	} 
	_,err = fSetup.GetUserClient(channelId,CHANNEL_DEFAULT_USER,orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgInovkeChainCode: GetUserClient failed: org=%s\n", orgName)
	}

	bytes,err := fSetup.ExecuteCC(ChainCodeID,fcn,paras,fSetup.Peers) 
	if err != nil {
		return nil,fmt.Errorf("OrgInovkeChainCode:invoke chaincode failed: org=%s\n", orgName)
	}
	var result sdkfabric.CCInvokeResult
	err = json.Unmarshal(bytes,&result)
	if err != nil {
		str := string(bytes)
		logger.Debug("OrgInovkeChainCode unmarshal error,result=")
		logger.Debug(str) 
		return str,fmt.Errorf("OrgInovkeChainCode unmarshal result error")
	}
	type rs struct {
		TransactionID string `json:"TransactionID"`
		TxValidationCode int `json:"TxValidationCode"`
		ChaincodeStatus int `json:"ChaincodeStatus"`
	} 
	r := rs{TransactionID: result.TransactionID,TxValidationCode: result.TxValidationCode,ChaincodeStatus: result.ChaincodeStatus}
	return r,nil
}

func OrgQueryChainCode(chainId string,orgName string,channelId string,ChainCodeID string,args []string)  (interface{},error) {
	if len(args)<2 {
		logger.Errorf("OrgQueryChainCode: len args<2 error")
		return nil,fmt.Errorf("OrgQueryChainCode: len args<2 error")
	}
	fcn := args[0]
	paras := args[1:]
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChainCode:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	var peer string
	for _,org := range obj.DeployNetCfg.PeerOrgs {
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer = p.Hostname + "." + org.Domain
				break
			}
			break
		}
	}
	err = fSetup.Initialize() 
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChainCode:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 

	err = fSetup.SetupCA(orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChainCode: set ca failed: org=%s\n", orgName,)
	}

	_,err = fSetup.GetRegisteredUser(CHANNEL_DEFAULT_USER,orgName,USER_DEFAULT_SECRET,USER_DEFAULT_TYPE)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChainCode: GetRegisteredUser failed: org=%s\n", orgName)
	} 
	_,err = fSetup.GetUserClient(channelId,CHANNEL_DEFAULT_USER,orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChainCode: GetUserClient failed: org=%s\n", orgName)
	}

	bytes,err := fSetup.QueryCC(ChainCodeID,fcn,paras,peer) 
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChainCode: query chaincode failed: org=%s\n", orgName)
	}
	str := string(bytes)
	logger.Debug("OrgQueryChainCode result=")
	logger.Debug(str) 
	return str,nil
}

func OrgQueryBlockChain(chainId string,orgName string,channelId string)  (interface{},error) {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockChain:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	var peer string
	for _,org := range obj.DeployNetCfg.PeerOrgs {
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer = p.Hostname + "." + org.Domain
				break
			}
			break
		}
	}
	err = fSetup.Initialize() 
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockChain:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 

	err = fSetup.SetupCA(orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockChain: set ca failed: org=%s\n", orgName,)
	}

	_,err = fSetup.GetRegisteredUser(CHANNEL_DEFAULT_USER,orgName,USER_DEFAULT_SECRET,USER_DEFAULT_TYPE)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockChain: GetRegisteredUser failed: org=%s\n", orgName)
	} 
	_,err = fSetup.GetUserClient(channelId,CHANNEL_DEFAULT_USER,orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockChain: GetUserClient failed: org=%s\n", orgName)
	}
	bci,err := fSetup.QueryChainInfo(CHANNEL_DEFAULT_USER,peer)  
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChainCode: query chain info failed: org=%s\n", orgName)
	}
	return bci,nil
}

func OrgQueryBlockById(chainId string,orgName string,channelId string,blockId uint64)  (interface{},error) {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockById:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	var peer string
	for _,org := range obj.DeployNetCfg.PeerOrgs {
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer = p.Hostname + "." + org.Domain
				break
			}
			break
		}
	}
	err = fSetup.Initialize()
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockById:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 

	err = fSetup.SetupCA(orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockById: set ca failed: org=%s\n", orgName,)
	}

	_,err = fSetup.GetRegisteredUser(CHANNEL_DEFAULT_USER,orgName,USER_DEFAULT_SECRET,USER_DEFAULT_TYPE)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockById: GetRegisteredUser failed: org=%s\n", orgName)
	} 
	_,err = fSetup.GetUserClient(channelId,CHANNEL_DEFAULT_USER,orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockById: GetUserClient failed: org=%s\n", orgName)
	}
	tx,err := fSetup.QueryBlockById(CHANNEL_DEFAULT_USER,blockId,peer)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryBlockById: query block info failed: org=%s\n", orgName)
	}
	return tx,nil
}

func OrgQueryTxById(chainId string,orgName string,channelId string,txId string)  (interface{},error) {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryTxById:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	var peer string
	for _,org := range obj.DeployNetCfg.PeerOrgs {
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer = p.Hostname + "." + org.Domain
				break
			}
			break
		}
	}
	err = fSetup.Initialize() 
	if err != nil {
		return nil,fmt.Errorf("OrgQueryTxById:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 

	err = fSetup.SetupCA(orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryTxById: set ca failed: org=%s\n", orgName,)
	}

	_,err = fSetup.GetRegisteredUser(CHANNEL_DEFAULT_USER,orgName,USER_DEFAULT_SECRET,USER_DEFAULT_TYPE)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryTxById: GetRegisteredUser failed: org=%s\n", orgName)
	} 
	_,err = fSetup.GetUserClient(channelId,CHANNEL_DEFAULT_USER,orgName)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryTxById: GetUserClient failed: org=%s\n", orgName)
	}
	tx,err := fSetup.QueryTransactionByID(CHANNEL_DEFAULT_USER,txId,peer)   
	if err != nil {
		return nil,fmt.Errorf("OrgQueryTxById: query tx info failed: org=%s\n", orgName)
	}
	return tx,nil
}

func OrgQueryChannel(chainId string,orgName string)  (interface{},error) {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChannel:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	var peer string
	for _,org := range obj.DeployNetCfg.PeerOrgs {
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer = p.Hostname + "." + org.Domain
				break
			}
			break
		}
	}
	err = fSetup.Initialize() 
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChannel:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 
	chs,err := fSetup.QueryChannels(peer)  
	if err != nil {
		return nil,fmt.Errorf("OrgQueryChannel: query channels failed: org=%s\n", orgName)
	}
	return chs,nil
}

func OrgQueryInstalledCC(chainId string,orgName string,channelId string)  (interface{},error) {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryInstalledCC:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	var peer string
	for _,org := range obj.DeployNetCfg.PeerOrgs { 
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer = p.Hostname + "." + org.Domain
				break
			}
			break
		}
	}
	err = fSetup.Initialize() 
	if err != nil {
		return nil,fmt.Errorf("OrgQueryInstalledCC:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 
	ccs,err := fSetup.QueryInstalledChaincodes(peer)   
	if err != nil {
		return nil,fmt.Errorf("OrgQueryInstalledCC: query failed: org=%s\n", orgName)
	}
	return ccs,nil
}

func OrgQueryInstantiateCC(chainId string,orgName string,channelId string)  (interface{},error) {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return nil,fmt.Errorf("OrgQueryInstantiateCC:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	var peer string
	for _,org := range obj.DeployNetCfg.PeerOrgs { 
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer = p.Hostname + "." + org.Domain
				break
			}
			break
		}
	}
	err = fSetup.Initialize() 
	if err != nil {
		return nil,fmt.Errorf("OrgQueryInstantiateCC:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 
	ccs,err := fSetup.QueryInstantiatedChaincodes(peer)  
	if err != nil {
		return nil,fmt.Errorf("OrgQueryInstantiateCC: query failed: org=%s\n", orgName)
	}
	return ccs,nil  
}