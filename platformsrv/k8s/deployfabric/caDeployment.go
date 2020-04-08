
package deployfabric
import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/certgenerate/fabric"
)

// type CaDeployMent struct {
// 	APIVersion string `json:"apiVersion"`
// 	Kind string `json:"kind"`
// 	Metadata struct {
// 		Labels struct {
// 			HyperledgerFabricService string `json:"hyperledger.fabric.service"`
// 		} `json:"labels"`
// 		Name string `json:"name"`
// 	} `json:"metadata"`
// 	Spec struct {
// 		Replicas int `json:"replicas"`
// 		Strategy struct {
// 			Type string `json:"type"`
// 		} `json:"strategy"`
// 		Template struct {
// 			Metadata struct {
// 				Labels struct {
// 					HyperledgerFabricService string `json:"hyperledger.fabric.service"`
// 				} `json:"labels"`
// 			} `json:"metadata"`
// 			Spec struct {
// 				Affinity struct {
// 					NodeAffinity struct {
// 						PreferredDuringSchedulingIgnoredDuringExecution []struct {
// 							Weight int `json:"weight"`
// 							Preference struct {
// 								MatchExpressions []struct {
// 									Key string `json:"key"`
// 									Operator string `json:"operator"`
// 									Values []string `json:"values"`
// 								} `json:"matchExpressions"`
// 							} `json:"preference"`
// 						} `json:"preferredDuringSchedulingIgnoredDuringExecution"`
// 					} `json:"nodeAffinity"`
// 				} `json:"affinity"`
// 				Containers []struct {
// 					Name string `json:"name"`
// 					Image string `json:"image"`
// 					ImagePullPolicy string `json:"imagePullPolicy"`
// 					Resources struct {
// 						Requests struct {
// 							Memory string `json:"memory"`
// 							CPU string `json:"cpu"`
// 						} `json:"requests"`
// 					} `json:"resources"`
// 					Args []string `json:"args"`
// 					Env []struct {
// 						Name string `json:"name"`
// 						Value string `json:"value"`
// 					} `json:"env"`
// 					Ports []struct {
// 						ContainerPort int `json:"containerPort"`
// 					} `json:"ports"`
// 					VolumeMounts []struct {
// 						MountPath string `json:"mountPath"`
// 						Name string `json:"name"`
// 						SubPath string `json:"subPath,omitempty"`
// 					} `json:"volumeMounts"`
// 				} `json:"containers"`
// 				RestartPolicy string `json:"restartPolicy"`
// 				ImagePullSecrets []struct {
// 					Name string `json:"name"`
// 				} `json:"imagePullSecrets"`
// 				Volumes []struct {
// 					Name string `json:"name"`
// 					EmptyDir struct {
// 					} `json:"emptyDir,omitempty"`
// 					Nfs struct {
// 						Server string `json:"server"`
// 						Path string `json:"path"`
// 					} `json:"nfs,omitempty"`
// 				} `json:"volumes"`
// 			} `json:"spec"`
// 		} `json:"template"`
// 	} `json:"spec"`
// }

type LabelsSt struct {
	HyperledgerFabricService string `json:"hyperledger.fabric.service"`
}

type StrategySt struct {
	Type string `json:"type"`
}

type MetadataTemplateSt struct {
	Labels LabelsSt `json:"labels"`
}

type MatchPreferenceExeAffinitySpecTemplateSt struct {
	Key string `json:"key"`
	Operator string `json:"operator"`
	Values []string `json:"values"`
}

type PreferenceExeAffinitySpecTemplateSt struct {
	MatchExpressions []MatchPreferenceExeAffinitySpecTemplateSt `json:"matchExpressions"`
} 

type ExeAffinitySpecTemplateSt struct {
	Weight int `json:"weight"`
	Preference PreferenceExeAffinitySpecTemplateSt `json:"preference"`
}

type NodeAffinityAffinitySpecTemplateSt struct {
	PreferredDuringSchedulingIgnoredDuringExecution []ExeAffinitySpecTemplateSt `json:"preferredDuringSchedulingIgnoredDuringExecution"`
}

type AffinitySpecTemplateSt struct {
	NodeAffinity NodeAffinityAffinitySpecTemplateSt `json:"nodeAffinity"`
}

type RequestsResourceContainerSpecTemplateSt struct {
	Memory string `json:"memory"`
	CPU string `json:"cpu"`
}

type ResourceContainerSpecTemplateSt struct {
	Requests RequestsResourceContainerSpecTemplateSt `json:"requests"`
}

type EnvContainerSpecTemplateSt struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type PortContainerSpecTemplateSt struct {
	ContainerPort int `json:"containerPort"`
}

type VolumeContainerSpecTemplateSt struct {
	MountPath string `json:"mountPath"`
	Name string `json:"name"`
	SubPath string `json:"subPath,omitempty"`
}

type ContainerSpecTemplateSt struct {
	Name string `json:"name"`
	Image string `json:"image"`
	ImagePullPolicy string `json:"imagePullPolicy"`
	Resources ResourceContainerSpecTemplateSt `json:"resources"`
	Args []string `json:"args"`
	Env []EnvContainerSpecTemplateSt `json:"env"`
	Ports []PortContainerSpecTemplateSt `json:"ports"`
	VolumeMounts []VolumeContainerSpecTemplateSt `json:"volumeMounts"`
} 

type ImagePullSecretSpecTemplateSt struct {
	Name string `json:"name"`
}

type NfsVolumeSpecTemplateSt struct {
	Server string `json:"server"`
	Path string `json:"path"`
}

type EmptyDirVolumeSpecTemplateSt struct {
}

type VolumeSpecTemplateSt struct {
	Name string `json:"name"`
	EmptyDir EmptyDirVolumeSpecTemplateSt `json:"emptyDir,omitempty"`
	Nfs NfsVolumeSpecTemplateSt `json:"nfs,omitempty"`
}

type SpecTemplateSt struct {
	Affinity AffinitySpecTemplateSt `json:"affinity"`
	Containers []ContainerSpecTemplateSt `json:"containers"`
	RestartPolicy string `json:"restartPolicy"`
	ImagePullSecrets []ImagePullSecretSpecTemplateSt `json:"imagePullSecrets"`
	Volumes []VolumeSpecTemplateSt `json:"volumes"`
}

type TemplateSt struct {
	Metadata MetadataTemplateSt `json:"metadata"`
	Spec SpecTemplateSt `json:"spec"`
} 

type SpecSt struct {
	Replicas int `json:"replicas"`
	Strategy StrategySt `json:"strategy"`
	Template TemplateSt `json:"template"`
}


type CaDeployMent struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata struct {
		Labels LabelsSt `json:"labels"`
		Name string `json:"name"`
	} `json:"metadata"`
	Spec SpecSt `json:"spec"`
}

func CreateCaDeployment(clusterId string,namespaceId string,chainId string,image string,caName string,orgDomain string,caPriKey string)([]byte,error) {
	caDeployMent :=  CaDeployMent {
		APIVersion: "v1",
		Kind: "Deployment",
		Metadata: struct{
			Labels LabelsSt `json:"labels"`
			Name string `json:"name"`
		}{
			Labels: LabelsSt{HyperledgerFabricService: caName}, 
			Name: caName,
		},
		Spec: SpecSt{
			Replicas: 1,
			Strategy: StrategySt{
				Type: "Recreate",
			},
			Template: TemplateSt{
				Metadata: MetadataTemplateSt{
					Labels: LabelsSt{
						HyperledgerFabricService: caName,
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
							Args: []string{"sh","-c","cp -a /var/data/" + chainId + "/. /cert; exec fabric-ca-server start -b admin:adminpw -d"},
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
									SubPath: clusterId,
								},
								{
									MountPath: "/cert",
									Name: "pod-cert",
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
								Path: utils.BAAS_CFG.NfsBasePath,
							},
						},
						{
							Name: "pod-cert",
							EmptyDir: EmptyDirVolumeSpecTemplateSt{}, 
						},
					},
				},
			},
		},
	}
	bytes, err := json.Marshal(caDeployMent)
	if err != nil {
		logger.Errorf("CreateCaDeployment: Marshal deployment error,%v",err) 
		return nil,fmt.Errorf("CreateCaDeployment: Marshal deployment error,%v",err)
	}
	bl,yamlStr := fabric.JsonToYaml(string(bytes))
	if !bl {
		logger.Errorf("CreateCaDeployment: json2yaml str error") 
		return nil,fmt.Errorf("CreateCaDeployment: json2yaml str error")
	}
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("CreateCaDeployment: cluster nil,cluster id =%s",clusterId)
		return nil,fmt.Errorf("CreateCaDeployment: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_APP + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.DEPLOYMENTS
	bytes,err = utils.RequestWithCertAndBody(reqUrl,utils.REQ_POST,cluster.Cert,cluster.Key,yamlStr)
	if err != nil { 
		logger.Errorf("CreateCaDeployment: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("CreateCaDeployment: RequestWithCertAndBody err,%v", err)
	}
	var result interface{} 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("CreateCaDeployment: create result err,%v", err)
		return nil,fmt.Errorf("CreateCaDeployment: create result err,%v", err)
	}
	logger.Debug("CreateCaDeployment: create success,result=,%v", result) 
	return nil,nil
}

func DeleteCaDeployment(clusterId string,namespaceId string,deployName string)([]byte,error) {
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("DeleteCaDeployment: cluser nil,cluser id =%s",clusterId)
		return nil,fmt.Errorf("DeleteCaDeployment: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_APP + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.DEPLOYMENTS + "/" + deployName
	bytes,err := utils.RequestWithCert(reqUrl,utils.REQ_DELETE,cluster.Cert,cluster.Key)
	if err != nil { 
		logger.Errorf("DeleteCaDeployment: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("DeleteCaDeployment: RequestWithCertAndBody err,%v", err)
	}
	var result interface{} 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("DeleteCaDeployment: delete result err,%v", err)
		return nil,fmt.Errorf("DeleteCaDeployment: delete result err,%v", err)
	}
	logger.Debug("DeleteCaDeployment: delete success,result=,%v", result) 
	return nil,nil
}


