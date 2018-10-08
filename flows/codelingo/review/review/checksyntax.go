package review

// TODO(waigani) refactor below to be called on lingo review. use `lingo
// review --dry-run` if user only wants to check syntax.

// package main

// import (
// 	"fmt"
// 	. "github.com/codelingo/clql/preprocessor"
// 	"github.com/fatih/color"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// )

// func main() {
// 	if len(os.Args) == 1 {
// 		fmt.Println("Codelingo syntax checker\nUsage: ./checksyntax [dot-lingo-file-path] [extra] [verbose]")
// 		return
// 	}

// 	verbose := 1
// 	if len(os.Args) >= 3 {
// 		verbose = 2
// 	}

// 	path, err := filepath.Abs(os.Args[1])
// 	if err == nil {
// 		fmt.Println("Reading the .lingo file at " + path)
// 		filebytes, err := ioutil.ReadFile(path)
// 		if err != nil {
// 			fmt.Println("Failed to open and read file " + path)
// 			return
// 		}

// 		m, er := Process(string(filebytes), verbose)

// 		if er == nil {
// 			color.Cyan("=== LINGO SYNTAX AND SEMANTICS PASS ALL CHECKS ===")
// 			color.Cyan("%s", path)
// 		} else {
// 			color.Red("=== LINGO SYNTAX AND SEMANTIC CHECK FAILED ===")
// 			color.Red("%s", path)
// 			color.Yellow(er.Error())
// 		}

// 		if len(os.Args) == 4 {
// 			// extra verbose, dump the map
// 			color.Cyan("\nYAML MAP FROM PREPROCESSOR:")
// 			fmt.Println(m)
// 		}

// 	} else {
// 		fmt.Println("Could not find file: " + os.Args[1])
// 		return
// 	}

// 	return
// }
