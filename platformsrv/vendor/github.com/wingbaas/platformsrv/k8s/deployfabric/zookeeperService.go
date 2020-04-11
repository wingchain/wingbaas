
package deployfabric

import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type ZookeeperService struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataService `json:"metadata"`
	Spec SpecService `json:"spec"`
}

func CreateZookeeperService(clusterId string,namespaceId string,chainId string,zkName string)([]byte,error) {
	zkService := ZookeeperService{
		APIVersion: "v1",
		Kind: "Service",
		Metadata: MetadataService{
			Name: zkName,
		},
		Spec: SpecService{
			//Type: "NodePort",    
			Ports: []PortSpecService{
				{
					Name: "client",
					Port: 2181,
					TargetPort: 2181,
				},
				{
					Name: "follower",
					Port: 2888,
					TargetPort: 2888,
				},
				{
					Name: "election",
					Port: 3888,
					TargetPort: 3888,
				},
			},
			Selector: SelectorSpecService{
				App: zkName,
			},
		},
	}
	bytes, err := CreateService(clusterId,namespaceId,zkService)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


