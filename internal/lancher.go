package internal

import (
	"context"
	"github.com/risy007/wx-gw/config"
	"github.com/risy007/wx-gw/internal/logger"
	"github.com/risy007/wx-gw/internal/mq"
	"github.com/risy007/wx-gw/internal/weixin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	config.Module,
	logger.Module,
	mq.Module,
	weixin.Module,
	fx.Invoke(launcher),
)

func launcher(
	lifecycle fx.Lifecycle,
	config *config.Config,
	logger *zap.Logger,
	mqService *mq.Service,
	wxService *weixin.Service,
) {

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("启动应用......")

			go func() {
				err := mqService.Start(ctx)
				if err != nil {
					logger.Error(err.Error())
				}
				wxService.Start()
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {

			return nil
		},
	})
}
