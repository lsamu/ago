package secret

import (
    "crypto/md5"
    "crypto/rand"
    "crypto/sha1"
    "crypto/sha256"
    "fmt"
    "strconv"
)

//MD5 MD5
func MD5(s string) string  {
    h := md5.New()
    h.Write([]byte(s))
    cipher := h.Sum(nil)
    return fmt.Sprintf("%x", cipher)
}

// SHA1 SHA1
func SHA1(b string) string {
    h := sha1.New()
    _, _ = h.Write([]byte(b))
    return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA256 加密
func SHA256(str string) string {
    h := sha256.New()
    h.Write([]byte(str))
    r := h.Sum(nil)
    return fmt.Sprintf("%x", r)
}

//GenRandomCode 随机字符串
func GenRandomCode(n int) int64 {
    letterBytes := "123456789"
    buf := make([]byte, n)
    if _, err := rand.Read(buf); err != nil {
        return 0
    }
    for i := 0; i < n; {
        idx := int(buf[i] & 0x3F)
        if idx < len(letterBytes) {
            buf[i] = letterBytes[idx]
            i++
        } else {
            if _, err := rand.Read(buf[i : i+1]); err != nil {
                return 0
            }
        }
    }
    ss, _ := strconv.ParseInt(string(buf), 10, 32)
    return ss
}
