package antefilters

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibcicstypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

// https://github.com/cosmos/ibc-go/commit/ac5b2ca4e37c252d28afd1b87e976dbd9954215b
const (
	MaximumReceiverLength = 2048  // maximum length of the receiver address in bytes (value chosen arbitrarily)
	MaximumMemoLength     = 32768 // maximum length of the memo in bytes (value chosen arbitrarily)
)

// IBCTransferReceiverLengthCheck is a filter that additionally checks validity of certain fields in IBC messages.
// This exists to backport IBC security patch (https://github.com/cosmos/ibc-go/commit/ac5b2ca4e37c252d28afd1b87e976dbd9954215b#diff-3f0867015c073be1949877f2c8ed217edb70ca70ecc29412e55add06d70a2bb5)
// to the particular version of cosmos/ibc-go that we are using.
//
// It is assumed that the rest of the checks are implemented in the respective ValidateBasic() methods.
// For this, this filter only implements certain extra checks that are known to be missing.
//
// NOTE: consider upgrading to cosmos/ibc-go/v8
func IBCTransferReceiverLengthCheck(_ sdk.Context, m sdk.Msg, _ bool) error {
	switch msg := m.(type) {

	case *ibctransfertypes.MsgTransfer:
		if len(msg.Receiver) > MaximumReceiverLength {
			return errors.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"recipient address must not exceed %d bytes",
				MaximumReceiverLength,
			)
		}

		if len(msg.Memo) > MaximumMemoLength {
			return errors.Wrapf(
				sdkerrors.ErrMemoTooLarge,
				"memo length must not exceed %d bytes",
				MaximumMemoLength,
			)
		}

		if _, err := sdk.AccAddressFromBech32(msg.Receiver); err != nil {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"receiver address is not bech32: %v",
				err,
			)
		}

	case *ibcicstypes.MsgSendTx:
		if len(msg.Owner) > MaximumReceiverLength {
			return errors.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"owner address must not exceed %d bytes",
				MaximumReceiverLength,
			)
		}

	case *ibcfeetypes.MsgRegisterCounterpartyPayee:
		if len(msg.CounterpartyPayee) > MaximumReceiverLength {
			return errors.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"counterparty address must not exceed %d bytes",
				MaximumReceiverLength,
			)
		}
	}

	return nil
}
