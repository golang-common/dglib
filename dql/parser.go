/**
 * @Author: daipengyuan
 * @Description:
 * @File:  parser
 * @Version: 1.0.0
 * @Date: 2021/8/17 15:46
 */

package dql

type Person struct {
	Uid      string   `json:"uid" db:"uid" dtype:"Person"`
	Name     string   `json:"name" db:"name,string" index:"index" token:"exact"`
	NameNick string   `json:"nick_name" db:"name|nick,string" token:"term"`
	Age      int      `json:"age" db:"age,int" index:"index" token:"int"`
	Friend   []Person `json:"friend" db:"friend,uid" index:"reverse,count,list"`
	FriendOf []Person `json:"friend_of" db:"~friend,uid"`
}

