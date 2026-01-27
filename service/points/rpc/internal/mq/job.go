package mq

import (
	"github.com/IBM/sarama"
	"github.com/shopspring/decimal"
)

type JobMsg struct {
	AccountID int64           `json:"account_id"`
	UserID    int64           `json:"user_id"`
	Amount    decimal.Decimal `json:"amount"`
	TraceID   string          `json:"traceing"`
}

type Job struct {
	Msg     *sarama.ConsumerMessage
	Session sarama.ConsumerGroupSession
}
