package main

import (
	"collector/RBC/ackpb"
	"collector/RBC/notifypb"
	"collector/RBC/outputpb"
	"collector/RBC/proposalpb"
	"collector/RBC/tcMsgpb"
	"collector/config"
	"collector/crypto/binaryquadraticform"
	"collector/crypto/timedCommitment"
	"collector/signature"
	"collector/util"
	"collector/watch"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	blsSig "go.dedis.ch/dela/crypto/bls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// var totalReceiveTC = make([]string, 0)
// var totalLocalTC = make(map[int][]string)
// var totalLocalSig = make(map[int][]string)
// var totalLocalFrom = make(map[int][]string)
// var totalGlobalTC = make([]string, 0)
// var totalOutputCount = make(map[string]int)
// var totalOutputSig = make(map[string][]string)
// var totalOutputFrom = make(map[string][]string)
// var ackChan = make(chan int, 5)
// var notifyChan = make(chan int, 5)

var totalReceiveTC []string
var totalLocalTC map[int][]string
var totalLocalSig map[int][]string
var totalLocalFrom map[int][]string
var totalGlobalTC []string
var totalOutputCount map[string]int
var totalOutputSig map[string][]string
var totalOutputFrom map[string][]string
var ackChan chan int
var notifyChan chan int

// Reset start a new round
func reset() {
	totalReceiveTC = make([]string, 0)

	totalLocalTC = make(map[int][]string)
	totalGlobalTC = make([]string, 0)

	totalLocalSig = make(map[int][]string)
	totalLocalFrom = make(map[int][]string)

	totalOutputCount = make(map[string]int)
	totalOutputSig = make(map[string][]string)
	totalOutputFrom = make(map[string][]string)

	ackChan = make(chan int, 5)
	notifyChan = make(chan int, 5)
}

type TC struct {
	MaskedMsg string
	HA        string
	HB        string
	HC        string
}

// tcMsgServer is used to implement tcMsgpb.TcMsgReceive
type tcMsgServer struct {
	tcMsgpb.UnimplementedTcMsgHandleServer
}

// tcMsgReceive implements tcMsgpb.TcMsgReceive
func (ts *tcMsgServer) TcMsgReceive(ctx context.Context, in *tcMsgpb.TcMsg) (*tcMsgpb.TcMsgResponse, error) {
	maskedMsg := new(big.Int)
	z := new(big.Int)
	maskedMsg.SetString(in.GetMaskedMsg(), 10)
	z.SetString(in.GetZ(), 10)

	h, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(in.GetHA()), big.NewInt(in.GetHB()), big.NewInt(in.GetHC()))
	fmt.Printf("===>[TcMsgReceive]The group element h is (a=%v,b=%v,c=%v,d=%v)\n", h.GetA(), h.GetB(), h.GetC(), h.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from TcMsgReceive]Generate new BQuadratic Form failed: %s", err))
	}
	M_k, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(in.GetMkA()), big.NewInt(in.GetMkB()), big.NewInt(in.GetMkC()))
	fmt.Printf("===>[TcMsgReceive]The group element M_K is (a=%v,b=%v,c=%v,d=%v)\n", M_k.GetA(), M_k.GetB(), M_k.GetC(), M_k.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from TcMsgReceive]Generate new BQuadratic Form failed: %s", err))
	}
	a1, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(in.GetA1A()), big.NewInt(in.GetA1B()), big.NewInt(in.GetA1C()))
	fmt.Printf("===>[TcMsgReceive]The group element a1 is (a=%v,b=%v,c=%v,d=%v)\n", a1.GetA(), a1.GetB(), a1.GetC(), a1.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from TcMsgReceive]Generate new BQuadratic Form failed: %s", err))
	}
	a2, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(in.GetA2A()), big.NewInt(in.GetA2B()), big.NewInt(in.GetA2C()))
	fmt.Printf("===>[TcMsgReceive]The group element a2 is (a=%v,b=%v,c=%v,d=%v)\n", a2.GetA(), a2.GetB(), a2.GetC(), a2.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from TcMsgReceive]Generate new BQuadratic Form failed: %s", err))
	}

	result := timedCommitment.VerifyTC(maskedMsg, h, M_k, a1, a2, z)
	if result {
		fmt.Println("===>[TcMsgReceive]new tc pass!!!")

		newTC := TC{
			MaskedMsg: maskedMsg.String(),
			HA:        h.GetA().String(),
			HB:        h.GetB().String(),
			HC:        h.GetC().String(),
		}

		marshaledTC, err := json.Marshal(newTC)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from TcMsgReceive]Marshal error : %s", err))
		}

		totalReceiveTC = append(totalReceiveTC, string(marshaledTC))
	}

	return &tcMsgpb.TcMsgResponse{}, nil
}

// proposalServer is used to implement proposalpb.ProposalReceive
type proposalServer struct {
	proposalpb.UnimplementedProposalHandleServer
}

// ProposalReceive implements proposalpb.ProposalReceive
func (ps *proposalServer) ProposalReceive(ctx context.Context, in *proposalpb.Proposal) (*proposalpb.ProposalResponse, error) {
	fmt.Printf("===>[ProposalReceive]Received: Round: %v Sender: %v TC: %v\n", in.GetRound(), in.GetSender(), in.GetTc())

	result := signature.VerifyMsgSig(1, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
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
	result := signature.VerifyMsgSig(2, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetFrom())
	fmt.Println("===>[AckReceive]Verify ", result)

	ackTC := in.GetTc()
	sender, _ := strconv.Atoi(in.GetSender())
	if reflect.DeepEqual(ackTC, totalLocalTC[sender]) {
		totalLocalSig[sender] = append(totalLocalSig[sender], in.GetSig())
		totalLocalFrom[sender] = append(totalLocalFrom[sender], in.GetFrom())
	}

	if len(totalLocalSig[sender]) == 2*f+1 && totalLocalTC[sender] != nil {
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
func (ns *notifyServer) NotifyReceive(ctx context.Context, in *notifypb.Notify) (*notifypb.NotifyResponse, error) {
	fmt.Printf("===>[notifyReceive]Received: Round: %v Sender: %v From :%v TC: %v\n", in.GetRound(), in.GetSender(), in.GetFrom(), in.GetTc())

	result := signature.VerifyMsgSig(3, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetFrom())
	fmt.Println("===>[notifyReceive]Verify ", result)

	if signature.VerifyAggMsgSig(in.GetRound(), in.GetSender(), in.GetTc(), in.GetAggsig(), in.GetAggfrom()) {
		totalGlobalTC = append(totalGlobalTC, in.GetTc()...)
	} else {
		fmt.Println("===>[notifyReceive]Agg Verify error")
	}

	return &notifypb.NotifyResponse{}, nil
}

// outputServer is used to implement outputpb.outputReceive
type outputServer struct {
	outputpb.UnimplementedOutputHandleServer
}

// outputReceive implements outputpb.outputReceive
func (os *outputServer) OutputReceive(ctx context.Context, in *outputpb.Output) (*outputpb.OutputResponse, error) {
	f := config.GetF()
	var aggSig string

	fmt.Printf("===>[outputReceive]Received: Round: %v From: %v From :%v RN: %v\n", in.GetRound(), in.GetFrom(), in.GetFrom(), in.GetRandomNumber())
	result := signature.VerifyOutputSig(4, in.GetRound(), in.GetRandomNumber(), in.GetSig(), in.GetFrom())
	fmt.Println("===>[outputReceive]Verify ", result)

	if result {
		if totalOutputSig[in.GetRandomNumber()] == nil {
			totalOutputCount[in.GetRandomNumber()] = 1
		} else {
			totalOutputCount[in.GetRandomNumber()]++
		}

		totalOutputFrom[in.GetRandomNumber()] = append(totalOutputFrom[in.GetRandomNumber()], in.GetFrom())
		totalOutputSig[in.GetRandomNumber()] = append(totalOutputSig[in.GetRandomNumber()], in.GetSig())

		if totalOutputCount[in.GetRandomNumber()] == 2*f+1 {
			aggSig = signature.AggSig(totalOutputSig[in.GetRandomNumber()])
			check := signature.VerifyAggOutputSig(in.GetRound(), in.GetRandomNumber(), aggSig, totalOutputFrom[in.GetRandomNumber()])
			fmt.Println("===>[outputReceive]Agg Verify", check)

			go watch.WriteFile("../output", in.GetRandomNumber())
		}
	}

	return &outputpb.OutputResponse{}, nil
}

type Collector struct {
	Round        int
	ID           int
	Address      string
	LocalTC      []string
	GlobalTC     []string
	OutputCh     chan string
	Signer       blsSig.Signer
	ReceiveTimer *time.Ticker
	RBCTimer     *time.Ticker
}

// initialize a new Collector
func NewCollector(id int) *Collector {
	c := &Collector{
		Round:        0,
		ID:           id,
		Address:      "127.0.0.1:" + strconv.Itoa(30000+id),
		LocalTC:      make([]string, 0),
		GlobalTC:     make([]string, 0),
		OutputCh:     make(chan string),
		Signer:       blsSig.NewSigner(),
		ReceiveTimer: time.NewTicker(5 * time.Second),
		RBCTimer:     time.NewTicker(5 * time.Second),
	}

	c.ReceiveTimer.Stop()
	c.RBCTimer.Stop()

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
	tcMsgpb.RegisterTcMsgHandleServer(ps, &tcMsgServer{})
	proposalpb.RegisterProposalHandleServer(ps, &proposalServer{})
	ackpb.RegisterAckHandleServer(ps, &ackServer{})
	notifypb.RegisterNotifyHandleServer(ps, &notifyServer{})
	outputpb.RegisterOutputHandleServer(ps, &outputServer{})
	go ps.Serve(lis)
	fmt.Printf("===>[Collector]Collector is listening at %v", lis.Addr())

	go watch.WatchOutput(collector.OutputCh, "../output")

	for {
		select {
		case <-collector.OutputCh:
			{
				collector.ReceiveTimer.Reset(10 * time.Second)
				reset()
				collector.Round++
				collector.LocalTC = collector.LocalTC[0:0]
				collector.GlobalTC = collector.GlobalTC[0:0]
				fmt.Println("\n===>[Collector]Round:", collector.Round)
			}

		case <-collector.ReceiveTimer.C:
			{
				collector.ReceiveTimer.Stop()
				collector.RBCTimer.Reset(10 * time.Second)

				collector.LocalTC = totalReceiveTC
				fmt.Println("===>[Collector]New TC is:", collector.LocalTC)
				if len(collector.LocalTC) == 0 {
					continue
				}

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

					sig := signature.GenerateMsgSig(1, strconv.Itoa(collector.Round), strconv.Itoa(collector.ID), collector.LocalTC, collector.Signer)
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

					sig := signature.GenerateMsgSig(2, strconv.Itoa(collector.Round), strconv.Itoa(newProposal), totalLocalTC[newProposal], collector.Signer)
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

					sig := signature.GenerateMsgSig(3, strconv.Itoa(collector.Round), strconv.Itoa(newNotify), totalLocalTC[newNotify], collector.Signer)
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

		case <-collector.RBCTimer.C:
			{
				collector.RBCTimer.Stop()
				// collector.GlobalTC = totalGlobalTC
				fmt.Println("===>[Collector]Total TC set is:", totalGlobalTC)
				collector.GlobalTC = util.RemoveRepeatElement(totalGlobalTC)
				fmt.Println("===>[Collector]Reduced Total TC set is:", collector.GlobalTC)

				// TODO: forceopen
				// randomNumber := strconv.Itoa(collector.Round)
				randomNumber := big.NewInt(0)
				for _, tc := range collector.GlobalTC {
					newTC := &TC{}
					err = json.Unmarshal([]byte(tc), newTC)
					if err != nil {
						fmt.Println("===>[!!!Collector]Failed to response:", err)
					}

					maskedMsg := new(big.Int)
					HA := new(big.Int)
					HB := new(big.Int)
					HC := new(big.Int)
					maskedMsg.SetString(newTC.MaskedMsg, 10)
					HA.SetString(newTC.HA, 10)
					HB.SetString(newTC.HB, 10)
					HC.SetString(newTC.HC, 10)

					h, err := binaryquadraticform.NewBQuadraticForm(HA, HB, HC)
					fmt.Printf("===>[Collector]The group element h is (a=%v,b=%v,c=%v,d=%v)\n", h.GetA(), h.GetB(), h.GetC(), h.GetDiscriminant())
					if err != nil {
						panic(fmt.Errorf("===>[ERROR from Collector]Generate new BQuadratic Form failed: %s", err))
					}

					openMsgString := timedCommitment.ForcedOpen(maskedMsg, h)
					openMsg := new(big.Int)
					openMsg.SetString(openMsgString, 10)

					randomNumber.Xor(randomNumber, openMsg)
				}

				// output phase
				ipList := util.GetIPAddress("../ipAddress")

				conn, err := grpc.Dial(ipList[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					fmt.Println("===>[!!!Collector]did not connect:", err)
					continue
				}

				oc := outputpb.NewOutputHandleClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)

				sig := signature.GenerateOutputSig(4, strconv.Itoa(collector.Round), randomNumber.String(), collector.Signer)
				_, err = oc.OutputReceive(ctx, &outputpb.Output{Type: 4, Round: strconv.Itoa(collector.Round), RandomNumber: randomNumber.String(), Sig: sig, From: strconv.Itoa(collector.ID)})
				if err != nil {
					fmt.Println("===>[!!!Collector]Failed to response:", err)
					continue
				}

				cancel()
			}

		}
	}
}
