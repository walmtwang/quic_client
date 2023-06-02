package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"log"
	"net/url"
	"quic_demo/quicConn"
	"quic_demo/rtmp"
	"strings"
)

func main() {

	var ip string
	var tcUrl string
	var streamName string
	var fileName string
	var port int
	flag.StringVar(&ip, "ip", "", "ip")
	flag.StringVar(&tcUrl, "tcUrl", "", "tcUrl")
	flag.StringVar(&streamName, "streamName", "", "streamName")
	flag.StringVar(&fileName, "fileName", "", "fileName")
	flag.IntVar(&port, "port", 443, "port, default 443")
	flag.Parse()
	if ip == "" || tcUrl == "" || streamName == "" || fileName == "" {
		log.Fatalln("ip == \"\" ||tcUrl == \"\" ||streamName == \"\" ||fileName == \"\"")

	}

	url2, err := url.Parse(tcUrl)
	if err != nil {
		log.Fatalf("url.Parse failed, err:%v", err)
	}
	domain := strings.Split(url2.Host, ":")[0]

	tlsConfig := tls.Config{
		ServerName: domain,
		NextProtos: []string{"rtmp over quic"},
	}
	quicSession, err := quic.DialAddr(fmt.Sprintf("%s:%d", ip, port), &tlsConfig, &quic.Config{
		Versions: []quic.VersionNumber{quic.VersionDraft29},
	})
	if err != nil {
		log.Fatalf("quic.DialAddr err:%v", err)
		return
	}
	log.Printf("tlsConfig:%v", tlsConfig)
	log.Printf("quicSession:%v", quicSession)
	log.Printf("quicSession.ConnectionState().TLS:%v", quicSession.ConnectionState().TLS)
	tlsInfo := quicSession.ConnectionState().TLS
	for i := 0; i < len(tlsInfo.PeerCertificates); i++ {
		cert := tlsInfo.PeerCertificates[i]
		log.Printf("index:%v, Issuer:%v, Subject:%v, NotBefore:%v, NotAfter:%v", i, cert.Issuer, cert.Subject, cert.NotBefore, cert.NotAfter)
	}

	quicStream, err := quicSession.OpenStreamSync(context.Background())

	qConn := quicConn.NewQuicConn(quicSession, quicStream)

	rtmpPlay := rtmp.NewRtmpPlay(qConn, fileName,
		tcUrl,
		streamName)
	if err := rtmpPlay.Start(); err != nil {
		log.Fatalf("rtmpPlay.Start err:%v", err)
		return
	}
}
