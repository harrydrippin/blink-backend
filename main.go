package main

import (
	pb "blink-backend/blink"
	"blink-backend/database"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.Println("Initialized")

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 8080))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	database := database.GetInstance()
	defer database.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterBlinkServer(grpcServer, &BlinkServer{
		queue: &TaskQueue{
			receiveRequestQueue: make(map[string][]pb.ReceiveRequest),
			channelMap:          make(map[string]chan int),
		},
	})
	log.Println("gRPC server started at port 8080")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
