
package public

const (
	SOLO_FABRIC     string = "SOLO_FABRIC"
	KAFKA_FABRIC    string = "KAFKA_FABRIC"
	RAFT_FABRIC     string = "RAFT_FABRIC"
	ZOOK_COUNT      int    = 3
	KAFKA_COUNT     int    = 4
)

type NodeSpec struct {
	Hostname string   `json:"Hostname"`
}
  
type UsersSpec struct {
    Count int `json:"Count"`
}
  
type OrgSpec struct { 
    Name          string       `json:"Name"`
	Domain        string       `json:"Domain"`
	DeployNode    string       `json:"DeployNode"`
    Specs         []NodeSpec   `json:"Specs"`
	Users         UsersSpec    `json:"Users"`
}
  
type DeployNetConfig struct {  
    OrdererOrgs []OrgSpec `json:"OrdererOrgs"`
	PeerOrgs    []OrgSpec `json:"PeerOrgs"`
	KafkaDeployNode	string `json:"KafkaDeployNode"`
	ZookeeperDeployNode string `json:"ZookeeperDeployNode"` 
}

type DeployNodeInfo map[string]string
type DeployNodeGroup map[string]DeployNodeInfo

type DeployPara struct {
	DeployNetCfg        DeployNetConfig    	`json:"DeployNetCfg"`  
	DeployHost       	DeployNodeGroup 	`json:"DeployHost"`
	DeployType       	string          	`json:"DeployType"`
	Version    		 	string          	`json:"Version"`
	CryptoType       	string          	`json:"CryptoType"`
	ClusterId        	string          	`json:"ClusterId"` 
}

type ServiceNodePortSt struct { 
	ServerName	string
	NodePort 	string
}


