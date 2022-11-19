package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	gRPC "github.com/emjakobsen1/dsys5/proto"
	"google.golang.org/grpc"
)

var bidderName = flag.String("name", "Anonymous user", "type -name <yourname> to set name of the bidder.")
var servers []gRPC.AuctionhouseClient
var ServerConn map[gRPC.AuctionhouseClient]*grpc.ClientConn
var _ports = [3]string{"9080", "9081", "9082"}

func main() {
	flag.Parse()
	ServerConn = make(map[gRPC.AuctionhouseClient]*grpc.ClientConn)
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()
	//c := gRPC.NewAuctionhouseClient(conn)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())
	timeContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, port := range _ports {
		log.Printf("client %s: Attempts to dial on port %s\n", *bidderName, port)
		conn, err := grpc.DialContext(timeContext, fmt.Sprintf(":%s", port), opts...)
		if err != nil {
			log.Printf("Fail to Dial : %v", err)
			continue
		}
		var s = gRPC.NewAuctionhouseClient(conn)
		servers = append(servers, s)
		ServerConn[s] = conn
		fmt.Println(conn.GetState().String())
	}
	/*
		for i := 0; i < 10; i++ {
			Bid(c, (int32(i)))
			Result(c)
			time.Sleep(5 * time.Second)
		}*/

}
func Result(c gRPC.AuctionhouseClient) {
	msg := gRPC.Void{}
	response, err := c.Result(context.Background(), &msg)
	if err != nil {
		log.Fatalf("Error calling Result RPC: %s ", err)
	}
	fmt.Printf("%s Highest bid: %d \n", response.Status, response.HighestBid)
}
func Bid(c gRPC.AuctionhouseClient, amountToBid int32) {
	msg := gRPC.Amount{Bidder: *bidderName, Amount: amountToBid}
	response, err := c.Bid(context.Background(), &msg)
	if err != nil {
		log.Fatalf("error calling Bid RPC %s ", err)
	}
	fmt.Printf("-- %s \n", response.Status)
}
