package build

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/cheggaaa/pb"
	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/util"
)

// TODO(waigani) This is not the right place for this struct.
type DriverWait struct {
	// waits until all tenets for this driver have started
	Start *sync.WaitGroup
	// waits until all tenets for this driver are built
	End *sync.WaitGroup
	// shows progress on tenets being built
	Bar *pb.ProgressBar
}

// Run builds each tenet for each lingofile from source for each driver
// specified in drivers.
func Run(drivers []string, lingofiles ...string) error {

	// set up drivers, waitgroups and progress bars
	allErrsc := make(chan error)
	allWaits := sync.WaitGroup{}
	dn := len(drivers)
	allWaits.Add(dn)

	driverWaits := make(map[string]*DriverWait, dn)
	fn := len(lingofiles)
	for _, driver := range drivers {
		bar := pb.New(fn)
		bar.Prefix(driver + " ")

		start := &sync.WaitGroup{}
		start.Add(fn)
		end := &sync.WaitGroup{}
		end.Add(fn)
		driverWaits[driver] = &DriverWait{
			Start: start,
			End:   end,
			Bar:   bar,
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
		go func(dw *DriverWait) {
			dw.Start.Wait()
			msg := "The %s tenet has started to build ..."
			if fn > 0 {
				msg = "All %s tenets have started to build ..."
			}
			util.Printf(msg, driver)
			dw.Bar.Start()
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
	return nil
}

func buildTenet(lingofilePath string, waits map[string]*DriverWait) []error {
	cfg, err := ReadLingoFile(lingofilePath)
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
			cfg.Binary.DW = binaryW
			if err := cfg.buildBinary(); err != nil {
				errsc <- err
			}
			binaryW.End.Done()
		}
	}()

	go func() {
		dockerW, ok := waits["docker"]
		if ok {
			defer wg.Done()
			cfg.Docker.DW = dockerW
			if err := cfg.buildDocker(); err != nil {
				errsc <- err
			}
			dockerW.End.Done()
		}
	}()

	go func() {
		remoteW, ok := waits["remote"]
		if ok {
			defer wg.Done()
			cfg.Remote.DW = remoteW
			if err := cfg.buildRemote(); err != nil {
				errsc <- err
			}
			remoteW.End.Done()
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
	DW *DriverWait
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
func ReadLingoFile(cfgPath string) (*lingofileCfg, error) {
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

func (cfg *lingofileCfg) PushBinary() error {
	return nil
}

func (cfg *lingofileCfg) PushDocker() error {
	dw := cfg.Binary.DW
	dw.Start.Done()
	// defer dw.end.Done()
	defer dw.Bar.Increment()

	// TODO(waigani) use client with auth
	// c, err := util.DockerClient()
	// if err != nil {
	// 	return errors.Trace(err)
	// }
	// opts := docker.PushImageOptions{
	// 	Registry:     cfg.Registry,
	// 	OutputStream: os.Stdout,
	// 	Name:         cfg.Owner + "/" + cfg.Name,
	// }
	// auth := docker.AuthConfiguration{
	// 	Username: cfg.Owner,
	// }
	// return c.PushImage(opts, auth)

	// TODO(waigani) we're not setting registry
	cmd := exec.Command("docker", "push", cfg.Owner+"/"+cfg.Name)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	// cmd.Stdin = os.Stdin
	return cmd.Start()
}

func (cfg *lingofileCfg) PushRemote() error {
	// dw := cfgRemote.dw
	// dw.start.Done()
	// dw.bar.Increment()
	// dw.end.Done()

	return errors.New("not implemented")
}
