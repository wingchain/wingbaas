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
        },
        {
            "metadata": {
                "name": "kube-node-lease",
                "selfLink": "/api/v1/namespaces/kube-node-lease",
                "uid": "d3747e6f-455f-47ad-abb9-ecbc67f8746f",
                "resourceVersion": "7",
                "creationTimestamp": "2020-03-26T09:07:01Z"
            },
            "spec": {
                "finalizers": [
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "metadata": {
                "name": "kube-public",
                "selfLink": "/api/v1/namespaces/kube-public",
                "uid": "ce324640-f2a4-4d2f-8638-b102f9a15b96",
                "resourceVersion": "6",
                "creationTimestamp": "2020-03-26T09:07:01Z"
            },
            "spec": {
                "finalizers": [
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "metadata": {
                "name": "kube-system",
                "selfLink": "/api/v1/namespaces/kube-system",
                "uid": "7d013cdd-cc98-4a0b-a180-7cd227653dd4",
                "resourceVersion": "5",
                "creationTimestamp": "2020-03-26T09:07:01Z"
            },
            "spec": {
                "finalizers": [
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "metadata": {
                "name": "kubernetes-dashboard",
                "selfLink": "/api/v1/namespaces/kubernetes-dashboard",
                "uid": "8ff22b3c-655b-4d09-92dc-ec00da60d7d8",
                "resourceVersion": "28389",
                "creationTimestamp": "2020-03-26T12:30:35Z"
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
