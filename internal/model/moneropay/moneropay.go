package moneropay

import "time"

type MpayHealthMsg struct {
	Health MoneropayHealth
}

type MoneropayStatus int

type MoneropayServices struct {
	Walletrpc  bool
	Sqlite     bool
	Postgresql bool
}

type MoneropayHealth struct {
	Status   MoneropayStatus
	Services MoneropayServices
}

type MoneropayReceive struct {
	Amount       MoneropayReceiveAmount `json:"amount"`
	Complete     bool                   `json:"complete"`
	Description  string                 `json:"description"`
	CreatedAt    time.Time              `json:"created_at"`
	Transactions []MoneropayReceiveTx     `json:"transactions"`
}

type MoneropayReceiveCovered struct {
	Total    uint64 `json:"total"`
	Unlocked uint64 `json:"unlocked"`
}

type MoneropayReceiveAmount struct {
	Expected uint64                  `json:"expected"`
	Covered  MoneropayReceiveCovered `json:"covered"`
}

type MoneropayReceiveTx struct {
	Amount          uint64    `json:"amount"`
	Confirmations   uint32    `json:"confirmations"`
	DoubleSpendSeen bool      `json:"double_spend_seen"`
	Fee             uint64    `json:"fee"`
	Timestamp       time.Time `json:"timestamp"`
	TxHash          string    `json:"tx_hash"`
	UnlockTime      uint64    `json:"unlock_time"`
	Locked          bool      `json:"locked"`
}
