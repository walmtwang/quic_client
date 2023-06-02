package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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
	flag.IntVar(&port, "port", 1935, "port, default 1935")
	flag.Parse()
	if ip == "" || tcUrl == "" || streamName == "" || fileName == "" {
		log.Fatalln("ip == \"\" ||tcUrl == \"\" ||streamName == \"\" ||fileName == \"\"")

	}

	tcUrl = strings.Replace(tcUrl, "rtmps://", "rtmp://", -1)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalf("tls.Dial failed, err:%v", err)
	}

	rtmpPublisher := rtmp.NewRtmpPublisher(conn, fileName,
		tcUrl,
		streamName)
	if err := rtmpPublisher.Start(); err != nil {
		log.Fatalf("rtmpPublisher.Start err:%v", err)
		return
	}
}
