// gendocs scans Go source code and generates a Markdown document describing
// all HTTP end-points and the parameters they accept.
//
// Usage:
//
//	go run ./tools/gendocs -root . -out docs/ENDPOINTS.md
//
// What it detects (static analysis only, nothing is executed):
//
//	Routes:
//	  r.GET("/path", handler)            (gin / echo / chi style: GET POST PUT DELETE PATCH HEAD OPTIONS Any)
//	  g := r.Group("/api"); g.POST(...)  (Group prefixes are resolved)
//	  r.Handle("GET", "/path", handler)
//	  http.HandleFunc("/path", handler)  (also "GET /path" Go 1.22 patterns)
//	  Route{Method: "GET", Path: "/x", Handler: h}   (route tables, e.g. []server.Route{...})
//
//	Parameters (inside the handler body):
//	  c.Query / DefaultQuery / GetQuery   -> query
//	  c.Param                             -> path
//	  c.PostForm / DefaultPostForm / FormValue -> form
//	  c.GetHeader / r.Header.Get          -> header
//	  c.ShouldBindJSON(&req) / Bind(&req) / json.NewDecoder(..).Decode(&req)
//	                                      -> request body, fields taken from struct `json`/`form`/`binding` tags
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// ---------- data model ----------

type Param struct {
	Name string
	In   string // query, path, form, header
}

type Field struct {
	Name     string
	Type     string
	Required bool
	Doc      string
}

type Body struct {
	In     string // "json body", "query struct", ...
	Type   string
	Fields []Field
}

type Endpoint struct {
	Service string
	Method  string
	Path    string
	Handler string
	Doc     string
	File    string
	Line    int
	Params  []Param
	Bodies  []Body
}

type funcEntry struct {
	decl *ast.FuncDecl
	dir  string
}

type structEntry struct {
	st *ast.StructType
}

var (
	fset    = token.NewFileSet()
	funcs   = map[string][]funcEntry{}
	structs = map[string]structEntry{}
)

var httpMethods = map[string]bool{
	"GET": true, "POST": true, "PUT": true, "DELETE": true,
	"PATCH": true, "HEAD": true, "OPTIONS": true, "Any": true,
}

// ---------- main ----------

func main() {
	root := flag.String("root", ".", "root directory of the repository")
	out := flag.String("out", "docs/ENDPOINTS.md", "output markdown file")
	title := flag.String("title", "API End-points", "document title")
	flag.Parse()

	files := parseTree(*root)
	collectDecls(files, *root)

	var eps []Endpoint
	for _, pf := range files {
		eps = append(eps, findRoutes(pf, *root)...)
	}
	sort.SliceStable(eps, func(i, j int) bool {
		a, b := eps[i], eps[j]
		if a.Service != b.Service {
			return a.Service < b.Service
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Method < b.Method
	})

	md := render(*title, eps)
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(*out, []byte(md), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("gendocs: wrote %d end-points to %s\n", len(eps), *out)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "gendocs:", err)
	os.Exit(1)
}

// ---------- parsing ----------

type parsedFile struct {
	path string
	file *ast.File
}

func parseTree(root string) []parsedFile {
	var res []parsedFile
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor", "node_modules", "gendocs", "testdata":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, perr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if perr != nil {
			fmt.Fprintf(os.Stderr, "gendocs: skip %s: %v\n", path, perr)
			return nil
		}
		res = append(res, parsedFile{path, f})
		return nil
	})
	if err != nil {
		fatal(err)
	}
	return res
}

func collectDecls(files []parsedFile, root string) {
	for _, pf := range files {
		dir := filepath.Dir(pf.path)
		for _, d := range pf.file.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				funcs[d.Name.Name] = append(funcs[d.Name.Name], funcEntry{d, dir})
			case *ast.GenDecl:
				for _, s := range d.Specs {
					if ts, ok := s.(*ast.TypeSpec); ok {
						if st, ok := ts.Type.(*ast.StructType); ok {
							structs[ts.Name.Name] = structEntry{st}
						}
					}
				}
			}
		}
	}
}

// ---------- route discovery ----------

func findRoutes(pf parsedFile, root string) []Endpoint {
	var eps []Endpoint
	dir := filepath.Dir(pf.path)
	service, _ := filepath.Rel(root, dir)
	if service == "." || service == "" {
		service = "(root)"
	}
	service = filepath.ToSlash(service)

	add := func(method, path string, handler ast.Expr, pos token.Pos) {
		ep := Endpoint{Service: service, Method: strings.ToUpper(method), Path: path}
		p := fset.Position(pos)
		ep.File, ep.Line = filepath.ToSlash(p.Filename), p.Line
		analyzeHandler(&ep, handler, dir)
		eps = append(eps, ep)
	}

	for _, decl := range pf.file.Decls {
		groups := map[string]string{} // variable name -> accumulated prefix
		ast.Inspect(decl, func(n ast.Node) bool {
			switch n := n.(type) {

			case *ast.AssignStmt: // g := r.Group("/api")
				if len(n.Lhs) == 1 && len(n.Rhs) == 1 {
					if id, ok := n.Lhs[0].(*ast.Ident); ok {
						if call, ok := n.Rhs[0].(*ast.CallExpr); ok && isGroupCall(call) {
							groups[id.Name] = prefixOf(call, groups)
						}
					}
				}

			case *ast.CallExpr:
				sel, ok := n.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				name := sel.Sel.Name
				switch {
				case httpMethods[name] && len(n.Args) >= 2:
					if p, ok := strLit(n.Args[0]); ok {
						m := name
						if m == "Any" {
							m = "ANY"
						}
						add(m, prefixOf(sel.X, groups)+p, n.Args[len(n.Args)-1], n.Pos())
					}
				case name == "Handle" && len(n.Args) >= 3:
					m, ok1 := strLit(n.Args[0])
					p, ok2 := strLit(n.Args[1])
					if ok1 && ok2 {
						add(m, prefixOf(sel.X, groups)+p, n.Args[len(n.Args)-1], n.Pos())
					}
				case (name == "HandleFunc" || name == "Handle") && len(n.Args) == 2:
					if p, ok := strLit(n.Args[0]); ok {
						method := "ANY"
						if i := strings.Index(p, " "); i > 0 { // Go 1.22 "GET /path"
							method, p = p[:i], strings.TrimSpace(p[i+1:])
						}
						add(method, prefixOf(sel.X, groups)+p, n.Args[1], n.Pos())
					}
				}

			case *ast.CompositeLit: // Route{Method: "GET", Path: "/x", Handler: h}
				var method, path string
				var handler ast.Expr
				for _, el := range n.Elts {
					kv, ok := el.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key, ok := kv.Key.(*ast.Ident)
					if !ok {
						continue
					}
					switch key.Name {
					case "Method":
						method = methodName(kv.Value)
					case "Path":
						path, _ = strLit(kv.Value)
					case "Handler":
						handler = kv.Value
					}
				}
				if method != "" && path != "" && handler != nil {
					add(method, path, handler, n.Pos())
				}
			}
			return true
		})
	}
	return eps
}

func isGroupCall(c *ast.CallExpr) bool {
	sel, ok := c.Fun.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Group" && len(c.Args) >= 1
}

// prefixOf resolves the accumulated Group prefix of an expression.
func prefixOf(e ast.Expr, groups map[string]string) string {
	switch e := e.(type) {
	case *ast.Ident:
		return groups[e.Name]
	case *ast.CallExpr:
		if isGroupCall(e) {
			sel := e.Fun.(*ast.SelectorExpr)
			p, _ := strLit(e.Args[0])
			return prefixOf(sel.X, groups) + p
		}
	}
	return ""
}

func strLit(e ast.Expr) (string, bool) {
	bl, ok := e.(*ast.BasicLit)
	if !ok || bl.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(bl.Value)
	return s, err == nil
}

func methodName(e ast.Expr) string {
	if s, ok := strLit(e); ok {
		return strings.ToUpper(s)
	}
	if sel, ok := e.(*ast.SelectorExpr); ok { // http.MethodGet
		return strings.ToUpper(strings.TrimPrefix(sel.Sel.Name, "Method"))
	}
	return ""
}

// ---------- handler analysis ----------

func analyzeHandler(ep *Endpoint, h ast.Expr, dir string) {
	var body *ast.BlockStmt
	var ftype *ast.FuncType

	switch h := h.(type) {
	case *ast.FuncLit:
		ep.Handler = "(inline func)"
		body, ftype = h.Body, h.Type
	case *ast.Ident:
		ep.Handler = h.Name
		body, ftype = lookupFunc(h.Name, ep, dir)
	case *ast.SelectorExpr:
		ep.Handler = types.ExprString(h)
		body, ftype = lookupFunc(h.Sel.Name, ep, dir)
	case *ast.CallExpr: // wrapped handler, e.g. mw(handler)
		for i := len(h.Args) - 1; i >= 0; i-- {
			analyzeHandler(ep, h.Args[i], dir)
			if ep.Handler != "" {
				return
			}
		}
		return
	default:
		ep.Handler = types.ExprString(h)
	}
	if body == nil {
		return
	}

	vars := map[string]string{} // local var -> type name
	if ftype != nil && ftype.Params != nil {
		for _, f := range ftype.Params.List {
			for _, n := range f.Names {
				vars[n.Name] = typeName(f.Type)
			}
		}
	}

	seen := map[string]bool{}
	addParam := func(name, in string) {
		k := in + ":" + name
		if !seen[k] {
			seen[k] = true
			ep.Params = append(ep.Params, Param{name, in})
		}
	}
	addBody := func(in string, arg ast.Expr) {
		id := identOf(arg)
		if id == "" {
			return
		}
		tn := vars[id]
		if tn == "" {
			return
		}
		for _, b := range ep.Bodies {
			if b.Type == tn && b.In == in {
				return
			}
		}
		ep.Bodies = append(ep.Bodies, Body{In: in, Type: tn, Fields: structFields(tn, 0)})
	}

	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ValueSpec: // var req T
			if n.Type != nil {
				for _, nm := range n.Names {
					vars[nm.Name] = typeName(n.Type)
				}
			}
		case *ast.AssignStmt: // req := T{} / &T{} / new(T)
			if n.Tok == token.DEFINE && len(n.Lhs) >= 1 && len(n.Rhs) >= 1 {
				if id, ok := n.Lhs[0].(*ast.Ident); ok {
					if tn := typeOfExpr(n.Rhs[0]); tn != "" {
						vars[id.Name] = tn
					}
				}
			}
		case *ast.CallExpr:
			sel, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			arg0, hasLit := "", false
			if len(n.Args) >= 1 {
				arg0, hasLit = strLit(n.Args[0])
			}
			switch sel.Sel.Name {
			case "Query", "DefaultQuery", "GetQuery", "QueryArray", "GetQueryArray":
				if hasLit {
					addParam(arg0, "query")
				}
			case "Param":
				if hasLit {
					addParam(arg0, "path")
				}
			case "PostForm", "DefaultPostForm", "GetPostForm", "PostFormArray":
				if hasLit {
					addParam(arg0, "form")
				}
			case "FormValue":
				if hasLit {
					addParam(arg0, "form/query")
				}
			case "GetHeader":
				if hasLit {
					addParam(arg0, "header")
				}
			case "Get": // r.URL.Query().Get("x") or r.Header.Get("x")
				if hasLit {
					switch x := sel.X.(type) {
					case *ast.CallExpr:
						if s, ok := x.Fun.(*ast.SelectorExpr); ok && s.Sel.Name == "Query" {
							addParam(arg0, "query")
						}
					case *ast.SelectorExpr:
						if x.Sel.Name == "Header" {
							addParam(arg0, "header")
						}
					}
				}
			case "ShouldBindJSON", "BindJSON", "ShouldBind", "Bind":
				if len(n.Args) >= 1 {
					addBody("body", n.Args[0])
				}
			case "ShouldBindQuery", "BindQuery":
				if len(n.Args) >= 1 {
					addBody("query", n.Args[0])
				}
			case "ShouldBindUri", "BindUri":
				if len(n.Args) >= 1 {
					addBody("path", n.Args[0])
				}
			case "Decode": // json.NewDecoder(r.Body).Decode(&req)
				if len(n.Args) == 1 {
					addBody("body", n.Args[0])
				}
			case "Unmarshal": // json.Unmarshal(data, &req)
				if len(n.Args) == 2 {
					addBody("body", n.Args[1])
				}
			}
		}
		return true
	})
}

func lookupFunc(name string, ep *Endpoint, dir string) (*ast.BlockStmt, *ast.FuncType) {
	cands := funcs[name]
	if len(cands) == 0 {
		return nil, nil
	}
	pick := cands[0]
	for _, c := range cands { // prefer same directory (package)
		if c.dir == dir {
			pick = c
			break
		}
	}
	ep.Doc = docText(pick.decl.Doc)
	return pick.decl.Body, pick.decl.Type
}

func identOf(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.UnaryExpr: // &req
		return identOf(e.X)
	}
	return ""
}

func typeName(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return typeName(e.X)
	case *ast.SelectorExpr:
		return e.Sel.Name
	}
	return ""
}

func typeOfExpr(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.CompositeLit:
		return typeName(e.Type)
	case *ast.UnaryExpr:
		return typeOfExpr(e.X)
	case *ast.CallExpr: // new(T)
		if id, ok := e.Fun.(*ast.Ident); ok && id.Name == "new" && len(e.Args) == 1 {
			return typeName(e.Args[0])
		}
	}
	return ""
}

// structFields returns the documented fields of a struct (embedded structs are flattened).
func structFields(name string, depth int) []Field {
	se, ok := structs[name]
	if !ok || depth > 3 {
		return nil
	}
	var out []Field
	for _, f := range se.st.Fields.List {
		if len(f.Names) == 0 { // embedded
			out = append(out, structFields(typeName(f.Type), depth+1)...)
			continue
		}
		var tag reflect.StructTag
		if f.Tag != nil {
			tag = reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
		}
		for _, n := range f.Names {
			if !n.IsExported() {
				continue
			}
			fname := n.Name
			for _, key := range []string{"json", "form", "uri", "query"} {
				if v := strings.Split(tag.Get(key), ",")[0]; v != "" && v != "-" {
					fname = v
					break
				}
			}
			if strings.HasPrefix(tag.Get("json"), "-") && tag.Get("json") == "-" {
				continue
			}
			req := strings.Contains(tag.Get("binding"), "required") ||
				strings.Contains(tag.Get("validate"), "required")
			doc := docText(f.Doc)
			if doc == "" {
				doc = docText(f.Comment)
			}
			out = append(out, Field{fname, types.ExprString(f.Type), req, doc})
		}
	}
	return out
}

func docText(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}
	return strings.TrimSpace(cg.Text())
}

// ---------- markdown rendering ----------

func cell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
}

func anchor(ep Endpoint) string {
	s := strings.ToLower(ep.Method + " " + ep.Path)
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}
	return b.String()
}

func render(title string, eps []Endpoint) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", title)
	b.WriteString("> Auto-generated by `tools/gendocs` – do not edit by hand.\n\n")
	if len(eps) == 0 {
		b.WriteString("_No end-points found._\n")
		return b.String()
	}

	// summary table
	b.WriteString("## Summary\n\n| Service | Method | Path | Handler |\n|---|---|---|---|\n")
	for _, ep := range eps {
		fmt.Fprintf(&b, "| %s | `%s` | `%s` | `%s` |\n", cell(ep.Service), ep.Method, cell(ep.Path), cell(ep.Handler))
	}
	b.WriteString("\n")

	// details per service
	cur := ""
	for _, ep := range eps {
		if ep.Service != cur {
			cur = ep.Service
			fmt.Fprintf(&b, "## Service: `%s`\n\n", cur)
		}
		fmt.Fprintf(&b, "### `%s %s`\n\n", ep.Method, ep.Path)
		if ep.Doc != "" {
			fmt.Fprintf(&b, "%s\n\n", ep.Doc)
		}
		fmt.Fprintf(&b, "- Handler: `%s`\n- Source: `%s:%d`\n\n", ep.Handler, ep.File, ep.Line)

		if len(ep.Params) > 0 {
			b.WriteString("**Parameters**\n\n| Name | In |\n|---|---|\n")
			for _, p := range ep.Params {
				fmt.Fprintf(&b, "| `%s` | %s |\n", cell(p.Name), p.In)
			}
			b.WriteString("\n")
		}
		for _, bd := range ep.Bodies {
			fmt.Fprintf(&b, "**Request %s** – `%s`\n\n", bd.In, bd.Type)
			if len(bd.Fields) == 0 {
				b.WriteString("_(structure not resolved)_\n\n")
				continue
			}
			b.WriteString("| Field | Type | Required | Description |\n|---|---|---|---|\n")
			for _, f := range bd.Fields {
				req := ""
				if f.Required {
					req = "yes"
				}
				fmt.Fprintf(&b, "| `%s` | `%s` | %s | %s |\n", cell(f.Name), cell(f.Type), req, cell(f.Doc))
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}
