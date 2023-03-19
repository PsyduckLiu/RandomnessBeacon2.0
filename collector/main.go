package main

import (
	"collector/RBC/proposalpb"
	"collector/config"
	"collector/signature"
	"collector/util"
	"collector/watch"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	blsSig "go.dedis.ch/dela/crypto/bls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// proposalServer is used to implement proposalpb.ProposalReceive
type proposalServer struct {
	proposalpb.UnimplementedProposalHandleServer
}

// ProposalReceive implements proposalpb.ProposalReceive
func (ps *proposalServer) ProposalReceive(ctx context.Context, in *proposalpb.Proposal) (*proposalpb.ProposalResponse, error) {
	fmt.Printf("===>[ProposalReceive]Received: Round: %v Sender: %v TC: %v\n", in.GetRound(), in.GetSender(), in.GetTc())

	result := signature.VerifySig(1, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig())
	fmt.Println("===>[ProposalReceive]Verify ", result)

	return &proposalpb.ProposalResponse{}, nil
}

type Collector struct {
	ID       int
	Address  string
	OutputCh chan string
	Signer   blsSig.Signer
	// pk,sk
}

// initialize a new Collector
func NewCollector(id int) *Collector {
	c := &Collector{
		ID:       id,
		Address:  "127.0.0.1:" + strconv.Itoa(2333+id),
		OutputCh: make(chan string),
		Signer:   blsSig.NewSigner(),
	}

	config.WriteKey(id, c.Signer.GetPublicKey())

	return c
}

func main() {
	id := os.Args[1]
	idInt, _ := strconv.Atoi(id)
	collecor := NewCollector(idInt)

	lis, err := net.Listen("tcp", collecor.Address)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from Collector]Failed to listen: %s", err))
	}

	ps := grpc.NewServer()
	proposalpb.RegisterProposalHandleServer(ps, &proposalServer{})
	go ps.Serve(lis)
	fmt.Printf("===>[Collector]Collector is listening at %v", lis.Addr())

	go watch.WatchOutput(collecor.OutputCh, "../output")

	for {
		select {
		case newOutput := <-collecor.OutputCh:
			{
				fmt.Println("===>[Collector]New output is:", newOutput)
				ipList := util.GetIPAddress("../ipAddress")
				fmt.Println("===>[Collector]IP list is:", ipList)

				for _, ipAddress := range ipList {
					if ipAddress != collecor.Address {
						conn, err := grpc.Dial(ipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
						if err != nil {
							fmt.Println("===>[!!!Collector]did not connect:", err)
							continue
						}

						pc := proposalpb.NewProposalHandleClient(conn)
						ctx, cancel := context.WithTimeout(context.Background(), time.Second)

						// tcList := [3]string{"1", "2", "3"}
						var tcList []string
						tcList = append(tcList, "1")
						tcList = append(tcList, "2")
						tcList = append(tcList, "3")

						sig := signature.GenerateSig(1, "1", tcList, collecor.Signer)
						fmt.Println("Signature:", sig)
						_, err = pc.ProposalReceive(ctx, &proposalpb.Proposal{Type: 1, Round: "1", Sender: id, Tc: tcList, Sig: sig})
						if err != nil {
							fmt.Println("===>[!!!Collector]Failed to response:", err)
							continue
						}

						cancel()
					}
				}
			}
		}
	}
}
