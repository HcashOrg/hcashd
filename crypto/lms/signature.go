package lms

import (
	"github.com/LoCCS/lms"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type Signature struct{
	hcashcrypto.SignatureAdapter
	lms.MerkleSig
}

func (s Signature) GetType() int {
	return pqcTypeLMS
}

func (s Signature) Serialize() []byte{
	sigBytes, err := s.MerkleSig.Serialize()
	if err != nil{
		return nil
	}
	return sigBytes
}
