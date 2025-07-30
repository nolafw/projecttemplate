package grpckit

import (
	"context"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// duration を人間が読みやすい形式でフォーマットする関数
func formatDuration(d time.Duration) string {
	// time.Duration.String() を使用すると以下のような形式になります：
	// - 1.5s (1.5秒)
	// - 250ms (250ミリ秒)
	// - 10µs (10マイクロ秒)
	// - 5ns (5ナノ秒)
	return d.String()
}

// より詳細なduration表示（小数点以下の精度を制御）
func formatDurationDetailed(d time.Duration) string {
	switch {
	case d >= time.Second:
		return d.Truncate(time.Millisecond).String()
	case d >= time.Millisecond:
		return d.Truncate(time.Microsecond).String()
	default:
		return d.String()
	}
}

// 拡張されたクライアント情報を取得するヘルパー関数
func getExtendedClientInfo(ctx context.Context) map[string]any {
	info := make(map[string]any)

	// peer情報からIPアドレスとポート情報を取得
	if p, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
			info["client_ip"] = tcpAddr.IP.String()
			info["client_port"] = tcpAddr.Port
			info["network"] = tcpAddr.Network()
		} else {
			info["client_addr"] = p.Addr.String()
		}

		// TLS情報（存在する場合）
		if authInfo := p.AuthInfo; authInfo != nil {
			info["auth_type"] = authInfo.AuthType()
		}
	}

	// gRPCメタデータから有用な情報を取得
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// User-Agent情報
		if userAgent := md.Get("user-agent"); len(userAgent) > 0 {
			info["user_agent"] = userAgent[0]
		} else {
			info["user_agent"] = "grpc-client"
		}

		// gRPCバージョン
		if grpcVersion := md.Get("grpc-version"); len(grpcVersion) > 0 {
			info["grpc_version"] = grpcVersion[0]
		}

		// Content-Type
		if contentType := md.Get("content-type"); len(contentType) > 0 {
			info["content_type"] = contentType[0]
		}

		// Accept-Encoding
		if acceptEncoding := md.Get("grpc-accept-encoding"); len(acceptEncoding) > 0 {
			info["accept_encoding"] = acceptEncoding[0]
		}

		// カスタムヘッダー（X-で始まるもの）
		customHeaders := make(map[string]string)
		for key, values := range md {
			if len(key) > 2 && (key[:2] == "x-" || key[:2] == "X-") && len(values) > 0 {
				customHeaders[key] = values[0]
			}
		}
		if len(customHeaders) > 0 {
			info["custom_headers"] = customHeaders
		}

		// リクエストID（存在する場合）
		if requestID := md.Get("x-request-id"); len(requestID) > 0 {
			info["request_id"] = requestID[0]
		}

		// トレースID（存在する場合）
		if traceID := md.Get("x-trace-id"); len(traceID) > 0 {
			info["trace_id"] = traceID[0]
		}
	}

	return info
}

// クライアント情報を取得するヘルパー関数（簡易版）
func getClientInfo(ctx context.Context) (clientIP, userAgent string) {
	// peer情報からIPアドレスを取得
	if p, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
			clientIP = tcpAddr.IP.String()
		} else {
			clientIP = p.Addr.String()
		}
	}

	// メタデータからUser-Agentを取得（存在する場合）
	// gRPCではHTTPのようなUser-Agentヘッダーは標準ではないが、
	// カスタムメタデータとして送信されることがある
	userAgent = "grpc-client"

	return clientIP, userAgent
}

// より詳細なクライアント情報を取得する関数
func getDetailedClientInfo(ctx context.Context) map[string]any {
	info := make(map[string]any)

	// peer情報から詳細な接続情報を取得
	if p, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
			info["client_ip"] = tcpAddr.IP.String()
			info["client_port"] = tcpAddr.Port
			info["network"] = tcpAddr.Network()
		} else {
			info["client_addr"] = p.Addr.String()
		}

		// TLS情報（存在する場合）
		if authInfo := p.AuthInfo; authInfo != nil {
			info["auth_type"] = authInfo.AuthType()
		}
	}

	return info
}

// TODO: 何を出力するかは、クライアントコード側で決められるようにしたい。
// callback関数を引数に取って、そこで出力する内容はクライアントコードで決めるようにするか?
// gRPCリクエストログ出力用のUnaryServerInterceptor
func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startTime := time.Now()
		clientInfo := getExtendedClientInfo(ctx)

		// 基本ログ項目
		logAttrs := []any{
			"method", info.FullMethod,
			"timestamp", startTime,
		}

		// クライアント情報を追加
		for key, value := range clientInfo {
			logAttrs = append(logAttrs, key, value)
		}

		// リクエスト開始ログ
		logger.Info("gRPC request started", logAttrs...)

		// 実際のハンドラーを実行
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		// レスポンスログの基本項目
		responseLogAttrs := []any{
			"method", info.FullMethod,
			"duration", duration.String(), // 人間が読みやすい形式（例: 1.5s, 250ms, 10µs）
		}

		// クライアント情報を再度追加
		for key, value := range clientInfo {
			responseLogAttrs = append(responseLogAttrs, key, value)
		}

		if err != nil {
			// エラーレスポンスのログ
			st, _ := status.FromError(err)
			responseLogAttrs = append(responseLogAttrs,
				"code", st.Code(),
				"error", err.Error(),
			)
			logger.Error("gRPC request failed", responseLogAttrs...)
		} else {
			// 成功レスポンスのログ
			responseLogAttrs = append(responseLogAttrs, "code", codes.OK)
			logger.Info("gRPC request completed", responseLogAttrs...)
		}

		return resp, err
	}
}

// TODO: 何を出力するかは、クライアントコード側で決められるようにしたい。
// callback関数を引数に取って、そこで出力する内容はクライアントコードで決めるようにするか?
// gRPCストリーミングログ出力用のStreamServerInterceptor
func LoggingStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now()
		clientInfo := getExtendedClientInfo(stream.Context())

		// 基本ログ項目
		logAttrs := []any{
			"method", info.FullMethod,
			"timestamp", startTime,
			"client_stream", info.IsClientStream,
			"server_stream", info.IsServerStream,
		}

		// クライアント情報を追加
		for key, value := range clientInfo {
			logAttrs = append(logAttrs, key, value)
		}

		// ストリーミング開始ログ
		logger.Info("gRPC stream started", logAttrs...)

		// 実際のハンドラーを実行
		err := handler(srv, stream)

		duration := time.Since(startTime)

		// レスポンスログの基本項目
		responseLogAttrs := []any{
			"method", info.FullMethod,
			"duration", duration.String(), // 人間が読みやすい形式（例: 1.5s, 250ms, 10µs）
		}

		// クライアント情報を再度追加
		for key, value := range clientInfo {
			responseLogAttrs = append(responseLogAttrs, key, value)
		}

		if err != nil {
			// エラーレスポンスのログ
			st, _ := status.FromError(err)
			responseLogAttrs = append(responseLogAttrs,
				"code", st.Code(),
				"error", err.Error(),
			)
			logger.Error("gRPC stream failed", responseLogAttrs...)
		} else {
			// 成功レスポンスのログ
			responseLogAttrs = append(responseLogAttrs, "code", codes.OK)
			logger.Info("gRPC stream completed", responseLogAttrs...)
		}

		return err
	}
}

// gRPCサーバー作成用のヘルパー関数
func NewGRPCServerWithLogging(logger *slog.Logger, opts ...grpc.ServerOption) *grpc.Server {
	// ログ出力interceptorを追加
	defaultOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(LoggingUnaryInterceptor(logger)),
		grpc.StreamInterceptor(LoggingStreamInterceptor(logger)),
	}

	// 追加のオプションをマージ
	allOpts := append(defaultOpts, opts...)

	return grpc.NewServer(allOpts...)
}

// | 用途           | ライブラリ                                                              |
// | ------------ | ------------------------------------------------------------------ |
// | 認証・認可        | [grpc\_auth](https://github.com/grpc-ecosystem/go-grpc-middleware) |
// | リカバリ（panic）  | `grpc_recovery`                                                    |

// ✅ 自作Interceptorの活用パターン
// 以下のような処理が共通化可能です：
// 認証（JWTトークンの検証）
// メタデータ（ヘッダー）の検査・追加
// リカバリ(panic)
