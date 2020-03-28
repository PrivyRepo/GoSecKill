package main

import (
	"database/sql"
	"encoding/json"
	"homework/common/mysql"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/homework?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

func main4() {
	http.HandleFunc("/testmysql", func(w http.ResponseWriter, rw *http.Request) {
		rows, err := db.Query("SELECT * FROM product")
		defer rows.Close()
		if err != nil {
			log.Println(err)
		}
		resultRows := mysql.GetResultRows(rows)
		bytes, err := json.Marshal(resultRows)
		if err != nil {
			log.Println(err)
		}
		w.Write(bytes)
	})
	http.ListenAndServe(":8888", nil)
}
