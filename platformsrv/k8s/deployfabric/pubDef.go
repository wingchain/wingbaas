
package deployfabric

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/certgenerate/fabric"
)

const (
	KUBERNETES_DEPLOY_FAILED	string	=	"Failure"
)

//service struct define

type MetadataService struct {
	Name string `json:"name"`
}

type PortSpecService struct {
	Name string `json:"name"`
	Port int `json:"port"`
	TargetPort int `json:"targetPort"`
}

type SelectorSpecService struct {
	App string `json:"app"`
}

type SpecService struct {
	Type string `json:"type"`
	Ports []PortSpecService `json:"ports"`
	Selector SelectorSpecService `json:"selector"`
}

//deployment struct define 

// type Deployment struct {
// 	APIVersion string `json:"apiVersion"`
// 	Kind string `json:"kind"`
// 	Metadata struct {
// 		Name string `json:"name"`
// 	} `json:"metadata"`
// 	Spec struct {
// 		Replicas int `json:"replicas"`
// 		Selector struct {
// 			MatchLabels struct {
// 				App string `json:"app"`
// 			} `json:"matchLabels"`
// 		} `json:"selector"`
// 		Strategy struct {
// 			Type string `json:"type"`
// 		} `json:"strategy"`
// 		Template struct {
// 			Metadata struct {
// 				Labels struct {
// 					App string `json:"app"`
// 				} `json:"labels"`
// 			} `json:"metadata"`
// 			Spec struct {
// 				Affinity struct {
// 					NodeAffinity struct {
// 						PreferredDuringSchedulingIgnoredDuringExecution []struct {
// 							Preference struct {
// 								MatchExpressions []struct {
// 									Key string `json:"key"`
// 									Operator string `json:"operator"`
// 									Values []string `json:"values"`
// 								} `json:"matchExpressions"`
// 							} `json:"preference"`
// 							Weight int `json:"weight"`
// 						} `json:"preferredDuringSchedulingIgnoredDuringExecution"`
// 					} `json:"nodeAffinity"`
// 				} `json:"affinity"`
// 				Containers []struct {
// 					Args []string `json:"args"`
// 					Env []struct {
// 						Name string `json:"name"`
// 						Value string `json:"value"`
// 					} `json:"env"`
// 					Image string `json:"image"`
// 					ImagePullPolicy string `json:"imagePullPolicy"`
// 					Name string `json:"name"`
// 					Ports []struct {
// 						ContainerPort int `json:"containerPort"`
// 					} `json:"ports"`
// 					Resources struct {
// 						Requests struct {
// 							CPU string `json:"cpu"`
// 							Memory string `json:"memory"`
// 						} `json:"requests"`
// 					} `json:"resources"`
// 					VolumeMounts []struct {
// 						MountPath string `json:"mountPath"`
// 						Name string `json:"name"`
// 					} `json:"volumeMounts"`
// 				} `json:"containers"`
// 				ImagePullSecrets interface{} `json:"imagePullSecrets"`
// 				RestartPolicy string `json:"restartPolicy"`
// 				Volumes []struct {
// 					Name string `json:"name"`
// 					Nfs struct {
// 						Path string `json:"path"`
// 						Server string `json:"server"`
// 					} `json:"nfs"`
// 				} `json:"volumes"`
// 			} `json:"spec"`
// 		} `json:"template"`
// 	} `json:"spec"`
// }


type LabelsSt struct {
	App string `json:"app"`
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
	WorkingDir string `json:"workingDir,omitempty"` 
	Command []string `json:"command,omitempty"`
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
	//EmptyDir EmptyDirVolumeSpecTemplateSt `json:"emptyDir,omitempty"`
	Nfs NfsVolumeSpecTemplateSt `json:"nfs,omitempty"` 
}

type NodeSelectorSpecTemplateSt struct {
	KubernetesIoHostname string `json:"kubernetes.io/hostname"`
}

type SpecTemplateSt struct {
	Affinity AffinitySpecTemplateSt `json:"affinity"`
	Containers []ContainerSpecTemplateSt `json:"containers"`
	RestartPolicy string `json:"restartPolicy"`
	ImagePullSecrets []ImagePullSecretSpecTemplateSt `json:"imagePullSecrets"`
	Hostname string `json:"hostname"`
	Volumes []VolumeSpecTemplateSt `json:"volumes"`
}

type TemplateSt struct {
	Metadata MetadataTemplateSt `json:"metadata"`
	Spec SpecTemplateSt `json:"spec"`
} 
type MatchLabelSt struct {
	App	string `json:"app"`
}

type SelectorSt struct {
	MatchLabels MatchLabelSt `json:"matchLabels"`
}

type SpecSt struct {
	Selector SelectorSt `json:"selector"`
	Replicas int `json:"replicas"`
	Strategy StrategySt `json:"strategy"`
	Template TemplateSt `json:"template"`
}

type MetadataDeployMent struct {
	Name string `json:"name"` 
}

type DeployResult struct {
	Kind string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata interface{} `json:"metadata"`
	Status interface{} `json:"status"`
	Message string `json:"message"`
	Reason string `json:"reason"`
	Details struct {
		Name string `json:"name"`
		Group string `json:"group"`
		Kind string `json:"kind"`
	} `json:"details"`
	Code int `json:"code"` 
}

func CreateDeployment(clusterId string,namespaceId string,deployObj interface{})([]byte,error) { 
	bytes, err := json.Marshal(deployObj) 
	if err != nil {
		logger.Errorf("CreateDeployment: Marshal deployment error,%v",err) 
		return nil,fmt.Errorf("CreateDeployment: Marshal deployment error,%v",err)
	}
	bl,yamlStr := fabric.JsonToYaml(string(bytes))
	if !bl {
		logger.Errorf("CreateDeployment: json2yaml str error") 
		return nil,fmt.Errorf("CreateDeployment: json2yaml str error")
	}
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("CreateDeployment: cluster nil,cluster id =%s",clusterId)
		return nil,fmt.Errorf("CreateDeployment: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_APP + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.DEPLOYMENTS
	bytes,err = utils.RequestWithCertAndBody(reqUrl,utils.REQ_POST,cluster.Cert,cluster.Key,yamlStr)
	if err != nil { 
		logger.Errorf("CreateDeployment: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("CreateDeployment: RequestWithCertAndBody err,%v", err)
	}
	var result DeployResult 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("CreateDeployment: create result err,%v", err)
		return nil,fmt.Errorf("CreateDeployment: create result err,%v", err)
	}
	status,ok := (result.Status).(string) 
	if ok {
		if status == KUBERNETES_DEPLOY_FAILED {
			logger.Errorf("CreateDeployment: create failed,result=%s",string(bytes))
			logger.Errorf("CreateDeployment: create failed,yamlstr=%s",yamlStr)
			return nil,fmt.Errorf("CreateDeployment: create failed")
		}
	}
	logger.Debug("CreateDeployment: create success,result str=")
	logger.Debug(string(bytes))
	return nil,nil
}

func DeleteDeployment(clusterId string,namespaceId string,deployName string)([]byte,error) {
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("DeleteDeployment: cluser nil,cluser id =%s",clusterId)
		return nil,fmt.Errorf("DeleteDeployment: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_APP + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.DEPLOYMENTS + "/" + deployName
	bytes,err := utils.RequestWithCert(reqUrl,utils.REQ_DELETE,cluster.Cert,cluster.Key)
	if err != nil { 
		logger.Errorf("DeleteDeployment: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("DeleteDeployment: RequestWithCertAndBody err,%v", err)
	}
	var result interface{} 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("DeleteDeployment: delete result err,%v", err)
		return nil,fmt.Errorf("DeleteDeployment: delete result err,%v", err)
	}
	logger.Debug("DeleteDeployment: delete success,result=,%v", result)
	return nil,nil 
}

func CreateService(clusterId string,namespaceId string,serviceObj interface{})([]byte,error) {
	bytes, err := json.Marshal(serviceObj)
	if err != nil {
		logger.Errorf("CreateService: Marshal service error,%v",err) 
		return nil,fmt.Errorf("CreateService: Marshal service error,%v",err)
	}
	bl,yamlStr := fabric.JsonToYaml(string(bytes))
	if !bl {
		logger.Errorf("CreateService: json2yaml str error") 
		return nil,fmt.Errorf("CreateService: json2yaml str error")
	}
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("CreateService: cluster nil,cluster id =%s",clusterId)
		return nil,fmt.Errorf("CreateService: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_V1 + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.SERVICES
	bytes,err = utils.RequestWithCertAndBody(reqUrl,utils.REQ_POST,cluster.Cert,cluster.Key,yamlStr)
	if err != nil { 
		logger.Errorf("CreateService: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("CreateService: RequestWithCertAndBody err,%v", err)
	}
	var result DeployResult 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("CreateService: create result err,%v", err)
		return nil,fmt.Errorf("CreateService: create result err")
	}
	status,ok := (result.Status).(string)
	if ok {
		if status == KUBERNETES_DEPLOY_FAILED { 
			logger.Errorf("CreateService: create failed,result=%s",string(bytes))
			logger.Errorf("CreateService: create failed,yamlstr=%s",yamlStr)
			return nil,fmt.Errorf("CreateService: create failed")
		}
	}
	logger.Debug("CreateService: create success,result str=")
	logger.Debug(string(bytes))
	return nil,nil 
}

func DeleteService(clusterId string,namespaceId string,serviceName string)([]byte,error) { 
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("DeleteService: cluser nil,cluser id =%s",clusterId)
		return nil,fmt.Errorf("DeleteService: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_V1 + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.SERVICES + "/" + serviceName
	bytes,err := utils.RequestWithCert(reqUrl,utils.REQ_DELETE,cluster.Cert,cluster.Key)
	if err != nil { 
		logger.Errorf("DeleteService: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("DeleteService: RequestWithCertAndBody err,%v", err)
	}
	var result interface{} 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("DeleteService: delete result err,%v", err)
		return nil,fmt.Errorf("DeleteService: delete result err,%v", err)
	}
	logger.Debug("DeleteService: delete success,result=,%v", result)
	return nil,nil
}


