package bliss

import (
	"github.com/LoCCS/bliss"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type PublicKey struct{
	hcashcrypto.PublicKeyAdapter
	bliss.PublicKey
}

func (p PublicKey) GetType() int {
	return pqcTypeBliss
}

func (p PublicKey) Serialize() []byte{
	return p.PublicKey.Serialize()
}

func (p PublicKey) SerializeCompressed() []byte{
	return p.Serialize()
}

func (p PublicKey) SerializeUnCompressed() []byte{
	return p.Serialize()
}

