package daemonrpc

type DaemonRPCMsg struct {
	Response DaemonRPCResponseWrapper
}

type DaemonRPCResponseWrapper any

type DaemonRPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Result  any    `json:"result"`
	Error   any    `json:"error"`
}

type DaemonRPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type GetVersionHardForks struct {
	Height    uint32 `json:"height"`
	HfVersion uint   `json:"hf_version"`
}

type DaemonResponseBodyGetVersion struct {
	CurrentHeight uint32                `json:"current_height"`
	HardForks     []GetVersionHardForks `json:"hard_forks"`
	Release       bool                  `json:"release"`
	Status        string                `json:"status"`
	Untrusted     bool                  `json:"untrusted"`
	Version       uint32                `json:"version"`
}

type DaemonResponseBodyGetInfo struct {
	AdjustedTime              uint32 `json:"adjusted_time"`
	AltBlocksCount            uint32 `json:"alt_blocks_count"`
	BlockSizeLimit            uint   `json:"block_size_limit"`
	BlockSizeMedian           uint   `json:"block_size_median"`
	BlockWeightLimit          uint   `json:"block_weight_limit"`
	BlockWeightMedian         uint   `json:"block_weight_median"`
	BootstrapDaemonAddress    string `json:"bootstrap_daemon_address"`
	BusySyncing               bool   `json:"busy_syncing"`
	Credits                   int    `json:"credits"`
	CumulativeDifficulty      uint64 `json:"cumulative_difficulty"`
	CumulativeDifficultyTop64 uint32 `json:"cumulative_difficulty_top64"`
	DatabaseSize              uint64 `json:"database_size"`
	Difficulty                uint64 `json:"difficulty"`
	DifficultyTop64           uint32 `json:"difficulty_top64"`
	FreeSpace                 uint64 `json:"free_space"`
	GreyPeerlistSize          int    `json:"grey_peerlist_size"`
	Height                    int    `json:"height"`
	HeightWithoutBootstrap    int    `json:"height_without_bootstrap"`
	IncomingConnectionsCount  int    `json:"incoming_connections_count"`
	Mainnet                   bool   `json:"mainnet"`
	Nettype                   string `json:"nettype"`
	Offline                   bool   `json:"offline"`
	OutgoingConnectionsCount  int    `json:"outgoing_connections_count"`
	RpcConnectionsCount       int    `json:"rpc_connections_count"`
	Stagenet                  bool   `json:"stagenet"`
	StartTime                 int    `json:"start_time"`
	Status                    string `json:"status"`
	Synchronized              bool   `json:"synchronized"`
	Target                    int    `json:"target"`
	TargetHeight              int    `json:"target_height"`
	Testnet                   bool   `json:"testnet"`
	TopBlockHash              string `json:"top_block_hash"`
	TopHash                   string `json:"top_hash"`
	TxCount                   int    `json:"tx_count"`
	TxPoolSize                int    `json:"tx_pool_size"`
	Untrusted                 bool   `json:"untrusted"`
	UpdateAvailable           bool   `json:"update_available"`
	Version                   string `json:"version"`
	WasBootstrapEverUsed      bool   `json:"was_bootstrap_ever_used"`
	WhitePeerlistSize         int    `json:"white_peerlist_size"`
	WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
	WideDifficulty            string `json:"wide_difficulty"`
}

func MakeDaemonRPCResponse(response DaemonRPCResponseWrapper) *DaemonRPCResponse {
	resp := &DaemonRPCResponse{
		Result: response,
	}
	return resp
}
