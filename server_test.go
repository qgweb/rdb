package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/garyburd/redigo/redis"
)

func testDial(t *testing.T) redis.Conn {
	c, err := redis.Dial("tcp", ":6380")
	if err != nil {
		t.Fatal("连接6380端口失败,", err)
	}
	return c
}

func TestSet(t *testing.T) {
	conn := testDial(t)
	Convey("SET", t, func() {
		b, _ := redis.String(conn.Do("SET", "name", 111))
		So(b, ShouldEqual, "OK")
		conn.Close()
	})
}

func TestGet(t *testing.T) {
	conn := testDial(t)
	Convey("GET", t, func() {
		b, _ := redis.String(conn.Do("GET", "name"))
		So(b, ShouldEqual, "111")
		conn.Close()
	})
}

func TestDel(t *testing.T) {
	conn := testDial(t)
	Convey("DEL", t, func() {
		b, _ := redis.Bool(conn.Do("DEL", "name"))
		So(b, ShouldBeTrue)
		conn.Close()
	})
}