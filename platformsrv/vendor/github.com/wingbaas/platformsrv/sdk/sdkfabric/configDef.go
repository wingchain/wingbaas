
package sdkfabric

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
	// Keyfile interface{} `json:"keyfile"`
	// Certfile interface{} `json:"certfile"`
	Key CryptoconfigSt `json:"key"`
	Cert CryptoconfigSt `json:"cert"`
}

type TLSCertsSt struct {
	SystemCertPool bool `json:"systemCertPool"`
	Client ClientTLSCertsSt `json:"client"`
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
	Peers map[string]PeerField `json:"peers"`
	Policies PoliciesSt `json:"policies"`
}

type OrgField struct {
	Mspid string `json:"mspid"`
	CryptoPath string `json:"cryptoPath"`
	Peers []string `json:"peers"`
	CertificateAuthorities []string `json:"certificateAuthorities"`
	Users map[string]ClientTLSCertsSt `json:"users"` 
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