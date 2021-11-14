package types

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	. "github.com/Decentr-net/decentr/testutil"
)

func TestGenesisState_Validate(t *testing.T) {
	addr := NewAccAddress()
	postUUID := uuid.Must(uuid.NewV1())
	p := DefaultParams()

	for _, tc := range []struct {
		desc     string
		genState *GenesisState
		valid    bool
	}{
		{
			desc:     "default",
			genState: DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid",
			genState: &GenesisState{
				Params: &p,
				Posts: []Post{
					{
						Uuid:         postUUID.String(),
						Owner:        addr,
						Title:        "Title",
						PreviewImage: "",
						Category:     0,
						Text:         "Fifteen symbols should be typed",
					},
				},
				Likes: []Like{
					{
						Owner:     NewAccAddress(),
						PostOwner: addr,
						PostUuid:  postUUID.String(),
						Weight:    0,
					},
				},
			},
			valid: true,
		},
		{
			desc: "self-like",
			genState: &GenesisState{
				Params: &p,
				Posts: []Post{
					{
						Uuid:         postUUID.String(),
						Owner:        addr,
						Title:        "Title",
						PreviewImage: "",
						Category:     0,
						Text:         "Fifteen symbols should be typed",
					},
				},
				Likes: []Like{
					{
						Owner:     addr,
						PostOwner: addr,
						PostUuid:  postUUID.String(),
						Weight:    0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "unknown post",
			genState: &GenesisState{
				Params: &p,
				Posts: []Post{
					{
						Uuid:         postUUID.String(),
						Owner:        addr,
						Title:        "Title",
						PreviewImage: "",
						Category:     0,
						Text:         "Fifteen symbols should be typed",
					},
				},
				Likes: []Like{
					{
						Owner:     addr,
						PostOwner: NewAccAddress(),
						PostUuid:  postUUID.String(),
						Weight:    0,
					},
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
