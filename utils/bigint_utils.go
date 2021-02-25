package utils

import (
	"crypto/elliptic"
	"errors"
	"math/big"

	"encoding/asn1"
	// "github.com/mathews/asn1"
	"github.com/mathews/gmsm/log"
)

//MarshalBigInt correctly encode a negative bigint,
//if the bigint is negative, don't use big.Int.Bytes() to encode the bigint
func MarshalBigInt(num *big.Int) ([]byte, error) {

	log.Logger.Debugf("marshalling BigInt %x", num)

	btxt, err := asn1.Marshal(num)
	if err != nil {
		return nil, err
	} else {
		log.Logger.Debugf("buf lenth %d, content %x", len(btxt), btxt)
		ret := btxt[2:]

		err = checkInteger(ret)
		if err != nil {
			log.Logger.Errorf("asn1 Marshal %x as %x", num, ret)
			ret = num.Bytes()
		}
		return ret, nil

	}
}

var bigOne = big.NewInt(1)

// checkInteger returns nil if the given bytes are a valid DER-encoded
// INTEGER and an error otherwise.
func checkInteger(bytes []byte) error {
	if len(bytes) == 0 {
		return errors.New("empty integer")
	}
	if len(bytes) == 1 {
		return nil
	}
	//FIXME
	if (bytes[0] == 0 && bytes[1]&0x80 == 0) || (bytes[0] == 0xff && bytes[1]&0x80 == 0x80) {
		log.Logger.Errorf("integer %x not minimally-encoded", bytes)
		return errors.New("integer not minimally-encoded")
	}
	return nil
}

// UnmarshalBigInt correctly decode a bigint, even a negative bigint
func UnmarshalBigInt(bytes []byte) (*big.Int, error) {
	log.Logger.Debugf("UnmarshalBigInt %x", bytes)
	if err := checkInteger(bytes); err != nil {
		ret := new(big.Int)
		ret.SetBytes(bytes)
		return ret, nil
	}
	ret := new(big.Int)
	if len(bytes) > 0 {
		// if len(bytes) > 32 && bytes[0] == 1 {
		// 	ret.SetBytes(bytes[1:])
		// 	ret = ret.Neg(ret)

		// } else
		if bytes[0]&0x80 == 0x80 {
			// This is a negative number.
			notBytes := make([]byte, len(bytes))
			for i := range notBytes {
				notBytes[i] = ^bytes[i]
			}
			ret.SetBytes(notBytes)
			ret.Add(ret, bigOne)
			ret = ret.Neg(ret)
			log.Logger.Debugf("Unmarshaled negative BigInt as %d\n", ret)
		} else {
			ret.SetBytes(bytes)
		}
	}
	log.Logger.Debugf("Unmarshaled BigInt as %d", ret)
	return ret, nil
}

// EllipticMarshal is revised from elliptic.Marshal to correct the negative bigint problem
func EllipticMarshal(curve elliptic.Curve, x, y *big.Int) ([]byte, error) {
	byteLen := (curve.Params().BitSize + 7) >> 3

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point

	xBytes, err := MarshalBigInt(x)
	if err != nil {
		return nil, err
	}
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	yBytes, err := MarshalBigInt(y)
	if err != nil {
		return nil, err
	}
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)

	log.Logger.Debugf("Marshal x,y points: %x", ret)

	return ret, nil
}

// EllipticUnmarshal is revised from elliptic.Unmarshal to correct the negative bigint problem
func EllipticUnmarshal(curve elliptic.Curve, data []byte) (x, y *big.Int, err error) {
	byteLen := (curve.Params().BitSize + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}
	if data[0] != 4 { // uncompressed form
		return
	}
	p := curve.Params().P
	x, err = UnmarshalBigInt(data[1 : 1+byteLen])
	if err != nil {
		return nil, nil, err
	}

	y, err = UnmarshalBigInt(data[1+byteLen:])
	if err != nil {
		return nil, nil, err
	}
	if x.Cmp(p) >= 0 || y.Cmp(p) >= 0 {
		return nil, nil, errors.New("x,y is bigger than P point")
	}
	if !curve.IsOnCurve(x, y) {
		return nil, nil, errors.New("x,y points not on curve")
	}
	return
}
