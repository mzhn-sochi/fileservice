package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"s3-service/api/s3"
	"s3-service/internal/handler"
	"s3-service/internal/interceptor"
	"s3-service/internal/service"
	"s3-service/internal/storage"
	"syscall"
)

var (
	debug bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug")
}

func main() {
	flag.Parse()

	if debug {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatalf("godotenv.Load: %s\n", err.Error())
		}
	}

	s, err := storage.New(false)
	if err != nil {
		log.Fatalf("storage.New: %s\n", err.Error())
	}

	serv := service.New(s)
	server := grpc.NewServer(grpc.ChainStreamInterceptor(interceptor.LoggingStreamInterceptor()))
	reflection.Register(server)
	adapter := handler.New(serv)

	s3.RegisterS3Server(
		server,
		adapter,
	)
	go func(s *grpc.Server) {
		port := os.Getenv("APP_PORT")
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			panic(fmt.Errorf("cannot bind port %s: %w", port, err))
		}
		fmt.Printf("\nServer started on %s port\n\n", port)
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}(server)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGSEGV, syscall.SIGINT)

	<-signalChan

	server.GracefulStop()
}
