package internal

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"wx-gw/config"
	"wx-gw/internal/logger"
	"wx-gw/internal/mq"
	"wx-gw/internal/weixin"
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
