package main

import (
	"log"
	"net"
	"os"

	"github.com/zv0n/webdav-proxy/configuration"
	"github.com/zv0n/webdav-proxy/webdavrpc"
	"google.golang.org/grpc"
)

const configPath = "/etc/webdav-proxy.conf"

func main() {
	config, err := configuration.ParseConfigFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Could not read config file: \"%s\" - %v", configPath, err)
	}

	if err := os.Remove(config.SocketPath); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to remove %s, error: %s", config.SocketPath, err.Error())
	}

	listener, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	webdavServer := webdavrpc.Server{
		Config: config,
	}

	grpcServer := grpc.NewServer()

	webdavrpc.RegisterMountServiceServer(grpcServer, &webdavServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
