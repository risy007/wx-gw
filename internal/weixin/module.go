package weixin

import "go.uber.org/fx"

var Module = fx.Module("weixin", fx.Provide(NewWxService))
