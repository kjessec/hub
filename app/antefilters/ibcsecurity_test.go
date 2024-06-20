package antefilters_test

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcicstypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/mars-protocol/hub/app/antefilters"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIBCSecurityMessageFilter(t *testing.T) {
	t.Run("ibctransfertypes.MsgTransfer", func(t *testing.T) {
		t.Run("receiver address must not exceed 2048 bytes", func(t *testing.T) {
			// must error
			assert.Error(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibctransfertypes.MsgTransfer{
					Receiver: createBech32AddressButVeryLong(2048),
				},
				false,
			))

			// within safe length - must NOT error
			assert.NoError(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibctransfertypes.MsgTransfer{
					Receiver: createBech32AddressButVeryLong(12),
				},
				false,
			))
		})

		t.Run("receiver address must be bech32", func(t *testing.T) {
			// must error
			assert.Error(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibctransfertypes.MsgTransfer{
					Receiver: "invalidbech32",
				},
				false,
			))

			// within safe length - must NOT error
			assert.NoError(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibctransfertypes.MsgTransfer{
					Receiver: createBech32AddressButVeryLong(12),
				},
				false,
			))
		})

		t.Run("memo length must not exceed 32768 bytes", func(t *testing.T) {
			// must error
			assert.Error(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibctransfertypes.MsgTransfer{
					// use a long address as memo
					Memo:     createBech32AddressButVeryLong(32768),
					Receiver: createBech32AddressButVeryLong(12),
				},
				false,
			))

			// empty memo - valid
			assert.NoError(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibctransfertypes.MsgTransfer{
					// use a long address as memo
					Memo:     "",
					Receiver: createBech32AddressButVeryLong(12),
				},
				false,
			))

			// short memo - valid
			assert.NoError(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibctransfertypes.MsgTransfer{
					// use a long address as memo
					Memo:     "sdfdff",
					Receiver: createBech32AddressButVeryLong(12),
				},
				false,
			))
		})
	})

	t.Run("ibcicstypes.MsgSendTx", func(t *testing.T) {
		t.Run("owner address must not exceed 2048 bytes", func(t *testing.T) {
			// must error
			assert.Error(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibcicstypes.MsgSendTx{
					Owner: createBech32AddressButVeryLong(2048),
				},
				false,
			))

			// within safe length - must NOT error
			assert.NoError(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibcicstypes.MsgSendTx{
					Owner: createBech32AddressButVeryLong(12),
				},
				false,
			))
		})
	})

	t.Run("ibcfeetypes.MsgRegisterCounterpartyPayee", func(t *testing.T) {
		t.Run("counterparty address must not exceed 2048 bytes", func(t *testing.T) {
			// must error
			assert.Error(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibcicstypes.MsgSendTx{
					Owner: createBech32AddressButVeryLong(2048),
				},
				false,
			))

			// within safe length - must NOT error
			assert.NoError(t, antefilters.IBCTransferReceiverLengthCheck(
				sdk.Context{},
				&ibcicstypes.MsgSendTx{
					Owner: createBech32AddressButVeryLong(12),
				},
				false,
			))

		})
	})
}

func createBech32AddressButVeryLong(minLength int) string {
	data := []byte(strings.Repeat("a", minLength*2))
	conv, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		panic(err)
	}

	b32addr, err := bech32.Encode(
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		conv,
	)
	if err != nil {
		panic(errors.Wrap(err, "bech32 encoding"))
	}

	return b32addr
}
