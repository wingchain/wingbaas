
package api

import (
	"os"
	"io"
	"time"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/logger" 
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/k8s/deployfabric"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
)

func upChainCodeV2(c echo.Context,cfg public.DeployPara) error {  
	logger.Debug("upChainCodeV2")
	chainId := c.FormValue("BlockChainId") 
	ccId := c.FormValue("ChainCodeId")
	ccVersion := c.FormValue("ChainCodeVersion") 
	chain,_ := k8s.GetChain(chainId,cfg.ClusterId)
	if chain == nil {
		msg := "upChainCodeV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	file, err := c.FormFile("file") 
	if err != nil {
		msg := "get upload file error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	src, err := file.Open()
	if err != nil {
		msg := "open upload file error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	defer src.Close()
	rootPath,_ := utils.GetProcessRunRoot() 
	dstDir := rootPath + "/tmp/" + chainId + "/src/" + ccId + ccVersion + "/"    
	bl,_ := utils.CreateDir(dstDir)
	if !bl {
		msg := "create cc dir error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}  
	dst, err := os.Create(dstDir + file.Filename) 
	if err != nil {
		msg := "create cc file error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	defer dst.Close()
	_,err = io.Copy(dst, src)
	if err != nil {
		msg := "write cc file error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	nfsDir := utils.BAAS_CFG.NfsLocalRootDir + chainId + "/src/" + ccId + ccVersion + "/" 
	bl,_ = utils.CreateDir(nfsDir) 
	//logger.Debug("nfs cc dir=%s",nfsDir)
	if !bl {
		msg := "create nfs cc dir error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	_,err = utils.CopyDir(dstDir,nfsDir)
	if err != nil {
		msg := "copy cc file to nfs error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var updatePara deployfabric.PathToolsDeploymentPara
	ccPath := "/var/data/src/" + ccId + ccVersion + "/*.go"
	ccDstPath := "/opt/gopath/src/github.com/chaincode/" + ccId + ccVersion
	ccRelativePath := "github.com/chaincode/" + ccId + ccVersion
	args := []string{"sh","-c"}
	cmd := "mkdir -p " + ccDstPath + "; cp -a /var/data/. /cert; cp " + ccPath + " " + ccDstPath + "; cp -a /data/ccbase/. " + ccDstPath + ";"
	cmd = cmd + " peer lifecycle chaincode package " + ccId + ccVersion + ".tar.gz --path " + ccRelativePath + " --label " + ccId + ccVersion + ";"
	cmd = cmd + " cp " + ccId + ccVersion + ".tar.gz" + " /var/data/src/;"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		break
	}
	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "get tools image failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	updatePara.ToolsImage = toolsImage
	_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,chainId,updatePara)
	if err != nil {
		msg := "up chaincode to cli failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func orgDeployCCV2(c echo.Context,cfg public.DeployPara,d FabricDeployCCPara) error { 
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "orgDeployCCV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "orgDeployCCV2 get tools image failed"
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
	ccPkgName := d.ChainCodeId + d.ChainCodeVersion + ".tar.gz"
	args := []string{"sh","-c"}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		var updatePara deployfabric.PathToolsDeploymentPara
		updatePara.ToolsImage = toolsImage
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		cmd := "cp -a /var/data/. /cert; cp /var/data/src/" + ccPkgName + " . ;"
		cmd = cmd + " peer lifecycle chaincode install " + ccPkgName + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgDeployCCV2 cc install update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}

		// cmd = cmd + " peer lifecycle chaincode approveformyorg -o " + orderAddr
		// cmd = cmd + " --tls --cafile " + orderCaFile
		// cmd = cmd + " --channelID " + d.ChannelId
		// cmd = cmd + " --name " + d.ChainCodeId
		// cmd = cmd + " --version " + d.ChainCodeVersion
		// cmd = cmd + " --init-required --sequence 1;"
		// cmd = cmd + " /bin/bash"
		// updatePara.Args = append(updatePara.Args,args...)
		// updatePara.Args = append(updatePara.Args,cmd)
		// _,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		// if err != nil {
		// 	msg := "orgDeployCCV2 cc approve update cli failed"
		// 	logger.Error(msg)
		// 	ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		// 	return c.JSON(http.StatusOK,ret)
		// }
	}
	time.Sleep(5*time.Second)
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		var updatePara deployfabric.PathToolsDeploymentPara
		updatePara.ToolsImage = toolsImage
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		cmd := "cp -a /var/data/. /cert; peer lifecycle chaincode commit -o " + orderAddr
		cmd = cmd + " --channelID " + d.ChannelId
		cmd = cmd + " --name " + d.ChainCodeId
		cmd = cmd + " --version " + d.ChainCodeVersion
		cmd = cmd + " --sequence 1 --init-required --tls --cafile " + orderCaFile
		cmd = cmd + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		logger.Debug("commit args=",updatePara.Args)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgDeployCCV2 commit update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		argStr,err := generateArgJsonStr(d.InitArgs)
		if err != nil {
			msg := "orgDeployCCV2 generateArgJsonStr failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		argStr = "'" + argStr + "'"
		initCmd := "cp -a /var/data/. /cert; peer chaincode invoke -o " + orderAddr
		initCmd = initCmd + " --tls --cafile " + orderCaFile
		initCmd = initCmd + " -C " + d.ChannelId
		initCmd = initCmd + " -n " + d.ChainCodeId
		initCmd = initCmd + " --isInit -c " + argStr
		initCmd = initCmd + ";"
		initCmd = initCmd + " /bin/bash"
		args = append(args,initCmd)
		updatePara.Args = args
		logger.Debug("init args=",updatePara.Args)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgDeployCCV2 init update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		break
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}