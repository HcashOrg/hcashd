// Copyright (c) 2015-2016 The Decred developers
// Copyright (c) 2017-2018 The Hcash developers @sammietocat
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Contents
// * TestCurvePointAdd
// * TestRecoverXBigInt
// * TestRecoverXFieldElement
// * TestScalarMult

package edwards

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestCurvePointAdd tests the addition on curve points
func TestCurvePointAdd(t *testing.T) {
	pointHexStrIdx := 0
	pointHexStrSet := []string{
		"4a3f2684abc42977fe50adbb158a9939cc31b210a7c6e6ea4856395ef3e51bf4",
		"eb4c9d80865dc40107846fdbc1e8b3ce3647615877a77b88720d2913adf1f0a4",
		"3e4e5f8276a802b3c8a4bf442a418cc3a435a71ed0f38aa784f9e460d1c57e2e",
		"647f67801b99bf7eb0c00efbc9b63f4246eba59ff21616e85ecf6139e1006f86",
		"213b22c341bc07cac961a066e137f1e43e671a497eba7add362e2abb15475de9",
		"f52f2b3860d5a0a86db4b786dc73d1f3fe29cbfea4b1b0600b58b072d6d25722",
		"4b6b53496e0d93222fb612f02f688914cb93dea6414d510c92a75812976082ed",
		"768fbf82d32f02e26fb7e51bce61abbb6081085026d049fdf42efc6fdc0715be",
		"b66308e6f1080b6f623d8ce4c2537ff3e60d6c5288ae00fdcc9e652ff770d193",
		"e5bbf34b1d308d393b8fee9e3bba1fd66726ceeba2d8be80d46fd74a4eb9187f",
		"a172dd13f4eaf3bd0e3501f3c2edca2ceddbbe05e6a5d503b114d6e9e14522e5",
		"413df123e96ffc6e8d033037cdbe40d70fb9ec17adf547d9f95a2c6e778bb3cb",
		"b3c815edc658038e31fef3e08190bfdfc63640df5e3b490fb50421cc0380bc21",
	}

	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	tpcv := mockUpPointConversionVectors(50)
	for i := range tpcv {
		if i == 0 {
			continue
		}

		x1, y1, err := curve.EncodedBytesToBigIntPoint(tpcv[i-1].bIn)
		// The random point wasn't on the curve.
		if err != nil {
			continue
		}

		x2, y2, err := curve.EncodedBytesToBigIntPoint(tpcv[i].bIn)
		// The random point wasn't on the curve.
		if err != nil {
			continue
		}

		x, y := curve.Add(x1, y1, x2, y2)
		pointEnc := BigIntPointToEncodedBytes(x, y)
		pointEncAsStr := hex.EncodeToString(pointEnc[:])
		pointHexStr := pointHexStrSet[pointHexStrIdx]
		// Assert our results.
		if pointEncAsStr != pointHexStr {
			t.Fatalf("expected %s, got %s", pointEncAsStr, pointHexStr)
		}
		pointHexStrIdx++
	}
}

// TestRecoverXBigInt tests recovering X from a random Y
//	where (X,Y) is a curve point
//  and Y is expressed in a big integer
func TestRecoverXBigInt(t *testing.T) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	for _, vector := range mockUpPointConversionVectors(1000) {
		isNegative := vector.bIn[31]>>7 == 1
		_, y, err := curve.EncodedBytesToBigIntPoint(vector.bIn)
		// flag indicates whether y is no curve
		notOnCurve := (nil != err)

		if notOnCurve {
			y = EncodedBytesToBigInt(vector.bIn)
		}

		x := curve.RecoverXBigInt(isNegative, y)
		if !curve.IsOnCurve(x, y) {
			if !notOnCurve {
				t.Fatalf("expected %v, got %v", true, notOnCurve)
			}
		} else {
			if notOnCurve {
				t.Fatalf("expected %v, got %v", false, notOnCurve)
			}
			b := BigIntPointToEncodedBytes(x, y)
			if !bytes.Equal(vector.bIn[:], b[:]) {
				t.Fatalf("want %s, got %s\n", hex.EncodeToString(vector.bIn[:]), hex.EncodeToString(b[:]))
			}
		}
	}
}

// TestRecoverXFieldElement tests recovering X from a random Y
//	where (X,Y) is a curve point
// and Y is expressed as a field element
func TestRecoverXFieldElement(t *testing.T) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	for _, vector := range mockUpPointConversionVectors(1000) {
		isNegative := vector.bIn[31]>>7 == 1
		_, y, err := curve.EncodedBytesToBigIntPoint(vector.bIn)
		// flag indicates whether y is no curve
		notOnCurve := (nil != err)

		if notOnCurve {
			y = EncodedBytesToBigInt(vector.bIn)
		}

		yFE := EncodedBytesToFieldElement(vector.bIn)
		x := curve.RecoverXFieldElement(isNegative, yFE)
		xBI := FieldElementToBigInt(x)
		if !curve.IsOnCurve(xBI, y) {
			if !notOnCurve {
				t.Fatalf("expected %v, got %v", true, notOnCurve)
			}
		} else {
			if notOnCurve {
				t.Fatalf("expected %v, got %v", false, notOnCurve)
			}

			b := BigIntPointToEncodedBytes(xBI, y)
			if !bytes.Equal(vector.bIn[:], b[:]) {
				t.Fatalf("want %s, got %s\n",
					hex.EncodeToString(vector.bIn[:]), hex.EncodeToString(b[:]))
			}
		}
	}
}

// TestScalarMult tests multiplication on curve
func TestScalarMult(t *testing.T) {
	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	for _, vector := range mockUpScalarMultVec() {
		x, y, _ := curve.EncodedBytesToBigIntPoint(vector.bIn)
		sBig := EncodedBytesToBigInt(vector.s) // We need big endian
		xMul, yMul := curve.ScalarMult(x, y, sBig.Bytes())
		finalPoint := BigIntPointToEncodedBytes(xMul, yMul)
		if !bytes.Equal(vector.bRes[:], finalPoint[:]) {
			t.Fatalf("want %s, got %s",
				hex.EncodeToString(vector.bRes[:]), hex.EncodeToString(finalPoint[:]))
		}
	}
}
