package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	p := Path("a/b/c/d/e/f/g")
	s := []interface{}{p.Parent(9)}
	for _, v := range s {
		fmt.Printf(">> %#v\n", v)
	}
	for _, v := range Path().Glob("*") {
		fmt.Printf(">>>> %#v\n", v.Path)
	}
}

// P 结构
type P struct {
	Path   string
	Name   string
	Stem   string
	Suffix string
}

// Path 获取新的path结构
func Path(paths ...interface{}) *P {
	pathlist := make([]string, 1)
	for _, v := range paths {
		switch v := v.(type) {
		case string:
			pathlist = append(pathlist, v)
		case []string:
			pathlist = append(pathlist, v...)
		case int:
			pathlist = append(pathlist, fmt.Sprint(v))
		case P:
			pathlist = append(pathlist, fmt.Sprint(v.Path))
		case *P:
			pathlist = append(pathlist, fmt.Sprint(v.Path))
		default:
			panic("Path 参数错误，必须[string, []string, int, P, *P]")
		}
	}
	pat := filepath.Join(pathlist...)
	return &P{
		Path:   pat,
		Name:   filepath.Base(pat),
		Stem:   strings.Split(filepath.Base(pat), ".")[0],
		Suffix: filepath.Ext(pat),
	}
}

// Item 返回字符串路径
func (p *P) Item() string {
	return p.Path
}

// Exists 存在
func (p *P) Exists() bool {
	_, err := os.Stat(p.Path)
	if err == nil {
		return true
	}
	return false
}

// Parent 获取n级父级目录结构(调用后返回绝对路径)
func (p *P) Parent(n int) *P {
	pat, err := filepath.Abs(p.Path)
	if err != nil {
		panic(err)
	}
	for i := 0; i < n; i++ {
		pat = filepath.Join(pat, "..")
	}
	return Path(pat)
}

// Makedirs 递归创建
func (p *P) Makedirs() {
	err := os.MkdirAll(p.Path, 0666)
	if err != nil {
		panic(err)
	}
}

// Glob 检索paths
func (p *P) Glob(pattern string) []*P {
	res := make([]*P, 0)
	lst, err := filepath.Glob(filepath.Join(p.Path, pattern))
	if err != nil {
		panic(err)
	}
	for _, v := range lst {
		res = append(res, Path(v))
	}
	// filepath.Walk(p.Path, func(pat string, info os.FileInfo, err error) error {
	// 	ok, err := filepath.Match(pattern, pat)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if ok {
	// 		res = append(res, Path(pat))
	// 	}
	// 	return nil
	// })
	return res
}
