
package txgeneratev2

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common" 
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/bccsp/factory"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/common/channelconfig"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/protolator"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/protolator/protoext/ordererext" 
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/encoder"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/genesisconfig" 
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/update"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/protoutil" 
	"github.com/pkg/errors"
)



func doOutputBlock(config *genesisconfig.Profile, channelID string, outputBlock string) error {
	pgen, err := encoder.NewBootstrapper(config)
	if err != nil {
		return errors.WithMessage(err, "could not create bootstrapper")
	}
	logger.Info("Generating genesis block")
	if config.Orderer == nil {
		return errors.Errorf("refusing to generate block which is missing orderer section")
	}
	if config.Consortiums == nil {
		logger.Warning("Genesis block does not contain a consortiums group definition.  This block cannot be used for orderer bootstrap.")
	}
	genesisBlock := pgen.GenesisBlockForChannel(channelID)
	logger.Info("Writing genesis block")
	err = writeFile(outputBlock, protoutil.MarshalOrPanic(genesisBlock), 0640)
	if err != nil {
		return fmt.Errorf("Error writing genesis block: %s", err)
	}
	return nil
}

func doOutputChannelCreateTx(conf, baseProfile *genesisconfig.Profile, channelID string, outputChannelCreateTx string) error {
	logger.Info("Generating new channel configtx")
	var configtx *cb.Envelope
	var err error
	if baseProfile == nil {
		configtx, err = encoder.MakeChannelCreationTransaction(channelID, nil, conf)
	} else {
		configtx, err = encoder.MakeChannelCreationTransactionWithSystemChannelContext(channelID, nil, conf, baseProfile)
	}
	if err != nil {
		return err
	}
	logger.Info("Writing new channel tx")
	err = writeFile(outputChannelCreateTx, protoutil.MarshalOrPanic(configtx), 0640)
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
	original, err := encoder.NewChannelGroup(conf)
	if err != nil {
		return errors.WithMessage(err, "error parsing profile as channel group")
	}
	original.Groups[channelconfig.ApplicationGroupKey].Version = 1
	updated := proto.Clone(original).(*cb.ConfigGroup)
	originalOrg, ok := original.Groups[channelconfig.ApplicationGroupKey].Groups[asOrg]
	if !ok {
		return errors.Errorf("org with name '%s' does not exist in config", asOrg)
	}
	if _, ok = originalOrg.Values[channelconfig.AnchorPeersKey]; !ok {
		return errors.Errorf("org '%s' does not have any anchor peers defined", asOrg)
	}
	delete(originalOrg.Values, channelconfig.AnchorPeersKey)
	updt, err := update.Compute(&cb.Config{ChannelGroup: original}, &cb.Config{ChannelGroup: updated})
	if err != nil {
		return errors.WithMessage(err, "could not compute update")
	}
	updt.ChannelId = channelID
	newConfigUpdateEnv := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: protoutil.MarshalOrPanic(updt),
	}
	updateTx, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, nil, newConfigUpdateEnv, 0, 0)
	logger.Info("Writing anchor peer update")
	err = writeFile(outputAnchorPeersUpdate, protoutil.MarshalOrPanic(updateTx), 0640)
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
	block, err := protoutil.UnmarshalBlock(data)
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
	env, err := protoutil.UnmarshalEnvelope(data)
	if err != nil {
		return fmt.Errorf("Error unmarshaling envelope: %s", err)
	}
	err = protolator.DeepMarshalJSON(os.Stdout, env)
	if err != nil {
		return fmt.Errorf("malformed transaction contents: %s", err)
	}
	return nil
}

func doPrintOrg(t *genesisconfig.TopLevel, printOrg string) error {
	for _, org := range t.Organizations {
		if org.Name == printOrg {
			og, err := encoder.NewOrdererOrgGroup(org)
			if err != nil {
				return errors.Wrapf(err, "bad org definition for org %s", org.Name)
			}
			if err := protolator.DeepMarshalJSON(os.Stdout, &ordererext.DynamicOrdererOrgGroup{ConfigGroup: og}); err != nil {
				return errors.Wrapf(err, "malformed org definition for org: %s", org.Name)
			}
			return nil
		}
	}
	return errors.Errorf("organization %s not found", printOrg)
}

func writeFile(filename string, data []byte, perm os.FileMode) error {
	dirPath := filepath.Dir(filename)
	exists, err := dirExists(dirPath)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(dirPath, 0750)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, data, perm)
}
func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateGenesisBlock(cfgTxPath string,channelName string,outputPath string)bool {
	factory.InitFactories(nil)
	var profileConfig *genesisconfig.Profile
	profileConfig = genesisconfig.Load("TwoOrgsOrdererGenesis",cfgTxPath) 
	err := doOutputBlock(profileConfig,"testchainid", outputPath + "/genesis.block")
	if err != nil {
		logger.Errorf("CreateGenesisBlock failed,err=%s",err.Error())
		return false
	}
	return true 
} 

func CreateChannelTx(cfgTxPath string,channelName string,outputPath string)bool{ 
	factory.InitFactories(nil)
	var profileConfig *genesisconfig.Profile
	//var baseProfile *genesisconfig.Profile
	profileConfig = genesisconfig.Load("TwoOrgsChannel",cfgTxPath)
	err := doOutputChannelCreateTx(profileConfig,nil,channelName, outputPath + "/" + channelName + ".tx")
	if err != nil {
		logger.Errorf("CreateChannelTx failed")
		return false
	}
	return true
}

func CreateAnchorTx(cfgTxPath string,channelName string,outputPath string,orgs []string)bool{
	factory.InitFactories(nil)
	var profileConfig *genesisconfig.Profile
	profileConfig = genesisconfig.Load("TwoOrgsChannel",cfgTxPath)
	for i:=0;i<len(orgs);i++ {
		err := doOutputAnchorPeersUpdate(profileConfig,channelName, outputPath+"/"+orgs[i]+"anchors.tx",orgs[i])
		if err != nil {
			logger.Errorf("CreateAnchorTx failed")
			return false
		}
	} 
	return true
}

