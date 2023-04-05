package main

import (
	"collector/RBC/newLeaderpb"
	"collector/RBC/outputpb"
	"collector/RBC/pkMsgpb"
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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	blsSig "go.dedis.ch/dela/crypto/bls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var totalReceiveTC []string
var totalLocalTC map[int][]string
var totalGlobalTC []string

// var totalLocalSig map[int][]string
// var totalLocalFrom map[int][]string
// var totalOutputCount map[string]int
// var totalOutputSig map[string][]string
// var totalOutputFrom map[string][]string
// var ackChan chan int
// var notifyChan chan int

// Reset start a new round
func reset() {
	totalReceiveTC = make([]string, 0)
	totalLocalTC = make(map[int][]string)
	totalGlobalTC = make([]string, 0)
}

type NewLeader struct {
	Round  string
	View   string
	Sender string
	Sig    string
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

// submitServer is used to implement submitpb.SubmitReceive
type submitServer struct {
	submitpb.UnimplementedSubmitHandleServer
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

// SubmitReceive implements submitpb.SubmitReceive
func (ps *submitServer) SubmitReceive(ctx context.Context, in *submitpb.Submit) (*submitpb.SubmitResponse, error) {
	fmt.Printf("===>[SubmitReceive]Received: Round: %v Sender: %v TC: %v\n", in.GetRound(), in.GetSender(), in.GetTc())

	result := signature.VerifySubmitSig(1, in.GetRound(), in.GetView(), in.GetSender(), in.GetTc(), in.GetSig(), in.GetSender())
	fmt.Println("===>[SubmitReceive]Verify ", result)

	sender, _ := strconv.Atoi(in.GetSender())
	if result && totalLocalTC[sender] == nil {
		totalLocalTC[sender] = in.GetTc()
		totalGlobalTC = append(totalGlobalTC, in.GetTc()...)
		fmt.Println(totalLocalTC[sender])
	}

	return &submitpb.SubmitResponse{}, nil
}

type Collector struct {
	Round          int
	View           int
	ID             int
	NewLeaderVote  int
	Address        string
	RandomNumber   string
	LocalTC        []string
	GlobalTC       []string
	OutputCh       chan string
	Signer         blsSig.Signer
	ReceiveTimer   *time.Ticker
	SubmitTimer    *time.Ticker
	ProposalTimer  *time.Ticker
	NewLeaderTimer *time.Ticker
	NewOutputTimer *time.Ticker
	RBCTimer       *time.Ticker
}

// initialize a new Collector
func NewCollector(id int) *Collector {
	c := &Collector{
		Round:          0,
		View:           0,
		ID:             id,
		NewLeaderVote:  0,
		Address:        "127.0.0.1:" + strconv.Itoa(30000+id),
		RandomNumber:   "",
		LocalTC:        make([]string, 0),
		GlobalTC:       make([]string, 0),
		OutputCh:       make(chan string),
		Signer:         blsSig.NewSigner(),
		ReceiveTimer:   time.NewTicker(5 * time.Second),
		SubmitTimer:    time.NewTicker(5 * time.Second),
		ProposalTimer:  time.NewTicker(10 * time.Second),
		NewLeaderTimer: time.NewTicker(5 * time.Second),
		NewOutputTimer: time.NewTicker(5 * time.Second),
		RBCTimer:       time.NewTicker(5 * time.Second),
	}

	c.ReceiveTimer.Stop()
	c.SubmitTimer.Stop()
	c.ProposalTimer.Stop()
	c.NewLeaderTimer.Stop()
	c.NewOutputTimer.Stop()
	c.RBCTimer.Stop()

	boardIp := config.GetBoardLisIP()
	biKey, _ := c.Signer.GetPublicKey().MarshalBinary()

	conn, err := grpc.Dial(boardIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("===>[!!!Collector]did not connect: %s", err))
	}

	pc := pkMsgpb.NewPkMsgHandleClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sig := signature.GenerateNewKeySig(0, strconv.Itoa(c.ID), hex.EncodeToString(biKey), c.Signer)
	_, err = pc.PkMsgReceive(ctx, &pkMsgpb.PkMsg{Type: 0, Sender: strconv.Itoa(c.ID), Pk: hex.EncodeToString(biKey), Sig: sig})
	if err != nil {
		fmt.Println("Send to", boardIp)
		fmt.Println("===>[!!!Collector]Failed to response:", err)
	}

	return c
}

func main() {
	id := os.Args[1]
	idInt, _ := strconv.Atoi(id)
	collector := NewCollector(idInt)

	go func() {
		time.Sleep(5 * time.Second)
		boardIP := config.GetBoardIP()
		filename1 := "http://" + boardIP + "/Key.yml"
		config.DownloadFile(filename1, "download/Key.yml")
		filename2 := "http://" + boardIP + "/IP.yml"
		config.DownloadFile(filename2, "download/IP.yml")
		filename3 := "http://" + boardIP + "/Config.yml"
		config.DownloadFile(filename3, "download/Config.yml")
	}()

	lis, err := net.Listen("tcp", collector.Address)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from Collector]Failed to listen: %s", err))
	}

	ps := grpc.NewServer()
	tcMsgpb.RegisterTcMsgHandleServer(ps, &tcMsgServer{})
	submitpb.RegisterSubmitHandleServer(ps, &submitServer{})
	go ps.Serve(lis)
	fmt.Printf("===>[Collector]Collector is listening at %v", lis.Addr())

	go watch.WatchOutput(collector.OutputCh, "output")

	for {
		select {
		// new output
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

		// after receive phase
		case <-collector.ReceiveTimer.C:
			{
				collector.ReceiveTimer.Stop()

				collector.NewLeaderVote = 0
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

				sig := signature.GenerateSubmitSig(1, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.LocalTC, collector.Signer)
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

		// after submit phase
		case <-collector.SubmitTimer.C:
			{
				collector.SubmitTimer.Stop()

				collector.GlobalTC = totalGlobalTC
				fmt.Println("===>[Leader][Collector]New total TC is:", collector.GlobalTC)
				if len(collector.GlobalTC) == 0 {
					continue
				}

				// proposal phase
				boardIp := config.GetBoardLisIP()
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
				if behaviour > 19 {
					sig := signature.GenerateSubmitSig(2, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.GlobalTC, collector.Signer)
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

		// after proposal phase
		case <-collector.ProposalTimer.C:
			{
				collector.ProposalTimer.Stop()

				boardIP := config.GetBoardIP()
				filename1 := "http://" + boardIP + "/TC"
				config.DownloadFile(filename1, "download/TC")

				tcSet := config.ReadFileStringArray("download/TC")
				subset := util.IsSubSet(collector.LocalTC, tcSet)
				if !subset {
					fmt.Println("===>[Collector]Message on Bulletin Board is wrong! Reboot")
					// TODO

					// send new-leader blame
					boardIp := config.GetBoardLisIP()
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
					collector.GlobalTC = tcSet
					fmt.Println("===>[Collector]Message on Bulletin Board is correct! Continue")
				}

				collector.NewLeaderTimer.Reset(5 * time.Second)
			}

		// after new leader blame phase
		case <-collector.NewLeaderTimer.C:
			{
				collector.NewLeaderTimer.Stop()

				f := config.GetF()
				boardIP := config.GetBoardIP()
				filename1 := "http://" + boardIP + "/NewLeader"
				config.DownloadFile(filename1, "download/NewLeader")

				blameVote := config.ReadFileStringArray("download/NewLeader")

				for _, vote := range blameVote {
					newVote := &NewLeader{}
					err = json.Unmarshal([]byte(vote), newVote)
					if err != nil {
						fmt.Println("===>[!!!Collector]Failed to response:", err)
					}

					if newVote.Round == strconv.Itoa(collector.Round) && newVote.View == strconv.Itoa(collector.View) && newVote.View == strconv.Itoa(collector.View) {
						result := signature.VerifyNewLeaderSig(3, newVote.Round, newVote.View, newVote.Sender, newVote.Sig, newVote.Sender)
						fmt.Println("===>[Collector]New Leader VoteVerify ", result)

						if result {
							collector.NewLeaderVote++
						}
					}
				}

				if collector.NewLeaderVote > f+1 {
					collector.View++

					// re-submit phase
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

					sig := signature.GenerateSubmitSig(1, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), strconv.Itoa(collector.ID), collector.LocalTC, collector.Signer)
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
				} else {
					fmt.Println("===>[Collector]Golabel TC set is:", collector.GlobalTC)
					// force-open
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

						collector.NewOutputTimer.Reset(5 * time.Second)
					}
					collector.RandomNumber = randomNumber.String()

					leaderId := util.GetLeader(collector.Round, collector.View)
					if collector.ID == leaderId {
						// output phase
						boardIp := config.GetBoardLisIP()
						fmt.Println("===>[Collector]Board IP is:", boardIp)

						conn, err := grpc.Dial(boardIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
						if err != nil {
							fmt.Println("===>[!!!Collector]did not connect:", err)
							continue
						}

						oc := outputpb.NewOutputHandleClient(conn)
						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()

						sig := signature.GenerateOutputSig(4, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), randomNumber.String(), strconv.Itoa(collector.ID), collector.Signer)
						_, err = oc.OutputReceive(ctx, &outputpb.Output{Type: 4, Round: strconv.Itoa(collector.Round), View: strconv.Itoa(collector.View), RandomNumber: randomNumber.String(), Sender: id, Sig: sig})
						if err != nil {
							fmt.Println("Send to", boardIp)
							fmt.Println("===>[!!!Collector]Failed to response:", err)
							continue
						}
					}
				}
			}

		// after new output blame phase
		case <-collector.NewOutputTimer.C:
			{
				collector.NewOutputTimer.Stop()

				boardIP := config.GetBoardIP()
				filename1 := "http://" + boardIP + "/outputCandidate"
				config.DownloadFile(filename1, "download/outputCandidate")

				outputFromBoard := config.ReadFile("download/outputCandidate")
				if outputFromBoard != collector.RandomNumber {
					boardIp := config.GetBoardLisIP()
					fmt.Println("===>[Collector]Board IP is:", boardIp)

					conn, err := grpc.Dial(boardIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("===>[!!!Collector]did not connect:", err)
						continue
					}

					oc := outputpb.NewOutputHandleClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					sig := signature.GenerateOutputSig(5, strconv.Itoa(collector.Round), strconv.Itoa(collector.View), collector.RandomNumber, strconv.Itoa(collector.ID), collector.Signer)
					_, err = oc.OutputReceive(ctx, &outputpb.Output{Type: 5, Round: strconv.Itoa(collector.Round), View: strconv.Itoa(collector.View), RandomNumber: collector.RandomNumber, Sender: id, Sig: sig})
					if err != nil {
						fmt.Println("Send to", boardIp)
						fmt.Println("===>[!!!Collector]Failed to response:", err)
						continue
					}
				}
			}
		}
	}
}
