package rewrite

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/juju/errors"
)

type SRCHunk struct {
	Filename    string
	StartOffset int64
	EndOffset   int64
	SRC         string
	Discard     bool
}

func toString(node interface{}) (string, error) {
	var b bytes.Buffer
	fset := &token.FileSet{}
	if err := format.Node(&b, fset, node); err != nil {
		return "", errors.Trace(err)
	}
	return b.String(), nil
}

func ClqlToSrc(clql string) (hunks map[string]string, err error) {
	return nil, errors.New("currently unavailable; WIP flow")

	// TODO(waigani) currently only supports one hunk / @rewrite.set decorator

	//hunks = make(map[string]string)
	//clqlQueryAST, err := inner.ParseString(clql)
	//if err != nil {
	//	//xxx.Print(clql)
	//	util.Logger.Infof("clql: %s", clql)
	//	return nil, errors.Trace(err)
	//}
	//
	//root := clqlQueryAST.GetFactTree()
	//
	//var finalASTs []interface{}
	//for _, child := range root.GetChildren() {
	//	childFact := child.GetFact()
	//	if namespace := childFact.ID.GetNamespace(); namespace != "go" {
	//		return nil, errors.Errorf("expected go namespace, got: %s", namespace)
	//
	//	}
	//	// build ast node from fact.
	//	parentKind := snakeCaseToCamelCase(childFact.ID.GetKind())
	//	astChildNode := astFactory[parentKind]()
	//
	//	// add all properties to node.
	//	for _, grandChild := range childFact.GetChildren() {
	//		processElm(astChildNode, childFact, grandChild, hunks)
	//	}
	//	// collect all nodes
	//	finalASTs = append(finalASTs, astChildNode)
	//}
	//
	//// print asts
	//for _, node := range finalASTs {
	//
	//	// funcDeclNode := node.(*ast.FuncDecl)
	//	// funcDeclNode.Body = new(ast.BlockStmt)
	//	// funcDeclNode.Type = new(ast.FuncType)
	//	//	funcDeclNode.Doc = new(ast.CommentGroup)
	//	// funcDeclNode.Recv = new(ast.FieldList)
	//	//	funcDeclNode.Doc = new(ast.CommentGroup)
	//
	//	//	xxx.Dump(node)
	//
	//	src, err := toString(node)
	//
	//	if err != nil {
	//		return nil, errors.Trace(err)
	//	}
	//
	//	// DEMOWARE(waigani) hardcoding to default key for now.
	//	hunks["_default"] = src
	//}
	//return hunks, nil
}

//func processElm(astParentNode interface{}, parent *inner.Fact, childElm *inner.Element, srcs map[string]string) {
//
//	parentKind := snakeCaseToCamelCase(parent.ID.GetKind())
//	var value interface{}
//	//var isFact bool
//	var name string
//
//	// TODO(waigani) transform name
//	if prop := childElm.GetProperty(); prop != nil {
//		name = snakeCaseToCamelCase(prop.Name)
//		value = propertyConvert(parentKind, name, prop)
//	}
//
//	if childFact := childElm.GetFact(); childFact != nil {
//
//		childFactKind := snakeCaseToCamelCase(childFact.ID.GetKind())
//		if factory, ok := astFactory[childFactKind]; ok {
//			value = factory()
//		} else {
//			util.Logger.Infof("need new fact: %s", childFactKind)
//		}
//
//		// process every child Elm of the fact, add properties to the node.
//		for _, grandChildElm := range childFact.GetChildren() {
//			processElm(value, childFact, grandChildElm, srcs)
//		}
//
//		// TODO(waigani) transform name
//		name = snakeCaseToCamelCase(childFact.GetID().GetKind())
//
//		// TODO(waigani) workout parent field name by matching fact types
//		// hardcoding for now
//		name = "Name"
//
//		//isFact = true
//	}
//
//	// xxx.Print("child: " + name)
//
//	// set the ast field
//	astParentNodeVal := reflect.ValueOf(astParentNode).Elem()
//
//	f := astParentNodeVal.FieldByName(name)
//	if f.IsValid() && f.CanSet() {
//		f.Set(reflect.ValueOf(value))
//	} else {
//		// xxx.Print("field not found: " + name + " for parent: " + parentKind)
//	}
//
//}

var astFactory = map[string]func() interface{}{
	"Decls": func() interface{} {
		return new([]ast.Decl)
	},
	"Decl": func() interface{} {
		return new(ast.Decl)
	},
	"FuncDecl": func() interface{} {
		return new(ast.FuncDecl)
	},
	"Ident": func() interface{} {
		return new(ast.Ident)
	},
	"AssignStmt": func() interface{} {
		return new(ast.AssignStmt)
	},
	"BinaryExpr": func() interface{} {
		return new(ast.BinaryExpr)
	},
	"IfStmt": func() interface{} {
		return new(ast.IfStmt)
	},
	"CallExpr": func() interface{} {
		return new(ast.CallExpr)
	},
}

//func propertyConvert(parent, field string, in *inner.Property) interface{} {
//
//	switch parent + "." + field {
//	case "Ident.NamePos":
//		return token.Pos(in.GetInt())
//	}
//
//	return in.GetString_()
//
//}

func snakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return

}
