package main

generatorName := cmdutil.GetFlagString(cmd, "generator")
generators := o.Generators("expose")
generator, found := generators[generatorName]
if {
        return cmdutil.UsageErrorf(cmd, "generator %q not found.", generatorName)
}
names := generator.ParamNames()

import "fmt"

func main() {
}
