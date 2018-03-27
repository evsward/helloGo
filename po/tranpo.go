package po

import "fmt"

type Tran struct {
	id               int
	block_number     int
	block_hash       string
	transaction_hash string
	block_timestamp  string
	from_account     string
	to_account       string
	tran_value       int
	tran_gas_price   int
	tran_gas_amount  int
	tran_nonce       int
}

func NewTran(id int, block_number int, tran_value int, tran_gas_amount int, tran_gas_price int, tran_nonce int, block_hash string, transaction_hash string, block_timestamp string, from_account string, to_account string) Tran {
	var result = new(Tran)
	result.id = id
	result.block_number = block_number
	result.block_hash = block_hash
	result.transaction_hash = transaction_hash
	result.block_timestamp = block_timestamp
	result.from_account = from_account
	result.to_account = to_account
	result.tran_value = tran_value
	result.tran_gas_price = tran_gas_price
	result.tran_gas_amount = tran_gas_amount
	result.tran_nonce = tran_nonce
	return *result
}

func (tran Tran) String() string {
	return fmt.Sprintf("block_hash:%v, transaction_hash:%v, %v%v%v%v%v%v%v%v%v", tran.block_hash, tran.transaction_hash, tran.id, tran.block_number, tran.block_timestamp, tran.from_account, tran.to_account, tran.tran_gas_amount, tran.tran_gas_price, tran.tran_nonce, tran.tran_value)
}
