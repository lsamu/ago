package jwt

import (
    "errors"
    "github.com/dgrijalva/jwt-go"
    "sync"
    "sync/atomic"
    "time"
)

const claimHistoryResetDuration = time.Hour * 24

type (
    //ParseOption ParseOption
    ParseOption func(parser *TokenParser)
    //TokenParser TokenParser
    TokenParser struct {
        resetTime     time.Duration
        resetDuration time.Duration
        history       sync.Map
    }
)

// NewTokenParser 新的TokenParser
func NewTokenParser(opts ...ParseOption) *TokenParser {
    parser := &TokenParser{
        resetTime:     time.Since(time.Now()),
        resetDuration: claimHistoryResetDuration,
    }
    for _, opt := range opts {
        opt(parser)
    }
    return parser
}

//ParseToken 格式化Token
func (tp *TokenParser) ParseToken(secret, prevSecret string) ( token *jwt.Token, err error) {
    if len(prevSecret) > 0 {

    }
    return
}

//CreateToken 生成Token
func (tp *TokenParser) CreateToken(secret, prevSecret string) ( token *jwt.Token, err error) {
    if len(prevSecret) > 0 {

    }
    return
}

// loadCount 个数
func (tp *TokenParser) loadCount(secret string) uint64 {
    value, ok := tp.history.Load(secret)
    if ok {
        return *value.(*uint64)
    }
    return 0
}

// incrementCount 计数
func (tp *TokenParser) incrementCount(secret string) {
    now := time.Since(time.Now())
    if tp.resetTime+tp.resetDuration < now {
        tp.history.Range(func(key, value interface{}) bool {
            tp.history.Delete(key)
            return true
        })
    }

    value, ok := tp.history.Load(secret)
    if ok {
        atomic.AddUint64(value.(*uint64), 1)
    } else {
        var count uint64 = 1
        tp.history.Store(secret, &count)
    }
}

var ks []byte

//InitJwt InitJwt
func InitJwt(key string) {
    ks = []byte(key)
}

//CustomClaims CustomClaims
type CustomClaims struct {
    ID uint32 `json:"id"`
    IP string `json:"ip"`
    EX string `json:"ex"`
}

//CreateToken CreateToken
func CreateToken(id CustomClaims, key ...interface{}) (token string, err error) {
    claim := jwt.MapClaims{
        "id": id.ID,
        "ip": id.IP,
        "ex": time.Now().String(),
    }
    tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    token, err = tokens.SignedString(ks)
    return token, err
}

//ParseToken ParseToken
func ParseToken(tokens string) (id CustomClaims, err error) {
    token, err := jwt.Parse(tokens, secret())
    if err != nil {
        return
    }
    claim, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        err = errors.New("cannot convert claim to mapclaim")
        return id, err
    }
    //验证token，如果token被修改过则为false
    if !token.Valid {
        err = errors.New("token is invalid")
        return id, err
    }
    id.ID = uint32(claim["id"].(float64))
    id.IP = claim["ip"].(string)
    id.EX = claim["ex"].(string)
    return id, err
}

func secret() jwt.Keyfunc {
    return func(token *jwt.Token) (interface{}, error) {
        return ks, nil
    }
}
