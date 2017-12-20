package mss


import (
	"io"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
	"github.com/LoCCS/mss"
	"fmt"
)

var pqcTypeMSS = 5

type mssDSA struct {

	// Private keys
	//newPrivateKey     func(s1, s2, a *poly.PolyArray) hcashcrypto.PrivateKey
	privKeyFromBytes  func(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey)
	privKeyBytesLen   func() int

	// Public keys
	//newPublicKey               func(a *poly.PolyArray) hcashcrypto.PublicKey
	parsePubKey                func(pubKeyStr []byte) (hcashcrypto.PublicKey, error)
	pubKeyBytesLen             func() int

	// Signatures
	//newSignature      func(z1, z2 *poly.PolyArray, c []uint32) hcashcrypto.Signature
	parseDERSignature func(sigStr []byte) (hcashcrypto.Signature, error)
	parseSignature    func(sigStr []byte) (hcashcrypto.Signature, error)
	recoverCompact    func(signature, hash []byte) (hcashcrypto.PublicKey, bool, error)

	//
	generateKey func(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey, error)
	sign        func(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error)
	verify      func(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool

	// Symmetric cipher encryption
	//generateSharedSecret func(privkey []byte, x, y *big.Int) []byte
	//encrypt              func(x, y *big.Int, in []byte) ([]byte, error)
	//decrypt              func(privkey []byte, in []byte) ([]byte, error)
}

// Private keys
//func (sp mssDSA) NewPrivateKey(s1, s2, a *poly.PolyArray) hcashcrypto.PrivateKey {
//	return sp.newPrivateKey(s1, s2, a)
//}
func (sp mssDSA) PrivKeyFromBytes(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey) {
	return sp.privKeyFromBytes(pk)
}
func (sp mssDSA) PrivKeyBytesLen() int {
	return sp.privKeyBytesLen()
}

// Public keys
//func (sp mssDSA) NewPublicKey(a *poly.PolyArray) hcashcrypto.PublicKey {
//	return sp.newPublicKey(a)
//}
func (sp mssDSA) ParsePubKey(pubKeyStr []byte) (hcashcrypto.PublicKey, error) {
	return sp.parsePubKey(pubKeyStr)
}
func (sp mssDSA) PubKeyBytesLen() int {
	return sp.pubKeyBytesLen()
}

// Signatures
//func (sp mssDSA) NewSignature(z1, z2 *poly.PolyArray, c []uint32) hcashcrypto.Signature {
//	return sp.newSignature(z1, z2, c)
//}
func (sp mssDSA) ParseDERSignature(sigStr []byte) (hcashcrypto.Signature, error) {
	return sp.parseDERSignature(sigStr)
}
func (sp mssDSA) ParseSignature(sigStr []byte) (hcashcrypto.Signature, error) {
	return sp.parseSignature(sigStr)
}
func (sp mssDSA) RecoverCompact(signature, hash []byte) (hcashcrypto.PublicKey, bool,
	error) {
	return sp.recoverCompact(signature, hash)
}

// MSS
func (sp mssDSA) GenerateKey(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey,
	error) {
	return sp.generateKey(rand)
}
func (sp mssDSA) Sign(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error) {
	return sp.sign(priv, hash)
}
func (sp mssDSA) Verify(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool {
	return sp.verify(pub, hash, sig)
}

func newMSSDSA() DSA {
	var mss DSA = &mssDSA{
		privKeyFromBytes: func(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey) {
			fmt.Println("privKeyFromBytes is called")
			return nil, nil
		},
		privKeyBytesLen: func() int {
			return MSSPrivKeyLen
		},
		parsePubKey: func(pubKeyStr []byte) (hcashcrypto.PublicKey, error) {
			return &PublicKey{
				root: pubKeyStr,
			}, nil
		},
		pubKeyBytesLen: func() int {
			return MSSPubKeyLen
		},
		parseDERSignature: func(sigStr []byte) (hcashcrypto.Signature, error) {
			sig := mss.DeserializeMerkleSig(sigStr)
			return &Signature{
				MerkleSig: *sig,
			}, nil
		},
		parseSignature: func(sigStr []byte) (hcashcrypto.Signature, error) {
			sig := mss.DeserializeMerkleSig(sigStr)
			return &Signature{
				MerkleSig: *sig,
			}, nil
		},
		recoverCompact: func(signature, hash []byte) (hcashcrypto.PublicKey, bool, error) {
			return nil, false, nil
		},
		generateKey: func(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey, error) {
			fmt.Println("genereate key is called")
			return nil, nil, nil
		},
		sign: func(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error) {
			mssPrv := priv.(PrivateKey).MerkleAgent
			_, sig, err := mss.Sign(&mssPrv, hash)
			if err != nil{
				return nil, err
			}
			return &Signature{
				MerkleSig: *sig,
			}, nil
		},

		verify: func(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool {
			pbBytes := pub.(PublicKey).root
			signature := sig.(*Signature)
			mssSig := signature.MerkleSig
			result := mss.Verify(pbBytes, hash, &mssSig)
			return result
		},
	}

	return mss.(DSA)
}