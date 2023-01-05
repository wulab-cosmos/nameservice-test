package nameservice

import (
	"encoding/json"

	error "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Msg & Handler，主要用来处理共识。每一个Msg会对应一个Handler。
//需要强调的是所有的Hander调用前都会执行一个叫anteHander的前置处理器，来做一些通用的工作，如：安全校验。

// SDK中Msg的命令约束是 Msg{.Action}（命名规范）

// MsgSetName defines a SetName message
type MsgSetName struct {
	Name  string
	Value string
	Owner sdk.AccAddress
}

// MsgBuyName defines the BuyName message
type MsgBuyName struct {
	Name  string
	Bid   sdk.Coins
	Buyer sdk.AccAddress
}

// NewMsgSetName 定义MsgSetName的构造函数
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,  // 所要设置的域名
		Value: value, // 要设置的域名解析值
		Owner: owner, // 域名的所有者
	}
}

// NewMsgBuyName 定义MsgBuyName的构造函数
func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route 需要返回模块的名称
func (msg MsgSetName) Route() string { return "nameservice" }
func (msg MsgBuyName) Route() string { return "nameservice" }

// Type 需要返回action
func (msg MsgSetName) Type() string { return "set_name" }
func (msg MsgBuyName) Type() string { return "buy_name" }

// ValidateBasic 对消息运行无状态检查（ValidateBasic用于对Msg的有效性进行一些基本的无状态检查。在此情形下，需要检查没有属性为空。请注意这里使用sdk.Error类型。 SDK提供了一组应用开发人员经常遇到的错误类型。）
func (msg MsgSetName) ValidateBasic() error.Error {
	if msg.Owner.Empty() {
		// todo：error类型报错
		return error.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Value cannot be empty")
	}
	return nil
}

// todo：将error包更改正确
func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes 对消息进行编码以进行签名
// GetSignBytes定义了如何编码Msg以进行签名。在大多数情形下，要编码成排好序的JSON。不应修改输出。
func (msg MsgSetName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgBuyName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners 定义需要谁的签名
// GetSigners定义一个Tx上需要哪些人的签名才能使其有效。在这种情形下，MsgSetName要求域名所有者在尝试重置域名解析值时要对该交易签名。
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
