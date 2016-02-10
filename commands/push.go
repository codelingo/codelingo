package commands

import (
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/commands/common"
	bt "github.com/lingo-reviews/lingo/tenet/build"
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

	lingofiles, err := getLingoFiles(ctx.Bool("all"))
	if err != nil {
		common.OSErrf(err.Error())
		return
	}
	fn := len(lingofiles)
	if fn == 0 {
		util.Printf("WARNING tenet not built. No %s found.\n", common.Lingofile)
		return
	}

	// set up drivers, waitgroups and progress bars
	allErrsc := make(chan error)
	allWaits := sync.WaitGroup{}
	dn := len(drivers)
	allWaits.Add(dn)

	driverWaits := make(map[string]*bt.DriverWait, dn)
	for _, driver := range drivers {
		bar := pb.New(fn)
		bar.Prefix(driver + " ")

		start := &sync.WaitGroup{}
		start.Add(fn)
		end := &sync.WaitGroup{}
		end.Add(fn)
		driverWaits[driver] = &bt.DriverWait{
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
		go func(dw *bt.DriverWait) {
			dw.Start.Wait()
			msg := "The %s tenet is being pushed ..."
			if fn > 0 {
				msg = "All %s tenets are being pushed ..."
			}
			util.Printf(msg, driver)
			dw.Bar.Start()
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

func pushTenet(lingofilePath string, waits map[string]*bt.DriverWait) []error {
	cfg, err := bt.ReadLingoFile(lingofilePath)
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
			if err := cfg.PushBinary(); err != nil {
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
			if err := cfg.PushDocker(); err != nil {
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
			if err := cfg.PushRemote(); err != nil {
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
