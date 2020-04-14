
package deployfabric
import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type SpecTemplateStOrder struct {
	//NodeSelector NodeSelectorSpecTemplateSt `json:"nodeSelector,omitempty"`
	Containers []ContainerSpecTemplateSt `json:"containers"`
	RestartPolicy string `json:"restartPolicy"`
	ImagePullSecrets []ImagePullSecretSpecTemplateSt `json:"imagePullSecrets"`
	Hostname string `json:"hostname"`
	Volumes []VolumeSpecTemplateSt `json:"volumes"`
}

type TemplateStOrder struct {
	Metadata MetadataTemplateSt `json:"metadata"` 
	Spec SpecTemplateStOrder `json:"spec"`
} 
type SpecStOrder struct {
	Selector SelectorSt `json:"selector"`
	Replicas int `json:"replicas"`
	Strategy StrategySt `json:"strategy"`
	Template TemplateStOrder `json:"template"`
}

type OrderDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataDeployMent `json:"metadata"`
	Spec SpecStOrder `json:"spec"`  
}  

func CreateOrderKafkaDeployment(clusterId string,namespaceId string,chainId string,image string,orderName string,orderDomain string)([]byte,error) {
	orderKafkaDeployMent :=  OrderDeployMent {
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: orderName,
		},
		Spec: SpecStOrder{
			Selector: SelectorSt{
				MatchLabels: MatchLabelSt{
					App: orderName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateStOrder{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: orderName,
					},
				},
				Spec: SpecTemplateStOrder{
					// NodeSelector: NodeSelectorSpecTemplateSt{
					// 	KubernetesIoHostname: "deploy host",
					// },
					Containers: []ContainerSpecTemplateSt{ 
						{
							Name: orderName,
							Image: image,
							ImagePullPolicy: "IfNotPresent",
							Resources: ResourceContainerSpecTemplateSt{
								Requests: RequestsResourceContainerSpecTemplateSt{
									Memory: "256Mi",
									CPU: "128m",
								},
							},
							Args: []string{"sh","-c","cp -a /var/data/. " + "/cert; exec orderer"},
							Env: []EnvContainerSpecTemplateSt{
								{
									Name: "CONFIGTX_ORDERER_KAFKA_BROKERS", 
									Value: "[kafka1:9092,kafka2:9092,kafka3:9092,kafka4:9092]",
								},
								{
									Name: "CONFIGTX_ORDERER_ORDERERTYPE",
									Value: "kafka",
								},
								{
									Name: "ORDERER_GENERAL_BATCHTIMEOUT",
									Value: "2s",
								},
								{
									Name: "ORDERER_GENERAL_GENESISFILE",
									Value: "/cert/channel-artifacts/genesis.block",
								},
								{
									Name: "ORDERER_GENERAL_GENESISMETHOD",
									Value: "file",
								},
								{
									Name: "ORDERER_GENERAL_LEDGERTYPE",
									Value: "file",
								},
								{
									Name: "ORDERER_GENERAL_LISTENADDRESS",
									Value: "0.0.0.0",
								},
								{
									Name: "ORDERER_GENERAL_LISTENPORT",
									Value: "7050",
								}, 
								{
									Name: "ORDERER_GENERAL_LOCALMSPDIR", 
									Value: "/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "." + orderDomain + "/msp",
								},
								{
									Name: "ORDERER_GENERAL_LOCALMSPID",
									Value: "OrdererMSP",
								},
								{
									Name: "ORDERER_GENERAL_LOGLEVEL",
									Value: "DEBUG", 
								},
								{
									Name: "ORDERER_GENERAL_MAXMESSAGECOUNT",
									Value: "10",
								},
								{
									Name: "ORDERER_GENERAL_MAXWINDOWSIZE",
									Value: "1000",
								},
								{
									Name: "ORDERER_FILELEDGER_LOCATION",
									Value: "/var/fabric/production/orderer",
								},
								{
									Name: "ORDERER_GENERAL_TLS_CERTIFICATE",
									Value: "/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "."  + orderDomain + "/tls/server.crt",
								},
								{
									Name: "ORDERER_GENERAL_TLS_PRIVATEKEY",
									Value:  "/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "."  + orderDomain + "/tls/server.key",
								}, 
								{
									Name: "ORDERER_GENERAL_TLS_ROOTCAS",
									Value: "[/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "."  + orderDomain + "tls/ca.crt]",
								},
								{
									Name: "ORDERER_KAFKA_RETRY_SHORTINTERVAL",
									Value: "1s",
								},
								{
									Name: "ORDERER_KAFKA_RETRY_SHORTTOTAL",
									Value: "30s",
								},
								{
									Name: "ORDERER_KAFKA_VERBOSE",
									Value: "true",
								}, 
							},
							Ports: []PortContainerSpecTemplateSt{
								{
									ContainerPort: 7050,
								},
							},
							WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric/orderers",
							VolumeMounts: []VolumeContainerSpecTemplateSt{ 
								{
									MountPath: "/var/data",
									Name: "order-data-store",
									//SubPath: clusterId + "/" + chainId + "/DataStore/" + orderName + "/data",
								}, 
							},    
						},
					},
					RestartPolicy: "Always",
					Volumes: []VolumeSpecTemplateSt{
						{
							Name: "order-data-store",
							Nfs: NfsVolumeSpecTemplateSt{
								Server: utils.BAAS_CFG.NfsInternalAddr,
								Path: utils.BAAS_CFG.NfsBasePath + "/" + chainId,
							},
						},
					},
				},
			},
		},
	}
	bytes, err := CreateDeployment(clusterId,namespaceId,orderKafkaDeployMent)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}

func CreateOrderSoloDeployment(clusterId string,namespaceId string,chainId string,image string,orderName string,orderDomain string)([]byte,error) {
	orderSoloDeployMent :=  OrderDeployMent {
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: orderName,
		},
		Spec: SpecStOrder{
			Selector: SelectorSt{
				MatchLabels: MatchLabelSt{
					App: orderName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateStOrder{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: orderName,
					},
				},
				Spec: SpecTemplateStOrder{
					// NodeSelector: NodeSelectorSpecTemplateSt{
					// 	KubernetesIoHostname: "deploy host",
					// },
					Containers: []ContainerSpecTemplateSt{ 
						{
							Name: orderName,
							Image: image,
							ImagePullPolicy: "IfNotPresent",
							Resources: ResourceContainerSpecTemplateSt{
								Requests: RequestsResourceContainerSpecTemplateSt{
									Memory: "256Mi",
									CPU: "128m",
								},
							},
							Args: []string{"sh","-c","cp -a /var/data/. " + "/cert; exec orderer"},
							Env: []EnvContainerSpecTemplateSt{
								{
									Name: "ORDERER_GENERAL_GENESISFILE",
									Value: "/cert/channel-artifacts/genesis.block",
								},
								{
									Name: "ORDERER_GENERAL_GENESISMETHOD",
									Value: "file",
								},
								{
									Name: "ORDERER_GENERAL_LISTENADDRESS",
									Value: "0.0.0.0",
								},
								{
									Name: "ORDERER_GENERAL_LISTENPORT",
									Value: "7050",
								}, 
								{
									Name: "ORDERER_GENERAL_LOCALMSPDIR", 
									Value: "/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "." + orderDomain + "/msp",
								},
								{
									Name: "ORDERER_GENERAL_LOCALMSPID",
									Value: "OrdererMSP",
								},
								{
									Name: "ORDERER_GENERAL_LOGLEVEL",
									Value: "DEBUG", 
								},
								{
									Name: "ORDERER_FILELEDGER_LOCATION",
									Value: "/var/fabric/production/orderer",
								},
								{
									Name: "ORDERER_GENERAL_TLS_CERTIFICATE",
									Value: "/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "."  + orderDomain + "/tls/server.crt",
								},
								{
									Name: "ORDERER_GENERAL_TLS_PRIVATEKEY",
									Value:  "/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "."  + orderDomain + "/tls/server.key",
								}, 
								{
									Name: "ORDERER_GENERAL_TLS_ROOTCAS",
									Value: "[/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "."  + orderDomain + "tls/ca.crt]",
								},
							},
							Ports: []PortContainerSpecTemplateSt{
								{
									ContainerPort: 7050,
								},
							},
							WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric/orderers",
							VolumeMounts: []VolumeContainerSpecTemplateSt{ 
								{
									MountPath: "/var/data",
									Name: "order-data-store",
									//SubPath: clusterId + "/" + chainId + "/DataStore/" + orderName + "/data",
								}, 
							},    
						},
					},
					RestartPolicy: "Always",
					Volumes: []VolumeSpecTemplateSt{
						{
							Name: "order-data-store",
							Nfs: NfsVolumeSpecTemplateSt{
								Server: utils.BAAS_CFG.NfsInternalAddr,
								Path: utils.BAAS_CFG.NfsBasePath + "/" + chainId,
							},
						},
					},
				},
			},
		},
	}
	bytes, err := CreateDeployment(clusterId,namespaceId,orderSoloDeployMent)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}

