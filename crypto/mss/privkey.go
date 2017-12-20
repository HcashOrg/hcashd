package mss


import (
	"github.com/LoCCS/mss"
	//wots "github.com/LoCCS/mss/ots/winternitz"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type PrivateKey struct{
	hcashcrypto.PrivateKeyAdapter
	mss.MerkleAgent
}


// Public returns the PublicKey corresponding to this private key.
func (p PrivateKey) PublicKey() (hcashcrypto.PublicKey) {
	root := p.MerkleAgent.Root()
	pk := &PublicKey{
		root:root,
	}
	return pk
}

// GetType satisfies the bliss PrivateKey interface.
func (p PrivateKey) GetType() int {
	return pqcTypeMSS
}

func (p PrivateKey) Serialize() []byte{
	return p.MerkleAgent.SerializeSecret()
}
