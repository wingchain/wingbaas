
package rest

import (
	"io/ioutil"
	"os"
	"reflect"
	"github.com/hyperledger/fabric/common/tools/protolator"
	_ "github.com/hyperledger/fabric/protos/common"
	_ "github.com/hyperledger/fabric/protos/msp"
	_ "github.com/hyperledger/fabric/protos/orderer"
	_ "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/common/tools/configtxlator/update"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/golang/protobuf/proto"
	"github.com/wingbaas/platformsrv/settings/fabric/txgenerate/pkg/errors"
	"github.com/wingbaas/platformsrv/logger"
)

func EncodeProto(msgName string, input, output *os.File) error {
	msgType := proto.MessageType(msgName)
	if msgType == nil {
		return errors.Errorf("message of type %s unknown", msgType)
	}
	msg := reflect.New(msgType.Elem()).Interface().(proto.Message)

	err := protolator.DeepUnmarshalJSON(input, msg)
	if err != nil {
		return errors.Wrapf(err, "error encoding input")
	}
	out, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrapf(err, "error marshaling")
	}
	_, err = output.Write(out)
	if err != nil {
		return errors.Wrapf(err, "error writing output")
	}
	return nil
}

func DecodeProto(msgName string, input, output *os.File) error {
	msgType := proto.MessageType(msgName)
	logger.Debug("decodeProto msgname="+msgName+" type=")
	logger.Debug(msgType)  
	if msgType == nil {
		return errors.Errorf("message of type %s unknown", msgType)
	}
	msg := reflect.New(msgType.Elem()).Interface().(proto.Message)
	in, err := ioutil.ReadAll(input)
	if err != nil {
		return errors.Wrapf(err, "error reading input")
	}
	err = proto.Unmarshal(in, msg)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling")
	}
	err = protolator.DeepMarshalJSON(output, msg)
	if err != nil {
		return errors.Wrapf(err, "error encoding output")
	}
	return nil
}

func ComputeUpdt(original, updated, output *os.File, channelID string) error {
	origIn, err := ioutil.ReadAll(original)
	if err != nil {
		return errors.Wrapf(err, "error reading original config")
	}
	origConf := &cb.Config{}
	err = proto.Unmarshal(origIn, origConf)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling original config")
	}
	updtIn, err := ioutil.ReadAll(updated)
	if err != nil {
		return errors.Wrapf(err, "error reading updated config")
	}
	updtConf := &cb.Config{}
	err = proto.Unmarshal(updtIn, updtConf)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling updated config")
	}
	cu, err := update.Compute(origConf, updtConf)
	if err != nil {
		return errors.Wrapf(err, "error computing config update")
	}
	cu.ChannelId = channelID
	outBytes, err := proto.Marshal(cu)
	if err != nil {
		return errors.Wrapf(err, "error marshaling computed config update")
	}
	_, err = output.Write(outBytes)
	if err != nil {
		return errors.Wrapf(err, "error writing config update to output")
	}
	return nil
}