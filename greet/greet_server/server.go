package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/wexel-nath/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()

	resp := &greetpb.GreetResponse{
		Result: "Hello " + firstName,
	}

	return resp, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		resp := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %s number %d", firstName, i),
		}

		if err := stream.Send(resp); err != nil {
			log.Info(err)
		}
		time.Sleep(time.Second)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function invoked with a stream request")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function invoked with a stream request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: fmt.Sprintf("Hello %s!", firstName),
		})
		if err != nil {
			log.Fatalf("error while sending data to client: %v", err)
		}
	}
}

func main() {
	address := "0.0.0.0:50051"
	fmt.Println("Listening for grpc on:", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
