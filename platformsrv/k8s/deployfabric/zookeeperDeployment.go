
package deployfabric
import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type ZookeeperDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataDeployMent `json:"metadata"`
	Spec SpecSt `json:"spec"`  
} 

func CreateZookeeperDeployment(clusterId string,namespaceId string,chainId string,zkId string,image string,zookeeperName string)([]byte,error) {
	zkDeployMent :=  ZookeeperDeployMent {
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: zookeeperName,
		},
		Spec: SpecSt{
			Selector: SelectorSt{
				MatchLabels: MatchLabelSt{
					App: zookeeperName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: zookeeperName,
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
							Name: zookeeperName,
							Image: image,
							ImagePullPolicy: "IfNotPresent",
							Resources: ResourceContainerSpecTemplateSt{
								Requests: RequestsResourceContainerSpecTemplateSt{
									Memory: "256Mi",
									CPU: "128m",
								},
							},
							Env: []EnvContainerSpecTemplateSt{
								{
									Name: "ZOO_MY_ID",
									Value: zkId,
								},
								{
									Name: "ZOO_SERVERS",
									Value: "server.1=zookeeper1:2888:3888 server.2=zookeeper2:2888:3888 server.3=zookeeper3:2888:3888",
								},
							},
							Ports: []PortContainerSpecTemplateSt{
								{
									ContainerPort: 2181,
								},
								{
									ContainerPort: 2888,
								},
								{
									ContainerPort: 3888,
								},	
							},
							VolumeMounts: []VolumeContainerSpecTemplateSt{
								{
									MountPath: "/data",
									Name: "zk-data-store",
									//SubPath: clusterId + "/" + chainId + "/DataStore/" + zookeeperName + "/data",
								}, 
							},
						},
					},
					Hostname: zookeeperName,
					RestartPolicy: "Always",
					Volumes: []VolumeSpecTemplateSt{
						{
							Name: "zk-data-store",
							Nfs: NfsVolumeSpecTemplateSt{
								Server: utils.BAAS_CFG.NfsInternalAddr,
								Path: utils.BAAS_CFG.NfsBasePath,
							},
						},
					},
				},
			},
		},
	}
	bytes, err := CreateDeployment(clusterId,namespaceId,zkDeployMent)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


