package main

import (
	"flag"

	"github.com/spf13/pflag"
)

type ConfigFlags struct {
	CertFile *string
}

const (
	flagGoodName = "good-name"
	flagBadName  = "bad_name"
)

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(f.CertFile, flagGoodName, *f.CertFile, "This is the description of a good flag")
	flags.StringVar(f.CertFile, flagBadName, *f.CertFile, "This is the description of a bad flag")

	var flagvar int
	flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")

	var ip *int = flags.Int(flagBadName, 1234, "This is the description of a bad flag")
}

