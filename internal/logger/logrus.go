package logger

import (
	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
	"wx-gw/config"
)

func NewLogrusLogger(conf *config.Config) (log *logrus.Logger, err error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,                  // 完整时间
		TimestampFormat: "2006-01-02 15:04:05", // 时间格式
		ForceColors:     true,                  //显示颜色
	})

	logger.SetLevel(logrus.DebugLevel)
	if !conf.Log.Development {
		logger.SetLevel(logrus.InfoLevel)
	}

	if conf.Log.ToFile {
		logpath := conf.Log.Directory
		hook, err := lumberjackrus.NewHook(
			&lumberjackrus.LogFile{
				Filename:   logpath + conf.Name + ".log",
				MaxSize:    100,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   false,
				LocalTime:  false,
			},
			logrus.DebugLevel,
			&logrus.TextFormatter{
				FullTimestamp:   true,                  // 完整时间
				TimestampFormat: "2006-01-02 15:04:05", // 时间格式
				ForceColors:     true,                  //显示颜色
			},
			&lumberjackrus.LogFileOpts{
				logrus.DebugLevel: &lumberjackrus.LogFile{
					Filename: logpath + "debug.log",
				},
				logrus.InfoLevel: &lumberjackrus.LogFile{
					Filename: logpath + "info.log",
				},
				logrus.ErrorLevel: &lumberjackrus.LogFile{
					Filename:   logpath + "error.log",
					MaxSize:    100,   // optional
					MaxBackups: 1,     // optional
					MaxAge:     1,     // optional
					Compress:   false, // optional
					LocalTime:  false, // optional
				},
			},
		)

		if err != nil {
			return nil, err
		}
		logger.AddHook(hook)
		return logger, nil
	}
	return logger, nil
}
