package commands

import (
	"os"
	"os/exec"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/util"
)

var PushCMD = cli.Command{
	Name:  "push",
	Usage: "push a tenet to its registry",
	Description: `
	
Call "lingo push" in the root directory of the source code of a tenet.
It will look for a .lingofile with instructions on how where to push your tenet and what to call it.
For example:

owner = "lingoreviews"
name = "simpleseed"
registry = "hub.docker.com"

At this stage only docker tenets can be pushed.
`[1:],
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "push every tenet found in every subdirectory",
		},
	},
	Action: push,
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

func push(ctx *cli.Context) {

	// TODO(waigani) a lot of this code is to support multiple drivers. It's
	// copied from build. It's left so we can easily add remote. We just hard
	// code to docker for now.

	if len(ctx.Args()) > 0 {
		common.OSErrf("push takes no arguments")
	}

	// drivers := ctx.Args()
	// if len(drivers) == 0 {
	// 	drivers = allDrivers
	// }
	drivers := []string{"docker"}

	lingofiles, err := getLingofiles(ctx)
	if err != nil {
		common.OSErrf(err.Error())
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
			msg := "The %s tenet is being pushed ..."
			if fn > 0 {
				msg = "All %s tenets are being pushed ..."
			}
			util.Printf(msg, driver)
			dw.bar.Start()
		}(dw)
	}

	// start building all tenets
	for _, f := range lingofiles {
		go func(f string) {
			errs := pushTenet(f, driverWaits)
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
		util.Printf("push error: %v\n", err)
	}

}

func pushTenet(lingofilePath string, waits map[string]*driverWait) []error {
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
			if err := cfg.pushBinary(); err != nil {
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
			if err := cfg.pushDocker(); err != nil {
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
			if err := cfg.pushRemote(); err != nil {
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

func (cfg *lingofileCfg) pushBinary() error {
	// dw := cfg.Binary.dw
	// dw.start.Done()
	// dw.bar.Increment()
	// dw.end.Done()

	return nil
}

func (cfg *lingofileCfg) pushDocker() error {
	dw := cfg.Binary.dw
	dw.start.Done()
	// defer dw.end.Done()
	defer dw.bar.Increment()

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

func (cfg *lingofileCfg) pushRemote() error {
	// dw := cfgRemote.dw
	// dw.start.Done()
	// dw.bar.Increment()
	// dw.end.Done()

	return errors.New("not implemented")
}
