package base

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"regexp"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"github.com/mergestat/timediff"
)

var (
	Dump    io.Writer
	config  map[string]any
	lastUpd time.Time = time.Unix(0, 0)
)

const addrPattern = "^4[0-9A-Za-z]{94}$"

// const (
// 	configLoc = "/home/nodo/variables/config.json"
// 	configBak = "/home/nodo/variables/config.back.json"
// )

const (
	configLoc = "config.json"
	configBak = "config.json.bak"
)

type ConfigSavedMsg struct{}

type ErrorMsg struct {
	err error
}

func Bool(b bool) string {
	if b {
		return "TRUE"
	} else {
		return "FALSE"
	}
}

func GetBool(s string) bool {
	return s == "TRUE"
}

func SaveConfigFile() tea.Msg {
	err := backup()
	if err != nil {
		spew.Fprintf(Dump, "SaveConfig: %v", err)
		return err
	}
	if config != nil {
		j, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return err
		}
		os.WriteFile(configLoc, j, 0o644)
	}
	return ConfigSavedMsg{}
}

func loadConfigFile() (map[string]any, error) {
	var j map[string]any
	data, err := os.ReadFile(configLoc)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func updateConfig() error {
	if config == nil || time.Since(lastUpd) > 2*time.Second {
		var err error
		config, err = loadConfigFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func SetMpayConfig(key string, value any) {
	err := updateConfig()
	if err != nil {
		spew.Fdump(Dump, err)
		return
	}
	if config["config"] == nil ||
		config["config"].(map[string]any)["moneropay"] == nil {
		return
	}
	config["config"].(map[string]any)["moneropay"].(map[string]any)[key] = value
	SaveConfigFile()
}

func SetBanlistConfig(key string, value bool) {
	err := updateConfig()
	if err != nil {
		spew.Fdump(Dump, err)
		return
	}
	if config["config"] == nil ||
		config["config"].(map[string]any)["banlists"] == nil {
		return
	}
	config["config"].(map[string]any)["banlists"].(map[string]any)[key] = Bool(value)
	SaveConfigFile()
}

func SetConfig(key string, value any) {
	err := updateConfig()
	if err != nil {
		spew.Fdump(Dump, err)
		return
	}
	if config["config"] == nil {
		return
	}
	switch value := value.(type) {
	case bool:
		config["config"].(map[string]any)[key] = Bool(value)
	default:
		config["config"].(map[string]any)[key] = value
	}
	SaveConfigFile()
}

func GetConfig() *(map[string]any) {
	err := updateConfig()
	if err != nil {
		spew.Fdump(Dump, err)
	}
	return &config
}

func GetVal(qr ...string) any {
	c := (*GetConfig())["config"].(map[string]any)
	for _, s := range qr {
		switch c[s].(type) {
		case map[string]any:
			c = c[s].(map[string]any)
		default:
			if c[s] == "TRUE" || c[s] == "FALSE" {
				return GetBool(c[s].(string))
			}
			return c[s]
		}
	}
	return c
}

func IsFirstBoot() bool {
	_, err := os.Stat("/home/nodo/variables/firstboot")
	err, ok := err.(*fs.PathError)
	return ok
}

func backup() error {
	_, err := os.Stat(configBak)
	if err != nil {
		os.Remove(configBak)
	}
	from, err := os.Open(configLoc)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.Create(configBak)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func UnixTimeRelative(unix int64) string {
	return timediff.TimeDiff(time.Unix(unix, 0))
}

func UnixTime(unix int64) string {
	t := time.Unix(unix, 0)
	tz, _ := GetVal("timezone").(string)
	loc, _ := time.LoadLocation(tz)
	return t.In(loc).Format("2 January 2006 15:04:05")
}

func ValidateAddr(addr string) bool {
	match, _ := regexp.MatchString(addrPattern, addr)
	return match
}
