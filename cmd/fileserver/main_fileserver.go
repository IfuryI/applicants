package main

import (
	constants "bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"fmt"
	"log"
	"net"

	"bitbucket.org/projectiu7/backend/src/master/internal/proto"
	fileServerGrpc "bitbucket.org/projectiu7/backend/src/master/internal/services/fileserver/delivery/grpc"
	"google.golang.org/grpc"
)

func main() {
	handler := fileServerGrpc.NewFileServerHandlerServer()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", constants.FileServerPort))

	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	proto.RegisterFileServerHandlerServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
