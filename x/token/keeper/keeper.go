package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	operations "github.com/Decentr-net/decentr/x/operations/types"
	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/Decentr-net/decentr/x/utils"
)

var totalSupplyKey = []byte("totalSupply")

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	storeKey      sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSubspace params.Subspace

	accountKeeper      auth.AccountKeeper
	distributionKeeper distribution.Keeper
}

// NewKeeper creates new instances of the token Keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace params.Subspace,
	accountKeeper auth.AccountKeeper,
	distributionKeeper distribution.Keeper,
) Keeper {
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		paramSubspace:      ps,
		accountKeeper:      accountKeeper,
		distributionKeeper: distributionKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// AddTokens adds token to the given owner
func (k Keeper) AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int) {
	if k.IsBanned(ctx, owner) {
		return
	}

	var balance sdk.Int

	if k.IsInitialBalanceSet(ctx, owner) {
		balance = k.GetBalance(ctx, owner).Add(amount)
	} else {
		k.Logger(ctx).Info(
			"account initialized with init balance",
			"account", owner.String())
		balance = utils.InitialTokenBalance().Add(amount)
	}

	k.SetBalance(ctx, owner, balance)
	k.IncBalanceDelta(ctx, owner, amount)
}

// addTotalSupply increase or decrease total supply with the given amount of tokens
func (k Keeper) addTotalSupply(ctx sdk.Context, amount sdk.Int) {
	balance := k.GetTotalSupply(ctx)
	balance = balance.Add(amount)
	ctx.KVStore(k.storeKey).Set(totalSupplyKey, k.cdc.MustMarshalBinaryBare(balance))
}

// GetBalance returns token balance for the given address
func (k Keeper) GetBalance(ctx sdk.Context, address sdk.AccAddress) sdk.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	if store.Has(address) {
		balance := store.Get(address)
		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryBare(balance, &amount)
		return amount
	}

	// if account exists, but initial token is not set
	if k.accountKeeper.GetAccount(ctx, address) != nil {
		return utils.InitialTokenBalance()
	}

	return sdk.ZeroInt()
}

func (k Keeper) IsInitialBalanceSet(ctx sdk.Context, owner sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	return store.Has(owner)
}

// SetBalance set balance to the given user
func (k Keeper) SetBalance(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	store.Set(owner, k.cdc.MustMarshalBinaryBare(amount))
	k.addTotalSupply(ctx, amount)
}

func (k Keeper) GetBalanceDelta(ctx sdk.Context, address sdk.AccAddress) sdk.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeltaPrefix)
	delta := sdk.ZeroInt()

	if store.Has(address) {
		k.cdc.MustUnmarshalBinaryBare(store.Get(address), &delta)
	}

	return delta
}

func (k Keeper) IncBalanceDelta(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeltaPrefix)

	inc := func(address sdk.AccAddress) {
		delta := k.GetBalanceDelta(ctx, address)
		store.Set(address, k.cdc.MustMarshalBinaryBare(delta.Add(amount)))
	}

	inc(address)
	inc(types.AccumulatedDelta)
}

func (k Keeper) DistributeRewards(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeltaPrefix)

	total := k.GetBalanceDelta(ctx, types.AccumulatedDelta)
	if total.IsNil() || total.IsZero() {
		return
	}

	part := k.distributionKeeper.GetFeePoolCommunityCoins(ctx).AmountOf(operations.DefaultDenom).QuoInt(total)
	if part.IsNil() || part.IsZero() {
		return
	}

	it := store.Iterator(nil, nil)
	defer it.Close()

	for it.Valid() {
		address := sdk.AccAddress(it.Key())

		delta := k.GetBalanceDelta(ctx, address)

		store.Delete(it.Key())
		it.Next()

		if types.AccumulatedDelta.Equals(address) {
			continue
		}

		amount := part.MulInt(delta).TruncateInt()
		if amount.IsNil() || amount.IsZero() {
			continue
		}

		coins := sdk.Coins{
			{
				Denom:  operations.DefaultDenom,
				Amount: amount,
			},
		}

		k.AddRewardDistribution(ctx, address, ctx.BlockHeight(), coins)
		if err := k.distributionKeeper.DistributeFromFeePool(ctx, coins, address); err != nil {
			panic(fmt.Errorf("failed to do accrual: %w", err))
		}
	}
}

// GetBalanceIterator gets an iterator over all balances in which the keys are the accounts and the values are their balance
func (k Keeper) GetBalanceIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	return sdk.KVStorePrefixIterator(store, nil)
}

// GetDeltasIterator gets an iterator over all balance deltas
func (k Keeper) GetDeltasIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeltaPrefix)
	return sdk.KVStorePrefixIterator(store, nil)
}

// GetRewardsDistributionIterator gets an iterator over all accounts' history
func (k Keeper) GetRewardsDistributionIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RewardsPrefix)
	return sdk.KVStorePrefixIterator(store, nil)
}

// GetTotalSupply returns total token supply
func (k Keeper) GetTotalSupply(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	totalSupply := store.Get(totalSupplyKey)
	if totalSupply == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.Int
	k.cdc.MustUnmarshalBinaryBare(totalSupply, &amount)
	return amount
}

func (k *Keeper) AddRewardDistribution(ctx sdk.Context, address sdk.AccAddress, height int64, coins sdk.Coins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RewardsPrefix)

	dd := k.GetRewardsDistributionHistory(ctx, address)
	b := k.cdc.MustMarshalBinaryBare(append(dd, types.RewardDistribution{
		Height: height,
		Coins:  coins,
	}))

	store.Set(address, b)
}

func (k *Keeper) GetRewardsDistributionHistory(ctx sdk.Context, address sdk.AccAddress) []types.RewardDistribution {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RewardsPrefix)

	if store.Has(address) {
		var dd []types.RewardDistribution
		k.cdc.MustUnmarshalBinaryBare(store.Get(address), &dd)
		return dd
	}

	return nil
}

// SetBan bans/unbans address
func (k Keeper) SetBan(ctx sdk.Context, address sdk.AccAddress, ban bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BanListPrefix)
	if ban {
		k.IncBalanceDelta(ctx, address, k.GetBalanceDelta(ctx, address).MulRaw(-1))
		store.Set(address, []byte{0x00})
	} else {
		store.Delete(address)
	}
}

// IsBanned returns is address banned
func (k Keeper) IsBanned(ctx sdk.Context, address sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BanListPrefix)

	return store.Has(address)
}

func (k Keeper) ResetAccount(ctx sdk.Context, addr sdk.AccAddress) {
	k.SetBalance(ctx, addr, utils.InitialTokenBalance())
	k.IncBalanceDelta(ctx, addr, k.GetBalanceDelta(ctx, addr).MulRaw(-1))
}
