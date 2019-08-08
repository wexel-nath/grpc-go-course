package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/wexel-nath/grpc-go-course/calculator/pb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'am a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer cc.Close()

	c := pb.NewCalculatorServiceClient(cc)
	//fmt.Println(fmt.Sprintf("created client: %f", c))

	//doUnary(c)

	//doServerStream(c)

	//doClientStream(c)

	doBiDiStream(c)
}

func doUnary(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")

	req := &pb.SumRequest{
		First:  10,
		Second: 3,
	}

	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculator with RPC: %v", err)
	}

	fmt.Printf("Response from Calculator: %v\n", resp.GetResult())
}

func doServerStream(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Stream RPC...")

	req := &pb.PrimesRequest{
		Number: 12345678912355,
	}

	stream, err := c.Primes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculator with RPC: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error with Primes RPC: %v", err)
		}

		fmt.Printf("Response from Calculator: %v\n", resp.GetPrime())
	}
}

func doClientStream(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Stream RPC...")

	numbers := []float32{ 1, 2, 3, 4 }

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error while calling Calculator with Client Stream RPC: %v", err)
	}

	for _, n := range numbers {
		fmt.Println("Sending a number:", n)
		err = stream.Send(&pb.AverageRequest{ Number: n })
		if err != nil {
			log.Fatalf("error while sending to Average: %v", err)
		}
		time.Sleep(time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving Average: %v", err)
	}

	fmt.Println("Average result:", resp.GetResult())
}

func doBiDiStream(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDi Stream RPC...")

	numbers := []int64{ 1, 5, 3, 6, 2, 20 }

	stream, err := c.Maximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling Calculator with BiDi Stream RPC: %v", err)
	}

	go func() {
		for _, n := range numbers {
			time.Sleep(time.Second)
			fmt.Println("Sending a number:", n)
			err = stream.Send(&pb.MaximumRequest{ Number: n })
			if err != nil {
				log.Fatalf("error while sending to Maximum: %v", err)
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
			log.Fatalf("error while receiving from Maximum: %v", err)
		}

		fmt.Println("Current Maximum:", resp.GetResult())
	}
}
