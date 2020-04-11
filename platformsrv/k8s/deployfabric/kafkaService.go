
package deployfabric

import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type KafkaService struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataService `json:"metadata"`
	Spec SpecService `json:"spec"`
}

func CreateKafkaService(clusterId string,namespaceId string,chainId string,kafkaName string)([]byte,error) {
	kafkaService := ZookeeperService{
		APIVersion: "v1",
		Kind: "Service",   
		Metadata: MetadataService{
			Name: kafkaName,
		},
		Spec: SpecService{
			Ports: []PortSpecService{
				{
					Name: kafkaName,
					Port: 9092,
					TargetPort: 9092,
				},
			},
			Selector: SelectorSpecService{
				App: kafkaName,
			},
		},
	}
	bytes, err := CreateService(clusterId,namespaceId,kafkaService)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


