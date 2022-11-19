package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	gRPC "github.com/emjakobsen1/dsys5/proto"
	"google.golang.org/grpc"
)

var server *Server
var duration = flag.Int("n", 120, "duration for the auction.")
var ports = [3]string{"9080", "9081", "9082"}

type Server struct {
	// an interface that the server needs to have
	gRPC.UnimplementedAuctionhouseServer
	port              string
	currentHighestBid int32
	running           bool
	// here you can impliment other fields that you want
}

func (s *Server) Bid(ctx context.Context, Amount *gRPC.Amount) (*gRPC.Ack, error) {
	log.Printf("Current highest bid %d \n", s.currentHighestBid)
	if Amount.Amount > s.currentHighestBid {
		s.currentHighestBid = Amount.Amount
		ack := &gRPC.Ack{Status: "SUCCESS, amount:" + fmt.Sprint(Amount.Amount) + " bidder: " + Amount.Bidder}
		return ack, nil
	} else {
		ack := &gRPC.Ack{Status: "FAIL: must be higher than current highest bid. " + Amount.Bidder}
		return ack, nil

	}

}

func (s *Server) Result(ctx context.Context, Void *gRPC.Void) (*gRPC.Outcome, error) {
	outcome := &gRPC.Outcome{Status: "RUNNING", HighestBid: s.currentHighestBid}
	return outcome, nil
}

func main() {
	flag.Parse()
	log.Printf("AUCTION started with duration: %d", *duration)
	go launch(ports[0])
	go launch(ports[1])
	launch(ports[2])

}

func getServer(serverPort string) *Server {
	s := &Server{
		port:              serverPort,
		currentHighestBid: 0,
		running:           true,
	}
	log.Println(s)
	return s

}

func launch(portAddress string) {

	list, err := net.Listen("tcp", "localhost:"+portAddress)
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	server = getServer((portAddress))
	gRPC.RegisterAuctionhouseServer(grpcServer, server)
	log.Printf("Server running on port %s", server.port)
	/*
		time.Sleep(time.Second * time.Duration(*duration))
		server.running = false
		log.Println("AUCTION has ended")
		//os.Exit(3)*/
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}
