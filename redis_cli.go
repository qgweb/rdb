// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// simple does nothing except block while running the service.
package main

import (
	"github.com/peterh/liner"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
	"io"
	"os"
	"github.com/qgweb/glib/convert"
	"golang.org/x/crypto/ssh"
	"time"
	"flag"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

var (
	shost = flag.String("shost", "192.168.1.199", "跳板机地址")
	sport = flag.String("sport", "22", "跳板机端口")
	suser = flag.String("suser", "root", "跳板机用户")
	spwd = flag.String("spwd", "qazwsxedc", "跳板机密码")
	config = flag.String("c", "config.json", "配置文件")
	rconn redis.Conn
)

func init() {
	flag.Parse()

	if d, err := ioutil.ReadFile(*config); err == nil {
		r := gjson.Parse(string(d))

		*shost = r.Get("shost").String()
		*sport = r.Get("sport").String()
		*suser = r.Get("suser").String()
		*spwd = r.Get("spwd").String()
		fmt.Println(*shost, *sport, *suser, *spwd)
	}
}
func main() {
	config := &ssh.ClientConfig{
		User: *suser,
		Auth: []ssh.AuthMethod{
			ssh.Password(*spwd),
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", *shost, *sport), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	fmt.Println("连接到跳板机成功，地址", *shost, ":", *sport)
	defer client.Close()

	line := liner.NewLiner()
	defer line.Close()
	line.SetMultiLineMode(false)
	line.SetCtrlCAborts(true)

	for {
		command, err := line.Prompt("zb" + "> ")
		line.AppendHistory(command)
		var cs = strings.Split(command, " ")
		var ncs = make([]interface{}, 0, len(cs))
		for _, v := range cs {
			if v == "" {
				continue
			}
			ncs = append(ncs, v)
		}

		if err == io.EOF {
			fmt.Fprintln(os.Stdout, "byte")
			return
		}
		if err == liner.ErrPromptAborted {
			fmt.Fprintln(os.Stdout, "prompt aborte")
			continue
		}
		if command == "exit" || command == "bye" {
			return
		}

		if convert.ToString(ncs[0]) == "connect" && len(ncs) >= 3 {
			cconn, err := client.Dial("tcp", fmt.Sprintf("%s:%s", ncs[1], ncs[2]))
			if err != nil {
				fmt.Fprintln(os.Stdout, "[error] " + err.Error())
				continue
			}
			if rconn != nil {
				fmt.Fprintln(os.Stdout, "老连接关闭")
				rconn.Close()
			}
			fmt.Fprintf(os.Stdout, "[%s:%s]连接成功\n", ncs[1], ncs[2])
			rconn = redis.NewConn(cconn, time.Second * 30, time.Second * 30)
			if len(ncs) == 4 {
				rconn.Do("AUTH", ncs[3])
			}
			defer rconn.Close()
			continue
		}

		if rconn == nil {
			fmt.Fprintln(os.Stdout, "[error] redis连接未连接")
			continue
		}

		var resp interface{}
		if len(ncs) > 1 {
			resp, err = rconn.Do(convert.ToString(ncs[0]), ncs[1:]...)
		} else {
			resp, err = rconn.Do(convert.ToString(ncs[0]))
		}

		if v, err := redis.Int(resp, err); err == nil {
			fmt.Fprintln(os.Stdout, v)
			continue
		}
		if v, err := redis.Int64(resp, err); err == nil {
			fmt.Fprintln(os.Stdout, v)
			continue
		}
		if v, err := redis.Ints(resp, err); err == nil {
			fmt.Fprintln(os.Stdout, v)
			continue
		}
		if v, err := redis.Uint64(resp, err); err == nil {
			fmt.Fprintln(os.Stdout, v)
			continue
		}
		if v, err := redis.String(resp, err); err == nil {
			fmt.Fprintln(os.Stdout, v)
			continue
		}
		if v, err := redis.Strings(resp, err); err == nil {
			for k, vv := range v {
				fmt.Fprintln(os.Stdout, fmt.Sprintf("[%d]", k), vv)
			}
			continue
		}
		if v, err := redis.Bool(resp, err); err == nil {
			fmt.Fprintln(os.Stdout, v)
			continue
		}
		if v, err := redis.Float64(resp, err); err == nil {
			fmt.Fprintln(os.Stdout, v)
			continue
		}
	}
}