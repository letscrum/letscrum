# letscrum

## Generate API

```Command for Windows
protoc -I . -I third_party --plugin="protoc-gen-ts=.\node_modules\.bin\protoc-gen-ts.cmd" --js_out="import_style=commonjs,binary:sdk" --ts_out="sdk" .\api\general\v1\common.proto .\api\general\v1\letscrum.proto .\api\letscrum\v1\letscrum.proto
```

```Command for Linux/OSX
protoc -I . -I third_party --plugin="protoc-gen-ts=.\node_modules\.bin\protoc-gen-ts.cmd" --js_out="import_style=commonjs,binary:sdk" --ts_out="sdk" .\api\general\v1\common.proto .\api\general\v1\letscrum.proto .\api\letscrum\v1\letscrum.proto
```
