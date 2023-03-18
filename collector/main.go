package main

import (
	"collector/RBC/proposalpb"
	"collector/util"
	"collector/watch"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// proposalServer is used to implement proposalpb.ProposalReceive
type proposalServer struct {
	proposalpb.UnimplementedProposalHandleServer
}

// ProposalReceive implements proposalpb.ProposalReceive
func (ps *proposalServer) ProposalReceive(ctx context.Context, in *proposalpb.Proposal) (*proposalpb.ProposalResponse, error) {
	log.Printf("Received: %v", in.GetRound())
	return &proposalpb.ProposalResponse{}, nil
}

func main() {
	ip := os.Args[1]
	fmt.Println(ip)
	adress := "127.0.0.1:" + ip
	outputCh := make(chan string)

	lis, err := net.Listen("tcp", adress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ps := grpc.NewServer()
	proposalpb.RegisterProposalHandleServer(ps, &proposalServer{})
	log.Printf("server listening at %v", lis.Addr())
	go ps.Serve(lis)
	// if err := ps.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }

	go watch.WatchOutput(outputCh, "../output")

	for {
		select {
		case newOutput := <-outputCh:
			{
				fmt.Println("===>[Collector]New output is:", newOutput)
				ipList := util.GetIPAdress("../ipAdress")
				fmt.Println(ipList)

				for _, ipAdress := range ipList {
					if ipAdress != adress {
						fmt.Println(ipAdress)
						conn, err := grpc.Dial(ipAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
						if err != nil {
							log.Fatalf("did not connect: %v", err)
						}
						pc := proposalpb.NewProposalHandleClient(conn)

						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						_, err = pc.ProposalReceive(ctx, &proposalpb.Proposal{Round: "1", Sender: "1", Sig: "1"})
						if err != nil {
							log.Fatalf("could not greet: %v", err)
						}
					}
				}

			}
		}
	}
}
