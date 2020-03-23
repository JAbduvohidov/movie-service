package token

type Payload struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}