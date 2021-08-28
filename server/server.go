package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/grpc-server-streaming/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s server) FetchResponse(in *pb.Request, srv pb.StreamService_FetchResponseServer) error {

	log.Printf("fetch response for id : %d", in.Id)

	// use wait group to allow process to be concurrent
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)

		// go func
		go func(count int64) {
			defer wg.Done()

			time.Sleep(time.Duration(count) * time.Second)
			resp := pb.Response{Result: fmt.Sprintf("Request #%d for Id:%d", count, in.Id)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("sending error : %v", err)
			}
			log.Printf("Finishing request number : %d", count)
		}(int64(i))
	}
	wg.Wait()
	return nil
}

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func createUser(userStore UserStore, userName string, password string, role string) error {
	user, err := NewUser(userName, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func seedUsers(userStore UserStore) error {
	//admin ...
	err := createUser(userStore, "admin", "password", "admin")
	if err != nil {
		log.Printf("failure creating seed user: admin")
		return err
	}

	//member ...
	return createUser(userStore, "member", "password", "member")
}

// path -> role
func accessibleRoles() map[string][]string {
	const streamServicePath = "/sapi.StreamService/"
	return map[string][]string{
		streamServicePath + "FetchResponse": {"admin", "member"},
	}
}

func main() {

	//create listener
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Setup auth ...
	userStore := NewInMemoryUserStore()
	jwtManager := NewJWTManager(secretKey, tokenDuration)

	// Setup interceptors
	interceptor := NewAuthInterceptor(jwtManager, accessibleRoles())

	//Setup some seed users ...
	if err := seedUsers(userStore); err != nil {
		log.Fatalf("Initialiing users Failed... %v", err)
	}

	// Streaming server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	// Auth service
	authServer := NewAuthServer(userStore, jwtManager)

	// Register all services to grpc server
	pb.RegisterAuthServiceServer(s, authServer)
	pb.RegisterStreamServiceServer(s, server{})

	//reflection...
	reflection.Register(s)

	log.Println("starting server")
	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
