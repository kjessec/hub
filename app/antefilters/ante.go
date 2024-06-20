package antefilters

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
				return ctx, errors.Wrap(err, "message filter invariant")
			}
		}
	}

	return next(ctx, tx, simulate)
}
