package flow

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/codelingo/lingo/app/util"
	"github.com/common-nighthawk/go-figure"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

type DecoratedResult struct {

	// Ctx are the flags and args parsed from the decorated fact
	Ctx *cli.Context

	// Payload is the struct returned form the flow server for the matched
	// fact.
	Payload proto.Message
}

type flowRunner struct {
	cliCtx           *cli.Context
	cliApp           *CLIApp
	decoratorApp     *DecoratorApp
	decoratedResultc chan *DecoratedResult
	errc             chan error
}

type CLIApp struct {
	cli.App
	Request func(*cli.Context) (chan proto.Message, <-chan *UserVar, chan error, func(), error)

	// Help data
	Tagline string
}

type DecoratorApp struct {
	cli.App
	ConfirmDecorated func(*cli.Context, proto.Message) (bool, error)
	SetUserVar       func(*UserVar)

	// Help info
	DecoratorUsage   string
	DecoratorExample string
}

// NewFlow creates a flowRunner from a CLIApp.
// For flows that query the flow server it overrides the action with command function, which
// runs the cliApp's Request function and listens on the channels it returns.
// TODO: move away from flowRunner model by explicitly defining the command function as the
// action in the flow.
func NewFlow(cliApp *CLIApp, decoratorApp *DecoratorApp) *flowRunner {

	// setBaseApp overrides Action with the help action. We don't want this
	// and either want nil or the action that was already set on
	// CLIApp.Action.
	action := cliApp.Action
	setBaseApp(cliApp)
	cliApp.Action = action

	fRunner := &flowRunner{
		cliApp:       cliApp,
		decoratorApp: decoratorApp,
		errc:         make(chan error),
	}
	if fRunner.decoratorApp != nil {
		fRunner.decoratorApp.Setup()
	}
	fRunner.setHelp()

	if fRunner.cliApp.Action == nil {

		// If Action is set, decoratedResultc is not used, nor is it closed.
		// Only set it if we're calling the internal flowRunner.action.
		fRunner.decoratedResultc = make(chan *DecoratedResult)
		fRunner.cliApp.Action = fRunner.action
	}
	return fRunner
}

func (f *flowRunner) setHelp() {
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		defer func() {

			// decoratedResultc is nil if the Flow is using it's own action.
			if f.decoratedResultc != nil {
				close(f.decoratedResultc)
			}
		}()
		figure.NewFigure("codelingo", "larry3d", false).Print()
		printHelp(w, CLIAPPHELPTMP, f.cliApp)
		if f.decoratorApp != nil {
			printHelp(w, DECAPPHELPTMP, f.decoratorApp)
		}
		printHelp(w, INFOTMP, f.cliApp)
	}

	f.cliApp.Commands = append(f.cliApp.Commands, cli.Command{
		Name: "_genJSONHelp",
		Action: func(ctx *cli.Context) {

			defer func() {
				// decoratedResultc is nil if the Flow is using it's own action.
				if f.decoratedResultc != nil {
					close(f.decoratedResultc)
				}
			}()

			var cliFlags, decFlags []string
			for _, flag := range f.cliApp.Flags {
				cliFlags = append(cliFlags, flag.String())
			}
			if f.decoratorApp != nil {
				for _, flag := range f.decoratorApp.Flags {
					decFlags = append(decFlags, flag.String())
				}
			}

			helpData := helpData{
				Name:         f.cliApp.Name,
				Tagline:      f.cliApp.Tagline,
				Description:  f.cliApp.Usage,
				LastCompiled: f.cliApp.Compiled.String(),
				Version:      f.cliApp.Version,

				Options: decFlags,
			}

			if f.decoratorApp != nil {
				helpData.Decorator = decHelp{
					Description: f.decoratorApp.Usage,
					// Example: f.decoratorApp.Example,
					Options: decFlags,
				}
			}

			b, err := json.Marshal(helpData)
			if err != nil {
				fmt.Println(err)
				return
			}

			outPath := ctx.Args()[0]
			if err := ioutil.WriteFile(outPath, b, 0644); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Success! JSON Help written to:", outPath)

		},
	})

}

func setBaseApp(cliApp *CLIApp) {

	app := CLIApp{
		App:     *cli.NewApp(),
		Request: cliApp.Request,
	}

	// base settings
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lingo-file",
			Usage: "the Tenet `FILE` to run rewrite over. If the flag is not set, codelingo.yaml files are read from the branch being rewritten.",
			//	Destination: &language,
		},
	}
	app.Compiled = time.Now()
	app.EnableBashCompletion = true

	// user settings
	app.Name = cliApp.Name
	app.Tagline = cliApp.Tagline
	app.Usage = cliApp.Usage
	app.Flags = append(app.Flags, cliApp.Flags...)
	app.Action = cliApp.Action
	app.Version = cliApp.Version

	*cliApp = app
}

// TODO(waigani) incorrect usage func

// Run runs the CLI app. We assume that cliApp's uses f.command as its action to query
// the flow server and stream back decorated/confirmed results on f.decoratedResultc.
// Special cliApps that don't use the flow server have their own custom actions, in which
// case the chans returned from here will not be closed.
func (f *flowRunner) Run() (chan *DecoratedResult, chan error) {

	go func() {
		defer close(f.errc)
		if err := f.cliApp.Run(os.Args); err != nil {
			f.errc <- err
		}
	}()

	return f.decoratedResultc, f.errc
}

func (f *flowRunner) action(ctx *cli.Context) {
	if err := f.command(ctx); err != nil {
		util.Logger.Debugw(errors.ErrorStack(err))
		util.FatalOSErr(err)
	}

}

func (f *flowRunner) command(ctx *cli.Context) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer func() {
		wg.Wait()
		close(f.decoratedResultc)
	}()
	defer wg.Done()

	resultc, userVarC, errc, cancel, err := f.RunCLI()
	if err != nil {
		return errors.Trace(err)
	}

	// If user is manually confirming results, set a long timeout.
	timeout := time.After(time.Hour * 1)
l:
	for {
		select {
		case err, ok := <-errc:
			if !ok {
				errc = nil
				break
			}

			util.Logger.Debugf("Result error: %s", errors.ErrorStack(err))
			return errors.Trace(err)
		case v, ok := <-userVarC:
			if !ok {
				userVarC = nil
				break
			}

			f.SetUserVar(v)
		case result, ok := <-resultc:
			if !ok {
				resultc = nil
				break
			}

			// TODO(waigani) this is brittle and expects the result struct to have a
			// DecoratorOptions field. We need to refactor the Flow server to return a
			// tuple of [<decorator>, <result>].
			decorator := reflect.Indirect(reflect.ValueOf(result)).FieldByName("DecoratorOptions").String()

			var keep bool
			if ctx.Bool("keep-all") {
				keep = true
			} else {
				keep, err = f.ConfirmDecorated(decorator, result)
				if err != nil {
					cancel()
					return errors.Trace(err)
				}
			}

			if keep {
				wg.Add(1)
				go func(string, proto.Message) {
					defer wg.Done()
					ctx, err := NewCtx(&f.decoratorApp.App, strings.Split(decorator, " ")[1:]...)
					if err != nil {
						cancel()
						f.errc <- err
						util.Logger.Fatalf("error getting decorated context: %q", err)
						return
					}
					util.Logger.Debug("sending result to Flow...")
					f.decoratedResultc <- &DecoratedResult{
						Ctx:     ctx,
						Payload: result,
					}
					util.Logger.Debug("...result sent to Flow")

				}(decorator, result)

			}

		case <-timeout:
			cancel()
			return errors.New("timed out waiting for issue")
		}
		if resultc == nil && errc == nil && userVarC == nil {
			break l
		}
	}

	return
}

func (f *flowRunner) CliCtx() (*cli.Context, error) {
	if f.cliCtx == nil {
		ctx, err := NewCtx(&f.cliApp.App, os.Args[1:]...)
		if err != nil {
			return nil, errors.Trace(err)
		}
		f.cliCtx = ctx
	}
	return f.cliCtx, nil
}

func (f *flowRunner) RunCLI() (chan proto.Message, <-chan *UserVar, chan error, func(), error) {
	ctx, err := f.CliCtx()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)
	}
	if ctx.Bool("debug") {
		err := util.SetDebugLogger()
		if err != nil {
			return nil, nil, nil, nil, errors.Trace(err)
		}
	}

	return f.cliApp.Request(ctx)
}

func (f *flowRunner) ConfirmDecorated(decorator string, payload proto.Message) (bool, error) {
	ctx, err := NewCtx(&f.decoratorApp.App, strings.Split(decorator, " ")...)
	if err != nil {
		return false, errors.Trace(err)
	}

	return f.decoratorApp.ConfirmDecorated(ctx, payload)
}

func (f *flowRunner) SetUserVar(userVar *UserVar) {
	f.decoratorApp.SetUserVar(userVar)
}

func NewCtx(app *cli.App, input ...string) (*cli.Context, error) {

	fSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	for _, flag := range app.Flags {
		flag.Apply(fSet)
	}

	if err := fSet.Parse(input); err != nil {
		return nil, errors.Trace(err)
	}
	if err := normalizeFlags(app.Flags, fSet); err != nil {
		return nil, errors.Trace(err)
	}

	ctx := cli.NewContext(nil, fSet, nil)

	return ctx, nil
}

func normalizeFlags(flags []cli.Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		parts := strings.Split(f.GetName(), ",")
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if !visited[name] {
				copyFlag(name, ff, set)
			}
		}
	}
	return nil
}

func copyFlag(name string, ff *flag.Flag, set *flag.FlagSet) {
	switch ff.Value.(type) {
	case *cli.StringSlice:
	default:
		set.Set(name, ff.Value.String())
	}
}
