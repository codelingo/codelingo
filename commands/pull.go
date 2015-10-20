package commands

import (
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/tenet"
)

var PullCMD = cli.Command{
	Name:  "pull",
	Usage: "pull tenet image(s) from docker hub",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  allFlg.String(),
			Usage: "pull all tenets in tenet.toml",
		}, cli.BoolFlag{
			Name:  updateFlg.String(),
			Usage: "look for a newer image on docker hub",
		},
	},
	Description: `

  pull takes one argument, the name of the docker image or a --all flag. If
  the flag is provided, 0 arguments are expected and all tenets in tenet.toml
  are pulled.

`[1:],
	Action: pull,
}

func pull(c *cli.Context) {
	all := c.Bool("all")
	expectedArgs := 1
	if all {
		expectedArgs = 0
	}
	if l := len(c.Args()); l != expectedArgs {
		oserrf("expected %d argument(s), got %d", expectedArgs, l)
		return
	}

	if all {
		if err := pullAll(c); err != nil {
			oserrf(err.Error())
		}
		return
	}
	if err := pullOne(c, c.Args().First()); err != nil {
		oserrf(err.Error())
	}
}

// Pull all tenets from config using assigned drivers.
func pullAll(c *cli.Context) error {
	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		return err
	}
	cfg, err := buildConfiguration(cfgPath, CascadeBoth)
	if err != nil {
		return err
	}

	ts, err := tenets(c, cfg)
	if err != nil {
		return err
	}

	for _, t := range ts {
		// TODO(waigani) don't return on err, collect errs and report at end
		err = t.InitDriver()
		if err != nil {
			return err
		}
		err = t.Pull()
		if err != nil {
			return err
		}
	}
	return nil
}

// Pull a tenet by name, assuming docker driver and default registry.
func pullOne(c *cli.Context, name string) error {
	// TODO: Add --driver, --registry flags for more info
	t, err := tenet.New(c, tenet.Config{Name: name})
	if err != nil {
		return errors.Trace(err)
	}

	err = t.InitDriver()
	if err != nil {
		return err
	}

	return t.Pull()
}

// // pullTenetImage attempts to pull the docker image <author>/<tenetName> from docker hub
// func pullTenetImage(c *cli.Context, imageName string) error {
// 	// // create new container from image
// 	// dockerArgs := []string{"run", "-i", "--name", containerName, imageName}
// 	// if haveContainer(containerName) {
// 	// 	if c.GlobalString(updateFlg.long) {

// 	// 	}
// 	// 	// start existing container
// 	// 	dockerArgs = []string{"start", "-i", containerName}
// 	// }

// 	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, "docker", dockerArgs...)
// 	if err != nil {
// 		log.Fatalf("Error running plugin: %s", err)
// 	}
// 	defer client.Close()

// 	// TODO(waigani) continue here, this is where we validate tenet. Write
// 	// validate as it's own command and pull it in here.

// 	// TODO(waigani) change plug to tenet. put in own package so it can be
// 	// imported by tenetseed also.
// 	p := plug{client}
// 	res, err := p.SayHi("master")
// 	if err != nil {
// 		log.Fatalf("error calling SayHi: %s", err)
// 	}
// 	log.Printf("Response from plugin: %q", res)

// 	res, err = p.SayBye("master")
// 	if err != nil {
// 		log.Fatalf("error calling SayBye: %s", err)
// 	}
// 	log.Printf("Response from plugin2: %q", res)
// 	return nil
// }

// legacy code for downloading a raw executable file.
// func downloadTenet(c *cli.Context, author, tenetName, container string) error {
// 	tenetPath := path.Join(tenetHome(c), author, tenetName)
// 	// TODO(waigani) check if file exists.
// 	// TODO(waigani) versioning.

// 	dir := path.Dir(tenetPath)
// 	if err := os.MkdirAll(dir, 0777); err != nil {
// 		return err
// 	}
// 	out, err := os.Create(tenetPath)
// 	defer out.Close()
// 	if err != nil {
// 		return errors.Trace(err)
// 	}

// 	// TODO(waigani) use container arg to optionally download tenet in docker container.
// 	url := lingoWeb("/tenets/" + author + "/" + tenetName + "/download")
// 	resp, err := http.Get(url.String())
// 	if err != nil {
// 		return errors.Trace(err)
// 	}
// 	defer resp.Body.Close()

// 	_, err = io.Copy(out, resp.Body)
// 	return err
// }
