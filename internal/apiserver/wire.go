//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"
	"github.com/yanking/miniblog/internal/apiserver/biz"
	"github.com/yanking/miniblog/internal/apiserver/store"
	"github.com/yanking/miniblog/internal/pkg/server"
	"github.com/yanking/miniblog/pkg/restserver/middlewares"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB, // 提供数据库实例
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(middlewares.UserRetriever), new(*UserRetriever)),
		),
	)
	return nil, nil
}
