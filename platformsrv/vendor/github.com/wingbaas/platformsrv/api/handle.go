
package api

import (
	"os"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/k8s/deployfabric"
	"github.com/goinggo/mapstructure"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
)

const (
	CODE_SUCCESS             int = 0
	CODE_ERROR_BODY          int = 101
	CODE_ERROR_MASHAL        int = 102
	CODE_ERROR_EXE           int = 103
	CODE_ERROR_CONFIG        int = 103
	MSG_SUCCESS              string = "success"
	MSG_ERROR                string = "error"
)

type Deploy struct {
	BlockChainName      string        `json:"BlockChainName"`
	BlockChainType      string        `json:"BlockChainType"`
	DeployCfg     		interface{}   `json:"DeployCfg"`       
}

type Delete struct{
	ClusterId 		string  `json:"ClusterId"`
	BlockChainName	string  `json:"BlockChainName"`
}

/* All interface of platformSrv return struct */
type ApiRet struct {
	Code    int         `json:"code"`    //error code, 0 if success
	Message string      `json:"message"` //description message
	Data    interface{} `json:"data"`    //return data in this member
}

func getApiRet (code int, msg string, data interface{}) ApiRet {
	var ret ApiRet
	ret.Code    = code
	ret.Message = msg
	ret.Data    = data
	return ret
}

func addCluster(c echo.Context) error {
	logger.Debug("addCluster")
	var cluster k8s.Cluster
	result, err := ioutil.ReadAll(c.Request().Body) 
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
    err = json.Unmarshal(result, &cluster)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	err = k8s.AddCluster(cluster)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret) 
}

func getClusters(c echo.Context) error {
	logger.Debug("getClusters")
	obj,err := k8s.GetClusters()
	if err!= nil {
		msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,obj)
	return c.JSON(http.StatusOK,ret)
}

func getHosts(c echo.Context) error {
	logger.Debug("getHosts")
	clusterId := c.Param("clusterid")
	obj,err := k8s.GetHostList(clusterId)
	if err!= nil {
		msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,obj)
	return c.JSON(http.StatusOK,ret)
}

func getNamespaces(c echo.Context) error {
	logger.Debug("getNamespaces")
	clusterId := c.Param("clusterid")
	obj,_ := k8s.GetClusterNamespaces(clusterId)
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,obj)
	return c.JSON(http.StatusOK,ret)
}  

func deployBlockChain(c echo.Context) error { 
	logger.Debug("deployBlockChain")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d Deploy
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var blockId string
	var clusterId string
	if d.BlockChainType == "fabric" {
		cfgMap, ok := d.DeployCfg.(map[string]interface{}) 
		if !ok {
			msg := "blockchain fabric deploy parameter error"
        	ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		var cfg public.DeployPara
		err = mapstructure.Decode(cfgMap,&cfg); 
		if err != nil {
			msg := "blockchain fabric decode config map error"
        	ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		cn,_ := k8s.GetChainByName(d.BlockChainName,cfg.ClusterId)
		if cn != nil {
			msg := "blockchainName already exsist: " + d.BlockChainName
        	ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
		}
		blockId,err = fabric.DeployFabric(cfg,d.BlockChainName,d.BlockChainType)
		clusterId = cfg.ClusterId
		if err != nil {
			ret := getApiRet(CODE_ERROR_BODY,err.Error(),nil)
			return c.JSON(http.StatusOK,ret)
		}
	}else if d.BlockChainType == "wingchain" {
		msg := "wingchain unsupported temporary!"
        ret := getApiRet(CODE_SUCCESS,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}else {
		msg := "unsupported blockchain type!"
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var chain k8s.Chain
	chain.BlockChainId = blockId
	chain.BlockChainName = d.BlockChainName
	chain.BlockChainType = d.BlockChainType
	chain.ClusterId = clusterId
	err = k8s.AddChain(chain)
	if err != nil {
		ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,chain)
	return c.JSON(http.StatusOK,ret) 
}

func getChains(c echo.Context) error {
	logger.Debug("getChains")
	clusterId := c.Param("clusterid")
	obj,err := k8s.GetChains(clusterId)
	if err!= nil {
		msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,obj)
	return c.JSON(http.StatusOK,ret)
}

func getBlockChainTypes(c echo.Context) error {
	logger.Debug("getBlockChainTypes")
	if utils.BLOCK_CFG_MAP == nil {
		msg := "blockchain type config not find"
        ret := getApiRet(CODE_ERROR_CONFIG,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,utils.BLOCK_CFG_MAP)
	return c.JSON(http.StatusOK,ret)
}

func deleteBlockChain(c echo.Context) error { 
	logger.Debug("deleteBlockChain")
	result, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var d Delete
    err = json.Unmarshal(result, &d)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	_,err = deployfabric.DeleteNamespace(d.ClusterId,d.BlockChainName)
	if err != nil {
        ret := getApiRet(CODE_ERROR_EXE,err.Error(),nil)
		return c.JSON(http.StatusOK,ret)
	}
	ch,_ := k8s.GetChainByName(d.BlockChainName,d.ClusterId) 
	if ch != nil {
		k8s.DeleteChain(*ch)
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + ch.BlockChainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + ch.BlockChainId
		cfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + ch.BlockChainId + ".json"
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		os.Remove(cfgFile)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret)
}
