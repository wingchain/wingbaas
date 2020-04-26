
package deployfabric
import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type OrderDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataDeployMent `json:"metadata"`
	Spec SpecSt `json:"spec"`  
}  

func CreateOrderKafkaDeployment(clusterId string,node string,namespaceId string,chainId string,image string,orderName string,orderDomain string)([]byte,error) {
	orderKafkaDeployMent :=  OrderDeployMent {
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: orderName,
			Labels: LabelsSt{
				App: orderName, 
			},
		},
		Spec: SpecSt{
			Selector: SelectorSt{
				MatchLabels: LabelsSt{ 
					App: orderName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: orderName,
					},
				},
				Spec: SpecTemplateSt{
					NodeName: node,
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
									Value: "[kafka0:9092,kafka1:9092,kafka2:9092,kafka3:9092]",
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
									Name: "ORDERER_GENERAL_TLS_ENABLED",
									Value: "true",
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
									Value: "[/cert/crypto-config/ordererOrganizations/" + orderDomain + "/orderers/" + orderName + "."  + orderDomain + "/tls/ca.crt]",
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
								{
									Name: "CORE_VM_DOCKER_ATTACHSTDOUT",
									Value: "true",
								},
								{
									Name: "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE",
									Value: "host",
								},
								{
									Name: "CORE_VM_ENDPOINT",
									Value: "unix:///host/var/run/docker.sock",
								}, 
								{
									Name: "CORE_VM_DOCKER_HOSTCONFIG_DNSSEARCH",
									Value: namespaceId + ".svc.cluster.local",
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
									Name: "order-cert",
									//SubPath: clusterId + "/" + chainId + "/DataStore/" + orderName + "/data",
								}, 
								{
									MountPath: "/host/var/run/",
									Name: "host-vol-var-run",
								}, 
							},    
						},
					},
					RestartPolicy: "Always", 
					Hostname: orderName,
					Volumes: []interface{}{
						// VolumeSpecTemplateSt{
						// 	Name: "order-cert",
						// 	Nfs: NfsVolumeSpecTemplateSt{
						// 		Server: utils.BAAS_CFG.NfsInternalAddr,
						// 		Path: utils.BAAS_CFG.NfsBasePath + "/" + chainId, 
						// 	},
						// },
						VolumeSpecTemplateHostSt{
							Name: "order-cert",
							HostPath: HostPathVolumeSpecTemplateSt{
								Path: "/home/nfs/" + chainId, 
							}, 
						},
						VolumeSpecTemplateHostSt{
							Name: "host-vol-var-run",
							HostPath: HostPathVolumeSpecTemplateSt{
								Path: "/var/run/",  
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
		Spec: SpecSt{
			Selector: SelectorSt{
				MatchLabels: LabelsSt{
					App: orderName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: orderName,
					},
				},
				Spec: SpecTemplateSt{
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
					Volumes: []interface{}{
						VolumeSpecTemplateSt{
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

