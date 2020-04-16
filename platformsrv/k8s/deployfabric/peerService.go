
package deployfabric

import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type PeerService struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataService `json:"metadata"`
	Spec SpecService `json:"spec"`
}

func CreatePeerService(clusterId string,namespaceId string,chainId string,peerName string)([]byte,error) {
	peerService := PeerService{
		APIVersion: "v1",
		Kind: "Service",   
		Metadata: MetadataService{
			Name: peerName,
		},
		Spec: SpecService{
			Type: "NodePort", 
			Ports: []PortSpecService{
				{
					Name: "api",
					Port: 7051,
					TargetPort: 7051,
				},
				{
					Name: "cclisten",
					Port: 7052,
					TargetPort: 7052,   
				},
				{
					Name: "events",
					Port: 7053,
					TargetPort: 7053,
				},
			},
			Selector: SelectorSpecService{
				App: peerName,
			},
		},
	}
	bytes, err := CreateService(clusterId,namespaceId,peerService)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


