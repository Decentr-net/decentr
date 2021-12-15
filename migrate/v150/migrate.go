package v150

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Decentr-net/decentr/config"
	"github.com/cosmos/cosmos-sdk/client"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	GenesisTime = time.Date(2021, 12, 30, 11, 0, 0, 0, time.UTC)
	ChainID     = "mainnet-2"
)

func Migrate(genDoc *tmtypes.GenesisDoc, ctx client.Context) (*tmtypes.GenesisDoc, error) {
	config.SetAddressPrefixes()

	var appState genutiltypes.AppMap
	var err error
	if err := json.Unmarshal(genDoc.AppState, &appState); err != nil {
		return nil, fmt.Errorf("failed to marchal app state from genesis doc:  %w", err)
	}

	migrateCosmosAppState(appState, ctx)
	migrateDecentrAppState(appState, ctx)

	genDoc.AppState, err = json.Marshal(appState)
	if err != nil {
		return nil, err
	}

	genDoc.GenesisTime = GenesisTime
	genDoc.ChainID = ChainID
	return genDoc, nil
}
