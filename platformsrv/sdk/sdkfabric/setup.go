
package sdkfabric

import (
	"fmt"
	"strings"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile      string
	OrdererID       string
	ChannelId		string
	OrgAdmin        string
	OrgName         string
	Peers           []string
	Sdk             *fabsdk.FabricSDK 
	ChClient        *channel.Client
	netAdmin        *resmgmt.Client
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
	EndorsePolicy		string
	InitArgs            []string 
}

type UserClient struct {
	Secret string
	IdentityTypeUser string
	RegistrarUserName string
	RegistrarPassword string
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error { 
	if setup.initialized {
		logger.Debug("sdk already inited")
		return nil
	}
	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		logger.Errorf("failed to create SDK,err=%v",err)
		return fmt.Errorf("failed to create SDK,err=%v",err) 
	}
	setup.Sdk = sdk
	resourceManagerClientContext := setup.Sdk.Context(fabsdk.WithUser(setup.OrgAdmin),fabsdk.WithOrg(setup.OrgName))
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)  
	if err != nil {
		logger.Errorf("failed to create channel management client from Admin identity,err=%v",err) 
		return fmt.Errorf("failed to create channel management client from Admin identity,err=%v",err)
	}
	setup.netAdmin = resMgmtClient 
	logger.Debug("rc client create success")
	setup.initialized = true
	logger.Debug("sdk init success")
	return nil
} 

func (setup *FabricSetup) CreateChannel(ch ChannnelSetup) error {
	 
	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(setup.Sdk.Context(), mspclient.WithOrg(setup.OrgName))
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
	txID, err := setup.netAdmin.SaveChannel(req,resmgmt.WithOrdererEndpoint(setup.OrdererID))
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
	//ccPolicy := cauthdsl.SignedByAnyMember([]string{cc.InitOrg})

	var endorse string
	if cc.EndorsePolicy == "" {
		endorse = "AND('" + setup.OrgName + "MSP.member')"
	}
	endorse = cc.EndorsePolicy  
	ccPolicy,err := setup.genPolicy(endorse)
	if err != nil {
		logger.Errorf("InstantiateCC: genPolicy failed,msg=%v",err.Error())
		return fmt.Errorf("InstantiateCC: genPolicy failed,msg=%v",err.Error())
	}
	args := packArgs(cc.InitArgs)
	req := resmgmt.InstantiateCCRequest{
		Name: cc.ChainCodeID,
		Path: cc.ChaincodeGoPath,
		Version: cc.ChaincodeVersion,
		Args: args,
		Policy: ccPolicy,
	}
	reqPeers := resmgmt.WithTargetEndpoints(setup.Peers...)

	resp, err := setup.netAdmin.InstantiateCC(ch.ChannelID,req,reqPeers)
	if err != nil || resp.TransactionID == "" {
		logger.Errorf("failed to instantiate the chaincode,err=%v",err)
		return fmt.Errorf("failed to instantiate the chaincode,err=%v",err)
	}
	logger.Debug("Chaincode instantiate success") 
	return nil
}

func (setup *FabricSetup) UpgradeCC(ch ChannnelSetup,cc ChaincodeSetup) error {
	//ccPolicy := cauthdsl.SignedByAnyMember([]string{cc.InitOrg})
	var endorse string
	if cc.EndorsePolicy == "" {
		endorse = "AND('" + setup.OrgName + "MSP.member')"
	}
	endorse = cc.EndorsePolicy  
	ccPolicy,err := setup.genPolicy(endorse)
	if err != nil {
		logger.Errorf("UpgradeCC: genPolicy failed,msg=%v",err.Error())
		return fmt.Errorf("UpgradeCC: genPolicy failed,msg=%v",err.Error())
	}
	args := packArgs(cc.InitArgs)
	req := resmgmt.UpgradeCCRequest{
		Name: cc.ChainCodeID,
		Path: cc.ChaincodeGoPath,
		Version: cc.ChaincodeVersion,
		Args: args,
		Policy: ccPolicy,
	}
	reqPeers := resmgmt.WithTargetEndpoints(setup.Peers...)

	resp, err := setup.netAdmin.UpgradeCC(ch.ChannelID,req,reqPeers)
	if err != nil || resp.TransactionID == "" {
		logger.Errorf("failed to UpgradeCC the chaincode,err=%v",err)
		return fmt.Errorf("failed to UpgradeCC the chaincode,err=%v",err)
	}
	logger.Debug("Chaincode UpgradeCC success") 
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
	setup.Sdk.Close() 
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}

func (setup *FabricSetup) GetRegisteredUser(userName string, orgName string,secret string,identityTypeUser string) (string, error) {
	ctxProvider := setup.Sdk.Context()
	client,err := mspclient.New(ctxProvider,mspclient.WithOrg(orgName))
	if err != nil {
		logger.Errorf("GetRegisteredUser: failed to create msp client: %s", err.Error())
		return "",fmt.Errorf("GetRegisteredUser: failed to create msp client: %s", err.Error())
	}
	enrollSecret,err := client.Register(&mspclient.RegistrationRequest{
					Name:        userName,
					Type:        identityTypeUser,
					Affiliation: "org1" + ".department1",
					Secret:      secret,
					CAName: strings.ToLower(orgName + "-ca"),
				})
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg,"is already registered") {
			logger.Debug("GetRegisteredUser: %s already registered",userName)
			return "",nil
		}
		logger.Errorf("GetRegisteredUser: %s", err.Error())
		return "",fmt.Errorf("GetRegisteredUser: failed to Register")
	}
	err = client.Enroll(userName,mspclient.WithSecret(enrollSecret))
	if err != nil {
		logger.Errorf("GetRegisteredUser: enroll %s failed: %s", userName, err)
		return "",fmt.Errorf("failed " + err.Error())
	}
	logger.Debug(userName + " enrolled Successfully")
	return "",nil
}

func (setup *FabricSetup) GetUserClient(channelName string,userName string,orgName string) (*channel.Client,error){
	clientChannelContext := setup.Sdk.ChannelContext(channelName,fabsdk.WithUser(userName),fabsdk.WithOrg(orgName))
	client,err := channel.New(clientChannelContext)
	if err != nil {
		logger.Errorf("GetUserClient: failed to create channel client: %s\n",err)
		return nil,fmt.Errorf("GetUserClient: failed to create channel client: %s\n",err)
	}
	setup.ChClient = client
	logger.Debug("GetUserClient success,user=%v  orgName=%v",userName,orgName) 
	return client,nil
}

func (setup *FabricSetup) SetupCA(orgName string) error {
	ctxProvider := setup.Sdk.Context()
	mpClient, err := mspclient.New(ctxProvider,mspclient.WithOrg(orgName))
	err = mpClient.Enroll("admin",mspclient.WithSecret("adminpw"))
	if err != nil {
		logger.Errorf("SetupCA: enroll registrar failed: %s", err)
		return fmt.Errorf("SetupCA: enroll registrar failed")
	}
	return nil
}


