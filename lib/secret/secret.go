package secret

import (
    "crypto/md5"
    "crypto/sha1"
    "fmt"
)

func MD5(s string) string  {
    h := md5.New()
    h.Write([]byte(s))
    cipher := h.Sum(nil)
    return fmt.Sprintf("%x", cipher)
}

// SHA1
func SHA1(b string) string {
    h := sha1.New()
    _, _ = h.Write([]byte(b))
    return fmt.Sprintf("%x", h.Sum(nil))
}