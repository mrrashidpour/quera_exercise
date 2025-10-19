package main

import "fmt"

// اینترفیس حساب
type Account interface {
	MonthlyInterest() int
	Transfer(receiver Account, amount int) string
	Deposit(amount int) string
	Withdraw(amount int) string
	CheckBalance() int
}

type SavingsAccount struct {
	balance int
}

type CheckingAccount struct {
	balance int
}

type InvestmentAccount struct {
	balance int
}

func NewSavingsAccount() *SavingsAccount {
	return &SavingsAccount{balance: 0}
}

func NewCheckingAccount() *CheckingAccount {
	return &CheckingAccount{balance: 0}
}

func NewInvestmentAccount() *InvestmentAccount {
	return &InvestmentAccount{balance: 0}
}

func (a *SavingsAccount) MonthlyInterest() int {
	return a.balance * 5 / 100 / 12
}

func (a *SavingsAccount) Deposit(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	a.balance += amount
	return "Success"
}

func (a *SavingsAccount) Withdraw(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	if a.balance < amount {
		return "Account balance is not enough"
	}
	a.balance -= amount
	return "Success"
}

func (a *SavingsAccount) CheckBalance() int {
	return a.balance
}

func (a *SavingsAccount) Transfer(receiver Account, amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	switch receiver.(type) {
	case *SavingsAccount, *CheckingAccount, *InvestmentAccount:
	default:
		return "Invalid receiver account"
	}
	if a.balance < amount {
		return "Account balance is not enough"
	}
	a.balance -= amount
	receiver.Deposit(amount)
	return "Success"
}

func (a *CheckingAccount) MonthlyInterest() int {
	return a.balance * 1 / 100 / 12
}

func (a *CheckingAccount) Deposit(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	a.balance += amount
	return "Success"
}

func (a *CheckingAccount) Withdraw(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	if a.balance < amount {
		return "Account balance is not enough"
	}
	a.balance -= amount
	return "Success"
}

func (a *CheckingAccount) CheckBalance() int {
	return a.balance
}

func (a *CheckingAccount) Transfer(receiver Account, amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	switch receiver.(type) {
	case *SavingsAccount, *CheckingAccount, *InvestmentAccount:
	default:
		return "Invalid receiver account"
	}
	if a.balance < amount {
		return "Account balance is not enough"
	}
	a.balance -= amount
	receiver.Deposit(amount)
	return "Success"
}

func (a *InvestmentAccount) MonthlyInterest() int {
	return a.balance * 2 / 100 / 12
}

func (a *InvestmentAccount) Deposit(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	a.balance += amount
	return "Success"
}

func (a *InvestmentAccount) Withdraw(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	if a.balance < amount {
		return "Account balance is not enough"
	}
	a.balance -= amount
	return "Success"
}

func (a *InvestmentAccount) CheckBalance() int {
	return a.balance
}

func (a *InvestmentAccount) Transfer(receiver Account, amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	switch receiver.(type) {
	case *SavingsAccount, *CheckingAccount, *InvestmentAccount:
	default:
		return "Invalid receiver account"
	}
	if a.balance < amount {
		return "Account balance is not enough"
	}
	a.balance -= amount
	receiver.Deposit(amount)
	return "Success"
}

// ----------------------
// تابع main برای تست ساده (اختیاری)
// ----------------------

func main() {
	s := NewSavingsAccount()
	c := NewCheckingAccount()

	fmt.Println(s.Deposit(1000))     // Success
	fmt.Println(s.MonthlyInterest()) // 4 (تقریباً)
	fmt.Println(s.Transfer(c, 500))  // Success
	fmt.Println(c.CheckBalance())    // 500
	fmt.Println(s.CheckBalance())    // 500
}
