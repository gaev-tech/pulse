package email

import (
	"context"

	"go.uber.org/zap"
)

type LogSender struct {
	log *zap.Logger
}

func NewLogSender(log *zap.Logger) *LogSender {
	return &LogSender{log: log}
}

func (sender *LogSender) SendMagicLink(_ context.Context, toEmail, link string) error {
	sender.log.Info("magic link (dev mode — email not sent)",
		zap.String("to", toEmail),
		zap.String("link", link),
	)
	return nil
}
