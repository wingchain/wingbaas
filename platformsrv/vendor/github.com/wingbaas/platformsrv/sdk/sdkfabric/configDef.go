
package sdkfabric

import (
	"time"
)

type LoggingSt struct {
	Level string `json:"level"`
}

type CryptoconfigSt struct {
	Path string `json:"path"`
}

type CryptoStoreSt struct {
	Path string `json:"path"`
}

type CredentialStoreSt struct {
	Path string `json:"path"`
	CryptoStore CryptoStoreSt `json:"cryptoStore"`
}

type DefaultSt struct {
	Provider string `json:"provider"`
}

type SecuritySt struct {
	Enabled bool `json:"enabled"`
	Default DefaultSt `json:"default"`
	HashAlgorithm string `json:"hashAlgorithm"`
	SoftVerify bool `json:"softVerify"`
	Level int `json:"level"`
}

type BCCSPSt struct {
	Security SecuritySt `json:"security"` 
}

type ClientTLSCertsSt struct {
	Key CryptoconfigSt `json:"key"`
	Cert CryptoconfigSt `json:"cert"`
}

type ClientCertSt struct {
	Keyfile string `json:"keyfile"`
	Certfile string `json:"certfile"`
}

type TLSCertsSt struct {
	SystemCertPool bool `json:"systemCertPool"`
	Client ClientTLSCertsSt `json:"client,omitempty"`
	//Client ClientCertSt`json:"client,omitempty"`
}

type ClientSt struct {
	Organization string `json:"organization"`
	Logging LoggingSt `json:"logging"`
	Cryptoconfig CryptoconfigSt `json:"cryptoconfig"`
	CredentialStore CredentialStoreSt `json:"credentialStore"`
	BCCSP BCCSPSt `json:"BCCSP"`
	TLSCerts TLSCertsSt `json:"tlsCerts"`
} 

type PeerField struct {
	EndorsingPeer bool `json:"endorsingPeer"`
	ChaincodeQuery bool `json:"chaincodeQuery"`
	LedgerQuery bool `json:"ledgerQuery"`
	EventSource bool `json:"eventSource"`
}

type RetryOptsSt struct {
	Attempts int `json:"attempts"`
	InitialBackoff string `json:"initialBackoff"`
	MaxBackoff string `json:"maxBackoff"`
	BackoffFactor float64 `json:"backoffFactor"`
}

type QueryChannelConfigSt struct {
	MinResponses int `json:"minResponses"`
	MaxTargets int `json:"maxTargets"`
	RetryOpts RetryOptsSt `json:"retryOpts"`
}

type PoliciesSt struct {
	QueryChannelConfig QueryChannelConfigSt `json:"queryChannelConfig"`
}

type ChannelField struct {
	Orderers []string `json:"orderers"`
	Peers map[string]PeerField `json:"peers"`
	Policies PoliciesSt `json:"policies"`
}

type OrgField struct {
	Mspid string `json:"mspid"`
	CryptoPath string `json:"cryptoPath"`
	Peers []string `json:"peers,omitempty"`
	CertificateAuthorities []string `json:"certificateAuthorities,omitempty"`
	//Users map[string]ClientTLSCertsSt `json:"users"` 
}

type GrpcOptionsSt struct {
	SslTargetNameOverride string `json:"ssl-target-name-override"`
	KeepAliveTime string `json:"keep-alive-time"`
	KeepAliveTimeout string `json:"keep-alive-timeout"`
	KeepAlivePermit bool `json:"keep-alive-permit"`
	FailFast bool `json:"fail-fast"`
	AllowInsecure bool `json:"allow-insecure"`
}

type TLSCACertsSt struct {
	Path string `json:"path"`
	Client ClientTLSCertsSt `json:"client,omitempty"`
}

type MemberField struct {
	URL string `json:"url"`
	EventURL string `json:"eventUrl,omitempty"`
	GrpcOptions GrpcOptionsSt `json:"grpcOptions"`
	TLSCACerts TLSCACertsSt `json:"tlsCACerts"`
}

type HTTPOptionsSt struct {
	Verify bool `json:"verify"`
} 

type RegistrarSt struct {
	EnrollID string `json:"enrollId"`
	EnrollSecret string `json:"enrollSecret"`
}

type CaField struct {
	URL string `json:"url"`
	HTTPOptions HTTPOptionsSt `json:"httpOptions"`
	Registrar RegistrarSt `json:"registrar"`
	CaName string `json:"caName"`
	TLSCACerts TLSCACertsSt `json:"tlsCACerts"`
}

type MatchFieldSt struct {
	Pattern string `json:"pattern,omitempty"`
	URLSubstitutionExp string `json:"urlSubstitutionExp,omitempty"`
	EventURLSubstitutionExp string `json:"eventUrlSubstitutionExp,omitempty"`
	SslTargetOverrideURLSubstitutionExp string `json:"sslTargetOverrideUrlSubstitutionExp,omitempty"`
	MappedHost string `json:"mappedHost,omitempty"`
}

type EntityMatchersSt struct {
	Peer []MatchFieldSt `json:"peer"`
	Orderer []MatchFieldSt`json:"orderer"`
	CertificateAuthorities []MatchFieldSt `json:"certificateAuthorities"` 
}

type SdkFabricCfg struct {
	Name string `json:"name"`
	Version string `json:"version"`
	Client ClientSt `json:"client"`
	Channels map[string]ChannelField `json:"channels"`
	Organizations map[string]OrgField `json:"organizations"`
	Orderers map[string]MemberField `json:"orderers"`
	Peers map[string]MemberField `json:"peers"`
	CertificateAuthorities map[string]CaField `json:"certificateAuthorities"`
	EntityMatchers EntityMatchersSt `json:"entityMatchers"` 
}

type CCInvokeResult struct {
	Proposal struct {
		TxnID string `json:"TxnID"`
		Header string `json:"header"`
		Payload string `json:"payload"`
	} `json:"Proposal"`
	Responses []struct {
		Endorser string `json:"Endorser"`
		Status int `json:"Status"`
		ChaincodeStatus int `json:"ChaincodeStatus"`
		Version int `json:"version"`
		Response struct {
			Status int `json:"status"`
		} `json:"response"`
		Payload string `json:"payload"`
		Endorsement struct {
			Endorser string `json:"endorser"`
			Signature string `json:"signature"` 
		} `json:"endorsement"`
	} `json:"Responses"`
	TransactionID string `json:"TransactionID"`
	TxValidationCode int `json:"TxValidationCode"`
	ChaincodeStatus int `json:"ChaincodeStatus"`
	Payload interface{} `json:"Payload"`
}

type BlockInfo struct {
	Header struct {
		Number int `json:"number"`
		PreviousHash string `json:"previous_hash"`
		DataHash string `json:"data_hash"`
	} `json:"header"`
	Transactions []struct {
		Signature string `json:"signature"`
		ChannelHeader struct {
			Type int `json:"type"`
			ChannelID string `json:"channel_id"`
			TxID string `json:"tx_id"`
			ChaincodeID struct {
				Name string `json:"name"`
			} `json:"chaincode_id"`
		} `json:"channel_header"`
		SignatureHeader struct {
			Certificate struct {
				Raw string `json:"Raw"`
				RawTBSCertificate string `json:"RawTBSCertificate"`
				RawSubjectPublicKeyInfo string `json:"RawSubjectPublicKeyInfo"`
				RawSubject string `json:"RawSubject"`
				RawIssuer string `json:"RawIssuer"`
				Signature string `json:"Signature"`
				SignatureAlgorithm int `json:"SignatureAlgorithm"`
				PublicKeyAlgorithm int `json:"PublicKeyAlgorithm"`
				PublicKey struct {
					Curve struct {
						P int64 `json:"P"`
						N int64 `json:"N"`
						B int64 `json:"B"`
						Gx int64 `json:"Gx"`
						Gy int64 `json:"Gy"`
						BitSize int `json:"BitSize"`
						Name string `json:"Name"`
					} `json:"Curve"`
					X int64 `json:"X"`
					Y int64 `json:"Y"`
				} `json:"PublicKey"`
				Version int `json:"Version"`
				SerialNumber int64 `json:"SerialNumber"`
				Issuer struct {
					Country []string `json:"Country"`
					Organization []string `json:"Organization"`
					OrganizationalUnit interface{} `json:"OrganizationalUnit"`
					Locality []string `json:"Locality"`
					Province []string `json:"Province"`
					StreetAddress interface{} `json:"StreetAddress"`
					PostalCode interface{} `json:"PostalCode"`
					SerialNumber string `json:"SerialNumber"`
					CommonName string `json:"CommonName"`
					Names []struct {
						Type []int `json:"Type"`
						Value string `json:"Value"`
					} `json:"Names"`
					ExtraNames interface{} `json:"ExtraNames"`
				} `json:"Issuer"`
				Subject struct {
					Country interface{} `json:"Country"`
					Organization interface{} `json:"Organization"`
					OrganizationalUnit []string `json:"OrganizationalUnit"`
					Locality interface{} `json:"Locality"`
					Province interface{} `json:"Province"`
					StreetAddress interface{} `json:"StreetAddress"`
					PostalCode interface{} `json:"PostalCode"`
					SerialNumber string `json:"SerialNumber"`
					CommonName string `json:"CommonName"`
					Names []struct {
						Type []int `json:"Type"`
						Value string `json:"Value"`
					} `json:"Names"`
					ExtraNames interface{} `json:"ExtraNames"`
				} `json:"Subject"`
				NotBefore time.Time `json:"NotBefore"`
				NotAfter time.Time `json:"NotAfter"`
				KeyUsage int `json:"KeyUsage"`
				Extensions []struct {
					ID []int `json:"Id"`
					Critical bool `json:"Critical"`
					Value string `json:"Value"`
				} `json:"Extensions"`
				ExtraExtensions interface{} `json:"ExtraExtensions"`
				UnhandledCriticalExtensions interface{} `json:"UnhandledCriticalExtensions"`
				ExtKeyUsage interface{} `json:"ExtKeyUsage"`
				UnknownExtKeyUsage interface{} `json:"UnknownExtKeyUsage"`
				BasicConstraintsValid bool `json:"BasicConstraintsValid"`
				IsCA bool `json:"IsCA"`
				MaxPathLen int `json:"MaxPathLen"`
				MaxPathLenZero bool `json:"MaxPathLenZero"`
				SubjectKeyID string `json:"SubjectKeyId"`
				AuthorityKeyID string `json:"AuthorityKeyId"`
				OCSPServer interface{} `json:"OCSPServer"`
				IssuingCertificateURL interface{} `json:"IssuingCertificateURL"`
				DNSNames []string `json:"DNSNames"`
				EmailAddresses interface{} `json:"EmailAddresses"`
				IPAddresses interface{} `json:"IPAddresses"`
				URIs interface{} `json:"URIs"`
				PermittedDNSDomainsCritical bool `json:"PermittedDNSDomainsCritical"`
				PermittedDNSDomains interface{} `json:"PermittedDNSDomains"`
				ExcludedDNSDomains interface{} `json:"ExcludedDNSDomains"`
				PermittedIPRanges interface{} `json:"PermittedIPRanges"`
				ExcludedIPRanges interface{} `json:"ExcludedIPRanges"`
				PermittedEmailAddresses interface{} `json:"PermittedEmailAddresses"`
				ExcludedEmailAddresses interface{} `json:"ExcludedEmailAddresses"`
				PermittedURIDomains interface{} `json:"PermittedURIDomains"`
				ExcludedURIDomains interface{} `json:"ExcludedURIDomains"`
				CRLDistributionPoints interface{} `json:"CRLDistributionPoints"`
				PolicyIdentifiers interface{} `json:"PolicyIdentifiers"`
			} `json:"Certificate"`
			Nonce string `json:"nonce"`
		} `json:"signature_header"`
		TxActionSignatureHeader struct {
			Certificate struct {
				Raw string `json:"Raw"`
				RawTBSCertificate string `json:"RawTBSCertificate"`
				RawSubjectPublicKeyInfo string `json:"RawSubjectPublicKeyInfo"`
				RawSubject string `json:"RawSubject"`
				RawIssuer string `json:"RawIssuer"`
				Signature string `json:"Signature"`
				SignatureAlgorithm int `json:"SignatureAlgorithm"`
				PublicKeyAlgorithm int `json:"PublicKeyAlgorithm"`
				PublicKey struct {
					Curve struct {
						P int64 `json:"P"`
						N int64 `json:"N"`
						B int64 `json:"B"`
						Gx int64 `json:"Gx"`
						Gy int64 `json:"Gy"`
						BitSize int `json:"BitSize"`
						Name string `json:"Name"`
					} `json:"Curve"`
					X int64 `json:"X"`
					Y int64 `json:"Y"`
				} `json:"PublicKey"`
				Version int `json:"Version"`
				SerialNumber int64 `json:"SerialNumber"`
				Issuer struct {
					Country []string `json:"Country"`
					Organization []string `json:"Organization"`
					OrganizationalUnit interface{} `json:"OrganizationalUnit"`
					Locality []string `json:"Locality"`
					Province []string `json:"Province"`
					StreetAddress interface{} `json:"StreetAddress"`
					PostalCode interface{} `json:"PostalCode"`
					SerialNumber string `json:"SerialNumber"`
					CommonName string `json:"CommonName"`
					Names []struct {
						Type []int `json:"Type"`
						Value string `json:"Value"`
					} `json:"Names"`
					ExtraNames interface{} `json:"ExtraNames"`
				} `json:"Issuer"`
				Subject struct {
					Country interface{} `json:"Country"`
					Organization interface{} `json:"Organization"`
					OrganizationalUnit []string `json:"OrganizationalUnit"`
					Locality interface{} `json:"Locality"`
					Province interface{} `json:"Province"`
					StreetAddress interface{} `json:"StreetAddress"`
					PostalCode interface{} `json:"PostalCode"`
					SerialNumber string `json:"SerialNumber"`
					CommonName string `json:"CommonName"`
					Names []struct {
						Type []int `json:"Type"`
						Value string `json:"Value"`
					} `json:"Names"`
					ExtraNames interface{} `json:"ExtraNames"`
				} `json:"Subject"`
				NotBefore time.Time `json:"NotBefore"`
				NotAfter time.Time `json:"NotAfter"`
				KeyUsage int `json:"KeyUsage"`
				Extensions []struct {
					ID []int `json:"Id"`
					Critical bool `json:"Critical"`
					Value string `json:"Value"`
				} `json:"Extensions"`
				ExtraExtensions interface{} `json:"ExtraExtensions"`
				UnhandledCriticalExtensions interface{} `json:"UnhandledCriticalExtensions"`
				ExtKeyUsage interface{} `json:"ExtKeyUsage"`
				UnknownExtKeyUsage interface{} `json:"UnknownExtKeyUsage"`
				BasicConstraintsValid bool `json:"BasicConstraintsValid"`
				IsCA bool `json:"IsCA"`
				MaxPathLen int `json:"MaxPathLen"`
				MaxPathLenZero bool `json:"MaxPathLenZero"`
				SubjectKeyID string `json:"SubjectKeyId"`
				AuthorityKeyID string `json:"AuthorityKeyId"`
				OCSPServer interface{} `json:"OCSPServer"`
				IssuingCertificateURL interface{} `json:"IssuingCertificateURL"`
				DNSNames []string `json:"DNSNames"`
				EmailAddresses interface{} `json:"EmailAddresses"`
				IPAddresses interface{} `json:"IPAddresses"`
				URIs interface{} `json:"URIs"`
				PermittedDNSDomainsCritical bool `json:"PermittedDNSDomainsCritical"`
				PermittedDNSDomains interface{} `json:"PermittedDNSDomains"`
				ExcludedDNSDomains interface{} `json:"ExcludedDNSDomains"`
				PermittedIPRanges interface{} `json:"PermittedIPRanges"`
				ExcludedIPRanges interface{} `json:"ExcludedIPRanges"`
				PermittedEmailAddresses interface{} `json:"PermittedEmailAddresses"`
				ExcludedEmailAddresses interface{} `json:"ExcludedEmailAddresses"`
				PermittedURIDomains interface{} `json:"PermittedURIDomains"`
				ExcludedURIDomains interface{} `json:"ExcludedURIDomains"`
				CRLDistributionPoints interface{} `json:"CRLDistributionPoints"`
				PolicyIdentifiers interface{} `json:"PolicyIdentifiers"`
			} `json:"Certificate"`
			Nonce string `json:"nonce"`
		} `json:"tx_action_signature_header"`
		ChaincodeSpec struct {
			Type int `json:"type"`
			ChaincodeID struct {
				Name string `json:"name"`
			} `json:"chaincode_id"`
			Input struct {
				Args []string `json:"Args"`
			} `json:"input"`
		} `json:"chaincode_spec"`
		Endorsements []struct {
			SignatureHeader struct {
				Certificate interface{} `json:"Certificate"`
				Nonce string `json:"nonce"`
			} `json:"signature_header"`
			Signature string `json:"signature"`
		} `json:"endorsements"`
		ProposalHash string `json:"proposal_hash"`
		Events struct {
		} `json:"events"`
		Response struct {
			Status int `json:"status"`
		} `json:"response"`
		NsReadWriteSet []struct {
			Namespace string `json:"Namespace"`
			KVRWSet struct {
				Reads []struct {
					Key string `json:"key"`
					Version struct {
						BlockNum int `json:"block_num"`
					} `json:"version"`
				} `json:"reads"`
				Writes []struct {
					Key string `json:"key"`
					Value string `json:"value"`
				} `json:"writes"`
			} `json:"KVRWSet"`
		} `json:"ns_read_write_Set"`
		ValidationCode int `json:"validation_code"`
		ValidationCodeName string `json:"validation_code_name"`
	} `json:"transactions"`
	BlockCreatorSignature struct {
		SignatureHeader struct {
			Certificate struct {
				Raw string `json:"Raw"`
				RawTBSCertificate string `json:"RawTBSCertificate"`
				RawSubjectPublicKeyInfo string `json:"RawSubjectPublicKeyInfo"`
				RawSubject string `json:"RawSubject"`
				RawIssuer string `json:"RawIssuer"`
				Signature string `json:"Signature"`
				SignatureAlgorithm int `json:"SignatureAlgorithm"`
				PublicKeyAlgorithm int `json:"PublicKeyAlgorithm"`
				PublicKey struct {
					Curve struct {
						P int64 `json:"P"`
						N int64 `json:"N"`
						B int64 `json:"B"`
						Gx int64 `json:"Gx"`
						Gy int64 `json:"Gy"`
						BitSize int `json:"BitSize"`
						Name string `json:"Name"`
					} `json:"Curve"`
					X int64 `json:"X"`
					Y int64 `json:"Y"`
				} `json:"PublicKey"`
				Version int `json:"Version"`
				SerialNumber int64 `json:"SerialNumber"`
				Issuer struct {
					Country []string `json:"Country"`
					Organization []string `json:"Organization"`
					OrganizationalUnit interface{} `json:"OrganizationalUnit"`
					Locality []string `json:"Locality"`
					Province []string `json:"Province"`
					StreetAddress interface{} `json:"StreetAddress"`
					PostalCode interface{} `json:"PostalCode"`
					SerialNumber string `json:"SerialNumber"`
					CommonName string `json:"CommonName"`
					Names []struct {
						Type []int `json:"Type"`
						Value string `json:"Value"`
					} `json:"Names"`
					ExtraNames interface{} `json:"ExtraNames"`
				} `json:"Issuer"`
				Subject struct {
					Country []string `json:"Country"`
					Organization interface{} `json:"Organization"`
					OrganizationalUnit interface{} `json:"OrganizationalUnit"`
					Locality []string `json:"Locality"`
					Province []string `json:"Province"`
					StreetAddress interface{} `json:"StreetAddress"`
					PostalCode interface{} `json:"PostalCode"`
					SerialNumber string `json:"SerialNumber"`
					CommonName string `json:"CommonName"`
					Names []struct {
						Type []int `json:"Type"`
						Value string `json:"Value"`
					} `json:"Names"`
					ExtraNames interface{} `json:"ExtraNames"`
				} `json:"Subject"`
				NotBefore time.Time `json:"NotBefore"`
				NotAfter time.Time `json:"NotAfter"`
				KeyUsage int `json:"KeyUsage"`
				Extensions []struct {
					ID []int `json:"Id"`
					Critical bool `json:"Critical"`
					Value string `json:"Value"`
				} `json:"Extensions"`
				ExtraExtensions interface{} `json:"ExtraExtensions"`
				UnhandledCriticalExtensions interface{} `json:"UnhandledCriticalExtensions"`
				ExtKeyUsage interface{} `json:"ExtKeyUsage"`
				UnknownExtKeyUsage interface{} `json:"UnknownExtKeyUsage"`
				BasicConstraintsValid bool `json:"BasicConstraintsValid"`
				IsCA bool `json:"IsCA"`
				MaxPathLen int `json:"MaxPathLen"`
				MaxPathLenZero bool `json:"MaxPathLenZero"`
				SubjectKeyID interface{} `json:"SubjectKeyId"`
				AuthorityKeyID string `json:"AuthorityKeyId"`
				OCSPServer interface{} `json:"OCSPServer"`
				IssuingCertificateURL interface{} `json:"IssuingCertificateURL"`
				DNSNames interface{} `json:"DNSNames"`
				EmailAddresses interface{} `json:"EmailAddresses"`
				IPAddresses interface{} `json:"IPAddresses"`
				URIs interface{} `json:"URIs"`
				PermittedDNSDomainsCritical bool `json:"PermittedDNSDomainsCritical"`
				PermittedDNSDomains interface{} `json:"PermittedDNSDomains"`
				ExcludedDNSDomains interface{} `json:"ExcludedDNSDomains"`
				PermittedIPRanges interface{} `json:"PermittedIPRanges"`
				ExcludedIPRanges interface{} `json:"ExcludedIPRanges"`
				PermittedEmailAddresses interface{} `json:"PermittedEmailAddresses"`
				ExcludedEmailAddresses interface{} `json:"ExcludedEmailAddresses"`
				PermittedURIDomains interface{} `json:"PermittedURIDomains"`
				ExcludedURIDomains interface{} `json:"ExcludedURIDomains"`
				CRLDistributionPoints interface{} `json:"CRLDistributionPoints"`
				PolicyIdentifiers interface{} `json:"PolicyIdentifiers"`
			} `json:"Certificate"`
			Nonce string `json:"nonce"`
		} `json:"signature_header"`
		Signature string `json:"signature"`
	} `json:"block_creator_signature"`
	LastConfigBlockNumber struct {
		SignatureData struct {
			SignatureHeader struct {
				Certificate struct {
					Raw string `json:"Raw"`
					RawTBSCertificate string `json:"RawTBSCertificate"`
					RawSubjectPublicKeyInfo string `json:"RawSubjectPublicKeyInfo"`
					RawSubject string `json:"RawSubject"`
					RawIssuer string `json:"RawIssuer"`
					Signature string `json:"Signature"`
					SignatureAlgorithm int `json:"SignatureAlgorithm"`
					PublicKeyAlgorithm int `json:"PublicKeyAlgorithm"`
					PublicKey struct {
						Curve struct {
							P int64 `json:"P"`
							N int64 `json:"N"`
							B int64 `json:"B"`
							Gx int64 `json:"Gx"`
							Gy int64 `json:"Gy"`
							BitSize int `json:"BitSize"`
							Name string `json:"Name"`
						} `json:"Curve"`
						X int64 `json:"X"`
						Y int64 `json:"Y"`
					} `json:"PublicKey"`
					Version int `json:"Version"`
					SerialNumber int64 `json:"SerialNumber"`
					Issuer struct {
						Country []string `json:"Country"`
						Organization []string `json:"Organization"`
						OrganizationalUnit interface{} `json:"OrganizationalUnit"`
						Locality []string `json:"Locality"`
						Province []string `json:"Province"`
						StreetAddress interface{} `json:"StreetAddress"`
						PostalCode interface{} `json:"PostalCode"`
						SerialNumber string `json:"SerialNumber"`
						CommonName string `json:"CommonName"`
						Names []struct {
							Type []int `json:"Type"`
							Value string `json:"Value"`
						} `json:"Names"`
						ExtraNames interface{} `json:"ExtraNames"`
					} `json:"Issuer"`
					Subject struct {
						Country []string `json:"Country"`
						Organization interface{} `json:"Organization"`
						OrganizationalUnit interface{} `json:"OrganizationalUnit"`
						Locality []string `json:"Locality"`
						Province []string `json:"Province"`
						StreetAddress interface{} `json:"StreetAddress"`
						PostalCode interface{} `json:"PostalCode"`
						SerialNumber string `json:"SerialNumber"`
						CommonName string `json:"CommonName"`
						Names []struct {
							Type []int `json:"Type"`
							Value string `json:"Value"`
						} `json:"Names"`
						ExtraNames interface{} `json:"ExtraNames"`
					} `json:"Subject"`
					NotBefore time.Time `json:"NotBefore"`
					NotAfter time.Time `json:"NotAfter"`
					KeyUsage int `json:"KeyUsage"`
					Extensions []struct {
						ID []int `json:"Id"`
						Critical bool `json:"Critical"`
						Value string `json:"Value"`
					} `json:"Extensions"`
					ExtraExtensions interface{} `json:"ExtraExtensions"`
					UnhandledCriticalExtensions interface{} `json:"UnhandledCriticalExtensions"`
					ExtKeyUsage interface{} `json:"ExtKeyUsage"`
					UnknownExtKeyUsage interface{} `json:"UnknownExtKeyUsage"`
					BasicConstraintsValid bool `json:"BasicConstraintsValid"`
					IsCA bool `json:"IsCA"`
					MaxPathLen int `json:"MaxPathLen"`
					MaxPathLenZero bool `json:"MaxPathLenZero"`
					SubjectKeyID interface{} `json:"SubjectKeyId"`
					AuthorityKeyID string `json:"AuthorityKeyId"`
					OCSPServer interface{} `json:"OCSPServer"`
					IssuingCertificateURL interface{} `json:"IssuingCertificateURL"`
					DNSNames interface{} `json:"DNSNames"`
					EmailAddresses interface{} `json:"EmailAddresses"`
					IPAddresses interface{} `json:"IPAddresses"`
					URIs interface{} `json:"URIs"`
					PermittedDNSDomainsCritical bool `json:"PermittedDNSDomainsCritical"`
					PermittedDNSDomains interface{} `json:"PermittedDNSDomains"`
					ExcludedDNSDomains interface{} `json:"ExcludedDNSDomains"`
					PermittedIPRanges interface{} `json:"PermittedIPRanges"`
					ExcludedIPRanges interface{} `json:"ExcludedIPRanges"`
					PermittedEmailAddresses interface{} `json:"PermittedEmailAddresses"`
					ExcludedEmailAddresses interface{} `json:"ExcludedEmailAddresses"`
					PermittedURIDomains interface{} `json:"PermittedURIDomains"`
					ExcludedURIDomains interface{} `json:"ExcludedURIDomains"`
					CRLDistributionPoints interface{} `json:"CRLDistributionPoints"`
					PolicyIdentifiers interface{} `json:"PolicyIdentifiers"`
				} `json:"Certificate"`
				Nonce string `json:"nonce"`
			} `json:"signature_header"`
			Signature string `json:"signature"`
		} `json:"signature_data"`
	} `json:"last_config_block_number"`
	TransactionFilter string `json:"transaction_filter"`
	OrdererKafkaMetadata struct {
		LastOffsetPersisted int64 `json:"last_offset_persisted"`
	} `json:"orderer_kafka_metadata"`
}