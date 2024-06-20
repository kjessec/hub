package antefilters

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageFilterDecorator(t *testing.T) {
	var passFilter MessageFilter = func(ctx sdk.Context, m sdk.Msg, simulate bool) error {
		return nil
	}
	var failFilter MessageFilter = func(ctx sdk.Context, m sdk.Msg, simulate bool) error {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fail")
	}

	t.Run("all pass", func(t *testing.T) {
		ante := sdk.ChainAnteDecorators(NewMessageFilterDecorator(passFilter))

		_, err := ante(sdk.Context{}, , false)
		assert.NoError(t, err)

	})

	t.Run("early exit", func(t *testing.T) {
		ante := sdk.ChainAnteDecorators(NewMessageFilterDecorator(passFilter, failFilter, passFilter))

		_, err := ante(sdk.Context{}, nil, false)
		assert.Error(t, err)
	})
}
