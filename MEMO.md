プロジェクトのディレクトリの構造を考察中


grpcurlテストコマンド

```sh
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 list User
grpcurl -plaintext -d '{"userId":"1"}' localhost:50051 User.GetUser
```