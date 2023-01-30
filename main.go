package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/test"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ipfs/go-log"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/json"
)

const (
	IdKey        = "id"
	PartiesKey   = "parties"
	ThresholdKey = "threshold"
	RoundKey     = "round"
	MnemonicKey  = "mnemonic"
	InputKey     = "input"
	OutputKey    = "output"
	FormatKey    = "format"
	MessageKey   = "message"
	KeyKey       = "key"
	SignatureKey = "signature"
	PublicKey    = "pub-key"
)

type Message struct {
	From        *tss.PartyID `json:"from"`
	IsBroadcast bool         `json:"is_broadcast"`
	Bytes       []byte       `json:"bytes"`
}

var outCh chan tss.Message
var endCh chan keygen.LocalPartySaveData
var errCh chan *tss.Error
var Output string
var P *keygen.LocalParty
var Id int

func main() {
	var rootCmd = &cobra.Command{
		Use:                        "tss-cli",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	rootCmd.AddCommand(Privgen())
	rootCmd.AddCommand(Pubgen())
	rootCmd.AddCommand(Sign())
	rootCmd.AddCommand(Verify())

	rootCmd.Execute()

	go MonitorEnd()

	for {
		buf := bufio.NewReader(os.Stdin)
		fmt.Println("Press any key to continue...")
		_, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			err = GeneratePrivateKeyUpdate()
			fmt.Println(err)
		}
	}
}

func MonitorEnd() {
	for {
		select {
		case err := <-errCh:
			fmt.Println("ERR:", err)
			return
		case end := <-endCh:
			fmt.Println("END:", end)
			return
		case msg := <-outCh:
			dest := msg.GetTo()

			if dest != nil && dest[0].Index == msg.GetFrom().Index {
				return
			}

			b, _, err := msg.WireBytes()

			if err != nil {
				fmt.Println("GetBytes error:", err)
			}

			message := Message{
				From:        msg.GetFrom(),
				IsBroadcast: msg.IsBroadcast(),
				Bytes:       b,
			}

			jsonStr, err := json.Marshal(message)

			if err != nil {
				fmt.Println("Marshal error:", err)
			}

			fmt.Println("File", Output)
			err = ioutil.WriteFile(Output, jsonStr, 0644)
			if err != nil {
				fmt.Println("Create file error:", err)
			}
			return
		}
	}
}

func PrepareParty(threshold, parties int) {
	preParams, _ := keygen.GeneratePreParams(1 * time.Minute)
	pIDs := tss.GenerateTestPartyIDs(parties)
	p2pCtx := tss.NewPeerContext(pIDs)
	params := tss.NewParameters(tss.S256(), p2pCtx, pIDs[Id-1], parties, threshold)

	errCh = make(chan *tss.Error, 4)
	outCh = make(chan tss.Message, 4)
	endCh = make(chan keygen.LocalPartySaveData, 4)

	P = keygen.NewLocalParty(params, outCh, endCh, *preParams).(*keygen.LocalParty)
}

func GeneratePrivateKeyUpdate() error {
	if err := log.SetLogLevel("tss-lib", "debug"); err != nil {
		return err
	}

	updater := test.SharedPartyUpdater
	fmt.Println("1")
	b, err := ioutil.ReadFile(Output)
	if err != nil {
		fmt.Println("Read file error:", err)
		return err
	}

	var message Message

	if err = json.Unmarshal(b, &message); err != nil {
		fmt.Println("Unmarshal error:", err)
		return err
	}

	msg, err := tss.ParseWireMessage(message.Bytes, message.From, message.IsBroadcast)
	if err != nil {
		return err
	}

	dest := msg.GetTo()
	if dest == nil {
		if Id == msg.GetFrom().Index {
			return nil
		}
		go updater(P, msg, errCh)
	}

	return nil
}

func GeneratePrivateKey(id, threshold, parties, round int, mnemonic, output string) error {
	if err := log.SetLogLevel("tss-lib", "debug"); err != nil {
		panic(err)
	}

	Id = id
	Output = output
	PrepareParty(parties, threshold)

	if err := P.Start(); err != nil {
		fmt.Println("Start error:", err)
		return err
	}

	return nil
}

func GeneratePublicKey(input, output, format string) error {
	return nil
}

func SignMessage(input, output, message, key string, round int) error {
	return nil
}

func VerifyMessage(message, signature, pubkey string) error {
	return nil
}
