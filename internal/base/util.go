package base

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var (
	Dump    io.Writer
	config  map[string]any
	lastUpd time.Time = time.Unix(0, 0)
)

// const (
// 	configLoc = "/home/nodo/variables/config.json"
// 	configBak = "/home/nodo/variables/config.back.json"
// )

const (
	configLoc = "config.json"
	configBak = "config.json.bak"
)

const Device string = "Nodo"

func Bool(b bool) string {
	if b {
		return "TRUE"
	} else {
		return "FALSE"
	}
}

func SaveConfigFile() error {
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
	return nil
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
			// spew.Fdump(Dump, err)
			return err
		}
	}
	return nil
}

func SetConfig(key string, value any) {
	err := updateConfig()
	if err != nil {
		spew.Fdump(Dump, err)
		return
	}
	if config["config"] == nil {
		spew.Fdump(Dump, config)
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
	c := GetConfig()
	for _, s := range qr {
		switch (*c)[s].(type) {
		case map[string]any:
			*c = (*c)[s].(map[string]any)
		default:
			return (*c)[s]
		}
	}
	return c
}

func IsFirstBoot() bool {
	j, err := os.Stat("/home/nodo/variables/firstboot")
	spew.Fdump(Dump, j, err)
	if err == nil {
		return false
	}
	return !errors.Is(err, os.ErrNotExist)
}

func backup() error {
	_, err := os.Stat(configBak)
	if err != nil {
		return err
	}
	os.Remove(configBak)
	from, err := os.Open(configLoc)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.Open(configBak)
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
