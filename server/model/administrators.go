package model

var SuperUser User

func InitSuperUser() {
	SuperUser = User{
		UserId:   -1,
		UserPwd:  "",
		UserName: "SuperUser",
	}
}
