package main

import proto "github.com/golang/protobuf/proto"

import "fmt"

type TestPB struct {
	CorrectName []byte `protobuf:"bytes,2,opt,name=correct_name,proto3" json:"correct_name,omitempty"`
	IncorrectName []byte `protobuf:"bytes,2,opt,name=incorrectname,proto3" json:"incorrect_name,omitempty"`
	AnotherIncorrectName []byte `protobuf:"bytes,2,opt,name=incorrect_name,proto3" json:"incorrect_name,omitempty"`
}