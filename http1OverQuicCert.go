package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"log"
)

func main() {

	var ip string
	var port int
	var domain string
	flag.StringVar(&domain, "domain", "", "server domain, default ''")
	flag.StringVar(&ip, "ip", "", "server ip")
	flag.IntVar(&port, "port", 80, "server port, default 80")
	flag.Parse()
	if ip == "" || domain == "" {
		log.Fatalln("ip == \"\" || domain == \"\"")
	}

	tlsConf := &tls.Config{
		ServerName: domain,
		NextProtos: []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddrEarlyContext(context.Background(), fmt.Sprintf("%s:%d", ip, port), tlsConf, nil)
	if err != nil {
		log.Fatalf("quic.DialAddrEarlyContext failed, err:%v", err)
	}
	tlsInfo := conn.ConnectionState().TLS
	for i := 0; i < len(tlsInfo.PeerCertificates); i++ {
		cert := tlsInfo.PeerCertificates[i]
		log.Printf("index:%v, Issuer:%v, Subject:%v, NotBefore:%v, NotAfter:%v", i, cert.Issuer, cert.Subject, cert.NotBefore, cert.NotAfter)
	}

}
