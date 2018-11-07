package snippets

import (
	"crypto/tls"
	"io"
	"net"
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
