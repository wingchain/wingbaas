
package deployfabric
import (
	"os"
	"fmt"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
)

type ToolsDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataDeployMent `json:"metadata"`
	Spec SpecSt `json:"spec"`  
}  

type ToolsDeploymentPara struct{
	ToolsImage	string
	PeerName 	string
	PeerDomain 	string
	OrgName 	string
} 

func CreateToolsDeployment(clusterId string,node string,namespaceId string,chainId string,p ToolsDeploymentPara)([]byte,error) {
	appName := "cli" + namespaceId
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("CreateToolsDeployment: get cluster failed,id=%s",clusterId)
		return nil,fmt.Errorf("CreateToolsDeployment: get cluster failed,id=%s",clusterId)
	}
	toolsDeployMent :=  ToolsDeployMent { 
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: appName,
			Labels: LabelsSt{
				App: appName, 
			},
		},
		Spec: SpecSt{
			Selector: SelectorSt{
				MatchLabels: LabelsSt{
					App: appName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: appName,
					},
				},
				Spec: SpecTemplateSt{
					NodeName: node,
					Containers: []ContainerSpecTemplateSt{ 
						{
							Name: appName, 
							Image: p.ToolsImage,
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
									Name: "FABRIC_LOGGING_SPEC", 
									Value: "INFO",
								},
								{
									Name: "GOPATH", 
									Value: "/opt/gopath",
								},
								{
									Name: "CORE_VM_ENDPOINT",
									Value: "unix:///host/var/run/docker.sock",
								},
								{
									Name: "CORE_PEER_ID", 
									Value: appName,
								},
								{
									Name: "CORE_PEER_ADDRESS",
									Value: p.PeerName + ":7051",
								},
								{
									Name: "CORE_PEER_LOCALMSPID",
									Value: p.OrgName + "MSP",
								},
								{
									Name: "CORE_PEER_NETWORKID",
									Value: "dev-" + namespaceId,
								},
								{
									Name: "CORE_PEER_MSPCONFIGPATH",
									Value: "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.PeerName + "."  + p.PeerDomain + "/msp", 
								},
								{
									Name: "CORE_PEER_TLS_ENABLED",
									Value: "true",
								},
								{
									Name: "CORE_PEER_TLS_CERT_FILE",
									Value: "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.PeerName + "."  + p.PeerDomain + "/tls/server.crt",
								},
								{
									Name: "CORE_PEER_TLS_KEY_FILE",
									Value:  "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.PeerName + "."  + p.PeerDomain + "/tls/server.key",
								}, 
								{
									Name: "CORE_PEER_TLS_ROOTCERT_FILE",
									Value: "/cert/crypto-config/peerOrganizations/" + p.PeerDomain + "/peers/" + p.PeerName + "."  + p.PeerDomain + "/tls/ca.crt",
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
									Name: "CORE_PEER_ADDRESSAUTODETECT",
									Value: "true", 
								},
								{
									Name: "CORE_VM_DOCKER_HOSTCONFIG_DNSSEARCH",
									Value: namespaceId + ".svc.cluster.local",
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
	bytes, err := CreateDeployment(clusterId,namespaceId,toolsDeployMent)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)    
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


