package utils

import (
	// "bytes"
  	"net/url"
  	"regexp"
  	"strings"
  	// "math/big"
	// "encoding/binary"
	"github.com/ethereum/go-ethereum/common"

	"contract_notify/common/crypto"
	"contract_notify/collector/types"
	// mepBytes "contract_notify/common/bytes"
)

// var (
// 	parsedLogPrefix = []byte("cpl")
// )

// func ParsedLogDataKey(txHash []byte) []byte {
// 	return append(parsedLogPrefix, txHash...)
// }


// func getDigestHash(log *types.Log, params []string, thirdSign []byte) []byte {
// 	bs := [][]byte{}
// 	for _, name := range params {
// 		v := log.ParsedEvent[name]
// 		var b []byte
// 		switch v.(type) {
// 			case string:
// 				b = mepBytes.LeftPadBytes([]byte(v.(string)), 32)
// 				// raw := v.(string)
// 				// if strings.HasPrefix(raw, "0x") {
// 				// 	raw = strings.ToUpper(raw)
// 				// 	b = mepBytes.LeftPadBytes(common.FromHex(raw), 32)
// 				// 	bs = append(bs, b)
// 				// 	continue
// 				// }
// 				// b = mepBytes.LeftPadBytes([]byte(raw), 32)
// 			case common.Address: 
// 				tmp := strings.ToUpper(v.(common.Address).String())
// 				b = mepBytes.LeftPadBytes(common.FromHex(tmp), 32)
// 			case common.Hash: 
// 				tmp := strings.ToUpper(v.(common.Address).String())
// 				b = mepBytes.LeftPadBytes(common.FromHex(tmp), 32)
// 			case *big.Int: 
// 				b = mepBytes.LeftPadBytes(v.(*big.Int).Bytes(), 32)
// 			case int8: b = mepBytes.LeftPadBytes(new(big.Int).SetUint64(uint64(v.(int8))).Bytes(), 32)
// 			case int16: b = mepBytes.LeftPadBytes(new(big.Int).SetUint64(uint64(v.(int16))).Bytes(), 32)
// 			case int32: b = mepBytes.LeftPadBytes(new(big.Int).SetUint64(uint64(v.(int32))).Bytes(), 32)
// 			case int64: 
// 				b = mepBytes.LeftPadBytes(new(big.Int).SetUint64(uint64(v.(int64))).Bytes(), 32)
// 			case []byte: b = mepBytes.LeftPadBytes(v.([]byte), 32)
// 			case bool: 
// 				buf := bytes.NewBuffer([]byte{})
// 				binary.Write(buf, binary.BigEndian, v.(bool))
// 				b = mepBytes.LeftPadBytes(buf.Bytes(), 32)
// 			default:
// 				panic(v)
// 		}
// 		bs = append(bs, b)
// 	}
// 	if thirdSign != nil {
// 		bs = append(bs, mepBytes.LeftPadBytes(thirdSign, 32))
// 	}
// 	hash := crypto.Keccak256Hash(bs...)
// 	signedHash := crypto.Keccak256Hash(append([]byte("\x19Ethereum Signed Message:\n32"), hash.Bytes()...))
// 	return signedHash.Bytes()
// }

// func Sign(log *types.Log, thirdSign []byte, params []string, privateKeyB []byte) ([]byte, error) {
// 	p := &crypto.PrivateKey{}
// 	p.NewPrivateKeyFromBytes(privateKeyB)
// 	signedHash := getDigestHash(log, params, thirdSign)
// 	signature, err := p.Sign(signedHash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return signature, nil
// }


// func VerifySign(log *types.Log, params []string, pk *crypto.PublicKey, signature []byte) bool {
// 	signedHash := getDigestHash(log, params, nil)
// 	return pk.VerifySignature(signedHash, signature)
// }

func getDigestHash(topics []string, data string) []byte {
	bs := [][]byte{}
	for _, t := range topics {
		bs = append(bs, common.FromHex(t))
	}
	bs = append(bs, common.FromHex(data))
	hash := crypto.Keccak256Hash(bs...)
	return hash.Bytes()
}

func Sign(
	topics []string,
	data string,
	p *crypto.PrivateKey,
	// p *crypto.PrivateKey,
) ([]byte, error) {
	signedHash := getDigestHash(topics, data)
	signature, e := p.SignEth(signedHash)
	if e != nil {
		panic(e)
	}
	return signature, nil
	
}

func VerifySign(log *types.Log, pk *crypto.PublicKey, signature []byte) bool {
	signedHash := getDigestHash(log.Topics, log.Data)
	return pk.VerifySEthignature(signedHash, signature)
}


func CheckHttpUrl(addr string) bool {
	if 0 == len(addr) {
		return false
	}
	parse, err := url.Parse(addr)
	if err != nil {
		return false
	}

	if "http" != parse.Scheme && "https" != parse.Scheme {
		return false
	}
	re := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$|^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)+([A-Za-z]|[A-Za-z][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	u := strings.Split(parse.Host, ":")
	if 2 == len(u) && 0 != len(u[1]) {
		result := re.FindAllStringSubmatch(u[0], -1)
		if result == nil {
			return false
		}
		return true
	}
	result := re.FindAllStringSubmatch(parse.Host, -1)
	if result == nil {
		return false
	}

	return true
}

// func SaveParsedLog(db db.Database, txHash string, paseredData []byte) error {
// 	h, e := hex.Decode(txHash)
// 	if e!=nil {
// 		return e
// 	}
// 	return db.Put(ParsedLogDataKey(h), paseredData)
// }

// func GetParsedLog(db db.Database, txHash string) (map[string]interface{}, error) {
// 	h, e := hex.Decode(txHash)
// 	if e!=nil {
// 		return nil, e
// 	}
// 	if rawData, e := db.Get(ParsedLogDataKey(h)); e!=nil {
// 		return nil, e
// 	} else {
// 		var ret map[string]interface{}
// 		if e := json.Unmarshal(rawData, &ret); e!=nil{
// 			return nil, e
// 		}
// 		return ret, nil
// 	}
// }
