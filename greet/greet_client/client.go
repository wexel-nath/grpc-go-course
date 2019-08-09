package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/wexel-nath/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'am a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	//doBiDiStreaming(c)

	doUnaryWithDeadline(c, 5 * time.Second)
	doUnaryWithDeadline(c, 1 * time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nathan",
			LastName:  "Welch",
		},
	}

	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet with RPC: %v", err)
	}

	fmt.Printf("Response from Greet: %v\n", resp.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nathan",
			LastName:  "Welch",
		},
	}

	resultStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes with RPC: %v", err)
	}

	for {
		resp, err := resultStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		fmt.Printf("Response from GreetManyTimes: %v\n", resp.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	greetings := []*greetpb.Greeting{
		{
			FirstName: "Nathan",
		},
		{
			FirstName: "Tristan",
		},
		{
			FirstName: "Rav",
		},
		{
			FirstName: "Alex",
		},
		{
			FirstName: "Callum",
		},
	}

	for _, greeting := range greetings {
		fmt.Println("Sending a greeting to", greeting.FirstName)
		err = stream.Send(&greetpb.LongGreetRequest{Greeting:greeting})
		if err != nil {
			log.Fatalf("error while sending to LongGreet: %v", err)
		}
		time.Sleep(time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving LongGreet: %v", err)
	}

	fmt.Printf("Response from LongGreet: %v\n", resp.GetResult())
}


func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone: %v", err)
	}

	greetings := []*greetpb.Greeting{
		{
			FirstName: "Nathan",
		},
		{
			FirstName: "Tristan",
		},
		{
			FirstName: "Rav",
		},
		{
			FirstName: "Alex",
		},
		{
			FirstName: "Callum",
		},
	}

	go func() {
		for _, greeting := range greetings {
			time.Sleep(time.Second)
			fmt.Println("Sending a greeting to", greeting.FirstName)
			err = stream.Send(&greetpb.GreetEveryoneRequest{ Greeting: greeting })
			if err != nil {
				log.Fatalf("error while sending to LongGreet: %v", err)
			}
		}
		stream.CloseSend()
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving LongGreet: %v", err)
		}
		fmt.Printf("Response from GreetEveryone: %v\n", resp.GetResult())
	}
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a Unary GreetWithDeadline RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nathan",
			LastName:  "Welch",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline exceeded")
			} else {
				fmt.Println("Unexpected error status error:", s.Message())
			}
		} else {
			log.Fatalf("error while calling Greet with RPC: %v", err)
		}
		return
	}

	fmt.Printf("Response from Greet: %v\n", resp.GetResult())
}
