package main

type Account struct {
	Nickname string
	Gmail    string
	Data     string
}

type GmailCode struct {
	Time int64
	Code string
}

type AccountToken struct {
	Login string
	Time  int64
}
