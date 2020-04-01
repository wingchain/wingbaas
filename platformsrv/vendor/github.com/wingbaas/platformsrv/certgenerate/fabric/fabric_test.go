
package fabric 

import (
	"testing"
	"fmt"
)

func Test_fabricCertGen(t *testing.T) {
	var str = `
	{
		"OrdererOrgs": [
			{
				"Name": "Orderer",
				"Domain": "orderer.baas.xyz",
				"Specs": [
					{
						"Hostname": "orderer0"
					},
					{
						"Hostname": "orderer1"
					}
				]
			}
		],
		"PeerOrgs": [
			{
				"Name": "Org1",
				"Domain": "Org1.fabric.baas.xyz",
				"Specs": [
					{
						"Hostname": "peer0"
					},
					{
						"Hostname": "peer1"
					}
				],
				"Users": {
					"Count": 1
				}
			},
			{
				"Name": "Org2",
				"Domain": "Org2.fabric.baas.xyz",
				"Specs": [
					{
						"Hostname": "peer0"
					},
					{
						"Hostname": "peer1"
					}
				],
				"Users": {
					"Count": 1
				}
			}
		]
	}
	`
	bl,_ := JsonToYaml(str)
	//bl = Generate(str,"crypto-config","SM2")
	bl = Generate(str,"crypto-config","ECDSA")
	if !bl {
		fmt.Println("generate cert failed")
	}else{
		fmt.Println("generate cert success")
	}


	// var extendStr = `
	// {
	// 	"OrdererOrgs": [
	// 		{
	// 			"Name": "Orderer",
	// 			"Domain": "orderer.baas.xyz",
	// 			"Specs": [
	// 				{
	// 					"Hostname": "orderer0"
	// 				},
	// 				{
	// 					"Hostname": "orderer1"
	// 				}
	// 			]
	// 		}
	// 	],
	// 	"PeerOrgs": [
	// 		{
	// 			"Name": "Org1",
	// 			"Domain": "Org1.fabric.baas.xyz",
	// 			"Specs": [
	// 				{
	// 					"Hostname": "peer0"
	// 				},
	// 				{
	// 					"Hostname": "peer1"
	// 				}
	// 			],
	// 			"Users": {
	// 				"Count": 1
	// 			}
	// 		},
	// 		{
	// 			"Name": "Org2",
	// 			"Domain": "Org2.fabric.baas.xyz",
	// 			"Specs": [
	// 				{
	// 					"Hostname": "peer0"
	// 				},
	// 				{
	// 					"Hostname": "peer1"
	// 				}
	// 			],
	// 			"Users": {
	// 				"Count": 1
	// 			}
	// 		},
	// 		{
	// 			"Name": "Org3",
	// 			"Domain": "Org3.fabric.baas.xyz",
	// 			"Specs": [
	// 				{
	// 					"Hostname": "peer0"
	// 				},
	// 				{
	// 					"Hostname": "peer1"
	// 				}
	// 			],
	// 			"Users": {
	// 				"Count": 1
	// 			}
	// 		}
	// 	]
	// }
	// `
	// // bl = Extend(extendStr,"crypto-config","SM2")
	// bl = Extend(extendStr,"crypto-config","ECDSA")
	// if !bl {
	// 	fmt.Println("Extend cert failed")
	// }else{
	// 	fmt.Println("Extend cert success")
	// }
}
