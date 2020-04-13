
package deployfabric

import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type OrderService struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataService `json:"metadata"`
	Spec SpecService `json:"spec"`
}

func CreateOrderService(clusterId string,namespaceId string,chainId string,orderName string)([]byte,error) {
	orderService := OrderService{
		APIVersion: "v1",
		Kind: "Service",   
		Metadata: MetadataService{
			Name: orderName,
		},
		Spec: SpecService{
			Type: "NodePort", 
			Ports: []PortSpecService{
				{
					Name: orderName,
					Port: 7050,
					TargetPort: 7050,
				},
			},
			Selector: SelectorSpecService{
				App: orderName,
			},
		},
	}
	bytes, err := CreateService(clusterId,namespaceId,orderService)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


