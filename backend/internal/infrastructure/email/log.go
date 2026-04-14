package email

import (
	"context"

	"go.uber.org/zap"
)

// LogSender prints magic links to stdout instead of sending email (non-production use).
type LogSender struct {
	log *zap.Logger
}

// NewLogSender creates a new LogSender.
func NewLogSender(log *zap.Logger) *LogSender {
	return &LogSender{log: log}
}

// SendMagicLink logs the magic link instead of sending an email.
func (logSender *LogSender) SendMagicLink(_ context.Context, toEmail, link string) error {
	logSender.log.Info("magic link (dev mode — email not sent)",
		zap.String("to", toEmail),
		zap.String("link", link),
	)
	return nil
}
