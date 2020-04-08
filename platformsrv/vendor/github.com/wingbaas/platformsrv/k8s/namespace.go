
package k8s
import (
	"fmt"
	"encoding/json"
	"time"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
)

const (
	NAMESPACES    string = "namespaces"
	DEPLOYMENTS   string = "deployments"
	SERVICES      string = "services"
	API_V1        string = "/api/v1/" 
	API_APP       string = "/apis/apps/v1/"   
)

type Namespace struct {
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

type Namespaces struct {
	Kind string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
		SelfLink string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []Namespace `json:"items"`
}

func GetClusterNamespaces(cluserId string) ([]Namespace,error) {
	cluster,_ := GetCluster(cluserId)
	if cluster == nil {
		logger.Errorf("GetClusterNamespaces: cluser nil,cluser id =%s",cluserId)
		return nil,fmt.Errorf("GetClusterNamespaces: cluser nil,cluser id =%s",cluserId)
	}
	bytes,err := utils.RequestWithCert(cluster.Addr + API_V1 + NAMESPACES,utils.REQ_GET,cluster.Cert,cluster.Key)
	if err != nil { 
		logger.Errorf("GetClusterNamespaces: RequestWithCert err,%v", err)
		return nil,nil
	}
	var namespaces Namespaces
	err = json.Unmarshal(bytes, &namespaces)
	if err != nil {
		logger.Errorf("GetClusterNamespaces: unmarshal namespaces err,%v", err)
		return nil,fmt.Errorf("GetClusterNamespaces: unmarshal namespaces err,%v", err)
	}
	return namespaces.Items,nil
}

