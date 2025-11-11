package i_lws

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/moneronodo/sshui/internal/model/lws"
)

const (
	prog = "/home/nodo/bin/monero-lws-admin"
	args = "--db-path=/media/monero/bitmonero/light_wallet_server"
)

func command(arguments ...string) ([]byte, error) {
	a := append([]string{args}, arguments...)
	c, err := exec.Command(prog, a...).CombinedOutput()
	if err != nil {
		if strings.HasPrefix(string(c), "View key has invalid hex") {
			return nil, &lws.LwsViewkeyInvalidErr{}
		}
		if strings.HasSuffix(string(c), "Address/viewkey mismatch") {
			return nil, &lws.LwsAddrKeyMismatchErr{}
		}
		if strings.HasPrefix(string(c), "Invalid base58 public address - wrong") {
			return nil, &lws.LwsBase58InvalidErr{}
		}
		return nil, nil
	}
	return c, nil
}

func ListAccounts() (lws.LwsListAccounts, error) {
	accs := lws.LwsListAccounts{}
	c, err := command("list_accounts")
	if err != nil {
		return accs, err
	}
	err = json.Unmarshal(c, &accs)
	return accs, err
}

func ListRequests() (lws.LwsListReqeusts, error) {
	accs := lws.LwsListReqeusts{}
	c, err := command("list_requests")
	if err != nil {
		return accs, err
	}
	err = json.Unmarshal(c, &accs)
	return accs, err
}

func AddAccount(address, viewkey string) error {
	_, err := command("add_account", address, viewkey)
	return err
}

func DeleteAccount(address string) error {
	_, err := command("modify_account_status", "hidden", address)
	return err
}

func DeactivateAccount(address string) error {
	_, err := command("modify_account_status", "inactive", address)
	return err
}

func ReactivateAccount(address string) error {
	_, err := command("modify_account_status", "active", address)
	return err
}

func Rescan(address string, height int) error {
	_, err := command("rescan", strconv.Itoa(height), address)
	return err
}

func AcceptRequest(address ...string) error {
	if len(address) == 0 {
		return nil
	}
	args := []string{"accept_requests", "create"}
	args = append(args, address...)
	_, err := command(args...)
	return err
}

func RejectRequest(address ...string) error {
	if len(address) == 0 {
		return nil
	}
	args := []string{"reject_requests", "create"}
	args = append(args, address...)
	_, err := command(args...)
	return err
}
