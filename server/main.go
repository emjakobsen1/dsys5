package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	gRPC "github.com/emjakobsen1/dsys5/proto"
	"google.golang.org/grpc"
)

var server *Server
var duration = flag.Int("n", 120, "duration for the auction.")
var currentWinner string

//var ports = [3]string{"9080", "9081", "9082"}

type Server struct {
	// an interface that the server needs to have
	gRPC.UnimplementedAuctionhouseServer
	port              string
	currentHighestBid int32
	running           bool
	// here you can impliment other fields that you want
}

func (s *Server) Bid(ctx context.Context, Amount *gRPC.Amount) (*gRPC.Ack, error) {
	log.Printf("SERVER: %s, highest bid: %d, user: %s, bids: %d \n", s.port, s.currentHighestBid, Amount.Bidder, Amount.Amount)
	if Amount.Amount > s.currentHighestBid {
		s.currentHighestBid = Amount.Amount
		currentWinner = Amount.Bidder
		ack := &gRPC.Ack{Status: "SUCCESS, amount:" + fmt.Sprint(Amount.Amount) + " bidder: " + Amount.Bidder, Id: Amount.Id}
		return ack, nil
	} else {
		ack := &gRPC.Ack{Status: "FAIL: must be higher than current highest bid. " + Amount.Bidder, Id: Amount.Id}
		return ack, nil

	}

}

func (s *Server) Result(ctx context.Context, Void *gRPC.Void) (*gRPC.Outcome, error) {
	outcome := &gRPC.Outcome{Status: "RUNNING", HighestBid: s.currentHighestBid, Id: Void.Id}
	return outcome, nil
}

func main() {
	flag.Parse()
	log.Printf("AUCTION started with duration: %d seconds", *duration)
	go launch("9080")
	go launch("9081")
	launch("9082")

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

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
func run() {
	time.Sleep(time.Second * time.Duration(*duration))
	server.running = false

	log.Println("AUCTION has ended")
	log.Printf("Winner was: %s", currentWinner)
}
