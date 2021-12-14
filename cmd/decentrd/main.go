package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/tendermint/spm/cosmoscmd"

	"github.com/Decentr-net/decentr/app"
	"github.com/Decentr-net/decentr/config"
	communitytypes "github.com/Decentr-net/decentr/x/community/types"
	operationstypes "github.com/Decentr-net/decentr/x/operations/types"
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		config.AppName,
		config.AccountAddressPrefix,
		app.DefaultNodeHome,
		config.AppName,
		app.ModuleBasics,
		app.New,
		cosmoscmd.AddSubCmd(
			AddGenesisModeratorsCmd(),
			AddGenesisSupervisorsCmd(),
		),
	)

	// overwrite cosmos migrate cmd
	for _, cmd := range rootCmd.Commands() {
		if strings.HasPrefix(cmd.Use, "migrate") {
			rootCmd.RemoveCommand(cmd)
			break
		}
	}
	rootCmd.AddCommand(MigrateGenesisCmd())

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}

// AddGenesisModeratorsCmd returns add-genesis-community-moderators cobra Command.
func AddGenesisModeratorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-moderator [moderator]",
		Short: "Add a genesis moderator to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			moderator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse moderator: %w", err)
			}

			genFile := server.GetServerContextFromCmd(cmd).Config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			communityGenState := communitytypes.GetGenesisStateFromAppState(cdc, appState)
			communityGenState.Params.Moderators = append(communityGenState.Params.Moderators, moderator)

			communityGenStateBz, err := cdc.MarshalJSON(&communityGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal community genesis state: %w", err)
			}

			appState[communitytypes.ModuleName] = communityGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	return cmd
}

// AddGenesisSupervisorsCmd returns add-genesis-supervisors cobra Command.
func AddGenesisSupervisorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-supervisor [moderator]",
		Short: "Add a genesis supervisor to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			supervisor, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to unmarshal [supervisor]: %w", err)
			}

			genFile := server.GetServerContextFromCmd(cmd).Config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			operationsGenState := operationstypes.GetGenesisStateFromAppState(cdc, appState)
			operationsGenState.Params.Supervisors = append(operationsGenState.Params.Supervisors, supervisor)

			operationsGenStateBz, err := cdc.MarshalJSON(&operationsGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal operations genesis state: %w", err)
			}

			appState[operationstypes.ModuleName] = operationsGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	return cmd
}
