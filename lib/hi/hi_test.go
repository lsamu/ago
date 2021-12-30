package hi

import (
    "fmt"
    "testing"
)

func TestHi(t *testing.T) {
    s, err := Hi(100, 1000, nil, func() (int, bool) {
        return 200, true
    })
    if err != nil {
        t.Log(err)
    }
    fmt.Println(s.Count())
    fmt.Println(s.CountSucc())
    fmt.Println(s.CostSussAvg())
    fmt.Println(s.TPS())
}
