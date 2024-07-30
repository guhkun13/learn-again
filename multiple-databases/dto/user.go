package dto

type CreateUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UpdateUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
