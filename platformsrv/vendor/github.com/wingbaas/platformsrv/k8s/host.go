
package k8s

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
)

const (
	NODES    string = "nodes"
)

type HostList struct {
	APIVersion string `json:"apiVersion"`
	Kind     string `json:"kind"`
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
	Items      []struct {
		Metadata struct {
			Name            string `json:"name"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			UID             string `json:"uid"`
			CreationTimestamp string `json:"creationTimestamp"`
			Annotations struct {
				Flannel_alpha_coreos_com_backend_data                  string `json:"flannel.alpha.coreos.com/backend-data"`
				Flannel_alpha_coreos_com_backend_type                  string `json:"flannel.alpha.coreos.com/backend-type"`
				Flannel_alpha_coreos_com_kube_subnet_manager           string `json:"flannel.alpha.coreos.com/kube-subnet-manager"`
				Flannel_alpha_coreos_com_public_ip                     string `json:"flannel.alpha.coreos.com/public-ip"`
				Kubeadm_alpha_kubernetes_io_cri_socket                 string `json:"kubeadm.alpha.kubernetes.io/cri-socket"`
				Node_alpha_kubernetes_io_ttl                           string `json:"node.alpha.kubernetes.io/ttl"`
				Volumes_kubernetes_io_controller_managed_attach_detach string `json:"volumes.kubernetes.io/controller-managed-attach-detach"`
			} `json:"annotations"`
			Labels            struct {
				Beta_kubernetes_io_arch        string `json:"beta.kubernetes.io/arch"`
				Beta_kubernetes_io_os          string `json:"beta.kubernetes.io/os"`
				Kubernetes_io_arch             string `json:"kubernetes.io/arch"`
				Kubernetes_io_hostname         string `json:"kubernetes.io/hostname"`
				Kubernetes_io_os               string `json:"kubernetes.io/os"`
				Node_role_kubernetes_io_master string `json:"node-role.kubernetes.io/master"`
			} `json:"labels"`
		} `json:"metadata"`
		Spec struct {
			PodCIDR  string   `json:"podCIDR"`
			PodCIDRs []string `json:"podCIDRs"`
			Taints   []struct {
				Effect string `json:"effect"`
				Key    string `json:"key"`
			} `json:"taints"`
		} `json:"spec"`
		Status struct {
			Addresses []struct {
				Address string `json:"address"`
				Type    string `json:"type"`
			} `json:"addresses"`
			Allocatable struct {
				CPU               string `json:"cpu"`
				Ephemeral_storage string `json:"ephemeral-storage"`
				Hugepages_1Gi     string `json:"hugepages-1Gi"`
				Hugepages_2Mi     string `json:"hugepages-2Mi"`
				Memory            string `json:"memory"`
				Pods              string `json:"pods"`
			} `json:"allocatable"`
			Capacity struct {
				CPU               string `json:"cpu"`
				Ephemeral_storage string `json:"ephemeral-storage"`
				Hugepages_1Gi     string `json:"hugepages-1Gi"`
				Hugepages_2Mi     string `json:"hugepages-2Mi"`
				Memory            string `json:"memory"`
				Pods              string `json:"pods"`
			} `json:"capacity"`
			Conditions []struct {
				LastHeartbeatTime  string `json:"lastHeartbeatTime"`
				LastTransitionTime string `json:"lastTransitionTime"`
				Message            string `json:"message"`
				Reason             string `json:"reason"`
				Status             string `json:"status"`
				Type               string `json:"type"`
			} `json:"conditions"`
			DaemonEndpoints struct {
				KubeletEndpoint struct {
					Port int64 `json:"Port"`
				} `json:"kubeletEndpoint"`
			} `json:"daemonEndpoints"`
			Images []struct {
				Names     []string `json:"names"`
				SizeBytes int64    `json:"sizeBytes"`
			} `json:"images"`
			NodeInfo struct {
				Architecture            string `json:"architecture"`
				BootID                  string `json:"bootID"`
				ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
				KernelVersion           string `json:"kernelVersion"`
				KubeProxyVersion        string `json:"kubeProxyVersion"`
				KubeletVersion          string `json:"kubeletVersion"`
				MachineID               string `json:"machineID"`
				OperatingSystem         string `json:"operatingSystem"`
				OsImage                 string `json:"osImage"`
				SystemUUID              string `json:"systemUUID"`
			} `json:"nodeInfo"`
		} `json:"status"`
	} `json:"items"`
}


 
func GetHostList(clusterId string)(*HostList,error) {
	logger.Debug("GetHostList")
	cluster,err := GetCluster(clusterId) 
	if err != nil {
		logger.Errorf("GetHostList: get clusters error,%v", err)
		return nil,fmt.Errorf("%v", err)
	}
	if cluster == nil {
		logger.Errorf("GetHostList: cluser nil,cluser id = %s",clusterId)
		return nil,fmt.Errorf("GetHostList: cluser nil,cluser id = %s",clusterId)
	}
	bytes,err := utils.RequestWithCert(cluster.Addr + API_V1 + NODES,utils.REQ_GET,cluster.Cert,cluster.Key)
	if err != nil { 
		logger.Errorf("GetHostList: RequestWithCert err,%v", err)
		return nil,nil
	}
	var hostList HostList
	err = json.Unmarshal(bytes, &hostList)
	if err != nil {
		logger.Errorf("GetHostList: unmarshal host list err,%v", err)
		return nil,fmt.Errorf("GetHostList: unmarshal host list err,%v", err)
	}
	return &hostList,nil 
}
