package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"log"
	"net"
	"os"
	"time"
)

// go run main.go --bind 0.0.0.0:99 --backend 127.0.0.1:8000
func main() {
	help := flag.Bool("help", false, "print usage")
	bind := flag.String("bind", "127.0.0.1:6000", "The address to bind to")
	backend := flag.String("backend", "", "The backend server address")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *backend == "" {
		flag.Usage()
		return
	}

	if *bind == "" {
		fmt.Println("use default bind")
	}

	success, err := RunProxy(*bind, *backend)
	if !success {
		fmt.Println("errrrrrr ", err)
		os.Exit(1)
	}
}

func RunProxy(bind, backend string) (bool, error) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return false, err
	}
	defer listener.Close()
	fmt.Println("tcp-proxy started.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("errrrrrr")
		} else {
			go ConnectionHandler(conn, backend)
		}
	}
}

func ConnectionHandler(conn net.Conn, backend string) {
	target, err := net.Dial("tcp", backend)
	defer conn.Close()
	if err != nil {
		fmt.Println("tcp-proxy start112312312ed.")
	} else {
		defer target.Close()
		closed := make(chan bool, 2)
		go Proxy(conn, target, closed)
		go Proxy(target, conn, closed)
		<-closed
	}
}

func Proxy(from net.Conn, to net.Conn, closed chan bool) {

	db, err := sql.Open("mysql",
		"testqcroo:1*UJIuJS@tcp(13.215.101.156:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tb_proxy_log(respone,writetime) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 4096)
	for {
		n1, err := from.Read(buffer)
		if err != nil {
			closed <- true
			return
		}
		resstr := string(buffer[0:n1])
		fmt.Printf("数据为>>>>>>>：\n\n\n %v", resstr)
		if resstr != "" {
			// GBK编码器
			gbkDecoder := simplifiedchinese.GBK.NewDecoder()

			// 使用GBK编码器进行解码
			decodedBytes, _, err := transform.Bytes(gbkDecoder, []byte(resstr))
			if err != nil {
				return
			}

			// 将解码后的字节转为字符串
			decodedText := string(decodedBytes)

			res, err := stmt.Exec(decodedText, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res)
		}

		n2, err := to.Write(buffer[:n1])
		fmt.Println(n2)
		if err != nil {
			closed <- true
			return
		}
	}
}
