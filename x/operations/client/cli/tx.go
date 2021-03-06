package cli

import (
	"bufio"
	"fmt"
	"github.com/Decentr-net/decentr/x/operations/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	operationsTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"op"},
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	operationsTxCmd.AddCommand(flags.PostCommands(
		GetCmdResetAccount(cdc),
		GetCmdDistributeReward(cdc),
	)...)

	return operationsTxCmd
}

func GetCmdDistributeReward(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "distribute-reward <receiver> <id> <reward>",
		Short: "distribute-reward",
		Aliases: []string{"dr"},
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			receiver, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse receiver: %w", err)
			}

			id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse id: %w", err)
			}

			reward, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse reward: %w", err)
			}

			msg := types.NewMsgDistributeRewards(cliCtx.FromAddress, []types.Reward{{
				Receiver: receiver,
				ID:       id,
				Reward:   reward,
			}})
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdResetAccount is the CLI command for sending a ResetAccount transaction
func GetCmdResetAccount(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reset-account <account owner>",
		Short: "reset account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			accountOwner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse account owner: %w", err)
			}

			msg := types.NewMsgResetAccount(cliCtx.GetFromAddress(), accountOwner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
