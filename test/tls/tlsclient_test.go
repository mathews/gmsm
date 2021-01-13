package tls

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/tjfoc/gmsm/gmtls"
	"github.com/tjfoc/gmsm/x509"
)

func TestTaSSLClient(t *testing.T) {
	// t.SkipNow()

	caArray, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/CA.pem")
	if err != nil {
		panic("failed to read root certificate")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caArray)
	if !ok {
		panic("failed to parse root certificate")
	}
	config := &gmtls.Config{
		GMSupport: &gmtls.GMSupport{},
		// CipherSuites: []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3, gmtls.GMTLS_ECDHE_SM2_WITH_SM4_SM3},
		CipherSuites: []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3},
		Rand:         rand.Reader,

		// RootCAs:            roots,
		InsecureSkipVerify: true,
	}
	var IDs []string
	for _, i := range config.CipherSuites {

		IDs = append(IDs, strconv.Itoa(int(i)))
	}

	fmt.Printf("ciphers: %s\n", strings.Join(IDs, ", "))
	conn, err := gmtls.Dial("tcp", "192.168.11.60:445", config)
	defer conn.Close()

	log.Printf("remote %s\n", conn.RemoteAddr().String())
	log.Printf("local %s\n", conn.LocalAddr().String())

	if err != nil {
		panic("failed to connect: " + err.Error())
	}

}

func TestHttpsClient(t *testing.T) {

	t.SkipNow()

	caArray, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/CA.pem")
	if err != nil {
		panic("failed to read root certificate")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caArray)
	if !ok {
		panic("failed to parse root certificate")
	}
	config := &gmtls.Config{
		GMSupport: &gmtls.GMSupport{},
		// CipherSuites: []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3, gmtls.GMTLS_ECDHE_SM2_WITH_SM4_SM3},
		// CipherSuites: []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3},
		Rand: rand.Reader,

		// RootCAs:            roots,
		InsecureSkipVerify: true,
	}
	var IDs []string
	for _, i := range config.CipherSuites {

		IDs = append(IDs, strconv.Itoa(int(i)))
	}

	fmt.Printf("ciphers: %s\n", strings.Join(IDs, ", "))

	// conn, err := gmtls.Dial("tcp", "192.168.11.60:445", config)
	// defer conn.Close()

	// log.Printf("remote %s\n", conn.RemoteAddr().String())
	// log.Printf("local %s\n", conn.LocalAddr().String())

	// if err != nil {
	// 	panic("failed to connect: " + err.Error())
	// }

	client := &http.Client{
		Transport: &http.Transport{
			// TLSClientConfig: &tls.Config{
			// 	InsecureSkipVerify: true,
			// 	CipherSuites:       []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3, gmtls.GMTLS_ECDHE_SM2_WITH_SM4_SM3},
			// 	Rand:               rand.Reader,
			// },
			// Dial: func(network, addr string) (net.Conn, error) {
			// 	return gmtls.Dial(network, addr, config)
			// },
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return gmtls.Dial(network, addr, config)
			},
		},
		// Timeout: time.Duration(6000) * time.Microsecond,
	}

	resp, err := client.Get("https://192.168.11.230:4433/")
	if err != nil {
		panic("failed to Get: " + err.Error())
	}
	fmt.Printf("resp status: %s\n", resp.Status)

	text := make([]byte, int(resp.ContentLength))

	resp.Body.Read(text)

	fmt.Printf("resp: %s\n", text)

}
