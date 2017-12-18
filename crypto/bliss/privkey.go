package bliss

import (
	"github.com/LoCCS/bliss"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type PrivateKey struct{
	hcashcrypto.PrivateKeyAdapter
	bliss.PrivateKey
}


// Public returns the PublicKey corresponding to this private key.
func (p PrivateKey) PublicKey() (hcashcrypto.PublicKey) {
	blissPkp := p.PrivateKey.PublicKey()
	pk := &PublicKey{
		PublicKey: *blissPkp,
	}
	return pk
}

// GetType satisfies the bliss PrivateKey interface.
func (p PrivateKey) GetType() int {
	return pqcTypeBliss
}

func (p PrivateKey) Serialize() []byte{
	return p.PrivateKey.Serialize()
}