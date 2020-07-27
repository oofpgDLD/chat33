package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

/*
	获取好友列表, 无论是否删除都返回
*/
func FindFriendsAfterTime(id string, tp int, time int64) ([]map[string]string, error) {
	var sqlStr string
	if tp == 3 {
		sqlStr = "select f.user_id as F_user_id,u.user_id as U_user_id,f.*,u.* from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.add_time >= ?"
		return conn.Query(sqlStr, id, time)
	} else {
		sqlStr = "select f.user_id as F_user_id,u.user_id as U_user_id,f.*,u.* from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.type = ? and f.add_time >= ?"
		return conn.Query(sqlStr, id, strconv.Itoa(tp), time)
	}
}

func GetFriendList(id string, tp, isDelete int) ([]map[string]string, error) {
	var sqlStr string
	if tp == 3 {
		sqlStr = "select f.user_id as F_user_id,u.user_id as U_user_id,f.*,u.* from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.is_delete = ?"
		return conn.Query(sqlStr, id, isDelete)
	} else {
		sqlStr = "select f.user_id as F_user_id,u.user_id as U_user_id,f.*,u.* from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.type = ? and f.is_delete = ?"
		return conn.Query(sqlStr, id, strconv.Itoa(tp), isDelete)
	}
}

/*
	添加好友申请
*/
func AddFriendRequest(tp int, userID, friendID, reason, remark string) (bool, error) {
	//如果发送过申请 更新请求时间和原因（多次发送或者删除好友之后重新加）
	const sqlStr = "insert into `apply` values(?, ?, ?, ?, ?, ?, ?) on duplicate key update `state` = 0, `apply_reason` = ?, `datetime` = ?,`remark` = ?"
	now := utility.NowMillionSecond()
	_, _, err := conn.Exec(sqlStr, tp, userID, friendID, reason, types.AwaitState, remark, now, reason, now, remark)
	if err != nil {
		return false, err
	}
	return true, nil
}

//添加好友请求
func AddApply(tp int, userID, friendID, reason, remark string, source string) (int, error) {
	sqlStr := "insert into `apply` (type,apply_user,target,apply_reason,state,remark,datetime,source)values(?, ?, ?, ?, ?, ?, ?, ?)"
	now := utility.NowMillionSecond()
	_, id, err := conn.Exec(sqlStr, tp, userID, friendID, reason, types.AwaitState, remark, now, source)
	return int(id), err
}

/*
	同意好友申请
*/
func AcceptFriend(userID, friendID, source string, addTime int64) error {
	tx, err := conn.NewTx()
	if err != nil {
		return err
	}

	const updateStatus = "update `apply` set `state` = ?,datetime = ? where apply_user = ? and target = ? and type = ?"
	_, _, err = tx.Exec(updateStatus, types.AcceptState, utility.NowMillionSecond(), friendID, userID, types.IsFriend)
	if err != nil {
		tx.RollBack()
		return err
	}

	//添加好友
	const addFriend = "insert into `friends`(user_id,friend_id,remark,add_time,DND,`top`,`type`,is_delete,source) values(?, ?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update add_time = ?, DND = ?, top = ?, type = ?, is_delete = ?,source = ?"
	// id friend_id remark add_time, DND, top

	_, _, err = tx.Exec(addFriend, userID, friendID, "", addTime, types.NoDisturbingOff, types.NotOnTop, types.UncommonUse, types.FriendIsNotDelete, source, addTime, types.NoDisturbingOff, types.NotOnTop, types.UncommonUse, types.FriendIsNotDelete, source)

	if err != nil {
		tx.RollBack()
		return err
	}

	//读取请求的备注信息
	readRemark := "select * from apply where apply_user = ? and target = ? and type = ?"
	res, err := tx.Query(readRemark, friendID, userID, types.IsFriend)
	if err != nil {
		tx.RollBack()
		return err
	}
	remark := res[0]["remark"]

	//添加好友
	const addFriend1 = "insert into `friends`(user_id,friend_id,remark,add_time,DND,`top`,`type`,is_delete,source) values(?, ?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update remark = ?, add_time = ?, DND = ?, top = ?, type = ?, is_delete = ?,source = ?"
	_, _, err = tx.Exec(addFriend1, friendID, userID, remark, addTime, types.NoDisturbingOff, types.NotOnTop, types.UncommonUse, types.FriendIsNotDelete, source, remark, addTime, types.NoDisturbingOff, types.NotOnTop, types.UncommonUse, types.FriendIsNotDelete, source)
	if err != nil {
		tx.RollBack()
		return err
	}

	return tx.Commit()
}

/*
	拒绝好友申请
*/
func RejectFriend(userID, friendID string) error {
	const sqlStr = "update `apply` set `state` = ?,datetime = ? where apply_user = ? and target = ? and type = ?"
	_, _, err := conn.Exec(sqlStr, types.RejectState, utility.NowMillionSecond(), friendID, userID, types.IsFriend)
	return err
}

/*
	同意好友申请  但是不更新时间
*/
func AcceptFriendApply(userId, friendId string) (int64, error) {
	const sqlStr = "update `apply` set `state` = ? where ((apply_user = ? and target = ?) or (apply_user = ? and target = ?)) and type = ?"
	num, _, err := conn.Exec(sqlStr, types.AcceptState, friendId, userId, userId, friendId, types.IsFriend)
	return num, err
}

func InsertFriend(userID, friendID, remark, extRemark string, dnd, top int, addTime int64) (err error) {
	//添加好友
	const addFriend = "insert into `friends`(user_id,friend_id,remark,ext_remark,add_time,DND,`top`,`type`,is_delete) values(?, ?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update add_time = ?, DND = ?, top = ?, type = ?, is_delete = ?"
	// id friend_id remark add_time, DND, top
	_, _, err = conn.Exec(addFriend, userID, friendID, remark, extRemark, addTime, dnd, top, types.UncommonUse, types.FriendIsNotDelete, addTime, dnd, top, types.UncommonUse, types.FriendIsNotDelete)
	return
}

/*
	设置好友备注
*/
func SetFriendRemark(userID, friendID, remark string) error {
	const sqlStr = "update `friends` set `remark` = ? where user_id = ? and friend_id = ?"
	_, _, err := conn.Exec(sqlStr, remark, userID, friendID)
	return err
}

//设置好友详细备注
func SetFriendExtRemark(userID, friendID, remark string) error {
	const sqlStr = "update `friends` set `ext_remark` = ? where user_id = ? and friend_id = ?"
	_, _, err := conn.Exec(sqlStr, remark, userID, friendID)
	return err
}

/*
	设置好友免打扰
*/
func SetFriendDND(userID, friendID string, DND int) error {
	const sqlStr = "update `friends` set `DND` = ? where user_id = ? and friend_id = ?"
	_, _, err := conn.Exec(sqlStr, DND, userID, friendID)
	return err
}

/*
	设置好友置顶
*/
func SetFriendTop(userID, friendID string, top int) error {
	const sqlStr = "update `friends` set `top` = ? where user_id = ? and friend_id = ?"
	_, _, err := conn.Exec(sqlStr, top, userID, friendID)
	return err
}

//删除好友
func DeleteFriend(userID, friendID string, alterTime int64) error {
	const sqlStr = "update `friends` set is_delete = ?,add_time = ? where user_id = ? and friend_id = ?"
	tx, err := conn.NewTx()
	if err != nil {
		return err
	}

	_, _, err = tx.Exec(sqlStr, types.FriendIsDelete, alterTime, userID, friendID)
	if err != nil {
		tx.RollBack()
		return err
	}
	_, _, err = tx.Exec(sqlStr, types.FriendIsDelete, alterTime, friendID, userID)
	if err != nil {
		tx.RollBack()
		return err
	}
	return tx.Commit()
}

// 检查是否是好友关系
func CheckFriend(userID, friendID string, isDelete int) (bool, error) {
	sqlStr := "select * from `friends` where user_id = ? and friend_id = ? and is_delete = ?"
	rows, err := conn.Query(sqlStr, userID, friendID, isDelete)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

//查询好友请求是否存在
func FindFriendRequest(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT `state` FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	return conn.Query(sqlStr, friendID, userID, types.IsFriend)
}

//查询好友请求信息
func FindFriendRequestInfo(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	return conn.Query(sqlStr, friendID, userID, types.IsFriend)
}

//查询好友请求数量
func FindApplyCount(userID, friendID string) (int32, error) {
	sqlStr := "SELECT *  FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	rows, err := conn.Query(sqlStr, userID, friendID, types.IsFriend)
	if err != nil {
		return 0, err
	}
	return int32(len(rows)), nil
}

//查询好友请求数量
func FindApplyId(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT id FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	return conn.Query(sqlStr, userID, friendID, types.IsFriend)
}

//查询好友来源
func FindApplySource(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT source  FROM apply WHERE apply_user = ? AND target = ? and type = ?"
	return conn.Query(sqlStr, userID, friendID, types.IsFriend)
}

//查看好友详情
func FindUserInfo(friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM `user` WHERE user_id = ?"
	return conn.Query(sqlStr, friendID)
}

//查看好友关系详情 备注等信息
func FindFriend(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT f.user_id as F_user_id,u.user_id as U_user_id,f.*,u.* FROM `friends` f LEFT JOIN user u ON f.friend_id = u.user_id WHERE f.user_id = ? AND f.friend_id = ?"
	return conn.Query(sqlStr, userID, friendID)
}

//TODO
//查询好友id头像备注（昵称）
func FindFriendInfoByUserId(uid string, fid string) ([]map[string]string, error) {
	sql := "SELECT f.friend_id,u.avatar, IF(ISNULL(f.remark)||LENGTH(f.remark)<1,u.username,f.remark) AS username  FROM friends f left join user u on f.friend_id = u.user_id WHERE f.user_id = ? and f.friend_id = ?"
	return conn.Query(sql, uid, fid)
}

//获取最新的消息id
func FindLastCatLogId(userID, friendID string) ([]map[string]string, error) {
	sqlStr := "SELECT MAX(`id`) FROM private_chat_log WHERE (sender_id = ? AND receive_id = ?) OR (sender_id = ? AND receive_id = ?) and is_delete = 2 "
	return conn.Query(sqlStr, userID, friendID, friendID, userID)
}

//查找特定类型消息记录
func FindTypicalChatLogs(userID, friendID, owner string, start int64, number int, queryType []string) ([]map[string]string, int64, error) {
	condition := ""
	if owner != "" {
		condition += " and sender_id = " + owner
	}
	if len(queryType) > 0 {
		condition += fmt.Sprintf(" and msg_type in(%s)", strings.Join(queryType, ","))
	}
	sqlStr := fmt.Sprintf(`SELECT * FROM private_chat_log WHERE ((sender_id = ? AND receive_id = ?) OR (sender_id = ? AND receive_id = ?)) %s AND id >= ? AND id <= ? AND is_delete = 2 AND 
send_time > (SELECT add_time FROM friends WHERE user_id = ? AND friend_id = ? ) ORDER BY id DESC LIMIT ?`, condition)
	var nextId int64 = -1
	maps, err := conn.Query(sqlStr, userID, friendID, friendID, userID, start-300, start, userID, friendID, number+1)
	if err != nil {
		return nil, -1, err
	}
	if len(maps) > number {
		nextId = utility.ToInt64(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	return maps, nextId, nil
}

func FindAllPrivateLogs() ([]map[string]string, error) {
	const sqlStr = `SELECT * FROM private_chat_log`
	return conn.Query(sqlStr)
}

//查找消息记录
func FindCatLog(userID, friendID string, start int64, number int) ([]map[string]string, int64, error) {
	sqlStr := `SELECT * FROM private_chat_log WHERE ((sender_id = ? AND receive_id = ?) OR (sender_id = ? AND receive_id = ?)) AND id <= ? AND is_delete = 2 AND 
send_time > (SELECT add_time FROM friends WHERE user_id = ? AND friend_id = ?) ORDER BY id DESC LIMIT ?`
	var nextId int64 = -1
	maps, err := conn.Query(sqlStr, userID, friendID, friendID, userID, start, userID, friendID, number+1)
	if err != nil {
		return nil, -1, err
	}
	if len(maps) > number {
		nextId = utility.ToInt64(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	return maps, nextId, nil
}

//查找消息记录，不需要判断时间是否是添加好友之后
func FindChatLogV2(userID, friendID string, start int64, number int) ([]map[string]string, int64, error) {
	sqlStr := `SELECT * FROM private_chat_log WHERE ((sender_id = ? AND receive_id = ?) OR (sender_id = ? AND receive_id = ?)) AND id >= ? AND id <= ? AND is_delete = 2 ORDER BY id DESC LIMIT ?`
	var nextId int64 = -1
	maps, err := conn.Query(sqlStr, userID, friendID, friendID, userID, start-300, start, number+1)
	if err != nil {
		return nil, -1, err
	}
	if len(maps) > number {
		nextId = utility.ToInt64(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	return maps, nextId, nil
}

//查询 群会话秘钥 通知消息
func FindSessionKeyAlert(userId string, endTime int64) ([]map[string]string, error) {
	const sqlStr = `SELECT * FROM private_chat_log WHERE 
receive_id = ? and msg_type = ? and send_time <= ? and content LIKE '%\"type\":19%'`
	return conn.Query(sqlStr, userId, types.Alert, endTime)
}

//查找user的名字头像
func SenderInfo(userID string) ([]map[string]string, error) {
	sqlStr := "SELECT username as name,avatar FROM user WHERE user_id = ?"
	return conn.Query(sqlStr, userID)
}

//查询该条好友聊天记录是否属于user
func CheckCatLogIsUser(userId, id string) (bool, error) {
	sqlStr := "SELECT * from private_chat_log WHERE id = ? and sender_id = ?"
	rows, err := conn.Query(sqlStr, id, userId)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

//删除好友聊天记录
func DeleteCatLog(id string) (int, error) {
	sqlStr := "update private_chat_log set is_delete = 1 WHERE id = ?"
	num, _, err := conn.Exec(sqlStr, id)
	return int(num), err
}

//根据uid查找用户id
func FindUserByMarkId(appId, uid string) ([]map[string]string, error) {
	sqlStr := "SELECT user_id FROM user WHERE app_id = ? and (account = ? or uid = ?)"
	return conn.Query(sqlStr, appId, uid, uid)
}

func FindUserByPhone(appId, uid string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM user WHERE app_id = ? and account LIKE '__" + uid + "'"
	return conn.Query(sqlStr, appId)
}

func FindUserByPhoneV2(appId, phone string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM user WHERE app_id = ? and phone = ?"
	return conn.Query(sqlStr, appId, phone)
}

func FindUserByToken(appId, token string) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM user u RIGHT JOIN token t ON u.user_id = t.user_id  WHERE u.app_id = ? AND t.token = ?"
	return conn.Query(sqlStr, appId, token)
}

//获取所有好友未读消息数
func GetAllFriendUnreadMsgCountByUserId(uid string, status int) (int32, error) {
	sqlStr := "SELECT * FROM friends f RIGHT JOIN private_chat_log c ON f.friend_id = c.sender_id  WHERE f.user_id = ? AND c.status = ? and is_delete = 2"
	rows, err := conn.Query(sqlStr, uid, status)
	if err != nil {
		return 0, err
	}
	return int32(len(rows)), nil
}

//查询所有未删除好友id
func FindFriendIdByUserId(uid string) ([]map[string]string, error) {
	sql := "SELECT friend_id  FROM friends  WHERE user_id = ? and is_delete = ?"
	return conn.Query(sql, uid, types.FriendIsNotDelete)
}

//查询未读消息数
func FindUnReadNum(uid, fid string) (int32, error) {
	sql := "SELECT *  FROM private_chat_log WHERE sender_id = ? AND receive_id = ? AND status = ? and is_delete = ?"
	rows, err := conn.Query(sql, fid, uid, types.NotRead, types.FriendMsgIsNotDelete)
	if err != nil {
		return 0, err
	}
	return int32(len(rows)), nil
}

func FindAllUnReadNum(uid string, status int) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE receive_id = ? AND status = ? and is_delete = 2 ORDER BY id ASC"
	return conn.Query(sql, uid, status)
}

func FindAllReaded(uid string, status int, time int64) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE (receive_id = ? or sender_id = ?) AND status = ? and is_delete = 2  and send_time >= ? ORDER BY id ASC"
	return conn.Query(sql, uid, uid, status, time)
}

func FindNotBurndLogAfter(uid string, isDel int, time int64) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE (receive_id = ? or sender_id = ?) AND status != ? and is_delete = ? and send_time > ? ORDER BY id ASC"
	return conn.Query(sql, uid, uid, types.HadBurnt, isDel, time)
}

func FindNotBurndLogBetween(uid string, isDel int, begin, end int64) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE (receive_id = ? or sender_id = ?) AND status != ? and is_delete = ? and send_time >= ? and send_time < ? ORDER BY id DESC"
	return conn.Query(sql, uid, uid, types.HadBurnt, isDel, begin, end)
}

//查找消息记录
func FindChatLogsNumberBetween(uid string, isDel int, begin, end int64) ([]map[string]string, error) {
	sql := "SELECT count(*) as count FROM private_chat_log WHERE (receive_id = ? or sender_id = ?) AND status != ? and is_delete = ? and send_time >= ? and send_time < ?"
	return conn.Query(sql, uid, uid, types.HadBurnt, isDel, begin, end)
}

func FindLastlogByUserId(uid string, status int) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE (receive_id = ? or sender_id = ?) AND status = ? and is_delete = 2 ORDER BY id DESC LIMIT 1"
	return conn.Query(sql, uid, uid, status)
}

//查询好友第一条聊天记录
func FindFirstMsg(uid, fid string) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log c WHERE c.sender_id = ? AND c.status = ? AND c.receive_id = ? and is_delete = 2  ORDER BY id DESC LIMIT 0,1"
	return conn.Query(sql, fid, types.NotRead, uid)
}

//添加私聊聊天记录
func AddPrivateChatLog(senderId, receiveId, msgId string, msgType, status, isSnap int, content, ext string, time int64) (int64, int64, error) {
	sql := "INSERT INTO private_chat_log (sender_id,receive_id,msg_id,msg_type,content,ext,status,send_time,is_delete,is_snap) VALUES (?,?,?,?,?,?,?,?,2,?)"
	return conn.Exec(sql, senderId, receiveId, msgId, msgType, content, ext, status, time, isSnap)
}

//修改聊天记录状态
func ChangePrivateChatLogStstus(id, status int) (int64, int64, error) {
	sql := "UPDATE private_chat_log SET status = ? WHERE id = ? and status != ?"
	return conn.Exec(sql, status, id, types.HadBurnt)
}

//查找聊天记录
func FindPrivateChatLog(senderId, receiveId string) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE sender_id = ? AND receive_id = ?"
	return conn.Query(sql, senderId, receiveId)
}

func FindPrivateChatLogById(logId string) ([]map[string]string, error) {
	sql := "SELECT * FROM `private_chat_log` LEFT JOIN `user` ON private_chat_log.sender_id = user.user_id WHERE `id` = ?"
	return conn.Query(sql, logId)
}

func FindPrivateChatLogByMsgId(senderId, MsgId string) ([]map[string]string, error) {
	sql := "SELECT * FROM `private_chat_log` LEFT JOIN `user` ON private_chat_log.sender_id = user.user_id WHERE `sender_id` = ? and msg_id = ?"
	return conn.Query(sql, senderId, MsgId)
}

func UpdatePrivateLogContentById(logId, content string) error {
	sql := "UPDATE private_chat_log SET content = ? WHERE id = ?"
	_, _, err := conn.Exec(sql, content, logId)
	return err
}

//通过状态查找聊天记录
func FindPrivateChatLogByStatus(senderId, receiveId string, status int) ([]map[string]string, error) {
	sql := "SELECT * FROM private_chat_log WHERE sender_id = ? AND receive_id = ? and status = ?"
	return conn.Query(sql, senderId, receiveId, status)
}

//修改聊天记录状态
func ChangePrivateChatLogStstusByUserAndFriendId(uid, fid string) (int64, int64, error) {
	sql := "UPDATE private_chat_log SET status = ? WHERE sender_id = ? and receive_id = ?"
	return conn.Exec(sql, types.HadRead, fid, uid)
}

//修改聊天记录状态
func UpdatePrivateLogStateById(logId string, state int) (int64, int64, error) {
	sql := "UPDATE private_chat_log SET status = ? WHERE id = ?"
	return conn.Exec(sql, state, logId)
}

//批量修改聊天记录状态
func ChangePrivateChatLogStstusByDatetime(userId string, datetime int64) (int64, int64, error) {
	sql := "UPDATE private_chat_log SET status = ? WHERE receive_id = ? and send_time <= ? and status != ?"
	return conn.Exec(sql, types.HadRead, userId, datetime, types.HadBurnt)
}

//查询群成员id
func FindRoomMemberIds(roomId int) ([]map[string]string, error) {
	sql := `SELECT user_id FROM room_user WHERE room_id = ? AND is_delete = ?`
	return conn.Query(sql, roomId, types.IsNotDelete)
}

//user1 是否对 user2 消息免打扰  是否消息免打扰 1免打扰 2关闭
func CheckDND(user1, user2 string) (bool, error) {
	sql := "select DND from friends where user_id = ? and friend_id = ? and is_delete = 1"
	result, err := conn.Query(sql, user1, user2)
	if err != nil {
		return true, err
	}
	if len(result) == 0 {
		return true, err
	}
	return result[0]["DND"] == "1", err
}

//查找加好友配置
func FindAddFriendConfByUserId(userId string) ([]map[string]string, error) {
	sql := "select * from add_friend_conf where user_id = ?"
	return conn.Query(sql, userId)
}

//修改加好友配置
//是否需要验证
func IsNeedConfirm(userId string, state int) error {
	sql := "INSERT INTO add_friend_conf (user_id,need_confirm,need_answer,question,answer) values (?,?,?,?,?)  ON DUPLICATE KEY UPDATE need_confirm = ?"
	_, _, err := conn.Exec(sql, userId, state, types.NotNeedAnswer, "", "", state)
	return err
}

//需要回答问题
func NeedAnswer(userId string, question, answer string) error {
	sql := "INSERT INTO add_friend_conf (user_id,need_confirm,need_answer,question,answer) values (?,?,?,?,?)  ON DUPLICATE KEY UPDATE need_answer = ?,question = ?,answer = ?"
	_, _, err := conn.Exec(sql, userId, types.NeedConfirm, types.NeedAnswer, question, answer, types.NeedAnswer, question, answer)
	return err
}

//不需要回答问题
func NotNeedAnswer(userId string) error {
	sql := "update add_friend_conf set need_answer = ? where user_id = ?"
	_, _, err := conn.Exec(sql, types.NotNeedAnswer, userId)
	return err
}

//设置问题和答案
func SetQuestionandAnswer(userId, question, answer string) error {
	sql := "INSERT INTO add_friend_conf (user_id,need_confirm,need_answer,question,answer) values (?,?,?,?,?)  ON DUPLICATE KEY UPDATE question = ?, answer = ?"
	_, _, err := conn.Exec(sql, userId, types.NeedConfirm, types.NotNeedAnswer, question, answer, question, answer)
	return err
}

//查询好友来源
func FindFriendSource(userId, friendId string) ([]map[string]string, error) {
	sql := "select source from friends where user_id = ? and friend_id = ?"
	return conn.Query(sql, userId, friendId)
}

//查询最新添加好友记录
func FindApplyOrderByTime(userId, friendId string) ([]map[string]string, error) {
	sql := "select * from apply where (apply_user = ? and target = ? and type = ?) or (target = ? and apply_user = ? and type = ?) order by datetime desc LIMIT 0,1"
	return conn.Query(sql, userId, friendId, types.IsFriend, userId, friendId, types.IsFriend)
}

//查询添加好友记录
func FindApplyByUserId(userId, friendId string) ([]map[string]string, error) {
	sql := "select * from apply where apply_user = ? and target = ? and type = ?"
	return conn.Query(sql, userId, friendId, types.IsFriend)
}

//修改是否加入黑名单标志
func SetFriendIsBlock(userId, friendId string, state int, alterTime int64) error {
	sql := "UPDATE `friends` SET is_blocked = ?,add_time = ? WHERE user_id = ? AND friend_id = ?"
	_, _, err := conn.Exec(sql, state, alterTime, userId, friendId)
	return err
}

//
func FindBlockedList(userId string) ([]map[string]string, error) {
	const sqlStr = "select f.user_id as F_user_id,u.user_id as U_user_id,f.*,u.* from `friends` as f left join `user` as u on f.friend_id=u.user_id where f.user_id = ? and f.is_blocked = ?"
	return conn.Query(sqlStr, userId, types.IsBlocked)
}
