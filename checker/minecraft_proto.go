package checker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func writeVarInt(buf *bytes.Buffer, value int32) {
	uval := uint32(value)
	for {
		temp := uval & 0x7F
		uval >>= 7
		if uval != 0 {
			temp |= 0x80
		}
		buf.WriteByte(byte(temp))
		if uval == 0 {
			break
		}
	}
}

func readVarInt(r io.Reader) (int32, error) {
	var result int32
	var numRead uint
	buf := make([]byte, 1)
	for {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		result |= int32(buf[0]&0x7F) << (7 * numRead)
		numRead++
		if numRead > 5 {
			return 0, fmt.Errorf("VarInt too big")
		}
		if buf[0]&0x80 == 0 {
			break
		}
	}
	return result, nil
}

func writeString(buf *bytes.Buffer, s string) {
	writeVarInt(buf, int32(len(s)))
	buf.WriteString(s)
}

func buildHandshakePacket(host string, port int) []byte {
	var data bytes.Buffer
	writeVarInt(&data, 0x00)
	writeVarInt(&data, 767)
	writeString(&data, host)
	binary.Write(&data, binary.BigEndian, uint16(port))
	writeVarInt(&data, 1)

	var packet bytes.Buffer
	writeVarInt(&packet, int32(data.Len()))
	packet.Write(data.Bytes())
	return packet.Bytes()
}

func buildStatusRequestPacket() []byte {
	var packet bytes.Buffer
	writeVarInt(&packet, 1)
	writeVarInt(&packet, 0x00)
	return packet.Bytes()
}

func buildPingPacket(payload int64) []byte {
	var data bytes.Buffer
	writeVarInt(&data, 0x01)
	binary.Write(&data, binary.BigEndian, payload)

	var packet bytes.Buffer
	writeVarInt(&packet, int32(data.Len()))
	packet.Write(data.Bytes())
	return packet.Bytes()
}

func readStatusResponse(conn net.Conn) (string, error) {
	length, err := readVarInt(conn)
	if err != nil {
		return "", fmt.Errorf("lecture longueur paquet : %w", err)
	}
	if length <= 0 || length > 1<<20 {
		return "", fmt.Errorf("longueur paquet invalide : %d", length)
	}

	packetData := make([]byte, length)
	if _, err := io.ReadFull(conn, packetData); err != nil {
		return "", fmt.Errorf("lecture donnees paquet : %w", err)
	}

	reader := bytes.NewReader(packetData)
	packetID, err := readVarInt(reader)
	if err != nil {
		return "", fmt.Errorf("lecture packet ID : %w", err)
	}
	if packetID != 0x00 {
		return "", fmt.Errorf("packet ID inattendu : 0x%02X", packetID)
	}

	strLen, err := readVarInt(reader)
	if err != nil {
		return "", fmt.Errorf("lecture longueur JSON : %w", err)
	}
	if strLen <= 0 || int(strLen) > reader.Len() {
		return "", fmt.Errorf("longueur JSON invalide : %d", strLen)
	}

	jsonBytes := make([]byte, strLen)
	if _, err := io.ReadFull(reader, jsonBytes); err != nil {
		return "", fmt.Errorf("lecture JSON : %w", err)
	}
	return string(jsonBytes), nil
}

func readPongResponse(conn net.Conn) (int64, error) {
	length, err := readVarInt(conn)
	if err != nil {
		return 0, err
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return 0, err
	}
	reader := bytes.NewReader(data)
	packetID, err := readVarInt(reader)
	if err != nil || packetID != 0x01 {
		return 0, fmt.Errorf("pong invalide")
	}
	var payload int64
	if err := binary.Read(reader, binary.BigEndian, &payload); err != nil {
		return 0, err
	}
	return payload, nil
}
