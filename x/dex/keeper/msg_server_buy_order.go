package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"interchange/x/dex/types"
)

func (k msgServer) SendBuyOrder(goCtx context.Context, msg *types.MsgSendBuyOrder) (*types.MsgSendBuyOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.AmountDenom, msg.PriceDenom)
	_, found := k.GetBuyOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgSendBuyOrderResponse{}, nil
	}

	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgSendBuyOrderResponse{}, err
	}

	if err := k.SafeBurn(ctx, msg.Port, msg.ChannelID, sender, msg.PriceDenom, msg.Amount*msg.Price); err != nil {
		return &types.MsgSendBuyOrderResponse{}, err
	}

	k.SaveVoucherDenom(ctx, msg.Port, msg.ChannelID, msg.PriceDenom)
	var packet types.BuyOrderPacketData
	packet.Buyer = msg.Creator

	return &types.MsgSendBuyOrderResponse{}, nil
}
