package crypto

import (
	"math/big"
	"crypto/ecdsa"
	"github.com/HcashOrg/hcashd/chaincfg/chainec"
)

type PrivateKey interface{
	chainec.PrivateKey
	PublicKey() PublicKey
}

type PublicKey interface{
	chainec.PublicKey
}

type Signature interface{
	chainec.Signature
}

type PublicKeyAdapter struct {
}

type PrivateKeyAdapter struct {
}

type SignatureAdapter struct {
}

//PrivateKeyAdapter
func (pa PrivateKeyAdapter) Serialize() []byte{
	return nil
}

func (pa PrivateKeyAdapter)SerializeSecret() []byte{
	return nil
}

func (pa PrivateKeyAdapter) Public() (*big.Int, *big.Int){
	return nil, nil
}

func (pa PrivateKeyAdapter) PublicKey() PublicKey{
	return nil
}

func (pa PrivateKeyAdapter) GetD() *big.Int{
	return nil
}

func (pa PrivateKeyAdapter) GetType() int{
	return 0
}

//PublicKeyAdapter
func (pa PublicKeyAdapter)Serialize() []byte{
	return nil
}

func (pa PublicKeyAdapter) SerializeUncompressed() []byte{
	return nil
}

func (pa PublicKeyAdapter) SerializeCompressed() []byte{
	return nil
}

func (pa PublicKeyAdapter) SerializeHybrid() []byte{
	return nil
}

func (pa PublicKeyAdapter) ToECDSA() *ecdsa.PublicKey{
	return nil
}

func (pa PublicKeyAdapter) GetCurve() interface{}{
	return nil
}

func (pa PublicKeyAdapter) GetX() *big.Int{
	return nil
}

func (pa PublicKeyAdapter) GetY() *big.Int{
	return nil
}

func (pa PublicKeyAdapter) GetType() int{
	return 0
}

//SignatureAdapter
func (s SignatureAdapter) Serialize() []byte{
	return nil
}

func (s SignatureAdapter) GetR() *big.Int{
	return nil
}

func (s SignatureAdapter ) GetS() *big.Int{
	return nil
}

func GetType() int{
	return 0
}