package zkp_test

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"testing"

	"github.com/mindaugasrukas/zkp_example/zkp"
	"github.com/mindaugasrukas/zkp_example/zkp/gen/zkp_pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)


func TestReadPacket(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct{
		input string
		expectedSize int
		expectedType string
	}{
		"RegisterRequest packet": {
			input: "register_request_packet.bin",
			expectedSize: 78,
			expectedType: "RegisterRequest",
		},
		"RegisterResponse packet": {
			input: "register_response_packet.bin",
			expectedSize: 69,
			expectedType: "RegisterResponse",
		},
		"RegisterResponse-error packet": {
			input: "register_response_packet-user_exists_error.bin",
			expectedSize: 88,
			expectedType: "RegisterResponse",
		},
		"AuthRequest packet": {
			input: "auth_request_packet.bin",
			expectedSize: 70,
			expectedType: "AuthRequest",
		},
		"ChallengeResponse packet": {
			input: "challenge_response_packet.bin",
			expectedSize: 72,
			expectedType: "ChallengeResponse",
		},
		"AnswerRequest packet": {
			input: "answer_request_packet.bin",
			expectedSize: 64,
			expectedType: "AnswerRequest",
		},
		"AuthResponse packet": {
			input: "auth_response_packet.bin",
			expectedSize: 61,
			expectedType: "AuthResponse",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server, client := net.Pipe()
			go func() {
				f, err := os.Open(path.Join("testdata/", test.input))
				defer f.Close()
				assert.NoError(err)
				data, err := ioutil.ReadAll(f)
				assert.NoError(err)
				server.Write(data)
				server.Close()
			}()

			bytes, err := zkp.ReadPacket(client)
			assert.NoError(err)
			client.Close()
			assert.Equal(test.expectedSize, len(bytes))

			var envelope zkp_pb.EnvelopeMessage
			err = proto.Unmarshal(bytes, &envelope)
			assert.NoError(err)
			assert.Equal(test.expectedType, envelope.Name)
		})
	}
}

func TestReadMessage(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct{
		input string
		expectedType proto.Message
	}{
		"RegisterRequest packet": {
			input: "register_request_packet.bin",
			expectedType: &zkp_pb.RegisterRequest{
				User: "max",
				Commits: []*zkp_pb.RegisterRequest_Commits{
					{
						Y1: []byte{0x10},
						Y2: []byte{0xc},
					},
				},
			},
		},
		"RegisterResponse packet": {
			input: "register_response_packet.bin",
			expectedType: &zkp_pb.RegisterResponse{
				Result: true,
			},
		},
		"RegisterResponse-error packet": {
			input: "register_response_packet-user_exists_error.bin",
			expectedType: &zkp_pb.RegisterResponse{
				Result: false,
				Error: "user already exists",
			},
		},
		"AuthRequest packet": {
			input: "auth_request_packet.bin",
			expectedType: &zkp_pb.AuthRequest{
				User: "max",
				Commits: []*zkp_pb.AuthRequest_Commits{
					{
						R1: []byte{0xd},
						R2: []byte{2},
					},
				},
			},
		},
		"ChallengeResponse packet": {
			input: "challenge_response_packet.bin",
			expectedType: &zkp_pb.ChallengeResponse{
				Challenge: []byte{1},
			},
		},
		"AnswerRequest packet": {
			input: "answer_request_packet.bin",
			expectedType: &zkp_pb.AnswerRequest{
				Answer: []byte{7},
			},
		},
		"AuthResponse packet": {
			input: "auth_response_packet.bin",
			expectedType: &zkp_pb.AuthResponse{
				Result: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server, client := net.Pipe()
			go func() {
				f, err := os.Open(path.Join("testdata/", test.input))
				defer f.Close()
				assert.NoError(err)
				data, err := ioutil.ReadAll(f)
				assert.NoError(err)
				server.Write(data)
				server.Close()
			}()

			message, err := zkp.ReadMessage(client)
			assert.NoError(err)
			client.Close()
			assert.Equal(test.expectedType.ProtoReflect().Type(), message.ProtoReflect().Type())
			assert.Equal(test.expectedType, message)
		})
	}
}

func TestSendMessage(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct{
		message     proto.Message
		expectedBin string
	}{
		"RegisterRequest packet": {
			expectedBin: "register_request_packet.bin",
			message: &zkp_pb.RegisterRequest{
				User: "max",
				Commits: []*zkp_pb.RegisterRequest_Commits{
					{
						Y1: []byte{0x10},
						Y2: []byte{0xc},
					},
				},
			},
		},
		"RegisterResponse packet": {
			expectedBin: "register_response_packet.bin",
			message: &zkp_pb.RegisterResponse{
				Result: true,
			},
		},
		"RegisterResponse-error packet": {
			expectedBin: "register_response_packet-user_exists_error.bin",
			message: &zkp_pb.RegisterResponse{
				Result: false,
				Error: "user already exists",
			},
		},
		"AuthRequest packet": {
			expectedBin: "auth_request_packet.bin",
			message: &zkp_pb.AuthRequest{
				User: "max",
				Commits: []*zkp_pb.AuthRequest_Commits{
					{
						R1: []byte{0xd},
						R2: []byte{2},
					},
				},
			},
		},
		"ChallengeResponse packet": {
			expectedBin: "challenge_response_packet.bin",
			message: &zkp_pb.ChallengeResponse{
				Challenge: []byte{1},
			},
		},
		"AnswerRequest packet": {
			expectedBin: "answer_request_packet.bin",
			message: &zkp_pb.AnswerRequest{
				Answer: []byte{7},
			},
		},
		"AuthResponse packet": {
			expectedBin: "auth_response_packet.bin",
			message: &zkp_pb.AuthResponse{
				Result: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server, client := net.Pipe()
			go func() {
				err := zkp.SendMessage(server, test.message)
				assert.NoError(err)
				server.Close()
			}()

			bytes, err := zkp.ReadPacket(client)
			assert.NoError(err)
			client.Close()

			f, err := os.Open(path.Join("testdata/", test.expectedBin))
			defer f.Close()
			assert.NoError(err)
			expectedBytes, err := ioutil.ReadAll(f)
			assert.NoError(err)

			assert.Equal(expectedBytes[4:], bytes)
		})
	}
}
