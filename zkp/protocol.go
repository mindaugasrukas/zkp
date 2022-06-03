package zkp

import (
	"encoding/binary"
	"io"
	"net"

	"google.golang.org/protobuf/proto"
)

func ReadMessage(conn net.Conn) ([]byte, error) {
	var size int
	in := make([]byte, 0, 10240)
	tmp := make([]byte, 4096)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		// calculate expected content size
		if n >= 4 && size == 0 {
			size = int(binary.LittleEndian.Uint32(tmp[:4]))
			n -= 4
			tmp = tmp[4:]
		}

		in = append(in, tmp[:n]...)

		size -= n
		// check if we received everything
		if size == 0 {
			break
		}
	}

	return in, nil
}


func ParseMessage(conn net.Conn, msg proto.Message) error {
	in, err := ReadMessage(conn)
	if err != nil {
		return err
	}

	if err = proto.Unmarshal(in, msg); err != nil {
		return err
	}
	return nil
}

func SendMessage(conn net.Conn, message proto.Message) error {
	out, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	// send the content size
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(out)))
	if _, err = conn.Write(size); err != nil {
		return err
	}
	// send the content
	if _, err = conn.Write(out); err != nil {
		return err
	}
	return nil
}
