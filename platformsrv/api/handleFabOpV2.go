
package api

import (
	"os"
	"time"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/logger" 
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/k8s/deployfabric"
	"github.com/wingbaas/platformsrv/settings/fabric"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
) 

func checkOrg(org string,peers []string)bool {
	for _,peer := range peers {
		if strings.Contains(peer,org) {
			return true
		}
	}
	return false
}

func callCCV2(c echo.Context,cfg public.DeployPara,d FabricCCCallPara) error { 
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	// locker := getChainOpLocker(d.BlockChainId) 
	// locker.Lock()
	// defer locker.Unlock()
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "callCCV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "callCCV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_CALLCC
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "callCCV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain)  

	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "callCCV2 get tools image failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var orderAddr string
	var orderCaFile string
	for _,org := range cfg.DeployNetCfg.OrdererOrgs {
		for _,spec :=  range org.Specs {
			orderAddr = spec.Hostname + ":7050"
			orderCaFile = "/cert/crypto-config/ordererOrganizations/" + org.Domain + "/orderers/" + spec.Hostname + "." + org.Domain + "/msp/"
			orderCaFile = orderCaFile + "tlscacerts/tlsca." + org.Domain + "-cert.pem"
			break
		}
		break
	}
	var updatePara deployfabric.PathToolsDeploymentPara 
	updatePara.LogLevel = "debug"
	args := []string{"sh","-c"}
	var peerAddr string
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		if checkOrg(org.Name,d.Peers) {
			for _,spec :=  range org.Specs {
				tmpPeerAddr :=  " --peerAddresses  " + spec.Hostname + ":7051 --tlsRootCertFiles /cert/crypto-config/peerOrganizations/"
				tmpPeerAddr = tmpPeerAddr + org.Domain + "/peers/"
				tmpPeerAddr = tmpPeerAddr + spec.Hostname + "." + org.Domain + "/tls/ca.crt"
				peerAddr = peerAddr + tmpPeerAddr
				break
			}
		}
	}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		if org.Name == d.OrgName {
			updatePara.ToolsImage = toolsImage
			updatePara.OrgName = org.Name
			updatePara.PeerDomain = org.Domain
			for _,spec :=  range org.Specs {
				updatePara.PeerName = spec.Hostname
				break
			}
			break
		}
	}	
	argStr,err := generateArgJsonStr(d.Args)
	if err != nil {
		msg := "callCCV2 generateArgJsonStr failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	randomStr := utils.GenerateRandomString(8)
	outFileName := d.ChainCodeId + "-" + updatePara.OrgName + "-call.txt"
	outFile := "/var/data/src/" + outFileName 
	argStr = "'" + argStr + "'"
	exeCmd := "peer chaincode invoke -o " + orderAddr
	exeCmd = exeCmd + " --tls --cafile " + orderCaFile
	exeCmd = exeCmd + " -C " + d.ChannelId
	exeCmd = exeCmd + " -n " + d.ChainCodeId
	exeCmd = exeCmd + peerAddr
	exeCmd = exeCmd + " -c " + argStr
	exeCmd = exeCmd + " > " + outFile + " 2>&1"
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + exeCmd + ")" + ";"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	//logger.Debug("v2 cc call args=",updatePara.Args)
	_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
	if err != nil {
		msg := "callCCV2 update cli failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	//check the init result start
	time.Sleep(2*time.Second)
	parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + outFileName
	flag := "status:200"
	_,err = getResultByFlag(parseFile,flag)  
	if err != nil {
		msg := "callCCV2 status failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
		return c.JSON(http.StatusOK,ret)
	}
	callResult,_ := getInvokeResult(parseFile,"invoke result: ") 
	logger.Debug("invoke result=") 
	logger.Debug(callResult)
	//check the init result end
	os.Remove(parseFile) 
	rObj,err := fabric.OrgQueryBlockChain(d.BlockChainId,updatePara.OrgName,d.ChannelId)
	if err != nil {
		msg := "callCCV2 get chain height failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type Tst struct {
		Height uint64 `json:"Height"`
		CurrentBlockHash string `json:"CurrentBlockHash"`
		PreviousBlockHash string `json:"PreviousBlockHash"`
	}
	var st Tst
	b1,_ := json.Marshal(rObj)
	err = json.Unmarshal(b1,&st)
	if err != nil {
		msg := "callCCV2 chain height transfer to obj failed"
		logger.Errorf(msg)
		logger.Error(rObj)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if st.Height >0 {
		st.Height = st.Height -1
	}
	rObj2,err := fabric.OrgQueryBlockById(d.BlockChainId,updatePara.OrgName,d.ChannelId,st.Height)
	if err != nil {
		msg := "callCCV2 get block info failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var bInfo sdkfabric.BlockInfo
	b2,_ := json.Marshal(rObj2)
	err = json.Unmarshal(b2,&bInfo)
	if err != nil {
		msg := "callCCV2 transfer to block info failed"
		logger.Errorf(msg)
		logger.Error(rObj2)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var txId string
	for _,tx := range bInfo.Transactions {
		txId = tx.ChannelHeader.TxID
		break
	}
	//os.Remove(parseFile)
	type tst struct {
		TransactionID string `json:"TransactionID"`
		TxValidationCode int `json:"TxValidationCode"`
		ChaincodeStatus int `json:"ChaincodeStatus"`
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,tst{TransactionID:txId,TxValidationCode:0,ChaincodeStatus:200})
	return c.JSON(http.StatusOK,ret)
}

func ccQueryV2(c echo.Context,cfg public.DeployPara,d FabricCCQueryPara) error { 
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "ccQueryV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "ccQueryV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	chain.Status = k8s.CHAIN_STATUS_QUERY
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "ccQueryV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain) 

	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "ccQueryV2 get tools image failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var updatePara deployfabric.PathToolsDeploymentPara
	args := []string{"sh","-c"}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		if org.Name == d.OrgName {
			updatePara.ToolsImage = toolsImage
			updatePara.OrgName = org.Name
			updatePara.PeerDomain = org.Domain
			for _,spec :=  range org.Specs {
				updatePara.PeerName = spec.Hostname
				break
			}
			break
		}
	}	
	argStr,err := generateArgJsonStr(d.Args)
	if err != nil {
		msg := "ccQueryV2 generateArgJsonStr failed"
		logger.Error(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	randomStr := utils.GenerateRandomString(8)
	outFileName := d.ChainCodeId + "-" + updatePara.OrgName + "-query.txt"
	outFile := "/var/data/src/" + outFileName
	argStr = "'" + argStr + "'"
	exeCmd := "peer chaincode query -C " + d.ChannelId
	exeCmd = exeCmd + " -n " + d.ChainCodeId
	exeCmd = exeCmd + " -C " + d.ChannelId
	exeCmd = exeCmd + " -n " + d.ChainCodeId
	exeCmd = exeCmd + " -c " + argStr
	exeCmd = exeCmd + " > " + outFile + " 2>&1"
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + exeCmd + ")" + ";"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	//logger.Debug("v2 cc query args=",updatePara.Args)
	_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
	if err != nil {
		msg := "ccQueryV2 update cli failed"
		logger.Error(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	//check the query result start
	time.Sleep(1*time.Second)
	parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + outFileName
	result,err := getResult(parseFile)  
	if err != nil {
		msg := "ccQueryV2 failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)  
		return c.JSON(http.StatusOK,ret)
	}
	//check the query result end 
	os.Remove(parseFile)
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,result)
	return c.JSON(http.StatusOK,ret)
}

func ccQueryInstatialV2(c echo.Context,cfg public.DeployPara,d FabricInstantialCCPara) error { 
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "ccQueryInstatialV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "ccQueryInstatialV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	chain.Status = k8s.CHAIN_STATUS_QUERY
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "ccQueryInstatialV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain) 

	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "ccQueryInstatialV2 get tools image failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var orderAddr string
	var orderCaFile string
	for _,org := range cfg.DeployNetCfg.OrdererOrgs {
		for _,spec :=  range org.Specs {
			orderAddr = spec.Hostname + ":7050"
			orderCaFile = "/cert/crypto-config/ordererOrganizations/" + org.Domain + "/orderers/" + spec.Hostname + "." + org.Domain + "/msp/"
			orderCaFile = orderCaFile + "tlscacerts/tlsca." + org.Domain + "-cert.pem"
			break
		}
		break
	}
	var updatePara deployfabric.PathToolsDeploymentPara
	args := []string{"sh","-c"}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		if org.Name == d.OrgName {
			updatePara.ToolsImage = toolsImage
			updatePara.OrgName = org.Name
			updatePara.PeerDomain = org.Domain
			for _,spec :=  range org.Specs {
				updatePara.PeerName = spec.Hostname
				break
			}
			break
		}
	}	
	randomStr := utils.GenerateRandomString(8)
	outFileName := d.ChannelId + "-" + updatePara.OrgName + "-ccinstatial.txt" 
	outFile := "/var/data/src/" + outFileName
	exeCmd := "peer lifecycle chaincode querycommitted -o " + orderAddr
	exeCmd = exeCmd + " --tls --cafile " + orderCaFile
	exeCmd = exeCmd + " --channelID " + d.ChannelId
	exeCmd = exeCmd + " --output json"
	exeCmd = exeCmd + " > " + outFile + " 2>&1"
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + exeCmd + ")" + ";"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	//logger.Debug("ccQueryInstatialV2 args=",updatePara.Args)
	_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
	if err != nil {
		msg := "ccQueryInstatialV2 update cli failed"
		logger.Error(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	//check the query result start
	time.Sleep(5*time.Second) 
	parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + outFileName
	result,err := getResultInstatialObj(parseFile) 
	if err != nil {
		msg := "ccQueryInstatialV2 failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)   
		return c.JSON(http.StatusOK,ret)
	}
	//check the query result end 
	os.Remove(parseFile)
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,result)
	return c.JSON(http.StatusOK,ret)
}