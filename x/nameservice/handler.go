package nameservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Msg & Handler，主要用来处理共识。每一个Msg会对应一个Handler。
//需要强调的是所有的Hander调用前都会执行一个叫anteHander的前置处理器，来做一些通用的工作，如：安全校验。

// NewHandler 返回“nameservice”类型消息的处理程序。
// NewHandler本质上是一个子路由（Switch-case），它将进入该模块的msg路由到正确的handler做处理。
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// Todo:这里应当有个Msg的包和Err的包出错，并不全在sdk/types中
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// SDK中handler的命名规范是handlerMsg{.Action}

// 处理set name消息
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // 检查消息发送人是否与当前所有者相同
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // 若不是，抛出一个错误
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // 若是，则将域名设置为消息中指定的值。
	return sdk.Result{}                      // return
}

// 处理buy name消息
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}
