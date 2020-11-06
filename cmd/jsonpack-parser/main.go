package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/hcl/strconv"
)

type importInfo struct {
	path string
	dir  string
}

type pkgFileInfo struct {
	node    *ast.File
	name    string
	types   map[string]ast.Node
	imports map[string]*importInfo
}

type packageInfo struct {
	node      *ast.Package
	name      string
	files     []*pkgFileInfo
	bigEndian bool
}

type SchemaDef struct {
	Type       string                `json:"type"`
	Properties map[string]*SchemaDef `json:"properties,omitempty"`
	Items      *SchemaDef            `json:"items,omitempty"`
	Order      []string              `json:"order,omitempty"`
}

type paramSet struct {
	srcDir     string
	structName string
	outputFile string
	bigEndian  bool
	pretty     bool
}

var params paramSet

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: jsonpack-parser -s <PACKAGE_DIR> -n <STRUCT_NAME> [options]\n")
		fmt.Fprintf(os.Stderr, "jsonpack-parser is a helper tool to generate schema defintion text from struct in package.\n\n")
		fmt.Fprintf(os.Stderr, "Available options:\n\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&params.srcDir, "s", "", "the package directory to parse")
	flag.StringVar(&params.structName, "n", "", "the struct name in package which want to generate schema definition")
	flag.StringVar(&params.outputFile, "o", "", "output file name for schema definition, program will output to stdout if this option not specified")
	flag.BoolVar(&params.bigEndian, "b", false, "generate number type with big-endian byte order, default is little-endian if this option not set")
	flag.BoolVar(&params.pretty, "p", false, "output JSON text of schema definition with indentation")

	flag.Parse()

	if params.srcDir == "" || params.structName == "" {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n-s and -n options are required.\n")
		os.Exit(2)
	}

	var err error
	params.srcDir, err = filepath.Abs(params.srcDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(params.outputFile) > 0 {
		params.outputFile, err = filepath.Abs(params.outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func (p *packageInfo) buildBultinType(data string) (*SchemaDef, bool) {
	const uintSize = 32 << (^uint(0) >> 32 & 1)

	switch data {
	case "string":
		return &SchemaDef{Type: data}, true
	case "bool":
		return &SchemaDef{Type: "boolean"}, true
	case "int":
		if uintSize == 64 {
			return &SchemaDef{Type: p.numEndian("int64")}, true
		} else {
			return &SchemaDef{Type: p.numEndian("int32")}, true
		}
	case "uint":
		if uintSize == 64 {
			return &SchemaDef{Type: p.numEndian("uint64")}, true
		} else {
			return &SchemaDef{Type: p.numEndian("uint32")}, true
		}
	case "byte":
		return &SchemaDef{Type: "uint8"}, true
	case "int8", "uint8":
		return &SchemaDef{Type: data}, true
	case "int16", "int32", "int64", "uint16", "uint32", "uint64", "float32", "float64":
		return &SchemaDef{Type: p.numEndian(data)}, true
	default:
		return nil, false
	}
}

func (p *packageInfo) parseNode(fileInfo *pkgFileInfo, node ast.Node) (*SchemaDef, error) {
	switch node := node.(type) {
	case *ast.Ident:
		data, ok := p.buildBultinType(node.Name)
		if ok {
			return data, nil
		}

		elemNode := p.lookupNode(node.Name)
		if elemNode == nil {
			return nil, fmt.Errorf("%s type doesn't exist", node.Name)
		}
		return p.parseNode(fileInfo, elemNode)

	case *ast.StructType:
		return p.parseStruct(fileInfo, node)

	case *ast.ArrayType:
		arrSt, err := p.parseNode(fileInfo, node.Elt)
		if err != nil {
			return nil, err
		}
		st := SchemaDef{Type: "array", Items: arrSt}
		return &st, nil

	case *ast.StarExpr:
		return p.parseNode(fileInfo, node.X)

	case *ast.MapType:
		return nil, fmt.Errorf("map type is not allowed in struct")

	case *ast.SelectorExpr:
		modName := node.X.(*ast.Ident).Name
		modType := node.Sel.Name
		extPkg, ok := fileInfo.imports[modName]
		if !ok {
			return nil, fmt.Errorf("package %s not found", modName)
		}

		extSch, err := parsePackageNode(extPkg.dir, modType)
		if err != nil {
			return nil, err
		}
		return extSch, nil
	}
	return nil, nil
}

func parseFieldName(name string, tag *ast.BasicLit) string {
	if tag == nil {
		return name
	}

	value := tag.Value[1 : len(tag.Value)-1]
	jsonTag, ok := lookupTag(value, "json")
	if !ok {
		return name
	}

	tagParts := strings.Split(jsonTag, ",")
	if tagParts[0] != "" {
		return tagParts[0]
	}
	return name
}

func (p *packageInfo) parseStruct(fileInfo *pkgFileInfo, stAst *ast.StructType) (*SchemaDef, error) {
	st := SchemaDef{Type: "object"}
	st.Properties = make(map[string]*SchemaDef)
	st.Order = make([]string, 0, len(stAst.Fields.List))

	for _, field := range stAst.Fields.List {
		var fieldName string
		if len(field.Names) > 0 && field.Names[0] != nil {
			fieldName = parseFieldName(field.Names[0].Name, field.Tag)
		}

		// embedded struct field
		if fieldName == "" {
			sch, err := p.parseNode(fileInfo, field.Type)
			if err != nil {
				return nil, err
			}
			switch typ := field.Type.(type) {
			case *ast.Ident:
				fieldName = parseFieldName(typ.Name, field.Tag)
				err = mergeEmbedField(fieldName, &st, sch)
				if err != nil {
					return nil, err
				}
			case *ast.SelectorExpr:
				fieldName = parseFieldName(typ.Sel.Name, field.Tag)
				err = mergeEmbedField(fieldName, &st, sch)
				if err != nil {
					return nil, err
				}
			}
			// fmt.Printf("  embedded field Name: %v Type: (%T)%+v Tag: %v\n", fieldName, field.Type, field.Type, field.Tag)
		} else {
			// fmt.Printf("  field Name: %v Type: (%T)%v Tag: %v\n", fieldName, field.Type, field.Type, field.Tag)
			// process field
			fieldProp, err := p.parseNode(fileInfo, field.Type)
			if err != nil {
				return nil, err
			}
			// skip invalid field
			if fieldProp == nil {
				continue
			}

			st.Properties[fieldName] = fieldProp
			st.Order = append(st.Order, fieldName)
		}
	}

	return &st, nil
}

func (p *packageInfo) lookupNode(name string) ast.Node {
	for _, file := range p.files {
		node := lookupNodeInFile(file, name)
		if node != nil {
			return node
		}
	}
	return nil
}

// borrow code from StructTag.Lookup function in golang reflect/type package
func lookupTag(tag string, key string) (value string, ok bool) {
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]
		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}

func lookupNodeInFile(fileInfo *pkgFileInfo, name string) ast.Node {
	node, ok := fileInfo.types[name]
	if ok {
		return node
	}
	return nil
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (p *packageInfo) numEndian(data string) string {
	if p.bigEndian {
		return data + "be"
	}
	return data + "le"
}

func (p *packageInfo) getNode(name string) (*pkgFileInfo, ast.Node) {
	for _, file := range p.files {
		if node, ok := file.types[name]; ok {
			return file, node
		}
	}
	return nil, nil
}

func parsePackageStruct(pkgDir string, name string) (*SchemaDef, error) {
	return _parsePackage(pkgDir, name, true)
}

func parsePackageNode(pkgDir string, name string) (*SchemaDef, error) {
	return _parsePackage(pkgDir, name, false)
}

func _parsePackage(pkgDir string, name string, findStruct bool) (*SchemaDef, error) {
	fset := token.NewFileSet() // positions are relative to fset
	// Parse src but stop after processing the imports.
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse package fail, error: %v\n", err)
		return nil, err
	}

	pkgInfo, err := newPackageInfo(pkgs, pkgDir, params.bigEndian)
	if err != nil {
		return nil, err
	}

	fileInfo, node := pkgInfo.getNode(name)
	if fileInfo == nil {
		return nil, fmt.Errorf("%s node not found", name)
	}
	if findStruct {
		_, ok := node.(*ast.StructType)
		if !ok {
			return nil, fmt.Errorf("struct %s not found", name)
		}
	}
	// parse  node here
	sch, err := pkgInfo.parseNode(fileInfo, node)
	if err != nil {
		return nil, fmt.Errorf("Parse type %s fail, err: %v", name, err)
	}
	return sch, nil
}

func mergeEmbedField(name string, dst *SchemaDef, embed *SchemaDef) error {
	// skip invalid field
	if embed == nil {
		return nil
	}
	if embed.Type == "object" {
		for _, fieldName := range embed.Order {
			if _, ok := dst.Properties[fieldName]; ok {
				return fmt.Errorf("ambiguous field %s in embedded struct", fieldName)
			}
			dst.Properties[fieldName] = embed.Properties[fieldName]
			dst.Order = append(dst.Order, fieldName)
		}
	} else {
		dst.Properties[name] = embed
		dst.Order = append(dst.Order, name)
	}
	return nil
}

func newPackageInfo(pkgs map[string]*ast.Package, pkgDir string, bigEndian bool) (*packageInfo, error) {
	// TODO: leverage cache
	pkgInfo := packageInfo{}
	pkgInfo.bigEndian = bigEndian
	if len(pkgs) > 1 {
		return nil, fmt.Errorf("package dir %s contains more than on packages", pkgDir)
	}
	for name, node := range pkgs {
		pkgInfo.name = name
		pkgInfo.node = node
	}

	for fileName, fileAst := range pkgInfo.node.Files {
		fileInfo := pkgFileInfo{
			node:    fileAst,
			name:    fileName,
			types:   make(map[string]ast.Node),
			imports: make(map[string]*importInfo),
		}

		// parse import
		for _, importNode := range fileAst.Imports {
			var importKey string
			pkgImport := importInfo{}
			pkgImport.path, _ = strconv.Unquote(importNode.Path.Value)
			// import alias name
			if importNode.Name != nil {
				importKey = importNode.Name.Name
			}

			ctx := build.Default
			ctx.Dir = pkgDir
			pkg, err := ctx.Import(pkgImport.path, "", 0)
			if err != nil {
				return nil, fmt.Errorf("Import package %s fail", pkg.ImportPath)
			}

			// no import alias here
			if importNode.Name == nil {
				importKey = pkg.Name
			}

			pkgImport.dir = pkg.Dir
			fileInfo.imports[importKey] = &pkgImport
		}
		// search type declarations
		ast.Inspect(fileAst, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				fileInfo.types[node.Name.Name] = node.Type
			}
			return true
		})

		pkgInfo.files = append(pkgInfo.files, &fileInfo)
	}
	return &pkgInfo, nil
}

func main() {
	var err error
	if !isDirectory(params.srcDir) {
		fmt.Fprintf(os.Stderr, "pacakge dir '%s' is not exist\n", params.srcDir)
		os.Exit(1)
	}

	sch, err := parsePackageStruct(params.srcDir, params.structName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	var schText []byte
	if params.pretty {
		schText, err = json.MarshalIndent(sch, "", "  ")
	} else {
		schText, err = json.Marshal(sch)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if params.outputFile != "" {
		f, fErr := os.OpenFile(params.outputFile, os.O_WRONLY|os.O_CREATE, 0644) //nolint:gosec
		if fErr != nil {
			fmt.Fprintf(os.Stderr, "Open output file %s fail, err: %v\n", params.outputFile, err)
			os.Exit(1)
		}

		_, fErr = f.Write(schText)
		if fErr != nil {
			fmt.Fprintf(os.Stderr, "Write schema definition to output file %s fail, err: %v\n", params.outputFile, err)
			os.Exit(1)
		}

		_ = f.Close()
		fmt.Printf("Schema definition has generated to output file %s\n", params.outputFile)
	} else {
		fmt.Fprintln(os.Stdout, string(schText))
	}

}
