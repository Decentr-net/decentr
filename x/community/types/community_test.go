package types

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	. "github.com/Decentr-net/decentr/testutil"
)

func TestPost_validatePost(t *testing.T) {
	validPost := Post{
		Uuid:         uuid.Must(uuid.NewV1()).String(),
		Owner:        NewAccAddress().String(),
		Title:        "test",
		PreviewImage: "https://tf-decentr-cerberusd-bucket-mainnet.s3.amazonaws.com/decentr1cehgxxyskvn39f9sj46jqjewp805x4prad6ds5/ce5b7147-211e-4f35-9732-ec4c200daf75",
		Category:     8,
		Text:         "I just need to type 15 symbols",
	}

	tt := []struct {
		name      string
		alterFunc func(p *Post)
		valid     bool
	}{
		{
			name:  "valid",
			valid: true,
		},
		{
			name: "invalid uuid",
			alterFunc: func(p *Post) {
				p.Uuid = ""
			},
			valid: false,
		},
		{
			name: "invalid uuid #2",
			alterFunc: func(p *Post) {
				p.Uuid = "123"
			},
			valid: false,
		},
		{
			name: "invalid owner",
			alterFunc: func(p *Post) {
				p.Owner = ""
			},
			valid: false,
		},
		{
			name: "invalid owner #2",
			alterFunc: func(p *Post) {
				p.Owner = "123"
			},
			valid: false,
		},
		{
			name: "invalid text",
			alterFunc: func(p *Post) {
				p.Text = "123"
			},
			valid: false,
		},
		{
			name: "invalid text #2",
			alterFunc: func(p *Post) {
				for i := 0; i < 64*1000; i++ {
					p.Text = p.Text + "s"
				}
			},
			valid: false,
		},
		{
			name: "invalid title",
			alterFunc: func(p *Post) {
				p.Title = ""
			},
			valid: false,
		},
		{
			name: "invalid title #2",
			alterFunc: func(p *Post) {
				for i := 0; i < 150; i++ {
					p.Title = p.Title + "s"
				}
			},
			valid: false,
		},
		{
			name: "invalid category",
			alterFunc: func(p *Post) {
				p.Category = -1
			},
			valid: false,
		},
		{
			name: "invalid category #2",
			alterFunc: func(p *Post) {
				p.Category = 10
			},
			valid: false,
		},
		{
			name: "invalid preview_image",
			alterFunc: func(p *Post) {
				p.PreviewImage = "https://"
			},
			valid: false,
		},
		{
			name: "invalid preview_image #2",
			alterFunc: func(p *Post) {
				p.PreviewImage = "decentr.xyz"
			},
			valid: false,
		},
		{
			name: "invalid preview_image #3",
			alterFunc: func(p *Post) {
				p.PreviewImage = "ftp://decentr.xyz/file.png"
			},
			valid: false,
		},
		{
			name: "empty post",
			alterFunc: func(p *Post) {
				*p = Post{}
			},
			valid: false,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := validPost
			if tc.alterFunc != nil {
				tc.alterFunc(&p)
			}

			if err := p.Validate(); tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestLike_Validate(t *testing.T) {
	validLike := Like{
		Owner:     NewAccAddress().String(),
		PostUuid:  uuid.Must(uuid.NewV1()).String(),
		PostOwner: NewAccAddress().String(),
		Weight:    LikeWeight_LIKE_WEIGHT_UP,
	}

	tt := []struct {
		name      string
		alterFunc func(l *Like)
		valid     bool
	}{
		{
			name:  "valid",
			valid: true,
		},
		{
			name: "empty owner",
			alterFunc: func(l *Like) {
				l.Owner = ""
			},
			valid: false,
		},
		{
			name: "invalid owner",
			alterFunc: func(l *Like) {
				l.Owner = "123"
			},
			valid: false,
		},
		{
			name: "empty post_uuid",
			alterFunc: func(l *Like) {
				l.PostUuid = ""
			},
			valid: false,
		},
		{
			name: "invalid post_uuid",
			alterFunc: func(l *Like) {
				l.PostUuid = "123"
			},
			valid: false,
		},
		{
			name: "empty post_owner",
			alterFunc: func(l *Like) {
				l.PostOwner = ""
			},
			valid: false,
		},
		{
			name: "invalid post_owner",
			alterFunc: func(l *Like) {
				l.PostOwner = "123"
			},
			valid: false,
		},
		{
			name: "invalid weight",
			alterFunc: func(l *Like) {
				l.Weight = -2
			},
			valid: false,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := validLike
			if tc.alterFunc != nil {
				tc.alterFunc(&l)
			}

			if err := l.Validate(); tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_ValidateFollowers(t *testing.T) {
	tt := []struct {
		name  string
		who   string
		whom  []string
		valid bool
	}{
		{
			name:  "valid",
			who:   NewAccAddress().String(),
			whom:  []string{NewAccAddress().String(), NewAccAddress().String()},
			valid: true,
		},
		{
			name:  "empty who",
			who:   "",
			whom:  []string{NewAccAddress().String(), NewAccAddress().String()},
			valid: false,
		},
		{
			name:  "invalid who",
			who:   "123",
			whom:  []string{NewAccAddress().String(), NewAccAddress().String()},
			valid: false,
		},
		{
			name:  "invalid_who",
			who:   NewAccAddress().String(),
			whom:  []string{NewAccAddress().String(), "123"},
			valid: false,
		},
		{
			name:  "invalid_who",
			who:   NewAccAddress().String(),
			whom:  nil,
			valid: false,
		},
		{
			name:  "empty",
			who:   "",
			whom:  nil,
			valid: false,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := ValidateFollowers(tc.who, tc.whom); tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
