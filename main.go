package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	// 传参
	var dbFile string
	flag.StringVar(&dbFile, "f", "", "file path")
	flag.Parse()

	// 创建posts文件夹
	err := mkdirPostsDirectory()
	if err != nil {
		return
	}

	fmt.Println("Database file:", dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT title,text FROM typecho_contents")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var title, text string
		err = rows.Scan(&title, &text)
		if err != nil {
			log.Fatal(err)
		}

		// 处理<!--markdown-->字符串
		text = strings.Replace(text, "<!--markdown-->", "<!--markdown-->\n\r", -1)

		// 数据写入文件
		err := writeFile(title, text, 0666)
		if err != nil {
			return
		}
	}
}

func mkdirPostsDirectory() error {
	if _, err := os.Stat("posts"); os.IsNotExist(err) {
		err := os.Mkdir("posts", 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(filename string, data string, perm os.FileMode) error {
	file, err := os.OpenFile("./posts/"+filename+".md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	n, err := file.Write([]byte(data))
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	if err1 := file.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
