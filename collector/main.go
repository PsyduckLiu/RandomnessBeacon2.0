package main

import (
	"collector/RBC/ackpb"
	"collector/RBC/notifypb"
	"collector/RBC/proposalpb"
	"collector/config"
	"collector/signature"
	"collector/util"
	"collector/watch"
	"context"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	blsSig "go.dedis.ch/dela/crypto/bls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var totalLocalTC [][]string
var totalGlobalTC []string
var totalLocalSig [][]string
var totalLocalFrom [][]string
var ackChan chan int
var notifyChan chan int

// Reset start a new round
func reset() {
	f := config.GetF()
	totalLocalTC = make([][]string, 3*f+1)
	totalGlobalTC = make([]string, 0)

	totalLocalSig = make([][]string, 3*f+1)
	totalLocalFrom = make([][]string, 3*f+1)

	ackChan = make(chan int, 5)
	notifyChan = make(chan int, 5)
}

// proposalServer is used to implement proposalpb.ProposalReceive
type proposalServer struct {
	proposalpb.UnimplementedProposalHandleServer
}

// ProposalReceive implements proposalpb.ProposalReceive
func (ps *proposalServer) ProposalReceive(ctx context.Context, in *proposalpb.Proposal) (*proposalpb.ProposalResponse, error) {
	fmt.Printf("===>[ProposalReceive]Received: Round: %v Sender: %v TC: %v\n", in.GetRound(), in.GetSender(), in.GetTc())

	result := signature.VerifySig(1, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
	fmt.Println("===>[ProposalReceive]Verify ", result)

	sender, _ := strconv.Atoi(in.GetSender())
	if result && totalLocalTC[sender] == nil {
		totalLocalTC[sender] = in.GetTc()
		fmt.Println(totalLocalTC[sender])
		go func() {
			ackChan <- sender
		}()
	}

	return &proposalpb.ProposalResponse{}, nil
}

// ackServer is used to implement ackpb.ackReceive
type ackServer struct {
	ackpb.UnimplementedAckHandleServer
}

// AckReceive implements ackpb.ackReceive
func (as *ackServer) AckReceive(ctx context.Context, in *ackpb.Ack) (*ackpb.AckResponse, error) {
	f := config.GetF()

	fmt.Printf("\n===>[AckReceive]Received: Round: %v Sender: %v From :%v TC: %v\n", in.GetRound(), in.GetSender(), in.GetFrom(), in.GetTc())
	result := signature.VerifySig(2, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetFrom())
	fmt.Println("===>[AckReceive]Verify ", result)

	ackTC := in.GetTc()
	sender, _ := strconv.Atoi(in.GetSender())
	if reflect.DeepEqual(ackTC, totalLocalTC[sender]) {
		totalLocalSig[sender] = append(totalLocalSig[sender], in.GetSig())
		totalLocalFrom[sender] = append(totalLocalFrom[sender], in.GetFrom())
	}

	if len(totalLocalSig[sender]) == 2*f+1 {
		fmt.Println("It's time for certificate")
		go func() {
			notifyChan <- sender
		}()
	}

	return &ackpb.AckResponse{}, nil
}

// notifyServer is used to implement notifypb.notifyReceive
type notifyServer struct {
	notifypb.UnimplementedNotifyHandleServer
}

// notifyReceive implements notifypb.notifyReceive
func (ps *notifyServer) NotifyReceive(ctx context.Context, in *notifypb.Notify) (*notifypb.NotifyResponse, error) {
	fmt.Printf("===>[notifyReceive]Received: Round: %v Sender: %v From :%v TC: %v\n", in.GetRound(), in.GetSender(), in.GetFrom(), in.GetTc())

	result := signature.VerifySig(3, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetFrom())
	fmt.Println("===>[notifyReceive]Verify ", result)

	if signature.VerifyAggSig(3, in.GetRound(), in.GetSender(), in.GetTc(), in.GetAggsig(), in.GetAggfrom()) {
		fmt.Println("GOt it!!")
	}

	return &notifypb.NotifyResponse{}, nil
}

type Collector struct {
	Round    int
	ID       int
	Address  string
	LocalTC  []string
	GlobalTC []string
	OutputCh chan string
	Signer   blsSig.Signer
}

// initialize a new Collector
func NewCollector(id int) *Collector {
	c := &Collector{
		Round:    0,
		ID:       id,
		Address:  "127.0.0.1:" + strconv.Itoa(2333+id),
		OutputCh: make(chan string),
		Signer:   blsSig.NewSigner(),
	}

	c.LocalTC = append(c.LocalTC, "1")
	c.LocalTC = append(c.LocalTC, "2")
	c.LocalTC = append(c.LocalTC, "3")

	config.WriteKey(id, c.Signer.GetPublicKey())

	return c
}

func main() {
	id := os.Args[1]
	idInt, _ := strconv.Atoi(id)
	collector := NewCollector(idInt)

	lis, err := net.Listen("tcp", collector.Address)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from Collector]Failed to listen: %s", err))
	}

	ps := grpc.NewServer()
	proposalpb.RegisterProposalHandleServer(ps, &proposalServer{})
	ackpb.RegisterAckHandleServer(ps, &ackServer{})
	notifypb.RegisterNotifyHandleServer(ps, &notifyServer{})
	go ps.Serve(lis)
	fmt.Printf("===>[Collector]Collector is listening at %v", lis.Addr())

	go watch.WatchOutput(collector.OutputCh, "../output")

	for {
		select {
		case newOutput := <-collector.OutputCh:
			{
				fmt.Println("===>[Collector]New output is:", newOutput)
				reset()
				collector.Round++

				// proposal phase
				ipList := util.GetIPAddress("../ipAddress")
				fmt.Println("===>[Collector]IP list is:", ipList)

				for _, ipAddress := range ipList {
					// if ipAddress != collector.Address {
					conn, err := grpc.Dial(ipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("===>[!!!Collector]did not connect:", err)
						continue
					}

					pc := proposalpb.NewProposalHandleClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)

					sig := signature.GenerateSig(1, strconv.Itoa(collector.Round), strconv.Itoa(collector.ID), collector.LocalTC, collector.Signer)
					_, err = pc.ProposalReceive(ctx, &proposalpb.Proposal{Type: 1, Round: strconv.Itoa(collector.Round), Sender: id, Tc: collector.LocalTC, Sig: sig})
					if err != nil {
						fmt.Println("Send to", ipAddress)
						fmt.Println("===>[!!!Collector]Failed to response:", err)
						continue
					}

					cancel()
					// }
				}
			}

		case newProposal := <-ackChan:
			{
				fmt.Println("===>[Collector]New proposal for:", newProposal)

				// forward phase
				ipList := util.GetIPAddress("../ipAddress")

				for _, ipAddress := range ipList {
					// if ipAddress != collector.Address {
					conn, err := grpc.Dial(ipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("===>[!!!Collector]did not connect:", err)
						continue
					}

					ac := ackpb.NewAckHandleClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)

					sig := signature.GenerateSig(2, strconv.Itoa(collector.Round), strconv.Itoa(newProposal), collector.LocalTC, collector.Signer)
					_, err = ac.AckReceive(ctx, &ackpb.Ack{Type: 2, Round: strconv.Itoa(collector.Round), Sender: strconv.Itoa(newProposal), From: strconv.Itoa(collector.ID), Tc: totalLocalTC[newProposal], Sig: sig})
					if err != nil {
						fmt.Println("===>[!!!Collector]Failed to response:", err)
						continue
					}

					cancel()
					// }
				}
			}

		case newNotify := <-notifyChan:
			{
				fmt.Println("===>[Collector]New Notify for:", newNotify)

				// forward phase
				ipList := util.GetIPAddress("../ipAddress")

				for _, ipAddress := range ipList {
					// if ipAddress != collector.Address {
					conn, err := grpc.Dial(ipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("===>[!!!Collector]did not connect:", err)
						continue
					}

					nc := notifypb.NewNotifyHandleClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)

					sig := signature.GenerateSig(3, strconv.Itoa(collector.Round), strconv.Itoa(newNotify), collector.LocalTC, collector.Signer)
					aggSig := signature.AggSig(totalLocalSig[newNotify])
					_, err = nc.NotifyReceive(ctx, &notifypb.Notify{Type: 3, Round: strconv.Itoa(collector.Round), Sender: strconv.Itoa(newNotify), Tc: totalLocalTC[newNotify], Sig: sig, From: strconv.Itoa(collector.ID), Aggsig: aggSig, Aggfrom: totalLocalFrom[newNotify]})
					if err != nil {
						fmt.Println("===>[!!!Collector]Failed to response:", err)
						continue
					}

					cancel()
					// }
				}
			}

		}
	}
}
