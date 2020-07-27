package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

//检查markId是否唯一
func CheckRoomMarkIdExist(id string) (bool, error) {
	const sqlStr = "select * from room where mark_id = ?"
	maps, err := conn.Query(sqlStr, id)
	if err != nil {
		return false, err
	}
	if len(maps) > 0 {
		return true, nil
	}
	return false, nil
}

// 获取群成员信息列表
func GetRoomMembers(roomId string, searchNumber int) ([]map[string]string, error) {
	if searchNumber == types.SearchAll {
		const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_user.room_id = ? and room_user.is_delete = ? order by room_user.`level` desc,room_user.create_time desc"
		return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
	} else {
		const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_user.room_id = ? and room_user.is_delete = ? order by room_user.`level` desc,room_user.create_time desc limit 0,?"
		return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted, searchNumber)
	}
}

// 获取群中管理员和群主信息
func GetRoomManagerAndMaster(roomId string) ([]map[string]string, error) {
	const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_user.room_id = ? and room_user.`level` > 1 and room_user.is_delete = ?"
	return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
}

//获取群中某成员信息
func GetRoomMemberInfo(roomId, userId string, isDel int) ([]map[string]string, error) {
	switch isDel {
	case types.RoomUserDeletedOrNot:
		//const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_user.room_id = ? and room_user.user_id = ?"
		const sqlStr = "select ru.id,ru.room_id,ru.user_id,ru.user_nickname,ru.`level`,ru.no_disturbing,ru.common_use,ru.room_top,ru.create_time,ru.is_delete,ru.source,u.username,u.avatar,u.public_key,u.identification,u.identification_info " +
			"from room_user as ru left join user as u on ru.user_id = u.user_id " +
			"where ru.room_id = ? and ru.user_id = ?"
		return conn.Query(sqlStr, roomId, userId)
	case types.RoomUserNotDeleted:
		fallthrough
	case types.RoomUserDeleted:
		//const sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_user.room_id = ? and room_user.user_id = ? and room_user.is_delete = ?"
		const sqlStr = "select ru.id,ru.room_id,ru.user_id,ru.user_nickname,ru.`level`,ru.no_disturbing,ru.common_use,ru.room_top,ru.create_time,ru.is_delete,ru.source,u.username,u.avatar,u.public_key,u.identification,u.identification_info " +
			"from room_user as ru left join user as u on ru.user_id = u.user_id " +
			"where ru.room_id = ? and ru.user_id = ? and ru.is_delete = ?"
		return conn.Query(sqlStr, roomId, userId, isDel)
	default:
		return nil, errors.New("获取群成员信息参数错误")
	}
}

//获取群中某成员信息
func GetRoomMemberInfoByName(roomId, name string) ([]map[string]string, error) {
	name = QueryStr(name)
	//var sqlStr = "select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_user.room_id = ? and room_user.is_delete = ? and (room_user.user_nickname LIKE '%" + name + "%' or `user`.username LIKE '%" + name + "%')"
	var sqlStr = fmt.Sprintf("select * from room_user left join `user` on room_user.user_id = `user`.user_id where room_user.room_id = ? and room_user.is_delete = ? and (room_user.user_nickname LIKE '%%%s%%' or `user`.username LIKE '%%%s%%')", name, name)
	return conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
}

// 获取群成员总数
func GetRoomMemberNumber(roomId string) (int64, error) {
	const sqlStr = "select count(*) as count from room_user where room_id = ? and is_delete = ?"
	maps, err := conn.Query(sqlStr, roomId, types.RoomUserNotDeleted)
	if err != nil || len(maps) == 0 {
		return 0, err
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

// 获取群中相应角色的数量
func GetRoomMemberNumberByLevel(roomId string, level int) (int64, error) {
	const sqlStr = "select count(*) as count from room_user as ru left join `user` as u on ru.user_id = u.user_id where ru.room_id = ? and ru.is_delete = ? and ru.`level` = ?"
	maps, err := conn.Query(sqlStr, roomId, types.RoomUserNotDeleted, level)
	if err != nil || len(maps) == 0 {
		return 0, err
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

//获取群列表
func GetRoomList(user string, Type int) ([]map[string]string, error) {
	var sql = "select * from room_user left join room on room_user.room_id = room.id where room_user.user_id = ? and room_user.is_delete = ? and room.is_delete = ?"
	switch Type {
	case 1:
		sql += " and common_use = " + utility.ToString(types.UncommonUse)
	case 2:
		sql += " and common_use = " + utility.ToString(types.CommonUse)
	case 3:
	}
	sql += " order by name"
	return conn.Query(sql, user, types.RoomUserNotDeleted, types.RoomNotDeleted)
}

// 获取所有群
func GetEnabledRooms() ([]map[string]string, error) {
	const sqlStr = "select * from room where is_delete = ?"
	return conn.Query(sqlStr, types.RoomNotDeleted)
}

// 获取群详情
func GetRoomsInfo(roomId string, isDel int) ([]map[string]string, error) {
	if isDel == types.RoomDeletedOrNot {
		const sqlStr = "select `id`,mark_id,`name`,avatar,master_id,create_time,encrypt,can_add_friend,join_permission,record_permission,admin_muted,master_muted,is_delete,close_until,identification,identification_info " +
			"from room where room.id = ?"
		return conn.Query(sqlStr, roomId)
	}
	const sqlStr = "select `id`,mark_id,`name`,avatar,master_id,create_time,encrypt,can_add_friend,join_permission,record_permission,admin_muted,master_muted,is_delete,close_until,identification,identification_info " +
		"from room where room.id = ? and is_delete = ?"
	return conn.Query(sqlStr, roomId, isDel)
}

// 获取群详情通过markId
func GetRoomsInfoByMarkId(markId string, isDel int) ([]map[string]string, error) {
	const sqlStr = "select * from room where room.mark_id = ? and is_delete = ?"
	return conn.Query(sqlStr, markId, isDel)
}

// 根据userId获取加入的所有群
func GetRoomsById(id string) ([]map[string]string, error) {
	const sqlStr = "select * from room_user left join `room` on room_id = `room`.id where room_user.user_id = ? and room_user.is_delete = ? and `room`.is_delete = ?"
	return conn.Query(sqlStr, id, types.RoomUserNotDeleted, types.RoomUserNotDeleted)
}

// 获取公告数目
func GetRoomSystemMsgNumber(roomId string) (int64, error) {
	const sqlStr = "select count(*) as count from room_msg_content where room_id = ? and msg_type = ? and is_delete = ?"
	maps, err := conn.Query(sqlStr, roomId, types.System, types.RoomUserNotDeleted)
	if err != nil || len(maps) == 0 {
		return 0, err
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

// 添加群聊聊天日志
func AppendRoomChatLog(senderId, receiveId, msgId string, msgType, isSnap int, content, ext string, time int64) (int64, int64, error) {
	const sqlStr = "insert into room_msg_content(room_id,sender_id,msg_id,msg_type,content,ext,datetime,is_delete,is_snap) values(?,?,?,?,?,?,?,?,?)"
	return conn.Exec(sqlStr, receiveId, senderId, msgId, msgType, content, ext, time, types.RoomMsgNotDelete, isSnap)
}

// 添加群聊接收日志
func AppendRoomMemberReceiveLog(roomLogId, receiver string, state int) (int64, int64, error) {
	const sqlStr = "insert into room_msg_receive(id,room_msg_id,receive_id,state) values(?,?,?,?)"
	return conn.Exec(sqlStr, nil, roomLogId, receiver, state)
}

// 根据LogId修改消息接收状态
func AlertRoomRevStateByRevId(revLogId string, state int) (int64, int64, error) {
	const sqlStr = "update room_msg_receive set state = ? where id = ?"
	return conn.Exec(sqlStr, state, revLogId)
}

func GetRoomChatLogsByUserId(roomId, owner string, startId, joinTime int64, number int, queryType []string) ([]map[string]string, string, error) {
	condition := ""
	if startId > 0 {
		condition += "and id <= " + utility.ToString(startId)
	}
	if owner != "" {
		condition += " and sender_id = " + owner
	}
	if len(queryType) > 0 {
		condition += fmt.Sprintf(" and msg_type in(%s)", strings.Join(queryType, ","))
	}
	nextLogId := "-1"
	sqlStr := fmt.Sprintf(
		"select * from room_msg_content left join `user` on sender_id = `user`.user_id where room_id = ? %s and is_delete = ? and datetime >= ? order by id desc limit 0,?",
		condition)
	maps, err := conn.Query(sqlStr, roomId, types.RoomMsgNotDelete, joinTime, number+1)
	if err != nil {
		return nil, nextLogId, err
	}
	if len(maps) > number {
		nextLogId = utility.ToString(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	return maps, nextLogId, nil
}

//获取聊天消息 startLogId 0:从最新一条消息开始 大于0:从startLogId开始
func GetChatlog(roomId string, startLogId, joinTime int64, number int) ([]map[string]string, string, error) {
	var sqlStr string
	var maps []map[string]string
	var err error
	if startLogId == 0 {
		sqlStr = "select * from room_msg_content left join `user` on sender_id = `user`.user_id where room_id = ? and is_delete = ? and datetime >= ? order by id desc limit 0,?"
		maps, err = conn.Query(sqlStr, roomId, types.RoomMsgNotDelete, joinTime, number+1)
	} else {
		sqlStr = "select * from room_msg_content left join `user` on sender_id = `user`.user_id where room_id = ? and id <= ? and is_delete = ? and datetime >= ? order by id desc limit 0,?"
		maps, err = conn.Query(sqlStr, roomId, startLogId, types.RoomMsgNotDelete, joinTime, number+1)
	}

	nextLogId := "-1"
	if err != nil {
		return nil, nextLogId, err
	}
	if len(maps) > number {
		nextLogId = utility.ToString(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	return maps, nextLogId, err
}

//获群公告
func GetSystemMsg(roomId string, startLogId int64, number int) ([]map[string]string, string, error) {
	var sqlStr string
	var maps []map[string]string
	var err error
	if startLogId == 0 {
		sqlStr = "SELECT * FROM `room_msg_content` where room_id = ? and msg_type = ? and is_delete = ? order by id desc limit 0,?"
		maps, err = conn.Query(sqlStr, roomId, types.System, types.RoomMsgNotDelete, number+1)
	} else {
		sqlStr = "SELECT * FROM `room_msg_content` where room_id = ? and msg_type = ? and id <= ? and is_delete = ? order by id desc limit 0,?"
		maps, err = conn.Query(sqlStr, roomId, types.System, startLogId, types.RoomMsgNotDelete, number+1)
	}

	nextLogId := "-1"
	if len(maps) > number {
		nextLogId = utility.ToString(maps[len(maps)-1]["id"])
		return maps[:len(maps)-1], nextLogId, err
	}
	return maps, nextLogId, err
}

//根据logId获取消息记录
func GetRoomMsgContent(id string) ([]map[string]string, error) {
	const sqlStr = "SELECT rmc.id,rmc.room_id,rmc.sender_id,rmc.is_snap,rmc.msg_id,rmc.msg_type,rmc.content,rmc.datetime,rmc.is_delete,u.username,u.avatar,u.uid " +
		"from room_msg_content as rmc left join `user` as u on rmc.sender_id = u.user_id " +
		"WHERE rmc.id = ?"
	return conn.Query(sqlStr, id)
}

//根据MsgId获取消息记录
func GetRoomMsgContentByMsgId(senderId, msgId string) ([]map[string]string, error) {
	const sqlStr = "SELECT rmc.id,rmc.room_id,rmc.sender_id,rmc.is_snap,rmc.msg_id,rmc.msg_type,rmc.content,rmc.datetime,rmc.is_delete,u.username,u.avatar " +
		"from room_msg_content as rmc left join `user` as u on rmc.sender_id = u.user_id " +
		"WHERE rmc.msg_id = ? and rmc.sender_id = ?"
	return conn.Query(sqlStr, msgId, senderId)
}

func GetRoomMsgContentAfter(rooms []string, time int64, isDel int) ([]map[string]string, error) {
	sqlStr := fmt.Sprintf("SELECT rmc.id,rmc.room_id,rmc.sender_id,rmc.is_snap,rmc.msg_id,rmc.msg_type,rmc.content,rmc.datetime,rmc.is_delete,u.username,u.avatar "+
		"from room_msg_content as rmc left join `user` as u on rmc.sender_id = u.user_id "+
		"WHERE rmc.room_id in(%s) and rmc.datetime > ? and rmc.is_delete = ? order by rmc.datetime ASC", strings.Join(rooms, ","))
	return conn.Query(sqlStr, time, isDel)
}

//reutrn [begin,end)
func GetRoomMsgContentBetween(rooms []string, begin, end int64, isDel int) ([]map[string]string, error) {
	sqlStr := fmt.Sprintf("SELECT rmc.id,rmc.room_id,rmc.sender_id,rmc.is_snap,rmc.msg_id,rmc.msg_type,rmc.content,rmc.datetime,rmc.is_delete,u.username,u.avatar "+
		"from room_msg_content as rmc left join `user` as u on rmc.sender_id = u.user_id "+
		"WHERE rmc.room_id in(%s) and rmc.datetime >= ? and rmc.datetime < ? and rmc.is_delete = ? order by rmc.datetime DESC", strings.Join(rooms, ","))
	return conn.Query(sqlStr, begin, end, isDel)
}

//reutrn [begin,end)
func GetRoomMsgContentNumberBetween(rooms []string, begin, end int64, isDel int) ([]map[string]string, error) {
	sqlStr := fmt.Sprintf("SELECT count(*) as count "+
		"from room_msg_content as rmc left join `user` as u on rmc.sender_id = u.user_id "+
		"WHERE rmc.room_id in(%s) and rmc.datetime >= ? and rmc.datetime < ? and rmc.is_delete = ?", strings.Join(rooms, ","))
	return conn.Query(sqlStr, begin, end, isDel)
}

func GetRoomRevMsgByLogId(logId, userId string) ([]map[string]string, error) {
	const sqlStr = "select * from room_msg_receive where room_msg_id = ? and receive_id = ?"
	return conn.Query(sqlStr, logId, userId)
}

func GetRoomMsgBurntNumber(logId string) (int, error) {
	const sqlStr = "select count(*) as count from room_msg_receive where room_msg_id = ? and state = ?"
	maps, err := conn.Query(sqlStr, logId, types.HadBurnt)
	count := 0
	if len(maps) > 1 {
		count = utility.ToInt(maps[0]["count"])
	}
	return count, err
}

// 获取群禁言数量
func GetRoomMutedListNumber(roomId string) ([]map[string]string, error) {
	const sqlStr = "SELECT COUNT(*) AS count FROM room_user_muted AS rm LEFT JOIN `room_user` ru ON rm.room_id = ru.room_id and rm.user_id = ru.user_id" +
		" WHERE rm.room_id = ? and ru.`level` = ? and rm.list_type != ? and ru.is_delete = ?"
	return conn.Query(sqlStr, roomId, types.RoomLevelNomal, types.AllSpeak, types.RoomUserNotDeleted)
}

// 获取群禁言数量 事务
func GetRoomMutedListNumberByTx(tx *mysql.MysqlTx, roomId string) ([]map[string]string, error) {
	const sqlStr = "SELECT COUNT(*) AS count FROM room_user_muted AS rm LEFT JOIN `room_user` ru ON rm.room_id = ru.room_id and rm.user_id = ru.user_id" +
		" WHERE rm.room_id = ? and ru.`level` = ? and rm.list_type != ? and ru.is_delete = ?"
	return tx.Query(sqlStr, roomId, types.RoomLevelNomal, types.AllSpeak, types.RoomUserNotDeleted)
}

// 获取群禁言类型
func GetRoomMutedType(roomId string) ([]map[string]string, error) {
	const sqlStr = "SELECT master_muted from room WHERE id = ?"
	return conn.Query(sqlStr, roomId)
}

// 设置群禁言类型
func SetRoomMutedType(tx *mysql.MysqlTx, roomId string, mutedType int) (int64, int64, error) {
	const sqlStr = "update room set master_muted = ? where id = ?"
	return tx.Exec(sqlStr, mutedType, roomId)
}

//取消成员禁言
func DelRoomUserMuted(tx *mysql.MysqlTx, roomId, userId string) (int64, int64, error) {
	const sqlStr = "update room_user_muted set list_type = ? where room_id = ? and user_id = ?"
	return tx.Exec(sqlStr, types.AllSpeak, roomId, userId)
}

//清空禁言
func ClearRoomMutedList(tx *mysql.MysqlTx, roomId string) (int64, int64, error) {
	const sqlStr = "update room_user_muted set list_type = ?, deadline = ? where room_id = ?"
	return tx.Exec(sqlStr, types.AllSpeak, 0, roomId)
}

//设置成员禁言
func AddRoomUserMuted(tx *mysql.MysqlTx, roomId, userId string, mutedType int, deadline int64) (int64, int64, error) {
	const sqlStr = "insert into room_user_muted values(?,?,?,?,?) ON DUPLICATE KEY UPDATE list_type = ?,deadline = ?"
	return tx.Exec(sqlStr, nil, roomId, userId, mutedType, deadline, mutedType, deadline)
}

//获取成员禁言信息
func GetRoomUserMuted(roomId, userId string) ([]map[string]string, error) {
	const sqlStr = "select rm.id,rm.room_id,rm.user_id,rm.list_type,rm.deadline " +
		"from room_user_muted rm left join `room_user` ru on rm.room_id = ru.room_id and rm.user_id = ru.user_id " +
		"where rm.room_id = ? and rm.user_id = ?  and ru.is_delete = ? and ru.`level` not in (?,?)"
	return conn.Query(sqlStr, roomId, userId, types.RoomUserNotDeleted, types.RoomLevelManager, types.RoomLevelMaster)
}

//获取群的某种禁言信息列表
func GetRoomUsersMutedInfo(roomId string, mutedType int) ([]map[string]string, error) {
	const sqlStr = "select rm.id,rm.room_id,rm.user_id,rm.list_type,rm.deadline " +
		"from room_user_muted rm left join `room_user` ru on rm.room_id = ru.room_id and rm.user_id = ru.user_id " +
		"where rm.room_id = ? and rm.list_type = ? and ru.is_delete = ? and ru.`level` not in (?,?)"
	return conn.Query(sqlStr, roomId, mutedType, types.RoomUserNotDeleted, types.RoomLevelManager, types.RoomLevelMaster)
}

func FindUserCreateRoomsNumber(userId string) (int, error) {
	const sqlStr = "select count(*) as count from room where master_id = ? and is_delete = ?"
	maps, err := conn.Query(sqlStr, userId, types.RoomNotDeleted)
	if err != nil {
		return 0, err
	}
	return utility.ToInt(maps[0]["count"]), nil
}

// 创建房间 返回 roomId
func CreateRoom(creater, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, members []string, randomRoomId string, createTime int64) (int64, error) {
	tx, err := conn.NewTx()
	if err != nil {
		return 0, err
	}
	const insertRoomSql = "insert into room(id,mark_id,name,avatar,master_id,create_time,encrypt,can_add_friend,join_permission,record_permission,admin_muted,master_muted,is_delete) values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, _, err = tx.Exec(insertRoomSql, nil, randomRoomId, roomName, roomAvatar, creater, createTime, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, types.RoomNotDeleted)
	if err != nil {
		tx.RollBack()
		return 0, err
	}
	maps, err := tx.Query("select LAST_INSERT_ID() as id")

	if err != nil || len(maps) < 1 {
		tx.RollBack()
		return 0, err
	}
	roomId := utility.ToInt64(maps[0]["id"])

	//
	const insertMemberSql = "insert into room_user(id,room_id,user_id,user_nickname,level,no_disturbing,common_use,room_top,create_time,is_delete,source) values(?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE is_delete = ?,create_time = ?"
	_, _, err = tx.Exec(insertMemberSql, nil, roomId, creater, "", types.RoomLevelMaster, types.NoDisturbingOff, types.UncommonUse, types.NotOnTop, createTime, types.RoomUserNotDeleted, "", types.RoomUserNotDeleted, createTime)
	if err != nil {
		tx.RollBack()
		return 0, err
	}
	for _, memberId := range members {
		_, _, err = tx.Exec(insertMemberSql, nil, roomId, memberId, "", types.RoomLevelNomal, types.NoDisturbingOff, types.UncommonUse, types.NotOnTop, createTime, types.RoomUserNotDeleted, "", types.RoomUserNotDeleted, createTime)
		if err != nil {
			tx.RollBack()
			return 0, err
		}
	}

	return roomId, tx.Commit()
}

// 创建房间 返回 roomId
func CreateRoomV2(tx *mysql.MysqlTx, creater, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, randomRoomId string, createTime int64) (int64, error) {
	const insertRoomSql = "insert into room(id,mark_id,name,avatar,master_id,create_time,encrypt,can_add_friend,join_permission,record_permission,admin_muted,master_muted,is_delete) values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, _, err := tx.Exec(insertRoomSql, nil, randomRoomId, roomName, roomAvatar, creater, createTime, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, types.RoomNotDeleted)
	if err != nil {
		return 0, err
	}
	maps, err := tx.Query("select LAST_INSERT_ID() as id")

	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, errors.New("can not find last id")
	}
	roomId := utility.ToInt64(maps[0]["id"])
	return roomId, nil
}

//删除群聊天记录
func DeleteRoomMsgContent(id string) (int, error) {
	sqlStr := "update room_msg_content set is_delete = ? where `id` = ?"
	num, _, err := conn.Exec(sqlStr, types.RoomMsgDeleted, id)
	return int(num), err
}

//删除群
func DeleteRoomById(roomId string) (int64, int64, error) {
	const sqlStr = "update room set is_delete = ? where id = ?"
	return conn.Exec(sqlStr, types.RoomDeleted, roomId)
}

// 入群申请，步骤1 添加user
func JoinRoomApproveStepInsert(tx *mysql.MysqlTx, roomId, userId, source string) (int64, int64, error) {
	createTime := utility.NowMillionSecond()
	const sqlStr = "insert into room_user values(?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE is_delete = ?,create_time = ?,source = ?"
	return tx.Exec(sqlStr, nil, roomId, userId, "", types.RoomLevelNomal, types.NoDisturbingOff, types.UncommonUse, types.NotOnTop, createTime, types.RoomUserNotDeleted, source, types.RoomUserNotDeleted, createTime, source)
}

// 入群申请，步骤2 更改状态
func JoinRoomApproveStepChangeState(tx *mysql.MysqlTx, logId int64, status int) (int64, int64, error) {
	const sqlStr = "update `apply` set `state` = ? where id = ?"
	return tx.Exec(sqlStr, status, logId)
}

// 添加群成员
func RoomAddMember(tx *mysql.MysqlTx, userId, roomId string, memberLevel int, createTime int64, source string) (int64, int64, error) {
	const sqlStr = "insert into room_user values(?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE is_delete = ?,create_time = ?,source = ?"
	return tx.Exec(sqlStr, nil, roomId, userId, "", memberLevel, types.NoDisturbingOff, types.UncommonUse, types.NotOnTop, createTime, types.RoomUserNotDeleted, source, types.RoomUserNotDeleted, createTime, source)
}

// 删除群成员
func DeleteRoomMemberById(userId, roomId string, time int64) (int64, int64, error) {
	const sqlStr = "update room_user set is_delete = ?,create_time = ?,user_nickname = ?,`level`= ?,no_disturbing = ?,common_use = ?,room_top = ? WHERE room_id = ? and user_id = ?"
	return conn.Exec(sqlStr, types.RoomUserDeleted, time, "", types.RoomLevelNomal, types.NoDisturbingOff, types.UncommonUse, types.NotOnTop, roomId, userId)
}

func GetJoinedRooms(userId string) ([]map[string]string, error) {
	const sqlStr = "select id,room_id,create_time,is_delete from room_user where user_id = ?"
	return conn.Query(sqlStr, userId)
}

// 修改是否可添加好友
func AlterRoomCanAddFriendPermission(roomId string, permisson int) error {
	if permisson != 0 {
		const sqlStr = "update room set can_add_friend = ? where id = ?"
		_, _, err := conn.Exec(sqlStr, permisson, roomId)
		if err != nil {
			return err
		}
	}
	return nil
}

// 修改群名称
func AlterRoomName(roomId, name string) (int64, int64, error) {
	const sqlStr = "update room set `name` = ? where id = ?"
	return conn.Exec(sqlStr, name, roomId)
}

// 修改群头像
func AlterRoomAvatar(roomId, avatar string) (int64, int64, error) {
	const sqlStr = "update room set avatar = ? where id = ?"
	return conn.Exec(sqlStr, avatar, roomId)
}

// 修改加入群权限
func AlterRoomJoinPermission(roomId string, permisson int) error {
	if permisson != 0 {
		const sqlStr = "update room set join_permission = ? where id = ?"
		_, _, err := conn.Exec(sqlStr, permisson, roomId)
		if err != nil {
			return err
		}
	}
	return nil
}

// 修改群拉取历史记录权限
func AlterRoomRecordPermission(roomId string, permisson int) error {
	if permisson != 0 {
		const sqlStr = "update room set record_permission = ? where id = ?"
		_, _, err := conn.Exec(sqlStr, permisson, roomId)
		if err != nil {
			return err
		}
	}
	return nil
}

// 设置群成员等级
func SetRoomMemberLevel(userId, roomId string, level int) (int64, int64, error) {
	const sqlStr = "update room_user set level = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, level, userId, roomId)
}

// 设置群免打扰
func SetRoomNoDisturbing(userId, roomId string, noDisturbing int) (int64, int64, error) {
	const sqlStr = "update room_user set no_disturbing = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, noDisturbing, userId, roomId)
}

// 设置群置顶
func SetRoomOnTop(userId, roomId string, onTop int) (int64, int64, error) {
	const sqlStr = "update room_user set room_top = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, onTop, userId, roomId)
}

// 群成员设置昵称
func SetMemberNickname(userId, roomId string, nickname string) (int64, int64, error) {
	const sqlStr = "update room_user set user_nickname = ? where user_id = ? and room_id = ?"
	return conn.Exec(sqlStr, nickname, userId, roomId)
}

// 转让群主
func SetNewMaster(master, userId, roomId string, level int) error {
	tx, err := conn.NewTx()
	if err != nil {
		return err
	}
	const setNewMasterSql = "update room_user set level = ? where user_id = ? and room_id = ?"
	_, _, err = tx.Exec(setNewMasterSql, level, userId, roomId)
	if err != nil {
		tx.RollBack()
		return err
	}
	const setRoomMasterSql = "update room set master_id = ? where id = ?"
	_, _, err = tx.Exec(setRoomMasterSql, userId, roomId)
	if err != nil {
		tx.RollBack()
		return err
	}
	const removeOldMasterSql = "update room_user set level = ? where user_id = ? and room_id = ?"
	_, _, err = tx.Exec(removeOldMasterSql, types.RoomLevelNomal, master, roomId)
	if err != nil {
		tx.RollBack()
		return err
	}
	return tx.Commit()
}

//查找群中消息免打扰的成员id列表   1：开启了免打扰，2：关闭
func FindRoomMemberSetNoDisturbing(roomId string, noDisturbing int) ([]map[string]string, error) {
	sql := "select user_id from room_user where room_id = ? and is_delete = ? and no_disturbing = ?"
	return conn.Query(sql, roomId, types.IsNotDelete, noDisturbing)
}

//查询群昵称 没有的话 返回用户名称
func FindRoomMemberName(roomId, userId string) ([]map[string]string, error) {
	sql := `SELECT IF(ISNULL(r.user_nickname)||LENGTH(r.user_nickname)<1,u.username,r.user_nickname) AS name FROM room_user r RIGHT JOIN user u ON r.user_id = u.user_id  WHERE r.room_id =? AND r.user_id = ?`
	return conn.Query(sql, roomId, userId)
}

//-------------------admin------------------//
//根据appId获取应用的所有未解散的群聊，包括被封群
func GetRoomCountInApp(appId string) (int64, error) {
	const sqlStr = "select count(*) as count from room left join user on room.master_id = user.user_id where room.is_delete = ? and user.app_id = ?"
	maps, err := conn.Query(sqlStr, types.RoomNotDeleted, appId)
	if err != nil || len(maps) == 0 {
		return 0, err
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

func GetCloseRoomCountInApp(appId string) (int64, error) {
	const sqlStr = "select count(*) as count from room left join user on room.master_id = user.user_id where room.is_delete = ? and user.app_id = ? and room.close_until > ?"
	maps, err := conn.Query(sqlStr, types.RoomNotDeleted, appId, utility.NowMillionSecond())
	if err != nil || len(maps) == 0 {
		return 0, err
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

//获取群的成员数上限
func GetRoomMembersLimit(appId string, level int) (int, error) {
	const sqlStr = "select * from room_config where level = ? and app_id = ?"
	maps, err := conn.Query(sqlStr, level, appId)
	if err != nil {
		return 0, err
	}
	if len(maps) <= 0 {
		return 0, nil
	}
	return utility.ToInt(maps[0]["number_limit"]), nil
}

//获取用户创建群个数上限
func GetCreateRoomsLimit(appId string, level int) (int, error) {
	const sqlStr = "select * from user_config where level = ? and app_id = ?"
	maps, err := conn.Query(sqlStr, level, appId)
	if err != nil {
		return 0, err
	}
	if len(maps) <= 0 {
		return 0, nil
	}
	return utility.ToInt(maps[0]["number_limit"]), nil
}

//设置用户创建群个数上限
func SetCreateRoomsLimit(appId string, level, limit int) error {
	const sqlStr = "insert into user_config(level,app_id,number_limit) values(?,?,?) ON DUPLICATE KEY UPDATE number_limit = ?"
	_, _, err := conn.Exec(sqlStr, level, appId, limit, limit)
	return err
}

//设置群的成员数上限
func SetRoomMembersLimit(appId string, level, limit int) error {
	const sqlStr = "insert into room_config(level,app_id,number_limit) values(?,?,?) ON DUPLICATE KEY UPDATE number_limit = ?"
	_, _, err := conn.Exec(sqlStr, level, appId, limit, limit)
	return err
}

//创建群个数 未解散
func FindCreateRoomNumbers(masterId string) (int, error) {
	const sqlStr = "select count(*) as count from room where master_id = ? and is_delete = ?"
	maps, err := conn.Query(sqlStr, masterId, types.RoomNotDeleted)
	if err != nil {
		return 0, err
	}
	if len(maps) <= 0 {
		return 0, nil
	}
	info := maps[0]
	return utility.ToInt(info["count"]), nil
}

func FindRoomsInAppQueryMarkId(appId, query string) ([]map[string]string, error) {
	sqlStr := ""
	if query != "" {
		query = QueryStr(query)
		//sqlStr = "SELECT `room`.*,`user`.user_id,`user`.account,`user`.uid from `room` left join `user` on room.master_id = user.user_id where is_delete = ? and user.app_id = ? and (room.mark_id LIKE '%" + query + "%' or room.name LIKE '%" + query + "%') order by room.create_time desc"
		sqlStr = fmt.Sprintf("SELECT `room`.*,`user`.user_id,`user`.account,`user`.uid from `room` left join `user` on room.master_id = user.user_id where is_delete = ? and user.app_id = ? and (room.mark_id LIKE '%%%s%%' or room.name LIKE '%%%s%%') order by room.create_time desc", query, query)
	} else {
		sqlStr = "SELECT `room`.*,`user`.user_id,`user`.account,`user`.uid from `room` left join `user` on room.master_id = user.user_id where is_delete = ? and user.app_id = ? order by room.create_time desc"
	}
	return conn.Query(sqlStr, types.RoomNotDeleted, appId)
}

//查找某个app所有下未封禁群
func FindRoomsInAppUnClose(appId string) ([]map[string]string, error) {
	sqlStr := "SELECT `room`.*,`user`.user_id,`user`.account,`user`.uid from `room` left join `user` on room.master_id = user.user_id where is_delete = ? and user.app_id = ? and room.close_until <= ? order by room.create_time desc"
	return conn.Query(sqlStr, types.RoomNotDeleted, appId, utility.NowMillionSecond())
}

//查找某个app所有下封禁群
func FindRoomsInAppClosed(appId string) ([]map[string]string, error) {
	sqlStr := "SELECT `room`.*,`user`.user_id,`user`.account,`user`.uid from `room` left join `user` on room.master_id = user.user_id where is_delete = ? and user.app_id = ? and room.close_until > ? order by room.create_time desc"
	return conn.Query(sqlStr, types.RoomNotDeleted, appId, utility.NowMillionSecond())
}

//根据群的发言人数递减排序
//查询区间：(datetime,now)
func RoomsOrderActiveMember(appId string, datetime int64) ([]map[string]string, error) {
	const sqlStr = "SELECT r.*,COUNT(DISTINCT rc.sender_id) " +
		"FROM `room_msg_content` AS rc LEFT JOIN room AS r ON rc.room_id = r.id LEFT JOIN `user` AS u ON u.user_id = r.master_id " +
		"WHERE u.app_id = ? AND rc.datetime > ? AND r.is_delete = ? AND r.close_until < ? " +
		"GROUP BY rc.room_id ORDER BY COUNT(DISTINCT rc.sender_id) DESC LIMIT 0,?"
	return conn.Query(sqlStr, appId, datetime, types.RoomNotDeleted, utility.NowMillionSecond(), 15)
}

//根据群的发言条数递减排序
func RoomsOrderActiveMsg(appId string, datetime int64) ([]map[string]string, error) {
	const sqlStr = "SELECT r.*, COUNT(rc.sender_id) " +
		"FROM `room_msg_content` AS rc LEFT JOIN room AS r ON rc.room_id = r.id LEFT JOIN `user` AS u ON u.user_id = r.master_id " +
		"WHERE u.app_id = ? AND rc.datetime > ? AND r.is_delete = ? AND r.close_until < ? " +
		"GROUP BY rc.room_id ORDER BY COUNT(rc.sender_id) DESC LIMIT 0,?"
	return conn.Query(sqlStr, appId, datetime, types.RoomNotDeleted, utility.NowMillionSecond(), 15)
}

//获取所有手动推荐群 0:非推荐群 1：推荐群
func FindAllRecommendRooms(appId string) ([]map[string]string, error) {
	const sqlStr = "SELECT r.*,u.user_id,u.account,u.uid FROM `room` AS r LEFT JOIN `user` AS u ON u.user_id = r.master_id " +
		"WHERE u.app_id = ? AND r.recommend = ? AND r.is_delete = ? AND r.close_until <= ? order by r.create_time desc"
	return conn.Query(sqlStr, appId, types.RoomRecommend, types.RoomNotDeleted, utility.NowMillionSecond())
}

//设置手动推荐群
func SetRecommendRoom(id string, recommend int) (int64, int64, error) {
	const sqlStr = "UPDATE `room` SET recommend = ? WHERE id = ?"
	return conn.Exec(sqlStr, recommend, id)
}

//设置为认证群
func SetRoomVerifyed(tx *mysql.MysqlTx, roomId, identificationInfo string) (int64, int64, error) {
	const sqlStr = "UPDATE `room` SET identification = ?, identification_info = ? WHERE id = ?"
	return tx.Exec(sqlStr, types.Verified, identificationInfo, roomId)
}
