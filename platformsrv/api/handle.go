
package api

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
)

const (
	CODE_SUCCESS             int = 0
	CODE_ERROR_BODY          int = 101
	CODE_ERROR_MASHAL        int = 102
	CODE_ERROR_EXE           int = 103
	MSG_SUCCESS              string = "success"
	MSG_ERROR                string = "error"
)

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

func getNamespaces(c echo.Context) error {
	logger.Debug("getNamespaces")
	clusterId := c.Param("clusterid")
	obj,_ := k8s.GetClusterNamespaces(clusterId)
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,obj)
	return c.JSON(http.StatusOK,ret)
}  
