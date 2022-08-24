package sneakyicmp

import (
	"log"
	"net"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ipv4HeaderLength = 20
	icmpHeader       = 8
)

func formatMessage(message string) []byte {
	x := []byte(message)
	return x
}

type Ping struct {
	Privileged bool
	Size       int
	DstIp      string
	SrcIp      string
	Protocol   string
}

type Message struct {
	payload       []Ping
	encryption    bool
	encryptionKey string
}

/*
Based on ip in string, return a pointer to object net.IPAddr
*/
func ipStringToIpAddr(ip string) *net.IPAddr {

	ipv4 := net.ParseIP(ip)
	targetIp := &net.IPAddr{
		IP: ipv4,
	}
	return targetIp
}

func NewPing(dstIp string, srcIp string, elevated bool) *Ping {
	p := Ping{
		Privileged: elevated,
		DstIp:      dstIp,
		SrcIp:      srcIp,
	}

	return p
}

func (p *Ping) privileged(elevated bool) {
	if elevated == true {
		p.Privileged = true
		p.Protocol = "ip:icmp"
	} else {
		p.Privileged = false
		p.Protocol = "udp4"
	}
}

func (p *Ping) makeConn() (net.PacketConn, error) {
	conn, err := net.ListenPacket(p.Protocol, p.DstIp)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *Ping) createICMP() ([]byte, error) {

	body := &icmp.Echo{
		ID:   123123,
		Seq:  100,
		Data: []byte("123213123"),
	}

	icmpPacket := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: body,
	}

	pkt, err := icmpPacket.Marshal(nil)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

func (p *Ping) SendICMP() {
	//Create icmp packet based on ping struct
	icmpPacket, err := p.createICMP()
	if err != nil {
		log.Fatalf("Failed to fetch message, %s", err)
	}
	//Make connection with host
	c, err := p.makeConn()
	if err != nil {
		log.Fatal(err)
	}

	//Send packet to host
	targetIp := ipStringToIpAddr(p.DstIp)
	_, err = c.WriteTo(icmpPacket, targetIp)
	if err != nil {
		log.Print(err)
	}
}

/*
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
*/
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
