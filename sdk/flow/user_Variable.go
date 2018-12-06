package flow

// A UserVar allows users and other external agents to set variable values while a query
// is being executed.
// TODO: UserVar code is copied from the Platform
type UserVar struct {
	VarC         chan<- string
	Name         string
	DefaultValue string
}

// Set sets the value of the variable
func (s *UserVar) Set(val string) {
	s.VarC <- val
	close(s.VarC)
}

// SetAsDefault sets the variable to its default value
func (s *UserVar) SetAsDefault() {
	s.Set(s.DefaultValue)
}
