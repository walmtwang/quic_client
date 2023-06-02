package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	quic "github.com/walmtwang/quic-go"
	"github.com/walmtwang/quic-go/h2quic"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	var ip string
	var port int
	var httpUrl string
	var requestTime int
	var sleepTime int
	var waitTime int
	flag.StringVar(&ip, "ip", "", "server ip")
	flag.StringVar(&httpUrl, "url", "", "http url, https://domain/live/stream.flv")
	flag.IntVar(&port, "port", 443, "server port, default 443")
	flag.IntVar(&requestTime, "requestTime", 10, "request time(s), default 10s")
	flag.IntVar(&sleepTime, "sleepTime", 60, "sleep time(s), default 60s")
	flag.IntVar(&waitTime, "waitTime", 60, "wait time(s), default 60s")
	flag.Parse()
	if ip == "" || httpUrl == "" {
		log.Fatalln("ip == \"\" ||  url == \"\"")
	}

	url2, err := url.Parse(httpUrl)
	if err != nil {
		log.Fatalf("url.Parse failed, err:%v", err)
	}
	domain := strings.Split(url2.Host, ":")[0]

	roundTripper := &h2quic.RoundTripper{
		QuicConfig: &quic.Config{
			Versions: []quic.VersionNumber{
				quic.VersionGQUIC43,
			},
		},
		TLSClientConfig: &tls.Config{
			ServerName: domain,
			NextProtos: []string{"rtmp over quic"},
		},
		Dial: func(network, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.Session, error) {
			return quic.DialAddr(fmt.Sprintf("%s:%d", ip, port), tlsCfg, cfg)
		},
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}
	go func() {

		resp, err := hclient.Get(httpUrl)
		if err != nil {
			log.Fatalf("hclient.Get err:%v", err)
		}
		tlsInfo := resp.TLS
		for i := 0; i < len(tlsInfo.PeerCertificates); i++ {
			cert := tlsInfo.PeerCertificates[i]
			log.Printf("index:%v, Issuer:%v, Subject:%v, NotBefore:%v, NotAfter:%v", i, cert.Issuer, cert.Subject, cert.NotBefore, cert.NotAfter)
		}

		fmt.Printf("http status:%v\n", resp.StatusCode)
		time.Sleep(time.Duration(requestTime) * time.Second)
		resp.Body.Close()
	}()

	time.Sleep(time.Duration(sleepTime) * time.Second)

	resp, err := hclient.Get(httpUrl)
	if err != nil {
		log.Fatalf("hclient.Get err:%v", err)
	}

	fmt.Printf("http status:%v\n", resp.StatusCode)
	resp.Body.Close()
	time.Sleep(time.Duration(waitTime) * time.Second)
}
