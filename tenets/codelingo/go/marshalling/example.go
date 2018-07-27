package main

import (
	"fmt"

	"encoding/json"

	"github.com/juju/errors"
)

// Suppose this were defined in the client
type ClientData struct {
	A string
}

type ServerData struct {
	A string
}

type Server struct {
	A string
}

func main() {
	err := SendData()
	if err != nil {
		panic(err.Error())
	}
}

func SendData() error {
	myServer := &Server{}

	data, err := MakeValidJSONRequest(myServer)
	if err != nil {
		return errors.Trace(err)
	}

	fmt.Printf("sending data %v", data)
	return nil
}

func MakeValidJSONRequest(s *Server) ([]byte, error) {
	cliStr := ServerData{
		A: "a fixed value",
	}

	data, err := json.Marshal(cliStr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return data, nil
}

func MakeInvalidJSONRequest(s *Server) ([]byte, error) {
	cliStr := ClientData{
		A: "a fixed value",
	}

	data, err := json.Marshal(cliStr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return data, nil
}
