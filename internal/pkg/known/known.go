package known

// 定义 HTTP/gRPC Header.
// gRPC 底层使用了 HTTP/2 作为传输协议，而 HTTP/2 的规范
// 规定 Header 的键必须是小写的。因此，在 gRPC 中，所有的 Header 键都会被强制转换为小写，
// 以符合 HTTP/2 的要求。HTTP/1.x 中，许多 HTTP/1.x 的实现会保留用户设置的大小写格式
// 但一些 HTTP 框架或工具库（比如某些 web 服务器或代理）可能会自动将 Header 转为小写，
// 以简化处理逻辑.
// 考虑兼容性这里统一将 Header 设置为小写.
// 另外，Header 的键以 x- 开头，说明是自定义 Header.
const (
	// XRequestID 用来定义上下文中的键，代表请求 ID.
	XRequestID = "x-request-id"

	// XUserID 用来定义上下文的键，代表请求用户 ID. UserID 整个用户生命周期唯一.
	XUserID = "x-user-id"
)
