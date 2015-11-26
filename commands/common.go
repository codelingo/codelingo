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
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/dev/tenet/log"

	"github.com/lingo-reviews/lingo/tenet"
	"github.com/lingo-reviews/lingo/tenet/driver"
	"github.com/lingo-reviews/lingo/util"
)

const (
	// TODO(waigani) move this into util
	defaultTenetCfgPath = ".lingo"
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
	Cascade     bool          `toml:"cascade"` // TODO: When switching from toml to viper, set this True by default
	Include     string        `toml:"include"`
	Template    string        `toml:"template"`
	TenetGroups []TenetGroup  `toml:"tenet_group"`
	allTenets   []TenetConfig // TODO(waigani) see comment in AllTenets
	// the root dir from which the config was build.
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

// newTenet returns a new tenet.Tenet built from a cfg.
func newTenet(ctx *cli.Context, tenetCfg TenetConfig) (tenet.Tenet, error) {
	return tenet.New(ctx, &driver.Base{
		Name:          tenetCfg.Name,
		Driver:        tenetCfg.Driver,
		Registry:      tenetCfg.Registry,
		Tag:           tenetCfg.Tag,
		ConfigOptions: tenetCfg.Options,
	})
}

func (c *config) AllTenets() []TenetConfig {
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

func (c *config) HasTenetGroup(name string) bool {
	for _, g := range c.TenetGroups {
		if g.Name == name {
			return true
		}
	}
	return false
}

func (c *config) AddTenetGroup(name string) {
	c.TenetGroups = append(c.TenetGroups, TenetGroup{Name: name})
}

func (c *config) RemoveTenetGroup(name string) {
	var groups []TenetGroup
	for _, g := range c.TenetGroups {
		if g.Name != name {
			groups = append(groups, g)
			break
		}
	}
	c.TenetGroups = groups
}

func (c *config) AddTenet(t TenetConfig, group string) error {
	if !c.HasTenetGroup(group) {
		c.AddTenetGroup(group)
	}
	g, err := c.FindTenetGroup(group)
	if err != nil {
		return errors.Trace(err)
	}
	// TODO(waigani) use pointers to avoid all this update crap
	g.Tenets = append(g.Tenets, t)
	c.UpdateTenetGroup(g)
	return nil
}

// TODO(waigani) This shouldn't be needed, move to pointers
func (c *config) UpdateTenetGroup(group TenetGroup) {
	var groups []TenetGroup
	for _, g := range c.TenetGroups {
		if g.Name != group.Name {
			groups = append(groups, g)
		}
	}
	groups = append(groups, group)
	c.TenetGroups = groups
}

func (c *config) FindTenetGroup(name string) (TenetGroup, error) {
	for _, g := range c.TenetGroups {
		if g.Name == name {
			return g, nil
		}
	}
	return TenetGroup{}, errors.Errorf("tenet group %q not found", name)
}

func (c *config) RemoveTenet(name string, group string) error {
	g, err := c.FindTenetGroup(group)
	if err != nil {
		return errors.Trace(err)
	}
	var tenets []TenetConfig
	for _, t := range c.AllTenets() {
		if t.Name != name {
			tenets = append(tenets, t)
		}
	}

	g.Tenets = tenets
	c.UpdateTenetGroup(g)
	return nil
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
	log.Printf(format, a...)
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

// Get a list of instantiated tenets from a config object.
func tenets(ctx *cli.Context, cfg *config) ([]tenet.Tenet, error) {
	var ts []tenet.Tenet
	for _, tenetCfg := range cfg.AllTenets() {
		tenet, err := tenet.New(ctx, &driver.Base{
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
func fileExtFilterForLang(lang string) (regex, glob string) {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return ".*\\.go", "*.go"
	}
	return ".*", "*"
}

// reviewQueue returns a map of all tenets waiting to review, grouped by
// config. int is the total number of tenets waiting to review.
func reviewQueue(dir string) (map[*config][]TenetConfig, int, error) {
	totalTenets := 0
	queue := make(map[*config][]TenetConfig)

	// Starting with initial dir
	// - read config for that dir with CascadeUp (buildConfig will handle cascade=false)
	// - use found cfg.Include to find files in that dir
	// - insert cfg->files into map
	// - keep count of total files for channel buffer
	err := filepath.Walk(dir, func(relPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// fmt.Println("dir:", relPath) // TODO: put behind a debug flag
			// TODO(waigani) CONTINUE HERE. Add cfg.Root() string which
			// returns root dir of cfg. Also use this for lingo which.
			cfg, err := buildConfig(path.Join(relPath, defaultTenetCfgPath), CascadeUp)
			if err != nil {
				return err
			}

			for _, tenetCfg := range cfg.AllTenets() {
				totalTenets++
				queue[cfg] = append(queue[cfg], tenetCfg)
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

	cfg := &config{buildRoot: filepath.Dir(startCfgPath)}

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
				log.Println("error reading file: %s", file)
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

				lHome, err := util.LingoHome()
				if err != nil {
					return "", errors.Trace(err)
				}
				defaultTenets := path.Join(usr.HomeDir, lHome, defaultTenetCfgPath)
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

func hasTenet(tenets []TenetConfig, imageName string) bool {
	for _, t := range tenets {
		if t.Name == imageName {
			return true
		}
	}
	return false
}

func (c *config) HasTenet(name string) bool {
	return hasTenet(c.AllTenets(), name)
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
