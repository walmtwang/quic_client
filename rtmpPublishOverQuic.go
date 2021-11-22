package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"log"
	"net"
	"net/url"
	"quic_demo/rtmp"
	"strings"
	"time"
)

type QuicConn struct {
	QuicSession quic.Session
	QuicStream  quic.Stream
}

func NewQuicConn(quicSession quic.Session, quicStream quic.Stream) *QuicConn {
	return &QuicConn{
		QuicSession: quicSession,
		QuicStream:  quicStream,
	}
}

func (q *QuicConn) Read(b []byte) (n int, err error) {
	return q.QuicStream.Read(b)
}

func (q *QuicConn) Write(b []byte) (n int, err error) {
	return q.QuicStream.Write(b)
}

func (q *QuicConn) Close() error {
	return q.QuicStream.Close()
}

func (q *QuicConn) LocalAddr() net.Addr {
	return q.QuicSession.LocalAddr()
}

func (q *QuicConn) RemoteAddr() net.Addr {
	return q.QuicSession.RemoteAddr()
}

func (q *QuicConn) SetDeadline(t time.Time) error {
	return q.QuicStream.SetDeadline(t)
}

func (q *QuicConn) SetReadDeadline(t time.Time) error {
	return q.QuicStream.SetReadDeadline(t)
}

func (q *QuicConn) SetWriteDeadline(t time.Time) error {
	return q.QuicStream.SetWriteDeadline(t)
}

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

	quicSession, err := quic.DialAddr(fmt.Sprintf("%s:%d", ip, port), &tls.Config{
		ServerName: domain,
		NextProtos: []string{"rtmp over quic"},
	}, &quic.Config{
		Versions: []quic.VersionNumber{quic.VersionDraft29},
	})
	if err != nil {
		log.Fatalf("quic.DialAddr err:%v", err)
		return
	}
	quicStream, err := quicSession.OpenStreamSync(context.Background())

	quicConn := NewQuicConn(quicSession, quicStream)

	rtmpPublisher := rtmp.NewRtmpPublisher(quicConn, fileName,
		tcUrl,
		streamName)
	if err := rtmpPublisher.Start(); err != nil {
		log.Fatalf("rtmpPublisher.Start err:%v", err)
		return
	}
}
