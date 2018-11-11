package snippets

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

func encryptListen() error {
	var certPem, keyPem []byte
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return err
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := tls.Listen("tcp", ":2000", cfg)
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}

func newService() {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache("cache-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.come", "www.example.come"),
	}
	router := httprouter.New()
	var userHandler = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var _ = "process"
	}
	var withMiddleware = func(h httprouter.Handle) httprouter.Handle {
		var _ = "intercept"
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			h(w, r, params)
		}
	}
	router.GET("/user/:id", withMiddleware(userHandler))
	server := &http.Server{
		Addr:         ":https",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    manager.TLSConfig(),
		Handler:      router,
	}
	server.ListenAndServeTLS("", "")
}

func withTrace(req *http.Request) *http.Request {
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
	}
	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
}
