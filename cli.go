package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Privgen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privgen",
		Short: "Generate private keys",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			id, err := cmd.Flags().GetString(IdKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			threshold, err := cmd.Flags().GetInt(ThresholdKey)
			if err != nil {
				return fmt.Errorf("invalid threshold: %w", err)
			}

			parties, err := cmd.Flags().GetInt(PartiesKey)
			if err != nil {
				return fmt.Errorf("invalid partiesy: %w", err)
			}

			round, err := cmd.Flags().GetInt(RoundKey)
			if err != nil {
				return fmt.Errorf("invalid partiesy: %w", err)
			}

			mnemonic, err := cmd.Flags().GetString(MnemonicKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			output, err := cmd.Flags().GetString(OutputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			return GeneratePrivateKey(id, threshold, parties, round, mnemonic, output)
		},
	}

	cmd.Flags().String("id", "", "id")
	cmd.MarkFlagRequired(IdKey)

	cmd.Flags().Int(ThresholdKey, 0, "threshold")
	cmd.MarkFlagRequired(ThresholdKey)

	cmd.Flags().Int(PartiesKey, 0, "parties")
	cmd.MarkFlagRequired(PartiesKey)

	cmd.Flags().Int(RoundKey, 0, "round")
	cmd.MarkFlagRequired(RoundKey)

	cmd.Flags().String(MnemonicKey, "", "mnemonic")
	cmd.MarkFlagRequired(MnemonicKey)

	cmd.Flags().String(OutputKey, "./", "output")

	return cmd
}

func Pubgen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pubgen",
		Short: "Generate public keys",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := cmd.Flags().GetString(InputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			output, err := cmd.Flags().GetString(OutputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			format, err := cmd.Flags().GetString(FormatKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			return GeneratePublicKey(input, output, format)
		},
	}

	cmd.Flags().String(InputKey, "", "mnemonic")
	cmd.MarkFlagRequired(InputKey)

	cmd.Flags().String(OutputKey, "./", "output")

	cmd.Flags().String(FormatKey, "string", "output")

	return cmd
}

func Sign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign message",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := cmd.Flags().GetString(InputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			output, err := cmd.Flags().GetString(OutputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			round, err := cmd.Flags().GetInt(RoundKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			message, err := cmd.Flags().GetString(MessageKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			key, err := cmd.Flags().GetString(KeyKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			return SignMessage(input, output, message, key, round)
		},
	}

	cmd.Flags().String(InputKey, "", "mnemonic")
	cmd.MarkFlagRequired(InputKey)

	cmd.Flags().String(OutputKey, "./", "output")

	cmd.Flags().Int(RoundKey, 0, "round")
	cmd.MarkFlagRequired(RoundKey)

	cmd.Flags().String(MessageKey, "", "message string/file")
	cmd.MarkFlagRequired(MessageKey)

	cmd.Flags().String(KeyKey, "", "key file")
	cmd.MarkFlagRequired(KeyKey)

	return cmd
}

func Verify() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign message",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			message, err := cmd.Flags().GetString(MessageKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			signature, err := cmd.Flags().GetString(SignatureKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			pubkey, err := cmd.Flags().GetString(PublicKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			return VerifyMessage(message, signature, pubkey)
		},
	}

	cmd.Flags().String(MessageKey, "", "message string/file")
	cmd.MarkFlagRequired(MessageKey)

	cmd.Flags().String(SignatureKey, "", "signature")
	cmd.MarkFlagRequired(SignatureKey)

	cmd.Flags().String(PublicKey, "", "pubkey string/file")
	cmd.MarkFlagRequired(PublicKey)

	return cmd
}
