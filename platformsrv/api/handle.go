
package api

import (
	"os"
	"io"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/k8s/deployfabric"
	"github.com/goinggo/mapstructure"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
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
	AllianceId	        string        `json:"AllianceId,omitempty"`
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

func upLoadKeyFile(c echo.Context) error { 
	logger.Debug("upLoadKeyFile")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	file, err := c.FormFile("file") 
	if err != nil {
		msg := "get upLoadKeyFile error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	src, err := file.Open()
	if err != nil {
		msg := "open upLoadKeyFile error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	defer src.Close() 
	dstDir := utils.BAAS_CFG.ClusterPkiBasePath
	fileId := utils.GenerateRandomString(16)
	dst, err := os.Create(dstDir + fileId)  
	if err != nil {
		msg := "create upLoadKeyFile error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	defer dst.Close()
	_,err = io.Copy(dst, src)
	if err != nil {
		msg := "copy key file error"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,fileId) 
	return c.JSON(http.StatusOK,ret) 
} 


func addCluster(c echo.Context) error {
	logger.Debug("addCluster")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
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
	if cluster.ClusterId == "" || cluster.ApiUrl == "" || cluster.HostDomain == "" || cluster.PublicIp == "" || cluster.Cert == "" || cluster.Key == "" {
		msg := "parameter error"
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	bl := sdkfabric.AddHosts(cluster.HostDomain,cluster.PublicIp)
	if !bl {
		msg := "add cluster domain to /etc/hosts failed"
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	clusterId := c.Param("clusterid")
	obj,_ := k8s.GetClusterNamespaces(clusterId)
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,obj)
	return c.JSON(http.StatusOK,ret)
}  

func deployBlockChain(c echo.Context) error { 
	logger.Debug("deployBlockChain")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
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
	if d.BlockChainName == "" {
		msg := "blockchainname error"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	var blockId string
	var clusterId string
	var version string
	if d.BlockChainType == public.BLOCK_CHAIN_TYPE_FABRIC {
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
		version = cfg.Version
		if strings.HasPrefix(cfg.Version,"1.") {
			if cfg.DeployType != public.KAFKA_FABRIC {
				msg := "version 1.x only support deploy type KAFKA_FABRIC"
        		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
				return c.JSON(http.StatusOK,ret)
			}
		}else if strings.HasPrefix(cfg.Version,"2."){
			if cfg.DeployType != public.RAFT_FABRIC {
				msg := "version 2.x only support deploy type RAFT_FABRIC"
        		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
				return c.JSON(http.StatusOK,ret)
			}
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
	chain.AllianceId = d.AllianceId
	chain.BlockChainId = blockId
	chain.BlockChainName = d.BlockChainName
	chain.BlockChainType = d.BlockChainType
	chain.ClusterId = clusterId
	chain.Version = version
	chain.Status = k8s.CHAIN_STATUS_FREE
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
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
		time.Sleep(60*time.Second)
		k8s.DeleteChain(*ch)
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + ch.BlockChainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + ch.BlockChainId
		cfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + ch.BlockChainId + ".json"
		keyPath := utils.BAAS_CFG.KeyStorePath + ch.BlockChainId 
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		os.Remove(cfgFile) 
		os.Remove(keyPath)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil) 
	return c.JSON(http.StatusOK,ret)
}

func getBlockchainCfg(c echo.Context) error {
	logger.Debug("getBlockchainCfg")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	blockchainId := c.Param("blockchainid")
	cfg,err := sdkfabric.LoadChainCfg(blockchainId)
	if err != nil {
		msg := "not find this chain"
		ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,cfg)
	return c.JSON(http.StatusOK,ret) 
} 
