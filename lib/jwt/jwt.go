package jwt

import (
    "errors"
    "github.com/dgrijalva/jwt-go"
    "time"
)

var ks []byte

func InitJwt(key string) {
    ks = []byte(key)
}

type CustomClaims struct {
    ID uint32    `json:"id"`
    IP string `json:"ip"`
    EX string `json:"ex"`
}

func CreateToken(id CustomClaims,key ...interface{}) (token string, err error) {
    claim := jwt.MapClaims{
        "id":  id.ID,
        "ip":  id.IP,
        "ex": time.Now().String(),
    }
    tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    token, err = tokens.SignedString(ks)
    return token, err
}

func ParseToken(tokens string) (id CustomClaims, err error){
    token,err:=jwt.Parse(tokens, secret())
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
