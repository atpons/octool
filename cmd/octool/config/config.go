package config

import (
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var Value Config

const (
	configFile = ".octool.toml"
)

type Config struct {
	OpenConnect     OpenConnect
	OcProxy         OcProxy
	StraightForward StraightForward
}

type OpenConnect struct {
	Juniper bool
	User    string
	Host    string
}

type OcProxy struct {
	Port string
}

type StraightForward struct {
	Enabled      bool
	Port         string
	User         string
	Host         string
	IdentityFile string
}

func Load(f string) {
	if f == "" {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		fp := filepath.Join(u.HomeDir, configFile)
		_, err = toml.DecodeFile(fp, &Value)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := toml.DecodeFile(f, &Value)
		if err != nil {
			panic(err)
		}
	}
}

func BuildOpenConnectOpts(ocproxyFileName string) []string {
	var opts []string
	opts = append(opts, "--script-tun")
	opts = append(opts, "-s")
	opts = append(opts, ocproxyFileName)
	if Value.OpenConnect.Juniper {
		opts = append(opts, "--juniper")
	}
	opts = append(opts, "--passwd-on-stdin")
	opts = append(opts, "-u")
	opts = append(opts, Value.OpenConnect.User)
	opts = append(opts, Value.OpenConnect.Host)
	return opts
}

func BuildStraightForwardOpts() []string {
	var opts []string
	opts = append(opts, "-i")
	opts = append(opts, Value.StraightForward.IdentityFile)
	opts = append(opts, "-u")
	opts = append(opts, Value.StraightForward.User)
	opts = append(opts, "-h")
	opts = append(opts, Value.StraightForward.Host)
	opts = append(opts, "-p")
	opts = append(opts, Value.StraightForward.Port)
	opts = append(opts, "-ocproxy")
	opts = append(opts, "localhost:"+Value.OcProxy.Port)
	return opts
}
