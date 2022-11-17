package main

import (
	"context"
	"log"
	"net"

	gRPC "github.com/emjakobsen1/dsys5/proto"
	"google.golang.org/grpc"
)

type Server struct {
	// an interface that the server needs to have
	gRPC.UnimplementedAuctionhouseServer
	currentHighestBid int32
	// here you can impliment other fields that you want
}

var server *Server

func (s *Server) Bid(ctx context.Context, Amount *gRPC.Amount) (*gRPC.Ack, error) {
	ack := &gRPC.Ack{Status: "Bid/ended"}
	return ack, nil
}

func (s *Server) Result(ctx context.Context, Void *gRPC.Void) (*gRPC.Outcome, error) {
	outcome := &gRPC.Outcome{Status: "Result/ended", HighestBid: s.currentHighestBid}
	return outcome, nil
}

func main() {
	launchServer()
}

func launchServer() {
	list, err := net.Listen("tcp", "localhost:9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	gRPC.RegisterAuctionhouseServer(grpcServer, &Server{currentHighestBid: 0})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}
