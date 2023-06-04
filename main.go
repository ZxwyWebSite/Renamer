package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// 测试阶段，没有错误处理，如有bug请上报
// ZxwyWebSite/Renamer-dev
// 一个快速对视频文件重命名的软件
// 为了解决挨个改名麻烦的问题
// 使用方式：
// 创建list.txt，每行一个文件名，最后一行不要回车
// 在需要改名的目录运行程序，检查操作是否正确，输入y运行
// 菜单样式：
// # 二次确认 | Renamer-CLI #
// '[Xxxx][01][720P][GB][MP4].mp4' => '第01话_Xxxx.mp4'
// ([y]确认|默认取消)> y
// [renamer] 操作成功完成
// [renamer] '[Xxxx][01][720P][GB][MP4].mp4' 重命名失败: err.Error()
// 附加功能：
// 1.指定扩展名(-e mp4)不加点
// 2.待开发...

var (
	extname    string // 扩展名
	runpath    string // 运行目录
	searchpath string // 搜索目录
)

func init() {
	runpath, _ = os.Getwd()
	flag.StringVar(&extname, `e`, ``, `指定文件扩展名(不加点)`)
	flag.StringVar(&searchpath, `p`, runpath, `指定扫描目录(相对位置)`)
	flag.Parse()
	if searchpath != runpath {
		// 输入 test，runpath=test!={path}，runpath={path}/test
		searchpath = runpath + `/` + searchpath
	}
}

func main() {
	var (
		file_list []string // 文件列表
		name_list []string // 名称列表
	)
	// 读取名称列表
	flist, err1 := os.Open(runpath + `/list.txt`)
	if err1 != nil {
		// fmt.Println(`请将 list.txt 放在运行目录`)
		// os.Exit(0)
		logs(`请将 list.txt 放在运行目录`)
	} else {
		defer flist.Close()
		br := bufio.NewReader(flist)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			name_list = append(name_list, string(a))
		}
	}
	// 读取搜索目录文件
	file, _ := os.ReadDir(searchpath)
	for _, f := range file {
		// 排除目录
		if !f.IsDir() {
			name := f.Name()
			// 是否限制扩展名
			if extname != `` {
				num := strings.LastIndex(name, `.`)
				if name[num+1:] != extname {
					continue
				}
			}
			// 将名称写入数组 (不符合的在前面用continue跳过)
			file_list = append(file_list, name)
		}
	}
	// 排序
	sort.Slice(
		file_list,
		func(i, j int) bool {
			return sortName(file_list[i]) < sortName(file_list[j])
		},
	)
	// 合并输出
	var out_list []string
	fmt.Print("\033[H\033[2J")
	fmt.Println(`# 二次确认 | Renamer-CLI #\n`)
	for num, _ := range name_list {
		out_list = append(out_list, file_list[num])
		fmt.Printf("'%v' => '%v'\n", file_list[num], name_list[num])
	}
	// 二次确认
	fmt.Print("\n确认修改?(y/*n)> ")
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	if string(data) != `y` {
		// fmt.Println(`取消操作`)
		// os.Exit(0)
		logs(`取消操作`)
	}
	// 文件重命名
	for i, v := range out_list {
		os.Rename(searchpath+`/`+v, searchpath+`/`+name_list[i])
	}
	//fmt.Println(`执行成功`)
	logs(`执行成功`)
	// 调试输出
	// fmt.Println(`文件列表：`)
	// for i, v := range file_list {
	// 	fmt.Printf("%v. %v\n", i+1, v)
	// }
	// fmt.Println(`名称列表：`)
	// for i, v := range name_list {
	// 	fmt.Printf("%v. %v\n", i+1, v)
	// }
}

// 文件按数字排序函数，来自 https://ask.csdn.net/questions/1025862
func sortName(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	i := len(name) - 1
	for ; i >= 0; i-- {
		if '0' > name[i] || name[i] > '9' {
			break
		}
	}
	i++
	b64 := make([]byte, 64/8)
	s64 := name[i:]
	if len(s64) > 0 {
		u64, err := strconv.ParseUint(s64, 10, 64)
		if err == nil {
			binary.BigEndian.PutUint64(b64, u64+1)
		}
	}
	return name[:i] + string(b64) + ext
}

// 退出
func logs(l string) {
	fmt.Println(`[renamer] ` + l)
	os.Exit(0)
}
