package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ipfs/go-log"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/ecdsa/signing"
	"github.com/binance-chain/tss-lib/test"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/json"
)

const (
	IdKey          = "id"
	PartiesKey     = "parties"
	QuorumKey      = "quorum"
	ThresholdKey   = "threshold"
	RoundKey       = "round"
	MnemonicKey    = "mnemonic"
	InputKey       = "input"
	OutputKey      = "output"
	FormatKey      = "format"
	MessageKey     = "message"
	MessageFileKey = "message-file"
	KeyKey         = "key"
	SignatureKey   = "signature"
	PublicKey      = "pub-key"
)

type Message struct {
	From        *tss.PartyID   `json:"from"`
	To          []*tss.PartyID `json:"to"`
	IsBroadcast bool           `json:"is_broadcast"`
	Bytes       []byte         `json:"bytes"`
	Type        string         `json:"type"`
}

var errCh = make(chan *tss.Error, 1)
var outCh = make(chan tss.Message, 1)
var endChG = make(chan keygen.LocalPartySaveData, 1)
var endChS = make(chan *signing.SignatureData, 1)

var Input string
var Output string
var PG *keygen.LocalParty
var PS *signing.LocalParty
var Id int

func main() {
	var rootCmd = &cobra.Command{
		Use:                        "tss-cli",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	rootCmd.AddCommand(Keygen())
	rootCmd.AddCommand(Sign())
	rootCmd.AddCommand(Verify())

	rootCmd.Execute()

	go checkChannels()

	for {
		buf := bufio.NewReader(os.Stdin)
		fmt.Println("Press any key to continue...")
		_, err := buf.ReadByte()

		if err != nil {
			fmt.Println(err)
		} else {
			err = Update()
			fmt.Println(err)
		}
	}
}

func checkChannels() {
outer:
	for {
		select {
		case err := <-errCh:
			fmt.Println("ERR:", err)
		case msg := <-outCh:
			fmt.Println("MSG:", msg)
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
				To:          msg.GetTo(),
				IsBroadcast: msg.IsBroadcast(),
				Bytes:       b,
				Type:        msg.Type(),
			}

			jsonStr, err := json.Marshal(message)

			if err != nil {
				fmt.Println("Marshal error:", err)
			}

			fmt.Println("recipients:", msg.GetTo())

			if len(msg.GetTo()) > 0 {
				for _, recipient := range msg.GetTo() {
					path := "./output/" + Output + "_from_" + msg.GetFrom().String() + "_to_" + recipient.String()
					fmt.Println("Saved file", path)
					err = ioutil.WriteFile(path, jsonStr, 0644)
					if err != nil {
						fmt.Println("Create file error:", err)
					}
				}
			} else {
				path := "./output/" + Output + "_from_" + msg.GetFrom().String() + "_to_all"
				fmt.Println("Saved file", path)
				err = ioutil.WriteFile(path, jsonStr, 0644)
				if err != nil {
					fmt.Println("Create file error:", err)
				}
			}
		case end := <-endChG:
			fmt.Println("END:", end)
			jsonStr, err := json.Marshal(end)
			if err != nil {
				fmt.Println("Marshal error:", err)
			}
			path := "./key/key"
			fmt.Println("Saved file", path)
			err = ioutil.WriteFile(path, jsonStr, 0644)
			if err != nil {
				fmt.Println("Create file error:", err)
			}
			break outer
		case endS := <-endChS:
			fmt.Println("END:", endS)
			break outer
		}

	}
}

func GeneratePartyIDs(count int) tss.SortedPartyIDs {
	ids := make(tss.UnSortedPartyIDs, 0, count)

	for i := 0; i < count; i++ {
		id := fmt.Sprintf("%d", i+1)
		mon := fmt.Sprintf("P[%d]", i+1)
		key, _ := new(big.Int).SetString(id, 10)
		ids = append(ids, tss.NewPartyID(id, mon, key))
	}

	return tss.SortPartyIDs(ids)
}

func PrepareGenParty(threshold, parties int) {
	preParams, _ := keygen.GeneratePreParams(1 * time.Minute)
	pIDs := GeneratePartyIDs(parties)
	p2pCtx := tss.NewPeerContext(pIDs)
	params := tss.NewParameters(p2pCtx, pIDs[Id-1], parties, threshold)

	PG = keygen.NewLocalParty(params, outCh, endChG, *preParams).(*keygen.LocalParty)
}

func PrepareSignParty(msg, key string, parties, quorum int) {
	msgInt := new(big.Int).SetBytes([]byte(msg))
	pIDs := GeneratePartyIDs(parties)
	p2pCtx := tss.NewPeerContext(pIDs)
	keyS, err := LoadKey(key)

	if err != nil {
		panic(err)
	}

	params := tss.NewParameters(p2pCtx, pIDs[Id-1], parties, quorum)

	PS = signing.NewLocalParty(msgInt, params, *keyS, outCh, endChS).(*signing.LocalParty)
}

func Update() error {
	if err := log.SetLogLevel("tss-lib", "debug"); err != nil {
		return err
	}

	updater := test.SharedPartyUpdater

	files, err := ioutil.ReadDir(Input)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		b, err := ioutil.ReadFile(Input + "/" + file.Name())
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

		if Id-1 == msg.GetFrom().Index {
			return errors.New("tried to send a message to itself")
		}

		if PG != nil {
			fmt.Println("Start PG update")
			go updater(PG, msg, errCh)
		} else if PS != nil {
			fmt.Println("Start PS update")
			go updater(PS, msg, errCh)
		}
	}

	return nil
}

func GenerateKey(id, threshold, parties int, input, output string) error {
	if err := log.SetLogLevel("tss-lib", "debug"); err != nil {
		panic(err)
	}

	Id = id
	Input = input
	Output = output
	PrepareGenParty(parties, threshold)

	go func(PG *keygen.LocalParty) {
		if err := PG.Start(); err != nil {
			errCh <- err
		}
	}(PG)

	return nil
}

func LoadKey(keyFile string) (*keygen.LocalPartySaveData, error) {
	b, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Read key file error:", err)
		return nil, err
	}

	var key = new(keygen.LocalPartySaveData)

	if err = json.Unmarshal(b, &key); err != nil {
		fmt.Println("Unmarshal key error:", err)
		return nil, err
	}

	return key, nil
}

func SignMessage(input, output, message, key string, id, parties, quorum int) error {
	if err := log.SetLogLevel("tss-lib", "debug"); err != nil {
		panic(err)
	}
	Id = id
	Input = input
	Output = output
	PrepareSignParty(message, key, parties, quorum)

	go func(PS *signing.LocalParty) {
		if err := PS.Start(); err != nil {
			errCh <- err
		}
	}(PS)

	return nil
}

func VerifyMessage(message, signature, pubkey string) error {
	return nil
}
