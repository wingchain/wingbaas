
package sdkfabric

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/certgenerate/fabric"
)

type GenerateParaSt struct {
	ClusterId			string
	NamespaceId			string
	BlockId				string
	ChannelName	 		string 
}

var svMap map[string][]public.ServiceNodePortSt

func GenerateCfg(netCfg public.DeployNetConfig,p GenerateParaSt)(string,error) { 
	svMap = make(map[string][]public.ServiceNodePortSt)
	err,sMap := k8s.GetServicesNodePort(p.ClusterId,p.NamespaceId,netCfg)
	if err != nil {
		logger.Errorf("GenerateCfg: GetServicesNodePort failed")
		return "",fmt.Errorf("GenerateCfg: GetServicesNodePort failed")
	}
	svMap = sMap
	channelMap := getChannelMap(netCfg,p) 
	orgMap := getOrgMap(netCfg,p)
	err,orderMap := getOrderMap(netCfg,p)
	if err != nil {
		logger.Errorf("GenerateCfg: getOrderMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("GenerateCfg: getOrderMap failed,id=%s",p.ClusterId)
	}
	err,peerMemberMap := getPeersMemberMap(netCfg,p)
	if err != nil {
		logger.Errorf("GenerateCfg: getPeersMemberMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("GenerateCfg: getPeersMemberMap failed,id=%s",p.ClusterId)
	}
	err,caMap := getCaMap(netCfg,p)
	if err != nil {
		logger.Errorf("GenerateCfg: getCaMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("GenerateCfg: getCaMap failed,id=%s",p.ClusterId)
	}

	err,pm := getPeersMatch(netCfg,p)
	if err != nil {
		logger.Errorf("GenerateCfg: getPeersMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("GenerateCfg: getPeersMatch failed,id=%s",p.ClusterId)
	}
	err,om := getOrdersMatch(netCfg,p)
	if err != nil {
		logger.Errorf("GenerateCfg: getOrdersMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("GenerateCfg: getOrdersMatch failed,id=%s",p.ClusterId)
	}
	err,cm := getCasMatch(netCfg,p)
	if err != nil {
		logger.Errorf("GenerateCfg: getCasMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("GenerateCfg: getCasMatch failed,id=%s",p.ClusterId)
	}

	var firstOrg public.OrgSpec
	for _,org := range netCfg.PeerOrgs { 
		firstOrg = org
		break
	}

	cfg := SdkFabricCfg {
		Name: "fabric-network",
		Version: "1.0.0",
	    Client: ClientSt {
			Organization: firstOrg.Name, 
			Logging: LoggingSt {
				Level: "info",
			},
			Cryptoconfig: CryptoconfigSt {
				Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config" ,
			},
			CredentialStore: CredentialStoreSt {
				Path: utils.BAAS_CFG.KeyStorePath + p.BlockId + "/credentialstore",
				CryptoStore: CryptoStoreSt {
					Path: utils.BAAS_CFG.KeyStorePath + p.BlockId + "/cryptostore", 
				}, 
			},
			BCCSP: BCCSPSt {
				Security: SecuritySt {
					Enabled: true,
					Default: DefaultSt {
						Provider: "SW",
					},
					HashAlgorithm: "SHA2",
					SoftVerify: true,
					Level: 256,
				},
			},
			TLSCerts: TLSCertsSt {
				SystemCertPool: false,
				Client: ClientTLSCertsSt {
					Keyfile: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + firstOrg.Domain + "/users/User1@" + firstOrg.Domain + "/tls/client.key",
					Certfile: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + firstOrg.Domain + "/users/User1@" + firstOrg.Domain + "/tls/client.crt",
				},
			},
		},
		Channels: channelMap,
		Organizations: orgMap,
		Orderers: orderMap,
		Peers: peerMemberMap, 
		CertificateAuthorities: caMap,
		EntityMatchers: EntityMatchersSt {
			Peer: pm,
			Orderer: om,
			CertificateAuthorities: cm,
		}, 
	}
	bytes, err := json.Marshal(cfg)
	if err != nil {
		logger.Errorf("GenerateCfg: Marshal fabric sdk config error") 
		return "",fmt.Errorf("GenerateCfg: Marshal fabric sdk config error")
	} 
	logger.Debug("GenerateCfg: json str=")
	logger.Debug(string(bytes))
	bl,yamlStr := fabric.JsonToYaml(string(bytes)) 
	if !bl {
		logger.Errorf("GenerateCfg: json2yaml error") 
		return "",fmt.Errorf("GenerateCfg: json2yaml error")
	}
	logger.Debug("GenerateCfg: yaml str=")
	logger.Debug(yamlStr)
	cfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config.yaml" 
	err = utils.WriteFile(cfgFile,yamlStr)
	if err != nil {
		logger.Errorf("GenerateCfg: write sdk config error")
		return "",fmt.Errorf("GenerateCfg: write sdk config error")
	}
	return "",nil
}

func getPeersMatch(netCfg public.DeployNetConfig,p GenerateParaSt)(error,[]MatchFieldSt) {
	cluster,_ := k8s.GetCluster(p.ClusterId)
	if cluster == nil {
		logger.Errorf("getPeersMatcher: get cluster failed,id=%s",p.ClusterId)
		return fmt.Errorf("getPeersMatcher: get cluster failed,id=%s",p.ClusterId),nil
	}
	var peersMatch []MatchFieldSt
	for _,org := range netCfg.PeerOrgs {
		for _,p := range org.Specs {
			var pm MatchFieldSt
			apiPort := k8s.GetNodePort(svMap,strings.ToLower(p.Hostname), "api")
			eventPort := k8s.GetNodePort(svMap,strings.ToLower(p.Hostname), "events")
			pm.Pattern = "(\\w*)" + p.Hostname + "." + org.Domain + "(\\w*)"
			pm.URLSubstitutionExp = cluster.PublicIp + ":" + apiPort
			pm.EventURLSubstitutionExp = cluster.PublicIp + ":" + eventPort
			pm.SslTargetOverrideURLSubstitutionExp = p.Hostname + "." + org.Domain
			pm.MappedHost = p.Hostname + "." + org.Domain
			peersMatch = append(peersMatch,pm)
		}
	}
	return nil,peersMatch
}

func getOrdersMatch(netCfg public.DeployNetConfig,p GenerateParaSt)(error,[]MatchFieldSt) {
	cluster,_ := k8s.GetCluster(p.ClusterId)
	if cluster == nil {
		logger.Errorf("getOrdersMatch: get cluster failed,id=%s",p.ClusterId)
		return fmt.Errorf("getOrdersMatch: get cluster failed,id=%s",p.ClusterId),nil
	}
	var ordersMatch []MatchFieldSt
	for _,org := range netCfg.OrdererOrgs {
		for _,p := range org.Specs {
			var om MatchFieldSt
			var orderPort = k8s.GetNodePort(svMap,p.Hostname, p.Hostname)
			om.Pattern = "(\\w*)" + p.Hostname + "." + org.Domain + "(\\w*)"
			om.URLSubstitutionExp = cluster.PublicIp + ":" + orderPort
			om.SslTargetOverrideURLSubstitutionExp = p.Hostname + "." + org.Domain
			om.MappedHost = p.Hostname + "." + org.Domain
			ordersMatch = append(ordersMatch,om)
		}
	}
	return nil,ordersMatch
}

func getCasMatch(netCfg public.DeployNetConfig,p GenerateParaSt)(error,[]MatchFieldSt) {
	cluster,_ := k8s.GetCluster(p.ClusterId)
	if cluster == nil {
		logger.Errorf("getCasMatch: get cluster failed,id=%s",p.ClusterId)
		return fmt.Errorf("getCasMatch: get cluster failed,id=%s",p.ClusterId),nil
	}
	var casMatch []MatchFieldSt
	for _,org := range netCfg.PeerOrgs {
		var cm MatchFieldSt
		key := strings.ToLower(org.Name + "-ca")
		var caPort = k8s.GetNodePort(svMap,key, key)
		cm.Pattern = "(\\w*)" + "ca." + org.Domain + "(\\w*)"
		cm.URLSubstitutionExp = cluster.PublicIp + ":" + caPort 
		cm.MappedHost = "ca." + org.Domain
		casMatch = append(casMatch,cm)
	}
	return nil,casMatch
}

func getCaMap(netCfg public.DeployNetConfig,p GenerateParaSt)(error,map[string]CaField) {
	cluster,_ := k8s.GetCluster(p.ClusterId)
	if cluster == nil {
		logger.Errorf("getCaMap: get cluster failed,id=%s",p.ClusterId)
		return fmt.Errorf("getCaMap: get cluster failed,id=%s",p.ClusterId),nil
	}
	m := make(map[string]CaField)
	var key string
	for _,org := range netCfg.PeerOrgs {
		var field CaField
		caKey := strings.ToLower(org.Name + "-ca")
		var caPort = k8s.GetNodePort(svMap,caKey, caKey)
		field.URL = "https://" + cluster.PublicIp + ":" + caPort
		field.HTTPOptions = HTTPOptionsSt {
			Verify: false,
		}
		field.Registrar = RegistrarSt {
			EnrollID: "admin",
      		EnrollSecret: "adminpw",
		}
		field.CaName = strings.ToLower(org.Name + "-ca")
		field.TLSCACerts = TLSCACertsSt {
			Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + org.Domain + "/ca/ca." + org.Domain + "-cert.pem",
		}
		key = "ca." + org.Domain
		m[key] = field
	}
	return nil,m
}

func getPeersMemberMap(netCfg public.DeployNetConfig,p GenerateParaSt) (error,map[string]MemberField) {
	cluster,_ := k8s.GetCluster(p.ClusterId)
	if cluster == nil {
		logger.Errorf("getPeersMemberMap: get cluster failed,id=%s",p.ClusterId)
		return fmt.Errorf("getPeersMemberMap: get cluster failed,id=%s",p.ClusterId),nil
	}
	m := make(map[string]MemberField)
	var key string
	for _,org := range netCfg.PeerOrgs {
		for _,member := range org.Specs {
			var field MemberField
			apiPort := k8s.GetNodePort(svMap,strings.ToLower(member.Hostname), "api")
			eventPort := k8s.GetNodePort(svMap,strings.ToLower(member.Hostname), "events")
			field.URL = cluster.PublicIp + ":" + apiPort
			field.EventURL = cluster.PublicIp + ":" + eventPort
			key = member.Hostname + "." + org.Domain
			field.GrpcOptions = GrpcOptionsSt {
				SslTargetNameOverride: member.Hostname + "." + org.Domain,
				KeepAliveTime: "0s",
				KeepAliveTimeout: "20s",
				KeepAlivePermit: false,
				FailFast: false,
				AllowInsecure: false,
			}
			field.TLSCACerts = TLSCACertsSt {
				Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + org.Domain + "/tlsca/tlsca." + org.Domain + "-cert.pem",
			}
			m[key] = field
		}
	}
	return nil,m
}

func getOrderMap(netCfg public.DeployNetConfig,p GenerateParaSt)(error,map[string]MemberField) {
	cluster,_ := k8s.GetCluster(p.ClusterId)
	if cluster == nil {
		logger.Errorf("getOrderMap: get cluster failed,id=%s",p.ClusterId)
		return fmt.Errorf("getOrderMap: get cluster failed,id=%s",p.ClusterId),nil
	}
	m := make(map[string]MemberField)
	var key string
	for _,org := range netCfg.OrdererOrgs {	
		for _,member := range org.Specs {
			var field MemberField
			key = member.Hostname + "." + org.Domain
			field.GrpcOptions = GrpcOptionsSt {
				SslTargetNameOverride: member.Hostname + "." + org.Domain,
				KeepAliveTime: "0s",
				KeepAliveTimeout: "20s",
				KeepAlivePermit: false,
				FailFast: false,
				AllowInsecure: false,
			}
			orderKey := member.Hostname
			orderPort := k8s.GetNodePort(svMap,orderKey,orderKey)
			field.URL = cluster.PublicIp + ":" + orderPort
			field.TLSCACerts = TLSCACertsSt {
				Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/ordererOrganizations/" + org.Domain + "/tlsca/tlsca." + org.Domain + "-cert.pem",
			}
			m[key] = field
		}
	}
	return nil,m
}

func getOrgMap(netCfg public.DeployNetConfig,p GenerateParaSt) map[string]OrgField {
	m := make(map[string]OrgField)
	for _,org := range netCfg.PeerOrgs {
		var field OrgField
		field.Mspid = org.Name + "MSP"
		field.CryptoPath = utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + org.Domain + "/users/User1@" + org.Domain + "/msp"
		for _,p := range org.Specs {
			field.Peers = append(field.Peers,p.Hostname + "." + org.Domain)
		}
		field.CertificateAuthorities = append(field.CertificateAuthorities,"ca." + org.Domain)
		m[org.Name] = field
	}
	return m
}

func getChannelMap(netCfg public.DeployNetConfig,p GenerateParaSt)map[string]ChannelField {
	peerMap := getPeerMap(netCfg)
	m := make(map[string]ChannelField)
	field := ChannelField {
		Peers: peerMap,
		Policies: PoliciesSt {
			QueryChannelConfig: QueryChannelConfigSt {
				MinResponses: 1,
				MaxTargets: 1,
				RetryOpts: RetryOptsSt {
					Attempts: 5,
					InitialBackoff: "500ms",
					MaxBackoff: "5s",
					BackoffFactor: 2.0,
				},
			},
		},
	}
	m[p.ChannelName] = field
	return m
}

func getPeerMap(netCfg public.DeployNetConfig)map[string]PeerField {
	m := make(map[string]PeerField)
	for _,org := range netCfg.PeerOrgs {
		for _,p := range org.Specs {
			field := PeerField {
				EndorsingPeer: true,
				ChaincodeQuery: true,
				LedgerQuery: true,
				EventSource: true,
			}
			key := p.Hostname + "." + org.Domain
			m[key] = field
		}
	}
	return m
}

