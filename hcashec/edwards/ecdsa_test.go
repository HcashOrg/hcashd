// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package edwards

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
)

// TestBadPubKey tests failed verification due to bad pubkeys
func TestBadPubKey(t *testing.T) {
	tRand := rand.New(rand.NewSource(54321))

	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()
	msg := []byte("Hello World in TestBadSignature")

	sks := mockUpSecKeysByScalars(curve, 50)
	for _, sk := range sks {
		r, s, err := Sign(curve, sk, msg)
		if nil != err {
			t.Fatalf("unexpected signing error: %s\n", err)
		}

		// derive the public key
		pkX, pkY := sk.Public()
		pk := NewPublicKey(curve, pkX, pkY)

		// Screw up a random bit in pk and
		// make sure it still fails.
		pkBadBytes := pk.Serialize()
		pos := tRand.Intn(31)
		bitPos := tRand.Intn(7)
		if pos == 0 {
			// 0th bit in first byte doesn't matter
			bitPos = tRand.Intn(6) + 1
		}
		pkBadBytes[pos] ^= 1 << uint8(bitPos)

		// parse and verify the bad pk
		if pkBad, err := ParsePubKey(curve, pkBadBytes); (nil == err) && (nil != pkBad) {

			if Verify(pkBad, msg, r, s) {
				t.Fatalf("verification should fail on %s with pk=%v,r=%v,s=%v\n",
					msg, pk, *r, *s)
			}
		}
	}
}

// TestBadSignature tests failed verification on bad signatures
func TestBadSignature(t *testing.T) {
	tRand := rand.New(rand.NewSource(54321))

	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()
	msg := []byte("Hello World in TestBadSignature")

	sks := mockUpSecKeysByScalars(curve, 50)
	for _, sk := range sks {
		r, s, err := Sign(curve, sk, msg)
		if nil != err {
			t.Fatalf("unexpected signing error: %s\n", err)
		}

		// derive the public key
		pkX, pkY := sk.Public()
		pk := NewPublicKey(curve, pkX, pkY)

		// Screw up a random bit in the signature and
		// make sure it still fails.
		sig := NewSignature(r, s)
		sigBadBytes := sig.Serialize()
		pos := tRand.Intn(63)
		bitPos := tRand.Intn(7)
		// flip a bit randomly on the a random byte
		sigBadBytes[pos] ^= 1 << uint8(bitPos)

		// parse a signature from sigBadBytes
		// and check by verification
		if badSig, err := ParseSignature(curve, sigBadBytes); nil == err {
			// should fail to verify
			if Verify(pk, msg, badSig.GetR(), badSig.GetS()) {
				t.Fatalf("verification should fail on %s with pk=%v,r=%v,s=%v",
					msg, pk, *badSig.GetR(), *badSig.GetS())
			}
		}
	}
}

// TestGolden tests the implemented TwistedEdwardsCurve against other implementation
func TestGolden(t *testing.T) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	// sign.input.gz is a selection of test cases from
	// http://ed25519.cr.yp.to/python/sign.input
	testDataZ, err := os.Open("testdata/sign.input.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer testDataZ.Close()
	testData, err := gzip.NewReader(testDataZ)
	if err != nil {
		t.Fatal(err)
	}
	defer testData.Close()

	in := bufio.NewReaderSize(testData, 1<<12)
	lineNo := 0
	for {
		lineNo++
		lineBytes, err := in.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("error reading test data: %s", err)
		}

		// each line is in inform of
		//	private-key-bytes:public-key-bytes:message:signature
		line := string(lineBytes)
		parts := strings.Split(line, ":")
		if len(parts) != 5 {
			t.Fatalf("bad number of parts on line %d (want %v, got %v)", lineNo,
				5, len(parts))
		}

		privBytes, _ := hex.DecodeString(parts[0])
		privArray := copyBytes64(privBytes)

		pubKeyBytes, _ := hex.DecodeString(parts[1])
		pubArray := copyBytes(pubKeyBytes)
		msg, _ := hex.DecodeString(parts[2])
		sig, _ := hex.DecodeString(parts[3])
		sigArray := copyBytes64(sig)
		// The signatures in the test vectors also include the message
		// at the end, but we just want R and S.
		sig = sig[:SignatureSize]

		if l := len(pubKeyBytes); l != PubKeyBytesLen {
			t.Fatalf("bad public key length on line %d: got %d bytes", lineNo, l)
		}

		var priv [PrivKeyBytesLen]byte
		copy(priv[:], privBytes)
		copy(priv[32:], pubKeyBytes)

		// Deserialize privkey and test functions.
		privkeyS1, pubkeyS1 := PrivKeyFromSecret(curve, priv[:32])
		privkeyS2, pubkeyS2 := PrivKeyFromBytes(curve, priv[:])
		pkS1 := privkeyS1.SerializeSecret()
		pkS2 := privkeyS2.SerializeSecret()
		pubkS1 := pubkeyS1.Serialize()
		pubkS2 := pubkeyS2.Serialize()
		cmp := bytes.Equal(pkS1[:], pkS2[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		cmp = bytes.Equal(privArray[:], copyBytes64(pkS1)[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		cmp = bytes.Equal(privArray[:], copyBytes64(pkS2)[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		cmp = bytes.Equal(pubkS1[:], pubkS2[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		cmp = bytes.Equal(pubArray[:], copyBytes(pubkS1)[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		cmp = bytes.Equal(pubArray[:], copyBytes(pubkS2)[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Deserialize pubkey and test functions.
		pubkeyP, err := ParsePubKey(curve, pubKeyBytes)
		pubkP := pubkeyP.Serialize()
		cmp = bytes.Equal(pubkS1[:], pubkP[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		cmp = bytes.Equal(pubkS2[:], pubkP[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		cmp = bytes.Equal(pubArray[:], copyBytes(pubkP)[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Deserialize signature and test functions.
		internalSig, err := ParseSignature(curve, sig)
		iSigSerialized := internalSig.Serialize()
		cmp = bytes.Equal(sigArray[:], copyBytes64(iSigSerialized)[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		sig2r, sig2s, err := Sign(curve, privkeyS2, msg)
		sig2 := &Signature{sig2r, sig2s}
		sig2B := sig2.Serialize()
		if !bytes.Equal(sig, sig2B[:]) {
			t.Errorf("different signature result on line %d: %x vs %x", lineNo,
				sig, sig2B[:])
		}

		var pubKey [PubKeyBytesLen]byte
		copy(pubKey[:], pubKeyBytes)
		if !Verify(pubkeyP, msg, sig2r, sig2s) {
			t.Errorf("signature failed to verify on line %d", lineNo)
		}
	}
}

// TestKeySerialization tests the serialization/deserialization
//	operations on keys
func TestKeySerialization(t *testing.T) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	sks := mockUpSecKeysByScalars(curve, 50)
	for _, sk := range sks {
		// derive a secret key from a scalar
		sk1, _, err := PrivKeyFromScalar(curve,
			copyBytes(sk.ecPk.D.Bytes())[:])
		if nil != err {
			t.Fatalf("unexpected signing error: %s\n", err)
		}
		// check equality
		if 0 != sk.GetD().Cmp(sk1.GetD()) {
			t.Fatalf("want %v, got %v\n", *sk.GetD(), *sk1.GetD())
		}

		//derive a secret key from a scalar
		sk2, _, err := PrivKeyFromScalar(curve, sk.Serialize())
		if nil != err {
			t.Fatalf("unexpected signing error: %s\n", err)
		}
		// check equality
		if 0 != sk.GetD().Cmp(sk2.GetD()) {
			t.Fatalf("want %v, got %v\n", *sk.GetD(), *sk2.GetD())
		}
	}
}

// TestStdSigning tests the standard signing/verifying operations
func TestStdSigning(t *testing.T) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()
	msg := []byte("Hello World in TestStdSigning")

	sks := mockUpSecKeysByScalars(curve, 50)
	for _, sk := range sks {
		r, s, err := Sign(curve, sk, msg)
		if nil != err {
			t.Fatalf("unexpected signing error: %s\n", err)
		}

		pkX, pkY := sk.Public()
		pk := NewPublicKey(curve, pkX, pkY)
		if !Verify(pk, msg, r, s) {
			t.Fatalf("verification failed on %s with pk=%v,r=%v,s=%v", msg, pk, r, s)
			//t.Fatalf("expected %v, got %v", true, ok)
		}
	}
}
