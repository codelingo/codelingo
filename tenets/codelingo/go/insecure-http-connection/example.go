package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

type zeroSource struct {}

func (zeroSource) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = 0
	}

	return len(b), nil
}

func main() {

	serverOne := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	serverOne.TLS = &tls.Config{
		Rand: zeroSource{},
	}

	serverOne.StartTLS()
	defer serverOne.Close()

	serverTwo := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	serverTwo.TLS = &tls.Config{ // Issue
		Rand: zeroSource{},
		InsecureSkipVerify: true,
	}

	serverTwo.StartTLS()
	defer serverTwo.Close()



	w := os.Stdout

	clientOne := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				KeyLogWriter: w,
				Rand: zeroSource{},
			},
		},
	}
	resp, err := clientOne.Get(serverOne.URL)
	if err != nil {
		log.Fatalf("Faile to get url: %v", err)
	}

	resp.Body.Close()

	clientTwo := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ // Issue
				KeyLogWriter: w,
				Rand: zeroSource{},
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err = clientTwo.Get(serverTwo.URL)
	if err != nil {
		log.Fatalf("Faile to get url: %v", err)
	}

	resp.Body.Close()


}
