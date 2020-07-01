
package sdkfabric

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
	"github.com/wingbaas/platformsrv/certgenerate/fabric" 
)

func UpdateChannelCfg(netCfg public.DeployNetConfig,p GenerateParaSt)(string,error) {  
	var srcCfg SdkFabricCfg 
	for _,org := range netCfg.PeerOrgs {
		srcCfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config-" + org.Name + ".json"
		bytes,err := utils.LoadFile(srcCfgFile)
		if err != nil {
			logger.Errorf("UpdateChannelCfg get srcCfgFile error,org="+org.Name)
			return "",fmt.Errorf("UpdateChannelCfg get srcCfgFile error,org="+org.Name)
		}
		err = json.Unmarshal(bytes,&srcCfg)
		if err != nil {
			logger.Errorf("UpdateChannelCfg unmarshal src cfg error,org="+org.Name)
			return "",fmt.Errorf("UpdateChannelCfg unmarshal src cfg error,org="+org.Name)
		}
		channelMap := getChannelMap(netCfg,p) 
		for k,v := range channelMap {
			if srcCfg.Channels == nil {
				srcCfg.Channels = make(map[string]ChannelField)
			}
			srcCfg.Channels[k] = v
		}
		bytes2, err := json.Marshal(srcCfg)
		if err != nil {
			logger.Errorf("UpdateChannelCfg: Marshal fabric sdk config error") 
			return "",fmt.Errorf("UpdateChannelCfg: Marshal fabric sdk config error")
		} 
		bl,yamlStr := fabric.JsonToYaml(string(bytes2)) 
		if !bl {
			logger.Errorf("UpdateChannelCfg: json2yaml error") 
			return "",fmt.Errorf("UpdateChannelCfg: json2yaml error")
		}
		cfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config-" + org.Name + ".yaml" 
		err = utils.WriteFile(cfgFile,yamlStr)
		if err != nil {
			logger.Errorf("UpdateChannelCfg: write sdk config error")
			return "",fmt.Errorf("UpdateChannelCfg: write sdk config error")
		}
		jsonFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config-" + org.Name + ".json"
		err = utils.WriteFile(jsonFile,string(bytes2))
		if err != nil {
			logger.Errorf("UpdateChannelCfg: write sdk json config error")
			return "",fmt.Errorf("UpdateChannelCfg: write sdk json config error")
		} 
	}
	
	//update the global config

	srcCfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config.json"
	bytes,err := utils.LoadFile(srcCfgFile)
	if err != nil {
		logger.Errorf("UpdateChannelCfg get srcCfgFile error")
		return "",fmt.Errorf("UpdateChannelCfg get srcCfgFile error")
	}
	err = json.Unmarshal(bytes,&srcCfg)
	if err != nil {
		logger.Errorf("UpdateChannelCfg unmarshal src cfg error")
		return "",fmt.Errorf("UpdateChannelCfg unmarshal src cfg error")
	}
	channelMap := getChannelMap(netCfg,p) 
	for k,v := range channelMap {
		if srcCfg.Channels == nil {
			srcCfg.Channels = make(map[string]ChannelField)
		}
		srcCfg.Channels[k] = v
	}
	bytes2, err := json.Marshal(srcCfg)
	if err != nil {
		logger.Errorf("UpdateChannelCfg: Marshal fabric sdk config error") 
		return "",fmt.Errorf("UpdateChannelCfg: Marshal fabric sdk config error")
	} 
	bl,yamlStr := fabric.JsonToYaml(string(bytes2)) 
	if !bl {
		logger.Errorf("UpdateChannelCfg: json2yaml error") 
		return "",fmt.Errorf("UpdateChannelCfg: json2yaml error")
	}
	cfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config.yaml" 
	err = utils.WriteFile(cfgFile,yamlStr)
	if err != nil {
		logger.Errorf("UpdateChannelCfg: write sdk config error")
		return "",fmt.Errorf("UpdateChannelCfg: write sdk config error")
	}
	jsonFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config.json"
	err = utils.WriteFile(jsonFile,string(bytes2))
	if err != nil {
		logger.Errorf("UpdateChannelCfg: write sdk json config error")
		return "",fmt.Errorf("UpdateChannelCfg: write sdk json config error")
	} 
	return "",nil 
}

func UpdateOrgCfg(netCfg public.DeployNetConfig,p GenerateParaSt,chList public.ChannelList)(string,error) {
	for _,org := range netCfg.PeerOrgs {
		var orgCfg public.DeployNetConfig
		orgCfg.OrdererOrgs = append(orgCfg.OrdererOrgs,netCfg.OrdererOrgs...)
		orgCfg.KafkaDeployNode = netCfg.KafkaDeployNode
		orgCfg.ZookeeperDeployNode = netCfg.ZookeeperDeployNode
		orgCfg.PeerOrgs = append(orgCfg.PeerOrgs,org)
		_,err := UpdateAddOrgCfg(orgCfg,p,chList)
		if err != nil {
			logger.Errorf("UpdateOrgCfg failed,org=%s",org.Name)
			return "",fmt.Errorf("UpdateOrgCfg failed,org=%s",org.Name)
		}
	}
	_,err := UpdateAddOrgAllCfg(netCfg,p,chList) 
	if err != nil {
		logger.Errorf("UpdateOrgCfg UpdateAddOrgAllCfg failed")
		return "",fmt.Errorf("UpdateOrgCfg UpdateAddOrgAllCfg failed")
	}
	return "",nil
}

func UpdateAddOrgCfg(netCfg public.DeployNetConfig,p GenerateParaSt,chList public.ChannelList)(string,error) { 
	svMap = make(map[string][]public.ServiceNodePortSt)
	err,sMap := k8s.GetServicesNodePort(p.ClusterId,p.NamespaceId,netCfg)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: GetServicesNodePort failed")
		return "",fmt.Errorf("UpdateAddOrgCfg: GetServicesNodePort failed")
	}
	svMap = sMap
	channelMap := getChannelListMap(netCfg,p,chList)
	orgMap := getOrgMap(netCfg,p)
	err,orderMap := getOrderMap(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: getOrderMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgCfg: getOrderMap failed,id=%s",p.ClusterId)
	}
	err,peerMemberMap := getPeersMemberMap(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: getPeersMemberMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgCfg: getPeersMemberMap failed,id=%s",p.ClusterId)
	}
	err,caMap := getCaMap(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: getCaMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgCfg: getCaMap failed,id=%s",p.ClusterId)
	}

	err,pm := getPeersMatch(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: getPeersMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgCfg: getPeersMatch failed,id=%s",p.ClusterId)
	}
	err,om := getOrdersMatch(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: getOrdersMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgCfg: getOrdersMatch failed,id=%s",p.ClusterId)
	}
	err,cm := getCasMatch(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: getCasMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgCfg: getCasMatch failed,id=%s",p.ClusterId)
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
				Path: utils.BAAS_CFG.KeyStorePath + p.BlockId + "/credentialstore/" + firstOrg.Name,
				CryptoStore: CryptoStoreSt {
					Path: utils.BAAS_CFG.KeyStorePath + p.BlockId + "/cryptostore/" + firstOrg.Name, 
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
					Key: CryptoconfigSt {
						Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + firstOrg.Domain + "/users/Admin@" + firstOrg.Domain + "/tls/client.key",
					},
					Cert: CryptoconfigSt {
						Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + firstOrg.Domain + "/users/Admin@" + firstOrg.Domain + "/tls/client.crt",
					},
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
		logger.Errorf("UpdateAddOrgCfg: Marshal fabric sdk config error") 
		return "",fmt.Errorf("UpdateAddOrgCfg: Marshal fabric sdk config error")
	} 
	// logger.Debug("GenerateCfg: json str=")
	// logger.Debug(string(bytes))
	bl,yamlStr := fabric.JsonToYaml(string(bytes)) 
	if !bl {
		logger.Errorf("UpdateAddOrgCfg: json2yaml error") 
		return "",fmt.Errorf("UpdateAddOrgCfg: json2yaml error")
	}
	// logger.Debug("GenerateCfg: yaml str=")
	// logger.Debug(yamlStr)
	cfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config-" + firstOrg.Name + ".yaml" 
	err = utils.WriteFile(cfgFile,yamlStr)
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: write sdk config error")
		return "",fmt.Errorf("UpdateAddOrgCfg: write sdk config error")
	}
	jsonFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config-" + firstOrg.Name + ".json"
	err = utils.WriteFile(jsonFile,string(bytes))
	if err != nil {
		logger.Errorf("UpdateAddOrgCfg: write sdk json config error")
		return "",fmt.Errorf("UpdateAddOrgCfg: write sdk json config error")
	} 
	return "",nil
}

func UpdateAddOrgAllCfg(netCfg public.DeployNetConfig,p GenerateParaSt,chList public.ChannelList)(string,error) { 
	svMap = make(map[string][]public.ServiceNodePortSt)
	err,sMap := k8s.GetServicesNodePort(p.ClusterId,p.NamespaceId,netCfg)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: GetServicesNodePort failed")
		return "",fmt.Errorf("UpdateAddOrgAllCfg: GetServicesNodePort failed")
	}
	svMap = sMap
	channelMap := getChannelListMap(netCfg,p,chList)
	orgMap := getOrgMap(netCfg,p)
	err,orderMap := getOrderMap(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: getOrderMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgAllCfg: getOrderMap failed,id=%s",p.ClusterId)
	}
	err,peerMemberMap := getPeersMemberMap(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: getPeersMemberMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgAllCfg: getPeersMemberMap failed,id=%s",p.ClusterId)
	}
	err,caMap := getCaMap(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: getCaMap,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgAllCfg: getCaMap failed,id=%s",p.ClusterId)
	}

	err,pm := getPeersMatch(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: getPeersMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgAllCfg: getPeersMatch failed,id=%s",p.ClusterId)
	}
	err,om := getOrdersMatch(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: getOrdersMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgAllCfg: getOrdersMatch failed,id=%s",p.ClusterId)
	}
	err,cm := getCasMatch(netCfg,p)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: getCasMatch,id=%s",p.ClusterId)
		return "",fmt.Errorf("UpdateAddOrgAllCfg: getCasMatch failed,id=%s",p.ClusterId)
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
				Path: utils.BAAS_CFG.KeyStorePath + p.BlockId + "/credentialstore/" + firstOrg.Name,
				CryptoStore: CryptoStoreSt {
					Path: utils.BAAS_CFG.KeyStorePath + p.BlockId + "/cryptostore/" + firstOrg.Name, 
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
					Key: CryptoconfigSt {
						Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + firstOrg.Domain + "/users/Admin@" + firstOrg.Domain + "/tls/client.key",
					},
					Cert: CryptoconfigSt {
						Path: utils.BAAS_CFG.BlockNetCfgBasePath + p.BlockId + "/crypto-config/peerOrganizations/" + firstOrg.Domain + "/users/Admin@" + firstOrg.Domain + "/tls/client.crt",
					},
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
		logger.Errorf("UpdateAddOrgAllCfg: Marshal fabric sdk config error") 
		return "",fmt.Errorf("UpdateAddOrgAllCfg: Marshal fabric sdk config error")
	} 
	// logger.Debug("UpdateAddOrgAllCfg: json str=")
	// logger.Debug(string(bytes))
	bl,yamlStr := fabric.JsonToYaml(string(bytes)) 
	if !bl {
		logger.Errorf("UpdateAddOrgAllCfg: json2yaml error") 
		return "",fmt.Errorf("UpdateAddOrgAllCfg: json2yaml error")
	}
	// logger.Debug("UpdateAddOrgAllCfg: yaml str=")
	// logger.Debug(yamlStr)
	cfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config.yaml" 
	err = utils.WriteFile(cfgFile,yamlStr)
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: write sdk config error")
		return "",fmt.Errorf("UpdateAddOrgAllCfg: write sdk config error")
	}
	jsonFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + p.BlockId + "/network-config.json"
	err = utils.WriteFile(jsonFile,string(bytes))
	if err != nil {
		logger.Errorf("UpdateAddOrgAllCfg: write sdk json config error")
		return "",fmt.Errorf("UpdateAddOrgAllCfg: write sdk json config error")
	} 
	return "",nil
}

