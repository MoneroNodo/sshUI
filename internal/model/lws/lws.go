package lws

type LwsAccountStatus int

const (
	LwsAccountActive   LwsAccountStatus = iota
	LwsAccountInactive LwsAccountStatus = iota
	LwsAccountHidden   LwsAccountStatus = iota
)

type LwsBase58InvalidErr struct{}
type LwsViewkeyInvalidErr struct{}
type LwsAddrKeyMismatchErr struct{}

func (e *LwsBase58InvalidErr) Error() string {
	return "Invalid base58 address"
}

func (e *LwsViewkeyInvalidErr) Error() string {
	return "Invalid hex"
}

func (e *LwsAddrKeyMismatchErr) Error() string {
	return "Address / viewkey mismatch"
}

type LwsAccount struct {
	Address    string `json:"address"`
	ScanHeight int32  `json:"scan_height"`
	AccessTime int64  `json:"access_time"`
	Status LwsAccountStatus
}

type LwsAddress string

type LwsUpdated struct {
	Updated []LwsAddress
}

type LwsListAccounts struct {
	Active   []LwsAccount `json:"active"`
	Inactive []LwsAccount `json:"inactive"`
	Hidden   []LwsAccount `json:"hidden"`
}

type LwsRequest struct {
	Address     string `json:"address"`
	StartHeight int32  `json:"start_height"`
}

type LwsListReqeusts struct {
	Create []LwsRequest `json:"create"`
	Import []LwsRequest `json:"import"`
}
