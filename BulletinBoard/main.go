package main

import (
	"board/RBC/newLeaderpb"
	"board/RBC/newOutputpb"
	"board/RBC/outputpb"
	"board/RBC/pkMsgpb"
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

var roundCurrent int
var viewCurrent int
var newOutputCount map[string]int
var newOutput string
var newLeaderBlame []string
var newOutputBlame []string
var outputChan chan bool

// Reset start a new round
func reset() {
	newOutputCount = make(map[string]int)
	newLeaderBlame = make([]string, 0)
	newOutputBlame = make([]string, 0)
	outputChan = make(chan bool)
}

// new leader blame
type NewLeader struct {
	Round  string
	View   string
	Sender string
	Sig    string
}

// new output blame
type NewOutput struct {
	Round  string
	View   string
	Sender string
	Number string
	Sig    string
}

// pkMsgServer is used to implement pkMsgpb.PkMsgReceive
type pkMsgServer struct {
	pkMsgpb.UnimplementedPkMsgHandleServer
}

// proposalServer is used to implement proposalpb.ProposalReceive
type proposalServer struct {
	proposalpb.UnimplementedProposalHandleServer
}

// newLeaderServer is used to implement newLeaderpb.NewLeaderReceive
type newLeaderServer struct {
	newLeaderpb.UnimplementedNewLeaderHandleServer
}

// outputServer is used to implement outputpb.outputReceive
type outputServer struct {
	outputpb.UnimplementedOutputHandleServer
}

// newOutputServer is used to implement newOutputpb.NewOutputReceive
type newOutputServer struct {
	newOutputpb.UnimplementedNewOutputHandleServer
}

// PkMsgReceive implements pkMsgpb.PkMsgReceive
func (ps *pkMsgServer) PkMsgReceive(ctx context.Context, in *pkMsgpb.PkMsg) (*pkMsgpb.PkMsgResponse, error) {
	fmt.Printf("===>[PkMsgReceive]Received: Sender: %v PK: %v\n", in.GetSender(), in.GetPk())

	result := signature.VerifyNewKeySig(0, in.GetSender(), in.GetPk(), in.GetSig(), in.GetSender())
	fmt.Println("===>[PkMsgReceive]Verify ", result)

	return &pkMsgpb.PkMsgResponse{}, nil
}

// ProposalReceive implements proposalpb.ProposalReceive
func (ps *proposalServer) ProposalReceive(ctx context.Context, in *proposalpb.Proposal) (*proposalpb.ProposalResponse, error) {
	fmt.Printf("===>[ProposalReceive]Received: Round: %v Sender: %v TC: %v\n", in.GetRound(), in.GetSender(), in.GetTc())

	result := signature.VerifyProposalSig(2, in.GetRound(), in.GetView(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
	fmt.Println("===>[ProposalReceive]Verify ", result)

	sender, _ := strconv.Atoi(in.GetSender())
	round, _ := strconv.Atoi(in.GetRound())
	view, _ := strconv.Atoi(in.GetView())
	leaderId := util.GetLeader(round, view)

	if result && sender == leaderId {
		viewCurrent = view
		roundCurrent = view
		config.WriteFileStringArray("/var/www/html/TC", in.GetTc())
	}

	return &proposalpb.ProposalResponse{}, nil
}

// NewLeaderReceive implements newLeaderpb.NewLeaderReceive
func (ns *newLeaderServer) NewLeaderReceive(ctx context.Context, in *newLeaderpb.NewLeader) (*newLeaderpb.NewLeaderResponse, error) {
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

		marshalednewLeader, err := json.Marshal(newLeader)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from NewLeaderReceive]Marshal error : %s", err))
		}

		newLeaderBlame = append(newLeaderBlame, string(marshalednewLeader))
	}

	return &newLeaderpb.NewLeaderResponse{}, nil
}

// outputReceive implements outputpb.outputReceive
func (os *outputServer) OutputReceive(ctx context.Context, in *outputpb.Output) (*outputpb.OutputResponse, error) {
	fmt.Printf("===>[outputReceive]Received: Round: %v From: %v RN: %v\n", in.GetRound(), in.GetSender(), in.GetRandomNumber())
	result := signature.VerifyOutputSig(4, in.GetRound(), in.GetView(), in.GetRandomNumber(), in.GetSender(), in.GetSig(), in.GetSender())
	fmt.Println("===>[outputReceive]Verify ", result)

	if result {
		outputChan <- result
		newOutput = in.GetRandomNumber()
		config.WriteFile("/var/www/html/outputCandidate", in.GetRandomNumber())
	}

	return &outputpb.OutputResponse{}, nil
}

// NewOutputReceive implements newOutputpb.NewOutputReceive
func (ns *newOutputServer) NewOutputReceive(ctx context.Context, in *newOutputpb.NewOutput) (*newOutputpb.NewOutputResponse, error) {
	fmt.Printf("===>[NewOutputReceive]Received: Round: %v Sender: %v\n", in.GetRound(), in.GetSender())
	result := signature.VerifyOutputSig(5, in.GetRound(), in.GetView(), in.GetOutput(), in.GetSender(), in.GetSig(), in.GetSender())
	round, _ := strconv.Atoi(in.GetRound())
	view, _ := strconv.Atoi(in.GetView())
	fmt.Println("===>[NewOutputReceive]Verify ", result)

	if result && round == roundCurrent && view == viewCurrent {

		newOutput := NewOutput{
			Round:  in.GetRound(),
			View:   in.GetView(),
			Sender: in.GetSender(),
			Number: in.GetOutput(),
			Sig:    in.GetSig(),
		}

		marshaledOutput, err := json.Marshal(newOutput)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from NewOutputReceive]Marshal error : %s", err))
		}

		newOutputCount[string(marshaledOutput)]++
		newOutputBlame = append(newOutputBlame, string(marshaledOutput))
	}

	return &newOutputpb.NewOutputResponse{}, nil
}

// bulletin board
type Board struct {
	ProposalTimer *time.Ticker
	OutputTimer   *time.Ticker
}

// initialize a new Board
func NewBoard() *Board {
	c := &Board{
		ProposalTimer: time.NewTicker(4 * time.Second),
		OutputTimer:   time.NewTicker(4 * time.Second),
	}

	c.ProposalTimer.Stop()
	c.OutputTimer.Stop()

	return c
}

func main() {
	// setup phase
	board := NewBoard()
	reset()
	lis, err := net.Listen("tcp", "127.0.0.1:40000")
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from board]Failed to listen: %s", err))
	}
	f := config.GetF()

	// setup server
	ps := grpc.NewServer()
	pkMsgpb.RegisterPkMsgHandleServer(ps, &pkMsgServer{})
	proposalpb.RegisterProposalHandleServer(ps, &proposalServer{})
	newLeaderpb.RegisterNewLeaderHandleServer(ps, &newLeaderServer{})
	outputpb.RegisterOutputHandleServer(ps, &outputServer{})
	newOutputpb.RegisterNewOutputHandleServer(ps, &newOutputServer{})
	go ps.Serve(lis)
	fmt.Printf("===>[Board]board is listening at %v", lis.Addr())

	// init phase
	config.InitGroup()
	time.Sleep(20 * time.Second)
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	rand.Read(b)
	rand_str := hex.EncodeToString(b)
	config.WriteFile("/var/www/html/output", rand_str)
	fmt.Println("\n===>[Board]New output:", rand_str)

	// timer for 4 delta - 2s
	board.ProposalTimer.Reset(18 * time.Second)

	for {
		select {
		// after proposal phase
		case <-board.ProposalTimer.C:
			{
				board.ProposalTimer.Stop()

				// at least one honest node disagrees with the TC set
				if len(newLeaderBlame) > f+1 {
					fmt.Println("===>[Board]publish newLeaderBlame messages")
					config.WriteFileStringArray("/var/www/html/NewLeader", newLeaderBlame)

					// timer for 3 delta - 2s
					board.ProposalTimer.Reset(13 * time.Second)
					reset()
				}
			}

		// leader send a new output candidate
		case <-outputChan:
			{
				// timer for 1 delta - 2s
				board.OutputTimer.Reset(3 * time.Second)
			}

		// after output phase
		case <-board.OutputTimer.C:
			{
				board.OutputTimer.Stop()

				for index, value := range newOutputCount {
					// at least one honest node disagrees with the output candidate
					if value > f+1 {
						newValue := &NewOutput{}
						err = json.Unmarshal([]byte(index), newValue)
						if err != nil {
							fmt.Println("===>[Board]Marshal fail:", err)
						}
						newOutput = newValue.Number

						fmt.Println("===>[Board]publish newOutputBlame messages")
						config.WriteFileStringArray("/var/www/html/NewOutput", newOutputBlame)
					}
				}

				// publish new output
				config.WriteFile("/var/www/html/output", newOutput)
				// timer for 4 delta - 2s
				board.ProposalTimer.Reset(18 * time.Second)
				reset()
			}
		}
	}
}
