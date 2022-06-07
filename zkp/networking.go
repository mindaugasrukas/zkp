package zkp

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ReadPacket reads message bytes from the TCP connection
//
// Packet structure:
// 4 bytes    | N bytes
// -----------+--------------------------------
// the size N | message content of the length N
//
// This is required as we reuse the same TCP connection during the communication process.
// At least we need to know the packet size to decode correctly.
//
// todo: add a full application packet header: version, size, type, etc.
func ReadPacket(conn net.Conn) ([]byte, error) {
	var psize, size int
	in := make([]byte, 0)
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
			psize = int(binary.LittleEndian.Uint32(tmp[:4]))
			size = psize
			n -= 4
			tmp = tmp[4:]
		}

		in = append(in, tmp[:n]...)

		size -= n
		// check if we received everything
		if size <= 0 {
			break
		}
	}
	return in[:psize], nil
}

// ReadMessage reads and parses the bytes from the TCP connection into proto Message
// Decode proto messages using envelope information.
func ReadMessage(conn net.Conn) (proto.Message, error) {
	in, err := ReadPacket(conn)
	if err != nil {
		return nil, err
	}

	var envelope zkp_pb.EnvelopeMessage
	if err = proto.Unmarshal(in, &envelope); err != nil {
		return nil, err
	}

	var msg proto.Message

	switch envelope.Name {
	case "RegisterRequest":
		msg = &zkp_pb.RegisterRequest{}
	case "RegisterResponse":
		msg = &zkp_pb.RegisterResponse{}
	case "AuthRequest":
		msg = &zkp_pb.AuthRequest{}
	case "AuthResponse":
		msg = &zkp_pb.AuthResponse{}
	case "AnswerRequest":
		msg = &zkp_pb.AnswerRequest{}
	case "ChallengeResponse":
		msg = &zkp_pb.ChallengeResponse{}
	}

	if err = envelope.Message.UnmarshalTo(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// SendMessage writes the proto Message to the TCP connection
// for packet structure see ReadPacket
// Envelope the proto messages for easier to decode them.
func SendMessage(conn net.Conn, message proto.Message) error {
	// Envelope Messages
	any, err := anypb.New(message)
	envelope := &zkp_pb.EnvelopeMessage{
		Name:    string(message.ProtoReflect().Descriptor().Name()),
		Message: any,
	}

	out, err := proto.Marshal(envelope)
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
