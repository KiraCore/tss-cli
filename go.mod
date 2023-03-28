module github.com/KiraCore/tss-cli

go 1.18

require (
	github.com/binance-chain/tss-lib v0.0.0-00010101000000-000000000000
	github.com/bnb-chain/tss-lib v1.3.5
	github.com/btcsuite/btcd v0.22.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.0
	github.com/tendermint/tendermint v0.34.20
	go.uber.org/zap v1.24.0
	golang.org/x/text v0.4.0
)

require (
	github.com/agl/ed25519 v0.0.0-20200225211852-fd4d107ace12 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 // indirect
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce // indirect
	github.com/decred/dcrd/dcrec/edwards/v2 v2.0.0 // indirect
	github.com/gogo/protobuf v1.3.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/ipfs/go-log v1.0.4 // indirect
	github.com/ipfs/go-log/v2 v2.1.1 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/otiai10/primes v0.0.0-20180210170552-f6d2a1ba97c4 // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.2.0 // indirect
	google.golang.org/protobuf v1.28.2-0.20220831092852-f930b1dc76e8 // indirect
)

replace (
	github.com/agl/ed25519 => github.com/binance-chain/edwards25519 v0.0.0-20200305024217-f36fc4b53d43
	//github.com/anyswap/FastMulThreshold-DSA => ./libs/smpc-lib
	//github.com/binance-chain/tss-lib => ./libs/thor-lib
	//github.com/bnb-chain/tss-lib => ./libs/bnb-lib
	//github.com/bnb-chain/tss-lib => ./libs/bnb-lib
	//github.com/binance-chain/tss-lib => gitlab.com/thorchain/tss/tss-lib v0.1.3
	//github.com/binance-chain/tss-lib => github.com/bnb-chain/tss-lib v1.3.5
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)
