package main

import (
	"log"
	"strings"

	"github.com/tidwall/redcon"
	"github.com/tidwall/buntdb"
	"runtime"
	"flag"
	"fmt"
)

var (
	addr = flag.String("addr", ":6380", "地址,xxxx:xxxx")
	dbpath = flag.String("path", "", "存储路径，默认内存")
	keycols int
)

func init() {
	flag.Parse()
	if *dbpath == "" {
		*dbpath = ":memory:"
	}
}

func main() {
	db, err := buntdb.Open(*dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	go log.Printf("started server at %s", *addr)

	err = redcon.ListenAndServe(*addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
			case "ping":
				conn.WriteString("PONG")
			case "quit":
				conn.WriteString("OK")
				conn.Close()
			case "set":
				CmdSet(conn, db, cmd)
			case "get":
				CmdGet(conn, db, cmd)
			case "del":
				CmdDel(conn, db, cmd)
			case "keys":
				CmdKeys(conn, db, cmd)
			case "info":
				CmdInfo(conn, db, cmd)
			}
		},
		func(conn redcon.Conn) bool {
			// use this function to accept or deny the connection.
			// log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// this is called when the connection has been closed
			log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func CmdGet(conn redcon.Conn, db *buntdb.DB, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}

	db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(string(cmd.Args[1]))
		if err != nil {
			conn.WriteNull()
			return err
		}
		conn.WriteBulkString(val)
		return nil
	})
}

func CmdDel(conn redcon.Conn, db *buntdb.DB, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(string(cmd.Args[1]))
		if err != nil {
			conn.WriteInt(0)
			return nil
		}
		keycols--
		if keycols <= 0 {
			keycols = 0
		}
		conn.WriteInt(1)
		return err
	})
}

func CmdSet(conn redcon.Conn, db *buntdb.DB, cmd redcon.Command) {
	if len(cmd.Args) != 3 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(string(cmd.Args[1]), string(cmd.Args[2]), nil)
		if err != nil {
			conn.WriteString("NO")
			return nil
		}
		keycols++
		conn.WriteString("OK")
		return err
	})
}

func CmdKeys(conn redcon.Conn, db *buntdb.DB, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	db.View(func(tx *buntdb.Tx) error {
		var vals = make([]string, 0, 100)
		tx.AscendKeys(string(cmd.Args[1]), func(key, value string) bool {
			vals = append(vals, key)
			return true
		})
		if len(vals) > 0 {
			conn.WriteArray(len(vals))
			for _, v := range vals {
				conn.WriteBulkString(v)
			}
		} else {
			conn.WriteNull()
		}
		return nil
	})
}

func CmdInfo(conn redcon.Conn, db *buntdb.DB, cmd redcon.Command) {
	db.View(func(tx *buntdb.Tx) error {
		//内存
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		conn.WriteBulkString(`{"mem":` + fmt.Sprintf("%v", mem.HeapAlloc / 1024) +
			`,"keysize":` + fmt.Sprintf("%v", keycols) + "}\n")
		return nil
	})
}