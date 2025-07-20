package errno

import (
	"github.com/yanking/miniblog/pkg/errorsx"
	"net/http"
)

// ErrPostNotFound 表示未找到指定的博客.
var ErrPostNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PostNotFound", Message: "Post not found."}
