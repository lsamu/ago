package strs

import (
    "fmt"
    "github.com/google/uuid"
    "math/rand"
    "os"
    "strings"
    "sync/atomic"
    "time"
)

// GenCode 获取随机码
func GenCode(width int) string {
    numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    r := len(numeric)
    rand.Seed(time.Now().UnixNano())

    var sb strings.Builder
    for i := 0; i < width; i++ {
        fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
    }
    return sb.String()
}

var (
    version string
    incrNum uint64
    pid     = os.Getpid()
)

//TraceID TraceID
func TraceID() string {
    return fmt.Sprintf("trace-id-%d-%s-%d",
        os.Getpid(),
        time.Now().Format("2006.01.02.15.04.05.999"),
        atomic.AddUint64(&incrNum, 1))
}

// UUID Define alias
type UUID = uuid.UUID

// NewUUID Create uuid
func NewUUID() (UUID, error) {
    return uuid.NewRandom()
}

// MustUUID Create uuid(Throw panic if something goes wrong)
func MustUUID() UUID {
    v, err := NewUUID()
    if err != nil {
        panic(err)
    }
    return v
}

// MustString Create uuid
func MustString() string {
    return MustUUID().String()
}

