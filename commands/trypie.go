package commands

// import (
// 	"log"
// 	"net/rpc"
// 	"net/rpc/jsonrpc"
// 	"os"

// 	"github.com/lingo-reviews/lingo/tenet"

// 	"github.com/codegangsta/cli"
// 	"github.com/natefinch/pie"
// )

// var TryPieCMD = cli.Command{
// 	Name:  "trypie",
// 	Usage: "try pie plugin",
// 	Action: func(c *cli.Context) {
// 		trypie()
// 	},
// }

// func trypie() {
// 	log.SetPrefix("[master log] ")

// 	t := tenet.New("lingoreviews/tenetseed")

// 	// create new container from image
// 	dockerArgs := []string{"run", "-i", "--name", containerName, image}
// 	if haveContainer(containerName) {
// 		// start existing container
// 		dockerArgs = []string{"start", "-i", containerName}
// 	}

// 	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, "docker", dockerArgs...)
// 	if err != nil {
// 		log.Fatalf("Error running plugin: %s", err)
// 	}
// 	defer client.Close()
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
// 	log.Printf("Response from plugin: %q", res)
// }

// type plug struct {
// 	client *rpc.Client
// }

// func (p plug) SayHi(name string) (result string, err error) {
// 	err = p.client.Call("Plugin.SayHi", name, &result)
// 	return result, err
// }

// func (p plug) SayBye(name string) (result string, err error) {
// 	err = p.client.Call("Plugin.SayBye", name, &result)
// 	return result, err
// }

// // CONTINUE HERE
// // get plugin working from within docker container
// // http://stackoverflow.com/questions/30653033/emulating-docker-run-using-the-golang-docker-api
// // --rm to remove container when finished
// // docker run -i -t --rm tenetseed
// //sh script:
// // #!/bin/bash

// // # use the local/shipyard-cli image (customized with pre-authentication)
// // sudo docker run --rm local/shipyard-cli "$@"
// // "run", "-i", "-a", "stdin", "tenetseed", "script", "-qc", `"/bin/bash"`, "/dev/null",
// // // docker ps -f name=tenetseed_container -a -q
