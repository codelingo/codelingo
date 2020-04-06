package sample

type W struct{}

func (W) Kill()       {}
func (W) Wait() error { return nil }
func (W) ohNo()       { go func() {}() }

type X struct{}

func (X) Krill()       {}
func (X) White() error { return nil }
func (X) ohYes()       { go func() {}() }
