package i_moneropay

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moneronodo/sshui/internal/base"
	"github.com/moneronodo/sshui/internal/model/moneropay"
)

const (
	connMoneropay = "file:///home/nodo/moneropay.sqlite?immutable=1"
	TxListSize    = 10
)

type MpayTxUpdateMsg struct {
	Transaction Transaction
	Index       int
}

type MpayTxListMsg struct {
	Transactions []Transaction
}

type Transaction struct {
	Subaddress  string
	Expected    uint64
	Covered     moneropay.MoneropayReceiveCovered
	TxIds       []moneropay.MoneropayReceiveTx
	Description string
	CreatedAt   time.Time
	Queried     bool // queried moneropay?
	Complete    bool
}

func GetTxList() []Transaction {
	txs := []Transaction{}
	db, err := sql.Open("sqlite3", connMoneropay)
	if err != nil {
		spew.Fdump(base.Dump, err)
		return txs
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := db.QueryContext(
		ctx,
		"SELECT address, expected_amount, description, created_at"+
			" FROM subaddresses, receivers"+
			" WHERE address_index=subaddress_index"+
			" ORDER BY address_index DESC"+
			" LIMIT ?;",
		TxListSize,
	)
	if err != nil {
		spew.Fdump(base.Dump, err)
		return txs
	}
	defer rows.Close()
	for i := 0; i < TxListSize && rows.Next(); i++ {
		tx := Transaction{}
		err := rows.Scan(&tx.Subaddress, &tx.Expected, &tx.Description, &tx.CreatedAt)
		if err != nil {
			spew.Fdump(base.Dump, err)
		} else {
			txs = append(txs, tx)
		}
	}
	return txs
}

func GetHealth(url string) *moneropay.MoneropayHealth {
	var j = &moneropay.MoneropayHealth{}
	c := &http.Client{Timeout: 3 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		spew.Fprintf(base.Dump, "http: %v\n", err)
		return j
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(j); err != nil {
		spew.Fprintf(base.Dump, "Decode: %v\n", err)
	}
	return j
}
