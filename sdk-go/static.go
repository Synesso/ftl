package sdkgo

import (
	stderrors "errors" //nolint:depguard
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
	"sync"

	"github.com/alecthomas/errors"
	"github.com/alecthomas/participle/v2"
	"github.com/iancoleman/strcase"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"

	"github.com/TBD54566975/ftl/schema"
	"github.com/TBD54566975/ftl/sdk-go/internal"
)

var fset = token.NewFileSet()

var contextIfaceType = once(func() *types.Interface {
	return mustLoadRef("context", "Context").Type().Underlying().(*types.Interface) //nolint:forcetypeassert
})
var errorIFaceType = once(func() *types.Interface {
	return mustLoadRef("builtin", "error").Type().Underlying().(*types.Interface) //nolint:forcetypeassert
})
var ftlCallFuncPath = "github.com/TBD54566975/ftl/sdk-go.Call"

// ExtractModule statically parses Go FTL module source into a schema.Module.
func ExtractModule(dir string) (*schema.Module, error) {
	pkgs, err := packages.Load(&packages.Config{
		Dir:  dir,
		Fset: fset,
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}, "./...")
	if err != nil {
		return &schema.Module{}, errors.WithStack(err)
	}
	if len(pkgs) == 0 {
		return &schema.Module{}, errors.Errorf("no packages found in %q, does \"go mod tidy\" need to be run?", dir)
	}
	module := &schema.Module{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			var verb *schema.Verb
			err := internal.Visit(file, func(node ast.Node, next func() error) (err error) {
				defer func() {
					if err != nil {
						err = errors.Wrap(err, fset.Position(node.Pos()).String())
					}
				}()
				switch node := node.(type) {
				case *ast.CallExpr:
					if err := visitCallExpr(module, verb, node, pkg); err != nil {
						return err
					}

				case *ast.File:
					if err := visitFile(module, node, pkg); err != nil {
						return err
					}

				case *ast.FuncDecl:
					verb, err = visitFuncDecl(pkg, module, node)
					if err != nil {
						return err
					}
					err = next()
					if err != nil {
						return err
					}
					verb = nil
					return nil

				case nil:
				default:
				}
				return next()
			})
			if err != nil {
				return nil, err
			}
		}
	}
	if module.Name == "" {
		return module, errors.Errorf("//ftl:module directive is required")
	}
	return module, schema.ValidateModule(module)
}

func visitCallExpr(module *schema.Module, verb *schema.Verb, node *ast.CallExpr, pkg *packages.Package) error {
	_, fn := deref[*types.Func](pkg, node.Fun)
	if fn == nil {
		return nil
	}
	if fn.FullName() != ftlCallFuncPath {
		return nil
	}
	if len(node.Args) != 3 {
		return errors.New("Call must have exactly three arguments")
	}
	_, verbFn := deref[*types.Func](pkg, node.Args[1])
	if verbFn == nil {
		return errors.Errorf("Call first argument must be a function but is %s", node.Args[1])
	}
	moduleName := verbFn.Pkg().Name()
	if moduleName == pkg.Name {
		moduleName = ""
	}
	ref := &schema.VerbRef{Module: moduleName, Name: strcase.ToLowerCamel(verbFn.Name())}
	verb.AddCall(ref)
	return nil
}

func visitFile(module *schema.Module, node *ast.File, pkg *packages.Package) error {
	if node.Doc == nil {
		return nil
	}
	directives, err := parseDirectives(fset, node.Doc)
	if err != nil {
		return errors.WithStack(err)
	}
	module.Comments = parseComments(node.Doc)
	for _, dir := range directives {
		switch dir.kind {
		case "module":
			if dir.id == "" {
				return errors.Errorf("%s: module not specified", dir)
			}
			if dir.id != pkg.Name {
				return errors.Errorf("%s: FTL module name %q does not match Go package name %q", dir, dir.id, pkg.Name)
			}
			module.Name = dir.id

		default:
			return errors.Errorf("%s: invalid directive", dir)
		}
	}
	return nil
}

func isType[T types.Type](t types.Type) bool {
	if _, ok := t.(*types.Named); ok {
		t = t.Underlying()
	}
	_, ok := t.(T)
	return ok
}

func checkSignature(sig *types.Signature) error {
	params := sig.Params()
	results := sig.Results()
	if params.Len() != 2 {
		return errors.Errorf("must have exactly two parameters in the form (context.Context, struct) but has %d", params.Len())
	}
	if results.Len() != 2 {
		return errors.Errorf("must have exactly two result values in the form (error, struct) but has %d", results.Len())
	}
	if !types.AssertableTo(contextIfaceType(), params.At(0).Type()) {
		return errors.Errorf("first parameter must be of type context.Context but is %s", params.At(0).Type())
	}
	if !isType[*types.Struct](params.At(1).Type()) {
		return errors.Errorf("second parameter must be a struct but is %s", params.At(1).Type())
	}
	if !types.AssertableTo(errorIFaceType(), results.At(1).Type()) {
		return errors.Errorf("first result must be an error but is %s", results.At(0).Type())
	}
	if !isType[*types.Struct](results.At(0).Type()) {
		return errors.Errorf("first result must be a struct but is %s", results.At(0).Type())
	}
	return nil
}

// "verbIndex" is the index into the Module.Decls of the verb that was parsed.
func visitFuncDecl(pkg *packages.Package, module *schema.Module, node *ast.FuncDecl) (verb *schema.Verb, err error) { //nolint:unparam
	if node.Doc == nil {
		return nil, nil
	}
	fnt := pkg.TypesInfo.Defs[node.Name].(*types.Func) //nolint:forcetypeassert
	sig := fnt.Type().(*types.Signature)               //nolint:forcetypeassert
	if sig.Recv() != nil {
		return nil, errors.Errorf("ftl:verb cannot be a method")
	}
	params := sig.Params()
	results := sig.Results()
	if err := checkSignature(sig); err != nil {
		return nil, err
	}
	req, err := parseStruct(pkg, module, params.At(1).Type())
	if err != nil {
		return nil, err
	}
	resp, err := parseStruct(pkg, module, results.At(0).Type())
	if err != nil {
		return nil, err
	}
	verb = &schema.Verb{
		Comments: parseComments(node.Doc),
		Name:     strcase.ToLowerCamel(node.Name.Name),
		Request:  req,
		Response: resp,
	}
	module.Decls = append(module.Decls, verb)
	return verb, nil
}

func parseComments(doc *ast.CommentGroup) []string {
	comments := []string{}
	if doc := doc.Text(); doc != "" {
		comments = strings.Split(strings.TrimSpace(doc), "\n")
	}
	return comments
}

func parseStruct(pkg *packages.Package, module *schema.Module, node types.Type) (*schema.DataRef, error) {
	named, ok := node.(*types.Named)
	if !ok {
		return nil, errors.Errorf("expected named type but got %s", node)
	}
	s, ok := node.Underlying().(*types.Struct)
	if !ok {
		return nil, errors.Errorf("expected struct but got %s", node)
	}
	out := &schema.Data{
		Name: named.Obj().Name(),
	}
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		ft, err := parseType(pkg, module, f.Type())
		if err != nil {
			return nil, errors.WithStack(err)
		}
		out.Fields = append(out.Fields, &schema.Field{Name: strcase.ToLowerCamel(f.Name()), Type: ft})
	}
	module.AddData(out)
	return &schema.DataRef{Name: out.Name}, nil
}

func parseType(pkg *packages.Package, module *schema.Module, node types.Type) (schema.Type, error) {
	switch node := node.Underlying().(type) {
	case *types.Basic:
		switch node.Kind() {
		case types.String:
			return &schema.String{}, nil

		case types.Int:
			return &schema.Int{}, nil

		case types.Bool:
			return &schema.Bool{}, nil

		case types.Float64:
			return &schema.Float{}, nil

		default:
			return nil, errors.Errorf("unsupported basic type %s", node)
		}

	case *types.Struct:
		ref, err := parseStruct(pkg, module, node)
		return ref, err

	case *types.Map:
		return parseMap(pkg, module, node)

	case *types.Slice:
		return parseSlice(pkg, module, node)

	default:
		return nil, errors.Errorf("unsupported type %s", node)
	}
}

func parseMap(pkg *packages.Package, module *schema.Module, node *types.Map) (*schema.Map, error) {
	key, err := parseType(pkg, module, node.Key())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	value, err := parseType(pkg, module, node.Elem())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &schema.Map{Key: key, Value: value}, nil
}

func parseSlice(pkg *packages.Package, module *schema.Module, node *types.Slice) (*schema.Array, error) {
	value, err := parseType(pkg, module, node.Elem())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &schema.Array{Element: value}, nil
}

type ftlDirective struct {
	kind  string
	id    string
	attrs map[string]directiveValue
}

func (f ftlDirective) String() string {
	out := &strings.Builder{}
	fmt.Fprintf(out, "//ftl:%s", f.kind)
	if f.id != "" {
		fmt.Fprintf(out, " %s", f.id)
	}
	keys := maps.Keys(f.attrs)
	slices.Sort(keys)
	for _, key := range keys {
		fmt.Fprintf(out, " %s=%s", key, f.attrs[key])
	}
	return out.String()
}

// A little parser for Go FTL comment-directives.
type directive struct {
	Kind  string          `parser:"'ftl' ':' @Ident"`
	ID    string          `parser:"( @(Ident | String)"`
	Attrs []directiveAttr `parser:"  @@* )?"`
}

type directiveAttr struct {
	Key   string         `parser:"@Ident '='"`
	Value directiveValue `parser:"@@"`
}

type directiveValue struct {
	Ident *string  `parser:"  @Ident"`
	Str   *string  `parser:"| @String"`
	Int   *int64   `parser:"| @Int"`
	Float *float64 `parser:"| @Float"`
	Bool  *dirBool `parser:"| @('true'|'false')"`
}

type dirBool bool

func (b *dirBool) UnmarshalText(d []byte) error {
	*b = dirBool(string(d) == "true")
	return nil
}

func (d directiveValue) String() string {
	switch {
	case d.Ident != nil:
		return *d.Ident
	case d.Str != nil:
		return strconv.Quote(*d.Str)
	case d.Int != nil:
		return strconv.FormatInt(*d.Int, 10)
	case d.Float != nil:
		return strconv.FormatFloat(*d.Float, 'g', 2, 64)
	case d.Bool != nil:
		return strconv.FormatBool(bool(*d.Bool))
	default:
		panic("??")
	}
}

var directiveParser = participle.MustBuild[directive](participle.Unquote())

func parseDirectives(fset *token.FileSet, docs *ast.CommentGroup) ([]ftlDirective, error) {
	if docs == nil {
		return nil, nil
	}
	directives := []ftlDirective{}
	for _, line := range docs.List {
		if !strings.HasPrefix(line.Text, "//ftl:") {
			continue
		}
		ast, err := directiveParser.ParseString("", line.Text[2:])
		if err != nil {
			// Adjust the Participle-reported position relative to the AST node.
			pos := fset.Position(line.Pos())
			var perr participle.Error
			if stderrors.As(err, &perr) {
				ppos := perr.Position()
				ppos.Filename = pos.Filename
				ppos.Column += pos.Column
				ppos.Line += pos.Line - 1
				err = participle.Errorf(ppos, "%s", perr.Message())
			} else {
				err = errors.Errorf("%s: %s", pos, err)
			}
			return nil, errors.Wrap(err, "invalid directive")
		}
		attrs := map[string]directiveValue{}
		for _, attr := range ast.Attrs {
			attrs[attr.Key] = attr.Value
		}
		directives = append(directives, ftlDirective{kind: ast.Kind, id: ast.ID, attrs: attrs})
	}
	return directives, nil
}

func once[T any](f func() T) func() T {
	var once sync.Once
	var t T
	return func() T {
		once.Do(func() { t = f() })
		return t
	}
}

// Lazy load the compile-time reference from a package.
func mustLoadRef(pkg, name string) types.Object {
	pkgs, err := packages.Load(&packages.Config{Fset: fset, Mode: packages.NeedTypes}, pkg)
	if err != nil {
		panic(err)
	}
	if len(pkgs) != 1 {
		panic("expected one package")
	}
	obj := pkgs[0].Types.Scope().Lookup(name)
	if obj == nil {
		panic("interface not found")
	}
	return obj
}

func deref[T types.Object](pkg *packages.Package, node ast.Expr) (string, T) {
	var obj T
	switch node := node.(type) {
	case *ast.Ident:
		obj, _ = pkg.TypesInfo.Uses[node].(T)
		return "", obj

	case *ast.SelectorExpr:
		x, ok := node.X.(*ast.Ident)
		if !ok {
			return "", obj
		}
		obj, _ = pkg.TypesInfo.Uses[node.Sel].(T)
		return x.Name, obj

	default:
		return "", obj
	}
}
