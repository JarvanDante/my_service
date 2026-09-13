package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/shared/aesapi"
)

const headerEncrypted = "X-Encrypted"
const headerDebugKey = "debugKey"

// ApiCrypto 对齐公司 H5 接口加解密: 请求体 / 响应信封 AES-128-ECB+Base64。
// 须挂在 MiddlewareHandlerResponse 外侧, 才能改到已经包好的 {code,message,data}。
// Authorization 保持明文, 上传/分片/对象流不加密。
func ApiCrypto(r *ghttp.Request) {
	if !apiCryptoEnabled(r) || !isFrontAPI(r.URL.Path) || isCryptoSkipPath(r.URL.Path) || r.Method == http.MethodOptions {
		r.Middleware.Next()
		return
	}
	if isApiDebug(r) {
		r.Middleware.Next()
		return
	}

	encryptResp := !apiAllowPlain(r) || r.Header.Get(headerEncrypted) == "1"
	if err := decryptFrontBody(r, &encryptResp); err != nil {
		r.SetError(err)
		return
	}

	r.Middleware.Next()

	if encryptResp {
		encryptFrontResponse(r)
	}
}

func decryptFrontBody(r *ghttp.Request, encryptResp *bool) error {
	if !hasRewritableBody(r) {
		return nil
	}
	body := r.GetBody()
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	plain, err := aesapi.DecryptBase64(string(body))
	if err != nil {
		if apiAllowPlain(r) && aesapi.LooksLikeJSON(body) {
			*encryptResp = *encryptResp && r.Header.Get(headerEncrypted) == "1"
			return nil
		}
		return gerror.NewCode(gcode.CodeInvalidParameter, "数据封装错误")
	}
	replaceJSONBody(r, plain)
	*encryptResp = true
	return nil
}

func encryptFrontResponse(r *ghttp.Request) {
	buf := bytes.TrimSpace(r.Response.Buffer())
	if !aesapi.LooksLikeJSON(buf) {
		return
	}
	enc, err := aesapi.EncryptBase64(buf)
	if err != nil {
		g.Log().Warningf(r.GetCtx(), "api_aes 响应加密失败: %v", err)
		return
	}
	r.Response.ClearBuffer()
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Del("Content-Length")
	r.Response.Write(enc)
}

func replaceJSONBody(r *ghttp.Request, plain []byte) {
	r.Body = io.NopCloser(bytes.NewReader(plain))
	r.ReloadParam()
	r.Body = io.NopCloser(bytes.NewReader(plain))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Length", strconv.Itoa(len(plain)))
	r.ContentLength = int64(len(plain))
}

func hasRewritableBody(r *ghttp.Request) bool {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
	default:
		return false
	}
	ct := strings.ToLower(r.Header.Get("Content-Type"))
	return !strings.Contains(ct, "multipart/")
}

func isFrontAPI(path string) bool {
	return strings.HasPrefix(path, "/front/v1")
}

func isCryptoSkipPath(path string) bool {
	switch path {
	case "/front/v1/media/upload", "/front/v1/media/object", "/front/v1/media/multipart/part":
		return true
	default:
		return false
	}
}

func isApiDebug(r *ghttp.Request) bool {
	got := strings.TrimSpace(r.Header.Get(headerDebugKey))
	if got == "" {
		return false
	}
	want := strings.TrimSpace(g.Cfg().MustGet(r.GetCtx(), "api_aes.debug_header", "").String())
	if want == "" {
		return false
	}
	return got == want
}

func apiCryptoEnabled(r *ghttp.Request) bool {
	return g.Cfg().MustGet(r.GetCtx(), "api_aes.enabled", true).Bool()
}

func apiAllowPlain(r *ghttp.Request) bool {
	return g.Cfg().MustGet(r.GetCtx(), "api_aes.allow_plain", false).Bool()
}
