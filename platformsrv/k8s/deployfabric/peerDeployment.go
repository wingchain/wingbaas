
package deployfabric
import (
	"os"
	"fmt"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
)

type PeerDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataDeployMent `json:"metadata"`
	Spec SpecSt `json:"spec"`  
}  

type PeerDeploymentPara struct{
	PeerImage 	string
	CcenvImage 	string
	BaseosImage string
	PeerName 	string
	RawPeerName string 
	AnchorPeer 	string
	PeerDomain 	string
	OrgName 	string
} 

func CreatePeerDeployment(clusterId string,node string,namespaceId string,chainId string,p PeerDeploymentPara)([]byte,error) {
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("CreatePeerDeployment: get cluster failed,id=%s",clusterId)
		return nil,fmt.Errorf("CreatePeerDeployment: get cluster failed,id=%s",clusterId)
	}
	peerDeployMent :=  PeerDeployMent { 
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: p.PeerName,
			Labels: LabelsSt{
				App: p.PeerName, 
			},
		},
		Spec: SpecSt{
			Selector: SelectorSt{
				MatchLabels: LabelsSt{
					App: p.PeerName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: p.PeerName,
					},
				},
				Spec: SpecTemplateSt{
					NodeName: node, 
					// NodeSelector: NodeSelectorSpecTemplateSt{
					// 	KubernetesIoHostname: "172-16-254-33", 
					// 	//Hostname: "172-16-254-33", 
					// },
					InitContainers: []ContainerSpecTemplateSt{
						{
							Name: "pre-pull-ccenv",
							Image: p.CcenvImage,
							ImagePullPolicy: "IfNotPresent",
							Resources: ResourceContainerSpecTemplateSt{
								Requests: RequestsResourceContainerSpecTemplateSt{
									Memory: "256Mi",
									CPU: "128m",
								},
							},
							Command: []string{"echo","SUCCESS"},
						},
						{
							Name: "pre-pull-baseos",
							Image: p.BaseosImage,
							ImagePullPolicy: "IfNotPresent",
							Resources: ResourceContainerSpecTemplateSt{
								Requests: RequestsResourceContainerSpecTemplateSt{
									Memory: "256Mi",
									CPU: "128m",
								},
							},
							Command: []string{"echo","SUCCESS"},
						},
					},
					Containers: []ContainerSpecTemplateSt{ 
						{
							Name: p.PeerName, 
							Image: p.PeerImage,
							ImagePullPolicy: "IfNotPresent",
							Resources: ResourceContainerSpecTemplateSt{
								Requests: RequestsResourceContainerSpecTemplateSt{
									Memory: "768Mi",
									CPU: "128m",
								},
							},
							Args: []string{"sh","-c","cp -a /var/data/. " + "/cert; exec peer node start"},
							Env: []EnvContainerSpecTemplateSt{
								{
									Name: "CORE_LOGGING_LEVEL", 
									Value: "debug",
								},
								{
									Name: "FABRIC_LOGGING_SPEC", 
									Value: "debug",
								},
								{
									Name: "CORE_PEER_ID", 
									Value: p.PeerName,
								},
								{
									Name: "CORE_PEER_ADDRESS",
									Value: p.PeerName + ":7051",
								},
								{
									Name: "CORE_PEER_NETWORKID",
									Value: "dev-" + namespaceId,
								},
								{
									Name: "CORE_PEER_FILESYSTEMPATH",
									Value: "/var/fabric/production",
								},
								{
									Name: "CORE_PEER_EVENTS_TIMEOUT",
									Value: "3000ms",
								},
								{
									Name: "CORE_PEER_GOSSIP_BOOTSTRAP",
									Value: p.AnchorPeer + ":7051",
								},
								{
									Name: "CORE_PEER_GOSSIP_EXTERNALENDPOINT",
									Value: p.PeerName + ":7051",
								},
								{
									Name: "CORE_PEER_GOSSIP_ORGLEADER",
									Value: "false",
								}, 
								{
									Name: "CORE_PEER_GOSSIP_SKIPHANDSHAKE", 
									Value: "true",
								},
								{
									Name: "CORE_PEER_GOSSIP_USELEADERELECTION",
									Value: "true",
								},
								{
									Name: "CORE_PEER_MSPCONFIGPATH",
									Value: "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.RawPeerName + "."  + p.PeerDomain + "/msp", 
									//Value: "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/users/Admin@"  + p.PeerDomain + "/msp",
								},
								{
									Name: "CORE_PEER_LOCALMSPID",
									Value: p.OrgName + "MSP",
								},
								{
									Name: "CORE_PEER_PROFILE_ENABLED",
									Value: "false",
								},
								{
									Name: "CORE_PEER_TLS_ENABLED",
									Value: "true",
								},
								{
									Name: "CORE_PEER_TLS_CERT_FILE",
									Value: "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.RawPeerName + "."  + p.PeerDomain + "/tls/server.crt",
								},
								{
									Name: "CORE_PEER_TLS_KEY_FILE",
									Value:  "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.RawPeerName + "."  + p.PeerDomain + "/tls/server.key",
								}, 
								{
									Name: "CORE_PEER_TLS_ROOTCERT_FILE",
									Value: "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.RawPeerName + "."  + p.PeerDomain + "/tls/ca.crt",
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
									Name: "CORE_PEER_ADDRESSAUTODETECT",
									Value: "true", 
								}, 
								{
									Name: "CORE_VM_DOCKER_HOSTCONFIG_DNS",
									Value: cluster.PublicIp,
								},
								{
									Name: "CORE_VM_DOCKER_HOSTCONFIG_DNSSEARCH",
									Value: namespaceId + ".svc.cluster.local",
								},
								{
									Name: "CORE_CHAINCODE_BUILDER",   
									Value: p.CcenvImage,
								}, 
								// {
								// 	Name: "CORE_PEER_CHAINCODELISTENADDRESS",
								// 	Value: p.PeerName + ":7052",
								// },
								{
									Name: "CORE_CHAINCODE_GOLANG_RUNTIME",
									Value: p.BaseosImage,
								},
								{
									Name: "CORE_CHAINCODE_LOGGING_LEVEL",
									Value: "debug",
								},
								{
									Name: "CORE_CHAINCODE_LOGGING_SHIM",
									Value: "debug",
								}, 

							},
							Ports: []PortContainerSpecTemplateSt{
								{
									ContainerPort: 7051,
								},
								{
									ContainerPort: 7052,
								},
								{
									ContainerPort: 7053,
								},
							},
							VolumeMounts: []VolumeContainerSpecTemplateSt{ 
								{
									MountPath: "/var/data",
									Name: "peer-cert",
								}, 
								{
									MountPath: "/host/var/run/",
									Name: "host-vol-var-run",
								}, 
							},    
						},
					},
					RestartPolicy: "Always",
					Hostname: p.PeerName,
					Volumes: []interface{}{
						VolumeSpecTemplateSt{
							Name: "peer-cert",
							Nfs: NfsVolumeSpecTemplateSt{
								Server: utils.BAAS_CFG.NfsInternalAddr,
								Path: utils.BAAS_CFG.NfsBasePath + "/" + chainId,
							},
						},
						// VolumeSpecTemplateHostSt{
						// 	Name: "peer-cert",
						// 	HostPath: HostPathVolumeSpecTemplateSt{
						// 		Path: "/home/nfs/" + chainId, 
						// 	}, 
						// },
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
	bytes, err := CreateDeployment(clusterId,namespaceId,peerDeployMent)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)    
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


