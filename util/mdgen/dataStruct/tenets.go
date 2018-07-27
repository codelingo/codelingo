package dataStruct

type HubTenets struct {
	Owners map[string]TenetsOwner
}

type TenetsOwner struct {
	Name    string
	Bundles map[string]Bundle
}

type Bundle struct {
	Name   string
	Tenets map[string]Tenet
}

type Tenet struct {
	Name string
}

//This is a dumb workaround for now
type Data struct {
	Owner  string
	Bundle string
	Tenet  string
}
