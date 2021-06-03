package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

var password = make(chan string) //创建管道，接收密码
var isOver = make(chan bool)     //判断是否退出
func main() {
	if len(os.Args)==1{
		fmt.Println("缺少解压包名称参数")
		return
	}

	rootPath,_ := os.Getwd()
	rarPath :=rootPath+"/"+os.Args[1]+".rar"
	b,err :=PathExists(rarPath)
	if err != nil {
		fmt.Printf("PathExists(%s),err(%v)\n", rarPath, err)
		return
	}
	if !b {
		fmt.Printf("path %s 存在\n", rarPath)
		return
	}
	passPath :=rootPath+"/pass.txt"

	savePassPath :=rootPath+"/save.txt"
	go passtxt(passPath)

Loop:
	for {
		select {
		case rarpass := <-password:
			go cmdshell(rarPath, rarpass,savePassPath)
		case <-time.After(time.Second * time.Duration(1)):
			break Loop
		case <-isOver:
			break Loop
		}
	}

}

func cmdshell(rarpath string, pass string,savepass string) {

	cmd := exec.Command("unrar", "e", "-p"+pass, rarpath)
	_, err := cmd.Output()
	if err != nil {
		return
	}

	fmt.Printf("密码为：%s \n", pass)

	isOver <- true // 成功后退出
	file, e := os.OpenFile(savepass, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	if e != nil {
		fmt.Println("密码保存失败")
		return
	}
	file.Write([]byte(pass))
}

func passtxt(passpath string) {
	fp, _ := os.OpenFile(passpath, os.O_RDONLY, 6)
	defer fp.Close()

	// 创建文件的缓存区
	r := bufio.NewReader(fp)
	for {
		pass, err2 := r.ReadBytes('\n')
		if err2 == io.EOF { //文件末尾
			break
		}
		pass = pass[:len(pass)-2] // 去除末尾 /n
		password <- string(pass)
	}
}
func PathExists(path string) (bool,error) {
	_,err := os.Stat(path)
	if err == nil {
		return true,nil
	}
	if os.IsNotExist(err) {
		return false,nil
	}
	return false,err
}