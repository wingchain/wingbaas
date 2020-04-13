
package deployfabric
import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)
 
type CaDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataDeployMent `json:"metadata"`
	Spec SpecSt `json:"spec"`  
}

func CreateCaDeployment(clusterId string,namespaceId string,chainId string,image string,caName string,orgDomain string,caPriKey string)([]byte,error) {
	caDeployMent :=  CaDeployMent {
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: caName,
		},
		Spec: SpecSt{
			Selector: SelectorSt{
				MatchLabels: MatchLabelSt{
					App: caName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: caName,
					},
				},
				Spec: SpecTemplateSt{
					Affinity: AffinitySpecTemplateSt{
						NodeAffinity: NodeAffinityAffinitySpecTemplateSt{
							PreferredDuringSchedulingIgnoredDuringExecution: []ExeAffinitySpecTemplateSt{
								{
									Weight: 10,
									Preference: PreferenceExeAffinitySpecTemplateSt{
										MatchExpressions: []MatchPreferenceExeAffinitySpecTemplateSt{
											{
												Key: "team",
												Operator: "In",
												Values: []string{"baas"},
											},
										},		
									},
								},	
							},
						},
					},
					Containers: []ContainerSpecTemplateSt{
						{
							Name: caName,
							Image: image,
							ImagePullPolicy: "IfNotPresent",
							Resources: ResourceContainerSpecTemplateSt{
								Requests: RequestsResourceContainerSpecTemplateSt{ 
									Memory: "256Mi",
									CPU: "128m",
								},
							},
							Args: []string{"sh","-c","cp -a /var/data/. " + "/cert; exec fabric-ca-server start -b admin:adminpw -d"},
							Env: []EnvContainerSpecTemplateSt{
								{
									Name: "FABRIC_CA_SERVER_CA_NAME", 
									Value: caName,
								},
								{
									Name: "FABRIC_CA_HOME",
									Value: "/etc/fabric-ca-server",
								},
								{
									Name: "FABRIC_CA_SERVER_CA_CERTFILE",
									Value: "/cert/crypto-config/peerOrganizations/" + orgDomain + "/ca/ca." + orgDomain + "-cert.pem",
								},
								{
									Name: "FABRIC_CA_SERVER_CA_KEYFILE",
									Value: "/cert/crypto-config/peerOrganizations/" + orgDomain + "/ca/" + caPriKey,
								},
								{
									Name: "FABRIC_CA_SERVER_TLS_ENABLED",
									Value: "true",
								},
								{
									Name: "FABRIC_CA_SERVER_TLS_CERTFILE",
									Value: "/cert/crypto-config/peerOrganizations/" + orgDomain + "/ca/ca." + orgDomain + "-cert.pem",
								},
								{
									Name: "FABRIC_CA_SERVER_TLS_KEYFILE",
									Value: "/cert/crypto-config/peerOrganizations/" + orgDomain + "/ca/" + caPriKey,
								},
							},
							Ports: []PortContainerSpecTemplateSt{
								{
									ContainerPort: 7054,
								},
							},
							VolumeMounts: []VolumeContainerSpecTemplateSt{
								{
									MountPath: "/var/data/",
									Name: "ca-data-store",
								},
							},
						},
					},
					RestartPolicy: "Always",
					Volumes: []VolumeSpecTemplateSt{
						{
							Name: "ca-data-store",
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
	bytes, err := CreateDeployment(clusterId,namespaceId,caDeployMent)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


