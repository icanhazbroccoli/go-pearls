package main

import "fmt"

var deposits = make(chan int)
var balances = make(chan int)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
			fmt.Printf("balance: %d\n", balance)
		case balances <- balance:
			fmt.Printf("publishing a balance: %d\n", balance)
		}
	}
}

func init() {
	go teller()
}

func main() {
	Deposit(5)
	Deposit(10)
	Deposit(15)
	fmt.Println(Balance())
}
