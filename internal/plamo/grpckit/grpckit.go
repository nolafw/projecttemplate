package grpckit

// TODO:

// | 用途           | ライブラリ                                                              |
// | ------------ | ------------------------------------------------------------------ |
// | 認証・認可        | [grpc\_auth](https://github.com/grpc-ecosystem/go-grpc-middleware) |
// | ログ出力         | `zap` + `grpc_zap`                                                 |
// | リカバリ（panic）  | `grpc_recovery`                                                    |

// ✅ 自作Interceptorの活用パターン
// 以下のような処理が共通化可能です：
// 認証（JWTトークンの検証）
// メタデータ（ヘッダー）の検査・追加
// ログ出力
// リカバリ(panic)
// トレース（OpenTelemetry）
