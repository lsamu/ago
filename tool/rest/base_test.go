package rest

import "testing"

func TestCreateTemplate(t *testing.T) {
    tests := []struct {
        name string
    }{
        // TODO: Add test cases.
        {
            name: "OK",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            CreateTemplate()
        })
    }
}
