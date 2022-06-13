# mjson2go

`mjson2go` is a tool that generates parameterized Go code usable by the MongoDB driver from a JSON pipeline source.

### [Try it in your browser](https://bbredesen.github.io/mjson2go)

Developing a MongoDB pipeline in Go is awkward, with lots of double braces and nested `bson.D` or `bson.A`s. Here is a very simple aggregation pipeline in JSON:

```json
[
    {
        "$match" : {
            "postType" : "Article"
        }
    },
    {
        "$group" : {
            "_id" : { "user" : "$userId" },
            "allPosts": { "$push" : "$$ROOT" },
            "count" : { "$sum" : 1 }
        }
    }
]
```

Translated to Go, this becomes:

```go
pipeline := bson.A{
    bson.D{
        { "$match", bson.D{
            { "postType", "Article" }
        },
    }},
    bson.D{
        { "$group", bson.D{
                { "_id", bson.D{
                    {"user", "$userId"}
                }},
                {"allPosts", bson.D{
                    {"$push", "$$ROOT"}
                }},
                {"count", bson.D{
                    {"$sum", 1}
                }},
        }
    }},
}
```

`mjson2go` allows you to export an aggregation pipeline from MongoDB Compass and/or directly develop your pipeline in JSON, instead of trying to copy and paste, translate into Go source, correct syntax, match braces, etc. etc.

`mjson2go` can be used from the command line or, ideally, via `go generate`. The tool produces a Go source file with a "Get___" function that returns a `bson.D` or `bson.A`, and allows you to specify input parameters from your code.

Finally, the tool will format and run `goimports` on the resulting source file.

# Usage

```
go install github.com/bbredesen/mjson2go@latest
```

## Via the Command Line
```
mjson2go -out=pipeline.go aggregation.json otherAggregation.json
```

The command above will produce Go source containing functions `GetAggregation()` and `GetOtherAggregation()`. If no files are provided, it will read from stdin and write to stdout. If a directory name is provided, it will process all files ending in .json in that directory.

## Via `go generate`
The tool will automatically put the resulting code in the same package as the `go:generate` annotation:

```go
package somepkg
//go:generate mjson2go -out=pipeline.go aggregations/
```

Note that `go generate` does not pass commands through a shell and will not expand globs.

# Pipeline Parameters

Input parameters can be specified as a string in your JSON file prefixed with "%%", and optionally with a Go type specification and ordering for the function arguments: 
```json
{
    "index": "%%intParam%int%2",
    "date": "%%dateParam%time.Time%1"
}
```

This will produce a function similar to 
```go
func GetAggregation(dateParam time.Time, intParam int) bson.D {
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
