package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "grpc-echo-server/protos"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedEchoServer
}

func (e *server) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("unary echoing message %q\n", req.Message)
	return &pb.EchoResponse{
		Message: fmt.Sprintf(req.Message),
	}, nil
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			fmt.Printf("server: error receiving from stream: %v\n", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Printf("bidi echoing message %q\n", in.Message)
		stream.Send(&pb.EchoResponse{Message: in.Message})
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at port %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}
