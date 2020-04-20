
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
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	//"github.com/pkg/errors" 
)

const (
	DefaultChannel	string	=	"mychannel" 
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile      string
	OrgID           string 
	OrdererID       string
	ChannelID       string
	ChainCodeID     string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	UserName        string
	client          *channel.Client
	admin           *resmgmt.Client
	sdk             *fabsdk.FabricSDK
	event           *event.Client
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
	logger.Debug("sdk=")
	logger.Debug(sdk)
	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	if err != nil {
		logger.Errorf("failed to load Admin identity,err=%v",err)
		return fmt.Errorf("failed to load Admin identity,err=%v",err)
	}
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext) 
	if err != nil {
		logger.Errorf("failed to create channel management client from Admin identity,err=%v",err)
		return fmt.Errorf("failed to create channel management client from Admin identity,err=%v",err)
	}
	setup.admin = resMgmtClient

	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		logger.Errorf("failed to create MSP client,err=%v",err)
		return fmt.Errorf("failed to create MSP client,err=%v",err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin) 
	if err != nil {
		logger.Errorf("failed to get admin signing identity,err=%v",err)
		return fmt.Errorf("failed to get admin signing identity,err=%v",err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfigPath: setup.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := setup.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		logger.Errorf("failed to save channel,err=%v",err)
		return fmt.Errorf("failed to save channel,err=%v",err)
	}
	logger.Debug("Channel created")
	// Make admin user join the previously created channel
	if err = setup.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
		logger.Errorf("failed to make admin join channel,err=%v",err)
		return fmt.Errorf("failed to make admin join channel,err=%v",err)
	}
	logger.Debug("Channel joined")
	logger.Debug("Initialization Successful")
	setup.initialized = true
	return nil
}

func (setup *FabricSetup) InstallAndInstantiateCC() error {

	// Create the chaincode package that will be sent to the peers
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		logger.Errorf("failed to create chaincode package,err=%v",err)
		return fmt.Errorf("failed to create chaincode package,err=%v",err)
	}
	logger.Debug("ccPkg created")

	// Install cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "0", Package: ccPkg}
	_, err = setup.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Errorf("failed to install chaincode,err=%v",err)
		return fmt.Errorf("failed to install chaincode,err=%v",err)
	}
	logger.Debug("Chaincode installed")

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.chainhero.io"}) 

	resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodeGoPath, Version: "0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		logger.Errorf("failed to instantiate the chaincode,err=%v",err)
		return fmt.Errorf("failed to instantiate the chaincode,err=%v",err)
	}
	logger.Debug("Chaincode instantiated") 

	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		logger.Errorf("failed to create new channel client,err=%v",err)
		return fmt.Errorf("failed to create new channel client,err=%v",err)
	}
	logger.Debug("Channel client created")

	// Creation of the client which will enables access to our channel events
	setup.event, err = event.New(clientContext)
	if err != nil {
		logger.Errorf("failed to create new event client,err=%v",err)
		return fmt.Errorf("failed to create new event client,err=%v",err)
	}
	logger.Debug("Event client created")
	logger.Debug("Chaincode Installation & Instantiation Successful")
	return nil
}

func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close() 
}
