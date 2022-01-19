package mysql

import "testing"

func init() {
    InitGormMysql("root", "root123", "127.0.0.1", "3306", "sass_mall", true)
}

func TestSchema_Generate(t *testing.T) {
    s := &Schema{}
    s.TempateBasePath = "./go/"
    s.OutBasePath = "/home/lauxinyi/Desktop/code/"
    if err := s.Generate("sass_mall", ""); err != nil {
        t.Errorf("Generate() error = %v", err)
    }
}
