package tenet

import "github.com/lingo-reviews/tenets/go/dev/api"

type Issue struct {
	api.Issue
	TenetName     string `json:"tenetName,omitempty"`
	Discard       bool   `json:"discard,omitempty"`
	DiscardReason string `json:"discardReason,omitempty"`
}
