# mjson2go

## Branch: gh-pages

This branch of the project exists only to host a demo version of the tool on Github Pages. See the main project at 

To build: `GOOS=js GOARCH=wasm go build -o mjson2go.wasm`

# Pipeline Parameters

Input parameters can be specified as a string in your JSON file prefixed with "%%", and optionally with a Go type specification and ordering for the function arguments: 
```json
{
    "index": "%%usingIndex%int%2",
    "date": "%%searchDate%time.Time%1"
}
```

This will produce a function similar to 
```go
func GetAggregation(searchDate time.Time, usingIndex int) bson.D {
    // ... 
}
```

## Field specification
All parts of the parameter specifier are optional, including the parameter name. You can simply write "%%" to create a parameter with all default values. Parameter names can be repeated to reuse the same value in the output.

`"%%<name>%<type>%<order>"`

`<name>` - The parameter name. Defaults to p\<index\> (p0, p1, etc.). Note that the tool will not stop you from using a Go keyword as a name.

`<type>` - The Go type specification. Defaults to string. The tool will run goimports after code generation and attempt to import non-primitive types.

`<order>` - Numeric key for the order of paramters in the function specification. Need not be sequential. Parameters default to source order, with the caveat that all explicitly ordered parameters are added first.

Usage examples: 

- `"%%articleType"` (defaults to string)
- `"%%userId%int"`
- `"%%beforeDate%time.Time"`
- `"%%mongoKey%primitive.ObjectID"`
- `"%%"` (defaults to string with a generated name)

## Command Line Flags

### `-fix`
Corrects three common JSON syntax errors that Compass (and Javascript) may allow, but which will cause unexpected JSON parsing behavior with this tool. 

* Corrects un-quoted object keys, i.e., changes `$match` to `"$match"`
* Change single-quote strings to double quotes
* Remove trailing commas after the last element of an array or object

`-fix` also formats and indents the resulting source with two spaces, and overwrites the original source file.

Defaults to true, so you will have to explicitly pass -fix=false to not fix (potential) source errors.

### `-package=pkgname`
By default, the generated code is in the main package or the value of the `$GOPACAKGE` environment variable (which is set by `go generate`). Setting this flag will override the default name.

### `-v`
(Slightly) more verbose output to `stderr`.

### `-out=filename.go`
Write the output to the provided filename. By default, output goes to `stdout` on the command line, or into 
`<current file>_mjson.go` if called via `go generate`. NOTE: `goimports` will not be run if output goes to 
`stdout`
