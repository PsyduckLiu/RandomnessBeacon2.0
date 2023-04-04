package main

import (
	"collector/RBC/ackpb"
	"collector/RBC/newLeaderpb"
	"collector/RBC/notifypb"
	"collector/RBC/outputpb"
	"collector/RBC/proposalpb"
	"collector/RBC/submitpb"
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
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	blsSig "go.dedis.ch/dela/crypto/bls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

// submitServer is used to implement submitpb.SubmitReceive
type submitServer struct {
	submitpb.UnimplementedSubmitHandleServer
}

// SubmitReceive implements submitpb.SubmitReceive
func (ps *submitServer) SubmitReceive(ctx context.Context, in *submitpb.Submit) (*submitpb.SubmitResponse, error) {
	fmt.Printf("===>[SubmitReceive]Received: Round: %v Sender: %v TC: %v\n", in.GetRound(), in.GetSender(), in.GetTc())

	result := signature.VerifySig(1, in.GetRound(), in.GetView(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
	fmt.Println("===>[SubmitReceive]Verify ", result)

	sender, _ := strconv.Atoi(in.GetSender())
	if result && totalLocalTC[sender] == nil {
		totalLocalTC[sender] = in.GetTc()
		totalGlobalTC = append(totalGlobalTC, in.GetTc()...)
		fmt.Println(totalLocalTC[sender])
	}

	return &submitpb.SubmitResponse{}, nil
}

// ackServer is used to implement ackpb.ackReceive
type ackServer struct {
	ackpb.UnimplementedAckHandleServer
}

// AckReceive implements ackpb.ackReceive
func (as *ackServer) AckReceive(ctx context.Context, in *ackpb.Ack) (*ackpb.AckResponse, error) {
	f := config.GetF()

	fmt.Printf("\n===>[AckReceive]Received: Round: %v Sender: %v From :%v TC: %v\n", in.GetRound(), in.GetSender(), in.GetFrom(), in.GetTc())
	// result := signature.VerifyMsgSig(2, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetFrom())
	result := signature.VerifySig(1, in.GetRound(), in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
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

	// result := signature.VerifyMsgSig(3, in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetFrom())
	result := signature.VerifySig(1, in.GetRound(), in.GetRound(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
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

			go config.WriteFile("../output", in.GetRandomNumber())
		}
	}

	return &outputpb.OutputResponse{}, nil
}

type Collector struct {
	Round          int
	View           int
	ID             int
	Address        string
	LocalTC        []string
	GlobalTC       []string
	OutputCh       chan string
	Signer         blsSig.Signer
	ReceiveTimer   *time.Ticker
	SubmitTimer    *time.Ticker
	ProposalTimer  *time.Ticker
	NewLeaderTimer *time.Ticker
	RBCTimer       *time.Ticker
}

// initialize a new Collector
func NewCollector(id int) *Collector {
	c := &Collector{
		Round:          0,
		View:           0,
		ID:             id,
		Address:        "127.0.0.1:" + strconv.Itoa(30000+id),
		LocalTC:        make([]string, 0),
		GlobalTC:       make([]string, 0),
		OutputCh:       make(chan string),
		Signer:         blsSig.NewSigner(),
		ReceiveTimer:   time.NewTicker(5 * time.Second),
		SubmitTimer:    time.NewTicker(5 * time.Second),
		ProposalTimer:  time.NewTicker(10 * time.Second),
		NewLeaderTimer: time.NewTicker(5 * time.Second),
		RBCTimer:       time.NewTicker(5 * time.Second),
	}

	c.ReceiveTimer.Stop()
	c.SubmitTimer.Stop()
	c.ProposalTimer.Stop()
	c.NewLeaderTimer.Stop()
	c.RBCTimer.Stop()

	config.DownloadFile("http://172.18.208.214/Key.yml", "download/Key.yml")
	config.DownloadFile("http://172.18.208.214/IP.yml", "download/IP.yml")
	config.DownloadFile("http://172.18.208.214/Config.yml", "download/Config.yml")
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
	submitpb.RegisterSubmitHandleServer(ps, &submitServer{})
	// ackpb.RegisterAckHandleServer(ps, &ackServer{})
	// notifypb.RegisterNotifyHandleServer(ps, &notifyServer{})
	// outputpb.RegisterOutputHandleServer(ps, &outputServer{})
	go ps.Serve(lis)
	fmt.Printf("===>[Collector]Collector is listening at %v", lis.Addr())

	go watch.WatchOutput(collector.OutputCh, "output")

	for {
		select {
		case <-collector.OutputCh:
			{
				collector.ReceiveTimer.Reset(5 * time.Second)
				reset()
				collector.Round++
				collector.View = 0
				collector.LocalTC = collector.LocalTC[0:0]
				collector.GlobalTC = collector.GlobalTC[0:0]
				fmt.Println("\n===>[Collector]Round:", collector.Round)
			}

		case <-collector.ReceiveTimer.C:
			{
				collector.ReceiveTimer.Stop()

				collector.LocalTC = totalReceiveTC
				fmt.Println("===>[Collector]New TC is:", collector.LocalTC)
				if len(collector.LocalTC) == 0 {
					continue
				}

				// submit phase
				ipList := config.GetPeerIP()
				leaderId := util.GetLeader(collector.Round, collector.View)
				leaderIp := ipList[leaderId]
				fmt.Println("===>[Collector]IP list is:", ipList)
				fmt.Println("===>[Collector]Leader IP is:", ipList[leaderId])

				conn, err := grpc.Dial(leaderIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					fmt.Println("===>[!!!Collector]did not connect:", err)
					continue
				}

				sc := submitpb.NewSubmitHandleClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				sig := signature.GenerateSig(1, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.LocalTC, collector.Signer)
				_, err = sc.SubmitReceive(ctx, &submitpb.Submit{Type: 1, Round: strconv.Itoa(collector.Round), View: strconv.Itoa(collector.View), Sender: id, Tc: collector.LocalTC, Sig: sig})
				if err != nil {
					fmt.Println("Send to", leaderIp)
					fmt.Println("===>[!!!Collector]Failed to response:", err)
					continue
				}

				if collector.ID == leaderId {
					collector.SubmitTimer.Reset(5 * time.Second)
				}
				collector.ProposalTimer.Reset(10 * time.Second)
			}

		case <-collector.SubmitTimer.C:
			{
				collector.SubmitTimer.Stop()

				collector.GlobalTC = totalGlobalTC
				fmt.Println("===>[Leader][Collector]New total TC is:", collector.GlobalTC)
				if len(collector.GlobalTC) == 0 {
					continue
				}

				// proposal phase
				boardIp := config.GetBoardIP()
				fmt.Println("===>[Collector]Board IP is:", boardIp)

				conn, err := grpc.Dial(boardIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					fmt.Println("===>[!!!Collector]did not connect:", err)
					continue
				}

				pc := proposalpb.NewProposalHandleClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				// simulate evil leader
				behaviour := rand.Intn(100)
				if behaviour > 99 {
					sig := signature.GenerateSig(2, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.GlobalTC, collector.Signer)
					_, err = pc.ProposalReceive(ctx, &proposalpb.Proposal{Type: 2, Round: strconv.Itoa(collector.Round), View: strconv.Itoa(collector.View), Sender: id, Tc: collector.GlobalTC, Sig: sig})
					if err != nil {
						fmt.Println("Send to", boardIp)
						fmt.Println("===>[!!!Collector]Failed to response:", err)
						continue
					}
				} else {
					fmt.Println("I am an adversary")
				}
			}
		case <-collector.ProposalTimer.C:
			{
				collector.ProposalTimer.Stop()

				config.DownloadFile("http://172.18.208.214/TC", "download/TC")
				tcSet := config.ReadFileStringArray("download/TC")
				subset := util.IsSubSet(collector.LocalTC, tcSet)
				if !subset {
					fmt.Println("===>[Collector]Message on Bulletin Board is wrong! Reboot")
					// TODO

					// send new-leader blame
					boardIp := config.GetBoardIP()
					fmt.Println("===>[Collector]Board IP is:", boardIp)

					conn, err := grpc.Dial(boardIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("===>[!!!Collector]did not connect:", err)
						continue
					}

					nc := newLeaderpb.NewNewLeaderHandleClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					sig := signature.GenerateNewLeaderSig(3, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.Signer)
					_, err = nc.NewLeaderReceive(ctx, &newLeaderpb.NewLeader{Type: 3, Round: strconv.Itoa(collector.Round), View: strconv.Itoa(collector.View), Sender: id, Sig: sig})
					if err != nil {
						fmt.Println("Send to", boardIp)
						fmt.Println("===>[!!!Collector]Failed to response:", err)
						continue
					}

					// timer
				} else {
					fmt.Println("===>[Collector]Message on Bulletin Board is correct! Continue")
				}
			}
		case newProposal := <-ackChan:
			{
				fmt.Println("===>[Collector]New proposal for:", newProposal)

				// forward phase
				ipList := config.GetPeerIP()

				for _, ipAddress := range ipList {
					// if ipAddress != collector.Address {
					conn, err := grpc.Dial(ipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("===>[!!!Collector]did not connect:", err)
						continue
					}

					ac := ackpb.NewAckHandleClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)

					// sig := signature.GenerateMsgSig(2, strconv.Itoa(collector.Round), strconv.Itoa(newProposal), totalLocalTC[newProposal], collector.Signer)
					sig := signature.GenerateSig(1, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.LocalTC, collector.Signer)
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
				ipList := config.GetPeerIP()

				for _, ipAddress := range ipList {
					// if ipAddress != collector.Address {
					conn, err := grpc.Dial(ipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("===>[!!!Collector]did not connect:", err)
						continue
					}

					nc := notifypb.NewNotifyHandleClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)

					// sig := signature.GenerateMsgSig(3, strconv.Itoa(collector.Round), strconv.Itoa(newNotify), totalLocalTC[newNotify], collector.Signer)
					sig := signature.GenerateSig(1, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.LocalTC, collector.Signer)
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
				ipList := config.GetPeerIP()

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
