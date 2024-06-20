package antefilters

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.AnteDecorator = MessageFilterDecorator{}

// MessageFilterDecorator is an AnteDecorator runs a list of filters before calling the next AnteHandler.
type MessageFilterDecorator struct {
	filters []MessageFilter
}
type MessageFilter func(ctx sdk.Context, msg sdk.Msg, simulate bool) error

func NewMessageFilterDecorator(
	filters ...MessageFilter,
) MessageFilterDecorator {
	return MessageFilterDecorator{
		filters,
	}
}

func (palc MessageFilterDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, m := range tx.GetMsgs() {
		for _, f := range palc.filters {
			if err := f(ctx, m, simulate); err != nil {
				return ctx, sdkerrors.Wrap(err, "filter invariant")
			}
		}
	}

	return next(ctx, tx, simulate)
}
