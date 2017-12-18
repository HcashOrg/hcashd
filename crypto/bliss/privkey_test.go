package bliss
import (
	"testing"
	_ "github.com/HcashOrg/hcashd/chaincfg/chainec"
	_ "github.com/HcashOrg/hcashd/crypto"
	"crypto/rand"
	"bytes"
)

func TestPrivateKey(t *testing.T) {


=======
>>>>>>> dev-test-postquantum
	sk, pk, err := Bliss.GenerateKey(rand.Reader)
	if err != nil{
		t.Fatal("Error in Generate keys")
	}

	pk2 := sk.PublicKey()
	pk3 := sk.(*PrivateKey).PrivateKey.PublicKey()
	pkBytes := pk.Serialize()
	pkBytes2 := pk2.Serialize()
	pkBytes3 := pk3.Serialize()
	skBytes := sk.Serialize()
	skBytes2 := sk.(*PrivateKey).PrivateKey.Serialize()

	if !bytes.Equal(pkBytes, pkBytes2){
		t.Fatal("Error in PublicKey(), the result is not same as the result of Generatekey()")
	}

	if !bytes.Equal(pkBytes, pkBytes3){
		t.Fatalf("Generated Public Key is not same as the result of the PublicKey() of bliss privateKey")
	}

	if !bytes.Equal(skBytes, skBytes2){
		t.Fatalf("Error in Serialization(), the result is not same as the result of function in Bliss")
	}

	prk, _ := Bliss.PrivKeyFromBytes(skBytes)
	skBytes3 := prk.Serialize()

	if !bytes.Equal(skBytes, skBytes3){
		t.Fatalf("serilization() and PrivKeyFromBytes() do not match")
	}

	tp := sk.GetType()
	if tp != pqcTypeBliss{
		t.Fatal("GetType() result not matched")
	}

}
