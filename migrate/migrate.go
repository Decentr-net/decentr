package migrate

// this file is adapted version of github.com/cosmos/cosmos-sdk/x/genutil/client/cli/migrate.go

import (
	"fmt"

	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	v150 "github.com/Decentr-net/decentr/migrate/v150"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate v1.5.0 /path/to/genesis.json --chain-id=mainnet-2 --genesis-time=2021-11-22T17:00:00Z
`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			importGenesis := args[0]

			oldGenDoc, err := tmtypes.GenesisDocFromFile(importGenesis)
			if err != nil {
				return fmt.Errorf("failed to read genesis document from file %s: %w", importGenesis, err)
			}

			newGenDoc, err := v150.Migrate(oldGenDoc, clientCtx)
			if err != nil {
				return fmt.Errorf("failed to run migration: %w", err)
			}

			if err := newGenDoc.ValidateAndComplete(); err != nil {
				return fmt.Errorf("failed to validate: %w", err)
			}

			bz, err := tmjson.Marshal(newGenDoc)
			if err != nil {
				return fmt.Errorf("failed to marshal genesis doc: %w", err)
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return fmt.Errorf("failed to sort JSON genesis doc: %w", err)
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}
}
