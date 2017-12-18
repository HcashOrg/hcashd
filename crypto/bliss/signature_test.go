package bliss

import (
	"testing"
	_ "github.com/HcashOrg/hcashd/chaincfg/chainec"
	_ "github.com/HcashOrg/hcashd/crypto"
	"crypto/rand"
	"bytes"
	"golang.org/x/crypto/sha3"
	"github.com/HcashOrg/hcashd/chaincfg/chainec"
)

func TestSignature(t *testing.T) {

	sk, _, err := Bliss.GenerateKey(rand.Reader)
	if err != nil{
		t.Fatal("Error in Generate keys")
	}

	message := make([]byte, 512)
	rand.Read(message)
	sha3.New512()
	hash := sha3.Sum512(message)

	var sig chainec.Signature
	sig, err = Bliss.Sign(sk, hash[:])
	if err != nil{
		t.Fatal("Error in Sign()")
	}

	sigBytes := sig.Serialize()
	restoredSig, err := Bliss.ParseSignature(sigBytes)
	if err != nil{
		t.Fatal("Error in ParseSignature")
	}
	sigBytes2 := restoredSig.Serialize()

	if !bytes.Equal(sigBytes, sigBytes2){
		t.Fatal("Serialization() and ParseSignature() do not match")
	}

	tp := sig.GetType()
	if tp != pqcTypeBliss{
		t.Fatal("GetType() result not matched")
	}

}
