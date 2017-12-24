package lms

import (
	"io"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
	"github.com/LoCCS/lms"
	"fmt"
	"golang.org/x/crypto/sha3"
)

var pqcTypeLMS = 5

type lmsDSA struct {

	// Private keys
	newPrivateKey     func() hcashcrypto.PrivateKey
	privKeyFromBytes  func(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey)
	privKeyBytesLen   func() int

	// Public keys
	newPublicKey               func() hcashcrypto.PublicKey
	parsePubKey                func(pubKeyStr []byte) (hcashcrypto.PublicKey, error)
	pubKeyBytesLen             func() int

	// Signatures
	newSignature      func() hcashcrypto.Signature
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
func (sp lmsDSA) NewPrivateKey() hcashcrypto.PrivateKey {
	return sp.newPrivateKey()
}
func (sp lmsDSA) PrivKeyFromBytes(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey) {
	return sp.privKeyFromBytes(pk)
}
func (sp lmsDSA) PrivKeyBytesLen() int {
	return sp.privKeyBytesLen()
}

// Public keys
func (sp lmsDSA) NewPublicKey() hcashcrypto.PublicKey {
	return sp.newPublicKey()
}
func (sp lmsDSA) ParsePubKey(pubKeyStr []byte) (hcashcrypto.PublicKey, error) {
	return sp.parsePubKey(pubKeyStr)
}
func (sp lmsDSA) PubKeyBytesLen() int {
	return sp.pubKeyBytesLen()
}

// Signatures
func (sp lmsDSA) NewSignature() hcashcrypto.Signature {
	return sp.newSignature()
}
func (sp lmsDSA) ParseDERSignature(sigStr []byte) (hcashcrypto.Signature, error) {
	return sp.parseDERSignature(sigStr)
}
func (sp lmsDSA) ParseSignature(sigStr []byte) (hcashcrypto.Signature, error) {
	return sp.parseSignature(sigStr)
}
func (sp lmsDSA) RecoverCompact(signature, hash []byte) (hcashcrypto.PublicKey, bool,
	error) {
	return sp.recoverCompact(signature, hash)
}

// LMS
func (sp lmsDSA) GenerateKey(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey,
	error) {
	return sp.generateKey(rand)
}
func (sp lmsDSA) Sign(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error) {
	return sp.sign(priv, hash)
}
func (sp lmsDSA) Verify(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool {
	return sp.verify(pub, hash, sig)
}

func newLMSDSA() DSA {
	var lms DSA = &lmsDSA{
		privKeyFromBytes: func(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey) {
			fmt.Println("privKeyFromBytes is called")
			return nil, nil
		},
		privKeyBytesLen: func() int {
			return LMSPrivKeyLen
		},
		parsePubKey: func(pubKeyStr []byte) (hcashcrypto.PublicKey, error) {
			return &PublicKey{
				root: pubKeyStr,
			}, nil
		},
		pubKeyBytesLen: func() int {
			return LMSPubKeyLen
		},
		parseDERSignature: func(sigStr []byte) (hcashcrypto.Signature, error) {
			sig := lms.DeserializeMerkleSig(sigStr)
			return &Signature{
				MerkleSig: *sig,
			}, nil
		},
		parseSignature: func(sigStr []byte) (hcashcrypto.Signature, error) {
			sig := lms.DeserializeMerkleSig(sigStr)
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

			sha3.New256()
			messageHash := sha3.Sum256(hash)

			lmsPrv := priv.(PrivateKey).MerkleAgent
			_, sig, err := lms.Sign(&lmsPrv, messageHash[:])
			if err != nil{
				return nil, err
			}
			return &Signature{
				MerkleSig: *sig,
			}, nil
		},

		verify: func(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool {
			sha3.New256()
			messageHash := sha3.Sum256(hash)
			pbBytes := pub.(*PublicKey).root
			signature := sig.(*Signature)
			lmsSig := signature.MerkleSig
			result := lms.Verify(pbBytes, messageHash[:], &lmsSig)
			return result
		},
	}

	return lms.(DSA)
}