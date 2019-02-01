## mimetypes

[![godoc](https://godoc.org/github.com/moriyoshi/mimetypes?status.svg)](https://godoc.org/github.com/moriyoshi/mimetypes)

[![travis-ci](https://travis-ci.com/moriyoshi/mimetypes.svg?branch=master)](https://travis-ci.com/moriyoshi/mimetypes)

A simple library for parsing Apache-style mime.types / XDG-defined media type globs.

This library was created to overcome the limitation of the implementation of `mime` package, which takes it for granted that an apache-style mime type file is installed in any of the hard-coded locations as for Unix-like platform.

### Usage

Be sure to import `github.com/moriyoshi/mimetypes/loaders` if you want to use the default loader implementations, like the following:

```go
import (
    _ "github.com/moriyoshi/mimetypes/loaders"
)
```

This ensures the default loaders to be compiled and incorporated into the build.


### License

MIT


