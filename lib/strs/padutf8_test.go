package strs

import (
    "testing"
    "testing/quick"
    "unicode/utf8"
)

func TestLeftEqualWithSameLengthUtf8(t *testing.T) {
    f := func(a string, pad string) bool {
        slen := utf8.RuneCountInString(a)
        padded := LeftUtf8(a, slen, pad)
        return padded == a
    }
    if err := quick.Check(f, nil); err != nil {
        t.Error(err)
    }
}

func TestRightEqualWithSameLengthUtf8(t *testing.T) {
    f := func(a string, pad string) bool {
        slen := utf8.RuneCountInString(a)
        padded := RightUtf8(a, slen, pad)
        return padded == a
    }
    if err := quick.Check(f, nil); err != nil {
        t.Error(err)
    }
}

func TestLeftEqualWithShorterLengthUtf8(t *testing.T) {
    f := func(a string, pad string) bool {
        slen := utf8.RuneCountInString(a)
        padded := LeftUtf8(a, slen-3, pad)
        return padded == a
    }
    if err := quick.Check(f, nil); err != nil {
        t.Error(err)
    }
}

func TestRightEqualWithShorterLengthUtf8(t *testing.T) {
    f := func(a string, pad string) bool {
        slen := utf8.RuneCountInString(a)
        padded := RightUtf8(a, slen-3, pad)
        return padded == a
    }
    if err := quick.Check(f, nil); err != nil {
        t.Error(err)
    }
}

func TestLeftEqualWithGreaterLengthUtf8(t *testing.T) {
    f := func(a string, pad string) bool {
        slen := utf8.RuneCountInString(a) + 3
        padded := LeftUtf8(a, slen, pad)
        return padded == times(pad, 3)+a
    }
    if err := quick.Check(f, nil); err != nil {
        t.Error(err)
    }
}

func TestRightEqualWithGreaterLengthUtf8(t *testing.T) {
    f := func(a string, pad string) bool {
        slen := utf8.RuneCountInString(a) + 3
        padded := RightUtf8(a, slen, pad)
        return padded == a+times(pad, 3)
    }
    if err := quick.Check(f, nil); err != nil {
        t.Error(err)
    }
}