package main

import (
	"container/list"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"hello/po"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(192.168.1.142:3306)/web3serve")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM ethereum_transactions")
	tranlist := list.New()
	log.Println(rows.Columns())
	for rows.Next() {
		var id int
		var block_number int
		var block_hash string
		var transaction_hash string
		var block_timestamp string
		var from_account string
		var to_account string
		var tran_value int
		var tran_gas_price int
		var tran_gas_amount int
		var tran_nonce int
		rows.Scan(&id)
		rows.Scan(&block_number)
		rows.Scan(&block_hash)
		rows.Scan(&transaction_hash)
		rows.Scan(&block_timestamp)
		rows.Scan(&from_account)
		rows.Scan(&to_account)
		rows.Scan(&tran_value)
		rows.Scan(&tran_gas_price)
		rows.Scan(&tran_gas_amount)
		rows.Scan(&tran_nonce)
		tran := po.NewTran(id, block_number, tran_value, tran_gas_amount, tran_gas_price, tran_nonce, block_hash, transaction_hash, block_timestamp, from_account, to_account)
		if err != nil {
			log.Fatalln(err)
		}
		tranlist.PushBack(tran)
		log.Printf("found row containing %q", tran)
	}
	rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
