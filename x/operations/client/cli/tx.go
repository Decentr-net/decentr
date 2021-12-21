package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/operations/types"
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewDistributeRewardsCmd(),
		NewResetAccountCmd(),
		NewMintCmd(),
		NewBurnCmd(),
	)

	return txCmd
}

func NewDistributeRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distribute-reward [receiver] [reward]",
		Short: "Add [reward] to [receiver]",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			receiver, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid receiver: %w", err)
			}

			reward, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return fmt.Errorf("invalid reward: %w", err)
			}

			msg := types.NewMsgDistributeRewards(clientCtx.GetFromAddress(), []types.Reward{
				{
					Receiver: receiver.String(),
					Reward:   sdk.DecProto{Dec: reward},
				},
			})

			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewResetAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-account [address]",
		Short: "Remove all account's activity and sets pdv to initial value by [address]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address: %w", err)
			}

			msg := types.NewMsgResetAccount(clientCtx.GetFromAddress(), address)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [coin]",
		Short: "Mint coin to the account in the '--from' flag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid coin: %w", err)
			}

			msg := types.NewMsgMint(clientCtx.GetFromAddress(), coin)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewBurnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [coin]",
		Short: "Burn coin from the account in the '--from' flag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid coin: %w", err)
			}

			msg := types.NewMsgBurn(clientCtx.GetFromAddress(), coin)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
