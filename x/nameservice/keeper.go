package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	//"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	// StoreKey的引用包位置改变
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
Keeper模块：定义数据库的交互
实现Set，Get等操作；可以设置迭代器，用于循环获取store中的键值对
*/

// Keeper 维护与数据存储的链接，并为状态机的各个部分公开 getter/setter 方法
type Keeper struct {
	coinKeeper keeper.Keeper       // bank模块的Keeper引用
	storeKey   storetypes.StoreKey // 从 sdk.Context 访问存储的未公开密钥
	cdc        *codec.Codec        // 访问编码解码器的指针（被Amino用于编码及解码二进制机构的编码解码器的指针）
}

// Keeper的构造函数
func NewKeeper(coinKeeper keeper.Keeper, storeKey storetypes.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Set方法，为指定域名设置解析字符串值
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	// Todo：cdc库出错
	// 存储只接受[]byte
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// Get方法：解析域名，查找域名对应的解析值
func (k Keeper) GetWhois(ctx sdk.Context, name string) Whois {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

// ResolveName - 返回解析得到的字符串
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - 重用Set方法设置解析值
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - 返回name是否已经有所有者
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

// GetOwner - 获取域名的当前所有者
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - 设置域名的当前所有者
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - 获取域名的当前价格。 如果价格还不存在，设置为 1 name token
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// SetPrice - 设置域名的当前价格
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// SDK 还有一个特性叫 sdk.Iterator，可以返回一个迭代器用于遍历指定 store 中的所有 <Key, Value> 对。
// 获取所有域名的迭代器，其中key是域名，value 是 whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
