##wingbaas API      

## all wingbaas api return struct as follow       
```json 
{
    "code": 0, //success:0, error: other code    
    "message": "success", //message   
    "data": json obj //data object       
}
``` 
 
## add already exsist kuberntes cluster into wingbaas                
URL：http://ip:port//api/v1/addcluster      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/addcluster   
request parameter:       
```json
{
  "ClusterId": "kubernetes-cluster1",    
  "Addr": "https://kubernetes:6443/api/v1/",      
  "Cert": "./conf/pki/cluster1/apiserver-kubelet-client.crt",      
  "Key": "./conf/pki/cluster1/apiserver-kubelet-client.key"       
}    
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "cluster add success",
    "data": null
}
```    

## get cluster list                    
URL：http://ip:port/api/v1/clusters              
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/clusters              
API RETURN：                      
```json     
{
    "code": 0,     
    "message": "success",      
    "data": [      
        {
            "ClusterId": "kubernetes-cluster1",     
            "Addr": "https://kubernetes:6443/api/v1/",     
            "Cert": "./conf/pki/cluster1/apiserver-kubelet-client.crt",     
            "Key": "./conf/pki/cluster1/apiserver-kubelet-client.key"     
        }    
    ]      
}    
```    

## get blockchain type support deploy                                         
URL：http://ip:port/api/v1/blockchaintypes                                                
METHOD：GET    
INPUT PARA:                          
RETURN：json obj       
example:        
request:http://localhost:9001/api/v1/blockchaintypes                                                    
API RETURN：                         
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "fabric": {
            "version": {
                "1.3.0": {
                    "baseos": "registry.wingbaas.cn:5000/fabric-baseos:amd64-0.4.13",
                    "ca": "registry.wingbaas.cn:5000/fabric-ca:1.3.0",
                    "ccenv": "registry.wingbaas.cn:5000/fabric-ccenv:1.3.0",
                    "kafka": "registry.wingbaas.cn:5000/fabric-kafka:0.4.10",
                    "orderer": "registry.wingbaas.cn:5000/fabric-orderer:1.3.0",
                    "peer": "registry.wingbaas.cn:5000/fabric-peer:1.3.0",
                    "zookeeper": "registry.wingbaas.cn:5000/fabric-zookeeper:0.4.10"
                }
            }
        },
        "wingchain": {
            "version": {
                "1.0.0": {}
            }
        }
    }
}
```         

## get namespaces in kubernetes cluster                               
URL：http://ip:port/api/v1/:clusterid/namespaces                             
METHOD：GET    
INPUT PARA: clusterid                   
RETURN：json obj       
example:        
request:http://localhost:9001/api/v1/kubernetes-cluster1/namespaces                                 
API RETURN：                         
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "metadata": {
                "name": "default",
                "selfLink": "/api/v1/namespaces/default",
                "uid": "2de9e8a2-7780-49fb-b5d8-6a7f5b6a5360",
                "resourceVersion": "152",
                "creationTimestamp": "2020-03-26T09:07:03Z"
            },
            "spec": {
                "finalizers": [
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        }
    ]
}
```         
## deploy blockchain network                 
URL：http://ip:port//api/v1/deploy      
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/deploy           
request parameter:       
```json
{
	"BlockChainName": "test-chainnetwork11",  
	"BlockChainType": "fabric", 
	"DeployCfg":{ 
		"DeployNetCfg": { 
			"OrdererOrgs": [
			{
				"Name": "Orderer",
				"Domain": "orderer.baas.xyz",
				"DeployNode": "172-16-254-33",
				"Specs": [
					{
						"Hostname": "orderer0"
					},
					{
						"Hostname": "orderer1"
					}
				]
			}
		],
			"PeerOrgs": [
			{
				"Name": "Org1",
				"Domain": "Org1.fabric.baas.xyz",
				"DeployNode": "172-16-254-130",
				"Specs": [
					{
						"Hostname": "peer0-org1"
					},
					{
						"Hostname": "peer1-org1"
					}
				],
				"Users": {
					"Count": 1
				}
			},
			{
				"Name": "Org2",
				"Domain": "Org2.fabric.baas.xyz",
				"DeployNode": "172-16-254-33",
				"Specs": [
					{
						"Hostname": "peer0-org2"
					},
					{
						"Hostname": "peer1-org2"
					}
				],
				"Users": {
					"Count": 1
				}
			}
		],
		"KafkaDeployNode": "172-16-254-33",
		"ZookeeperDeployNode": "172-16-254-130"
		},
		"DeployType": "KAFKA_FABRIC",
		"Version": "1.3.0",
		"CryptoType": "ECDSA",
		"ClusterId": "test-cluster1"
	}
}
``` 
API RETURN:                  
```json     
{
    "code": 0,       
    "message": "success",              
    "data": {
        "BlockChainId": "Ys55rfOqDUJpU5QAAN9xvk18JpP0gmpY", //blockchain id      
		"BlockChainName": "chainnetwork1",         
		"ClusterId": "cluster1"         
    }         
}
```      

## get blockchains deploy in kubernetes cluster                                 
URL：http://ip:port/api/v1/:clusterid/blockchains                                       
METHOD：GET    
INPUT PARA: clusterid                   
RETURN：json obj       
example:        
request:http://localhost:9001/api/v1/cluster1/blockchains                                           
API RETURN：                         
```json     
{
    "code": 0,      
    "message": "success",       
    "data": [       
        {
            "BlockChainId": "GAZrm1ppCH7oZqRsRlPDZjSVoF9cvbbH",       
            "BlockChainName": "chainnetwork2",          
            "BlockChainType": "fabric",          
            "ClusterId": "cluster1"       
        }      
    ]      
}      
```         

## get host list in kubernetes cluster                                 
URL：http://ip:port/api/v1/:clusterid/hosts                                                 
METHOD：GET       
INPUT PARA: clusterid                        
RETURN：json obj       
example:        
request:http://localhost:9001/api/v1/cluster1/hosts                                                     
API RETURN：                         
```json     
 "message": "success",
    "data": {
        "apiVersion": "v1",
        "kind": "NodeList",
        "metadata": {
            "resourceVersion": "2365800",
            "selfLink": "/api/v1/nodes"
        },
        "items": [
            {
                "metadata": {
                    "name": "172-16-254-130",
                    "resourceVersion": "2365660",
                    "selfLink": "/api/v1/nodes/172-16-254-130",
                    "uid": "68db3792-e499-42b2-9250-b501aa852601",
                    "creationTimestamp": "2020-03-26T09:45:20Z",
                    "annotations": {
                        "flannel.alpha.coreos.com/backend-data": "{\"VtepMAC\":\"4e:bd:a5:26:ec:3f\"}",
                        "flannel.alpha.coreos.com/backend-type": "vxlan",
                        "flannel.alpha.coreos.com/kube-subnet-manager": "true",
                        "flannel.alpha.coreos.com/public-ip": "172.16.254.130",
                        "kubeadm.alpha.kubernetes.io/cri-socket": "/var/run/dockershim.sock",
                        "node.alpha.kubernetes.io/ttl": "0",
                        "volumes.kubernetes.io/controller-managed-attach-detach": "true"
                    },
                    "labels": {
                        "beta.kubernetes.io/arch": "amd64",
                        "beta.kubernetes.io/os": "linux",
                        "kubernetes.io/arch": "amd64",
                        "kubernetes.io/hostname": "172-16-254-130",
                        "kubernetes.io/os": "linux",
                        "node-role.kubernetes.io/master": ""
                    }
                },
                "spec": {
                    "podCIDR": "10.244.1.0/24",
                    "podCIDRs": [
                        "10.244.1.0/24"
                    ],
                    "taints": null
                },
                "status": {
                    "addresses": [
                        {
                            "address": "172.16.254.130",
                            "type": "InternalIP"
                        },
                        {
                            "address": "172-16-254-130",
                            "type": "Hostname"
                        }
                    ],
                    "allocatable": {
                        "cpu": "4",
                        "ephemeral-storage": "19316971898",
                        "hugepages-1Gi": "0",
                        "hugepages-2Mi": "0",
                        "memory": "7856632Ki",
                        "pods": "110"
                    },
                    "capacity": {
                        "cpu": "4",
                        "ephemeral-storage": "20469Mi",
                        "hugepages-1Gi": "0",
                        "hugepages-2Mi": "0",
                        "memory": "7959032Ki",
                        "pods": "110"
                    },
                    "conditions": [
                        {
                            "lastHeartbeatTime": "2020-03-26T09:53:46Z",
                            "lastTransitionTime": "2020-03-26T09:53:46Z",
                            "message": "Flannel is running on this node",
                            "reason": "FlannelIsUp",
                            "status": "False",
                            "type": "NetworkUnavailable"
                        },
                        {
                            "lastHeartbeatTime": "2020-04-07T06:37:04Z",
                            "lastTransitionTime": "2020-03-26T09:45:20Z",
                            "message": "kubelet has sufficient memory available",
                            "reason": "KubeletHasSufficientMemory",
                            "status": "False",
                            "type": "MemoryPressure"
                        },
                        {
                            "lastHeartbeatTime": "2020-04-07T06:37:04Z",
                            "lastTransitionTime": "2020-03-26T09:45:20Z",
                            "message": "kubelet has no disk pressure",
                            "reason": "KubeletHasNoDiskPressure",
                            "status": "False",
                            "type": "DiskPressure"
                        },
                        {
                            "lastHeartbeatTime": "2020-04-07T06:37:04Z",
                            "lastTransitionTime": "2020-03-26T09:45:20Z",
                            "message": "kubelet has sufficient PID available",
                            "reason": "KubeletHasSufficientPID",
                            "status": "False",
                            "type": "PIDPressure"
                        },
                        {
                            "lastHeartbeatTime": "2020-04-07T06:37:04Z",
                            "lastTransitionTime": "2020-03-26T09:53:50Z",
                            "message": "kubelet is posting ready status",
                            "reason": "KubeletReady",
                            "status": "True",
                            "type": "Ready"
                        }
                    ],
                    "daemonEndpoints": {
                        "kubeletEndpoint": {
                            "Port": 10250
                        }
                    },
                    "images": [
                        {
                            "names": [
                                "registry.cn-hangzhou.aliyuncs.com/google_containers/kube-proxy@sha256:32cb1dbf68cc45401c73a92c13da3cd283797f63ffc3ab3e614963c6c9929e4a",
                                "k8s.gcr.io/kube-proxy:v1.17.4",
                                "registry.cn-hangzhou.aliyuncs.com/google_containers/kube-proxy:v1.17.4"
                            ],
                            "sizeBytes": 115964919
                        },
                        {
                            "names": [
                                "docker.io/kubernetesui/dashboard@sha256:fc90baec4fb62b809051a3227e71266c0427240685139bbd5673282715924ea7",
                                "docker.io/kubernetesui/dashboard:v2.0.0-beta8"
                            ],
                            "sizeBytes": 90835427
                        },
                        {
                            "names": [
                                "quay-mirror.qiniu.com/coreos/flannel@sha256:6d451d92c921f14bfb38196aacb6e506d4593c5b3c9d40a8b8a2506010dc3e10",
                                "quay.io/coreos/flannel@sha256:6d451d92c921f14bfb38196aacb6e506d4593c5b3c9d40a8b8a2506010dc3e10",
                                "quay-mirror.qiniu.com/coreos/flannel:v0.12.0-amd64",
                                "quay.io/coreos/flannel:v0.12.0-amd64"
                            ],
                            "sizeBytes": 52767393
                        },
                        {
                            "names": [
                                "quay.io/coreos/flannel@sha256:88f2b4d96fae34bfff3d46293f7f18d1f9f3ca026b4a4d288f28347fcb6580ac",
                                "quay.io/coreos/flannel:v0.10.0-amd64"
                            ],
                            "sizeBytes": 44598861
                        },
                        {
                            "names": [
                                "docker.io/kubernetesui/metrics-scraper@sha256:35fcae4fd9232a541a8cb08f2853117ba7231750b75c2cb3b6a58a2aaa57f878",
                                "docker.io/kubernetesui/metrics-scraper:v1.0.1"
                            ],
                            "sizeBytes": 40101504
                        },
                        {
                            "names": [
                                "docker.io/registry@sha256:7d081088e4bfd632a88e3f3bcd9e007ef44a796fddfe3261407a3f9f04abe1e7",
                                "docker.io/registry:latest"
                            ],
                            "sizeBytes": 25769582
                        },
                        {
                            "names": [
                                "registry.cn-hangzhou.aliyuncs.com/google_containers/pause@sha256:759c3f0f6493093a9043cc813092290af69029699ade0e3dbe024e968fcb7cca",
                                "k8s.gcr.io/pause:3.1",
                                "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.1"
                            ],
                            "sizeBytes": 742472
                        }
                    ],
                    "nodeInfo": {
                        "architecture": "amd64",
                        "bootID": "33c81c88-9222-43f8-890c-f9cff80958db",
                        "containerRuntimeVersion": "docker://1.13.1",
                        "kernelVersion": "4.19.0-6.el7.ucloud.x86_64",
                        "kubeProxyVersion": "v1.17.4",
                        "kubeletVersion": "v1.17.4",
                        "machineID": "32cccb85ec1f07509239bd8287881d1e",
                        "operatingSystem": "linux",
                        "osImage": "CentOS Linux 7 (Core)",
                        "systemUUID": "9264e13a-a63e-450c-bbda-8137fb821827"
                    }
                }
            }
        ]
    }
}
```       

## delete blockchain network                     
URL：http://ip:port//api/v1/delete                 
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/delete                  
request parameter:            
```json    
{    
	"BlockChainName": "test-chainnetwork11",         
	"ClusterId": "test-cluster1"     
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": null    
}
```        
