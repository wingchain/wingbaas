
package k8s
import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
)

type Service struct {
	Metadata struct {
		Name string `json:"name"`
		Namespace string `json:"namespace"`
		SelfLink string `json:"selfLink"`
		UID string `json:"uid"`
		ResourceVersion string `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Spec struct {
		Ports []struct {
			Name string `json:"name"`
			Protocol string `json:"protocol"`
			Port int `json:"port"`
			TargetPort int `json:"targetPort"`
			NodePort int `json:"nodePort"`
		} `json:"ports"`
		Selector struct {
			App string `json:"app"`
		} `json:"selector"`
		ClusterIP string `json:"clusterIP"`
		Type string `json:"type"`
		SessionAffinity string `json:"sessionAffinity"`
		ExternalTrafficPolicy string `json:"externalTrafficPolicy"`
	} `json:"spec"`
	Status struct {
		LoadBalancer struct {
		} `json:"loadBalancer"`
	} `json:"status"`
}

type Services struct {
	Kind string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
		SelfLink string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []Service `json:"items"`
}

func GetNamespaceServices(clusterId string,namespaceId string) ([]Service,error) {
	cluster,_ := GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("GetNamespaceServices: cluster nil,cluster id =%s",clusterId)
		return nil,fmt.Errorf("GetNamespaceServices: cluster nil,cluster id =%s",clusterId)
	}
	cp := utils.BAAS_CFG.ClusterPkiBasePath
	bytes,err := utils.RequestWithCert(cluster.ApiUrl + API_V1 + NAMESPACES + "/" + namespaceId + "/" + SERVICES,utils.REQ_GET,cp + cluster.Cert,cp + cluster.Key)
	if err != nil { 
		logger.Errorf("GetNamespaceServices: RequestWithCert err,%v", err)
		return nil,nil
	}
	var services Services
	err = json.Unmarshal(bytes, &services)
	if err != nil {
		logger.Errorf("GetNamespaceServices: unmarshal services err,%v", err)
		return nil,fmt.Errorf("GetNamespaceServices: unmarshal services err,%v", err)
	}
	if services.Items == nil || len(services.Items) < 1{
		logger.Errorf("GetNamespaceServices: service items nil")
		return nil,fmt.Errorf("GetNamespaceServices: service items nil")
	}
	// logger.Debug("service list=")
	// logger.Debug(services.Items)
	return services.Items,nil
} 

func getNodePortByName(s []Service,name string)(error,[]public.ServiceNodePortSt) {
	var np []public.ServiceNodePortSt
	for _,ss := range s {
		if ss.Metadata.Name == name {
			for _,p := range ss.Spec.Ports {
				var obj public.ServiceNodePortSt
				obj.ServerName = p.Name
				obj.NodePort = strconv.Itoa(p.NodePort)
				np = append(np,obj)
			}
			if len(np) <= 0 {
				logger.Errorf("getNodePortByName: ports is null,name=%s",name)
				return fmt.Errorf("getNodePortByName: ports is null,name=%s",name),nil
			}
			return nil,np
		}
	}
	logger.Errorf("getNodePortByName: not find,name=%s",name)
	return fmt.Errorf("getNodePortByName: not find,name=%s",name),nil
}

func GetServicesNodePort(clusterId string,namespaceId string,netCfg public.DeployNetConfig) (error,map[string][]public.ServiceNodePortSt) {
	s,_ := GetNamespaceServices(clusterId,namespaceId)
	if s == nil {
		logger.Errorf("GetServicesNodePort: GetNamespaceServices error")
		return fmt.Errorf("GetServicesNodePort: GetNamespaceServices error"),nil
	}
	m := make(map[string][]public.ServiceNodePortSt)
	for _,org := range netCfg.OrdererOrgs {
		for _,h := range org.Specs {
			_,nps := getNodePortByName(s,h.Hostname)
			if nps != nil {
				m[h.Hostname] = nps
			}
		}
	}
	for _,org := range netCfg.PeerOrgs {
		for _,h := range org.Specs {
			key := strings.ToLower(h.Hostname)
			_,nps := getNodePortByName(s,key)
			if nps != nil {
				m[key] = nps
			}
		}
		caKey := strings.ToLower(org.Name + "-ca")
		_,caps := getNodePortByName(s,caKey)
		if caps != nil {
			m[caKey] = caps
		}
	}
	// logger.Debug("GetServicesNodePort ret map=")
	// logger.Debug(m)
	return nil,m
}

func GetNodePort(m map[string][]public.ServiceNodePortSt,firstKey string,secondKey string) string {
	obj := m[firstKey]
	for _,po := range obj {
		if po.ServerName == secondKey {
			return po.NodePort 
		}
	}
	logger.Debug("GetNodePort: not find,firstKey=%s,secondKey=%s",firstKey,secondKey)
	return ""
}
