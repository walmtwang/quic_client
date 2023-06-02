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
	var ignoreSSL int
	flag.StringVar(&ip, "ip", "", "ip")
	flag.StringVar(&tcUrl, "tcUrl", "", "tcUrl")
	flag.StringVar(&streamName, "streamName", "", "streamName")
	flag.StringVar(&fileName, "fileName", "", "fileName")
	flag.IntVar(&port, "port", 443, "port, default 443")
	flag.IntVar(&ignoreSSL, "ignoreSSL", 0, "0: open SSL check, 1: ignore SSL check. default 0")
	flag.Parse()
	if ip == "" || tcUrl == "" || streamName == "" || fileName == "" {
		log.Fatalln("ip == \"\" ||tcUrl == \"\" ||streamName == \"\" ||fileName == \"\"")

	}
	insecureSkipVerify := false
	if ignoreSSL == 1 {
		insecureSkipVerify = true
	}

	url2, err := url.Parse(tcUrl)
	if err != nil {
		log.Fatalf("url.Parse failed, err:%v", err)
	}
	domain := strings.Split(url2.Host, ":")[0]

	quiConn, err := quic.DialAddr(fmt.Sprintf("%s:%d", ip, port), &tls.Config{
		ServerName:         domain,
		NextProtos:         []string{"rtmp over quic"},
		InsecureSkipVerify: insecureSkipVerify,
	}, &quic.Config{
		//Versions: []quic.VersionNumber{quic.VersionDraft29},
	})
	if err != nil {
		log.Fatalf("quic.DialAddr err:%v", err)
		return
	}
	tlsInfo := quiConn.ConnectionState().TLS
	for i := 0; i < len(tlsInfo.PeerCertificates); i++ {
		cert := tlsInfo.PeerCertificates[i]
		log.Printf("index:%v, Issuer:%v, Subject:%v, NotBefore:%v, NotAfter:%v", i, cert.Issuer, cert.Subject, cert.NotBefore, cert.NotAfter)
	}

	quicStream, err := quiConn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalf("quiConn.OpenStreamSync err:%v", err)
		return
	}

	qConn := quicConn.NewQuicConn(quiConn, quicStream)

	rtmpPublisher := rtmp.NewRtmpPublisher(qConn, fileName,
		tcUrl,
		streamName)
	if err := rtmpPublisher.Start(); err != nil {
		log.Fatalf("rtmpPublisher.Start err:%v", err)
		return
	}
}
