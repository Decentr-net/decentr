package types

import (
	"fmt"
	"net/url"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"
)

const (
	maxTitleLength = 150
	minTitleLength = 1
	maxPostLength  = 64 * 1000
	minPostLength  = 15
	maxURLLength   = 4 * 1024
)

func (p Post) Address() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Uuid)
}

func (p Post) Validate() error {
	if err := sdk.VerifyAddressFormat(p.Owner); err != nil {
		return fmt.Errorf("invalid owner: %w", err)
	}

	if _, err := uuid.FromString(p.Uuid); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	if p.Category < Category_CATEGORY_UNDEFINED || p.Category > Category_CATEGORY_SPORTS {
		return fmt.Errorf("invalid category: %d", p.Category)
	}

	if utf8.RuneCountInString(p.Title) > maxTitleLength || len(p.Title) < minTitleLength {
		return fmt.Errorf("invalid title: title length should be in [%d;%d]", minTitleLength, maxTitleLength)
	}

	if !isPreviewImageValid(p.PreviewImage) {
		return fmt.Errorf("invalid preview_image")
	}

	if utf8.RuneCountInString(p.Text) < minPostLength || len(p.Text) > maxPostLength {
		return fmt.Errorf("invalid text: text length should be in [%d;%d]", minPostLength, maxPostLength)
	}

	return nil
}

func (l Like) Validate() error {
	if err := sdk.VerifyAddressFormat(l.Owner); err != nil {
		return fmt.Errorf("invalid owner: %w", err)
	}

	if err := sdk.VerifyAddressFormat(l.PostOwner); err != nil {
		return fmt.Errorf("invalid post_owner: %w", err)
	}

	if l.Owner.Equals(l.PostOwner) {
		return fmt.Errorf("invalid post_owner: self-like")
	}

	if _, err := uuid.FromString(l.PostUuid); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	if l.Weight < LikeWeight_LIKE_WEIGHT_DOWN || l.Weight > LikeWeight_LIKE_WEIGHT_UP {
		return fmt.Errorf("invalid weight")
	}

	return nil
}

func ValidateFollowers(who string, whom []sdk.AccAddress) error {
	if _, err := sdk.AccAddressFromBech32(who); err != nil {
		return fmt.Errorf("invalid follower: %w", err)
	}

	if len(whom) == 0 {
		return fmt.Errorf("invalid followee: empty array")
	}

	for i, v := range whom {
		if err := sdk.VerifyAddressFormat(v); err != nil {
			return fmt.Errorf("invalid followee #%d: %w", i+1, err)
		}
	}

	return nil
}

func isPreviewImageValid(str string) bool {
	if len(str) > maxURLLength {
		return false
	}

	if str == "" {
		return true
	}

	url, err := url.Parse(str)
	if err != nil {
		return false
	}

	if len(url.Host) == 0 {
		return false
	}

	return url.Scheme == "http" || url.Scheme == "https"
}
