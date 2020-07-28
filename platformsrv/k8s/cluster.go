
package k8s

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
)

type Cluster struct {
	AllianceId	 string `json:"AllianceId,omitempty"` 
	Creator      string `json:"Creator,omitempty"`  
	ClusterId    string `json:"ClusterId"`
	ApiUrl       string `json:"ApiUrl"`
	HostDomain   string `json:"HostDomain"`
	PublicIp	 string `json:"PublicIp"` 
	InterIp      string `json:"InterIp"`  
	Cert         string `json:"Cert"`
	Key          string `json:"Key"`
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
	}
	logger.Errorf("GetClusters: cluster config not exsist")
	return nil,fmt.Errorf("GetClusters: cluster config not exsist")
}

func GetClustersByUser(user string)([]Cluster,error) {
	cfgPath := utils.BAAS_CFG.ClusterCfgPath  + CLUSTER_CFG_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var clusters []Cluster
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&clusters)
			if err != nil {
				logger.Errorf("GetClustersByUser: unmarshal clusters error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			var cas []Cluster
			for _,c := range clusters {
				tmpCluster := c
				if tmpCluster.Creator == user {
					cas = append(cas,tmpCluster)
				}
			}
			return cas,nil
		}else {
			logger.Errorf("GetClustersByUser: load cluster config error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Errorf("GetClustersByUser: cluster config not exsist")
	return nil,fmt.Errorf("GetClustersByUser: cluster config not exsist")
}

func GetCluster(clusterId string) (*Cluster,error) {
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

func DeleteCluster(clusterId string,creator string)error {
	cfgPath := utils.BAAS_CFG.ClusterCfgPath  + CLUSTER_CFG_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var clusters []Cluster
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&clusters)
			if err != nil {
				logger.Errorf("DeleteCluster: unmarshal clusters error,%v", err)
				return fmt.Errorf("%v", err)
			}
			var cas []Cluster
			for _,c := range clusters {
				var tmpCluster Cluster
				tmpCluster = c
				if tmpCluster.ClusterId == clusterId {
					if tmpCluster.Creator != creator {
						logger.Errorf("DeleteCluster: user is not cluster creator")
						return fmt.Errorf("DeleteCluster: user is not cluster creator")
					}
					chs,_ := GetChains(clusterId)
					if chs != nil {
						if len(chs) > 0{
							logger.Errorf("DeleteCluster: this cluster has blockchain,please delete blockchain first")
							return fmt.Errorf("DeleteCluster: this cluster has blockchain,please delete blockchain first")
						}
					}
					key := utils.BAAS_CFG.ClusterPkiBasePath + tmpCluster.Key
					cert := utils.BAAS_CFG.ClusterPkiBasePath + tmpCluster.Cert
					os.RemoveAll(key)
					os.RemoveAll(cert) 
				}else{
					cas = append(cas,tmpCluster)
				}
			}
			bytes, err := json.Marshal(cas)
			if err != nil {
				logger.Errorf("DeleteCluster: marshal clusters error,%v", err)
				return fmt.Errorf("%v", err)
			}
			err = utils.WriteFile(cfgPath,string(bytes))
			if err != nil {
				logger.Errorf("DeleteCluster: Write cluster config file error,%v", err)
				return fmt.Errorf("%v", err)
			}
			return nil
		}else {
			logger.Errorf("DeleteCluster: load cluster config error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	logger.Errorf("DeleteCluster: cluster config not exsist")
	return fmt.Errorf("DeleteCluster: cluster config not exsist")
}
