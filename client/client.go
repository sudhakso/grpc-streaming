package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/grpc-server-streaming/api"

	"google.golang.org/grpc"
)

const (
	username        = "admin"
	password        = "password"
	refreshDuration = 30 * time.Second

	// Add more paths
	streamServicePath = "/sapi.StreamService/"
)

func authMethods() map[string]bool {
	return map[string]bool{
		streamServicePath + "FetchResponse": true,
	}
}

func main() {

	// dial server
	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	defer conn.Close()
	// create auth client
	authClient := NewAuthClient(conn, username, password)
	interceptor, err := NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor:", err)
	}

	conn2, err := grpc.Dial(
		":50005",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	defer conn2.Close()

	// create stream connection with interceptor
	client := pb.NewStreamServiceClient(conn2)
	in := &pb.Request{Id: 1}
	stream, err := client.FetchResponse(context.Background(), in)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //means stream is finished
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			log.Printf("Resp received: %s", resp.Result)
		}
	}()

	<-done //we will wait until all response is received
	log.Printf("finished")
}
