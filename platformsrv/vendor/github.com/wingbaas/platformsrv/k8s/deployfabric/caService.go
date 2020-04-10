
package deployfabric

import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type CaService struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataService `json:"metadata"`
	Spec SpecService `json:"spec"`
}

func CreateCaService(clusterId string,namespaceId string,chainId string,caName string)([]byte,error) {
	caService := CaService{
		APIVersion: "v1",
		Kind: "Service",
		Metadata: MetadataService{
			Name: caName,
		},
		Spec: SpecService{
			Type: "NodePort",
			Ports: []PortSpecService{
				{
					Name: caName,
					Port: 7054,
					TargetPort: 7054,
				},
			},
			Selector: SelectorSpecService{
				App: caName,
			},
		},
	}
	bytes, err := CreateService(clusterId,namespaceId,caService)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


