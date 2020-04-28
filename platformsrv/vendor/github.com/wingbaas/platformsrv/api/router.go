
package api

import (
	"github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
	"github.com/wingbaas/platformsrv/logger"
)

func Start(port string) {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.POST("/api/v1/addcluster",addCluster) 
	r.GET("/api/v1/clusters",getClusters)
	r.GET("/api/v1/:clusterid/hosts",getHosts)
	r.GET("/api/v1/:clusterid/namespaces",getNamespaces)
	r.GET("/api/v1/blockchaintypes",getBlockChainTypes)
	r.POST("/api/v1/deploy",deployBlockChain)
	r.GET("/api/v1/:clusterid/blockchains",getChains)
	r.POST("/api/v1/delete",deleteBlockChain)
	r.POST("/api/v1/orgcreatechannel",orgCreateChannel)  
	r.POST("/api/v1/orgjoinchannel",orgJoinChannel) 
	r.POST("/api/v1/uploadcc",upChainCode)
	r.POST("/api/v1/orgdeploycc",orgDeployCC)
	r.POST("/api/v1/callcc",chaincodeCall)
	r.POST("/api/v1/querycc",chaincodeQuery) 

	// Start server
	logger.Debug("start wing baas api server")
	r.Logger.Fatal(r.Start(":"+port))
	logger.Debug("stop wing baas api server")
}
