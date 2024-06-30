package models

type Owner struct {
	Id       int64
	Email    string
	Login    string
	PassHash []byte
}

// NewOwner TODO: хочу сделать паттерн конструктора через функции
// -> как в книжке 100 ошибок в го - хорошая практика будет
func NewOwner() Owner {
	return Owner{}
}
