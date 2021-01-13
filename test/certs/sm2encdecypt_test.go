package certs

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"testing"

	sm509 "github.com/tjfoc/gmsm/x509"
)

func TestSm2EncDecrypt(t *testing.T) {

	keyArray, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SE.key")

	key, err := sm509.ReadPrivateKeyFromPem(keyArray, nil)

	if err != nil {
		println("error")
	}

	privKey := ecdsa.PrivateKey{ecdsa.PublicKey{key.PublicKey.Curve, key.PublicKey.X, key.PublicKey.Y}, key.D}
	// signedMsg, err := key.Sign(nil, []byte("fsdfdsgfdgdgh"), nil)

	// var randSign = "22220316zafes20180lk7zafes20180619zafepikas"
	// reader:=strings.NewReader(randSign)

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte("fsdfdsgfdgdgh"))
	// r, s, err := sm2.Sm2Sign(key, []byte("fsdfdsgfdgdgh"), []byte{}, nil)
	fmt.Printf("signature is %s, %s \n", r, s)

	// key, err := x509.ParseECPrivateKey(keyArray)
	// if err != nil {
	// 	println("error %s\n", err.Error())
	// }

	// key := keyIntf.(*ecdsa.PrivateKey)

	// r, s, err := ecdsa.Sign(nil, key, []byte("fsdfdsgfdgdgh"))

	if err != nil {
		println("error")
	}
	// fmt.Printf("signed byte[] is %s \n", signedMsg)

	pubArray, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SE.pem")

	// block, _ := pem.Decode(pubArray)
	// if block == nil || block.Type != "PUBLIC KEY" {
	// 	log.Fatal("PEM block has bot publickey")
	// }
	// println("key array: %s\n", block.Bytes)

	// cert, err := sm509.ParseSm2PublicKey(block.Bytes)
	// if err != nil {
	// 	println("error: %s\n", err.Error())
	// }
	// valid := cert.Verify([]byte("fsdfdsgfdgdgh"), signedMsg)

	// publicKey := cert.ToX509Certificate()

	// pub, err := sm509.ParseSm2PublicKey(pubArray)
	pub, err := sm509.ReadCertificateFromPem(pubArray)
	if err != nil {
		println("error: %s \n", err.Error())
	}
	println("DNS: %s\n", pub.DNSNames)
	println("subject: %s\n", pub.RawSubject)

	// valid := pub.Verify([]byte("fsdfdsgfdgdgh"), signedMsg)

	// sign := sm2Signature{}

	// _, err = asn1.Unmarshal(signedMsg, &sign)
	// if err != nil {
	// 	println("error")
	// }
	publicKey := pub.PublicKey.(*ecdsa.PublicKey)

	valid := ecdsa.Verify(publicKey, []byte("fsdfdsgfdgdgh"), r, s)

	// valid := pub.Verify([]byte("fsdfdsgfdgdgh"), signedMsg)

	// valid := sm2.Sm2Verify(publicKey, []byte("fsdfdsgfdgdgh"), []byte{}, r, s)

	if valid {
		println("valid")
	} else {
		println("fail")
	}

}
