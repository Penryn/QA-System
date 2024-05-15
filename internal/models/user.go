package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AdminType int    `json:"admin_type"` //1:普通管理员	2:超级管理员
}
