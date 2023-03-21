package main

import (
	"context"
	"fmt"
	"generator/config"
	"generator/crypto/timedCommitment"
	"generator/tcMsgpb"
	"generator/util"
	"generator/watch"
	"math/big"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Generator struct {
	ID       int
	Address  string
	OutputCh chan string
}

// initialize a new Collector
func NewGenerator(id int) *Generator {
	g := &Generator{
		ID:       id,
		Address:  "127.0.0.1:" + strconv.Itoa(3333+id),
		OutputCh: make(chan string),
	}

	return g
}

func main() {
	id := os.Args[1]
	idInt, _ := strconv.Atoi(id)
	generator := NewGenerator(idInt)

	index := new(big.Int)
	f := config.GetF()

	go watch.WatchOutput(generator.OutputCh, "../output")

	for {
		select {
		case <-generator.OutputCh:
			{
				time.Sleep(2 * time.Second)
				maskedMsg, h, M_k, a1, a2, z := timedCommitment.GenerateTC()
				fmt.Println(timedCommitment.VerifyTC(maskedMsg, h, M_k, a1, a2, z))

				// send TC
				ipList := util.GetIPAddress("../ipAddress")
				index.Mod(maskedMsg, big.NewInt(int64(3*f+1)))
				fmt.Println(index)
				ipAddress := ipList[index.Int64()]
				// ipAddress := ipList[idInt]

				conn, err := grpc.Dial(ipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					fmt.Println("===>[!!!Generator]did not connect:", err)
					continue
				}

				tc := tcMsgpb.NewTcMsgHandleClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)

				_, err = tc.TcMsgReceive(ctx, &tcMsgpb.TcMsg{MaskedMsg: maskedMsg.String(), HA: h.GetA().Int64(), HB: h.GetB().Int64(), HC: h.GetC().Int64(), MkA: M_k.GetA().Int64(), MkB: M_k.GetB().Int64(), MkC: M_k.GetC().Int64(), A1A: a1.GetA().Int64(), A1B: a1.GetB().Int64(), A1C: a1.GetC().Int64(), A2A: a2.GetA().Int64(), A2B: a2.GetB().Int64(), A2C: a2.GetC().Int64(), Z: z.String()})
				if err != nil {
					fmt.Println("Send to", ipAddress)
					fmt.Println("===>[!!!Collector]Failed to response:", err)
					continue
				}

				cancel()
			}
		}
	}

	// timedCommitment.ForcedOpen(maskedMsg, h)

	// for {
	// 	rand.Seed(time.Now().UnixNano())
	// 	b := make([]byte, 6)
	// 	rand.Read(b)
	// 	rand_str := hex.EncodeToString(b)

	// 	watch.WriteFile("../output", rand_str)

	// 	time.Sleep(5 * time.Second)
	// }
}
