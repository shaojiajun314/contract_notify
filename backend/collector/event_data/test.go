package eventdata

import (
	"contract_notify/types"
	"contract_notify/common"
	"math/big"
)

func test() types.Transactions {
	txs := types.Transactions([]*types.Transaction{
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("0xeef8ea6b237fcf422d10772c9800ea83f8f2fc8616cbb107918d00cca9cf6fab"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("0x1651fea220327ba522a8eb6fc3d20eadf1745a9d97b2d48f83a7e8af9a7ae7c2"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("0x73d7c1b54b9f5fadf962d553f0690d805f315f27d72241be6151ce769782849e"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09d4c"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09a4a"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09dca"),
					ThirdSign: "0xc01e299cde8f9960c3446e9b35164e9add80ec0d59b6b070057eac75834f65124a4440cb7aaa79db2129beffbf202ec54c7cc29adf57cf652a5eb7d36e9379ba00",
					Data: types.TargetData {
						ParsedEvent: []byte("{\"arg0\": 0, \"arg1\": 5, \"arg2\": 97, \"arg3\": \"0x997aa678bb7c561f608815255ad25d3108c239b2\", \"arg4\": \"0x997aa678bb7c561f608815255ad25d3108c239b2\", \"arg5\": \"0x35aee3799c88340a9210815c7b102706f0e9d59b\", \"arg6\": \"0xf85a14f2355dfeeb0312347e5f6583703c07ec4a\", \"arg7\": 1000000000000000000, \"arg8\": 1000000000000000000, \"arg9\": \"0x3acb4ec04d9586f7e37b635da5b22c074e884bd7\"}"),
						Topics: []string{"0x615603b669fbccbe237b413b2efaca060fbdc19c9a64b78664c20d095a40ccaf"},
						Data: "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000061000000000000000000000000997aa678bb7c561f608815255ad25d3108c239b2000000000000000000000000997aa678bb7c561f608815255ad25d3108c239b200000000000000000000000035aee3799c88340a9210815c7b102706f0e9d59b000000000000000000000000f85a14f2355dfeeb0312347e5f6583703c07ec4a0000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000003acb4ec04d9586f7e37b635da5b22c074e884bd7",
					},
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09d9a"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09d1a"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09d5a"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09d4a"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09d3a"),
				},
			},
			Signature: "", // warnning
		},
		&types.Transaction{
			RawTransaction: &types.RawTransaction{
				Sender: common.Address{}, // warnning
				Nonce: 0, // warnning
				Fee: new(big.Int).SetUint64(0), // warnning
				Payload: &types.EventTaskData {
					TaskHash: common.HexToAddress("f55c0472647ad1cb5e3e056ffde3183f295a469347e0fd34f033f1f46ec09d4f"),
				},
			},
			Signature: "", // warnning
		},

	})
	return txs
}