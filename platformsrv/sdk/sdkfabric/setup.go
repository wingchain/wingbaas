
package sdkfabric

import (
	"fmt"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	//"github.com/pkg/errors" 
)

const (
	DefaultChannel	string	=	"mychannel" 
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile      string
	OrdererID       string
	OrgAdmin        string
	OrgName         string
	UserName        string
	sdk             *fabsdk.FabricSDK
	netAdmin        *resmgmt.Client
	chClient        *channel.Client
	chEventClient   *event.Client
	initialized     bool
}

type ChannnelSetup struct {
	ChannelID       string
	ChannelConfig   string
}

type ChaincodeSetup struct {
	ChainCodeID     	string
	ChaincodeVersion	string
	ChaincodeGoPath 	string
	ChaincodePath   	string
	InitOrg				string
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error { 

	// Add parameters for the initialization
	if setup.initialized {
		logger.Debug("sdk already initialized")
		return nil
	}
	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		logger.Errorf("failed to create SDK,err=%v",err)
		return fmt.Errorf("failed to create SDK,err=%v",err) 
	}
	setup.sdk = sdk
	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)  
	if err != nil {
		logger.Errorf("failed to create channel management client from Admin identity,err=%v",err) 
		return fmt.Errorf("failed to create channel management client from Admin identity,err=%v",err)
	}
	setup.netAdmin = resMgmtClient 
	setup.initialized = true
	logger.Debug("sdk initialize success")
	return nil
} 

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}

func (setup *FabricSetup) CreateChannel(ch ChannnelSetup) error {
	 
	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(setup.sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		logger.Errorf("failed to create MSP client,err=%v",err)
		return fmt.Errorf("failed to create MSP client,err=%v",err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin) 
	if err != nil {
		logger.Errorf("failed to get admin signing identity,err=%v",err)
		return fmt.Errorf("failed to get admin signing identity,err=%v",err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: ch.ChannelID, ChannelConfigPath: ch.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := setup.netAdmin.SaveChannel(req,/*resmgmt.WithRetry(retry.DefaultResMgmtOpts),*/resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		logger.Errorf("failed to save channel,err=%v",err)
		return fmt.Errorf("failed to save channel,err=%v",err)
	}
	logger.Debug("Channel created success")
	return nil
}

func  (setup *FabricSetup)JoinChannel(ch ChannnelSetup)error {
	err := setup.netAdmin.JoinChannel(ch.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil {
		logger.Errorf("failed to make admin join channel,err=%v",err)
		return fmt.Errorf("failed to make admin join channel,err=%v",err)
	}
	logger.Debug("join Channel success")
	return nil
}

func (setup *FabricSetup) InstallCC(cc ChaincodeSetup) error {

	// Create the chaincode package that will be sent to the peers
	ccPkg, err := packager.NewCCPackage(cc.ChaincodePath, cc.ChaincodeGoPath)
	if err != nil {
		logger.Errorf("failed to create chaincode package,err=%v",err)
		return fmt.Errorf("failed to create chaincode package,err=%v",err)
	}
	logger.Debug("ccPkg created")

	// Install cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: cc.ChainCodeID, Path: cc.ChaincodePath, Version: cc.ChaincodeVersion, Package: ccPkg}
	_, err = setup.netAdmin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Errorf("failed to install chaincode,err=%v",err)
		return fmt.Errorf("failed to install chaincode,err=%v",err)
	}
	logger.Debug("Chaincode install success")
	return nil
}

func (setup *FabricSetup) InstantiateCC(ch ChannnelSetup,cc ChaincodeSetup) error { 
	// Set up chaincode policy
	//endorser := "OR('Org1MSP.member','Org2MSP.member')"
	
	ccPolicy := cauthdsl.SignedByAnyMember([]string{cc.InitOrg})
	args := packArgs([]string{"init", "a", "100", "b", "200"})
	req := resmgmt.InstantiateCCRequest{
		Name: cc.ChainCodeID,
		Path: cc.ChaincodeGoPath,
		Version: cc.ChaincodeVersion,
		Args:    args,
		Policy:  ccPolicy,
	}
	peers := []string{"peer0-org1.Org1.fabric.baas.xyz","peer1-org1.Org1.fabric.baas.xyz"}
	reqPeers := resmgmt.WithTargetEndpoints(peers...)

	resp, err := setup.netAdmin.InstantiateCC(ch.ChannelID,req,reqPeers)
	if err != nil || resp.TransactionID == "" {
		logger.Errorf("failed to instantiate the chaincode,err=%v",err)
		return fmt.Errorf("failed to instantiate the chaincode,err=%v",err)
	}
	logger.Debug("Chaincode instantiate success") 

	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(ch.ChannelID, fabsdk.WithUser(setup.UserName))
	setup.chClient, err = channel.New(clientContext)
	if err != nil {
		logger.Errorf("failed to create new channel client,err=%v",err)
		return fmt.Errorf("failed to create new channel client,err=%v",err)
	}
	logger.Debug("Channel client create success")
	event, err := event.New(clientContext)
	if err != nil {
		logger.Debug("failed to create new channel event client")
		return fmt.Errorf("failed to create new channel event client")
	}
	setup.chEventClient = event
	logger.Debug("channel event client create success")
	return nil
}

func (setup *FabricSetup) genPolicy(p string) (*common.SignaturePolicyEnvelope, error) {
	if p == "ANY" {
		envelop := cauthdsl.SignedByAnyMember([]string{setup.OrgName})
		return envelop,nil
	}
	return cauthdsl.FromString(p)
}

func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close() 
}
