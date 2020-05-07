
package sdkfabric
import (
	"fmt"
	"bytes"
	"unsafe"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"github.com/golang/protobuf/proto"
	"github.com/wingbaas/platformsrv/logger"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	ledgerUtil "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/util"
	fcb "github.com/hyperledger/fabric/protos/common"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"
	pbmsp "github.com/hyperledger/fabric/protos/msp"
	ab "github.com/hyperledger/fabric/protos/orderer"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
)

// Prettyprint print json
func Prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func deserializeIdentity(serializedID []byte) (*x509.Certificate, error) {
	sID := &pbmsp.SerializedIdentity{}
	err := proto.Unmarshal(serializedID,sID)
	if err != nil {
		logger.Errorf("Could not deserialize a SerializedIdentity, err %s",err)
		return nil,fmt.Errorf("Could not deserialize a SerializedIdentity")
	}
	bl, _ := pem.Decode(sID.IdBytes)
	if bl == nil {
		logger.Errorf("Could not decode the PEM structure")
		return nil,fmt.Errorf("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		logger.Errorf("ParseCertificate failed %s",err)
		return nil,fmt.Errorf("ParseCertificate failed")
	}
	return cert, nil
}

func copyChannelHeaderToLocalChannelHeader(localChannelHeader *ChannelHeader,chHeader *fcb.ChannelHeader,chaincodeHeaderExtension *pb.ChaincodeHeaderExtension) {
	localChannelHeader.Type = chHeader.Type
	localChannelHeader.Version = chHeader.Version
	localChannelHeader.Timestamp =  (*timestamp.Timestamp)(unsafe.Pointer(chHeader.Timestamp))
	localChannelHeader.ChannelID = chHeader.ChannelId
	localChannelHeader.TxID = chHeader.TxId
	localChannelHeader.Epoch = chHeader.Epoch
	localChannelHeader.ChaincodeID = chaincodeHeaderExtension.ChaincodeId
}

func copyChaincodeSpecToLocalChaincodeSpec(localChaincodeSpec *ChaincodeSpec,chaincodeSpec *pb.ChaincodeSpec) {
	localChaincodeSpec.Type = chaincodeSpec.Type
	localChaincodeSpec.ChaincodeID = chaincodeSpec.ChaincodeId
	localChaincodeSpec.Timeout = chaincodeSpec.Timeout
	chaincodeInput := &ChaincodeInput{}
	for _, input := range chaincodeSpec.Input.Args {
		chaincodeInput.Args = append(chaincodeInput.Args, string(input))
	}
	localChaincodeSpec.Input = chaincodeInput
}

func copyEndorsementToLocalEndorsement(localTransaction *Transaction,allEndorsements []*pb.Endorsement) {
	for _, endorser := range allEndorsements {
		endorsement := &Endorsement{}
		endorserSignatureHeader := &cb.SignatureHeader{}
		if err := proto.Unmarshal(endorser.Endorser, endorserSignatureHeader); err != nil {
			logger.Errorf("unmarshaling endorser signature: %s\n", err)
		}
		endorsement.SignatureHeader = getSignatureHeaderFromBlockData(endorserSignatureHeader)
		endorsement.Signature = endorser.Signature
		localTransaction.Endorsements = append(localTransaction.Endorsements,endorsement)
	}
}

func getValueFromBlockMetadata(block *cb.Block,index cb.BlockMetadataIndex) []byte {
	valueMetadata := &cb.Metadata{}
	if index == cb.BlockMetadataIndex_LAST_CONFIG {
		if err := proto.Unmarshal(block.Metadata.Metadata[index],valueMetadata); err != nil {
			logger.Errorf("getValueFromBlockMetadata: unmarshaling block.Metadata err= %s\n",err)
			return nil
		}
		lastConfig := &cb.LastConfig{}
		if err := proto.Unmarshal(valueMetadata.Value,lastConfig); err != nil {
			logger.Errorf("getValueFromBlockMetadata: unmarshaling valueMetadata.Value err= %s\n", err)
			return nil
		}
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(lastConfig.Index))
		return b
	} else if index == cb.BlockMetadataIndex_ORDERER {
		if err := proto.Unmarshal(block.Metadata.Metadata[index],valueMetadata); err != nil {
			logger.Errorf("getValueFromBlockMetadata: unmarshaling block.Metadata.Metadata err= %s\n",err)
			return nil
		}
		kafkaMetadata := &ab.KafkaMetadata{}
		if err := proto.Unmarshal(valueMetadata.Value,kafkaMetadata); err != nil {
			logger.Errorf("getValueFromBlockMetadata: unmarshaling kafkaMetadata err= %s\n",err)
			return nil
		}
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(kafkaMetadata.LastOffsetPersisted))
		return b
	} else if index == cb.BlockMetadataIndex_TRANSACTIONS_FILTER {
		return block.Metadata.Metadata[index]
	}
	return valueMetadata.Value
}

func getSignatureHeaderFromBlockMetadata(block *cb.Block,index cb.BlockMetadataIndex) (*SignatureMetadata, error) {
	signatureMetadata := &cb.Metadata{}
	if err := proto.Unmarshal(block.Metadata.Metadata[index],signatureMetadata); err != nil {
		logger.Errorf("getSignatureHeaderFromBlockMetadata: unmarshaling signatureMetadata err= %s\n",err)
		return nil, err
	}
	localSignatureHeader := &cb.SignatureHeader{}
	if len(signatureMetadata.Signatures) > 0 {
		if err := proto.Unmarshal(signatureMetadata.Signatures[0].SignatureHeader,localSignatureHeader); err != nil {
			logger.Errorf("getSignatureHeaderFromBlockMetadata: unmarshaling localSignatureHeader err= %s\n",err)
			return nil, err
		}
		localSignatureMetadata := &SignatureMetadata{}
		localSignatureMetadata.SignatureHeader = getSignatureHeaderFromBlockData(localSignatureHeader)
		localSignatureMetadata.Signature = signatureMetadata.Signatures[0].Signature
		return localSignatureMetadata, nil
	}
	return nil, nil
}

func getSignatureHeaderFromBlockData(header *cb.SignatureHeader) *SignatureHeader {
	signatureHeader := &SignatureHeader{}
	//signatureHeader.Certificate, _ = deserializeIdentity(header.Creator)
	signatureHeader.Nonce = header.Nonce
	return signatureHeader

}

// This method add transaction validation information from block TransactionFilter struct
func addTransactionValidation(block *Block,tran *Transaction, txIdx int) error {
	if len(block.TransactionFilter) > txIdx {
		tran.ValidationCode = uint8(block.TransactionFilter[txIdx])
		tran.ValidationCodeName = pb.TxValidationCode_name[int32(tran.ValidationCode)]
		return nil
	}
	logger.Errorf("Invalid index or transaction filler. Index: %d", txIdx)
	return fmt.Errorf("Invalid index or transaction filler")
}

// ProcessBlock parse block to json format
func processBlock(block *cb.Block) interface{} {
	var localBlock Block

	localBlock.Header = block.Header
	localBlock.TransactionFilter = ledgerUtil.NewTxValidationFlags(len(block.Data.Data))

	// process block metadata before data
	localBlock.BlockCreatorSignature, _ = getSignatureHeaderFromBlockMetadata(block,cb.BlockMetadataIndex_SIGNATURES)
	lastConfigBlockNumber := &LastConfigMetadata{}
	lastConfigBlockNumber.LastConfigBlockNum = binary.LittleEndian.Uint64(getValueFromBlockMetadata(block, cb.BlockMetadataIndex_LAST_CONFIG))
	lastConfigBlockNumber.SignatureData, _ = getSignatureHeaderFromBlockMetadata(block, cb.BlockMetadataIndex_LAST_CONFIG)
	localBlock.LastConfigBlockNumber = lastConfigBlockNumber

	txBytes := getValueFromBlockMetadata(block, cb.BlockMetadataIndex_TRANSACTIONS_FILTER)
	for index, b := range txBytes {
		localBlock.TransactionFilter[index] = b
	}

	ordererKafkaMetadata := &OrdererMetadata{}
	ordererKafkaMetadata.LastOffsetPersisted = binary.BigEndian.Uint64(getValueFromBlockMetadata(block, cb.BlockMetadataIndex_ORDERER))
	ordererKafkaMetadata.SignatureData, _ = getSignatureHeaderFromBlockMetadata(block, cb.BlockMetadataIndex_ORDERER)
	localBlock.OrdererKafkaMetadata = ordererKafkaMetadata

	for txIndex, data := range block.Data.Data {
		//Get envelope which is stored as byte array in the data field.
		envelope, err := utils.GetEnvelopeFromBlock(data)
		if err != nil {
			logger.Errorf("processBlock: getting envelope error=%s\n",err)
		}
		localTransaction := processTransaction(envelope)
		// add the transaction validation a
		addTransactionValidation(&localBlock, localTransaction, txIndex)
		//append the transaction
		localBlock.Transactions = append(localBlock.Transactions,localTransaction)
	}
	blockJSON, _ := json.Marshal(localBlock)
	blockJSONString, _ := Prettyprint(blockJSON)

	var rsObj interface{}
	err := json.Unmarshal(blockJSONString,&rsObj)
	if err != nil {
	  logger.Errorf("processBlock: unmarshal to json obj err=%s",err)
	  return nil
	}
	return rsObj
}

func processTransaction(envelope *fcb.Envelope) *Transaction {
	localTransaction := &Transaction{}
	//Get envelope which is stored as byte array in the data field.
	localTransaction.Signature = envelope.Signature
	//Get payload from envelope struct which is stored as byte array.
	payload, err := utils.GetPayload(envelope)
	if err != nil {
		logger.Errorf("processTransaction: getting payload from envelope: %s\n",err)
	}
	chHeader, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		logger.Errorf("processTransaction: unmarshaling channel header: %s\n",err)
	}
	headerExtension := &pb.ChaincodeHeaderExtension{}
	if err := proto.Unmarshal(chHeader.Extension,headerExtension); err != nil {
		logger.Errorf("processTransaction: unmarshaling chaincode header extension: %s\n",err)
	}
	localChannelHeader := &ChannelHeader{}
	copyChannelHeaderToLocalChannelHeader(localChannelHeader,chHeader,headerExtension)

	localTransaction.ChannelHeader = localChannelHeader
	localSignatureHeader := &cb.SignatureHeader{}
	if err := proto.Unmarshal(payload.Header.SignatureHeader,localSignatureHeader); err != nil {
		logger.Errorf("processTransaction: unmarshaling signature header: %s\n",err)
	}
	localTransaction.SignatureHeader = getSignatureHeaderFromBlockData(localSignatureHeader)
	//localTransaction.SignatureHeader.Nonce = localSignatureHeader.Nonce
	//localTransaction.SignatureHeader.Certificate, _ = deserializeIdentity(localSignatureHeader.Creator)
	transaction := &pb.Transaction{}
	if err := proto.Unmarshal(payload.Data,transaction); err != nil {
		logger.Errorf("processTransaction: unmarshaling transaction: %s\n",err)
	}
	chaincodeActionPayload, chaincodeAction, err := utils.GetPayloads(transaction.Actions[0])
	if err != nil {
		logger.Errorf("processTransaction: getting payloads from transaction actions: %s\n",err)
	}
	localSignatureHeader = &cb.SignatureHeader{}
	if err := proto.Unmarshal(transaction.Actions[0].Header, localSignatureHeader); err != nil {
		logger.Errorf("processTransaction: unmarshaling signature header: %s\n", err)
	}
	localTransaction.TxActionSignatureHeader = getSignatureHeaderFromBlockData(localSignatureHeader)
	//signatureHeader = &SignatureHeader{}
	//signatureHeader.Certificate, _ = deserializeIdentity(localSignatureHeader.Creator)
	//signatureHeader.Nonce = localSignatureHeader.Nonce
	//localTransaction.TxActionSignatureHeader = signatureHeader

	chaincodeProposalPayload := &pb.ChaincodeProposalPayload{}
	if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload,chaincodeProposalPayload); err != nil {
		logger.Errorf("processTransaction: unmarshaling chaincode proposal payload: %s\n",err)
	}
	chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{}
	if err := proto.Unmarshal(chaincodeProposalPayload.Input,chaincodeInvocationSpec); err != nil {
		logger.Errorf("processTransaction: unmarshaling chaincode invocationSpec: %s\n",err)
	}
	localChaincodeSpec := &ChaincodeSpec{}
	copyChaincodeSpecToLocalChaincodeSpec(localChaincodeSpec,chaincodeInvocationSpec.ChaincodeSpec)
	localTransaction.ChaincodeSpec = localChaincodeSpec
	copyEndorsementToLocalEndorsement(localTransaction,chaincodeActionPayload.Action.Endorsements)
	proposalResponsePayload := &pb.ProposalResponsePayload{}
	if err := proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload,proposalResponsePayload); err != nil {
		logger.Errorf("processTransaction: unmarshaling proposal response payload: %s\n",err)
	}
	localTransaction.ProposalHash = proposalResponsePayload.ProposalHash
	localTransaction.Response = chaincodeAction.Response
	events := &pb.ChaincodeEvent{}
	if err := proto.Unmarshal(chaincodeAction.Events,events); err != nil {
		logger.Errorf("processTransaction: unmarshaling chaincode action events:%s\n",err)
	}
	localTransaction.Events = events

	txReadWriteSet := &rwset.TxReadWriteSet{}
	if err := proto.Unmarshal(chaincodeAction.Results,txReadWriteSet); err != nil {
		logger.Errorf("processTransaction: unmarshaling chaincode action results: %s\n",err)
	}

	if len(chaincodeAction.Results) != 0 {
		for _, nsRwset := range txReadWriteSet.NsRwset {
			nsReadWriteSet := &NsReadWriteSet{}
			kvRWSet := &kvrwset.KVRWSet{}
			nsReadWriteSet.Namespace = nsRwset.Namespace
			if err := proto.Unmarshal(nsRwset.Rwset, kvRWSet); err != nil {
				logger.Errorf("processTransaction: unmarshaling tx read write set: %s\n",err)
			}
			nsReadWriteSet.KVRWSet = kvRWSet
			localTransaction.NsRwset = append(localTransaction.NsRwset, nsReadWriteSet)
		}
	}
	return localTransaction
}

func processTransaction2(envelope *fcb.Envelope) *Transaction {
	localTransaction := &Transaction{}
	//Get envelope which is stored as byte array in the data field.
	localTransaction.Signature = envelope.Signature
	//Get payload from envelope struct which is stored as byte array.
	payload, err := utils.GetPayload(envelope) 
	if err != nil {
		logger.Errorf("processTransaction2: getting payload from envelope: %s\n",err)
	}
	chHeader, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		logger.Errorf("processTransaction2: unmarshaling channel header: %s\n",err) 
	}
	headerExtension := &pb.ChaincodeHeaderExtension{}
	if err := proto.Unmarshal(chHeader.Extension,headerExtension); err != nil {
		logger.Errorf("processTransaction2: unmarshaling chaincode header extension: %s\n",err)
	}
	localChannelHeader := &ChannelHeader{}
	copyChannelHeaderToLocalChannelHeader(localChannelHeader,chHeader,headerExtension)

	localTransaction.ChannelHeader = localChannelHeader
	localSignatureHeader := &cb.SignatureHeader{}
	if err := proto.Unmarshal(payload.Header.SignatureHeader,localSignatureHeader); err != nil {
		logger.Errorf("processTransaction2: unmarshaling signature header: %s\n",err)
	}
	localTransaction.SignatureHeader = getSignatureHeaderFromBlockData(localSignatureHeader)
	//localTransaction.SignatureHeader.Nonce = localSignatureHeader.Nonce
	//localTransaction.SignatureHeader.Certificate, _ = deserializeIdentity(localSignatureHeader.Creator)
	transaction := &pb.Transaction{}
	if err := proto.Unmarshal(payload.Data, transaction); err != nil {
		logger.Errorf("processTransaction2: unmarshaling transaction: %s\n", err)
	}
	chaincodeActionPayload, chaincodeAction, err := utils.GetPayloads(transaction.Actions[0])
	if err != nil {
		logger.Errorf("processTransaction2: getting payloads from transaction actions: %s\n",err)
	}
	localSignatureHeader = &cb.SignatureHeader{}
	if err := proto.Unmarshal(transaction.Actions[0].Header, localSignatureHeader); err != nil {
		logger.Errorf("processTransaction2: unmarshaling signature header: %s\n",err)
	}
	localTransaction.TxActionSignatureHeader = getSignatureHeaderFromBlockData(localSignatureHeader)
	//signatureHeader = &SignatureHeader{}
	//signatureHeader.Certificate, _ = deserializeIdentity(localSignatureHeader.Creator)
	//signatureHeader.Nonce = localSignatureHeader.Nonce
	//localTransaction.TxActionSignatureHeader = signatureHeader

	chaincodeProposalPayload := &pb.ChaincodeProposalPayload{}
	if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload,chaincodeProposalPayload); err != nil {
		logger.Errorf("processTransaction2: unmarshaling chaincode proposal payload: %s\n", err)
	}
	chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{}
	if err := proto.Unmarshal(chaincodeProposalPayload.Input, chaincodeInvocationSpec); err != nil {
		logger.Errorf("processTransaction2: unmarshaling chaincode invocationSpec: %s\n",err)
	}
	localChaincodeSpec := &ChaincodeSpec{}
	copyChaincodeSpecToLocalChaincodeSpec(localChaincodeSpec,chaincodeInvocationSpec.ChaincodeSpec)
	localTransaction.ChaincodeSpec = localChaincodeSpec
	copyEndorsementToLocalEndorsement(localTransaction,chaincodeActionPayload.Action.Endorsements)
	proposalResponsePayload := &pb.ProposalResponsePayload{}
	if err := proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload,proposalResponsePayload); err != nil {
		logger.Errorf("processTransaction2: unmarshaling proposal response payload: %s\n",err)
	}
	localTransaction.ProposalHash = proposalResponsePayload.ProposalHash
	localTransaction.Response = chaincodeAction.Response
	events := &pb.ChaincodeEvent{}
	if err := proto.Unmarshal(chaincodeAction.Events, events); err != nil {
		logger.Errorf("processTransaction2: unmarshaling chaincode action events:%s\n",err)
	}
	localTransaction.Events = events

	txReadWriteSet := &rwset.TxReadWriteSet{}
	if err := proto.Unmarshal(chaincodeAction.Results, txReadWriteSet); err != nil {
		logger.Errorf("processTransaction2: unmarshaling chaincode action results: %s\n",err)
	}

	if len(chaincodeAction.Results) != 0 {
		for _, nsRwset := range txReadWriteSet.NsRwset {
			nsReadWriteSet := &NsReadWriteSet{}
			kvRWSet := &kvrwset.KVRWSet{}
			nsReadWriteSet.Namespace = nsRwset.Namespace
			if err := proto.Unmarshal(nsRwset.Rwset,kvRWSet); err != nil {
				logger.Errorf("processTransaction2: unmarshaling tx read write set: %s\n",err)
			}
			nsReadWriteSet.KVRWSet = kvRWSet
			localTransaction.NsRwset = append(localTransaction.NsRwset,nsReadWriteSet)
		}
	}
	return localTransaction
}
