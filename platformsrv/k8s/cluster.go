
package k8s

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
)

type Cluster struct {
	ClusterId     string `json:"ClusterId"`
	Addr          string `json:"Addr"`
	PublicIp	  string `json:"PublicIp"`
	Cert          string `json:"Cert"`
	Key           string `json:"Key"`
}

const (
	CLUSTER_CFG_FILE  string = "cluster.json"
)

func AddCluster(cluster Cluster) error {
	cfgPath := utils.BAAS_CFG.ClusterCfgPath + CLUSTER_CFG_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var clusters []Cluster
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&clusters)
			if err != nil {
				logger.Errorf("AddCluster: unmarshal clusters error,%v", err)
				return fmt.Errorf("%v", err)
			}
		}else {
			logger.Errorf("AddCluster: load cluster config error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	for _,c := range clusters {
		if c.ClusterId == cluster.ClusterId {
			logger.Errorf("%s%s","AddCluster: cluster id already exsist ",cluster.ClusterId)
			return fmt.Errorf("%s%s","AddCluster: cluster id already exsist ",cluster.ClusterId)
		}
	}
	clusters = append(clusters,cluster)
	bytes, err := json.Marshal(clusters)
	if err != nil {
		logger.Errorf("AddCluster: marshal clusters error,%v", err)
		return fmt.Errorf("%v", err)
	}
	err = utils.WriteFile(cfgPath,string(bytes))
	if err != nil {
		logger.Errorf("AddCluster: Write cluster config file error,%v", err)
		return fmt.Errorf("%v", err)
	}
	return nil
}
 
func GetClusters()([]Cluster,error) {
	cfgPath := utils.BAAS_CFG.ClusterCfgPath  + CLUSTER_CFG_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var clusters []Cluster
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&clusters)
			if err != nil {
				logger.Errorf("GetClusters: unmarshal clusters error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			return clusters,nil
		}else {
			logger.Errorf("GetClusters: load cluster config error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}else {
		logger.Errorf("GetClusters: load cluster config error")
		return nil,fmt.Errorf("%s", "GetClusters:  not find config file")
	}
}

func GetCluster(clusterId string) (*Cluster,error) { 
	logger.Debug("GetCluster")
	clusters,err := GetClusters()
	if err != nil {
		logger.Errorf("GetCluster: get clusters error,%v", err)
		return nil,fmt.Errorf("%v", err)
	}
	for _,cluster := range clusters {
		if cluster.ClusterId == clusterId {
			return &cluster,nil
		}
	}
	return nil,nil
}
