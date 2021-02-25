package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"golang.org/x/crypto/cryptobyte"
)

func TestCheckByteArray(t *testing.T) {
	// t.SkipNow()
	bStr := "01075bea74ced3d3b42cc9184efbe14c13e0307d4b1fbe7653973c74daa174b11e"
	// bStr := "fffe383c3cd8ef6c9dae70fc5bcf58c321687454aaf0505acfd648e610b50649"
	// bStr := "00b1a6dea4fb248782b6ea3672a891344788fe5222b9a0dfbebcd7e7700976cd80"
	byteArray, err := hex.DecodeString(bStr)
	if err != nil {
		fmt.Printf("error parse hex %s\n", err.Error())
		t.FailNow()
	}
	bint, err := UnmarshalBigInt(byteArray)
	if err != nil {
		fmt.Printf("error UnmarshalBigInt %s\n", err.Error())
		t.FailNow()
	}
	fmt.Printf("hex: %x\n", bint)
	fmt.Printf("val: %d\n", bint)

	// bbStr := "b1a6dea4fb248782b6ea3672a891344788fe5222b9a0dfbebcd7e7700976cd80"
	// bbStr := "fffe383c3cd8ef6c9dae70fc5bcf58c321687454aaf0505acfd648e610b50649"
	bbStr := "075bea74ced3d3b42cc9184efbe14c13e0307d4b1fbe7653973c74daa174b11e"
	bbyteArray, err := hex.DecodeString(bbStr)
	if err != nil {
		fmt.Printf("error parse hex %s\n", err.Error())
		t.FailNow()
	}
	b := new(big.Int).SetBytes(bbyteArray)
	fmt.Printf("b hex: %x\n", b)
	fmt.Printf("b val: %d\n", b)

	bbb, err := MarshalBigInt(bint)
	if err != nil {
		fmt.Printf("error MarshalBigInt %s\n", err.Error())
		t.FailNow()
	}
	fmt.Printf("MarshalBigInt %x\n", bbb)
}

func TestAddASN1BigInt(t *testing.T) {
	t.SkipNow()
	x := big.NewInt(-1)
	var b cryptobyte.Builder
	b.AddASN1BigInt(x)
	got, err := b.Bytes()
	if err != nil {
		t.Fatalf("unexpected error adding -1: %v", err)
	}
	fmt.Printf("out: %x \n", got)

	s := cryptobyte.String(got)
	var y big.Int
	ok := s.ReadASN1Integer(&y)
	if !ok || x.Cmp(&y) != 0 {
		t.Errorf("unexpected bytes %v, want %v", &y, x)
	}
	fmt.Printf("out: %x \n", &y)

}

func TestReadASN1IntegerInvalid(t *testing.T) {
	t.SkipNow()
	testData := []cryptobyte.String{
		[]byte{3, 1, 0}, // invalid tag
		// truncated
		[]byte{2, 1},
		[]byte{2, 2, 0},
		// not minimally encoded
		[]byte{2, 2, 0, 1},
		[]byte{2, 2, 0xff, 0xff},
	}

	for i, test := range testData {
		var out int64
		if test.ReadASN1Integer(&out) {
			t.Errorf("#%d: in.ReadASN1Integer() = true, want false (out = %d)", i, out)
		}
		fmt.Printf("out: %x \n", out)
	}

	for i, test := range testData {
		var out big.Int
		if test.ReadASN1Integer(&out) {
			t.Errorf("#%d: in.ReadASN1Integer() = true, want false (out = %d)", i, &out)
		}
		fmt.Printf("out: %x \n", &out)
	}
}
