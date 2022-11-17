package main

import (
	"context"
	"fmt"
	"log"
	"time"

	gRPC "github.com/emjakobsen1/dsys5/proto"
	"google.golang.org/grpc"
)

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()
	c := gRPC.NewAuctionhouseClient(conn)

	for {
		SendVoid(c)
		time.Sleep(5 * time.Second)
	}

}
func SendVoid(c gRPC.AuctionhouseClient) {
	msg := gRPC.Void{}
	response, err := c.Result(context.Background(), &msg)
	if err != nil {
		log.Fatalf("Error calling Result RPC: %s ", err)
	}
	fmt.Printf("%s Highest bid: %d \n", response.Status, response.HighestBid)
}
