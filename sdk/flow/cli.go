package flow

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codelingo/lingo/app/util"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
)

type DecoratedResult struct {

	// Ctx are the flags and args parsed from the decorated fact
	Ctx *cli.Context

	// Payload is the struct returned form the flow server for the matched
	// fact.
	Payload proto.Message
}

type flowRunner struct {
	cliCMD       *CLICommand
	decoratorCMD *DecoratorCommand
}

type CLICommand struct {
	cli.Command
	Request func(*cli.Context) (chan proto.Message, chan error, func(), error)
}

type DecoratorCommand struct {
	cli.Command
	ConfirmDecorated func(*cli.Context, proto.Message) (bool, error)
}

func NewFlow(cliCMD *CLICommand, decoratorCMD *DecoratorCommand) *flowRunner {

	return &flowRunner{
		cliCMD:       cliCMD,
		decoratorCMD: decoratorCMD,
	}
}

func (f *flowRunner) Run() (decoratedResultc chan *DecoratedResult, err error) {
	defer func() {
		if err != nil {
			close(decoratedResultc)
		}
	}()

	resultc, errc, cancel, err := f.RunCLI()
	if err != nil {
		return nil, errors.Trace(err)
	}

	// If user is manually confirming reviews, set a long timeout.
	timeout := time.After(time.Hour * 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Wait()
		close(decoratedResultc)
	}()

l:
	for {
		select {
		case err, ok := <-errc:
			if !ok {
				errc = nil
				break
			}

			util.Logger.Debugf("Review error: %s", errors.ErrorStack(err))
			return nil, errors.Trace(err)
		case result, ok := <-resultc:
			if !ok {
				resultc = nil
				break
			}
			wg.Add(1)

			// TODO(waigani) this is brittle and expects the result struct to have a
			// DecoratorOptions field. We need to refactor the Flow server to return a
			// tuple of [<decorator>, <result>].
			decorator := reflect.ValueOf(result).FieldByName("DecoratorOptions").String()
			keep, err := f.ConfirmDecorated(decorator, result)
			if err != nil {
				cancel()
				return nil, errors.Trace(err)
			}
			if keep {

				go func(string, proto.Message) {
					defer wg.Done()
					ctx, err := NewCtx(f.decoratorCMD.Command, strings.Split(decorator, " ")[1:])
					if err != nil {
						cancel()
						util.Logger.Fatalf("error getting decorated context: %q", err)
						return
					}
					decoratedResultc <- &DecoratedResult{
						Ctx:     ctx,
						Payload: result,
					}
				}(decorator, result)

			}

		case <-timeout:
			cancel()
			return nil, errors.New("timed out waiting for issue")
		}
		if resultc == nil && errc == nil {
			break l
		}
	}
	wg.Done()

	return decoratedResultc, nil
}

// TODO(waigani) move this to codelingo/sdk/flow
func (f *flowRunner) RunCLI() (chan proto.Message, chan error, func(), error) {
	cmd := *f.cliCMD
	ctx, err := NewCtx(cmd.Command, os.Args[1:])
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}
	return cmd.Request(ctx)
}

func (f *flowRunner) ConfirmDecorated(decorator string, payload proto.Message) (bool, error) {
	cmd := *f.decoratorCMD
	ctx, err := NewCtx(cmd.Command, strings.Split(decorator, " "))
	if err != nil {
		return false, errors.Trace(err)
	}

	return cmd.ConfirmDecorated(ctx, payload)
}

func NewCtx(cmd cli.Command, input []string) (*cli.Context, error) {

	fSet := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	for _, flag := range cmd.Flags {
		flag.Apply(fSet)
	}

	if err := fSet.Parse(input); err != nil {
		return nil, errors.Trace(err)
	}
	if err := normalizeFlags(cmd.Flags, fSet); err != nil {
		return nil, errors.Trace(err)
	}

	ctx := cli.NewContext(nil, fSet, nil)

	if ctx.Bool("debug") {
		err := util.SetDebugLogger()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

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
