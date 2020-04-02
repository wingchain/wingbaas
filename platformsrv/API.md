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
	"BlockChainName": "chainnetwork1",
	"BlockChainType": "fabric",
	"DeployCfg":{ 
		"DeployNetCfg": {
			"OrdererOrgs": [
			{
				"Name": "Orderer",
				"Domain": "orderer.baas.xyz",
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
				"Specs": [
					{
						"Hostname": "peer0"
					},
					{
						"Hostname": "peer1"
					}
				],
				"Users": {
					"Count": 1
				}
			},
			{
				"Name": "Org2",
				"Domain": "Org2.fabric.baas.xyz",
				"Specs": [
					{
						"Hostname": "peer0"
					},
					{
						"Hostname": "peer1"
					}
				],
				"Users": {
					"Count": 1
				}
			}
		]
		},
		"DeployType": "KAFKA_FABRIC",
		"Version": "1.3.0",
		"CryptoType": "ECDSA",
		"ClusterId": "cluster1"
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
