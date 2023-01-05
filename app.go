package app

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	dbm "github.com/tendermint/tendermint/light/store/db"
)

const (
	appName = "nameservice"
)

type nameServiceApp struct {
	*bam.BaseApp
}

// todo：DB包出错，需要重新定位
func NewNameServiceApp(logger log.Logger, db dbm.DB) *nameServiceApp {

	// 定义模块共享的最上册编码解码器
	cdc := MakeCodec()

	// BaseApp 通过 ABCI 协议处理与 Tendermint 的交互
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &nameServiceApp{
		BaseApp: bApp,
		cdc:     cdc,
	}

	return app
}
