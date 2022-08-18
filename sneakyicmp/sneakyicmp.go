package sneakyicmp

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/go-ping/ping"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ipv4HeaderLength = 20
	icmpHeader       = 8
)

var body = &icmp.Echo{
	ID:   123123,
	Seq:  100,
	Data: []byte("213"),
}

var xd = icmp.Message{
	Type: ipv4.ICMPTypeEcho,
	Code: 0,
	Body: body,
}

func formatMessage(message string) []byte {
	x := []byte(message)
	return x
}

/*
Based on ip in string, return a pointer to object net.IPAddr
*/
func ipStringToIpAddr(ip string) *net.IPAddr {
	ipOctects := strings.Split(ip, ".")
	var ipByteOctets [4]byte
	for i := 0; i < len(ipOctects); i++ {
		ipInt, err := strconv.Atoi(ipOctects[i])
		if err != nil {
			log.Fatalf("Wrong address!")
		}
		ipByteOctets[i] = byte(ipInt)
	}
	ipv4 := net.IPv4(ipByteOctets[0], ipByteOctets[1], ipByteOctets[2], ipByteOctets[3])
	targetIp := &net.IPAddr{
		IP: ipv4,
	}
	return targetIp
}

func SendICMP(msg icmp.Message, dst string) {
	c, err := net.ListenPacket("ip:1", dst)
	if err != nil {
		log.Fatal(err)
	}
	pkt, err := msg.Marshal(nil)
	if err != nil {
		log.Fatal("Failed to fetch message")
	}

	targetIp := ipStringToIpAddr(dst)
	x, err := c.WriteTo(pkt, targetIp)
	log.Println(err)
	log.Println(x)
}

func SendSneakyMessage(addr string, message string) error {
	p := ping.New(addr)

	fMessage := formatMessage(message)
	for i := 0; i < len(message); i++ {
		p.Count = 1
		p.Size = int(fMessage[i])
		p.Run()
	}

	//Send size 256 to finish transmission
	p.Count = 1
	p.Size = 256
	p.Run()

	return nil
}

func RecvSneakyMessage(url string) []byte {
	var sneakyMessage []byte
	conn, err := icmp.ListenPacket("ip4:icmp", url)
	if err != nil {
		log.Fatal(err)
	}

	bb := make([]byte, 256)
	cf := ipv4.FlagTTL | ipv4.FlagInterface
	msg := []ipv4.Message{
		{
			Buffers: [][]byte{bb},
			OOB:     ipv4.NewControlMessage(cf),
		},
	}

	pkt := conn.IPv4PacketConn()

	for {
		nrOfPkt, err := pkt.ReadBatch(msg, 0)
		for i := 0; i < nrOfPkt; i++ {
			if (msg[0].N) < 256 {
				sneakyMessage = append(sneakyMessage, byte((msg[0].N - ipv4HeaderLength - icmpHeader)))
			} else {
				return sneakyMessage
			}
		}
		if err != nil {
			log.Printf("%s, error could affect the message", err)
			continue
		}
	}

}
