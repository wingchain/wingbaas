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
	"Args":["transfer","a","b","10"]
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
	"ChainCodeID": "cctest",
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
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": {
        "chaincodes": [
            {
                "name": "cctest",
                "version": "10",
                "path": "cctest10",
                "input": "<nil>",
                "escc": "escc",
                "vscc": "vscc"
            }
        ]
    }
}
```         

## query all installed chaincode list                                     
URL：http://ip:port/api/v1/queryinstalledcc                                 
METHOD：POST   
RETURN：json object           
example:        
request:http://localhost:9001/api/v1/queryinstalledcc                                           
request parameter:            
```json    
{
	"BlockChainId": "6f2Wq0Hy5LSYpSt0yvmTd1XD2XEn3Hdr",
	"OrgName": "Org1",
	"ChannelId": "mychannel" 
}
``` 
API RETURN:                  
```json     
{
    "code": 0,
    "message": "success",    
    "data": "null"
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
    "data": "null"
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

