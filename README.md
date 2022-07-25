# letscrum

## Generate API

```text
protoc -I . -I third_party --plugin="protoc-gen-ts=.\node_modules\.bin\protoc-gen-ts.cmd" --js_out="import_style=commonjs,binary:sdk" --ts_out="sdk" .\api\general\v1\common.proto .\api\general\v1\letscrum.proto .\api\letscrum\v1\letscrum.proto .\api\project\v1\project.proto
```
