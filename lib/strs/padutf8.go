package strs

import "unicode/utf8"

func timesUtf8(str string, n int) (out string) {
    for i := 0; i < n; i++ {
        out += str
    }
    return
}

// LeftUtf8 left-pads the string with pad up to len runes
// len may be exceeded if
func LeftUtf8(str string, len int, pad string) string {
    return timesUtf8(pad, len-utf8.RuneCountInString(str)) + str
}

// RightUtf8 right-pads the string with pad up to len runes
func RightUtf8(str string, len int, pad string) string {
    return str + timesUtf8(pad, len-utf8.RuneCountInString(str))
}
