package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"interchange/x/dex/types"
)

func (k msgServer) CancelBuyOrder(goCtx context.Context, msg *types.MsgCancelBuyOrder) (*types.MsgCancelBuyOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pairIndex := types.OrderBookIndex(msg.Port, msg.Channel, msg.AmountDenom, msg.PriceDenom)
	b, found := k.GetBuyOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgCancelBuyOrderResponse{}, errors.New("the pair doesn't exist")
	}

	order, err := b.Book.GetOrderFromId(msg.OrderID)
	if err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}

	if order.Creator != msg.Creator {
		return &types.MsgCancelBuyOrderResponse{}, errors.New("canceller must be creator")
	}

	if err := b.Book.RemoveOrderFromId(msg.OrderID); err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}

	k.SetBuyOrderBook(ctx, b)

	buyer, err := sdk.AccAddressFromBech32(order.Creator)
	if err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}

	if err := k.SafeMint(ctx, msg.Port, msg.Channel, buyer, msg.PriceDenom, order.Amount*order.Price); err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}

	return &types.MsgCancelBuyOrderResponse{}, nil
}
