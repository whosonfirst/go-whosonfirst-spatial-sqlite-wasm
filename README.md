# go-whosonfirst-spatial-sqlite-wasm

Experimental Go package to expose the [go-whosonfirst-spatial-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite/blob/main/cmd/query/main.go) `query` function as a WASM binary.

## Important

This is work in progress. It also doesn't work yet.

As of this writing it is failing with the following compile-time errors:

```
> make wasm
GOOS=js GOARCH=wasm go build -mod vendor -o static/wasm/query.wasm cmd/query/main.go
# github.com/psanford/sqlite3vfs
vendor/github.com/psanford/sqlite3vfs/sqlite3vfs.go:28:9: undefined: newVFS
# github.com/whosonfirst/walk
vendor/github.com/whosonfirst/walk/symlink.go:61:6: undefined: IsAbs
vendor/github.com/whosonfirst/walk/walk.go:287:12: undefined: volumeNameLen
vendor/github.com/whosonfirst/walk/walk.go:471:15: undefined: volumeNameLen
make: *** [wasm] Error 2
```

Any help or suggestions would be appreciated.

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite