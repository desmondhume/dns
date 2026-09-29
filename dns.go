package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Header struct {
	Id     uint16
	Qr     bool
	Opcode uint16 // 4bit
	AA     bool
	TC     bool
	RD     bool
	RA     bool
	Z      uint16 // 3 bit
	Rcode  uint16 // 4 bit

	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type Message struct {
	Header      Header
	Questions   []Question
	Answers     []ResourceRecord
	Authorities []ResourceRecord
	Additional  []ResourceRecord
}

type Question struct {
	Name  []byte
	Type  uint16 // 2 octets
	Class uint16 // 2 octets
}

type ResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

func readUint16(b []byte) (uint16, []byte) {
	return binary.BigEndian.Uint16(b), b[2:]
}

func (m *Message) unpackMsg(msg []byte) error {
	id, rest := readUint16(msg)
	bits, rest := readUint16(rest)
	qdCount, rest := readUint16(rest)
	anCount, rest := readUint16(rest)
	nsCount, rest := readUint16(rest)
	arCount, rest := readUint16(rest)

	m.Header.Id = id
	m.Header.Qr = (bits>>15)&1 != 0
	m.Header.Opcode = bits >> 11 & 0x0F
	m.Header.AA = (bits>>10)&1 != 0
	m.Header.TC = (bits>>9)&1 != 0
	m.Header.RD = (bits>>8)&1 != 0
	m.Header.RA = (bits>>7)&1 != 0
	m.Header.Z = (bits >> 4) & 0x07
	m.Header.Rcode = bits & 0x0F
	m.Header.QDCount = qdCount
	m.Header.ANCount = anCount
	m.Header.NSCount = nsCount
	m.Header.ARCount = arCount

	// TODO: Parse rest of the message

	return nil
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Wrong address")
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Connection could not be established")
		return
	}

	for {
		// Messages carried by UDP are restricted to 512 bytes (not counting the IP or UDP headers).
		// Longer messages are truncated and the TC bit is set in the header.
		// https://www.rfc-editor.org/info/rfc1035/#section-4.2.1
		buffer := make([]byte, 512)

		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error")
		}

		message := Message{}
		message.unpackMsg(buffer)

		fmt.Printf("%x", buffer)
		fmt.Println("message")
		fmt.Printf("%#v", message.Header)
	}
}
