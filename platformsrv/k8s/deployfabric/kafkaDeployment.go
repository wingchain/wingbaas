
package deployfabric
import (
	"os"
	"github.com/wingbaas/platformsrv/utils"
)

type KafkaDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataDeployMent `json:"metadata"`
	Spec SpecSt `json:"spec"`  
} 

func CreateKafkaDeployment(clusterId string,namespaceId string,chainId string,kafkaId string,image string,kafkaName string)([]byte,error) {
	kafkaDeployMent :=  KafkaDeployMent {
		APIVersion: "apps/v1",
		Kind: "Deployment",
		Metadata: MetadataDeployMent{
			Name: kafkaName,
		},
		Spec: SpecSt{
			Selector: SelectorSt{
				MatchLabels: MatchLabelSt{
					App: kafkaName,
				},
			},
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						App: kafkaName,
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
							Name: kafkaName,
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
									Name: "KAFKA_BROKER_ID",
									Value: kafkaId,
								},
								{
									Name: "KAFKA_DEFAULT_REPLICATION_FACTOR",
									Value: "3",
								},
								{
									Name: "KAFKA_LOG_RETENTION_MS",
									Value: "-1",
								},
								{
									Name: "KAFKA_MESSAGE_MAX_BYTES",
									Value: "103809024",
								},
								{
									Name: "KAFKA_MIN_INSYNC_REPLICAS",
									Value: "2",
								},
								{
									Name: "KAFKA_REPLICA_FETCH_MAX_BYTES",
									Value: "103809024",
								},
								{
									Name: "KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE",
									Value: "false",
								},
								{
									Name: "KAFKA_ZOOKEEPER_CONNECT",
									Value: "zookeeper1:2181,zookeeper2:2181,zookeeper3:2181",
								}, 
							},
							Ports: []PortContainerSpecTemplateSt{
								{
									ContainerPort: 9092,
								},
							},
							VolumeMounts: []VolumeContainerSpecTemplateSt{
								{
									MountPath: "/data",
									Name: "kafka-data-store",
									//SubPath: clusterId + "/" + chainId + "/DataStore/" + kafkaName + "/data",
								}, 
							},    
						},
					},
					RestartPolicy: "Always",
					Volumes: []VolumeSpecTemplateSt{
						{
							Name: "kafka-data-store",
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
	bytes, err := CreateDeployment(clusterId,namespaceId,kafkaDeployMent)
	if err != nil {
		certPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
		nfsPath :=  utils.BAAS_CFG.NfsLocalRootDir + chainId
		os.RemoveAll(certPath)
		os.RemoveAll(nfsPath)
		DeleteNamespace(clusterId,namespaceId)
	}
	return bytes,err
}


