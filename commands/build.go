package commands

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/util"
)

const (
	lingofile = ".lingofile"
)

var allDrivers = []string{"binary", "docker"}

var BuildCMD = cli.Command{
	Name:  "build",
	Usage: "build a tenet from source",
	Description: `
	
Call "lingo build" in the root directory of the source code of a tenet.
It will look for a .lingofile with instructions on how to build your tenet.
For example:

language = "go"
owner = "lingoreviews"
name = "simpleseed"

[binary]
  build=false

[docker]
  overwrite_dockerfile=true

You can specify which driver to build. If no arguments are supplied, lingo
will try to build every driver. 
`[1:],
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "build every tenet found in every subdirectory",
		},
		// TODO(waigani) implement
		// cli.BoolFlag{
		// 	Name:  "watch",
		// 	Usage: "watch for any changes and attempt to rebuild the tenet",
		// },
	},
	Action: build,
	BashComplete: func(c *cli.Context) {
		// This will complete if no args are passed
		if len(c.Args()) > 0 {
			return
		}

		for _, d := range allDrivers {
			util.Println(d)
		}
	},
}

var c int

type driverWait struct {
	// waits until all tenets for this driver have started
	start *sync.WaitGroup
	// waits until all tenets for this driver are built
	end *sync.WaitGroup
	// shows progress on tenets being built
	bar *pb.ProgressBar
}

// TODO(waigani) add more user output
func build(ctx *cli.Context) {
	drivers := ctx.Args()
	if len(drivers) == 0 {
		drivers = allDrivers
	}

	lingofiles, err := getLingofiles(ctx)
	if err != nil {
		oserrf(err.Error())
		return
	}
	fn := len(lingofiles)
	if fn == 0 {
		util.Printf("WARNING tenet not built. No %s found.\n", lingofile)
		return
	}

	// set up drivers, waitgroups and progress bars
	allErrsc := make(chan error)
	allWaits := sync.WaitGroup{}
	dn := len(drivers)
	allWaits.Add(dn)

	driverWaits := make(map[string]*driverWait, dn)
	for _, driver := range drivers {
		bar := pb.New(fn)
		bar.Prefix(driver + " ")

		start := &sync.WaitGroup{}
		start.Add(fn)
		end := &sync.WaitGroup{}
		end.Add(fn)
		driverWaits[driver] = &driverWait{
			start: start,
			end:   end,
			bar:   bar,
		}

		errsc := make(chan error)
		go func(wg *sync.WaitGroup, errc chan error) {
			wg.Wait()
			bar.FinishPrint("Success! All " + driver + " tenets built.")

			close(errsc)
			allWaits.Done()

		}(end, errsc)

		go func(errsc chan error) {
			for err := range errsc {
				allErrsc <- err
			}

		}(errsc)
	}

	// print progress bars to user
	for driver, dw := range driverWaits {
		go func(dw *driverWait) {
			dw.start.Wait()
			msg := "The %s tenet has started to build ..."
			if fn > 0 {
				msg = "All %s tenets have started to build ..."
			}
			util.Printf(msg, driver)
			dw.bar.Start()
		}(dw)
	}

	// start building all tenets
	for _, f := range lingofiles {
		go func(f string) {
			errs := buildTenet(f, driverWaits)
			for _, err := range errs {
				allErrsc <- err
			}
		}(f)
	}

	// wait for all results and print errors at the end
	go func() {
		allWaits.Wait()
		close(allErrsc)
	}()

	var errs []error
	for err := range allErrsc {
		errs = append(errs, err)
	}
	// util.NoPrint = false
	for _, err := range errs {
		util.Printf("build error: %v\n", err)
	}

}

func getLingofiles(ctx *cli.Context) (lingofiles []string, err error) {
	if ctx.Bool("all") {
		// util.NoPrint = true

		lingofiles, err = allLingofiles(".")
		if err != nil {
			return nil, err
		}
	} else {
		if _, err := os.Stat(lingofile); err == nil {
			lingofiles = []string{lingofile}
		}
	}
	return
}

func allLingofiles(rootDir string) (lingofiles []string, err error) {
	err = filepath.Walk(rootDir, func(relPath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(relPath, lingofile) {
			lingofiles = append(lingofiles, relPath)
		}
		return nil
	})
	return
}

func buildTenet(lingofilePath string, waits map[string]*driverWait) []error {
	cfg, err := readLingoFile(lingofilePath)
	if err != nil {
		return []error{errors.Trace(err)}
	}

	wg := sync.WaitGroup{}
	errsc := make(chan error)
	wg.Add(len(waits))
	go func() {
		wg.Wait()
		close(errsc)
	}()

	go func() {
		// TODO(waigani) use a map of driver to builder interface
		binaryW, ok := waits["binary"]
		if ok {
			defer wg.Done()
			cfg.Binary.dw = binaryW
			if err := cfg.buildBinary(); err != nil {
				errsc <- err
			}
			binaryW.end.Done()
		}
	}()

	go func() {
		dockerW, ok := waits["docker"]
		if ok {
			defer wg.Done()
			cfg.Docker.dw = dockerW
			if err := cfg.buildDocker(); err != nil {
				errsc <- err
			}
			dockerW.end.Done()
		}
	}()

	go func() {
		remoteW, ok := waits["remote"]
		if ok {
			defer wg.Done()
			cfg.Remote.dw = remoteW
			if err := cfg.buildRemote(); err != nil {
				errsc <- err
			}
			remoteW.end.Done()
		}
	}()

	var errs []error
	for err := range errsc {
		errs = append(errs, err)
	}

	return errs
}

// TODO(waigani) move this into lingofile
type lingofileCfg struct {
	Language string         `toml:"language"`
	Owner    string         `toml:"owner"`
	Name     string         `toml:"name"`
	Registry string         `toml:"registry"`
	Binary   binaryBuildCfg `toml:"binary"`
	Docker   dockerBuildCfg `toml:"docker"`
	Remote   remoteBuildCfg `toml:"remote"`
}

type baseBuildCfg struct {
	// the tenet owner
	Owner string
	// the name of the tenet
	Name string
	// working directory
	dir string
	// a waitgroup and progress bar
	dw *driverWait
}

func (cfg *lingofileCfg) buildBinary() error {
	binCfg := cfg.Binary
	lang := strings.ToLower(cfg.Language)

	switch lang {
	case "go", "golang":
		return binCfg.BuildGo()
	}
	return errors.Errorf("unknown language %q", lang)
}

func (cfg *lingofileCfg) buildDocker() error {
	dockerCfg := cfg.Docker
	lang := strings.ToLower(cfg.Language)

	switch lang {
	case "go", "golang":
		return dockerCfg.BuildGo()
	}
	return errors.Errorf("unknown language %q", lang)
}

type remoteBuildCfg struct {
	Build bool `toml:"build"`
	*baseBuildCfg
}

func (cfg *lingofileCfg) buildRemote() error {
	remoteCfg := cfg.Remote
	if !remoteCfg.Build {
		return nil
	}

	return errors.New("not implemented")
}

// Read a .lingofile into a cfg object.
func readLingoFile(cfgPath string) (*lingofileCfg, error) {
	cfg := &lingofileCfg{}

	// TODO(waigani) also support yaml and json
	_, err := toml.DecodeFile(cfgPath, cfg)
	if err != nil {
		return nil, err
	}

	// set working dir
	absPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil, errors.Trace(err)
	}
	dir := filepath.Dir(absPath)

	base := &baseBuildCfg{
		Name:  cfg.Name,
		Owner: cfg.Owner,
		dir:   dir,
	}

	cfg.Binary.baseBuildCfg = base
	cfg.Docker.baseBuildCfg = base
	cfg.Remote.baseBuildCfg = base

	return cfg, nil
}
