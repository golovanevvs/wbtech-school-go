package main

import "fmt"

// modern interface for payment
type PaymentProcessor interface {
	Pay(amount int)
}

// modern bank implementing PaymentProcessor
type ModernBank struct {
	Name string
}

func (b *ModernBank) Pay(amount int) {
	fmt.Printf("[%s] Payment through modern bank: %d RUB\n", b.Name, amount/100)
}

// legacy bank with incompatible interface
type OldBank struct {
	Name string
}

func (b *OldBank) TransferFunds(amount int) {
	fmt.Printf("[%s] Payment through old bank: %d RUB\n", b.Name, amount/100)
}

// adapter
type OldBankAdapter struct {
	OldBank *OldBank
}

func (a *OldBankAdapter) Pay(amount int) {
	a.OldBank.TransferFunds(amount)
}

// client (works with PaymentProcessor interface only)
type Client struct{}

func (c *Client) MakePayment(p PaymentProcessor, amount int) {
	p.Pay(amount)
}

func main() {
	client := &Client{}

	modernBank := &ModernBank{Name: "WB-bank"}
	client.MakePayment(modernBank, 5500000)

	oldBank := &OldBank{Name: "Sber"}
	oldBankAdapter := &OldBankAdapter{OldBank: oldBank}
	client.MakePayment(oldBankAdapter, 500)
}
