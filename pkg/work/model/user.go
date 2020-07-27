package model

type Enterprise struct {
	Code   string
	Name   string
	Remark string
}

type UserInfo struct {
	AppId          string
	Uid            string
	Name           string
	Code           string
	EnterpriseCode string
	EnterpriseName string
}

type UsersList struct {
	Enterprise Enterprise
	Users      []*UserInfo
}
