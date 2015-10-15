// CONTINUE HERE: https://docs.docker.com/docker-hub/builds/

// this file houses functions common to multiple commands
package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/tenet"
)

const (
	defaultTenetCfgPath = "tenet.toml"
)

type tenetCfg struct {
	Configs []tenet.Config `toml:"tenet"`
}

// stderr is a var for mocking in tests
var stderr io.Writer = os.Stderr

// exiter is a var for mocking in tests
var exiter = func(code int) {
	os.Exit(code)
}

// TODO(waigani) write osoutf, replace all fmt.Print

func oserrf(format string, a ...interface{}) {
	fmt.Fprintf(stderr, "error: "+format, a...)
	exiter(1)
}

func lingoWeb(uri string) url.URL {
	return url.URL{
		Scheme: "http",
		// Opaque   string    // encoded opaque data
		// User     *Userinfo // username and password information
		Host: "localhost:3000", // will be lingo.reviews in prod
		Path: uri,
		// RawQuery string // encoded query values, without '?'
		// Fragment string // fragment for references, without '#'
	}
}

func tenetCfgs(c *cli.Context) []tenet.Config {
	cfg, err := readTenetCfgFile(c)
	if err != nil {
		oserrf("reading config file: %s", err.Error())
		return nil
	}
	return cfg.Configs
}

func tenets(c *cli.Context) []tenet.Tenet {
	cfgs := tenetCfgs(c)
	var ts []tenet.Tenet
	for _, cfg := range cfgs {
		tenet, err := tenet.New(c, cfg)
		if err != nil {
			oserrf("could not create tenet '%s': %s", cfg.Name, err.Error())
			return nil
		}
		ts = append(ts, tenet)
	}

	return ts
}

// pathToCfg can be either a local file system path or a URL.
func readTenetCfgFile(c *cli.Context) (*tenetCfg, error) {
	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		return nil, errors.Trace(err)
	}
	cfg := &tenetCfg{}

	// TODO(waigani) also support yaml and json
	_, err = toml.DecodeFile(cfgPath, cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return cfg, nil
}

// pathToCfg can be either a local file system path or a URL.
func writeTenetCfgFile(c *cli.Context, cfg *tenetCfg) error {
	fPath, err := tenetCfgPath(c)
	if err != nil {
		return errors.Trace(err)
	}
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err = enc.Encode(cfg)
	if err != nil {
		return errors.Trace(err)
	}

	return ioutil.WriteFile(fPath, buf.Bytes(), 0644)
}

// desiredTenetCfgPath returns the tenet config path found in 1. local flag
// or 2. global flag. It falls back to returning defaultTenetCfgPath
func desiredTenetCfgPath(c *cli.Context) string {
	flgName := tenetCfgFlg.long
	var cfgName string
	// 1. grab the config name from local flag
	if cfgName = c.String(flgName); cfgName != "" {
		return cfgName
	}
	if cfgName = c.GlobalString(flgName); cfgName != "" {
		return cfgName
	}
	// TODO(waigani) shouldn't need this - should fallback to default in flags.
	return defaultTenetCfgPath
}

func tenetCfgPath(c *cli.Context) (string, error) {
	cfgPath := desiredTenetCfgPath(c)
	return tenetCfgPathRecusive(cfgPath)
}

func userHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func defaultLingoHome() string {
	home, err := userHome()
	if err != nil {
		oserrf(err.Error())
		return ""
	}
	return path.Join(home, ".lingo")
}

func tenetHome(c *cli.Context) string {
	home := c.GlobalString(lingoHomeFlg.long)
	return path.Join(home, "tenets")
}

// writeFileAll writes the given file and any missing dirs in it's path.
func writeFileAll(filePath string, data []byte, perm os.FileMode) error {
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, perm); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, data, perm)
}

// tenetCfgPathRecusive looks for a config file at cfgPath. If the config
// file name is equal to defaultTenetCfgPath, the func recursively searches the
// parent directory until a file with that name is found. In the case that
// none is found "" is retuned.
func tenetCfgPathRecusive(cfgPath string) (string, error) {
	var err error
	cfgPath, err = filepath.Abs(cfgPath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		dir, file := path.Split(cfgPath)
		if file == defaultTenetCfgPath {
			if dir == "/" {
				// we've reached the end of the line. Fall back to default:
				usr, err := user.Current()
				if err != nil {
					return "", err
				}
				defaultTenets := path.Join(usr.HomeDir, ".lingo", defaultTenetCfgPath)
				if _, err := os.Stat(defaultTenets); err != nil {
					return "", err
				}
				return defaultTenets, nil
			}

			parent := path.Dir(path.Dir(dir))
			return tenetCfgPathRecusive(parent + "/" + defaultTenetCfgPath)
		}
		return "", err
	}
	return cfgPath, nil
}

func hasTenet(cfg *tenetCfg, imageName string) bool {
	for _, config := range cfg.Configs {
		if config.Name == imageName {
			return true
		}
	}
	return false
}

// func authorAndNameFromArg(arg string) (author, tenetName string, err error) {
// 	parts := strings.Split(arg, "/")

// 	// TODO(waigani) when publishing tenet, don't allow ":" char in name.
// 	if len(parts) != 2 {
// 		return "", "", errors.New(`expected argument, to be of form "<author>/<tenet>"`)
// 	}
// 	return parts[0], parts[1], nil
// }

func expectedArgs(c *cli.Context, expected int) error {
	if l := len(c.Args()); l != expected {
		return errors.Errorf("expected %d argument(s), got %d", expected, l)
	}
	return nil
}

// func CheckerConfig(configFile string) *Config {
// 	c := &Config{}

// 	// TODO(waigani) also support yaml and json
// 	_, err := toml.DecodeFile(configFile, c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return c
// }

// func ConfigFromURL(URL string) (*Config, error) {
// 	resp, err := http.Get(URL)
// 	if err != nil {
// 		return nil, errors.Trace(err)
// 	}
// 	var c Config
// 	defer resp.Body.Close()
// 	_, err = toml.DecodeReader(resp.Body, &c)
// 	// _, err := toml.DecodeFile("./pedantic.toml", &c)
// 	if err != nil {
// 		return nil, errors.Trace(err)
// 	}
// 	return &c, nil
// }
