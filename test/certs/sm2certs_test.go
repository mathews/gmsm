package certs

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mathews/gmsm/x509"
)

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func TestSm2Cert(t *testing.T) {

	/*

			certPEMBlock, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SE.crt")
			// certPEMBlock, err := ioutil.ReadFile("/home/mathews/桌面/sm/CFCA_CS_CA.cer")
		    if err != nil {
				println("error!")
		    }
		    //获取下一个pem格式证书数据 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
			// certDERBlock, restPEMBlock := pem.Decode(certPEMBlock)
			certDERBlock, _ := pem.Decode(certPEMBlock)
		    if certDERBlock == nil {
		        println("error!")
		    }
			pub, err := sm2.ParsePKIXPublicKey(certPEMBlock)
			if err!=nil{
				fmt.Errorf("error parsing public key %s", err)
			}
			if pub ==nil{
				fmt.Errorf("publickey is nil!")
			}
			fmt.Println(reflect.TypeOf(pub))
			println(typeof(pub))
	*/
	pubArr, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SE.pem")
	if err != nil {
		println("error")
	}
	cert, err := x509.ReadCertificateFromPem(pubArr)
	if err != nil {
		println("error")
	}
	iName := cert.Issuer.Names
	fmt.Printf("%s\n", iName)
	caArray, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/CA.pem")
	if err != nil {
		panic("failed to read root certificate")
	}
	cacert, err := x509.ReadCertificateFromPem(caArray)
	if err != nil {
		println("error")
	}
	err = cert.CheckSignatureFrom(cacert)

	if err != nil {
		println("CheckSignatureFrom error")
	} else {
		println("CheckSignatureFrom ok!")
	}
	keyArr, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SE.key")
	key, err := x509.ReadPrivateKeyFromPem(keyArr, nil)
	if err != nil {
		println("error")
	}
	fmt.Printf("x=%s\n", key.X)

}
