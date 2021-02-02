/*
Copyright Suzhou Tongji Fintech Research Institute 2017 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package x509

import (
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"testing"
	"time"

	"encoding/asn1"

	"github.com/mathews/gmsm/sm2"
)

func TestX509(t *testing.T) {
	priv, err := sm2.GenerateKey(nil) // 生成密钥对
	if err != nil {
		t.Fatal(err)
	}
	privPem, err := WritePrivateKeyToPem(priv, nil) // 生成密钥文件
	if err != nil {
		t.Fatal(err)
	}
	pubKey, _ := priv.Public().(*sm2.PublicKey)
	pubkeyPem, err := WritePublicKeyToPem(pubKey) // 生成公钥文件

	// err = ioutil.WriteFile("/home/mathews/ssm2.pem", pubkeyPem, 0777)

	if err != nil {
		t.Fatal(err)
	}
	privKey, err := ReadPrivateKeyFromPem(privPem, nil) // 读取密钥
	if err != nil {
		t.Fatal(err)
	}
	pubKey, err = ReadPublicKeyFromPem(pubkeyPem) // 读取公钥
	if err != nil {
		t.Fatal(err)
	}
	templateReq := CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "test.qloudchain.com",
			Organization: []string{"QloudChain"},
		},
		//		SignatureAlgorithm: ECDSAWithSHA256,
		SignatureAlgorithm: SM2WithSM3,
	}
	reqPem, err := CreateCertificateRequestToPem(&templateReq, privKey)
	if err != nil {
		t.Fatal(err)
	}
	req, err := ReadCertificateRequestFromPem(reqPem)
	if err != nil {
		t.Fatal(err)
	}
	err = req.CheckSignature()
	if err != nil {
		t.Fatalf("Request CheckSignature error:%v", err)
	} else {
		fmt.Printf("CheckSignature ok\n")
	}
	testExtKeyUsage := []ExtKeyUsage{ExtKeyUsageClientAuth, ExtKeyUsageServerAuth}
	testUnknownExtKeyUsage := []asn1.ObjectIdentifier{[]int{1, 2, 3}, []int{2, 59, 1}}
	extraExtensionData := []byte("extra extension")
	commonName := "test.qloudchain.com"
	template := Certificate{
		// SerialNumber is negative to ensure that negative
		// values are parsed. This is due to the prevalence of
		// buggy code that produces certificates with negative
		// serial numbers.
		SerialNumber: big.NewInt(-1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"QloudChain"},
			Country:      []string{"China"},
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  []int{2, 5, 4, 42},
					Value: "test",
				},
				// This should override the Country, above.
				{
					Type:  []int{2, 5, 4, 6},
					Value: "CN",
				},
			},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Date(2022, time.December, 10, 12, 1, 1, 1, time.UTC),

		//		SignatureAlgorithm: ECDSAWithSHA256,
		SignatureAlgorithm: SM2WithSM3,

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     KeyUsageCertSign,

		ExtKeyUsage:        testExtKeyUsage,
		UnknownExtKeyUsage: testUnknownExtKeyUsage,

		BasicConstraintsValid: true,
		IsCA:                  true,

		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},

		DNSNames:       []string{"test.qloudchain.com"},
		EmailAddresses: []string{"admin@qloudchain.com"},
		IPAddresses:    []net.IP{net.IPv4(127, 0, 0, 1).To4(), net.ParseIP("::1")},

		PolicyIdentifiers:   []asn1.ObjectIdentifier{[]int{1, 2, 3}},
		PermittedDNSDomains: []string{".qloudchain.com", "qloudchain.com"},

		CRLDistributionPoints: []string{"http://crl1.example.com/ca1.crl", "http://crl2.example.com/ca1.crl"},

		ExtraExtensions: []pkix.Extension{
			{
				Id:    []int{1, 2, 3, 4},
				Value: extraExtensionData,
			},
			// This extension should override the SubjectKeyId, above.
			{
				Id:       oidExtensionSubjectKeyId,
				Critical: false,
				Value:    []byte{0x04, 0x04, 4, 3, 2, 1},
			},
		},
	}
	caArray, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/CA.pem")
	caCert, err := ReadCertificateFromPem(caArray)
	if err != nil {
		println("error: %s \n", err.Error())
	}
	pubKey, _ = priv.Public().(*sm2.PublicKey)
	certpem, err := CreateCertificateToPem(&template, caCert, pubKey, privKey)
	if err != nil {
		t.Fatal("failed to create cert file")
	}

	err = ioutil.WriteFile("/home/mathews/sm2cert.pem", certpem, 0655)
	if err != nil {
		t.Fatal("failed to write cert file")
	}
	cert, err := ReadCertificateFromPem(certpem)
	if err != nil {
		t.Fatal("failed to read cert file")
	}
	err = cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("CheckSignature ok\n")
	}
}
