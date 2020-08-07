##wingbaas API      

## all wingbaas api return struct as follow       
```json 
{
    "code": 0, //success:0, error: other code    
    "message": "success", //message   
    "data": json obj //data object       
}
``` 
## user register                    
URL：http://ip:port/api/v1/register      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/register      
request parameter:       
```json
{
	"Mail": "testuser@163.com", 
    "Password": "password123456", 
    "VerifyCode": "123456"
}
``` 
API RETURN:                  
{
    "code": 0,
    "message": "success", 
    "data": null
}
```    

## user login                    
URL：http://ip:port/api/v1/login      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/login      
request parameter:       
```json
{
	"Mail": "testuser@163.com", 
    "Password": "password123456"
}
``` 
API RETURN:                  
{
    "code": 0,
    "message": "success", 
    "data": null
}
```    

## create alliance                   
URL：http://ip:port/api/v1/createalliance      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/createalliance      
request parameter:       
```json
{
	"Name": "食品联盟",
	"Describe": "食品行业联盟",
	"Creator": "testuser@163.com"
}
``` 
API RETURN:                  
{
    "code": 0,
    "message": "success",
    "data": "5bGiWNF0ycI7eTbH" //alliance id     
}
```    
## get alliance list                    
URL：http://ip:port/api/v1/alliances              
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/alliances              
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "Id": "5bGiWNF0ycI7eTbH",
            "Name": "食品联盟",
            "Describe": "食品行业联盟",
            "Creator": "testuser@163.com"
        }
    ]
}
```    

## get alliance by id                    
URL：http://ip:port/api/v1/:allianceid/alliance              
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/5bGiWNF0ycI7eTbH/alliance                 
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "Id": "5bGiWNF0ycI7eTbH",
        "Name": "食品联盟",
        "Describe": "食品行业联盟",
        "Creator": "testuser@163.com"
    }
}
```    
## user add to alliance                   
URL：http://ip:port/api/v1/useraddalliance      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/useraddalliance      
request parameter:       
```json
{
	"Mail": "testuser2@163.com",
	"Alliance":{
		 "Id": "5bGiWNF0ycI7eTbH",
         "Name": "食品联盟",
         "Describe": "食品行业联盟",
         "Creator": "testuser@163.com" 
	}
}
``` 
API RETURN:                  
{
    "code": 0,
    "message": "success",
    "data": null
}
```    

## get user joined alliances                      
URL：http://ip:port/api/v1/:usermail/joinedalliances              
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/testuser1@163.com/joinedalliances                 
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "Id": "dKKW2DZlbAnd4n5x",
            "Name": "食品联盟",
            "Describe": "食品行业联盟",
            "Creator": "testuser@163.com"
        }
    ]
}
```    

## get user created alliances                      
URL：http://ip:port/api/v1/:usermail/createdalliances                 
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/testuser@163.com/createdalliances                 
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "Id": "dKKW2DZlbAnd4n5x",
            "Name": "食品联盟",
            "Describe": "食品行业联盟",
            "Creator": "testuser@163.com"
        }
    ]
}
```    

## get users by alliance id                           
URL：http://ip:port/api/v1/:allianceid/users                 
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/FEWnQZAVNOBvQYRK/users                       
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "Mail": "testuser@163.com",
            "Password": "password123456",
            "VerifyCode": "123456",
            "Alliances": [
                {
                    "Id": "FEWnQZAVNOBvQYRK",
                    "Name": "食品联盟",
                    "Describe": "食品行业联盟",
                    "Creator": "testuser@163.com"
                }
            ]
        },
        {
            "Mail": "testuser1@163.com",
            "Password": "password123456",
            "VerifyCode": "123456",
            "Alliances": [
                {
                    "Id": "FEWnQZAVNOBvQYRK",
                    "Name": "食品联盟",
                    "Describe": "食品行业联盟",
                    "Creator": "testuser@163.com"
                }
            ]
        },
        {
            "Mail": "testuser2@163.com",
            "Password": "password123456",
            "VerifyCode": "123456",
            "Alliances": [
                {
                    "Id": "FEWnQZAVNOBvQYRK",
                    "Name": "食品联盟",
                    "Describe": "食品行业联盟",
                    "Creator": "testuser@163.com"
                }
            ]
        }
    ]
}
```    

## user delete alliance                   
URL：http://ip:port/api/v1/deleteAlliance      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/deleteAlliance      
request parameter:       
```json
{
	"Mail": "testuser2@163.com",
	"AllianceId": "kdikdkkdik"
}
``` 
API RETURN:                  
{
    "code": 0,
    "message": "success",
    "data": null
}
```    

## delete alliance user                       
URL：http://ip:port/api/v1/deleteallianceuser       
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/deleteallianceuser         
request parameter:       
```json
{
	"Mail": "testuser2@163.com",
	"AllianceId": "kdikdkkdik"
}
``` 
API RETURN:                  
{
    "code": 0,
    "message": "success",
    "data": null
}
```    

## upload already exsist kuberntes cluster cert or key file into wingbaas                 
URL：http://ip:port/api/v1/upLoadKeyFile      
METHOD：POST(formdata)     
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/upLoadKeyFile          
request parameter:(form-data file)   
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",
    "data": "hNl4075WtUIvbAKz"//file id        
}
```      
 
## add already exsist kuberntes cluster into wingbaas                
URL：http://ip:port/api/v1/addcluster      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/addcluster   
request parameter:       
```json
{
	"AllianceId": "oZXdf4olZ3nMOXGA",
	"ClusterId": "test-cluster1", 
    "ApiUrl": "https://kubernetes:6443", 
    "HostDomain": "kubernetes",
	"InterIp": "172.16.254.33",
    "PublicIp": "106.75.51.138",
    "Cert": "hEU9HNb6Mb3JJ8fR",//file id
    "Key": "hNl4075WtUIvbAKz"//file id
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

## user delete cluster                      
URL：http://ip:port/api/v1/deletecluster      
METHOD：POST   
RETURN：json object           

example:        
request:http://localhost:9001/api/v1/deletecluster      
request parameter:       
```json
{
	"Creator": "testuser@163.com",
	"ClusterId": "kdikdkkdik"
}
``` 
API RETURN:                  
{
    "code": 0,
    "message": "success",
    "data": null
}
```    

## get alliance created clusters                       
URL：http://ip:port/api/v1/:allianceid/allianceclusters                 
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/oZXdf4olZ3nMOXGA/allianceclusters                      
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "AllianceId": "oZXdf4olZ3nMOXGA",
            "ClusterId": "test-cluster1",
            "ApiUrl": "https://kubernetes:6443",
            "HostDomain": "kubernetes",
            "PublicIp": "106.75.51.138",
            "Cert": "Eol658oR9Q5fHP8g",
            "Key": "00STWLq5DqReXukG"
        }
    ]
}
```    

## get alliance created blockchains                           
URL：http://ip:port/api/v1/:allianceid/alliancechains                      
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/oZXdf4olZ3nMOXGA/alliancechains                       
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "AllianceId": "QlVwTJjMA3r90TMA",
            "BlockChainId": "0dXFptQ0qEqsZg7eqzUp9dReCKLuIH1Q",
            "BlockChainName": "test-chainnetwork11",
            "BlockChainType": "fabric",
            "ClusterId": "test-cluster1",
            "Version": "1.3.0",
            "Status": 0
        }
    ]
}
```    

## get chain by chian id                       
URL：http://ip:port/api/v1/:blockchainid/blockchain              
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/NitnGgWzSs7zAX6ktXQRbVtCHaYnFDCc/blockchain                 
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "AllianceId": "QlVwTJjMA3r90TMA",
        "BlockChainId": "NitnGgWzSs7zAX6ktXQRbVtCHaYnFDCc",
        "BlockChainName": "test-chainnetwork11",
        "BlockChainType": "fabric",
        "ClusterId": "test-cluster1",
        "Version": "1.3.0",
        "Status": 101//0:创建成功,101:创建中,102:创建失败
    }
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
	"AllianceId": "QlVwTJjMA3r90TMA",
	"BlockChainName": "test-chainnetwork11",  
	"BlockChainType": "fabric", 
	"DeployCfg":{ 
		"DeployNetCfg": { 
			"OrdererOrgs": [
			{
				"Name": "Orderer",//const paramerter    
				"Domain": "orderer.baas.xyz",//const parameter    
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
						"Hostname": "peer0-org1" //must abide by this rule,contains host symbol,link character '-' and org name       
					},
					{
						"Hostname": "peer1-org1"//must abide by this rule,contains host symbol,link character '-' and org name    
					}
				],
				"Users": {
					"Count": 1//const parameter     
				}
			},
			{
				"Name": "Org2",
				"Domain": "Org2.fabric.baas.xyz",
				"DeployNode": "172-16-254-33",
				"Specs": [
					{
						"Hostname": "peer0-org2"//must abide by this rule,contains host symbol,link character '-' and org name     
					},
					{
						"Hostname": "peer1-org2"//must abide by this rule,contains host symbol,link character '-' and org name     
					}
				],
				"Users": {
					"Count": 1//const parameter    
				}
			}
		],
		"KafkaDeployNode": "172-16-254-33",//only DeployType is KAFKA_FABRIC need 
		"ZookeeperDeployNode": "172-16-254-130",//only DeployType is KAFKA_FABRIC need   
		"ToolsDeployNode": "172-16-254-33"//only DeployType is RAFT_FABRIC need    
		},
		"DeployType": "KAFKA_FABRIC",//KAFKA_FABRIC or RAFT_FABRIC         
		"Version": "1.3.0",//if 1.x, parameter DeployType only KAFKA_FABRIC, 2.x, parameter DeployType only RAFT_FABRIC         
		"CryptoType": "ECDSA",//const parameter       
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
URL：http://ip:port/api/v1/delete                 
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

## org create channel in blockchain                            
URL：http://ip:port/api/v1/orgcreatechannel                     
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/orgcreatechannel                           
request parameter:            
```json    
{
	"BlockChainId": "KSpAvrMQhZHjoTDX8zqsv027qcuNZLHm",
	"OrgName": "Org1",
	"ChannelId": "mychannel" 
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

## org join into blockchain                           
URL：http://ip:port/api/v1/orgjoinchannel                    
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/orgjoinchannel                         
request parameter:            
```json    
{
	"BlockChainId": "KSpAvrMQhZHjoTDX8zqsv027qcuNZLHm",
	"OrgName": "Org1",
	"ChannelId": "mychannel" 
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

## upload chaincode src code into blockchain                           
URL：http://ip:port/api/v1/uploadcc                    
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/uploadcc                                 
request parameter(form-data):            
```json    
{
	"BlockChainId": "KSpAvrMQhZHjoTDX8zqsv027qcuNZLHm",
	"ChainCodeId": "cctest",
    "ChainCodeVersion": "1.0",
    "file": file
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

## deploy already uploaded chaincode blockchain                         
URL：http://ip:port/api/v1/orgdeploycc                         
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/orgdeploycc                             
request parameter:            
```json    
{
	"BlockChainId": "KSpAvrMQhZHjoTDX8zqsv027qcuNZLHm",
	"OrgName": "Org1",
	"ChannelId": "mychannel",
	"ChainCodeId": "cctest",
	"ChainCodeVersion": "1.0",
	"EndorsePolicy": "AND('Org1MSP.member','Org2MSP.member')",
	"InitArgs":["init","a","200","b","200"]
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

## upgrade chaincode                              
URL：http://ip:port/api/v1/orgupgradecc                         
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/orgupgradecc                             
request parameter:            
```json    
{
	"BlockChainId": "FVCZFbOddroBZbKGvsCgzX4yktITr8kS",
	"OrgName": "Org1",
	"ChannelId": "baaschannel",
	"ChainCodeId": "cctest",
	"ChainCodeVersion": "100", 
	"ChaincodeSeq": "2",
	"EndorsePolicy": "AND('Org1MSP.member','Org2MSP.member','Org3MSP.member')",
	"InitArgs":["init","a1","300","b1","300"]
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

## call chaincode                             
URL：http://ip:port/api/v1/callcc                             
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/callcc                                       
request parameter:            
```json    
{
	"BlockChainId": "KSpAvrMQhZHjoTDX8zqsv027qcuNZLHm",
	"OrgName": "Org1",
	"ChannelId": "mychannel",
	"ChainCodeId": "cctest",
	"Args":["transfer","a","b","10"],
	"Peers":["peer0-org1.Org1.fabric.baas.xyz","peer0-org2.Org2.fabric.baas.xyz"]
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": {
        "TransactionID": "4cb2d7db491f8f9c517e89e1332f539f4d5d8c335d6def0a4c18facb32a53c23",
        "TxValidationCode": 0,
        "ChaincodeStatus": 200
    }
}
```       

## query chaincode                                  
URL：http://ip:port/api/v1/querycc                                 
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/querycc                                           
request parameter:            
```json    
{
	"BlockChainId": "6f2Wq0Hy5LSYpSt0yvmTd1XD2XEn3Hdr",
	"OrgName": "Org1",
	"ChannelId": "mychannel" ,
	"ChainCodeId": "cctest",
	"Args": ["query","a"]
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": "190"
}
```       

## query all instatial chaincode list                                     
URL：http://ip:port/api/v1/queryinstatialcc                                 
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/queryinstatialcc                                           
request parameter:            
```json    
{
	"BlockChainId": "6f2Wq0Hy5LSYpSt0yvmTd1XD2XEn3Hdr",
	"OrgName": "Org1",
	"ChannelId": "mychannel" 
}
``` 
API RETURN (fabric v1.x):                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": {
        "chaincodes": [ 
            {
                "name": "cctest",
				"version": "10",
				"sequence": 1,
                "path": "cctest10",
                "input": "<nil>",
                "escc": "escc",
                "vscc": "vscc"
            }
        ]
    }
}
or   
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "CCName": "cctest",
            "CCVersion": "10",
            "CreateTime": "2020-07-04 09:28:08",
			"UpdateTime": "2020-07-04 09:28:08",
			"Status": 0:success,101:deploying,102:deploy failed
        }
    ]
}
```         

## query all channel list                                       
URL：http://ip:port/api/v1/querychannel                                 
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/querychannel                                           
request parameter:            
```json    
{
	"BlockChainId": "6f2Wq0Hy5LSYpSt0yvmTd1XD2XEn3Hdr",
	"OrgName": "Org1"
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": {
        "channels": [
            {
                "channel_id": "mychannel"
            }
        ]
    }
}
```         

## query transaction info                                             
URL：http://ip:port/api/v1/querytxinfo                                    
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/querytxinfo                                                  
request parameter:            
```json    
{
	"BlockChainId": "1m4MJ3U0avy8VLWkQeXCtgN8M3384TAd",
	"OrgName": "Org1",
	"ChannelId": "mychannel",
	"TxId": "4cb2d7db491f8f9c517e89e1332f539f4d5d8c335d6def0a4c18facb32a53c23" 
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
     "data": {
        "signature": "MEUCIQDkm4uM4KO67eiPKQjF7Bir9ycFo+N4WZznAfqwbxRRHQIgdAZI/+g1oJjPXcqgBcgQBlfdjrTfnDH/JD16vllCKBE=",
        "channel_header": {
            "type": 3,
            "channel_id": "mychannel",
            "tx_id": "4cb2d7db491f8f9c517e89e1332f539f4d5d8c335d6def0a4c18facb32a53c23",
            "chaincode_id": {
                "name": "cctest"
            }
        },
        "signature_header": {
            "Certificate": {
                "Raw": "MIICvDCCAmKgAwIBAgIUURtISOgT2ZyyoVYdpaicqfdgw3cwCgYIKoZIzj0EAwIwezELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xHTAbBgNVBAoTFE9yZzEuZmFicmljLmJhYXMueHl6MSAwHgYDVQQDExdjYS5PcmcxLmZhYnJpYy5iYWFzLnh5ejAeFw0yMDA1MDYwMTE0MDBaFw0yMTA1MDYwMTE5MDBaMEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABEoxdRDaQflbiibhX/L2eyaJyila9DhFjZrP0DuQg6mU+leKSdmMZXngJDGAf3CkU22BkRM6Wki33j/GG8lvXgajgf0wgfowDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFN46QltTS47uQacqRnJfdcYwEaENMCsGA1UdIwQkMCKAIDIcSAqtJUYNoKJUIuN652nHnULplRCzxv1XSscEbuI+MCUGA1UdEQQeMByCGmNxeWNzeGZkZU1hY0Jvb2stUHJvLmxvY2FsMGcGCCoDBAUGBwgBBFt7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9uIjoib3JnMS5kZXBhcnRtZW50MSIsImhmLkVucm9sbG1lbnRJRCI6ImNodXNlciIsImhmLlR5cGUiOiJ1c2VyIn19MAoGCCqGSM49BAMCA0gAMEUCIQClqEDU7TVZwT3D0xCLprb11lVE9QrBz1ycbKaPsCLqEwIgaVrICuRgc/4HdT7Y0TubJGVfOwXXDMrSODPRHzwaOsk=",
                "RawTBSCertificate": "MIICYqADAgECAhRRG0hI6BPZnLKhVh2lqJyp92DDdzAKBggqhkjOPQQDAjB7MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEdMBsGA1UEChMUT3JnMS5mYWJyaWMuYmFhcy54eXoxIDAeBgNVBAMTF2NhLk9yZzEuZmFicmljLmJhYXMueHl6MB4XDTIwMDUwNjAxMTQwMFoXDTIxMDUwNjAxMTkwMFowQTEuMAsGA1UECxMEdXNlcjALBgNVBAsTBG9yZzEwEgYDVQQLEwtkZXBhcnRtZW50MTEPMA0GA1UEAxMGY2h1c2VyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESjF1ENpB+VuKJuFf8vZ7JonKKVr0OEWNms/QO5CDqZT6V4pJ2YxleeAkMYB/cKRTbYGREzpaSLfeP8YbyW9eBqOB/TCB+jAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU3jpCW1NLju5BpypGcl91xjARoQ0wKwYDVR0jBCQwIoAgMhxICq0lRg2golQi43rnacedQumVELPG/VdKxwRu4j4wJQYDVR0RBB4wHIIaY3F5Y3N4ZmRlTWFjQm9vay1Qcm8ubG9jYWwwZwYIKgMEBQYHCAEEW3siYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiJvcmcxLmRlcGFydG1lbnQxIiwiaGYuRW5yb2xsbWVudElEIjoiY2h1c2VyIiwiaGYuVHlwZSI6InVzZXIifX0=",
                "RawSubjectPublicKeyInfo": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESjF1ENpB+VuKJuFf8vZ7JonKKVr0OEWNms/QO5CDqZT6V4pJ2YxleeAkMYB/cKRTbYGREzpaSLfeP8YbyW9eBg==",
                "RawSubject": "MEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcg==",
                "RawIssuer": "MHsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMR0wGwYDVQQKExRPcmcxLmZhYnJpYy5iYWFzLnh5ejEgMB4GA1UEAxMXY2EuT3JnMS5mYWJyaWMuYmFhcy54eXo=",
                "Signature": "MEUCIQClqEDU7TVZwT3D0xCLprb11lVE9QrBz1ycbKaPsCLqEwIgaVrICuRgc/4HdT7Y0TubJGVfOwXXDMrSODPRHzwaOsk=",
                "SignatureAlgorithm": 10,
                "PublicKeyAlgorithm": 3,
                "PublicKey": {
                    "Curve": {
                        "P": 1.15792089210356248762697446949407573530086143415290314195533631308867097853951e+77,
                        "N": 1.15792089210356248762697446949407573529996955224135760342422259061068512044369e+77,
                        "B": 4.1058363725152142129326129780047268409114441015993725554835256314039467401291e+76,
                        "Gx": 4.8439561293906451759052585252797914202762949526041747995844080717082404635286e+76,
                        "Gy": 3.6134250956749795798585127919587881956611106672985015071877198253568414405109e+76,
                        "BitSize": 256,
                        "Name": "P-256"
                    },
                    "X": 3.3558534260002794477127624873681564286798929101483351952104778506033896335764e+76,
                    "Y": 1.13232882272434803347155020456148948598479428090050511070716500975209826704902e+77
                },
                "Version": 3,
                "SerialNumber": 4.63036669450492795693207409130645271472699917175e+47,
                "Issuer": {
                    "Country": [
                        "US"
                    ],
                    "Organization": [
                        "Org1.fabric.baas.xyz"
                    ],
                    "OrganizationalUnit": null,
                    "Locality": [
                        "San Francisco"
                    ],
                    "Province": [
                        "California"
                    ],
                    "StreetAddress": null,
                    "PostalCode": null,
                    "SerialNumber": "",
                    "CommonName": "ca.Org1.fabric.baas.xyz",
                    "Names": [
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                6
                            ],
                            "Value": "US"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                8
                            ],
                            "Value": "California"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                7
                            ],
                            "Value": "San Francisco"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                10
                            ],
                            "Value": "Org1.fabric.baas.xyz"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                3
                            ],
                            "Value": "ca.Org1.fabric.baas.xyz"
                        }
                    ],
                    "ExtraNames": null
                },
                "Subject": {
                    "Country": null,
                    "Organization": null,
                    "OrganizationalUnit": [
                        "user",
                        "org1",
                        "department1"
                    ],
                    "Locality": null,
                    "Province": null,
                    "StreetAddress": null,
                    "PostalCode": null,
                    "SerialNumber": "",
                    "CommonName": "chuser",
                    "Names": [
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                11
                            ],
                            "Value": "user"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                11
                            ],
                            "Value": "org1"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                11
                            ],
                            "Value": "department1"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                3
                            ],
                            "Value": "chuser"
                        }
                    ],
                    "ExtraNames": null
                },
                "NotBefore": "2020-05-06T01:14:00Z",
                "NotAfter": "2021-05-06T01:19:00Z",
                "KeyUsage": 1,
                "Extensions": [
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            15
                        ],
                        "Critical": true,
                        "Value": "AwIHgA=="
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            19
                        ],
                        "Critical": true,
                        "Value": "MAA="
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            14
                        ],
                        "Critical": false,
                        "Value": "BBTeOkJbU0uO7kGnKkZyX3XGMBGhDQ=="
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            35
                        ],
                        "Critical": false,
                        "Value": "MCKAIDIcSAqtJUYNoKJUIuN652nHnULplRCzxv1XSscEbuI+"
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            17
                        ],
                        "Critical": false,
                        "Value": "MByCGmNxeWNzeGZkZU1hY0Jvb2stUHJvLmxvY2Fs"
                    },
                    {
                        "Id": [
                            1,
                            2,
                            3,
                            4,
                            5,
                            6,
                            7,
                            8,
                            1
                        ],
                        "Critical": false,
                        "Value": "eyJhdHRycyI6eyJoZi5BZmZpbGlhdGlvbiI6Im9yZzEuZGVwYXJ0bWVudDEiLCJoZi5FbnJvbGxtZW50SUQiOiJjaHVzZXIiLCJoZi5UeXBlIjoidXNlciJ9fQ=="
                    }
                ],
                "ExtraExtensions": null,
                "UnhandledCriticalExtensions": null,
                "ExtKeyUsage": null,
                "UnknownExtKeyUsage": null,
                "BasicConstraintsValid": true,
                "IsCA": false,
                "MaxPathLen": -1,
                "MaxPathLenZero": false,
                "SubjectKeyId": "3jpCW1NLju5BpypGcl91xjARoQ0=",
                "AuthorityKeyId": "MhxICq0lRg2golQi43rnacedQumVELPG/VdKxwRu4j4=",
                "OCSPServer": null,
                "IssuingCertificateURL": null,
                "DNSNames": [
                    "cqycsxfdeMacBook-Pro.local"
                ],
                "EmailAddresses": null,
                "IPAddresses": null,
                "URIs": null,
                "PermittedDNSDomainsCritical": false,
                "PermittedDNSDomains": null,
                "ExcludedDNSDomains": null,
                "PermittedIPRanges": null,
                "ExcludedIPRanges": null,
                "PermittedEmailAddresses": null,
                "ExcludedEmailAddresses": null,
                "PermittedURIDomains": null,
                "ExcludedURIDomains": null,
                "CRLDistributionPoints": null,
                "PolicyIdentifiers": null
            },
            "nonce": "JU2AL0bE28bNwBEP/vlB3riu/Js6Tgz6"
        },
        "tx_action_signature_header": {
            "Certificate": {
                "Raw": "MIICvDCCAmKgAwIBAgIUURtISOgT2ZyyoVYdpaicqfdgw3cwCgYIKoZIzj0EAwIwezELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xHTAbBgNVBAoTFE9yZzEuZmFicmljLmJhYXMueHl6MSAwHgYDVQQDExdjYS5PcmcxLmZhYnJpYy5iYWFzLnh5ejAeFw0yMDA1MDYwMTE0MDBaFw0yMTA1MDYwMTE5MDBaMEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABEoxdRDaQflbiibhX/L2eyaJyila9DhFjZrP0DuQg6mU+leKSdmMZXngJDGAf3CkU22BkRM6Wki33j/GG8lvXgajgf0wgfowDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFN46QltTS47uQacqRnJfdcYwEaENMCsGA1UdIwQkMCKAIDIcSAqtJUYNoKJUIuN652nHnULplRCzxv1XSscEbuI+MCUGA1UdEQQeMByCGmNxeWNzeGZkZU1hY0Jvb2stUHJvLmxvY2FsMGcGCCoDBAUGBwgBBFt7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9uIjoib3JnMS5kZXBhcnRtZW50MSIsImhmLkVucm9sbG1lbnRJRCI6ImNodXNlciIsImhmLlR5cGUiOiJ1c2VyIn19MAoGCCqGSM49BAMCA0gAMEUCIQClqEDU7TVZwT3D0xCLprb11lVE9QrBz1ycbKaPsCLqEwIgaVrICuRgc/4HdT7Y0TubJGVfOwXXDMrSODPRHzwaOsk=",
                "RawTBSCertificate": "MIICYqADAgECAhRRG0hI6BPZnLKhVh2lqJyp92DDdzAKBggqhkjOPQQDAjB7MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEdMBsGA1UEChMUT3JnMS5mYWJyaWMuYmFhcy54eXoxIDAeBgNVBAMTF2NhLk9yZzEuZmFicmljLmJhYXMueHl6MB4XDTIwMDUwNjAxMTQwMFoXDTIxMDUwNjAxMTkwMFowQTEuMAsGA1UECxMEdXNlcjALBgNVBAsTBG9yZzEwEgYDVQQLEwtkZXBhcnRtZW50MTEPMA0GA1UEAxMGY2h1c2VyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESjF1ENpB+VuKJuFf8vZ7JonKKVr0OEWNms/QO5CDqZT6V4pJ2YxleeAkMYB/cKRTbYGREzpaSLfeP8YbyW9eBqOB/TCB+jAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU3jpCW1NLju5BpypGcl91xjARoQ0wKwYDVR0jBCQwIoAgMhxICq0lRg2golQi43rnacedQumVELPG/VdKxwRu4j4wJQYDVR0RBB4wHIIaY3F5Y3N4ZmRlTWFjQm9vay1Qcm8ubG9jYWwwZwYIKgMEBQYHCAEEW3siYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiJvcmcxLmRlcGFydG1lbnQxIiwiaGYuRW5yb2xsbWVudElEIjoiY2h1c2VyIiwiaGYuVHlwZSI6InVzZXIifX0=",
                "RawSubjectPublicKeyInfo": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESjF1ENpB+VuKJuFf8vZ7JonKKVr0OEWNms/QO5CDqZT6V4pJ2YxleeAkMYB/cKRTbYGREzpaSLfeP8YbyW9eBg==",
                "RawSubject": "MEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcg==",
                "RawIssuer": "MHsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMR0wGwYDVQQKExRPcmcxLmZhYnJpYy5iYWFzLnh5ejEgMB4GA1UEAxMXY2EuT3JnMS5mYWJyaWMuYmFhcy54eXo=",
                "Signature": "MEUCIQClqEDU7TVZwT3D0xCLprb11lVE9QrBz1ycbKaPsCLqEwIgaVrICuRgc/4HdT7Y0TubJGVfOwXXDMrSODPRHzwaOsk=",
                "SignatureAlgorithm": 10,
                "PublicKeyAlgorithm": 3,
                "PublicKey": {
                    "Curve": {
                        "P": 1.15792089210356248762697446949407573530086143415290314195533631308867097853951e+77,
                        "N": 1.15792089210356248762697446949407573529996955224135760342422259061068512044369e+77,
                        "B": 4.1058363725152142129326129780047268409114441015993725554835256314039467401291e+76,
                        "Gx": 4.8439561293906451759052585252797914202762949526041747995844080717082404635286e+76,
                        "Gy": 3.6134250956749795798585127919587881956611106672985015071877198253568414405109e+76,
                        "BitSize": 256,
                        "Name": "P-256"
                    },
                    "X": 3.3558534260002794477127624873681564286798929101483351952104778506033896335764e+76,
                    "Y": 1.13232882272434803347155020456148948598479428090050511070716500975209826704902e+77
                },
                "Version": 3,
                "SerialNumber": 4.63036669450492795693207409130645271472699917175e+47,
                "Issuer": {
                    "Country": [
                        "US"
                    ],
                    "Organization": [
                        "Org1.fabric.baas.xyz"
                    ],
                    "OrganizationalUnit": null,
                    "Locality": [
                        "San Francisco"
                    ],
                    "Province": [
                        "California"
                    ],
                    "StreetAddress": null,
                    "PostalCode": null,
                    "SerialNumber": "",
                    "CommonName": "ca.Org1.fabric.baas.xyz",
                    "Names": [
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                6
                            ],
                            "Value": "US"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                8
                            ],
                            "Value": "California"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                7
                            ],
                            "Value": "San Francisco"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                10
                            ],
                            "Value": "Org1.fabric.baas.xyz"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                3
                            ],
                            "Value": "ca.Org1.fabric.baas.xyz"
                        }
                    ],
                    "ExtraNames": null
                },
                "Subject": {
                    "Country": null,
                    "Organization": null,
                    "OrganizationalUnit": [
                        "user",
                        "org1",
                        "department1"
                    ],
                    "Locality": null,
                    "Province": null,
                    "StreetAddress": null,
                    "PostalCode": null,
                    "SerialNumber": "",
                    "CommonName": "chuser",
                    "Names": [
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                11
                            ],
                            "Value": "user"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                11
                            ],
                            "Value": "org1"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                11
                            ],
                            "Value": "department1"
                        },
                        {
                            "Type": [
                                2,
                                5,
                                4,
                                3
                            ],
                            "Value": "chuser"
                        }
                    ],
                    "ExtraNames": null
                },
                "NotBefore": "2020-05-06T01:14:00Z",
                "NotAfter": "2021-05-06T01:19:00Z",
                "KeyUsage": 1,
                "Extensions": [
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            15
                        ],
                        "Critical": true,
                        "Value": "AwIHgA=="
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            19
                        ],
                        "Critical": true,
                        "Value": "MAA="
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            14
                        ],
                        "Critical": false,
                        "Value": "BBTeOkJbU0uO7kGnKkZyX3XGMBGhDQ=="
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            35
                        ],
                        "Critical": false,
                        "Value": "MCKAIDIcSAqtJUYNoKJUIuN652nHnULplRCzxv1XSscEbuI+"
                    },
                    {
                        "Id": [
                            2,
                            5,
                            29,
                            17
                        ],
                        "Critical": false,
                        "Value": "MByCGmNxeWNzeGZkZU1hY0Jvb2stUHJvLmxvY2Fs"
                    },
                    {
                        "Id": [
                            1,
                            2,
                            3,
                            4,
                            5,
                            6,
                            7,
                            8,
                            1
                        ],
                        "Critical": false,
                        "Value": "eyJhdHRycyI6eyJoZi5BZmZpbGlhdGlvbiI6Im9yZzEuZGVwYXJ0bWVudDEiLCJoZi5FbnJvbGxtZW50SUQiOiJjaHVzZXIiLCJoZi5UeXBlIjoidXNlciJ9fQ=="
                    }
                ],
                "ExtraExtensions": null,
                "UnhandledCriticalExtensions": null,
                "ExtKeyUsage": null,
                "UnknownExtKeyUsage": null,
                "BasicConstraintsValid": true,
                "IsCA": false,
                "MaxPathLen": -1,
                "MaxPathLenZero": false,
                "SubjectKeyId": "3jpCW1NLju5BpypGcl91xjARoQ0=",
                "AuthorityKeyId": "MhxICq0lRg2golQi43rnacedQumVELPG/VdKxwRu4j4=",
                "OCSPServer": null,
                "IssuingCertificateURL": null,
                "DNSNames": [
                    "cqycsxfdeMacBook-Pro.local"
                ],
                "EmailAddresses": null,
                "IPAddresses": null,
                "URIs": null,
                "PermittedDNSDomainsCritical": false,
                "PermittedDNSDomains": null,
                "ExcludedDNSDomains": null,
                "PermittedIPRanges": null,
                "ExcludedIPRanges": null,
                "PermittedEmailAddresses": null,
                "ExcludedEmailAddresses": null,
                "PermittedURIDomains": null,
                "ExcludedURIDomains": null,
                "CRLDistributionPoints": null,
                "PolicyIdentifiers": null
            },
            "nonce": "JU2AL0bE28bNwBEP/vlB3riu/Js6Tgz6"
        },
        "chaincode_spec": {
            "type": 1,
            "chaincode_id": {
                "name": "cctest"
            },
            "input": {
                "Args": [
                    "invoke",
                    "a",
                    "b",
                    "10"
                ]
            }
        },
        "endorsements": [
            {
                "signature_header": {
                    "Certificate": null,
                    "nonce": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLVENDQWRDZ0F3SUJBZ0lRRUY1ajIwWHBWdkhWWnBjd09UMUpyekFLQmdncWhrak9QUVFEQWpCN01Rc3cKQ1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ4TU5VMkZ1SUVaeQpZVzVqYVhOamJ6RWRNQnNHQTFVRUNoTVVUM0puTVM1bVlXSnlhV011WW1GaGN5NTRlWG94SURBZUJnTlZCQU1UCkYyTmhMazl5WnpFdVptRmljbWxqTG1KaFlYTXVlSGw2TUI0WERUSXdNRFV3TmpBeE1Ea3hNVm9YRFRNd01EVXcKTkRBeE1Ea3hNVm93WkRFTE1Ba0dBMVVFQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V4RmpBVQpCZ05WQkFjVERWTmhiaUJHY21GdVkybHpZMjh4S0RBbUJnTlZCQU1USDNCbFpYSXdMVzl5WnpFdVQzSm5NUzVtCllXSnlhV011WW1GaGN5NTRlWG93V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVF1eFZCTC8zNVAKdytNbHRrVCtnZkM3NEZKb1dWNisweVBWM0dDbzVteXRHbHhqd2U4aE9ZR2RlRFZETHRNTE8vRUwwZmZBMnF4QgpORnFCWlRhOWhLU2RvMDB3U3pBT0JnTlZIUThCQWY4RUJBTUNCNEF3REFZRFZSMFRBUUgvQkFJd0FEQXJCZ05WCkhTTUVKREFpZ0NBeUhFZ0tyU1ZHRGFDaVZDTGpldWRweDUxQzZaVVFzOGI5VjBySEJHN2lQakFLQmdncWhrak8KUFFRREFnTkhBREJFQWlCYW5TQktEdTB6K3FNR3AyNnNHQVVVV2ZJWmZBK1k1OXpBNURiZXFtL04zZ0lnRzZUdwpta2thcnhhTzZkS1hWb2lpaUJiQVc5dTRKNXIwU2VjRk11czcrR2M9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
                },
                "signature": "MEQCIFxfz2C1YqbQQt7eQwSNGrNcygMb64U8pt+V8xXlKAG/AiAv7YK27caRlSWZfyp+MNGG6OW7OHOpdjGRkACuLciHAg=="
            },
            {
                "signature_header": {
                    "Certificate": null,
                    "nonce": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLVENDQWRDZ0F3SUJBZ0lRZHJ4QmdmeTFZLzlFdzlyUXFtYjJSakFLQmdncWhrak9QUVFEQWpCN01Rc3cKQ1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ4TU5VMkZ1SUVaeQpZVzVqYVhOamJ6RWRNQnNHQTFVRUNoTVVUM0puTVM1bVlXSnlhV011WW1GaGN5NTRlWG94SURBZUJnTlZCQU1UCkYyTmhMazl5WnpFdVptRmljbWxqTG1KaFlYTXVlSGw2TUI0WERUSXdNRFV3TmpBeE1Ea3hNVm9YRFRNd01EVXcKTkRBeE1Ea3hNVm93WkRFTE1Ba0dBMVVFQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V4RmpBVQpCZ05WQkFjVERWTmhiaUJHY21GdVkybHpZMjh4S0RBbUJnTlZCQU1USDNCbFpYSXhMVzl5WnpFdVQzSm5NUzVtCllXSnlhV011WW1GaGN5NTRlWG93V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVNNUjNKbmFPckcKSm1CYXpUeE5rbHdHcWhSdTZKUGtya2ExKzJlazV6WmdZY1FvVlI5YXhnM3R0cFdjU2RqZFVGOEtUc3FVL2ZnWgpSL3pTY1JWSVNMNUhvMDB3U3pBT0JnTlZIUThCQWY4RUJBTUNCNEF3REFZRFZSMFRBUUgvQkFJd0FEQXJCZ05WCkhTTUVKREFpZ0NBeUhFZ0tyU1ZHRGFDaVZDTGpldWRweDUxQzZaVVFzOGI5VjBySEJHN2lQakFLQmdncWhrak8KUFFRREFnTkhBREJFQWlCc2FiZzNBK01sVWdNdTRnRndyUDlZTFVINnc5YnE4ME5RODRHVkFudFhwQUlnQWdxRQpMQ0MyRWRYaWRadTBNZzMrbFRydDk5U3NEZWM3UnRQazQxNStXMGs9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
                },
                "signature": "MEUCIQCyfJcZMNzRTpGCRNGd4ILGFo+wxB/cWRlkLud+QWZ4AwIgL1UoDJjEK7R+r6b/4iHYJGX2cfOyQ6GFn49vuN3zmJ8="
            }
        ],
        "proposal_hash": "LjF6cJEmGwzN2r2z79lnA8VBoM8HuRxaJUKsHBdUCks=",
        "events": {},
        "response": {
            "status": 200
        },
        "ns_read_write_Set": [
            {
                "Namespace": "cctest",
                "KVRWSet": {
                    "reads": [
                        {
                            "key": "a",
                            "version": {
                                "block_num": 3
                            }
                        },
                        {
                            "key": "b",
                            "version": {
                                "block_num": 3
                            }
                        }
                    ],
                    "writes": [
                        {
                            "key": "a",
                            "value": "MTcw"
                        },
                        {
                            "key": "b",
                            "value": "MjMw"
                        }
                    ]
                }
            },
            {
                "Namespace": "lscc",
                "KVRWSet": {
                    "reads": [
                        {
                            "key": "cctest",
                            "version": {
                                "block_num": 1
                            }
                        }
                    ]
                }
            }
        ],
        "validation_code": 0,
        "validation_code_name": "VALID"
    }
}
```        

## query block info                                           
URL：http://ip:port/api/v1/queryblockinfo                                     
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/queryblockinfo                                                 
request parameter:            
```json    
{
	"BlockChainId": "1m4MJ3U0avy8VLWkQeXCtgN8M3384TAd",
	"OrgName": "Org1",
	"ChannelId": "mychannel", 
	"BlockId": 2
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": {
        "header": {
		"number": 2,
		"previous_hash": "9DGubLEM4GBczl2ajzCeGK9/t/YJ6UG87XfkHDIH7Ig=",
		"data_hash": "AFVoh4IZJ/ajzdNA6vb/jC2OBP3LS6xjSIlPJIyPOOw="
	    },
        "transactions": [{
		"signature": "MEUCIQC4OQ8ScCWlcCOe2OncFE+U3eCfKHIXk2I46H/URVwYhgIgZsFu2p0AsRpO8vxw00yYuj8/MjYCQe0fRogNVkEwXa0=",
		"channel_header": {
			"type": 3,
			"channel_id": "mychannel",
			"tx_id": "c521725e1d03c17ca24129078c97993e267c96122766a643f0226bd9ffafb3ae",
			"chaincode_id": {
				"name": "cctest"
			}
		},
		"signature_header": {
			"Certificate": {
				"Raw": "MIICpzCCAk2gAwIBAgIUQVMJ4WKoU+LYbS5LAnoFrX4gIIcwCgYIKoZIzj0EAwIwezELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xHTAbBgNVBAoTFE9yZzEuZmFicmljLmJhYXMueHl6MSAwHgYDVQQDExdjYS5PcmcxLmZhYnJpYy5iYWFzLnh5ejAeFw0yMDA1MDcwMTM4MDBaFw0yMTA1MDcwMTQzMDBaMEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHf+bXMjzWx4F2rx0/cvC1hM4qWuHQTeopZspHPEn4hCBd4Seo0JEfHohw1puNcuZ4DBv5t/e49dUIR2IcjvuyejgegwgeUwDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFDsJ8DFe++aMQ6MakoleMRsTP5+GMCsGA1UdIwQkMCKAIHsN9tC8T3lYNBoZyvqKMibAW8nJd9RLhqS4JoW9CgT+MBAGA1UdEQQJMAeCBWJvZ29uMGcGCCoDBAUGBwgBBFt7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9uIjoib3JnMS5kZXBhcnRtZW50MSIsImhmLkVucm9sbG1lbnRJRCI6ImNodXNlciIsImhmLlR5cGUiOiJ1c2VyIn19MAoGCCqGSM49BAMCA0gAMEUCIQCF6SaJY1iTWQ/KcQNJkOaajZza+6iARvdFOWPr92rMTAIgSvhqCyOuPhksc7YCYZn7pfzU8qT2WBL+hoZpO7sS+X0=",
				"RawTBSCertificate": "MIICTaADAgECAhRBUwnhYqhT4thtLksCegWtfiAghzAKBggqhkjOPQQDAjB7MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEdMBsGA1UEChMUT3JnMS5mYWJyaWMuYmFhcy54eXoxIDAeBgNVBAMTF2NhLk9yZzEuZmFicmljLmJhYXMueHl6MB4XDTIwMDUwNzAxMzgwMFoXDTIxMDUwNzAxNDMwMFowQTEuMAsGA1UECxMEdXNlcjALBgNVBAsTBG9yZzEwEgYDVQQLEwtkZXBhcnRtZW50MTEPMA0GA1UEAxMGY2h1c2VyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEd/5tcyPNbHgXavHT9y8LWEzipa4dBN6ilmykc8SfiEIF3hJ6jQkR8eiHDWm41y5ngMG/m397j11QhHYhyO+7J6OB6DCB5TAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUOwnwMV775oxDoxqSiV4xGxM/n4YwKwYDVR0jBCQwIoAgew320LxPeVg0GhnK+ooyJsBbycl31EuGpLgmhb0KBP4wEAYDVR0RBAkwB4IFYm9nb24wZwYIKgMEBQYHCAEEW3siYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiJvcmcxLmRlcGFydG1lbnQxIiwiaGYuRW5yb2xsbWVudElEIjoiY2h1c2VyIiwiaGYuVHlwZSI6InVzZXIifX0=",
				"RawSubjectPublicKeyInfo": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEd/5tcyPNbHgXavHT9y8LWEzipa4dBN6ilmykc8SfiEIF3hJ6jQkR8eiHDWm41y5ngMG/m397j11QhHYhyO+7Jw==",
				"RawSubject": "MEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcg==",
				"RawIssuer": "MHsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMR0wGwYDVQQKExRPcmcxLmZhYnJpYy5iYWFzLnh5ejEgMB4GA1UEAxMXY2EuT3JnMS5mYWJyaWMuYmFhcy54eXo=",
				"Signature": "MEUCIQCF6SaJY1iTWQ/KcQNJkOaajZza+6iARvdFOWPr92rMTAIgSvhqCyOuPhksc7YCYZn7pfzU8qT2WBL+hoZpO7sS+X0=",
				"SignatureAlgorithm": 10,
				"PublicKeyAlgorithm": 3,
				"PublicKey": {
					"Curve": {
						"P": 115792089210356248762697446949407573530086143415290314195533631308867097853951,
						"N": 115792089210356248762697446949407573529996955224135760342422259061068512044369,
						"B": 41058363725152142129326129780047268409114441015993725554835256314039467401291,
						"Gx": 48439561293906451759052585252797914202762949526041747995844080717082404635286,
						"Gy": 36134250956749795798585127919587881956611106672985015071877198253568414405109,
						"BitSize": 256,
						"Name": "P-256"
					},
					"X": 54274763530378513350013542584945633205425191936301182671024633048682935060546,
					"Y": 2653931826697633232510707187333873436358028209562682241786600065930298899239
				},
				"Version": 3,
				"SerialNumber": 372936222660271904720974107094723451794319483015,
				"Issuer": {
					"Country": [
						"US"
					],
					"Organization": [
						"Org1.fabric.baas.xyz"
					],
					"OrganizationalUnit": null,
					"Locality": [
						"San Francisco"
					],
					"Province": [
						"California"
					],
					"StreetAddress": null,
					"PostalCode": null,
					"SerialNumber": "",
					"CommonName": "ca.Org1.fabric.baas.xyz",
					"Names": [{
							"Type": [
								2,
								5,
								4,
								6
							],
							"Value": "US"
						},
						{
							"Type": [
								2,
								5,
								4,
								8
							],
							"Value": "California"
						},
						{
							"Type": [
								2,
								5,
								4,
								7
							],
							"Value": "San Francisco"
						},
						{
							"Type": [
								2,
								5,
								4,
								10
							],
							"Value": "Org1.fabric.baas.xyz"
						},
						{
							"Type": [
								2,
								5,
								4,
								3
							],
							"Value": "ca.Org1.fabric.baas.xyz"
						}
					],
					"ExtraNames": null
				},
				"Subject": {
					"Country": null,
					"Organization": null,
					"OrganizationalUnit": [
						"user",
						"org1",
						"department1"
					],
					"Locality": null,
					"Province": null,
					"StreetAddress": null,
					"PostalCode": null,
					"SerialNumber": "",
					"CommonName": "chuser",
					"Names": [{
							"Type": [
								2,
								5,
								4,
								11
							],
							"Value": "user"
						},
						{
							"Type": [
								2,
								5,
								4,
								11
							],
							"Value": "org1"
						},
						{
							"Type": [
								2,
								5,
								4,
								11
							],
							"Value": "department1"
						},
						{
							"Type": [
								2,
								5,
								4,
								3
							],
							"Value": "chuser"
						}
					],
					"ExtraNames": null
				},
				"NotBefore": "2020-05-07T01:38:00Z",
				"NotAfter": "2021-05-07T01:43:00Z",
				"KeyUsage": 1,
				"Extensions": [{
						"Id": [
							2,
							5,
							29,
							15
						],
						"Critical": true,
						"Value": "AwIHgA=="
					},
					{
						"Id": [
							2,
							5,
							29,
							19
						],
						"Critical": true,
						"Value": "MAA="
					},
					{
						"Id": [
							2,
							5,
							29,
							14
						],
						"Critical": false,
						"Value": "BBQ7CfAxXvvmjEOjGpKJXjEbEz+fhg=="
					},
					{
						"Id": [
							2,
							5,
							29,
							35
						],
						"Critical": false,
						"Value": "MCKAIHsN9tC8T3lYNBoZyvqKMibAW8nJd9RLhqS4JoW9CgT+"
					},
					{
						"Id": [
							2,
							5,
							29,
							17
						],
						"Critical": false,
						"Value": "MAeCBWJvZ29u"
					},
					{
						"Id": [
							1,
							2,
							3,
							4,
							5,
							6,
							7,
							8,
							1
						],
						"Critical": false,
						"Value": "eyJhdHRycyI6eyJoZi5BZmZpbGlhdGlvbiI6Im9yZzEuZGVwYXJ0bWVudDEiLCJoZi5FbnJvbGxtZW50SUQiOiJjaHVzZXIiLCJoZi5UeXBlIjoidXNlciJ9fQ=="
					}
				],
				"ExtraExtensions": null,
				"UnhandledCriticalExtensions": null,
				"ExtKeyUsage": null,
				"UnknownExtKeyUsage": null,
				"BasicConstraintsValid": true,
				"IsCA": false,
				"MaxPathLen": -1,
				"MaxPathLenZero": false,
				"SubjectKeyId": "OwnwMV775oxDoxqSiV4xGxM/n4Y=",
				"AuthorityKeyId": "ew320LxPeVg0GhnK+ooyJsBbycl31EuGpLgmhb0KBP4=",
				"OCSPServer": null,
				"IssuingCertificateURL": null,
				"DNSNames": [
					"bogon"
				],
				"EmailAddresses": null,
				"IPAddresses": null,
				"URIs": null,
				"PermittedDNSDomainsCritical": false,
				"PermittedDNSDomains": null,
				"ExcludedDNSDomains": null,
				"PermittedIPRanges": null,
				"ExcludedIPRanges": null,
				"PermittedEmailAddresses": null,
				"ExcludedEmailAddresses": null,
				"PermittedURIDomains": null,
				"ExcludedURIDomains": null,
				"CRLDistributionPoints": null,
				"PolicyIdentifiers": null
			},
			"nonce": "hhYht3y96VuL2h6/uZFcGIm/NaMkltdj"
		},
		"tx_action_signature_header": {
			"Certificate": {
				"Raw": "MIICpzCCAk2gAwIBAgIUQVMJ4WKoU+LYbS5LAnoFrX4gIIcwCgYIKoZIzj0EAwIwezELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xHTAbBgNVBAoTFE9yZzEuZmFicmljLmJhYXMueHl6MSAwHgYDVQQDExdjYS5PcmcxLmZhYnJpYy5iYWFzLnh5ejAeFw0yMDA1MDcwMTM4MDBaFw0yMTA1MDcwMTQzMDBaMEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHf+bXMjzWx4F2rx0/cvC1hM4qWuHQTeopZspHPEn4hCBd4Seo0JEfHohw1puNcuZ4DBv5t/e49dUIR2IcjvuyejgegwgeUwDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFDsJ8DFe++aMQ6MakoleMRsTP5+GMCsGA1UdIwQkMCKAIHsN9tC8T3lYNBoZyvqKMibAW8nJd9RLhqS4JoW9CgT+MBAGA1UdEQQJMAeCBWJvZ29uMGcGCCoDBAUGBwgBBFt7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9uIjoib3JnMS5kZXBhcnRtZW50MSIsImhmLkVucm9sbG1lbnRJRCI6ImNodXNlciIsImhmLlR5cGUiOiJ1c2VyIn19MAoGCCqGSM49BAMCA0gAMEUCIQCF6SaJY1iTWQ/KcQNJkOaajZza+6iARvdFOWPr92rMTAIgSvhqCyOuPhksc7YCYZn7pfzU8qT2WBL+hoZpO7sS+X0=",
				"RawTBSCertificate": "MIICTaADAgECAhRBUwnhYqhT4thtLksCegWtfiAghzAKBggqhkjOPQQDAjB7MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEdMBsGA1UEChMUT3JnMS5mYWJyaWMuYmFhcy54eXoxIDAeBgNVBAMTF2NhLk9yZzEuZmFicmljLmJhYXMueHl6MB4XDTIwMDUwNzAxMzgwMFoXDTIxMDUwNzAxNDMwMFowQTEuMAsGA1UECxMEdXNlcjALBgNVBAsTBG9yZzEwEgYDVQQLEwtkZXBhcnRtZW50MTEPMA0GA1UEAxMGY2h1c2VyMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEd/5tcyPNbHgXavHT9y8LWEzipa4dBN6ilmykc8SfiEIF3hJ6jQkR8eiHDWm41y5ngMG/m397j11QhHYhyO+7J6OB6DCB5TAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUOwnwMV775oxDoxqSiV4xGxM/n4YwKwYDVR0jBCQwIoAgew320LxPeVg0GhnK+ooyJsBbycl31EuGpLgmhb0KBP4wEAYDVR0RBAkwB4IFYm9nb24wZwYIKgMEBQYHCAEEW3siYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiJvcmcxLmRlcGFydG1lbnQxIiwiaGYuRW5yb2xsbWVudElEIjoiY2h1c2VyIiwiaGYuVHlwZSI6InVzZXIifX0=",
				"RawSubjectPublicKeyInfo": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEd/5tcyPNbHgXavHT9y8LWEzipa4dBN6ilmykc8SfiEIF3hJ6jQkR8eiHDWm41y5ngMG/m397j11QhHYhyO+7Jw==",
				"RawSubject": "MEExLjALBgNVBAsTBHVzZXIwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDzANBgNVBAMTBmNodXNlcg==",
				"RawIssuer": "MHsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMR0wGwYDVQQKExRPcmcxLmZhYnJpYy5iYWFzLnh5ejEgMB4GA1UEAxMXY2EuT3JnMS5mYWJyaWMuYmFhcy54eXo=",
				"Signature": "MEUCIQCF6SaJY1iTWQ/KcQNJkOaajZza+6iARvdFOWPr92rMTAIgSvhqCyOuPhksc7YCYZn7pfzU8qT2WBL+hoZpO7sS+X0=",
				"SignatureAlgorithm": 10,
				"PublicKeyAlgorithm": 3,
				"PublicKey": {
					"Curve": {
						"P": 115792089210356248762697446949407573530086143415290314195533631308867097853951,
						"N": 115792089210356248762697446949407573529996955224135760342422259061068512044369,
						"B": 41058363725152142129326129780047268409114441015993725554835256314039467401291,
						"Gx": 48439561293906451759052585252797914202762949526041747995844080717082404635286,
						"Gy": 36134250956749795798585127919587881956611106672985015071877198253568414405109,
						"BitSize": 256,
						"Name": "P-256"
					},
					"X": 54274763530378513350013542584945633205425191936301182671024633048682935060546,
					"Y": 2653931826697633232510707187333873436358028209562682241786600065930298899239
				},
				"Version": 3,
				"SerialNumber": 372936222660271904720974107094723451794319483015,
				"Issuer": {
					"Country": [
						"US"
					],
					"Organization": [
						"Org1.fabric.baas.xyz"
					],
					"OrganizationalUnit": null,
					"Locality": [
						"San Francisco"
					],
					"Province": [
						"California"
					],
					"StreetAddress": null,
					"PostalCode": null,
					"SerialNumber": "",
					"CommonName": "ca.Org1.fabric.baas.xyz",
					"Names": [{
							"Type": [
								2,
								5,
								4,
								6
							],
							"Value": "US"
						},
						{
							"Type": [
								2,
								5,
								4,
								8
							],
							"Value": "California"
						},
						{
							"Type": [
								2,
								5,
								4,
								7
							],
							"Value": "San Francisco"
						},
						{
							"Type": [
								2,
								5,
								4,
								10
							],
							"Value": "Org1.fabric.baas.xyz"
						},
						{
							"Type": [
								2,
								5,
								4,
								3
							],
							"Value": "ca.Org1.fabric.baas.xyz"
						}
					],
					"ExtraNames": null
				},
				"Subject": {
					"Country": null,
					"Organization": null,
					"OrganizationalUnit": [
						"user",
						"org1",
						"department1"
					],
					"Locality": null,
					"Province": null,
					"StreetAddress": null,
					"PostalCode": null,
					"SerialNumber": "",
					"CommonName": "chuser",
					"Names": [{
							"Type": [
								2,
								5,
								4,
								11
							],
							"Value": "user"
						},
						{
							"Type": [
								2,
								5,
								4,
								11
							],
							"Value": "org1"
						},
						{
							"Type": [
								2,
								5,
								4,
								11
							],
							"Value": "department1"
						},
						{
							"Type": [
								2,
								5,
								4,
								3
							],
							"Value": "chuser"
						}
					],
					"ExtraNames": null
				},
				"NotBefore": "2020-05-07T01:38:00Z",
				"NotAfter": "2021-05-07T01:43:00Z",
				"KeyUsage": 1,
				"Extensions": [{
						"Id": [
							2,
							5,
							29,
							15
						],
						"Critical": true,
						"Value": "AwIHgA=="
					},
					{
						"Id": [
							2,
							5,
							29,
							19
						],
						"Critical": true,
						"Value": "MAA="
					},
					{
						"Id": [
							2,
							5,
							29,
							14
						],
						"Critical": false,
						"Value": "BBQ7CfAxXvvmjEOjGpKJXjEbEz+fhg=="
					},
					{
						"Id": [
							2,
							5,
							29,
							35
						],
						"Critical": false,
						"Value": "MCKAIHsN9tC8T3lYNBoZyvqKMibAW8nJd9RLhqS4JoW9CgT+"
					},
					{
						"Id": [
							2,
							5,
							29,
							17
						],
						"Critical": false,
						"Value": "MAeCBWJvZ29u"
					},
					{
						"Id": [
							1,
							2,
							3,
							4,
							5,
							6,
							7,
							8,
							1
						],
						"Critical": false,
						"Value": "eyJhdHRycyI6eyJoZi5BZmZpbGlhdGlvbiI6Im9yZzEuZGVwYXJ0bWVudDEiLCJoZi5FbnJvbGxtZW50SUQiOiJjaHVzZXIiLCJoZi5UeXBlIjoidXNlciJ9fQ=="
					}
				],
				"ExtraExtensions": null,
				"UnhandledCriticalExtensions": null,
				"ExtKeyUsage": null,
				"UnknownExtKeyUsage": null,
				"BasicConstraintsValid": true,
				"IsCA": false,
				"MaxPathLen": -1,
				"MaxPathLenZero": false,
				"SubjectKeyId": "OwnwMV775oxDoxqSiV4xGxM/n4Y=",
				"AuthorityKeyId": "ew320LxPeVg0GhnK+ooyJsBbycl31EuGpLgmhb0KBP4=",
				"OCSPServer": null,
				"IssuingCertificateURL": null,
				"DNSNames": [
					"bogon"
				],
				"EmailAddresses": null,
				"IPAddresses": null,
				"URIs": null,
				"PermittedDNSDomainsCritical": false,
				"PermittedDNSDomains": null,
				"ExcludedDNSDomains": null,
				"PermittedIPRanges": null,
				"ExcludedIPRanges": null,
				"PermittedEmailAddresses": null,
				"ExcludedEmailAddresses": null,
				"PermittedURIDomains": null,
				"ExcludedURIDomains": null,
				"CRLDistributionPoints": null,
				"PolicyIdentifiers": null
			},
			"nonce": "hhYht3y96VuL2h6/uZFcGIm/NaMkltdj"
		},
		"chaincode_spec": {
			"type": 1,
			"chaincode_id": {
				"name": "cctest"
			},
			"input": {
				"Args": [
					"invoke",
					"a",
					"b",
					"10"
				]
			}
		},
		"endorsements": [{
				"signature_header": {
					"Certificate": null,
					"nonce": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLakNDQWRDZ0F3SUJBZ0lRVW8vNzNlRXQzYVgydGJNcU9ab2Y0akFLQmdncWhrak9QUVFEQWpCN01Rc3cKQ1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ4TU5VMkZ1SUVaeQpZVzVqYVhOamJ6RWRNQnNHQTFVRUNoTVVUM0puTVM1bVlXSnlhV011WW1GaGN5NTRlWG94SURBZUJnTlZCQU1UCkYyTmhMazl5WnpFdVptRmljbWxqTG1KaFlYTXVlSGw2TUI0WERUSXdNRFV3TnpBeE16UXpPRm9YRFRNd01EVXcKTlRBeE16UXpPRm93WkRFTE1Ba0dBMVVFQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V4RmpBVQpCZ05WQkFjVERWTmhiaUJHY21GdVkybHpZMjh4S0RBbUJnTlZCQU1USDNCbFpYSXdMVzl5WnpFdVQzSm5NUzVtCllXSnlhV011WW1GaGN5NTRlWG93V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVEyKy9uVVhaRXcKSEZMVzJzcGdmYWpUalB2Z201OGdydGl6V0wrZGdleTV2WTRLMWhZQ0N6TzByZm81SU1UZnpPSWRnWUZKRm1OUgovSGtVdkJKcEJVazlvMDB3U3pBT0JnTlZIUThCQWY4RUJBTUNCNEF3REFZRFZSMFRBUUgvQkFJd0FEQXJCZ05WCkhTTUVKREFpZ0NCN0RmYlF2RTk1V0RRYUdjcjZpakltd0Z2SnlYZlVTNGFrdUNhRnZRb0UvakFLQmdncWhrak8KUFFRREFnTklBREJGQWlFQTgwVU1IYzRNU3FNNisrOURnc3JJWjV3RW9OZW9JZklQQ1JpQXVSMVBCc1VDSUg2egpjWTM5UW1lMHdLbERySmgrckpCVUJLdkQ0ZEdlUUNHeE9kaDZvVHVkCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
				},
				"signature": "MEUCIQDwa7j2RqmiQBJmRpuUn0fKEyTL3Eeo8LXpPXlB19OefAIgMJAKBoBqwfSjuH1Gtk4Bz2ezaVfTqzuvOcHgQq9q0Og="
			},
			{
				"signature_header": {
					"Certificate": null,
					"nonce": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNLekNDQWRHZ0F3SUJBZ0lSQVBhVGM0aDloeUEzNzh6WjM1ZDlwdm93Q2dZSUtvWkl6ajBFQXdJd2V6RUwKTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmhiR2xtYjNKdWFXRXhGakFVQmdOVkJBY1REVk5oYmlCRwpjbUZ1WTJselkyOHhIVEFiQmdOVkJBb1RGRTl5WnpFdVptRmljbWxqTG1KaFlYTXVlSGw2TVNBd0hnWURWUVFECkV4ZGpZUzVQY21jeExtWmhZbkpwWXk1aVlXRnpMbmg1ZWpBZUZ3MHlNREExTURjd01UTTBNemhhRncwek1EQTEKTURVd01UTTBNemhhTUdReEN6QUpCZ05WQkFZVEFsVlRNUk13RVFZRFZRUUlFd3BEWVd4cFptOXlibWxoTVJZdwpGQVlEVlFRSEV3MVRZVzRnUm5KaGJtTnBjMk52TVNnd0pnWURWUVFERXg5d1pXVnlNUzF2Y21jeExrOXlaekV1ClptRmljbWxqTG1KaFlYTXVlSGw2TUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFOGducnAzeVgKY2JhMXNpeHJSUEVuMGNuRFRvSDg3ZVZQZGp1UVNuaTI4S09WazBnNDN2c3c3dTByUE5hQ0UxdWMxU0tvQm1BYgo1VCtZeDJZYkVMMkF6Nk5OTUVzd0RnWURWUjBQQVFIL0JBUURBZ2VBTUF3R0ExVWRFd0VCL3dRQ01BQXdLd1lEClZSMGpCQ1F3SW9BZ2V3MzIwTHhQZVZnMEdobksrb295SnNCYnljbDMxRXVHcExnbWhiMEtCUDR3Q2dZSUtvWkkKemowRUF3SURTQUF3UlFJaEFKMXcxUll1RmZoeGRBNUtRQ2RxeU5jWU9rNzJ0OWNqUUZxVnEzOGFYT1A3QWlCcQpxNnJCQjd1VkxrRUxUM0tKOWwxakU1TDY2UGErWCtqb00yR05IWDJpWmc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
				},
				"signature": "MEUCIQD9PMO8i9wwU7tcJjCfQ1Xfrxma4bmMY84obQHCFif/8wIga+o5SIm/4ObUnbxA8DBVnNSEYClwsHhJMlS0msmPTh4="
			}
		],
		"proposal_hash": "q6GGG4Kk1TtD6nQQ4HtW6F9JCrZb0RPdD5dKv+17WiE=",
		"events": {},
		"response": {
			"status": 200
		},
		"ns_read_write_Set": [{
				"Namespace": "cctest",
				"KVRWSet": {
					"reads": [{
							"key": "a",
							"version": {
								"block_num": 1
							}
						},
						{
							"key": "b",
							"version": {
								"block_num": 1
							}
						}
					],
					"writes": [{
							"key": "a",
							"value": "MTkw"
						},
						{
							"key": "b",
							"value": "MjEw"
						}
					]
				}
			},
			{
				"Namespace": "lscc",
				"KVRWSet": {
					"reads": [{
						"key": "cctest",
						"version": {
							"block_num": 1
						}
					}]
				}
			}
		],
		"validation_code": 0,
		"validation_code_name": "VALID"
	}],
	"block_creator_signature": {
		"signature_header": {
			"Certificate": {
				"Raw": "MIICHDCCAcOgAwIBAgIRAKiEDmvKK6qlCwIN1CqLQmUwCgYIKoZIzj0EAwIwczELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xGTAXBgNVBAoTEG9yZGVyZXIuYmFhcy54eXoxHDAaBgNVBAMTE2NhLm9yZGVyZXIuYmFhcy54eXowHhcNMjAwNTA3MDEzNDM4WhcNMzAwNTA1MDEzNDM4WjBeMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEiMCAGA1UEAxMZb3JkZXJlcjAub3JkZXJlci5iYWFzLnh5ejBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABIFhj50QgM7UgbrZvwpHw381mL7pBWAEteGyOissqOehaXCWtJyxWvScJno/WUMpUcEzcuilWwnSa7moB+BwOGGjTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAIC7sPGE458Kdpk/rGR0/FMe+ylWrtBrI4+VXVsEZeil4MAoGCCqGSM49BAMCA0cAMEQCIF4NFG8yKbZ+/YiN2PxqCEoPIYgFv/2qZlUU4x+SwLCnAiAzH4eI26ovXkSiC5JytoVUEN0AZgaG1x9vcq2AKLOfHw==",
				"RawTBSCertificate": "MIIBw6ADAgECAhEAqIQOa8orqqULAg3UKotCZTAKBggqhkjOPQQDAjBzMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEZMBcGA1UEChMQb3JkZXJlci5iYWFzLnh5ejEcMBoGA1UEAxMTY2Eub3JkZXJlci5iYWFzLnh5ejAeFw0yMDA1MDcwMTM0MzhaFw0zMDA1MDUwMTM0MzhaMF4xCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMSIwIAYDVQQDExlvcmRlcmVyMC5vcmRlcmVyLmJhYXMueHl6MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgWGPnRCAztSButm/CkfDfzWYvukFYAS14bI6Kyyo56FpcJa0nLFa9Jwmej9ZQylRwTNy6KVbCdJruagH4HA4YaNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgLuw8YTjnwp2mT+sZHT8Ux77KVau0Gsjj5VdWwRl6KXg=",
				"RawSubjectPublicKeyInfo": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgWGPnRCAztSButm/CkfDfzWYvukFYAS14bI6Kyyo56FpcJa0nLFa9Jwmej9ZQylRwTNy6KVbCdJruagH4HA4YQ==",
				"RawSubject": "MF4xCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMSIwIAYDVQQDExlvcmRlcmVyMC5vcmRlcmVyLmJhYXMueHl6",
				"RawIssuer": "MHMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmRlcmVyLmJhYXMueHl6MRwwGgYDVQQDExNjYS5vcmRlcmVyLmJhYXMueHl6",
				"Signature": "MEQCIF4NFG8yKbZ+/YiN2PxqCEoPIYgFv/2qZlUU4x+SwLCnAiAzH4eI26ovXkSiC5JytoVUEN0AZgaG1x9vcq2AKLOfHw==",
				"SignatureAlgorithm": 10,
				"PublicKeyAlgorithm": 3,
				"PublicKey": {
					"Curve": {
						"P": 115792089210356248762697446949407573530086143415290314195533631308867097853951,
						"N": 115792089210356248762697446949407573529996955224135760342422259061068512044369,
						"B": 41058363725152142129326129780047268409114441015993725554835256314039467401291,
						"Gx": 48439561293906451759052585252797914202762949526041747995844080717082404635286,
						"Gy": 36134250956749795798585127919587881956611106672985015071877198253568414405109,
						"BitSize": 256,
						"Name": "P-256"
					},
					"X": 58520732816702070349095293814958830062554359573759139005892235715231473002401,
					"Y": 47691776103742245732617629914324541929131401205796413962283138087879588198497
				},
				"Version": 3,
				"SerialNumber": 223995978970909065685438874041224741477,
				"Issuer": {
					"Country": [
						"US"
					],
					"Organization": [
						"orderer.baas.xyz"
					],
					"OrganizationalUnit": null,
					"Locality": [
						"San Francisco"
					],
					"Province": [
						"California"
					],
					"StreetAddress": null,
					"PostalCode": null,
					"SerialNumber": "",
					"CommonName": "ca.orderer.baas.xyz",
					"Names": [{
							"Type": [
								2,
								5,
								4,
								6
							],
							"Value": "US"
						},
						{
							"Type": [
								2,
								5,
								4,
								8
							],
							"Value": "California"
						},
						{
							"Type": [
								2,
								5,
								4,
								7
							],
							"Value": "San Francisco"
						},
						{
							"Type": [
								2,
								5,
								4,
								10
							],
							"Value": "orderer.baas.xyz"
						},
						{
							"Type": [
								2,
								5,
								4,
								3
							],
							"Value": "ca.orderer.baas.xyz"
						}
					],
					"ExtraNames": null
				},
				"Subject": {
					"Country": [
						"US"
					],
					"Organization": null,
					"OrganizationalUnit": null,
					"Locality": [
						"San Francisco"
					],
					"Province": [
						"California"
					],
					"StreetAddress": null,
					"PostalCode": null,
					"SerialNumber": "",
					"CommonName": "orderer0.orderer.baas.xyz",
					"Names": [{
							"Type": [
								2,
								5,
								4,
								6
							],
							"Value": "US"
						},
						{
							"Type": [
								2,
								5,
								4,
								8
							],
							"Value": "California"
						},
						{
							"Type": [
								2,
								5,
								4,
								7
							],
							"Value": "San Francisco"
						},
						{
							"Type": [
								2,
								5,
								4,
								3
							],
							"Value": "orderer0.orderer.baas.xyz"
						}
					],
					"ExtraNames": null
				},
				"NotBefore": "2020-05-07T01:34:38Z",
				"NotAfter": "2030-05-05T01:34:38Z",
				"KeyUsage": 1,
				"Extensions": [{
						"Id": [
							2,
							5,
							29,
							15
						],
						"Critical": true,
						"Value": "AwIHgA=="
					},
					{
						"Id": [
							2,
							5,
							29,
							19
						],
						"Critical": true,
						"Value": "MAA="
					},
					{
						"Id": [
							2,
							5,
							29,
							35
						],
						"Critical": false,
						"Value": "MCKAIC7sPGE458Kdpk/rGR0/FMe+ylWrtBrI4+VXVsEZeil4"
					}
				],
				"ExtraExtensions": null,
				"UnhandledCriticalExtensions": null,
				"ExtKeyUsage": null,
				"UnknownExtKeyUsage": null,
				"BasicConstraintsValid": true,
				"IsCA": false,
				"MaxPathLen": -1,
				"MaxPathLenZero": false,
				"SubjectKeyId": null,
				"AuthorityKeyId": "Luw8YTjnwp2mT+sZHT8Ux77KVau0Gsjj5VdWwRl6KXg=",
				"OCSPServer": null,
				"IssuingCertificateURL": null,
				"DNSNames": null,
				"EmailAddresses": null,
				"IPAddresses": null,
				"URIs": null,
				"PermittedDNSDomainsCritical": false,
				"PermittedDNSDomains": null,
				"ExcludedDNSDomains": null,
				"PermittedIPRanges": null,
				"ExcludedIPRanges": null,
				"PermittedEmailAddresses": null,
				"ExcludedEmailAddresses": null,
				"PermittedURIDomains": null,
				"ExcludedURIDomains": null,
				"CRLDistributionPoints": null,
				"PolicyIdentifiers": null
			},
			"nonce": "0YRJ57plpqyU35VqAcfSXsaFSwJMOZf/"
		},
		"signature": "MEUCIQD4YW+G4NT/tP67eg5LATXn+JeDzieejBY7dLDxBPjGfAIgGfd5kQ0zmjGI5vCgmMcNy6Zd5HsV9+mFLE4JRh3N7Rs="
	},
	"last_config_block_number": {
		"signature_data": {
			"signature_header": {
				"Certificate": {
					"Raw": "MIICHDCCAcOgAwIBAgIRAKiEDmvKK6qlCwIN1CqLQmUwCgYIKoZIzj0EAwIwczELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xGTAXBgNVBAoTEG9yZGVyZXIuYmFhcy54eXoxHDAaBgNVBAMTE2NhLm9yZGVyZXIuYmFhcy54eXowHhcNMjAwNTA3MDEzNDM4WhcNMzAwNTA1MDEzNDM4WjBeMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEiMCAGA1UEAxMZb3JkZXJlcjAub3JkZXJlci5iYWFzLnh5ejBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABIFhj50QgM7UgbrZvwpHw381mL7pBWAEteGyOissqOehaXCWtJyxWvScJno/WUMpUcEzcuilWwnSa7moB+BwOGGjTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAIC7sPGE458Kdpk/rGR0/FMe+ylWrtBrI4+VXVsEZeil4MAoGCCqGSM49BAMCA0cAMEQCIF4NFG8yKbZ+/YiN2PxqCEoPIYgFv/2qZlUU4x+SwLCnAiAzH4eI26ovXkSiC5JytoVUEN0AZgaG1x9vcq2AKLOfHw==",
					"RawTBSCertificate": "MIIBw6ADAgECAhEAqIQOa8orqqULAg3UKotCZTAKBggqhkjOPQQDAjBzMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEZMBcGA1UEChMQb3JkZXJlci5iYWFzLnh5ejEcMBoGA1UEAxMTY2Eub3JkZXJlci5iYWFzLnh5ejAeFw0yMDA1MDcwMTM0MzhaFw0zMDA1MDUwMTM0MzhaMF4xCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMSIwIAYDVQQDExlvcmRlcmVyMC5vcmRlcmVyLmJhYXMueHl6MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgWGPnRCAztSButm/CkfDfzWYvukFYAS14bI6Kyyo56FpcJa0nLFa9Jwmej9ZQylRwTNy6KVbCdJruagH4HA4YaNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgLuw8YTjnwp2mT+sZHT8Ux77KVau0Gsjj5VdWwRl6KXg=",
					"RawSubjectPublicKeyInfo": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgWGPnRCAztSButm/CkfDfzWYvukFYAS14bI6Kyyo56FpcJa0nLFa9Jwmej9ZQylRwTNy6KVbCdJruagH4HA4YQ==",
					"RawSubject": "MF4xCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMSIwIAYDVQQDExlvcmRlcmVyMC5vcmRlcmVyLmJhYXMueHl6",
					"RawIssuer": "MHMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmRlcmVyLmJhYXMueHl6MRwwGgYDVQQDExNjYS5vcmRlcmVyLmJhYXMueHl6",
					"Signature": "MEQCIF4NFG8yKbZ+/YiN2PxqCEoPIYgFv/2qZlUU4x+SwLCnAiAzH4eI26ovXkSiC5JytoVUEN0AZgaG1x9vcq2AKLOfHw==",
					"SignatureAlgorithm": 10,
					"PublicKeyAlgorithm": 3,
					"PublicKey": {
						"Curve": {
							"P": 115792089210356248762697446949407573530086143415290314195533631308867097853951,
							"N": 115792089210356248762697446949407573529996955224135760342422259061068512044369,
							"B": 41058363725152142129326129780047268409114441015993725554835256314039467401291,
							"Gx": 48439561293906451759052585252797914202762949526041747995844080717082404635286,
							"Gy": 36134250956749795798585127919587881956611106672985015071877198253568414405109,
							"BitSize": 256,
							"Name": "P-256"
						},
						"X": 58520732816702070349095293814958830062554359573759139005892235715231473002401,
						"Y": 47691776103742245732617629914324541929131401205796413962283138087879588198497
					},
					"Version": 3,
					"SerialNumber": 223995978970909065685438874041224741477,
					"Issuer": {
						"Country": [
							"US"
						],
						"Organization": [
							"orderer.baas.xyz"
						],
						"OrganizationalUnit": null,
						"Locality": [
							"San Francisco"
						],
						"Province": [
							"California"
						],
						"StreetAddress": null,
						"PostalCode": null,
						"SerialNumber": "",
						"CommonName": "ca.orderer.baas.xyz",
						"Names": [{
								"Type": [
									2,
									5,
									4,
									6
								],
								"Value": "US"
							},
							{
								"Type": [
									2,
									5,
									4,
									8
								],
								"Value": "California"
							},
							{
								"Type": [
									2,
									5,
									4,
									7
								],
								"Value": "San Francisco"
							},
							{
								"Type": [
									2,
									5,
									4,
									10
								],
								"Value": "orderer.baas.xyz"
							},
							{
								"Type": [
									2,
									5,
									4,
									3
								],
								"Value": "ca.orderer.baas.xyz"
							}
						],
						"ExtraNames": null
					},
					"Subject": {
						"Country": [
							"US"
						],
						"Organization": null,
						"OrganizationalUnit": null,
						"Locality": [
							"San Francisco"
						],
						"Province": [
							"California"
						],
						"StreetAddress": null,
						"PostalCode": null,
						"SerialNumber": "",
						"CommonName": "orderer0.orderer.baas.xyz",
						"Names": [{
								"Type": [
									2,
									5,
									4,
									6
								],
								"Value": "US"
							},
							{
								"Type": [
									2,
									5,
									4,
									8
								],
								"Value": "California"
							},
							{
								"Type": [
									2,
									5,
									4,
									7
								],
								"Value": "San Francisco"
							},
							{
								"Type": [
									2,
									5,
									4,
									3
								],
								"Value": "orderer0.orderer.baas.xyz"
							}
						],
						"ExtraNames": null
					},
					"NotBefore": "2020-05-07T01:34:38Z",
					"NotAfter": "2030-05-05T01:34:38Z",
					"KeyUsage": 1,
					"Extensions": [{
							"Id": [
								2,
								5,
								29,
								15
							],
							"Critical": true,
							"Value": "AwIHgA=="
						},
						{
							"Id": [
								2,
								5,
								29,
								19
							],
							"Critical": true,
							"Value": "MAA="
						},
						{
							"Id": [
								2,
								5,
								29,
								35
							],
							"Critical": false,
							"Value": "MCKAIC7sPGE458Kdpk/rGR0/FMe+ylWrtBrI4+VXVsEZeil4"
						}
					],
					"ExtraExtensions": null,
					"UnhandledCriticalExtensions": null,
					"ExtKeyUsage": null,
					"UnknownExtKeyUsage": null,
					"BasicConstraintsValid": true,
					"IsCA": false,
					"MaxPathLen": -1,
					"MaxPathLenZero": false,
					"SubjectKeyId": null,
					"AuthorityKeyId": "Luw8YTjnwp2mT+sZHT8Ux77KVau0Gsjj5VdWwRl6KXg=",
					"OCSPServer": null,
					"IssuingCertificateURL": null,
					"DNSNames": null,
					"EmailAddresses": null,
					"IPAddresses": null,
					"URIs": null,
					"PermittedDNSDomainsCritical": false,
					"PermittedDNSDomains": null,
					"ExcludedDNSDomains": null,
					"PermittedIPRanges": null,
					"ExcludedIPRanges": null,
					"PermittedEmailAddresses": null,
					"ExcludedEmailAddresses": null,
					"PermittedURIDomains": null,
					"ExcludedURIDomains": null,
					"CRLDistributionPoints": null,
					"PolicyIdentifiers": null
				},
				"nonce": "aMUTuJ9gN0z93Yv2n9lNP77zVvoszRcU"
			},
			"signature": "MEUCIQCel9JDyshT+EbnuGA9oG8IGVWZiXicwo3+9oXv88WnXQIgXlusYKejZKAoNubzOKwMfB4aUDJvq9q3dtjyXNmKRbE="
		}
	},
	"transaction_filter": "AA==",
	"orderer_kafka_metadata": {
		"last_offset_persisted": 432345564227567616
	}
    }
}
```          

## query block hight                                               
URL：http://ip:port/api/v1/queryblock                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/queryblock                                                       
request parameter:            
```json    
{
	"BlockChainId": "1m4MJ3U0avy8VLWkQeXCtgN8M3384TAd",
	"OrgName": "Org1",
	"ChannelId": "mychannel"
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",     
     "data": {
        "Height": 5,
        "CurrentBlockHash": "9e96727d60e134b9e2b3f9147451a7bef8df6087aa947d43b28205fbcf4104d7",
        "PreviousBlockHash": "d329e68c1a07eafa5bec591a8a6f2d900cccb62ee7b3dc8496636cd5e1996e77"
    }
}
```          

## query block hight and total tx by blockchain id and channel id                                                
URL：http://ip:port/api/v1/blocktx                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/blocktx                                                       
request parameter:            
```json    
{
	"BlockChainId": "1m4MJ3U0avy8VLWkQeXCtgN8M3384TAd",
	"OrgName": "Org1",
	"ChannelId": "mychannel"
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "BlockChainId": "1m4MJ3U0avy8VLWkQeXCtgN8M3384TAd",
        "ChainnnelId": "mychannel",
        "Height": 2,
        "TxCount": 1
    }
}
```         

## get blockchain by id                    
URL：http://ip:port/api/v1/:blockchainid/blockchain              
METHOD：GET        
INPUT PARA:         
RETURN：json object              
example:          
request:http://localhost:9001/api/v1/1m4MJ3U0avy8VLWkQeXCtgN8M3384TAd/blockchain                 
API RETURN：                      
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "AllianceId": "5bGiWNF0ycI7eTbH",
        "BlockChainId": "1m4MJ3U0avy8VLWkQeXCtgN8M3384TAd",
        "BlockChainName": "test-network",
        "BlockChainType": "fabric",
		"ClusterId": "test-cluster1",
        "Version": "v1",
        "Status": 0
    }
}
```     

## query block and tx in this block                                                    
URL：http://ip:port/api/v1/queryblocktx                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/queryblocktx                                                       
request parameter:            
```json    
{
	"Start": 0,
	"End": 3,
	"BlockChainId": "4HBA0t8aBGLPw9yXkc7FEz7ZqH9tCV0G",
	"OrgName": "Org1",
	"ChannelId": "baaschannel"
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "BlockId": 1,
            "BlockHash": "peRdcWp5sKSdUXU7JhH4vD8Uu6WghNu7yI9+wkwvlCs=",
            "PreHash": "uvXcGRBtnYN+O/MXwYSGEZknsMWQeZCTwPHn9ZKxNog=",
            "Timestamp": "2006-01-02 15:04:05",
            "Txs": [
                {
                    "TxId": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
                    "Timestamp": "2006-01-02 15:04:05",
                    "Signature": "MEUCIQDKHuF2gGEvcevbshpKVuFsFpBHrwwMsn2i3B/zJuKHOgIgKKcL+7EygQowKU1SJ6H72VDVaAnVmcsSugm9i5V1mwU=",
                    "BlockId": 1
                }
            ]
        },
        {
            "BlockId": 2,
            "BlockHash": "9iJMOqH85FojoWQPEcL+Q/WYOFyrXPDtdSXq9baAPl8=",
            "PreHash": "yq80LDrGdrSocsshZx6/E/m5aMYbqFTtlLqhg2zUYso=",
            "Timestamp": "2006-01-02 15:04:05",
            "Txs": [
                {
                    "TxId": "4a0ec826f1d58747db8a0912fd182beb809b0d36017cf354d39fc3fda6c012d3",
                    "Timestamp": "2006-01-02 15:04:05",
                    "Signature": "MEUCIQDW4MygQ+OD8CmNIUP96Zkxd7jQmv5zyt6/apHLfPksYAIge6QUgHO5bLuYvECRY2PO1SBQSTB+r2LlTk50tR7pJi8=",
                    "BlockId": 2
                }
            ]
        }
    ]
}
```       

## query tx from db                                                        
URL：http://ip:port/api/v1/querytx                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/querytx                                                       
request parameter:            
```json    
{
	"TxId": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
	"BlockChainId": "4HBA0t8aBGLPw9yXkc7FEz7ZqH9tCV0G",
	"OrgName": "Org1",
	"ChannelId": "baaschannel"
}
``` 
API RETURN:               
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "TxId": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
        "Timestamp": "2020-07-28 11:02:45",
        "Signature": "MEUCIQDKHuF2gGEvcevbshpKVuFsFpBHrwwMsn2i3B/zJuKHOgIgKKcL+7EygQowKU1SJ6H72VDVaAnVmcsSugm9i5V1mwU=",
        "BlockId": 1
    }
}
```       

## search block or tx from db                                                        
URL：http://ip:port/api/v1/search                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/search                                                       
request parameter:            
```json    
{
	"Key": "1",
	"BlockChainId": "4HBA0t8aBGLPw9yXkc7FEz7ZqH9tCV0G",
	"OrgName": "Org1",
	"ChannelId": "baaschannel"
}
``` 
API RETURN:               
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "BlockId": 1,
        "BlockHash": "peRdcWp5sKSdUXU7JhH4vD8Uu6WghNu7yI9+wkwvlCs=",
        "PreHash": "uvXcGRBtnYN+O/MXwYSGEZknsMWQeZCTwPHn9ZKxNog=",
        "Timestamp": "2020-07-28 11:02:45",
        "Txs": [
            {
                "TxId": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
                "Timestamp": "2020-07-28 11:02:45",
                "Signature": "MEUCIQDKHuF2gGEvcevbshpKVuFsFpBHrwwMsn2i3B/zJuKHOgIgKKcL+7EygQowKU1SJ6H72VDVaAnVmcsSugm9i5V1mwU=",
                "BlockId": 1
            }
        ]
    }
}
```       

request parameter:            
```json    
{
	"Key": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
	"BlockChainId": "4HBA0t8aBGLPw9yXkc7FEz7ZqH9tCV0G",
	"OrgName": "Org1",
	"ChannelId": "baaschannel"
}
``` 
API RETURN:               
```json     
{
    "code": 0,
    "message": "success",
    "data": {
        "TxId": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
        "Timestamp": "2020-07-28 11:02:45",
        "Signature": "MEUCIQDKHuF2gGEvcevbshpKVuFsFpBHrwwMsn2i3B/zJuKHOgIgKKcL+7EygQowKU1SJ6H72VDVaAnVmcsSugm9i5V1mwU=",
        "BlockId": 1
    }
}
```       

## search pass tx by day                                                         
URL：http://ip:port/api/v1/querypasstx                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/querypasstx                                                       
request parameter:            
```json    
{
	"BlockChainId": "4HBA0t8aBGLPw9yXkc7FEz7ZqH9tCV0G",
	"OrgName": "Org1",
	"ChannelId": "baaschannel",
	"Day": -1 //only -1 ~ -14,pass 14 day tx
}
``` 
API RETURN:               
```json     
{
    "code": 0,
    "message": "success",
    "data": [
		{
        "TxId": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
        "Timestamp": "2020-07-28 11:02:45",
        "Signature": "MEUCIQDKHuF2gGEvcevbshpKVuFsFpBHrwwMsn2i3B/zJuKHOgIgKKcL+7EygQowKU1SJ6H72VDVaAnVmcsSugm9i5V1mwU=",
        "BlockId": 1
    	}
	]
}
```       

## search sure day tx by start and end time                                                            
URL：http://ip:port/api/v1/querydaytx                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/querydaytx                                                       
request parameter:            
```json    
{
	"BlockChainId": "4HBA0t8aBGLPw9yXkc7FEz7ZqH9tCV0G",
	"OrgName": "Org1",
	"ChannelId": "baaschannel",
	"StartTime": "2006-01-02 00:00:00",
	"EndTime": "2006-01-02 23:59:59"
}
``` 
API RETURN:               
```json     
{
    "code": 0,
    "message": "success",
    "data": [
		{
        "TxId": "f61c9653b318f4b4d78b5ce6782341e159cff67358528f8ab336dc5040dd9be4",
        "Timestamp": "2020-07-28 11:02:45",
        "Signature": "MEUCIQDKHuF2gGEvcevbshpKVuFsFpBHrwwMsn2i3B/zJuKHOgIgKKcL+7EygQowKU1SJ6H72VDVaAnVmcsSugm9i5V1mwU=",
        "BlockId": 1
    	}
	]
}
```       

## search every day tx by start and days                                                            
URL：http://ip:port/api/v1/queryeverydaytx                                          
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/queryeverydaytx                                                       
request parameter:            
```json    
{
	"StartTime": "2020-07-30 00:00:00",
	"Days": 1, //1 ~ 14
	"BlockChainId": "brEMt2bNgFKKNrJy6Q9I3aseamaskpVY",
	"OrgName": "Org",
	"ChannelId": "mychannel"
}
``` 
API RETURN:               
```json     
{
    "code": 0,
    "message": "success",
	"data": {
		"Txs":[
			{
        		"TimeStamp": "2020-07-30 00:00:00",
       			 "TxCount": 14
			},
			{
        		"TimeStamp": "2020-07-31 00:00:00",
        		"TxCount": 9
    		}
		],
		"Max":14,
		"Min":9
	}
}
```       