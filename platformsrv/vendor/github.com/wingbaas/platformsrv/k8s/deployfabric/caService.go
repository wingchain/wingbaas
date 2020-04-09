
package deployfabric
import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/certgenerate/fabric"
)

type MetadataCaService struct {
	Name string `json:"name"`
}

type PortSpecCaService struct {
	Name string `json:"name"`
	Port int `json:"port"`
	TargetPort int `json:"targetPort"`
}

type SelectorSpecCaService struct {
	App string `json:"app"`
}

type SpecCaService struct {
	Type string `json:"type"`
	Ports []PortSpecCaService `json:"ports"`
	Selector SelectorSpecCaService `json:"selector"`
} 

type CaService struct {
	APIVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Metadata MetadataCaService `json:"metadata"`
	Spec SpecCaService `json:"spec"`
}

func CreateCaService(clusterId string,namespaceId string,caName string)([]byte,error) {
	caService := CaService{
		APIVersion: "v1",
		Kind: "Service",
		Metadata: MetadataCaService{
			Name: caName,
		},
		Spec: SpecCaService{
			Type: "NodePort",
			Ports: []PortSpecCaService{
				{
					Name: caName,
					Port: 7054,
					TargetPort: 7054,
				},
			},
			Selector: SelectorSpecCaService{
				App: caName,
			},
		},
	}
	bytes, err := json.Marshal(caService)
	if err != nil {
		logger.Errorf("CreateCaService: Marshal service error,%v",err) 
		return nil,fmt.Errorf("CreateCaService: Marshal service error,%v",err)
	}
	bl,yamlStr := fabric.JsonToYaml(string(bytes))
	if !bl {
		logger.Errorf("CreateCaService: json2yaml str error") 
		return nil,fmt.Errorf("CreateCaService: json2yaml str error")
	}
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("CreateCaService: cluster nil,cluster id =%s",clusterId)
		return nil,fmt.Errorf("CreateCaService: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_V1 + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.SERVICES
	bytes,err = utils.RequestWithCertAndBody(reqUrl,utils.REQ_POST,cluster.Cert,cluster.Key,yamlStr)
	if err != nil { 
		logger.Errorf("CreateCaService: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("CreateCaService: RequestWithCertAndBody err,%v", err)
	}
	var result interface{} 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("CreateCaService: create result err,%v", err)
		return nil,fmt.Errorf("CreateCaService: create result err,%v", err)
	}
	logger.Debug("CreateCaService: create success,result=,%v", result) 
	logger.Debug("CreateCaService: create success,yamlstr=")
	logger.Debug(yamlStr)  
	return nil,nil
}

func DeleteCaService(clusterId string,namespaceId string,serviceName string)([]byte,error) {
	cluster,_ := k8s.GetCluster(clusterId)
	if cluster == nil {
		logger.Errorf("DeleteCaService: cluser nil,cluser id =%s",clusterId)
		return nil,fmt.Errorf("DeleteCaService: cluster nil,cluster id =%s",clusterId)
	}
	reqUrl := cluster.Addr + k8s.API_V1 + k8s.NAMESPACES + "/" + namespaceId + "/" + k8s.SERVICES + "/" + serviceName
	bytes,err := utils.RequestWithCert(reqUrl,utils.REQ_DELETE,cluster.Cert,cluster.Key)
	if err != nil { 
		logger.Errorf("DeleteCaService: RequestWithCertAndBody err,%v", err)
		return nil,fmt.Errorf("DeleteCaService: RequestWithCertAndBody err,%v", err)
	}
	var result interface{} 
	err = json.Unmarshal(bytes, &result)
	if err != nil { 
		logger.Errorf("DeleteCaService: delete result err,%v", err)
		return nil,fmt.Errorf("DeleteCaService: delete result err,%v", err)
	}
	logger.Debug("DeleteCaService: delete success,result=,%v", result) 
	return nil,nil
}


