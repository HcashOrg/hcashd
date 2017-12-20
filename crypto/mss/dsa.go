package mss


import (
	"io"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type DSA interface {

	// ----------------------------------------------------------------------------
	// Private keys
	//
	// NewPrivateKey instantiates a new private key for the given data

	// PrivKeyFromBytes calculates the public key from serialized bytes,
	// and returns both it and the private key.
	PrivKeyFromBytes(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey)

	// PrivKeyBytesLen returns the length of a serialized private key.
	PrivKeyBytesLen() int

	// ----------------------------------------------------------------------------
	// Public keys
	//
	// NewPublicKey instantiates a new public key (point) for the given data.
	//NewPublicKey(a *poly.PolyArray) hcashcrypto.PublicKey

	// ParsePubKey parses a serialized public key for the given
	// curve and returns a public key.
	ParsePubKey(pubKeyStr []byte) (hcashcrypto.PublicKey, error)

	// PubKeyBytesLen returns the length of the default serialization
	// method for a public key.
	PubKeyBytesLen() int

	// ----------------------------------------------------------------------------
	// Signatures
	//
	// NewSignature instantiates a new signature
	//NewSignature(z1, z2 *poly.PolyArray, c []uint32) hcashcrypto.Signature

	// ParseDERSignature parses a DER encoded signature .
	// If the method doesn't support DER signatures, it
	// just parses with the default method.
	ParseDERSignature(sigStr []byte) (hcashcrypto.Signature, error)

	// ParseSignature a default encoded signature
	ParseSignature(sigStr []byte) (hcashcrypto.Signature, error)

	// RecoverCompact recovers a public key from an encoded signature
	// and message, then verifies the signature against the public
	// key.
	RecoverCompact(signature , hash []byte) (hcashcrypto.PublicKey, bool, error)

	// ----------------------------------------------------------------------------
	// MSS
	//
	// GenerateKey generates a new private and public keypair from the
	// given reader.
	GenerateKey(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey, error)

	// Sign produces a MSS signature using a private key and a message.
	Sign(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error)

	// Verify verifies a MSS signature against a given message and
	// public key.
	Verify(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool
}

const (
	MSSTypeMSS = 5

	MSSVersion = 1

	MSSPubKeyLen = 32

	MSSPrivKeyLen = 417
)

var MSS = newMSSDSA()