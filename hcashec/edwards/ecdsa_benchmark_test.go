// Copyright (c) 2015-2016 The Decred developers
// Copyright (c) 2017-2018 The Hcash developers@sammietocat
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package edwards

import (
	"math/rand"
	"testing"
)

// BenchmarkSigning benchmarks the signing operation of ECC on TwistedEdwardsCurve
func BenchmarkSigning(b *testing.B) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	r := rand.New(rand.NewSource(54321))
	msg := []byte{
		0xbe, 0x13, 0xae, 0xf4,
		0xe8, 0xa2, 0x00, 0xb6,
		0x45, 0x81, 0xc4, 0xd1,
		0x0c, 0xf4, 0x1b, 0x5b,
		0xe1, 0xd1, 0x81, 0xa7,
		0xd3, 0xdc, 0x37, 0x55,
		0x58, 0xc1, 0xbd, 0xa2,
		0x98, 0x2b, 0xd9, 0xfb,
	}

	numKeys := 1024
	secKeys := mockUpSecKeysByBytes(curve, numKeys)

	for n := 0; n < b.N; n++ {
		randIndex := r.Intn(numKeys - 1)
		//_, _, err := Sign(curve, privKeyList[randIndex], msg)
		_, _, err := Sign(curve, secKeys[randIndex], msg)
		if err != nil {
			panic("sign failure")
		}
	}
}

// BenchmarkSigningNonStandard benchmarks signing operations
// with secret keys derived random scalar
func BenchmarkSigningNonStandard(b *testing.B) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	r := rand.New(rand.NewSource(54321))
	msg := []byte{
		0xbe, 0x13, 0xae, 0xf4,
		0xe8, 0xa2, 0x00, 0xb6,
		0x45, 0x81, 0xc4, 0xd1,
		0x0c, 0xf4, 0x1b, 0x5b,
		0xe1, 0xd1, 0x81, 0xa7,
		0xd3, 0xdc, 0x37, 0x55,
		0x58, 0xc1, 0xbd, 0xa2,
		0x98, 0x2b, 0xd9, 0xfb,
	}

	numKeys := 250
	privKeyList := mockUpSecKeysByScalars(curve, numKeys)

	for n := 0; n < b.N; n++ {
		randIndex := r.Intn(numKeys - 1)
		_, _, err := Sign(curve, privKeyList[randIndex], msg)
		if err != nil {
			panic("sign failure")
		}
	}
}

// BenchmarkVerification benchmarks the verification operations
// on TwistedEdwardsCurve
func BenchmarkVerification(b *testing.B) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()
	r := rand.New(rand.NewSource(54321))

	numSigs := 1024
	sigList := mockUpSigList(curve, numSigs)

	for n := 0; n < b.N; n++ {
		randIndex := r.Intn(numSigs - 1)

		if !Verify(sigList[randIndex].pubkey,
			sigList[randIndex].msg,
			sigList[randIndex].sig.R,
			sigList[randIndex].sig.S) {

			b.Fatalf("verification on %s with pubkey=%s,r=%v,s=%v",
				sigList[randIndex].msg,
				sigList[randIndex].pubkey,
				*sigList[randIndex].sig.R,
				*sigList[randIndex].sig.S)
		}
	}
}
