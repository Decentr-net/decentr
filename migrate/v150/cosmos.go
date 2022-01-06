package v150

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	v040genutil "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v040"
	v043genutil "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v043"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	v039slashing "github.com/cosmos/cosmos-sdk/x/slashing/legacy/v039"
	ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
	"github.com/cosmos/ibc-go/modules/core/exported"
	ibctypes "github.com/cosmos/ibc-go/modules/core/types"
)

func patchSlashingParams(appState genutiltypes.AppMap) genutiltypes.AppMap {
	v039Codec := codec.NewLegacyAmino()

	if appState[v039slashing.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var slashingGenState v039slashing.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039slashing.ModuleName], &slashingGenState)

		slashingGenState.Params.SignedBlocksWindow = 600

		appState[v039slashing.ModuleName] = v039Codec.MustMarshalJSON(slashingGenState)
	}

	return appState
}

func migrateCosmosAppState(appState genutiltypes.AppMap, clientCtx client.Context) genutiltypes.AppMap {
	appState = patchSlashingParams(appState)

	appState = v040genutil.Migrate(appState, clientCtx)
	appState = v043genutil.Migrate(appState, clientCtx)

	ibcTransferGenesis := ibctransfertypes.DefaultGenesisState()
	ibcGenesis := ibctypes.DefaultGenesisState()
	capabilityGenesis := capabilitytypes.DefaultGenesis()
	evidenceGenesis := evidencetypes.DefaultGenesisState()

	ibcTransferGenesis.Params.ReceiveEnabled = true
	ibcTransferGenesis.Params.SendEnabled = true
	ibcGenesis.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}

	v040Codec := clientCtx.Codec
	appState[ibctransfertypes.ModuleName] = v040Codec.MustMarshalJSON(ibcTransferGenesis)
	appState[ibchost.ModuleName] = v040Codec.MustMarshalJSON(ibcGenesis)
	appState[capabilitytypes.ModuleName] = v040Codec.MustMarshalJSON(capabilityGenesis)
	appState[evidencetypes.ModuleName] = v040Codec.MustMarshalJSON(evidenceGenesis)

	return appState
}