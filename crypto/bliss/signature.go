package bliss

import (
	"github.com/LoCCS/bliss"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type Signature struct{
	hcashcrypto.SignatureAdapter
	bliss.Signature
}

func (s Signature) GetType() int {
	return pqcTypeBliss
}

func (s Signature) Serialize() []byte{
	return s.Signature.Serialize()
}