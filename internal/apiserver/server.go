package apiserver

import (
	"context"
	"github.com/onexstack/onexstack/pkg/store/where"
	"github.com/yanking/miniblog/internal/apiserver/biz"
	"github.com/yanking/miniblog/internal/apiserver/model"
	"github.com/yanking/miniblog/internal/apiserver/store"
	"github.com/yanking/miniblog/internal/pkg/contextx"
	"github.com/yanking/miniblog/internal/pkg/known"
	"github.com/yanking/miniblog/internal/pkg/log"
	"github.com/yanking/miniblog/internal/pkg/server"
	"github.com/yanking/miniblog/pkg/restserver/middlewares"
	"github.com/yanking/miniblog/pkg/token"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"

	genericoptions "github.com/yanking/miniblog/pkg/options"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器.
	GRPCServerMode = "grpc"
	// GRPCServerMode 定义 gRPC + HTTP 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器.
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置.
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项.
type Config struct {
	ServerMode   string
	JWTKey       string
	Expiration   time.Duration
	HTTPOptions  *genericoptions.HTTPOptions
	GRPCOptions  *genericoptions.GRPCOptions
	MySQLOptions *genericoptions.MySQLOptions
}

// UnionServer 定义一个联合服务器. 根据 ServerMode 决定要启动的服务器类型.
//
// 联合服务器分为以下 2 大类：
//  1. Gin 服务器：由 Gin 框架创建的标准的 REST 服务器。根据是否开启 TLS，
//     来判断启动 HTTP 或者 HTTPS；
//  2. GRPC 服务器：由 gRPC 框架创建的标准 RPC 服务器
//  3. HTTP 反向代理服务器：由 grpc-gateway 框架创建的 HTTP 反向代理服务器。
//     根据是否开启 TLS，来判断启动 HTTP 或者 HTTPS；
//
// HTTP 反向代理服务器依赖 gRPC 服务器，所以在开启 HTTP 反向代理服务器时，会先启动 gRPC 服务器.
type UnionServer struct {
	srv server.Server
}

// ServerConfig 包含服务器的核心依赖和配置.
type ServerConfig struct {
	cfg       *Config
	biz       biz.IBiz
	retriever middlewares.UserRetriever
}

// NewUnionServer 根据配置创建联合服务器.
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	// 一些初始化代码
	// 注册租户解析函数，通过上下文获取用户 ID
	//nolint: gocritic
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	// 初始化 token 包的签名密钥、认证 Key 及 Token 默认过期时间
	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)

	// 创建服务配置，这些配置可用来创建服务器
	serverConfig, err := cfg.NewServerConfig()
	if err != nil {
		return nil, err
	}

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	// 根据服务模式创建对应的服务实例
	// 实际企业开发中，可以根据需要只选择一种服务器模式.
	// 这里为了方便给你展示，通过 cfg.ServerMode 同时支持了 Gin 和 GRPC 2 种服务器模式.
	// 默认为 gRPC 服务器模式.
	var srv server.Server
	switch cfg.ServerMode {
	case GinServerMode:
		srv, err = serverConfig.NewGinServer(), nil
	default:
		srv, err = serverConfig.NewGRPCServerOr()
	}
	if err != nil {
		return nil, err
	}

	return &UnionServer{srv: srv}, nil
}

// Run 运行应用.
func (s *UnionServer) Run() error {
	go s.srv.RunOrDie()

	// 创建一个 os.Signal 类型的 channel，用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞程序，等待从 quit channel 中接收到信号
	<-quit

	log.Infow("Shutting down server ...")

	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	s.srv.GracefulStop(ctx)

	log.Infow("Server exited")
	return nil
}

// NewServerConfig 创建一个 *ServerConfig 实例.
// 进阶：这里其实可以使用依赖注入的方式，来创建 *ServerConfig.
func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	// 初始化数据库连接
	db, err := cfg.NewDB()
	if err != nil {
		return nil, err
	}
	store := store.NewStore(db)
	return &ServerConfig{
		cfg:       cfg,
		biz:       biz.NewBiz(store),
		retriever: &UserRetriever{store: store},
	}, nil

}

// NewDB 创建一个 *gorm.DB 实例.
func (cfg *Config) NewDB() (*gorm.DB, error) {
	return cfg.MySQLOptions.NewDB()
}

// UserRetriever 定义一个用户数据获取器. 用来获取用户信息.
type UserRetriever struct {
	store store.IStore
}

// GetUser 根据用户 ID 获取用户信息.
func (r *UserRetriever) GetUser(ctx context.Context, userID string) (*model.UserM, error) {
	return r.store.User().Get(ctx, where.F("userID", userID))
}
