
package api

import (
	"os"
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
)

func orgCreateChannel(c echo.Context) error {
	logger.Debug("orgCreateChannel")
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
	d.ChannelId = sdkfabric.DefaultChannel
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
	d.ChannelId = sdkfabric.DefaultChannel
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
	chainId := c.FormValue("BlockChainId")
	ccId := c.FormValue("ChainCodeId")
	ccVersion := c.FormValue("ChainCodeVersion")
	file, err := c.FormFile("file") 
	if err != nil {
		// logger.Debug("para value=")
		// logger.Debug("%s %s %s\n",chainId,ccId,ccVersion)
		logger.Debug("err=%s",err)
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

func orgDeployCC(c echo.Context) error {
	logger.Debug("orgDeployCC")
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
		ChainCodeId string `json:"ChainCodeId"`
		ChainCodeVersion string `json:"ChainCodeVersion"`
		InitArgs []string `json:"InitArgs"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	d.ChannelId = sdkfabric.DefaultChannel
	err = fabric.OrgDeployChaiCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeId,d.ChainCodeVersion,d.InitArgs)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}

	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func chaincodeCall(c echo.Context) error {
	logger.Debug("chaincodeCall")
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
		ChainCodeID string `json:"ChainCodeID"`
		Args []string `json:"Args"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	d.ChannelId = sdkfabric.DefaultChannel
	err = fabric.OrgInvokeChainCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeID,d.Args)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}

func chaincodeQuery(c echo.Context) error {
	logger.Debug("chaincodeQuery")
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
		ChainCodeID string `json:"ChainCodeID"`
		Args []string `json:"Args"`
	}
	var d ReqPara
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	d.ChannelId = sdkfabric.DefaultChannel
	err = fabric.OrgQueryChainCode(d.BlockChainId,d.OrgName,d.ChannelId,d.ChainCodeID,d.Args)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil) 
	return c.JSON(http.StatusOK,ret)
}