package main

import (
	constants "bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"fmt"
	"log"
	"net"

	"bitbucket.org/projectiu7/backend/src/master/internal/proto"
	sessionsGrpc "bitbucket.org/projectiu7/backend/src/master/internal/services/sessions/delivery/grpc"
	sessionsRepo "bitbucket.org/projectiu7/backend/src/master/internal/services/sessions/repository"
	sessionsUC "bitbucket.org/projectiu7/backend/src/master/internal/services/sessions/usecase"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", constants.RedisPort),
		Password: "",
		DB:       0,
	})
	fmt.Println("localhost:%s", constants.RedisPort)
	p, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to create redis client", p, err)
	}
	fmt.Println(p)

	repo := sessionsRepo.NewRedisRepository(rdb)
	usecase := sessionsUC.NewUseCase(repo)
	handler := sessionsGrpc.NewAuthHandlerServer(usecase)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", constants.AuthPort))

	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	proto.RegisterAuthHandlerServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
