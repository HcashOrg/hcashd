package mss

import (
	"github.com/LoCCS/mss"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type Signature struct{
	hcashcrypto.SignatureAdapter
	mss.MerkleSig
}

func (s Signature) GetType() int {
	return pqcTypeMSS
}

func (s Signature) Serialize() []byte{
	return s.MerkleSig.Serialize()
}
