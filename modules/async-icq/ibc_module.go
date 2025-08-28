package icq

import (
	"strings"

	"github.com/cosmos/ibc-apps/modules/async-icq/v10/keeper"
	"github.com/cosmos/ibc-apps/modules/async-icq/v10/types"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

var _ porttypes.IBCModule = &IBCModule{}

// IBCModule implements the ICS26 interface for interchain query host chains
type IBCModule struct {
	keeper keeper.Keeper
}

// NewIBCModule creates a new IBCModule given the associated keeper
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	_ []string,
	portID string,
	channelID string,
	_ channeltypes.Counterparty,
	version string,
) (string, error) {
	if !im.keeper.IsHostEnabled(ctx) {
		return "", types.ErrHostDisabled
	}

	if err := ValidateICQChannelParams(ctx, im.keeper, order, portID, channelID); err != nil {
		return "", err
	}

	if strings.TrimSpace(version) == "" {
		version = types.Version
	}

	if version != types.Version {
		return "", errors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", version, types.Version)
	}

	// Claim channel capability passed back by IBC module
	// if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
	// 	return "", err
	// }

	return version, nil
}

func ValidateICQChannelParams(
	ctx sdk.Context,
	keeper keeper.Keeper,
	order channeltypes.Order,
	portID string,
	_ string,
) error {
	if order != channeltypes.UNORDERED {
		return errors.Wrapf(channeltypes.ErrInvalidChannelOrdering, "expected %s channel, got %s", channeltypes.UNORDERED, order)
	}

	boundPort := keeper.GetPort(ctx)
	if portID != boundPort {
		return errors.Wrapf(types.ErrInvalidHostPort, "expected %s, got %s", boundPort, portID)
	}
	return nil
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	_ []string,
	portID,
	channelID string,
	_ channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	if !im.keeper.IsHostEnabled(ctx) {
		return "", types.ErrHostDisabled
	}

	if err := ValidateICQChannelParams(ctx, im.keeper, order, portID, channelID); err != nil {
		return "", err
	}

	if counterpartyVersion != types.Version {
		return "", errors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", counterpartyVersion, types.Version)
	}

	// Module may have already claimed capability in OnChanOpenInit in the case of crossing hellos
	// (ie chainA and chainB both call ChanOpenInit before one of them calls ChanOpenTry)
	// If module can already authenticate the capability then module already owns it so we don't need to claim
	// Otherwise, module does not have channel capability and we must claim it from IBC
	// if !im.keeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
	// 	// Only claim channel capability passed back by IBC module if we do not already own it
	// 	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
	// 		return "", err
	// 	}
	// }

	return types.Version, nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	_ string,
	_ string,
	_ string,
	counterpartyVersion string,
) error {
	if !im.keeper.IsHostEnabled(ctx) {
		return types.ErrHostDisabled
	}

	if counterpartyVersion != types.Version {
		return errors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", counterpartyVersion, types.Version)
	}
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	_ string,
	_ string,
) error {
	if !im.keeper.IsHostEnabled(ctx) {
		return types.ErrHostDisabled
	}
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	_ sdk.Context,
	_ string,
	_ string,
) error {
	// Ensure channels cannot be closed by users
	return errors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	_ sdk.Context,
	_ string,
	_ string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) ibcexported.Acknowledgement {
	if !im.keeper.IsHostEnabled(ctx) {
		return channeltypes.NewErrorAcknowledgement(types.ErrHostDisabled)
	}

	txResponse, err := im.keeper.OnRecvPacket(ctx, packet)
	if err != nil {
		// Emit an event including the error msg
		keeper.EmitWriteErrorAcknowledgementEvent(ctx, packet, err)

		return channeltypes.NewErrorAcknowledgement(err)
	}

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return channeltypes.NewResultAcknowledgement(txResponse)
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	_ sdk.Context,
	_ string,
	_ channeltypes.Packet,
	_ []byte,
	_ sdk.AccAddress,
) error {
	return errors.Wrap(types.ErrInvalidChannelFlow, "cannot receive acknowledgement on a host channel end, a host chain does not send a packet over the channel")
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	_ sdk.Context,
	channelVersion string,
	_ channeltypes.Packet,
	_ sdk.AccAddress,
) error {
	return errors.Wrap(types.ErrInvalidChannelFlow, "cannot cause a packet timeout on a host channel end, a host chain does not send a packet over the channel")
}
