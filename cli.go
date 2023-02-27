package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Keygen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate keys",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			id, err := cmd.Flags().GetInt(IdKey)
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

			input, err := cmd.Flags().GetString(InputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			output, err := cmd.Flags().GetString(OutputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			return GenerateKey(id, threshold, parties, input, output)
		},
	}

	cmd.Flags().Int(IdKey, 0, "id")
	cmd.MarkFlagRequired(IdKey)

	cmd.Flags().Int(ThresholdKey, 0, "threshold")
	cmd.MarkFlagRequired(ThresholdKey)

	cmd.Flags().Int(PartiesKey, 0, "parties")
	cmd.MarkFlagRequired(PartiesKey)

	cmd.Flags().String(InputKey, "./input", "input")
	cmd.Flags().String(OutputKey, "./output", "output")

	return cmd
}

func Sign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign message",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cmd.Flags().GetInt(IdKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			input, err := cmd.Flags().GetString(InputKey)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			parties, err := cmd.Flags().GetInt(PartiesKey)
			if err != nil {
				return fmt.Errorf("invalid partiesy: %w", err)
			}

			quorum, err := cmd.Flags().GetInt(QuorumKey)
			if err != nil {
				return fmt.Errorf("invalid quorum: %w", err)
			}

			output, err := cmd.Flags().GetString(OutputKey)
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

			return SignMessage(input, output, message, key, id, parties, quorum)
		},
	}

	cmd.Flags().Int(IdKey, 0, "id")
	cmd.MarkFlagRequired(IdKey)

	cmd.Flags().Int(PartiesKey, 0, "parties")
	cmd.MarkFlagRequired(PartiesKey)

	cmd.Flags().Int(QuorumKey, 0, "quorum")
	cmd.MarkFlagRequired(QuorumKey)

	cmd.Flags().String(InputKey, "./input", "input")
	cmd.Flags().String(OutputKey, "./output", "output")

	cmd.Flags().String(MessageKey, "", "message")
	cmd.Flags().String(MessageFileKey, "./message/message", "message-file")
	cmd.MarkFlagsMutuallyExclusive(MessageKey, MessageFileKey)

	cmd.Flags().String(KeyKey, "./key/key", "key file")

	return cmd
}

func Verify() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify message",
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
