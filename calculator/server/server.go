package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/wexel-nath/grpc-go-course/calculator/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {}

func (*server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	fmt.Println(fmt.Sprintf("Sum was invoked with %v", req))

	sum := req.GetFirst() + req.GetSecond()

	resp := &pb.SumResponse{
		Result: sum,
	}

	return resp, nil
}

func (*server) Primes(req *pb.PrimesRequest, stream pb.CalculatorService_PrimesServer) error {
	fmt.Println(fmt.Sprintf("Primes was invoked with %v", req))

	number := req.GetNumber()
	factor := int64(2)

	for number > 1 {
		if number % factor == 0 {
			resp := &pb.PrimesResponse{
				Prime: factor,
			}

			err := stream.Send(resp)
			if err != nil {
				fmt.Println(fmt.Sprintf("failed to send to stream: %v", req))
			}

			number = number / factor
		} else {
			factor++
		}
	}
	return nil
}

func (*server) Average(stream pb.CalculatorService_AverageServer) error {
	fmt.Println("Average was invoked")

	sum := float32(0)
	count := float32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if count == 0 {
				// don't divide by zero
				count = 1
			}
			if err = stream.SendAndClose(&pb.AverageResponse{
				Result: sum / count,
			}); err != nil {
				log.Fatalf("failed to send and close: %v/n", err)
			}
			break
		}
		if err != nil {
			fmt.Println("failed to receive request:", err)
		}

		sum += req.GetNumber()
		count++
	}
	return nil
}

func (*server) Maximum(stream pb.CalculatorService_MaximumServer) error {
	fmt.Println("Maximum was invoked")

	max := int64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println("failed to receive request:", err)
		}

		if req.GetNumber() > max {
			max = req.GetNumber()

			err = stream.Send(&pb.MaximumResponse{ Result: max })
			if err != nil {
				log.Fatalf("failed to send: %v/n", err)
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error) {
	fmt.Println("SquareRoot was invoked")

	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Received a negative number: %d",
			number,
		)
	}

	return &pb.SquareRootResponse{ Result: math.Sqrt(float64(number)) }, nil
}

func main() {
	address := "0.0.0.0:50051"
	fmt.Println("Listening for grpc on:", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
