
package public

const (
	SYS_CHANNEL        				string = "testchainid"
	BLOCK_CHAIN_TYPE_FABRIC 		string = "fabric"
	SOLO_FABRIC     				string = "SOLO_FABRIC"
	KAFKA_FABRIC    				string = "KAFKA_FABRIC"
	RAFT_FABRIC     				string = "RAFT_FABRIC"
	ZOOK_COUNT      				int    = 3
	KAFKA_COUNT     				int    = 4 
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
	ToolsDeployNode string `json:"ToolsDeployNode"`
}

type DeployPara struct { 
	DeployNetCfg        DeployNetConfig    	`json:"DeployNetCfg"` 
	DeployType       	string          	`json:"DeployType"`
	Version    		 	string          	`json:"Version"`
	CryptoType       	string          	`json:"CryptoType"`
	ClusterId        	string          	`json:"ClusterId"` 
	
}

type ServiceNodePortSt struct { 
	ServerName	string
	NodePort 	string
}

type AddOrgConfig struct {
	BlockChainId   string      `json:"BlockChainId"`
	PeerOrgs       []OrgSpec   `json:"PeerOrgs"`
}

type ChannelList struct {
	Channels []struct {
		ChannelID string `json:"channel_id"`
	} `json:"channels"`
}

type ChannelTxSt struct {
	BlockChainId    string `json:"BlockChainId"`
	ChainnnelId     string `json:"ChainnnelId"`
	Height      	uint64 `json:"Height"`
	TxCount     	uint64 `json:"TxCount"`
}

type TxRecordSt struct {
	BlockChainId string `json:"BlockChainId"`
	ChTx []ChannelTxSt  `json:"ChTx"` 
}

type CCSt struct {
	BlockChainId    string `json:"BlockChainId"`
	BlockChainName  string `json:"BlockChainName"`
	CCName			string `json:"CCName"`
	CCVersion		string `json:"CCVersion"`
	CreateTime		string `json:"CreateTime"`
	UpdateTime		string `json:"UpdateTime"`
}

type ChannelCCSt struct {
	ChainnnelId     string `json:"ChainnnelId"`
	CCRecord	    []CCSt `json:"CCRecord"`
}

type CCRecordSt struct {
	BlockChainId   string `json:"BlockChainId"`
	BlockChainName string `json:"BlockChainName"`
	ChCC []ChannelCCSt  `json:"ChCC"`
}

