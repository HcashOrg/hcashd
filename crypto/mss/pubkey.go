package mss

import (
	"github.com/LoCCS/bliss"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type PublicKey struct{
	hcashcrypto.PublicKeyAdapter
	root []byte
}

func (p PublicKey) GetType() int {
	return pqcTypeMSS
}

func (p PublicKey) Serialize() []byte{
	return p.root
}

func (p PublicKey) SerializeCompressed() []byte{
	return p.root
}

func (p PublicKey) SerializeUnCompressed() []byte{
	return p.root
}