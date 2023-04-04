package main

import (
	"board/RBC/newLeaderpb"
	"board/RBC/proposalpb"
	"board/config"
	"board/signature"
	"board/util"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type NewLeader struct {
	Round  string
	View   string
	Sender string
	Sig    string
}

var totalReceiveNewLeader []string

// Reset start a new round
func reset() {
	totalReceiveNewLeader = make([]string, 0)
}

// proposalServer is used to implement proposalpb.ProposalReceive
type proposalServer struct {
	proposalpb.UnimplementedProposalHandleServer
}

// ProposalReceive implements proposalpb.ProposalReceive
func (ps *proposalServer) ProposalReceive(ctx context.Context, in *proposalpb.Proposal) (*proposalpb.ProposalResponse, error) {
	fmt.Printf("===>[ProposalReceive]Received: Round: %v Sender: %v TC: %v\n", in.GetRound(), in.GetSender(), in.GetTc())

	result := signature.VerifySig(2, in.GetRound(), in.GetView(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
	fmt.Println("===>[ProposalReceive]Verify ", result)

	sender, _ := strconv.Atoi(in.GetSender())
	round, _ := strconv.Atoi(in.GetRound())
	view, _ := strconv.Atoi(in.GetView())
	leaderId := util.GetLeader(round, view)

	if result && sender == leaderId {
		config.WriteFileStringArray("/var/www/html/TC", in.GetTc())
	}

	return &proposalpb.ProposalResponse{}, nil
}

// newLeaderServer is used to implement newLeaderpb.NewLeaderReceive
type newLeaderServer struct {
	newLeaderpb.UnimplementedNewLeaderHandleServer
}

// NewLeaderReceive implements newLeaderpb.NewLeaderReceive
func (ps *newLeaderServer) NewLeaderReceive(ctx context.Context, in *newLeaderpb.NewLeader) (*newLeaderpb.NewLeaderResponse, error) {
	fmt.Printf("===>[NewLeaderReceive]Received: Round: %v Sender: %v\n", in.GetRound(), in.GetSender())
	result := signature.VerifyNewLeaderSig(3, in.GetRound(), in.GetView(), in.GetSender(), in.GetSig(), in.GetSender())
	fmt.Println("===>[NewLeaderReceive]Verify ", result)

	if result {
		newLeader := NewLeader{
			Round:  in.GetRound(),
			View:   in.GetView(),
			Sender: in.GetSender(),
			Sig:    in.GetSig(),
		}

		marshaledTC, err := json.Marshal(newLeader)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from NewLeaderReceive]Marshal error : %s", err))
		}

		totalReceiveNewLeader = append(totalReceiveNewLeader, string(marshaledTC))
	}

	return &newLeaderpb.NewLeaderResponse{}, nil
}

type Board struct {
	ProposalTimer *time.Ticker
}

// initialize a new Board
func NewBoard() *Board {
	c := &Board{
		ProposalTimer: time.NewTicker(5 * time.Second),
	}

	c.ProposalTimer.Stop()

	return c
}

func main() {
	board := NewBoard()
	reset()

	lis, err := net.Listen("tcp", "127.0.0.1:40000")
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from board]Failed to listen: %s", err))
	}

	f := config.GetF()

	ps := grpc.NewServer()
	proposalpb.RegisterProposalHandleServer(ps, &proposalServer{})
	newLeaderpb.RegisterNewLeaderHandleServer(ps, &newLeaderServer{})
	go ps.Serve(lis)
	fmt.Printf("===>[Board]board is listening at %v", lis.Addr())

	config.InitGroup()
	time.Sleep(30 * time.Second)
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	rand.Read(b)
	rand_str := hex.EncodeToString(b)

	config.WriteFile("/var/www/html/output", rand_str)

	board.ProposalTimer.Reset(20 * time.Second)

	for {
		select {
		case <-board.ProposalTimer.C:
			{
				board.ProposalTimer.Stop()

				if len(totalReceiveNewLeader) > f+1 {
					fmt.Println("===>[Board]publish newLeader messages")
					config.WriteFileStringArray("/var/www/html/NewLeader", totalReceiveNewLeader)

					board.ProposalTimer.Reset(15 * time.Second)
					reset()
				}
			}
		}
	}

}
