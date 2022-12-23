package main

import (
	"crypto/elliptic"
	"fmt"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"math/big"
	"time"
)

const (
	IdKey        = "id"
	PartiesKey   = "parties"
	ThresholdKey = "threshold"
	RoundKey     = "round"
	MnemonicKey  = "mnemonic"
	InputKey     = "input-dir"
	OutputKey    = "output-dir"
	FormatKey    = "format"
	MessageKey   = "message"
	KeyKey       = "key"
	SignatureKey = "signature"
	PublicKey    = "pub-key"
)

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
}

func GeneratePrivateKey(id string, threshold, parties, round int, mnemonic, output string) error {
	preParams, _ := keygen.GeneratePreParams(1 * time.Minute)
	pIDs := tss.GenerateTestPartyIDs(parties)
	p2pCtx := tss.NewPeerContext(pIDs)
	outCh := make(chan tss.Message, parties)
	endCh := make(chan keygen.LocalPartySaveData, parties)
	curve := elliptic.P256()
	pid := tss.NewPartyID(id, id, big.NewInt(0))
	params := tss.NewParameters(curve, p2pCtx, pid, parties, threshold)
	P := keygen.NewLocalParty(params, outCh, endCh, *preParams)

	go func() {
		err := P.Start()
		fmt.Println(err)
	}()

	fmt.Printf("outCh: %v+", outCh)
	fmt.Printf("endCh: %v+", endCh)

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
