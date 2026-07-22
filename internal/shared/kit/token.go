package kit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
)

// TODO: 生产从配置读取密钥, 并替换为 JWT / gtoken。
const tokenSecret = "my_service_token_secret_change_me"

// IssueToken 生成简单签名 token: "<userId>.<sign>"。
func IssueToken(userId int64) string {
	sign := gmd5.MustEncryptString(fmt.Sprintf("%d:%s", userId, tokenSecret))
	return fmt.Sprintf("%d.%s", userId, sign)
}

// ParseToken 校验并解出 userId。
func ParseToken(token string) (int64, error) {
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return 0, gerror.New("invalid token")
	}
	uid, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, gerror.New("invalid token")
	}
	if gmd5.MustEncryptString(fmt.Sprintf("%d:%s", uid, tokenSecret)) != parts[1] {
		return 0, gerror.New("invalid token")
	}
	return uid, nil
}
