
package api

import (
	"os"
	"io"
	"strings"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/fabquery"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
)

func orgCreateChannel(c echo.Context) error { 
	logger.Debug("orgCreateChannel")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.OrgCreateChannel(d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func orgJoinChannel(c echo.Context) error {
	logger.Debug("orgJoinChannel")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	// err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	// if err != nil {
	// 	ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
	// 	return c.JSON(http.StatusOK,ret)
	// }
	err = fabric.OrgJoinChannel(d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func upChainCode(c echo.Context) error {  
	logger.Debug("upChainCode")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	chainId := c.FormValue("BlockChainId")
	err := k8s.CheckChainStatus(chainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ccId := c.FormValue("ChainCodeId")
	ccVersion := c.FormValue("ChainCodeVersion")
	cfg,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		msg := "upChainCode: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"2.") {
		return upChainCodeV2(c,cfg)
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
		msg := "copy cc file error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func singleOrgDeployCC(c echo.Context) error { 
	logger.Debug("singleOrgDeployCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricDeployCCPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "singleOrgDeployCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"2.") {
		return singleOrgInstallCCV2(c,cfg,d) 
	}
	err = fabric.SingleOrgInstallChaiCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.ChainCodeVersion)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func orgInstantialCC(c echo.Context) error { 
	logger.Debug("orgInstantialCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricDeployCCPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "orgInstantialCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"2.") {
		// msg := "orgInstantialCC: favric v2 not support this api"
		// ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		// return c.JSON(http.StatusOK,ret)
		return orgInitCCV2(c,cfg,d)
	}
	err = fabric.OrgInstantialChaiCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.ChainCodeVersion,d.EndorsePolicy,d.InitArgs)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil) 
	return c.JSON(http.StatusOK,ret) 
}

func orgDeployCC(c echo.Context) error { 
	logger.Debug("orgDeployCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricDeployCCPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "orgDeployCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"2.") {
		return orgDeployCCV2(c,cfg,d)
	}
	// err = fabric.OrgDeployChaiCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.ChainCodeVersion,d.EndorsePolicy,d.InitArgs)
	// if err != nil {
    //     ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
	// 	return c.JSON(http.StatusOK,ret)
	// } 
	
	go fabric.OrgDeployChaiCodeRoutine(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.ChainCodeVersion,d.EndorsePolicy,d.InitArgs) 

	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func orgUpgradeCC(c echo.Context) error { 
	logger.Debug("orgUpgradeCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricDeployCCPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "orgUpgradeCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"2.") {
		return orgUpgradeCCV2(c,cfg,d)
	}
	// err = fabric.OrgUpgradeChaiCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.ChainCodeVersion,d.EndorsePolicy,d.InitArgs)
	// if err != nil {
    //     ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil) 
	// 	return c.JSON(http.StatusOK,ret)
	// }
	
	go fabric.OrgUpgradeChaiCodeRoutine(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.ChainCodeVersion,d.EndorsePolicy,d.InitArgs)

	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil) 
	return c.JSON(http.StatusOK,ret) 
}

func chaincodeCall(c echo.Context) error {
	logger.Debug("chaincodeCall")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricCCCallPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "fabric chaincodeCall: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"2.") {
		return callCCV2(c,cfg,d)
	} 
	cr,err := fabric.OrgInvokeChainCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.Args,d.Peers)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,cr)
	return c.JSON(http.StatusOK,ret)
}

func chaincodeQuery(c echo.Context) error {
	logger.Debug("chaincodeQuery")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricCCQueryPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	// cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	// if err != nil {
	// 	msg := "fabric chaincodeQuery: not find this chain"
	// 	ret := getApiRet(CODE_ERROR_EXE,msg,nil)
	// 	return c.JSON(http.StatusOK,ret)
	// }
	// if strings.HasPrefix(cfg.Version,"2.") {
	// 	return ccQueryV2(c,cfg,d)
	// }
	qr,err := fabric.OrgQueryChainCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.Args)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret)
}
/*
func queryInstatialCC(c echo.Context) error {
	logger.Debug("queryInstatialCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricInstantialCCPara 
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "fabric queryInstatialCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"2.") {
		return ccQueryInstatialV2(c,cfg,d) 
	}
	qr,err := fabric.OrgQueryInstantiateCC(d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret)
}
*/

func queryInstatialCC(c echo.Context) error {
	logger.Debug("queryInstatialCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricInstantialCCPara 
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabric.GetCCRecord(d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret)
}

func queryInstalledCC(c echo.Context) error {
	logger.Debug("queryInstalledCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabric.OrgQueryInstalledCC(d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret) 
}

func queryChannel(c echo.Context) error {
	logger.Debug("queryChannel")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabric.OrgQueryChannel(d.BlockChainId,d.OrgName)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret) 
}

func queryTxInfo(c echo.Context) error {
	logger.Debug("queryTxInfo")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`
		TxId string `json:"TxId"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabric.OrgQueryTxById(d.BlockChainId,d.OrgName,d.ChannelId,d.TxId)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret) 
}

func queryBlockInfo(c echo.Context) error {
	logger.Debug("queryBlockInfo")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`
		BlockId uint64 `json:"BlockId"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabric.OrgQueryBlockById(d.BlockChainId,d.OrgName,d.ChannelId,d.BlockId)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil) 
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret) 
}

func queryBlock(c echo.Context) error {
	logger.Debug("queryBlock")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"` 
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabric.OrgQueryBlockChain(d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret) 
}

func addOrg(c echo.Context) error {
	logger.Debug("addOrg")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d public.AddOrgConfig
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChainAddOrg(d.BlockChainId,d)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,"") 
	return c.JSON(http.StatusOK,ret) 
} 

func orgApproveCC(c echo.Context) error { 
	logger.Debug("orgApproveCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricDeployCCPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "orgApproveCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"1.") {
		msg := "orgApproveCC: this api only for v2"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	return orgApproveCCV2(c,cfg,d) 
}

func orgCommitCC(c echo.Context) error { 
	logger.Debug("orgCommitCC")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d FabricDeployCCPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	cfg,err := sdkfabric.LoadChainCfg(d.BlockChainId)
	if err != nil {
		msg := "orgCommitCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if strings.HasPrefix(cfg.Version,"1.") {
		msg := "orgCommitCC: this api only for v2"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	return orgCommitCCV2(c,cfg,d)
}

func queryBlockTx(c echo.Context) error {
	logger.Debug("queryBlockTx")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"` 
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	logger.Debug("queryBlockTx: para=")
	logger.Debug(d)
	if d.BlockChainId == "" || d.ChannelId == "" {
		msg := "parameter err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabquery.GetBlockTx(d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil {
		msg := "orgCommitCC: not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr)  
	return c.JSON(http.StatusOK,ret) 
}

func queryBlockAndTx(c echo.Context) error {
	logger.Debug("queryBlockAndTx")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		Start uint64 `json:"Start"`
		End uint64 `json:"End"`
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"` 
	}
	var d ReqPara 
    err = json.Unmarshal(result, &d) 
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	if d.BlockChainId == "" || d.ChannelId == "" || d.OrgName == "" || d.Start >= d.End {
		msg := "parameter err"
        ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabquery.GetBlockAndTxList(d.Start,d.End,d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret)
}

func queryTxFromDb(c echo.Context) error {
	logger.Debug("queryTxFromDb")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"` 
		TxId string `json:"TxId"`
	}
	var d ReqPara 
    err = json.Unmarshal(result, &d) 
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	if d.BlockChainId == "" || d.ChannelId == "" || d.OrgName == "" || d.TxId == "" {
		msg := "parameter err"
        ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabquery.GetTxInfo(d.BlockChainId,d.TxId)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret)
}

func searchFromDb(c echo.Context) error {
	logger.Debug("searchFromDb")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"` 
		Key string `json:"Key"`
	}
	var d ReqPara 
    err = json.Unmarshal(result, &d) 
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	if d.BlockChainId == "" || d.ChannelId == "" || d.OrgName == "" || d.Key == "" {
		msg := "parameter err"
        ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabquery.GetSearchResult(d.BlockChainId,d.Key)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret)
}

func queryPassTx(c echo.Context) error {
	logger.Debug("queryPassTx")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`  
		Day int `json:"Day"`
	}
	var d ReqPara 
    err = json.Unmarshal(result, &d) 
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	if d.BlockChainId == "" || d.ChannelId == "" || d.OrgName == "" || d.Day >= 0 || d.Day < -14 {
		msg := "parameter err"
        ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabquery.GetPassDayTxCount(d.BlockChainId,d.OrgName,d.ChannelId,d.Day)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr) 
	return c.JSON(http.StatusOK,ret)
}

func queryDayTx(c echo.Context) error {
	logger.Debug("queryDayTx")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`  
		StartTime string `json:"StartTime"`
		EndTime string `json:"EndTime"`
	}
	var d ReqPara 
    err = json.Unmarshal(result, &d) 
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	if d.BlockChainId == "" || d.ChannelId == "" || d.OrgName == "" || d.StartTime == "" || d.EndTime == ""{
		msg := "parameter err"
        ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabquery.GetDayTxCount(d.StartTime,d.EndTime,d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr)  
	return c.JSON(http.StatusOK,ret)
}

func queryEveryDayTx(c echo.Context) error {
	logger.Debug("queryEveryDayTx")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	type ReqPara struct {
		BlockChainId string `json:"BlockChainId"`
		OrgName string `json:"OrgName"`
		ChannelId string `json:"ChannelId"`  
		StartTime string `json:"StartTime"`
		Days int `json:"Days"`
	}
	var d ReqPara 
    err = json.Unmarshal(result, &d) 
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	if d.BlockChainId == "" || d.ChannelId == "" || d.OrgName == "" || d.StartTime == "" || d.Days <= 0 || d.Days > 14{
		msg := "parameter err"
        ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = k8s.CheckChainStatus(d.BlockChainId)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = fabric.ChannelCheck(d.BlockChainId,d.ChannelId,d.OrgName)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	qr,err := fabquery.GetEveryDayTxCount(d.StartTime,d.Days,d.BlockChainId,d.OrgName,d.ChannelId)
	if err != nil { 
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,qr)  
	return c.JSON(http.StatusOK,ret)
}