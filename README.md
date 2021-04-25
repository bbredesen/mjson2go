# mjson2go

`mjson2go` is a tool that converts JSON to parameterized Go code usable by the MongoDB driver.

Writing a or converting a MongoDB pipeline in Go is awkward, with lots of nested braces. 
Here is a very simple aggregation pipeline in JSON:

```json
[
    {
        "$match" : {
            "postType" : "Article"
        }
    },
    {
        "$group" : {
            "_id" : {"user": "$userId"},
            "allPosts": { "$push" : "$$ROOT" },
            "count" : { "$sum" : 1 }
        }
    }
]
```

Translated to Go (and `go fmt`-ed), this becomes:

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

`mjson2go` allows you to export an aggregation pipeline from MongoDB Compass and/or directly develop your pipeline in JSON, instead of trying to copy, paste into Go source, correct syntax, match braces, etc. etc.

`mjson2go` can be used from the command line or, ideally, via `go generate`. The tool produces a Go source file with a "Get___" function that returns a `bson.D` or `bson.A`, and allows you to specify input parameters from your code.

Finally, the tool will format and run `goimports` on the resulting source file.

# Usage

```
go get github.com/bbredesen/mjson2go
```

## Via the Command Line
```
mjson2go -out=pipeline.go aggregation.json otherAggregation.json
```

The command above will produce Go source containing functions `GetAggregation()` and `GetOtherAggregation()`. If no files are provided, it will read from stdin. If a directory name is provided, it will process all files ending in .json in that directory.

## Via `go generate`
The tool will automatically put the resulting code in the same package as the `go:generate` annotation:

```go
package somepkg
//go:generate mjson2go -out=pipeline.go aggregations/
```

Note that `go generate` does not pass commands through a shell and thus will not expand globs.

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
All parts of the parameter specifier are optional, though you must at least provide a name after "%%" to create a parameter with default values.

`"%%<name>%<type>%<order>"`

`<name>` - The parameter name. Defaults to `p<index>` (p0, p1, etc.). Note that the tool will not stop you from using a Go keyword as a name.

`<type>` - The Go type specification. Defaults to `string`. Non-primitive types will be imported to the file via goimports

`<order>` - Numeric key for the order of paramters in the function specification. Need not be sequential.

For example: `%%articleType` (defaults to string), `%%userId%int`, `%%beforeDate%time.Time`, or `%%mongoKey%primitive.ObjectID`

## Command Line Flags

### `-fix`
Corrects three common JSON syntax errors that Compass (and Javascript) may allow, but which will cause unexpected JSON parsing behavior with this tool. 

* Corrects un-quoted object keys, i.e., changes `$match` to `"$match"`
* Change single-quote strings to double quotes
* Remove trailing commas after the last element of an array or object

`-fix` also formats and indents the resulting source with two spaces, and overwrites the original source file. `-fix` defaults to true, so you will have to explicitly say -fix=false to not fix potentila source errors.

### `-package=pkgname`
By default, the generated code is in the main package or the value of the `$GOPACAKGE` environment variable (which is set by `go generate`). Setting this flag will override the default name.

### `-v`
(Slightly) more verbose output to `stderr`.

### `-out=filename.go`
Write the output to the provided filename. By default, output goes to `stdout` on the command line, or into 
`<current file>_mjson.go` if called via `go generate`. WARNING: `goimports` will not be run if output goes to 
`stdout`
