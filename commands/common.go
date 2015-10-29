// CONTINUE HERE: https://docs.docker.com/docker-hub/builds/

// this file houses functions common to multiple commands
package commands

import (
	"bytes"
	"encoding/json"
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
	"github.com/lingo-reviews/lingo/tenet/driver"
)

const (
	defaultTenetCfgPath = "tenet.toml"
)

type CascadeDirection int

// Down and Both are intended to be used only with specific commands, and not exposed to the CLI user
const (
	CascadeNone CascadeDirection = iota // Only read config in the current working directory
	CascadeUp                           // Walk up parent directories and include found tenets
	CascadeDown                         // Search subdirectories recursively and include found tenets
	CascadeBoth                         // Combine CascadeUp and CascadeDown
)

type config struct {
	Include  string         `toml:"include"`
	Template string         `toml:"template"`
	Cascade  bool           `toml:"cascade"` // TODO: When switching from toml to viper, set this True by default
	Tenets   []tenet.Config `toml:"tenet"`
}

// stderr is a var for mocking in tests
var stderr io.Writer = os.Stderr

// exiter is a var for mocking in tests
var exiter = func(code int) {
	os.Exit(code)
}

// TODO(waigani) write osoutf, replace all fmt.Print

func oserrf(format string, a ...interface{}) {
	format = fmt.Sprintf("error: %s\n", format)
	fmt.Fprintf(stderr, format, a...)
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

// TODO: Better solution for logging: optionally to file, -v flag etc.
func log(format string, a ...interface{}) {
	fmt.Fprintf(stderr, format+"\n", a...)
}

// Get a list of instantiated tenets from a config object.
func tenets(ctx *cli.Context, cfg *config) ([]tenet.Tenet, error) {
	var ts []tenet.Tenet
	for _, tenetData := range cfg.Tenets {
		tenet, err := tenet.New(ctx, tenetData)
		if err != nil {
			message := fmt.Sprintf("could not create tenet '%s': %s", tenetData.Name, err.Error())
			return nil, errors.Annotate(err, message)
		}
		ts = append(ts, tenet)
	}

	return ts, nil
}

func allConfigs(dir string) (map[*config][]string, int, error) {
	totalTenets := 0
	queue := make(map[*config][]string)

	// Starting with initial dir
	// - read config for that dir with CascadeUp (buildConfig will handle cascade=false)
	// - use found cfg.Include to find files in that dir
	// - insert cfg->files into map
	// - keep count of total files for channel buffer
	err := filepath.Walk(dir, func(relPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// fmt.Println("dir:", relPath) // TODO: put behind a debug flag
			cfg, err := buildConfig(path.Join(relPath, defaultTenetCfgPath), CascadeUp)
			if err != nil {
				return err
			}

			// TODO: Use cfg.Include glob/regex?
			files, err := filepath.Glob(path.Join(relPath, "*.go"))
			if err != nil { // Non-fatal
				return nil
			}

			fileList := []string{}
			for _, f := range files {
				file, err := os.Open(f)
				if err != nil { // Non-fatal
					break
				}
				if fi, err := file.Stat(); err == nil && !fi.IsDir() {
					// fmt.Println("adding", f) // TODO: put behind a debug flag
					fileList = append(fileList, f)
				}
			}

			if len(fileList) > 0 {
				totalTenets += len(cfg.Tenets)

				queue[cfg] = fileList
			}
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return queue, totalTenets, nil
}

// Combine cascaded configuration files into a single config object.
func buildConfig(startCfgPath string, cascadeDir CascadeDirection) (*config, error) {
	if cascadeDir == CascadeNone {
		return readConfigFile(startCfgPath)
	}

	cfg := &config{}

	switch cascadeDir {
	case CascadeUp, CascadeDown:
		buildConfigRecursive(startCfgPath, cascadeDir, cfg)
		return cfg, nil
	case CascadeBoth:
		buildConfigRecursive(startCfgPath, CascadeUp, cfg)
		buildConfigRecursive(startCfgPath, CascadeDown, cfg)
		return cfg, nil
	}

	return nil, errors.New("invalid cascade direction")
}

// Build up a config object by following directories up or down.
func buildConfigRecursive(cfgPath string, cascadeDir CascadeDirection, cfg *config) {
	cfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return
	}

	currentCfg, err := readConfigFile(cfgPath)
	if err == nil {
		// Add the non-tenet properties - always when cascading down, otherwise
		// only if not already specified
		// TODO: Use reflection here to avoid forgotten values?
		if cascadeDir == CascadeDown || cfg.Include == "" {
			cfg.Include = currentCfg.Include
		}
		if cascadeDir == CascadeDown || cfg.Template == "" {
			cfg.Template = currentCfg.Template
		}

	DupeCheck:
		for _, t := range currentCfg.Tenets {
			// Don't duplicate tenets
			// TODO: handle case of same tenet but different options
			for _, existing := range cfg.Tenets {
				if existing.Name == t.Name {
					continue DupeCheck
				}
			}
			cfg.Tenets = append(cfg.Tenets, t)
		}
	} else if !os.IsNotExist(err) {
		// Just leave the current state of cfg on encountering an error
		log("error reading file: %s", cfgPath)
		return
	}

	currentDir, filename := path.Split(cfgPath)
	switch cascadeDir {
	case CascadeUp:
		if currentDir == "/" || (currentCfg != nil && !currentCfg.Cascade) {
			return
		}

		parent := path.Dir(path.Dir(currentDir))

		buildConfigRecursive(path.Join(parent, filename), cascadeDir, cfg)
	case CascadeDown:
		files, err := filepath.Glob(path.Join(currentDir, "*"))
		if err != nil {
			return
		}

		for _, f := range files {
			file, err := os.Open(f)
			if err != nil {
				log("error reading file: %s", file)
				return
			}
			if fi, err := file.Stat(); err == nil && fi.IsDir() {
				buildConfigRecursive(path.Join(f, filename), cascadeDir, cfg)
			}
		}
	default:
		oserrf("invalid cascade direction")
	}
}

// Read a single config file into a config object.
func readConfigFile(cfgPath string) (*config, error) {
	cfg := &config{}

	// TODO(waigani) also support yaml and json
	_, err := toml.DecodeFile(cfgPath, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Write a config file in the current directory from a config object.
func writeConfigFile(c *cli.Context, cfg *config) error {
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

// parseOptions returns a map of tenet names to Options from the command line.
func parseOptions(c *cli.Context) (map[string]driver.Options, error) {
	commandOptions := map[string]driver.Options{}
	// Parse command line specified options
	if commandOptionsJson := c.String("options"); commandOptionsJson != "" {
		err := json.Unmarshal([]byte(commandOptionsJson), &commandOptions)
		if err != nil {
			return nil, err
		}
	}
	return commandOptions, nil
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

// TODO: TECHDEBT Check if commented code will be needed and prune as appropriate
// func tenetHome(c *cli.Context) string {
// 	home := c.GlobalString(lingoHomeFlg.long)
// 	return path.Join(home, "tenets")
// }

// writeFileAll writes the given file and any missing dirs in it's path.
// func writeFileAll(filePath string, data []byte, perm os.FileMode) error {
// 	dir := path.Dir(filePath)
// 	if err := os.MkdirAll(dir, perm); err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(filePath, data, perm)
// }

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

func hasTenet(cfg *config, imageName string) bool {
	for _, config := range cfg.Tenets {
		if config.Name == imageName {
			return true
		}
	}
	return false
}

// Return a string representation of a CascadeDirection
func (c CascadeDirection) String() string {
	switch c {
	case CascadeNone:
		return "none"
	case CascadeUp:
		return "up"
	case CascadeDown:
		return "down"
	case CascadeBoth:
		return "both"
	}
	return "unknown"
}

// func authorAndNameFromArg(arg string) (author, tenetName string, err error) {
// 	parts := strings.Split(arg, "/")

// 	// TODO(waigani) when publishing tenet, don't allow ":" char in name.
// 	if len(parts) != 2 {
// 		return "", "", errors.New(`expected argument, to be of form "<author>/<tenet>"`)
// 	}
// 	return parts[0], parts[1], nil
// }

func exactArgs(c *cli.Context, expected int) error {
	if l := len(c.Args()); l != expected {
		return errors.Errorf("expected %d argument(s), got %d", expected, l)
	}
	return nil
}

func maxArgs(c *cli.Context, max int) error {
	if l := len(c.Args()); l > max {
		return errors.Errorf("expected up to %d argument(s), got %d", max, l)
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
