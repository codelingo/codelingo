package main

import (
    "github.com/spf13/pflag"
)

type ConfigFlags struct {
    CertFile         *string
}

const (
    flagGoodName      = "good-name"
    flagBadName     = "bad_name"
)

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
    flags.StringVar(f.CertFile, flagGoodName, *f.CertFile, "This is the discription of a good flag")
    flags.StringVar(f.CertFile, flagBadName, *f.CertFile, "This is the discription of a bad flag")
}
