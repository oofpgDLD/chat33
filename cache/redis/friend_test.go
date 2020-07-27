package redis

import (
	"fmt"
	"testing"
)

func TestMytest1(t *testing.T) {
	conn := GetConn()
	v, err := conn.Do("EXISTS", "user14")
	fmt.Println(v)
	fmt.Println(err)
}

//删除好友
func TestDeleteFriend(t *testing.T) {
	err := c.DeleteFriend("14", "1")
	fmt.Println(err)
}

//添加好友
func TestAddFriend(t *testing.T) {
	err := c.AddFriend("14", "7")
	fmt.Println(err)
}

//获取好友信息
func TestGetFriendInfo(t *testing.T) {
	v, err := c.GetFriendInfo("14", "11")
	fmt.Println(v)
	fmt.Println(err)
}

//修改好友信息
func TestUpdateFriendInfo(t *testing.T) {
	err := c.UpdateFriendInfo("14", "11", "top", "qqq")
	fmt.Println(err)
}
