package sneakyicmp

import (
	"log"

	"github.com/go-ping/ping"
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
