package mq

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"wx-gw/config"
	"wx-gw/internal/weixin"
)

type (
	mqInParams struct {
		fx.In
		Logger    *zap.Logger
		WXService *weixin.Service
	}
	mqOutResult struct {
		fx.Out
		MicroService *Service
	}

	Service struct {
		log   *zap.Logger
		cfg   *config.Config
		wecom *weixin.Service

		once   sync.Once
		cancel context.CancelFunc
		nc     *nats.Conn
	}
)

func (s *Service) Start(ctx context.Context) error {
	s.once.Do(func() {
		nc, err := nats.Connect(s.cfg.Nats.Address)
		if err != nil {
			return
		}
		s.nc = nc
		s.log.Info("nats 初始化连接成功", zap.String("address", s.cfg.Nats.Address))
		ctx, s.cancel = context.WithCancel(ctx)
		go s.Worker(ctx)
	})
	return nil
}

func (s *Service) Stop() error {
	s.cancel()
	if s.nc != nil {
		s.nc.Close()
	}
	return nil
}

// newService - Create service
func NewService(cfg *config.Config, in mqInParams) mqOutResult {
	var (
		out mqOutResult
	)
	srv := &Service{
		log:    in.Logger,
		cfg:    cfg,
		cancel: func() {},
		wecom:  in.WXService,
	}
	out.MicroService = srv
	return out
}

func (s *Service) Worker(ctx context.Context) {
	if len(s.cfg.Nats.Subscribes) <= 0 {
		s.log.Error("nats subscribes is empty！")
		return
	}

	msgCh := make(chan *nats.Msg, nats.DefaultMaxChanLen)
	defer close(msgCh)

	for _, ss := range s.cfg.Nats.Subscribes {
		sub, err := s.nc.ChanSubscribe(ss, msgCh)
		if err != nil {
			s.log.Fatal("订阅NATS消息失败", zap.Error(err))
			return
		}
		s.log.Info("订阅 NATS 消息成功！", zap.String("Subject", sub.Subject))

	}

	for {
		select {
		case <-ctx.Done():
			//s.nc.Flush()
			//s.nc.Drain()
			if len(msgCh) == 0 {
				break
			}
		case msg := <-msgCh:
			//s.log.Info(msg.Subject, zap.ByteString("msg", msg.Data))
			s.MsgHandler(ctx, msg)
		}
	}
	s.log.Info("停止微服务成功！")
}

func (s *Service) MsgHandler(ctx context.Context, msg *nats.Msg) {
	if msg.Subject == "prtg.alert" {
		//处理prtg.alert消息
		var plmsg weixin.PrtgAlertMessage
		err := json.Unmarshal(msg.Data, &plmsg)
		if err != nil {
			s.log.Error(err.Error())
			return
		}
		wxmsg := plmsg.ConvertToWorkWxMDMsg()
		s.wecom.SendMarkDownMessage(wxmsg.MarkDown.Content)
		//s.log.Debug(wxmsg.MarkDown.Content)
	} else if msg.Subject == "radius.auth.reject" {
		//处理radius.auth.reject消息
		var ramsg
	}
	//s.log.Info(fmt.Sprintf("接收到 NATS Message ，Subject = %s", msg.Subject), zap.ByteString("msg", msg.Data))
	//s.wecom.SendTextMessage(msg.Subject)
}
