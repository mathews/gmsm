package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/mathews/gmsm/gmtls"
	"github.com/mathews/gmsm/x509"
)

// func TestTaSSLServer(t *testing.T) {
func main() {
	caArray, err := ioutil.ReadFile("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/CA.pem")
	if err != nil {
		panic("failed to read root certificate")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caArray)
	if !ok {
		panic("failed to parse root certificate")
	}

	ssCert, err := gmtls.LoadGMX509KeyPair("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SS.pem", "/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SS.key")
	if err != nil {
		panic("failed to load certificate")
	}
	seCert, err := gmtls.LoadGMX509KeyPair("/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SE.pem", "/home/mathews/dev/sm-workspace/TASSL-1.1.1b/tassl_demo/cert/certs/SE.key")
	if err != nil {
		panic("failed to load certificate")
	}
	config := &gmtls.Config{
		GMSupport: &gmtls.GMSupport{},
		// CipherSuites: []uint16{gmtls.GMTLS_SM2_WITH_SM1_SM3, gmtls.GMTLS_ECDHE_SM2_WITH_SM4_SM3},
		Rand: rand.Reader,

		RootCAs: roots,
		// InsecureSkipVerify: true,
		Certificates: []gmtls.Certificate{ssCert, seCert},
	}
	var IDs []string
	for _, i := range config.CipherSuites {

		IDs = append(IDs, strconv.Itoa(int(i)))
	}

	fmt.Printf("ciphers: %s\n", strings.Join(IDs, ", "))

	// var wg sync.WaitGroup

	listerner, err := gmtls.Listen("tcp", "192.168.11.230:4433", config)

	if err != nil {
		panic("failed to start Listener: " + err.Error())
	}
	defer listerner.Close()

	// wg.Add(1)
	// wg.Wait()

	// conn, err := listerner.Accept()
	// if err != nil {
	// 	log.Println(err)
	// }
	// go handleConnection(conn)

	//Create the default mux
	mux := http.NewServeMux()

	//Handling the /v1/teachers. The handler is a function here
	mux.HandleFunc("/", teacherHandler)
	//Create the server.
	s := &http.Server{
		// Addr:    "192.168.11.230:4433",
		Handler: mux,
	}
	err = s.Serve(listerner)
	if err != nil {
		panic("failed to start Server: " + err.Error())
	}

	defer s.Close()

	// http.HandleFunc("/", helloServer)
	// err = http.Serve(listerner, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	// t.Cleanup(func() {
	// 	wg.Done()
	// })

}

func teacherHandler(res http.ResponseWriter, req *http.Request) {
	data := []byte("V1 of teacher's called")
	res.WriteHeader(200)
	res.Header().Set("Content-Type", "text/plain")
	res.Write(data)
}

func helloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// r := bufio.NewReader(conn)
	for {
		// msg, err := r.ReadString('\n')
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }

		// println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
