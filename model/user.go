/*******
* @Author:qingmeng
* @Description:
* @File:user
* @Date:2022/7/25
 */

package model

import (
	"strconv"
)


// 用于测试的结构体
type UserService struct {
	Id   int
	Name string
	Age int
}

// 用于测试查询用户的方法
func (u *UserService)QueryUser(uid int) (*UserService, bool) {
	users := make([]UserService,10)
	for i := 0; i < len(users); i++ {
		users[i]=UserService{
			Id:   i,
			Name: "用户"+strconv.Itoa(i),
			Age:  20+i,
		}
	}

	for _, user := range users {
		if user.Id==uid{
			return &user,true
		}
	}

	return &UserService{}, false
}

func (u *UserService) Test() string {
	return "hello"
}