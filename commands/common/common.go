// CONTINUE HERE: https://docs.docker.com/docker-hub/builds/

// this file houses functions common to multiple commands
package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"

	"github.com/lingo-reviews/lingo/tenet"
	"github.com/lingo-reviews/lingo/tenet/driver"
	"github.com/lingo-reviews/lingo/util"
)

const (
	DefaultTenetCfgPath = ".lingo"
	ConfigFile          = "services.yaml"
)

// TODO(waigani) The do server relies on this error message to know when to make a pull request.
var ErrMissingDotLingo = errors.New("No .lingo configuration found. Run `lingo init` to create a .lingo file in the current directory")

type CascadeDirection int

// Down and Both are intended to be used only with specific commands, and not exposed to the CLI user
const (
	CascadeNone CascadeDirection = iota // Only read config in the current working directory
	CascadeUp                           // Walk up parent directories and include found tenets
	CascadeDown                         // Search subdirectories recursively and include found tenets
	CascadeBoth                         // Combine CascadeUp and CascadeDown
)

type Config struct {
	Cascade     bool          `toml:"cascade"` // TODO: When switching from toml to viper, set this True by default
	Include     string        `toml:"include"`
	Template    string        `toml:"template"`
	TenetGroups []TenetGroup  `toml:"tenet_group"`
	allTenets   []TenetConfig // TODO(waigani) see comment in AllTenets
	// the root dir from which the config was built.
	buildRoot string
}

type TenetGroup struct {
	Name     string        `toml:"name"`
	Template string        `toml:"template"`
	Tenets   []TenetConfig `toml:"tenet"`
}

type TenetConfig struct {
	Name string `toml:"name"`

	// Name of the driver in use
	Driver string `toml:"driver"`

	// Registry server to pull the image from
	Registry string `toml:"registry"`

	// Tag server to pull the image from
	Tag string `toml:"tag"` // TODO(waigani) if we don't need this to pull, get rid of it.

	// Config options for tenet
	Options map[string]interface{} `toml:"options"`
}

// Provide a means to compare TenetConfig for equality as maps aren't inherently comparable.
func (c *TenetConfig) Hash() string {
	hash := strings.Join([]string{c.Name, c.Driver, c.Registry, c.Tag}, ",")
	for k, v := range c.Options {
		hash += k + v.(string)
	}
	return hash
}

// newTenet returns a new tenet.Tenet built from a cfg.
func NewTenet(tenetCfg TenetConfig) (tenet.Tenet, error) {
	return tenet.New(&driver.Base{
		Name:          tenetCfg.Name,
		Driver:        tenetCfg.Driver,
		Registry:      tenetCfg.Registry,
		Tag:           tenetCfg.Tag,
		ConfigOptions: tenetCfg.Options,
	})
}

func (c *Config) AllTenets() []TenetConfig {
	// TODO(waigani) quick work around. allTenets are the tenets built up
	// after a cascade read of cfgs. Rework things so it's clear that we are
	// either getting all tenets for one cfg or all tenets for all cfgs.
	// We may need a new struct AllCfg or something?
	s := seer{seen: map[string]bool{}}
	var tenets []TenetConfig
	for _, t := range c.allTenets {
		if !s.Seen(t.Name) {
			tenets = append(tenets, t)
		}
	}

	for _, g := range c.TenetGroups {
		for _, t := range g.Tenets {
			if !s.Seen(t.Name) {
				tenets = append(tenets, t)
			}
		}
	}
	return tenets
}

type seer struct {
	seen map[string]bool
}

func (s *seer) Seen(name string) (seen bool) {
	seen = s.seen[name]
	s.seen[name] = true
	return
}

func (c *Config) HasTenetGroup(name string) bool {
	for _, g := range c.TenetGroups {
		if g.Name == name {
			return true
		}
	}
	return false
}

func (c *Config) AddTenetGroup(name string) {
	if !c.HasTenetGroup(name) {
		c.TenetGroups = append(c.TenetGroups, TenetGroup{Name: name})
	}
}

func (c *Config) RemoveTenetGroup(name string) {
	var groups []TenetGroup
	for _, g := range c.TenetGroups {
		if g.Name != name {
			groups = append(groups, g)
		}
	}
	c.TenetGroups = groups
}

func (c *Config) AddTenet(t TenetConfig, group string) error {
	c.AddTenetGroup(group)
	g, err := c.FindTenetGroup(group)
	if err != nil {
		return errors.Trace(err)
	}

	// Return an error if a tenet of this name already exists in group
	for _, e := range g.Tenets {
		if e.Name == t.Name {
			return errors.Errorf("tenet %q already exists in group %q", t.Name, group)
		}
	}

	g.Tenets = append(g.Tenets, t)

	return nil
}

// FindTenetGroup returns a direct reference to the named group.
func (c *Config) FindTenetGroup(name string) (*TenetGroup, error) {
	for i := range c.TenetGroups {
		if c.TenetGroups[i].Name == name {
			return &c.TenetGroups[i], nil
		}
	}
	return nil, errors.Errorf("tenet group %q not found", name)
}

func (c *Config) RemoveTenet(name string, group string) error {
	g, err := c.FindTenetGroup(group)
	if err != nil {
		return errors.Trace(err)
	}
	var tenets []TenetConfig
	err = errors.Errorf("tenet %q not found", name)
	for _, t := range g.Tenets {
		if t.Name != name {
			tenets = append(tenets, t)
		} else {
			err = nil
		}
	}

	g.Tenets = tenets

	return err
}

// TODO(waigani) write osoutf, replace all fmt.Print

func OSErrf(format string, a ...interface{}) {
	format = fmt.Sprintf("error: %s\n", format)
	errStr := fmt.Sprintf(format, a...)
	log.Print(errStr)
	Stderr.Write([]byte(errStr))
	Exiter(1)
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

// Get a list of instantiated tenets from a config object.
func Tenets(ctx *cli.Context, cfg *Config) ([]tenet.Tenet, error) {
	var ts []tenet.Tenet
	for _, tenetCfg := range cfg.AllTenets() {
		tenet, err := tenet.New(&driver.Base{
			Name:          tenetCfg.Name,
			Driver:        tenetCfg.Driver,
			Registry:      tenetCfg.Registry,
			ConfigOptions: tenetCfg.Options,
		})
		if err != nil {
			message := fmt.Sprintf("could not create tenet '%s': %s", tenetCfg.Name, err.Error())
			return nil, errors.Annotate(err, message)
		}
		ts = append(ts, tenet)
	}

	return ts, nil
}

// TODO(waigani) make this externally extendable.
func FileExtFilterForLang(lang string) (regex, glob string) {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return ".*\\.go", "*.go"
	}
	return ".*", "*"
}

// Combine cascaded configuration files into a single config object.
func BuildConfig(startCfgPath string, cascadeDir CascadeDirection) (*Config, error) {
	if cascadeDir == CascadeNone {
		return ReadConfigFile(startCfgPath)
	}

	cfg := &Config{buildRoot: filepath.Dir(startCfgPath)}

	switch cascadeDir {
	case CascadeUp, CascadeDown:
		if err := buildConfigRecursive(startCfgPath, cascadeDir, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	case CascadeBoth:
		if err := buildConfigRecursive(startCfgPath, CascadeUp, cfg); err != nil {
			return nil, err
		}
		if err := buildConfigRecursive(startCfgPath, CascadeDown, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	return nil, errors.New("invalid cascade direction")
}

// Build up a config object by following directories up or down.
func buildConfigRecursive(cfgPath string, cascadeDir CascadeDirection, cfg *Config) error {
	cfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil
	}

	currentCfg, err := ReadConfigFile(cfgPath)
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

		// DEMOWARE: Need a better way to assign tenets to groups without duplication
		for _, g := range currentCfg.TenetGroups {
		DupeCheck:
			for _, t := range g.Tenets {
				for _, e := range cfg.allTenets {
					if e.Name == t.Name {
						continue DupeCheck
					}
				}
				cfg.AddTenet(t, g.Name)
				cfg.allTenets = append(cfg.allTenets, t)
			}
		}
		// Asign group properties
		for _, g := range currentCfg.TenetGroups {
			for i, cg := range cfg.TenetGroups {
				if cg.Name == g.Name {
					if g.Template != "" {
						cfg.TenetGroups[i].Template = g.Template
					}
				}
			}
		}
	} else if !os.IsNotExist(err) {
		// Just leave the current state of cfg on encountering an error
		log.Println("error reading file: %s", cfgPath)
		return nil
	}

	currentDir, filename := path.Split(cfgPath)
	switch cascadeDir {
	case CascadeUp:
		if currentDir == "/" || (currentCfg != nil && !currentCfg.Cascade) {
			return nil
		}

		parent := path.Dir(path.Dir(currentDir))

		if err := buildConfigRecursive(path.Join(parent, filename), cascadeDir, cfg); err != nil {
			return err
		}
	case CascadeDown:
		files, err := filepath.Glob(path.Join(currentDir, "*"))
		if err != nil {
			return nil
		}

		for _, f := range files {
			file, err := os.Open(f)
			if err != nil {
				log.Println("error reading file: %s", file)
				return nil
			}
			defer file.Close()
			if fi, err := file.Stat(); err == nil && fi.IsDir() {
				if err := buildConfigRecursive(path.Join(f, filename), cascadeDir, cfg); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("invalid cascade direction")
	}
	return nil
}

// Read a single config file into a config object.
func ReadConfigFile(cfgPath string) (*Config, error) {
	cfg := &Config{}

	// TODO(waigani) also support yaml and json
	_, err := toml.DecodeFile(cfgPath, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Write a config file in the current directory from a config object.
func WriteConfigFile(c *cli.Context, cfg *Config) error {
	fPath, err := TenetCfgPath(c)
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

// DesiredTenetCfgPath returns the tenet config path found in 1. local flag
// or 2. global flag. It falls back to returning DefaultTenetCfgPath
func DesiredTenetCfgPath(c *cli.Context) string {
	flgName := TenetCfgFlg.Long
	var cfgName string
	// 1. grab the config name from local flag
	if cfgName = c.String(flgName); cfgName != "" {
		return cfgName
	}
	if cfgName = c.GlobalString(flgName); cfgName != "" {
		return cfgName
	}
	// TODO(waigani) shouldn't need this - should fallback to default in flags.
	return DefaultTenetCfgPath
}

func TenetCfgPath(c *cli.Context) (string, error) {
	cfgPath := DesiredTenetCfgPath(c)
	return TenetCfgPathRecusive(cfgPath)
}

// TODO: TECHDEBT Check if commented code will be needed and prune as appropriate
// func tenetHome(c *cli.Context) string {
// 	home := c.GlobalString(lingoHomeFlg.Long)
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

// TenetCfgPathRecusive looks for a config file at cfgPath. If the config
// file name is equal to DefaultTenetCfgPath, the func recursively searches the
// parent directory until a file with that name is found. In the case that
// none is found "" is retuned.
func TenetCfgPathRecusive(cfgPath string) (string, error) {
	var err error
	cfgPath, err = filepath.Abs(cfgPath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		dir, file := path.Split(cfgPath)
		if file == DefaultTenetCfgPath {
			if dir == "/" {
				// we've reached the end of the line. Fall back to default:
				usr, err := user.Current()
				if err != nil {
					return "", err
				}

				lHome, err := util.LingoHome()
				if err != nil {
					return "", errors.Trace(err)
				}
				defaultTenets := path.Join(usr.HomeDir, lHome, DefaultTenetCfgPath)
				if _, err := os.Stat(defaultTenets); err != nil {
					return "", err
				}
				return defaultTenets, nil
			}
			parent := path.Dir(path.Dir(dir))
			return TenetCfgPathRecusive(parent + "/" + DefaultTenetCfgPath)
		}
		return "", err
	}
	return cfgPath, nil
}

func HasTenet(tenets []TenetConfig, imageName string) bool {
	for _, t := range tenets {
		if t.Name == imageName {
			return true
		}
	}
	return false
}

func (c *Config) HasTenet(name string) bool {
	return HasTenet(c.AllTenets(), name)
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

func ExactArgs(c *cli.Context, expected int) error {
	if l := len(c.Args()); l != expected {
		return errors.Errorf("expected %d argument(s), got %d", expected, l)
	}
	return nil
}

func MaxArgs(c *cli.Context, max int) error {
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
