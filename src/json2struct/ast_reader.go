package json2struct

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/hyperf/go-stringable/stringable"
)

type AstReader struct {
}

func NewAstReader() *AstReader {
	return &AstReader{}
}

func (r *AstReader) Run(code string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		// 查找类型声明
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// 检查是否是结构体
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		r.InspectObject(structType)

		return true
	})

	var buf bytes.Buffer
	err = format.Node(&buf, fset, node)
	if err != nil {
		panic(err)
	}

	return string(buf.Bytes())
}

func (r *AstReader) InspectObject(structType *ast.StructType) {
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				continue // 跳过嵌入字段
			}

			t, ok := field.Type.(*ast.ArrayType)
			if ok {
				r.InspectArray(t)
			}

			for _, ident := range field.Names {
				ident.Name = stringable.Studly(stringable.Snake(ident.Name))
			}
		}
	}
}

func (r *AstReader) InspectArray(t *ast.ArrayType) {
	t2, ok := t.Elt.(*ast.StructType)
	if ok {
		r.InspectObject(t2)
	}
}
