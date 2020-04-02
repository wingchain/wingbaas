
package txgenerate

import (
	"fmt"
	"io/ioutil"
	"os"
	"bufio"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/wingbaas/platformsrv/settings/fabric/txgenerate/encoder" 
	genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/fabric/common/tools/protolator"
	cb "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/wingbaas/platformsrv/logger"
)

type AnchorPeer struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
}

type Organization struct {
	Name        string `yaml:"Name"`
	ID          string `yaml:"ID"`
	MSPDir      string `yaml:"MSPDir"`
	AnchorPeers []AnchorPeer `yaml:"AnchorPeers,omitempty"`
}

type Kafkas struct {
	Brokers []string `yaml:"Brokers"`
}

type Orderers struct {
	OrdererType  string   `yaml:"OrdererType"`
	Addresses    []string `yaml:"Addresses"`
	BatchTimeout string   `yaml:"BatchTimeout"`
	BatchSize    struct {
		MaxMessageCount   int    `yaml:"MaxMessageCount"`
		AbsoluteMaxBytes  string `yaml:"AbsoluteMaxBytes"`
		PreferredMaxBytes string `yaml:"PreferredMaxBytes"`
	} `yaml:"BatchSize"`
	Kafka Kafkas `yaml:"Kafka"`
	Organizations interface{} `yaml:"Organizations"`
}


type ConfigTxYaml struct {  
	Organizations []Organization `yaml:"Organizations"`
	Orderer Orderers `yaml:"Orderer"`
	Application struct {
		Organizations interface{} `yaml:"Organizations"`
	} `yaml:"Application"`
	Profiles struct {
		TwoOrgsOrdererGenesis struct {
			Orderer Orderers `yaml:"Orderer"`
			Consortiums struct {
				SampleConsortium struct {
					Organizations []Organization `yaml:"Organizations"`
				} `yaml:"SampleConsortium"`
			} `yaml:"Consortiums"`
		} `yaml:"TwoOrgsOrdererGenesis"`
		TwoOrgsChannel struct {
			Consortium  string `yaml:"Consortium"`
			Application struct {
				Organizations []Organization `yaml:"Organizations"`
			} `yaml:"Application"`
		} `yaml:"TwoOrgsChannel"`
	} `yaml:"Profiles"`
} 

func doOutputBlock(config *genesisconfig.Profile, channelID string, outputBlock string) error {
	pgen := encoder.New(config)
	logger.Info("Generating genesis block")
	if config.Consortiums == nil {
		logger.Warning("Genesis block does not contain a consortiums group definition.  This block cannot be used for orderer bootstrap.")
	}
	genesisBlock := pgen.GenesisBlockForChannel(channelID)
	logger.Info("Writing genesis block")
	err := ioutil.WriteFile(outputBlock, utils.MarshalOrPanic(genesisBlock), 0644)
	if err != nil {
		return fmt.Errorf("Error writing genesis block: %s", err)
	}
	return nil
}

func doOutputChannelCreateTx(conf *genesisconfig.Profile, channelID string, outputChannelCreateTx string) error {
	logger.Info("Generating new channel configtx")

	//configtx, err := encoder.MakeChannelCreationTransaction(channelID, nil, nil, conf) //for v1.2
	configtx, err := encoder.MakeChannelCreationTransaction(channelID, nil, conf) //for v1.3
	if err != nil {
		return err
	}

	logger.Info("Writing new channel tx")
	err = ioutil.WriteFile(outputChannelCreateTx, utils.MarshalOrPanic(configtx), 0644)
	if err != nil {
		return fmt.Errorf("Error writing channel create tx: %s", err)
	}
	return nil
}

func doOutputAnchorPeersUpdate(conf *genesisconfig.Profile, channelID string, outputAnchorPeersUpdate string, asOrg string) error {
	logger.Info("Generating anchor peer update")
	if asOrg == "" {
		return fmt.Errorf("Must specify an organization to update the anchor peer for")
	}

	if conf.Application == nil {
		return fmt.Errorf("Cannot update anchor peers without an application section")
	}

	var org *genesisconfig.Organization
	for _, iorg := range conf.Application.Organizations {
		if iorg.Name == asOrg {
			org = iorg
		}
	}

	if org == nil {
		return fmt.Errorf("No organization name matching: %s", asOrg)
	}

	anchorPeers := make([]*pb.AnchorPeer, len(org.AnchorPeers))
	for i, anchorPeer := range org.AnchorPeers {
		anchorPeers[i] = &pb.AnchorPeer{
			Host: anchorPeer.Host,
			Port: int32(anchorPeer.Port),
		}
	}

	configUpdate := &cb.ConfigUpdate{
		ChannelId: channelID,
		WriteSet:  cb.NewConfigGroup(),
		ReadSet:   cb.NewConfigGroup(),
	}

	// Add all the existing config to the readset
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey] = cb.NewConfigGroup()
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Version = 1
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].ModPolicy = channelconfig.AdminsPolicyKey
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name] = cb.NewConfigGroup()
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Values[channelconfig.MSPKey] = &cb.ConfigValue{}
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.ReadersPolicyKey] = &cb.ConfigPolicy{}
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.WritersPolicyKey] = &cb.ConfigPolicy{}
	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.AdminsPolicyKey] = &cb.ConfigPolicy{}

	// Add all the existing at the same versions to the writeset
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey] = cb.NewConfigGroup()
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Version = 1
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].ModPolicy = channelconfig.AdminsPolicyKey
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name] = cb.NewConfigGroup()
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Version = 1
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].ModPolicy = channelconfig.AdminsPolicyKey
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Values[channelconfig.MSPKey] = &cb.ConfigValue{}
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.ReadersPolicyKey] = &cb.ConfigPolicy{}
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.WritersPolicyKey] = &cb.ConfigPolicy{}
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.AdminsPolicyKey] = &cb.ConfigPolicy{}
	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Values[channelconfig.AnchorPeersKey] = &cb.ConfigValue{
		Value:     utils.MarshalOrPanic(channelconfig.AnchorPeersValue(anchorPeers).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}

	configUpdateEnvelope := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: utils.MarshalOrPanic(configUpdate),
	}

	update := &cb.Envelope{
		Payload: utils.MarshalOrPanic(&cb.Payload{
			Header: &cb.Header{
				ChannelHeader: utils.MarshalOrPanic(&cb.ChannelHeader{
					ChannelId: channelID,
					Type:      int32(cb.HeaderType_CONFIG_UPDATE),
				}),
			},
			Data: utils.MarshalOrPanic(configUpdateEnvelope),
		}),
	}

	logger.Info("Writing anchor peer update")
	err := ioutil.WriteFile(outputAnchorPeersUpdate, utils.MarshalOrPanic(update), 0644)
	if err != nil {
		return fmt.Errorf("Error writing channel anchor peer update: %s", err)
	}
	return nil
}

func doInspectBlock(inspectBlock string) error {
	logger.Info("Inspecting block")
	data, err := ioutil.ReadFile(inspectBlock)
	if err != nil {
		return fmt.Errorf("Could not read block %s", inspectBlock)
	}

	logger.Info("Parsing genesis block")
	block, err := utils.UnmarshalBlock(data)
	if err != nil {
		return fmt.Errorf("error unmarshaling to block: %s", err)
	}
	err = protolator.DeepMarshalJSON(os.Stdout, block)
	if err != nil {
		return fmt.Errorf("malformed block contents: %s", err)
	}
	return nil
}

func doInspectChannelCreateTx(inspectChannelCreateTx string) error {
	logger.Info("Inspecting transaction")
	data, err := ioutil.ReadFile(inspectChannelCreateTx)
	if err != nil {
		return fmt.Errorf("could not read channel create tx: %s", err)
	}

	logger.Info("Parsing transaction")
	env, err := utils.UnmarshalEnvelope(data)
	if err != nil {
		return fmt.Errorf("Error unmarshaling envelope: %s", err)
	}

	err = protolator.DeepMarshalJSON(os.Stdout, env)
	if err != nil {
		return fmt.Errorf("malformed transaction contents: %s", err)
	}

	return nil
}

func CreateGenesisBlock(cfgTxPath string,channelName string,outputPath string)bool {
	factory.InitFactories(nil)
	var profileConfig *genesisconfig.Profile
	profileConfig = genesisconfig.Load("TwoOrgsOrdererGenesis",cfgTxPath)
	err := doOutputBlock(profileConfig,"testchainid", outputPath + "/genesis.block")
	if err!=nil {
		logger.Debug("CreateGenesisBlock failed")
		return false
	}
	return true
} 

func CreateChannelTx(cfgTxPath string,channelName string,outputPath string)bool{ 
	factory.InitFactories(nil)
	var profileConfig *genesisconfig.Profile
	profileConfig = genesisconfig.Load("TwoOrgsChannel",cfgTxPath)
	err := doOutputChannelCreateTx(profileConfig,channelName, outputPath + "/" + channelName + ".tx")
	if err!=nil {
		logger.Debug("CreateChannelTx failed")
		return false
	}
	return true
}

func CreateAnchorTx(cfgTxPath string,channelName string,outputPath string,orgs []string)bool{
	factory.InitFactories(nil)
	var profileConfig *genesisconfig.Profile
	profileConfig = genesisconfig.Load("TwoOrgsChannel",cfgTxPath)
	for i:=0;i<len(orgs);i++ {
		err := doOutputAnchorPeersUpdate(profileConfig, channelName, outputPath+"/"+orgs[i]+"anchors.tx",orgs[i])
		if err!=nil {
			logger.Debug("CreateAnchorTx failed")
			return false
		}
	} 
	return true
}

func doGenerateOrgJson(t *genesisconfig.TopLevel, orgName string,outFile string) (bool,string) {
	var msg string
	for _, org := range t.Organizations {
		if org.Name == orgName {
			og, err := encoder.NewOrdererOrgGroup(org)
			if err != nil {
				msg = "doGenerateOrgJson: bad org definition for org " + org.Name
				logger.Debug(msg)
				return false,msg
			}
			file, err := os.Create(outFile)
	        if err != nil {
			   msg = "doGenerateOrgJson: create out File Failed"
		       logger.Debug(msg)
		       return false,msg
			}
			defer file.Close()
			w := bufio.NewWriter(file)
			if err := protolator.DeepMarshalJSON(w, &cb.DynamicConsortiumOrgGroup{ConfigGroup: og}); err != nil {
				msg = "doGenerateOrgJson: malformed org definition for org: " + org.Name
				logger.Debug(msg)
				return false,msg
			}
			msg = "doGenerateOrgJson: success"
			return true,msg
		}
	}
	msg = "doGenerateOrgJson: organization not found,org= " + orgName
	return false,msg
}

func CreateAddOrgCfg(txCfgPath string,outFile string,orgName string)bool{
	factory.InitFactories(nil)
	var topLevelConfig *genesisconfig.TopLevel
	if txCfgPath != "" {
		topLevelConfig = genesisconfig.LoadTopLevel(txCfgPath)
	} else {
		topLevelConfig = genesisconfig.LoadTopLevel()
	}
	bl,msg := doGenerateOrgJson(topLevelConfig,orgName,outFile)
	if !bl {
		logger.Debug(msg) 
		return false
	}
    return true 
}
