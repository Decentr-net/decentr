package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	p := DefaultParams()
	return &GenesisState{
		Params:    &p,
		Posts:     []Post{},
		Likes:     []Like{},
		Following: map[string]GenesisState_AddressList{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	pm := map[string]struct{}{}
	for _, v := range gs.Posts {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("invalid post %s/%s: %w", v.Owner, v.Uuid, err)
		}
		pm[v.Address()] = struct{}{}
	}

	for _, v := range gs.Likes {
		p := Post{Owner: v.PostOwner, Uuid: v.PostUuid}.Address()
		if _, ok := pm[p]; !ok {
			return fmt.Errorf("invalid like %s by %s: post not found", p, v.Owner)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("invalid like %s by %s: %w", p, v.Owner, err)
		}
	}

	for who, whom := range gs.Following {
		if err := ValidateFollowers(who, whom.Address); err != nil {
			return fmt.Errorf("invalid following for %s: %w", who, err)
		}
	}

	return nil
}
