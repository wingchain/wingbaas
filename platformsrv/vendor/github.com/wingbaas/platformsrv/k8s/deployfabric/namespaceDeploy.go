
package deployfabric
import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/certgenerate/fabric/gopkg.in/yaml.v2"
)

type NamespaceDeploy struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}

type NamespaceDeployResult struct {
	Kind string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
		Name string `json:"name"`
		SelfLink string `json:"selfLink"`
		UID string `json:"uid"`
		ResourceVersion string `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Spec struct {
		Finalizers []string `json:"finalizers"`
	} `json:"spec"`
	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`
}

type NamespaceDeployResultError struct {
	Kind string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
	} `json:"metadata"`
	Status string `json:"status"`
	Message string `json:"message"`
	Reason string `json:"reason"`
	Details struct {
	} `json:"details"`
	Code int `json:"code"`
}

func CreateNamespace(clusterId string,namespaceId string)([]byte,error) {
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("CreateNamespace: cluster nil,cluster id =%s",clusterId)
		return nil,fmt.Errorf("CreateNamespace: cluster nil,cluster id =%s",clusterId)
	}
	var obj NamespaceDeploy
	obj.ApiVersion = "v1"
	obj.Kind = "Namespace"
	obj.Metadata.Name = namespaceId
	bytes,err := yaml.Marshal(obj)
	if err != nil {
		logger.Errorf("CreateNamespace: YAML obj marshal, err: %v", err)
		return nil,fmt.Errorf("CreateNamespace: YAML obj marshal, err: %v", err)
	}
	reqUrl := cluster.Addr + k8s.API_V1 + k8s.NAMESPACES
	bytes,err = utils.RequestWithCertAndBody(reqUrl,utils.REQ_POST,cluster.Cert,cluster.Key,string(bytes))
	if err != nil { 
		logger.Errorf("CreateNamespace: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("CreateNamespace: RequestWithCertAndBody err,%v", err)
	}
	var result NamespaceDeployResult 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("CreateNamespace: create result err,%v", err)
		return nil,fmt.Errorf("CreateNamespace: create result err,%v", err)
	}
	return nil,nil
}

func DeleteNamespace(clusterId string,namespaceId string)([]byte,error) {
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("DeleteNamespace: cluser nil,cluser id =%s",clusterId)
		return nil,fmt.Errorf("DeleteNamespace: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_V1 + k8s.NAMESPACES + "/" + namespaceId
	bytes,err := utils.RequestWithCert(reqUrl,utils.REQ_DELETE,cluster.Cert,cluster.Key)
	if err != nil { 
		logger.Errorf("DeleteNamespace: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("DeleteNamespace: RequestWithCertAndBody err,%v", err)
	}
	var result NamespaceDeployResult 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("DeleteNamespace: delete result err,%v", err)
		return nil,fmt.Errorf("DeleteNamespace: delete result err,%v", err)
	}
	return nil,nil
}


