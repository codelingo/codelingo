package somepackage

// Thing is an example
type Thing struct {
	Name   string
	Number int
}

// Example is an example
type Example struct {
	Count    int
	TheThing *Thing
}

// GetThing returns the Example's Thing
func (e *Example) GetThing() (thing *Thing) {
	thing = e.TheThing
	return
}

// GetThingAsValue returns the Example's Thing as value
func (e *Example) GetThingAsValue() (thing Thing) {
	thing = *e.TheThing
	return
}

// GetCount returns the count field of the Thing
func (e *Example) GetCount() (count int) {
	count = e.Count
	return
}

// GetAll returns all the Thing
func (e *Example) GetAll() (thing Thing, count int) {
	thing = *e.TheThing
	count = e.Count
	return
}

// Extract a Thing
func Extract(e *Example) (thing Thing) {
	thing = *e.TheThing
	return
}
