
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chainId := c.FormValue("BlockChainId") 
	ccId := c.FormValue("ChainCodeId")
	ccVersion := c.FormValue("ChainCodeVersion") 
	chain,_ := k8s.GetChain(chainId,cfg.ClusterId)
	if chain == nil {
		msg := "upChainCodeV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	// locker := getChainOpLocker(chainId)
	// locker.Lock()
	// defer locker.Unlock()

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "upChainCodeV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_PKGCC
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "upChainCodeV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain)  

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
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		break
	}
	outFileName := ccId + ccVersion + "-" + updatePara.OrgName + "-pkg.txt"
	outFile := "/var/data/src/" + ccId + ccVersion + "/" + outFileName
	ccPath := "/var/data/src/" + ccId + ccVersion + "/*.go"
	ccDstPath := "/opt/gopath/src/github.com/chaincode/" + ccId + ccVersion
	ccRelativePath := "github.com/chaincode/" + ccId + ccVersion
	args := []string{"sh","-c"}
	exeCmd := "peer lifecycle chaincode package " + ccId + ccVersion + ".tar.gz --path " + ccRelativePath + " --label " + ccId + ccVersion
	exeCmd = exeCmd + " > " + outFile + " 2>&1"
	cmd := "mkdir -p " + ccDstPath + "; cp -a /var/data/. /cert; cp " + ccPath + " " + ccDstPath + "; cp -a /data/ccbase/. " + ccDstPath + ";"
	cmd = cmd + " $(" + exeCmd + ")" + ";"
	cmd = cmd + " cp " + ccId + ccVersion + ".tar.gz" + " /var/data/src/;"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
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
	time.Sleep(5*time.Second)
	//check the pkg result start
	parseFile := utils.BAAS_CFG.NfsLocalRootDir + chainId + "/src/" + ccId + ccVersion + "/" + outFileName
	_,err = getResultPkg(parseFile)
	if err != nil {
		msg := "pkg cc file failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	//check the pkg result end
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func orgDeployCCV2(c echo.Context,cfg public.DeployPara,d FabricDeployCCPara) error { 
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "orgDeployCCV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "orgDeployCCV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_DEPLOYCC
	err := k8s.UpdateChainStatus(*chain) 
	if err != nil {
		msg := "orgDeployCCV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain)  

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
	args := []string{"sh","-c"}
	ccPkgName := d.ChainCodeId + d.ChainCodeVersion + ".tar.gz"
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		var updatePara deployfabric.PathToolsDeploymentPara
		curOrg := org.Name
		updatePara.ToolsImage = toolsImage
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + curOrg + "-install.txt"
		outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		exeCmd := "peer lifecycle chaincode install " + ccPkgName + " > " + outFile + " 2>&1" 
		cmd := "cp -a /var/data/. /cert; cp /var/data/src/" + ccPkgName + " . ;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash" 
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgDeployCCV2 cc install update cli failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		time.Sleep(5*time.Second)
		parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName 
		flag := d.ChainCodeId + d.ChainCodeVersion + ":"
		pkgIdent,err := getResultByIdentifier(parseFile,flag) 
		logger.Debug("cc identifier=",pkgIdent)
		if err != nil {
			msg := "orgDeployCCV2 get cc identifier failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
		outFileName = d.ChainCodeId + d.ChainCodeVersion + "-" + curOrg + "-approve.txt"
		outFile = "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		pkgIdent = flag + pkgIdent
		exeCmd = "peer lifecycle chaincode approveformyorg -o " + orderAddr
		exeCmd = exeCmd + " --tls --cafile " + orderCaFile
		exeCmd = exeCmd + " --channelID " + d.ChannelId
		exeCmd = exeCmd + " --name " + d.ChainCodeId
		exeCmd = exeCmd + " --version " + d.ChainCodeVersion 
		exeCmd = exeCmd + " --init-required --sequence 1 --package-id " + pkgIdent
		//exeCmd = exeCmd + " --signature-policy \"" + d.EndorsePolicy + "\"" 
		exeCmd = exeCmd + " > " + outFile + " 2>&1"
		cmd = "cp -a /var/data/. /cert;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = nil
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgDeployCCV2 cc approve update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		time.Sleep(5*time.Second)
		parseFile = utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		flag = "(VALID)"
		_,err = getResultByFlag(parseFile,flag) 
		if err != nil {
			msg := "orgDeployCCV2 get cc approve status failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
	}
	time.Sleep(5*time.Second)
	var peerAddr string
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		for _,spec :=  range org.Specs {
			tmpPeerAddr :=  " --peerAddresses  " + spec.Hostname + ":7051 --tlsRootCertFiles /cert/crypto-config/peerOrganizations/"
			tmpPeerAddr = tmpPeerAddr + org.Domain + "/peers/"
			tmpPeerAddr = tmpPeerAddr + spec.Hostname + "." + org.Domain + "/tls/ca.crt"
			peerAddr = peerAddr + tmpPeerAddr
			break
		}
	}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		var updatePara deployfabric.PathToolsDeploymentPara
		updatePara.ToolsImage = toolsImage
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.OrgName + "-commit.txt"
		outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		exeCmd := "peer lifecycle chaincode commit -o " + orderAddr
		exeCmd = exeCmd + " --channelID " + d.ChannelId
		exeCmd = exeCmd + " --name " + d.ChainCodeId
		exeCmd = exeCmd + " --version " + d.ChainCodeVersion 
		exeCmd = exeCmd + peerAddr
		//exeCmd = exeCmd + " --signature-policy \"" + d.EndorsePolicy + "\""  
		exeCmd = exeCmd + " --sequence 1 --init-required --tls --cafile " + orderCaFile
		exeCmd = exeCmd + " > " + outFile + " 2>&1"
		cmd := "cp -a /var/data/. /cert;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		//logger.Debug("commit args=",updatePara.Args)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgDeployCCV2 commit update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		//check the commit result start
		time.Sleep(5*time.Second)
		parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		flag := "(VALID)"
		_,err = getResultByFlag(parseFile,flag) 
		if err != nil {
			msg := "orgDeployCCV2 get cc commit status failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
		//check the commit result end

		argStr,err := generateArgJsonStr(d.InitArgs)
		if err != nil {
			msg := "orgDeployCCV2 generateArgJsonStr failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		outFileName = d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.OrgName + "-init.txt"
		outFile = "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		argStr = "'" + argStr + "'"
		exeCmd = "peer chaincode invoke -o " + orderAddr
		exeCmd = exeCmd + " --tls --cafile " + orderCaFile
		exeCmd = exeCmd + " -C " + d.ChannelId
		exeCmd = exeCmd + " -n " + d.ChainCodeId
		exeCmd = exeCmd + " --isInit -c " + argStr
		exeCmd = exeCmd + " > " + outFile + " 2>&1"
		cmd = "cp -a /var/data/. /cert;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = nil
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd) 
		//logger.Debug("init args=",updatePara.Args)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgDeployCCV2 init update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		//check the init result start 
		time.Sleep(5*time.Second)
		parseFile = utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		flag = "status:200"
		_,err = getResultByFlag(parseFile,flag) 
		if err != nil {
			msg := "orgDeployCCV2 get cc init status failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
		//check the init result end
		break
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func singleOrgInstallCCV2(c echo.Context,cfg public.DeployPara,d FabricDeployCCPara) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "singleOrgInstallCCV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "singleOrgInstallCCV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_DEPLOYCC
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "singleOrgInstallCCV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain)  

	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "singleOrgInstallCCV2 get tools image failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var updatePara deployfabric.PathToolsDeploymentPara
	args := []string{"sh","-c"}
	ccPkgName := d.ChainCodeId + d.ChainCodeVersion + ".tar.gz"
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
	outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.OrgName + "-install.txt"
	outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
	exeCmd := "peer lifecycle chaincode install " + ccPkgName + " > " + outFile + " 2>&1" 
	cmd := "cp -a /var/data/. /cert; cp /var/data/src/" + ccPkgName + " . ;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + exeCmd + ")" + ";"
	cmd = cmd + " /bin/bash" 
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
	if err != nil {
		msg := "singleOrgInstallCCV2 cc install update cli failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	time.Sleep(5*time.Second)
	parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName 
	flag := d.ChainCodeId + d.ChainCodeVersion + ":"
	pkgIdent,err := getResultByIdentifier(parseFile,flag) 
	logger.Debug("cc identifier=",pkgIdent)
	if err != nil {
		msg := "singleOrgInstallCCV2 get cc identifier failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,flag + pkgIdent)
	return c.JSON(http.StatusOK,ret) 
}

func orgApproveCCV2(c echo.Context,cfg public.DeployPara,d FabricDeployCCPara) error { 
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "orgApproveCCV2:not find this chain" 
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()
	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "orgApproveCCV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_DEPLOYCC
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "orgApproveCCV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE 
	defer k8s.UpdateChainStatus(*chain)  
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
	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "orgApproveCCV2 get tools image failed"
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
	randomStr := utils.GenerateRandomString(8)
	outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.OrgName + "-approve.txt"
	outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
	exeCmd := "peer lifecycle chaincode approveformyorg -o " + orderAddr
	exeCmd = exeCmd + " --tls --cafile " + orderCaFile
	exeCmd = exeCmd + " --channelID " + d.ChannelId
	exeCmd = exeCmd + " --name " + d.ChainCodeId
	exeCmd = exeCmd + " --version " + d.ChainCodeVersion 
	exeCmd = exeCmd + " --sequence " + d.ChaincodeSeq
	exeCmd = exeCmd + " --init-required --package-id " + d.ChaincodePkg
	//exeCmd = exeCmd + " --signature-policy \"" + d.EndorsePolicy + "\"" 
	exeCmd = exeCmd + " > " + outFile + " 2>&1"
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + exeCmd + ")" + ";"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
	if err != nil {
		msg := "orgApproveCCV2 cc approve update cli failed"
		logger.Error(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	time.Sleep(5*time.Second)
	parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
	flag := "(VALID)"
	_,err = getResultByFlag(parseFile,flag) 
	if err != nil {
		msg := "orgApproveCCV2 get cc approve status failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret) 
}

func orgCommitCCV2(c echo.Context,cfg public.DeployPara,d FabricDeployCCPara) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "orgCommitCCV2:not find this chain"  
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()
	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "orgCommitCCV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_DEPLOYCC
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "orgCommitCCV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE 
	defer k8s.UpdateChainStatus(*chain)  
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
	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "orgCommitCCV2 get tools image failed"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var peerAddr string
	var updatePara deployfabric.PathToolsDeploymentPara
	args := []string{"sh","-c"}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		if org.Name == d.OrgName {
			updatePara.ToolsImage = toolsImage
			updatePara.OrgName = org.Name
			updatePara.PeerDomain = org.Domain
			for _,spec :=  range org.Specs {
				updatePara.PeerName = spec.Hostname
				tmpPeerAddr :=  " --peerAddresses  " + spec.Hostname + ":7051 --tlsRootCertFiles /cert/crypto-config/peerOrganizations/"
				tmpPeerAddr = tmpPeerAddr + org.Domain + "/peers/"
				tmpPeerAddr = tmpPeerAddr + spec.Hostname + "." + org.Domain + "/tls/ca.crt"
				peerAddr = peerAddr + tmpPeerAddr
				break
			}
			break
		}
	}
	randomStr := utils.GenerateRandomString(8)
	outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.OrgName + "-commit.txt"
	outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
	exeCmd := "peer lifecycle chaincode commit -o " + orderAddr
	exeCmd = exeCmd + " --channelID " + d.ChannelId
	exeCmd = exeCmd + " --name " + d.ChainCodeId
	exeCmd = exeCmd + " --version " + d.ChainCodeVersion
	exeCmd = exeCmd + " --sequence " + d.ChaincodeSeq
	exeCmd = exeCmd + peerAddr
	//exeCmd = exeCmd + " --signature-policy \"" + d.EndorsePolicy + "\""
	exeCmd = exeCmd + " --init-required --tls --cafile " + orderCaFile
	exeCmd = exeCmd + " > " + outFile + " 2>&1"
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + exeCmd + ")" + ";"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	//logger.Debug("commit args=",updatePara.Args)
	_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
	if err != nil {
		msg := "orgCommitCCV2 commit update cli failed"
		logger.Error(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	//check the commit result start
	time.Sleep(5*time.Second)
	parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
	flag := "(VALID)"
	_,err = getResultByFlag(parseFile,flag) 
	if err != nil {
		msg := "orgCommitCCV2 get cc commit status failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
		return c.JSON(http.StatusOK,ret)
	}
	//check the commit result end
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func orgInitCCV2(c echo.Context,cfg public.DeployPara,d FabricDeployCCPara) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "orgInitCCV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "orgInitCCV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_DEPLOYCC
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "orgInitCCV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain)  

	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "orgInitCCV2 get tools image failed"
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
	argStr,err := generateArgJsonStr(d.InitArgs) 
	if err != nil {
		msg := "orgInitCCV2 generateArgJsonStr failed"
		logger.Errorf(msg)
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	argStr = "'" + argStr + "'"
	args := []string{"sh","-c"}
	var updatePara deployfabric.PathToolsDeploymentPara
	updatePara.ToolsImage = toolsImage
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		if d.OrgName == org.Name {
			updatePara.OrgName = org.Name
			updatePara.PeerDomain = org.Domain
			for _,spec :=  range org.Specs {
				updatePara.PeerName = spec.Hostname
				randomStr := utils.GenerateRandomString(8)
				outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.PeerName + "-init.txt"
				outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
				exeCmd := "peer chaincode invoke -o " + orderAddr
				exeCmd = exeCmd + " --tls --cafile " + orderCaFile
				exeCmd = exeCmd + " -C " + d.ChannelId
				exeCmd = exeCmd + " -n " + d.ChainCodeId
				exeCmd = exeCmd + " --isInit -c " + argStr
				exeCmd = exeCmd + " > " + outFile + " 2>&1"
				cmd := "cp -a /var/data/. /cert;"
				cmd = cmd + " echo " + randomStr + ";"
				cmd = cmd + " $(" + exeCmd + ")" + ";" 
				cmd = cmd + " /bin/bash"
				updatePara.Args = nil
				updatePara.Args = append(updatePara.Args,args...)
				updatePara.Args = append(updatePara.Args,cmd) 
				//logger.Debug("init args=",updatePara.Args)
				_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
				if err != nil {
					msg := "orgInitCCV2 init update cli failed"
					logger.Error(msg)
					ret := getApiRet(CODE_ERROR_EXE,msg,nil)
						return c.JSON(http.StatusOK,ret)
				}
				//check the init result start
				time.Sleep(5*time.Second)
				parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
				flag := "status:200"
				_,err = getResultByFlag(parseFile,flag) 
				if err != nil {
					msg := "orgInitCCV2 get cc init status failed"
					logger.Errorf(msg)
					ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
					return c.JSON(http.StatusOK,ret)
				}
				//check the init result end
				break
			}	
			ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
			return c.JSON(http.StatusOK,ret) 
		}
	}
	msg := "orgName not exsist"
	ret := getApiRet(CODE_ERROR_EXE,msg,nil)
	return c.JSON(http.StatusOK,ret) 
}

func orgUpgradeCCV2(c echo.Context,cfg public.DeployPara,d FabricDeployCCPara) error { 
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chain,_ := k8s.GetChain(d.BlockChainId,cfg.ClusterId)
	if chain == nil {
		msg := "orgUpgradeCCV2:not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	// locker := getChainOpLocker(d.BlockChainId)
	// locker.Lock()
	// defer locker.Unlock()

	if chain.Status >k8s.CHAIN_STATUS_FREE {
		msg := "orgUpgradeCCV2: chain status not support operation"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	chain.Status = k8s.CHAIN_STATUS_DEPLOYCC
	err := k8s.UpdateChainStatus(*chain)
	if err != nil {
		msg := "orgUpgradeCCV2: set chain status error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	} 
	chain.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*chain)  

	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "orgUpgradeCCV2 get tools image failed"
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
	args := []string{"sh","-c"}
	ccPkgName := d.ChainCodeId + d.ChainCodeVersion + ".tar.gz"
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		var updatePara deployfabric.PathToolsDeploymentPara
		curOrg := org.Name
		updatePara.ToolsImage = toolsImage
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + curOrg + "-install.txt"
		outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		exeCmd := "peer lifecycle chaincode install " + ccPkgName + " > " + outFile + " 2>&1" 
		cmd := "cp -a /var/data/. /cert; cp /var/data/src/" + ccPkgName + " . ;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash" 
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgUpgradeCCV2 cc install update cli failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		time.Sleep(5*time.Second)
		parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName 
		flag := d.ChainCodeId + d.ChainCodeVersion + ":"
		pkgIdent,err := getResultByIdentifier(parseFile,flag) 
		logger.Debug("cc identifier=",pkgIdent)
		if err != nil {
			msg := "orgUpgradeCCV2 get cc identifier failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
		outFileName = d.ChainCodeId + d.ChainCodeVersion + "-" + curOrg + "-approve.txt"
		outFile = "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		pkgIdent = flag + pkgIdent
		exeCmd = "peer lifecycle chaincode approveformyorg -o " + orderAddr
		exeCmd = exeCmd + " --tls --cafile " + orderCaFile
		exeCmd = exeCmd + " --channelID " + d.ChannelId
		exeCmd = exeCmd + " --name " + d.ChainCodeId
		exeCmd = exeCmd + " --version " + d.ChainCodeVersion 
		exeCmd = exeCmd + " --init-required "
		exeCmd = exeCmd + " --package-id " + pkgIdent
		exeCmd = exeCmd + " --sequence " + d.ChaincodeSeq
		exeCmd = exeCmd + " > " + outFile + " 2>&1"
		cmd = "cp -a /var/data/. /cert;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = nil
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		//logger.Debug("approve args=",updatePara.Args)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgUpgradeCCV2 cc approve update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		time.Sleep(5*time.Second)
		parseFile = utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		flag = "(VALID)"
		_,err = getResultByFlag(parseFile,flag) 
		if err != nil {
			msg := "orgUpgradeCCV2 get cc approve status failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
	}
	time.Sleep(5*time.Second)
	var peerAddr string
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		for _,spec :=  range org.Specs {
			tmpPeerAddr :=  " --peerAddresses  " + spec.Hostname + ":7051 --tlsRootCertFiles /cert/crypto-config/peerOrganizations/"
			tmpPeerAddr = tmpPeerAddr + org.Domain + "/peers/"
			tmpPeerAddr = tmpPeerAddr + spec.Hostname + "." + org.Domain + "/tls/ca.crt"
			peerAddr = peerAddr + tmpPeerAddr
			break
		}
	}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		var updatePara deployfabric.PathToolsDeploymentPara
		updatePara.ToolsImage = toolsImage
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		outFileName := d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.OrgName + "-commit.txt"
		outFile := "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		exeCmd := "peer lifecycle chaincode commit -o " + orderAddr
		exeCmd = exeCmd + " --channelID " + d.ChannelId
		exeCmd = exeCmd + " --name " + d.ChainCodeId
		exeCmd = exeCmd + " --version " + d.ChainCodeVersion
		exeCmd = exeCmd + " --sequence " + d.ChaincodeSeq
		exeCmd = exeCmd + peerAddr
		exeCmd = exeCmd + " --init-required --tls --cafile " + orderCaFile
		exeCmd = exeCmd + " > " + outFile + " 2>&1"
		cmd := "cp -a /var/data/. /cert;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd)
		//logger.Debug("commit args=",updatePara.Args)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgUpgradeCCV2 commit update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		//check the commit result start
		time.Sleep(5*time.Second)
		parseFile := utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		flag := "(VALID)"
		_,err = getResultByFlag(parseFile,flag) 
		if err != nil {
			msg := "orgUpgradeCCV2 get cc commit status failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
		//check the commit result end

		argStr,err := generateArgJsonStr(d.InitArgs)
		if err != nil {
			msg := "orgUpgradeCCV2 generateArgJsonStr failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		outFileName = d.ChainCodeId + d.ChainCodeVersion + "-" + updatePara.OrgName + "-init.txt"
		outFile = "/var/data/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		argStr = "'" + argStr + "'"
		exeCmd = "peer chaincode invoke -o " + orderAddr
		exeCmd = exeCmd + " --tls --cafile " + orderCaFile
		exeCmd = exeCmd + " -C " + d.ChannelId
		exeCmd = exeCmd + " -n " + d.ChainCodeId
		exeCmd = exeCmd + " --isInit -c " + argStr
		exeCmd = exeCmd + " > " + outFile + " 2>&1"
		cmd = "cp -a /var/data/. /cert;"
		cmd = cmd + " $(" + exeCmd + ")" + ";"
		cmd = cmd + " /bin/bash"
		updatePara.Args = nil
		updatePara.Args = append(updatePara.Args,args...)
		updatePara.Args = append(updatePara.Args,cmd) 
		//logger.Debug("init args=",updatePara.Args)
		_,err = deployfabric.PatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,chain.BlockChainName,d.BlockChainId,updatePara)
		if err != nil {
			msg := "orgUpgradeCCV2 init update cli failed"
			logger.Error(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		//check the init result start
		time.Sleep(5*time.Second)
		parseFile = utils.BAAS_CFG.NfsLocalRootDir + d.BlockChainId + "/src/" + d.ChainCodeId + d.ChainCodeVersion + "/" + outFileName
		flag = "status:200"
		_,err = getResultByFlag(parseFile,flag) 
		if err != nil {
			msg := "orgUpgradeCCV2 get cc init status failed"
			logger.Errorf(msg)
			ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
			return c.JSON(http.StatusOK,ret)
		}
		//check the init result end
		break
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil) 
	return c.JSON(http.StatusOK,ret)
}