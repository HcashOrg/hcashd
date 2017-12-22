package lms

import (
	"io"
	hcashcrypto "github.com/HcashOrg/hcashd/crypto"
)

type DSA interface {

	// ----------------------------------------------------------------------------
	// Private keys
	//
	// NewPrivateKey instantiates a new private key for the given data
	NewPrivateKey() hcashcrypto.PrivateKey

	// PrivKeyFromBytes calculates the public key from serialized bytes,
	// and returns both it and the private key.
	PrivKeyFromBytes(pk []byte) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey)

	// PrivKeyBytesLen returns the length of a serialized private key.
	PrivKeyBytesLen() int

	// ----------------------------------------------------------------------------
	// Public keys
	//
	// NewPublicKey instantiates a new public key (point) for the given data.
	NewPublicKey() hcashcrypto.PublicKey

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
	NewSignature() hcashcrypto.Signature

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
	// LMS
	//
	// GenerateKey generates a new private and public keypair from the
	// given reader.
	GenerateKey(rand io.Reader) (hcashcrypto.PrivateKey, hcashcrypto.PublicKey, error)

	// Sign produces a LMS signature using a private key and a message.
	Sign(priv hcashcrypto.PrivateKey, hash []byte) (hcashcrypto.Signature, error)

	// Verify verifies a LMS signature against a given message and
	// public key.
	Verify(pub hcashcrypto.PublicKey, hash []byte, sig hcashcrypto.Signature) bool
}

const (
	LMSTypeLMS = 5

	LMSVersion = 1

	LMSPubKeyLen = 32

	LMSPrivKeyLen = 4691
)

var LMS = newLMSDSA()