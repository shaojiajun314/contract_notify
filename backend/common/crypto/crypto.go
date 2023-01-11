package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	btc_ecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

const SignatureLength = 64 + 1 // 64 bytes ECDSA signature + 1 byte recovery id

type PrivateKey struct {
	prv *btcec.PrivateKey
}

func GenerateKey() (*PrivateKey, error) {
	prv, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil
	}

	return &PrivateKey{
		prv: prv,
	}, nil
}

func (p *PrivateKey) NewPrivateKeyFromHex(hexKey string) error {
	pk, err := hex.DecodeString(hexKey)
	if err != nil {
		return err
	}
	prv, _ := btcec.PrivKeyFromBytes(pk)

	p.prv = prv
	return nil
}

func (p *PrivateKey) NewPrivateKeyFromBytes(pk []byte) {
	prv, _ := btcec.PrivKeyFromBytes(pk)

	p.prv = prv
}

func (p *PrivateKey) Sign(hash []byte) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}

	if p.prv.ToECDSA().Curve != btcec.S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}

	sig, err := btc_ecdsa.SignCompact(p.prv, hash, false)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (p *PrivateKey) SignEth(hash []byte) ([]byte, error) {
	if signature, e := p.Sign(hash); e!=nil{
		return nil, e
	} else {
		v := signature[0] - 27
		copy(signature, signature[1:])
		signature[64] = v
		return signature, nil
	}
}

func (p *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{
		pubKey: p.prv.PubKey(),
	}
}

func (p *PrivateKey) Bytes() []byte {
	return p.prv.Serialize()
}

var secp256k1halfN = new(big.Int).Rsh(btcec.S256().N, 1)

type PublicKey struct {
	pubKey *btcec.PublicKey
}

func (p *PublicKey) VerifySignature(hash, signature []byte) bool {
	if len(signature) != 65 {
		return false
	}
	pk, _, e := btc_ecdsa.RecoverCompact(signature, hash)
	if e!=nil || !bytes.Equal(pk.SerializeUncompressed(), p.pubKey.SerializeUncompressed()) {// hex.EncodeToString() != hex.EncodeToString() {
		return false
	}
	return true
}

func (p *PublicKey) VerifySEthignature(hash, signature []byte) bool {
	v := signature[64] + 27
	copy(signature[1:], signature)
	signature[0] = v
	return p.VerifySignature(hash, signature)
}

func (p *PublicKey) CompressPubkey() []byte {
	return p.pubKey.SerializeCompressed()
}

func (p *PublicKey) UncompressPubkey() []byte {
	return p.pubKey.SerializeUncompressed()
}

func (p *PublicKey) ToHexString() string {
	b := p.CompressPubkey()

	return hex.EncodeToString(b)
}

func (p *PublicKey) DecompressPubkey(pubkey []byte) error {
	if len(pubkey) != 65 && len(pubkey) != 33 {
		return errors.New("invalid compressed public key length")
	}
	key, err := btcec.ParsePubKey(pubkey)
	if err != nil {
		return err
	}
	p.pubKey = key
	return nil
}

func (p *PublicKey) SigToPub(hash, sig []byte) error {
	if len(sig) != SignatureLength {
		return errors.New("invalid signature")
	}

	pub, _, err := btc_ecdsa.RecoverCompact(sig, hash)
	if err != nil {
		return err
	}

	p.pubKey = pub
	return nil
}

func (p *PublicKey) Ecrecover(hash, sig []byte) ([]byte, error) {
	err := p.SigToPub(hash, sig)
	if err != nil {
		return nil, err
	}
	return p.pubKey.SerializeUncompressed(), nil
}

func (p *PublicKey) EcdsaPubKey() *ecdsa.PublicKey {
	return p.pubKey.ToECDSA()
}
