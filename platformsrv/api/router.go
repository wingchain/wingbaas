
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
	r.POST("/api/v1/uploadkeyfile",upLoadKeyFile) 
	r.POST("/api/v1/addcluster",addCluster) 
	r.POST("/api/v1/deletecluster",deleteCluster) 
	r.GET("/api/v1/clusters",getClusters)
	r.GET("/api/v1/:usermail/clusters",getUserClusters) 
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
	r.POST("/api/v1/queryinstatialcc",queryInstatialCC) 
	//r.POST("/api/v1/queryinstalledcc",queryInstalledCC) 
	r.POST("/api/v1/querychannel",queryChannel)
	r.POST("/api/v1/querytxinfo",queryTxInfo)
	r.POST("/api/v1/queryblockinfo",queryBlockInfo)
	r.POST("/api/v1/queryblock",queryBlock)   
	r.POST("/api/v1/queryblocktx",queryBlockAndTx) 
	r.POST("/api/v1/querytx",queryTxFromDb) 
	r.POST("/api/v1/search",searchFromDb)
	r.POST("/api/v1/querypasstx",queryPassTx) 
	r.POST("/api/v1/querydaytx",queryDayTx) 
	r.POST("/api/v1/queryeverydaytx",queryEveryDayTx)
	r.POST("/api/v1/addorg",addOrg)
	r.POST("/api/v1/singleorgdeploycc",singleOrgDeployCC) 
	r.POST("/api/v1/orginstantialcc",orgInstantialCC) 
	r.POST("/api/v1/orgupgradecc",orgUpgradeCC) 
	r.POST("/api/v1/orgapprovecc",orgApproveCC)
	r.POST("/api/v1/orgcommitcc",orgCommitCC)

	r.POST("/api/v1/register",userRegister) 
	r.POST("/api/v1/login",userLogin) 
	r.POST("/api/v1/createalliance",creatAlliance)  
	r.GET("/api/v1/alliances",getAlliance)
	r.GET("/api/v1/:allianceid/alliance",getAllianceById)
	r.POST("/api/v1/useraddalliance",userAddAlliance) 
	r.GET("/api/v1/:usermail/joinedalliances",getJoinedAlliance)
	r.GET("/api/v1/:usermail/createdalliances",getCreatedAlliance) 
	r.POST("/api/v1/deletealliance",deleteAlliance)  
	r.GET("/api/v1/:allianceid/allianceclusters",getAllianceClusters)  
	r.GET("/api/v1/:allianceid/alliancechains",getAllianceChains) 
	r.GET("/api/v1/:blockchainid/cfg",getBlockchainCfg)    
	r.POST("/api/v1/blocktx",queryBlockTx)  
	r.GET("/api/v1/:blockchainid/blockchain",getBlockChain) 
	r.GET("/api/v1/:allianceid/users",getUsersByAllianceId)
	r.POST("/api/v1/deleteallianceuser",deleteAllianceUser)  

	//options

	r.OPTIONS("/api/v1/uploadkeyfile",upLoadKeyFile) 
	r.OPTIONS("/api/v1/addcluster",addCluster) 
	r.OPTIONS("/api/v1/deletecluster",deleteCluster) 
	r.OPTIONS("/api/v1/clusters",getClusters) 
	r.OPTIONS("/api/v1/:usermail/clusters",getUserClusters) 
	r.OPTIONS("/api/v1/:clusterid/hosts",getHosts)
	r.OPTIONS("/api/v1/:clusterid/namespaces",getNamespaces)
	r.OPTIONS("/api/v1/blockchaintypes",getBlockChainTypes)
	r.OPTIONS("/api/v1/deploy",deployBlockChain)
	r.OPTIONS("/api/v1/:clusterid/blockchains",getChains)
	r.OPTIONS("/api/v1/delete",deleteBlockChain)
	r.OPTIONS("/api/v1/orgcreatechannel",orgCreateChannel)  
	r.OPTIONS("/api/v1/orgjoinchannel",orgJoinChannel) 
	r.OPTIONS("/api/v1/uploadcc",upChainCode)
	r.OPTIONS("/api/v1/orgdeploycc",orgDeployCC) 
	r.OPTIONS("/api/v1/callcc",chaincodeCall)
	r.OPTIONS("/api/v1/querycc",chaincodeQuery)
	r.OPTIONS("/api/v1/queryinstatialcc",queryInstatialCC)
	//r.OPTIONS("/api/v1/queryinstalledcc",queryInstalledCC) 
	r.OPTIONS("/api/v1/querychannel",queryChannel)
	r.OPTIONS("/api/v1/querytxinfo",queryTxInfo)
	r.OPTIONS("/api/v1/queryblockinfo",queryBlockInfo)
	r.OPTIONS("/api/v1/queryblock",queryBlock) 
	r.OPTIONS("/api/v1/queryblocktx",queryBlockAndTx) 
	r.OPTIONS("/api/v1/querytx",queryTxFromDb)
	r.OPTIONS("/api/v1/search",searchFromDb) 
	r.OPTIONS("/api/v1/querypasstx",queryPassTx) 
	r.OPTIONS("/api/v1/querydaytx",queryDayTx)
	r.OPTIONS("/api/v1/queryeverydaytx",queryEveryDayTx)
	r.OPTIONS("/api/v1/addorg",addOrg)
	r.OPTIONS("/api/v1/singleorgdeploycc",singleOrgDeployCC) 
	r.OPTIONS("/api/v1/orginstantialcc",orgInstantialCC) 
	r.OPTIONS("/api/v1/orgupgradecc",orgUpgradeCC) 
	r.OPTIONS("/api/v1/orgapprovecc",orgApproveCC)
	r.OPTIONS("/api/v1/orgcommitcc",orgCommitCC)


	r.OPTIONS("/api/v1/register",userRegister)
	r.OPTIONS("/api/v1/login",userLogin) 
	r.OPTIONS("/api/v1/createalliance",creatAlliance) 
	r.OPTIONS("/api/v1/alliances",getAlliance)
	r.OPTIONS("/api/v1/:allianceid/alliance",getAllianceById)
	r.OPTIONS("/api/v1/useraddalliance",userAddAlliance) 
	r.OPTIONS("/api/v1/:usermail/joinedalliances",getJoinedAlliance)
	r.OPTIONS("/api/v1/:usermail/createdalliances",getCreatedAlliance) 
	r.OPTIONS("/api/v1/deletealliance",deleteAlliance)
	r.OPTIONS("/api/v1/:allianceid/allianceclusters",getAllianceClusters) 
	r.OPTIONS("/api/v1/:allianceid/alliancechains",getAllianceChains) 
	r.OPTIONS("/api/v1/:blockchainid/cfg",getBlockchainCfg) 
	r.OPTIONS("/api/v1/blocktx",queryBlockTx) 
	r.OPTIONS("/api/v1/:blockchainid/blockchain",getBlockChain)
	r.OPTIONS("/api/v1/:allianceid/users",getUsersByAllianceId)
	r.OPTIONS("/api/v1/deleteallianceuser",deleteAllianceUser)

	// Start server
	logger.Debug("start wing baas api server") 
	r.Logger.Fatal(r.Start(":"+port))
	logger.Debug("stop wing baas api server")
}
