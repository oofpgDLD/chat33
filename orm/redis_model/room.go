package redis_model

import (
	"github.com/33cn/chat33/cache"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var l = log15.New("module", "redis_model")

// 判断用户是否为管理员或者群主
func GetRoomUserLevel(roomId, userId string, isDel int) (int, error) {
	memberInfo, err := cache.Cache.FindRoomMemberInfo(roomId, userId)
	if err != nil {
		l.Warn("find member info from redis failed", "err", err)
	}
	if memberInfo == nil {
		//从mysql 获取并更新缓存
		member, err := mysql.FindRoomMemberById(roomId, userId, isDel) // db.GetRoomMemberInfo(roomId, userId, isDel)
		if member == nil {
			return types.RoomLevelNotExist, err
		}
		info := &types.RoomMember{
			Id:           member.Id,
			RoomId:       member.RoomId,
			UserId:       member.RoomMember.UserId,
			UserNickname: member.UserNickname,
			Level:        member.Level,
			NoDisturbing: member.NoDisturbing,
			CommonUse:    member.CommonUse,
			RoomTop:      member.RoomTop,
			CreateTime:   member.RoomMember.CreateTime,
			IsDelete:     member.IsDelete,
			Source:       member.Source,
		}
		err = cache.Cache.SaveRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			l.Warn("redis can not save room member")
		}
		return member.Level, nil
	}
	return memberInfo.Level, nil
}

//增加群成员,同时入群申请的第一步 添加user也用这个方法
func AddMember(tx types.Tx, userId, roomId string, memberLevel int, id, createTime int64, source string) error {
	info := &types.RoomMember{
		Id:           utility.ToString(id),
		RoomId:       roomId,
		UserId:       userId,
		UserNickname: "",
		Level:        memberLevel,
		NoDisturbing: types.NoDisturbingOff,
		CommonUse:    types.UncommonUse,
		RoomTop:      types.NotOnTop,
		CreateTime:   createTime,
		IsDelete:     types.RoomUserNotDeleted,
		Source:       source,
	}
	//存储user信息
	err := cache.Cache.SaveRoomUser(roomId, &[]*types.RoomMember{info})
	if err != nil {
		l.Warn("redis can not save room member")
		return err
	}
	//room-user关系也存储到redis
	return cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{info})
}

func GetRoomMemberName(roomId, userId string) (string, error) {
	info, err := FindRoomMemberById(roomId, userId, types.RoomUserNotDeleted)
	if info == nil {
		return "", err
	}
	if info.RoomMember.UserNickname == "" || len(info.RoomMember.UserNickname) < 1 {
		return info.User.Username, nil
	}
	return info.RoomMember.UserNickname, nil
}

//set三种permission,写在一起
func SetPermission(roomId string, canAddFriend, joinPermission, recordPermission, isDel int) error {
	//设置更新操作，不管redis存不存在，都要取mysql
	info, err := mysql.SearchRoomInfo(roomId, isDel)
	if err != nil {
		return err
	}
	roominfo, err := cache.Cache.FindRoomInfo(roomId)
	if err != nil {
		l.Warn("find room info from redis failed", "err", err)
	}
	//redis里不存在,直接存，存在，更新
	if roominfo == nil {
		err = cache.Cache.SaveRoomInfo(info)
		if err != nil {
			l.Warn("redis can not save room Info", "err", err)
			return err
		}
	}
	if canAddFriend != 0 {
		err = cache.Cache.UpdateRoomInfo(roomId, "canAddFriend", utility.ToString(canAddFriend))
		if err != nil {
			l.Warn("redis can not update room info", "err", err)
			return err
		}
	}
	if joinPermission != 0 {
		err = cache.Cache.UpdateRoomInfo(roomId, "joinPermission", utility.ToString(joinPermission))
		if err != nil {
			l.Warn("redis can not update room info", "err", err)
			return err
		}
	}
	if recordPermission != 0 {
		err = cache.Cache.UpdateRoomInfo(roomId, "recordPermision", utility.ToString(recordPermission))
		if err != nil {
			l.Warn("redis can not update room info", "err", err)
			return err
		}
	}
	return err

}

//设置群名
func SetRoomName(roomId, name string, isDel int) (bool, error) {
	var b = true
	info, err := mysql.SearchRoomInfo(roomId, isDel)
	//判断redis有没有群消息
	roominfo, err := cache.Cache.FindRoomInfo(roomId)
	if err != nil {
		l.Warn("find room info from redis failed", "err", err)
		b = false
	}
	//redis里不存在,直接存
	if roominfo == nil {
		err = cache.Cache.SaveRoomInfo(info)
		if err != nil {
			l.Warn("redis can not save room Info", "err", err)
			b = false
		}
		return b, err
	}
	err = cache.Cache.UpdateRoomInfo(roomId, "name", name)
	if err != nil {
		l.Warn("redis can not update room Info", "err", err)
		b = false
	}
	return b, err
}

func SetAvatar(roomId, avatar string, isDel int) (bool, error) {
	var b = true
	info, err := mysql.SearchRoomInfo(roomId, isDel)
	//判断redis有没有群消息
	roominfo, err := cache.Cache.FindRoomInfo(roomId)
	if err != nil {
		l.Warn("find room info from redis failed", "err", err)
		b = false
	}
	//redis里不存在,直接存
	if roominfo == nil {
		err = cache.Cache.SaveRoomInfo(info)
		if err != nil {
			l.Warn("redis can not save room Info", "err", err)
			b = false
		}
		return b, err
	}
	err = cache.Cache.UpdateRoomInfo(roomId, "avatar", avatar)
	if err != nil {
		l.Warn("redis can not update room Info", "err", err)
		b = false
	}
	return b, err
}

//设置成员等级
func SetMemberLevel(userId, roomId string, level, isDel int) error {
	memberInfo, err := cache.Cache.FindRoomMemberInfo(roomId, userId)
	if err != nil {
		l.Warn("find member info from redis failed", "err", err)
	}
	//redis里有数据，更新，无数据，直接存
	if memberInfo == nil {
		//先取mysql数据
		member, err := mysql.FindRoomMemberById(roomId, userId, isDel) // db.GetRoomMemberInfo(roomId, userId, isDel)
		if member == nil {
			l.Warn("Room Member info can not find")
			return err
		}
		info := &types.RoomMember{
			Id:           member.Id,
			RoomId:       member.RoomId,
			UserId:       member.RoomMember.UserId,
			UserNickname: member.UserNickname,
			Level:        member.Level,
			NoDisturbing: member.NoDisturbing,
			CommonUse:    member.CommonUse,
			RoomTop:      member.RoomTop,
			CreateTime:   member.RoomMember.CreateTime,
			IsDelete:     member.IsDelete,
			Source:       member.Source,
		}
		err = cache.Cache.SaveRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			l.Warn("redis can not save room member")
			return err
		}
		//room-user关系也存储到redis
		err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{info})
	}
	err = cache.Cache.UpdateRoomUserInfo(roomId, userId, "level", utility.ToString(level))
	if err != nil {
		l.Warn("redis can not update room member")
		return err
	}
	//UpdateRoomUser用来更新关系表里面的level
	return cache.Cache.UpdateRoomUser(roomId, userId, utility.ToString(level))
}

//转让群主  master原群主id userid新群主id
func SetNewMaster(master, userId, roomId string, level int) error {
	//设置群的master_id
	err := cache.Cache.UpdateRoomInfo(roomId, "masterId", userId)
	if err != nil {
		l.Warn("Update masterId failed", "err", err)
		return err
	}
	//设置新群主在群里的level
	err = SetMemberLevel(userId, roomId, level, types.RoomUserDeletedOrNot)
	if err != nil {
		l.Warn("Set new master Level failed", "err", err)
		return err
	}
	//设置老群主在群里的level
	err = SetMemberLevel(master, roomId, types.RoomLevelNomal, types.RoomUserDeletedOrNot)
	if err != nil {
		l.Warn("Set old master Level failed", "err", err)
		return err
	}
	return nil
}

func SetNoDisturbing(userId, roomId string, permission, isDel int) error {
	//先取mysql数据
	member, err := mysql.FindRoomMemberById(roomId, userId, isDel) // db.GetRoomMemberInfo(roomId, userId, isDel)
	if member == nil {
		l.Warn("Room Member info can not find")
		return err
	}
	memberInfo, err := cache.Cache.FindRoomMemberInfo(roomId, userId)
	if err != nil {
		l.Warn("find member info from redis failed", "err", err)
	}
	//redis里有数据，更新，无数据，直接存
	if memberInfo == nil {
		info := &types.RoomMember{
			Id:           member.Id,
			RoomId:       member.RoomId,
			UserId:       member.RoomMember.UserId,
			UserNickname: member.UserNickname,
			Level:        member.Level,
			NoDisturbing: member.NoDisturbing,
			CommonUse:    member.CommonUse,
			RoomTop:      member.RoomTop,
			CreateTime:   member.RoomMember.CreateTime,
			IsDelete:     member.IsDelete,
			Source:       member.Source,
		}
		err = cache.Cache.SaveRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			l.Warn("redis can not save room member")
			return err
		}
		//room-user关系也存储到redis
		err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			return err
		}
	}
	err = cache.Cache.UpdateRoomUserInfo(roomId, userId, "noDisturbing", utility.ToString(permission))
	if err != nil {
		l.Warn("redis can not update room member")
		return err
	}
	return nil
}

func SetOnTop(userId, roomId string, permission, isDel int) error {
	//先取mysql数据
	member, err := mysql.FindRoomMemberById(roomId, userId, isDel) // db.GetRoomMemberInfo(roomId, userId, isDel)
	if member == nil {
		l.Warn("Room Member info can not find")
		return err
	}
	memberInfo, err := cache.Cache.FindRoomMemberInfo(roomId, userId)
	if err != nil {
		l.Warn("find member info from redis failed", "err", err)
	}
	//redis里有数据，更新，无数据，直接存
	if memberInfo == nil {
		info := &types.RoomMember{
			Id:           member.Id,
			RoomId:       member.RoomId,
			UserId:       member.RoomMember.UserId,
			UserNickname: member.UserNickname,
			Level:        member.Level,
			NoDisturbing: member.NoDisturbing,
			CommonUse:    member.CommonUse,
			RoomTop:      member.RoomTop,
			CreateTime:   member.RoomMember.CreateTime,
			IsDelete:     member.IsDelete,
			Source:       member.Source,
		}
		err = cache.Cache.SaveRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			l.Warn("redis can not save room member")
			return err
		}
		//room-user关系也存储到redis
		err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			return err
		}
	}
	err = cache.Cache.UpdateRoomUserInfo(roomId, userId, "roomTop", utility.ToString(permission))
	if err != nil {
		l.Warn("redis can not update room member")
		return err
	}
	return nil
}

func SetMemberNickname(userId, roomId string, nickname string, isDel int) error {
	//先取mysql数据
	member, err := mysql.FindRoomMemberById(roomId, userId, isDel) // db.GetRoomMemberInfo(roomId, userId, isDel)
	if member == nil {
		l.Warn("Room Member info can not find")
		return err
	}
	memberInfo, err := cache.Cache.FindRoomMemberInfo(roomId, userId)
	if err != nil {
		l.Warn("find member info from redis failed", "err", err)
	}
	//redis里有数据，先清除再存，无数据，直接存
	if memberInfo == nil {
		info := &types.RoomMember{
			Id:           member.Id,
			RoomId:       member.RoomId,
			UserId:       member.RoomMember.UserId,
			UserNickname: member.UserNickname,
			Level:        member.Level,
			NoDisturbing: member.NoDisturbing,
			CommonUse:    member.CommonUse,
			RoomTop:      member.RoomTop,
			CreateTime:   member.RoomMember.CreateTime,
			IsDelete:     member.IsDelete,
			Source:       member.Source,
		}
		err = cache.Cache.SaveRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			l.Warn("redis can not save room member")
			return err
		}
		//room-user关系也存储到redis
		err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{info})
		if err != nil {
			return err
		}
	}
	err = cache.Cache.UpdateRoomUserInfo(roomId, userId, "userNickname", utility.ToString(nickname))
	if err != nil {
		l.Warn("redis can not update room member")
		return err
	}
	return nil
}

//新增成员禁言
func AddMutedMember(id int64, roomId, userId string, mutedType int, deadline int64) error {
	info := &types.RoomUserMuted{
		Id:       utility.ToString(id),
		RoomId:   roomId,
		UserId:   userId,
		ListType: mutedType,
		Deadline: deadline,
	}
	err := cache.Cache.SaveRoomUserMued(roomId, &[]*types.RoomUserMuted{info})
	if err != nil {
		l.Warn("redis can not save Muted member")
		return err
	}
	return nil
}

//取消某个成员禁言，其实是更新那个成员禁言状态
func DelMemberMuted(roomId, userId string) error {
	//取得已经更新的禁言状态
	mutedType, deadline, _ := mysql.GetRoomUserMuted(roomId, userId)

	member, err := mysql.FindRoomMemberById(roomId, userId, types.RoomUserNotDeleted)
	//如果群成员等级大于1，不需要再修改，返回
	if member.Level > 1 {
		l.Warn("The Member is master or admin, can not be Muted")
		return nil
	}
	mudtedInfo, err := cache.Cache.GetRoomUserMuted(roomId, userId)
	if err != nil {
		l.Warn("find Muted info from redis failed", "err", err)
	}
	//先判断禁言信息是否存在，不存在先新增，存在则更新
	if mudtedInfo == nil {
		var id string
		context, err := mysql.GetMutedListByType(roomId, mutedType)
		if err != nil {
			l.Warn("Get Muted List from mysql failed", "err", err)
		}
		//取出Id
		for _, v := range context {
			if v.UserId == userId {
				id = v.Id
			}
		}
		return AddMutedMember(utility.ToInt64(id), roomId, userId, mutedType, deadline)
	}
	info := &types.RoomUserMuted{
		Id:       mudtedInfo.Id,
		RoomId:   roomId,
		UserId:   userId,
		ListType: mutedType,
		Deadline: deadline,
	}

	err = cache.Cache.UpdateRoomUserMuted(roomId, info)
	if err != nil {
		l.Warn("redis can not delete Muted member")
		return err
	}
	return nil
}

func GetRoomMutedType(roomId string, isDel int) (int, error) {
	roominfo, err := SearchRoomInfo(roomId, isDel)
	if err != nil {
		return 0, err
	}
	if roominfo == nil {
		return 0, nil
	}
	return roominfo.MasterMuted, nil
}

func SetRoomMutedType(roomId string, mutedType, isDel int) error {
	_, err := SearchRoomInfo(roomId, isDel)
	if err != nil {
		return err
	}
	roomInfo, err := cache.Cache.FindRoomInfo(roomId)
	if err != nil {
		l.Warn("find room info from redis failed", "err", err)
	}
	//如果是已经删除状态  不需要更新
	if roomInfo != nil {
		err = cache.Cache.UpdateRoomInfo(roomId, "masterMuted", utility.ToString(mutedType))
		if err != nil {
			l.Warn("redis can not update room Info")
			return err
		}
	}

	return nil
}

func GetRoomUserMuted(roomId, userId string) (mutedType int, deadline int64) {
	info, err := cache.Cache.GetRoomUserMuted(roomId, userId)
	if err != nil {
		l.Warn("find Room User Muted from redis failed", "err", err)
	}
	if info == nil {
		var id string

		mutedType, deadline, _ := mysql.GetRoomUserMuted(roomId, userId)
		context, err := mysql.GetMutedListByType(roomId, mutedType)
		if err != nil {
			l.Warn("Get Muted List from mysql failed", "err", err)
		}
		//取出Id
		for _, v := range context {
			if v.UserId == userId {
				id = v.Id
			}
		}
		err = AddMutedMember(utility.ToInt64(id), roomId, userId, mutedType, deadline)
		if err != nil {
			l.Warn("AddMutedMember to redis failed")
		}
		return mutedType, deadline
	}
	return info.ListType, info.Deadline
}

func ClearMutedList(tx types.Tx, roomId string) error {
	roomUserMued := &types.RoomUserMuted{
		ListType: types.AllSpeak,
		Deadline: 0,
	}
	info, err := cache.Cache.GetRoomUser(roomId)
	if err != nil {
		l.Warn("find Room_User from redis failed", "err", err)
	}
	//如果关系不存在，要先存关系，然后直接操作数据库
	if info == nil {
		//通过roomid找到user和room_user表的所有数据 找出有关的数据
		RoomMemberInfo, err := mysql.FindRoomMembers(roomId, types.SearchAll)
		if err != nil {
			l.Warn("Room Member info can not find", err)
		}
		//MemberInfo:=&types.RoomMember{}
		for _, v := range RoomMemberInfo {
			roomUserMued.UserId = v.RoomMember.UserId
			err := cache.Cache.UpdateRoomUserMuted(roomId, roomUserMued)
			if err != nil {
				l.Warn("redis can not Update Muted member")
				return err
			}
			//循环存储到redis
			err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{v.RoomMember})
			if err != nil {
				return err
			}
		}
	}
	//key是userid ,val是level(这里用不到)
	for userid := range info {
		roomUserMued.UserId = userid
		err := cache.Cache.UpdateRoomUserMuted(roomId, roomUserMued)
		if err != nil {
			l.Warn("redis can not Update Muted member")
			return err
		}
	}
	return nil
}

//获取某种禁言信息
func GetMutedListByType(roomId string, mutedType int) ([]*types.RoomUserMuted, error) {
	res := make([]*types.RoomUserMuted, 0)

	roomUserMuted := &types.RoomUserMuted{}

	info, err := cache.Cache.GetRoomUser(roomId)
	if err != nil {
		l.Warn("find Room_User from redis failed", "err", err)
	}
	if info == nil {
		//通过roomid找到user和room_user表的所有数据 找出没有被删除的user数据
		RoomMemberInfo, err := mysql.FindRoomMembers(roomId, types.SearchAll)
		if err != nil {
			l.Warn("Room Member info can not find", err)
		}
		//MemberInfo:=&types.RoomMember{}
		for _, v := range RoomMemberInfo {
			//循环存储到redis
			err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{v.RoomMember})
			if err != nil {
				return nil, err
			}
		}
		return mysql.GetMutedListByType(roomId, mutedType)
	}
	//key是userid ,val是level
	for userid, level := range info {
		//不是群主或者管理员
		if level == types.RoomLevelManager || level == types.RoomLevelMaster {
			continue
		}
		Mutedinfo, _ := cache.Cache.GetRoomUserMuted(roomId, userid)
		//先判断禁言信息是否存在，不存在先新增，存在则更新
		if Mutedinfo == nil {
			var id string
			mutedType, deadline, _ := mysql.GetRoomUserMuted(roomId, userid)
			context, err := mysql.GetMutedListByType(roomId, mutedType)
			if err != nil {
				l.Warn("Get Muted List from mysql failed", "err", err)
			}
			//取出Id
			for _, v := range context {
				if v.UserId == userid {
					id = v.Id
				}
			}
			err = AddMutedMember(utility.ToInt64(id), roomId, userid, mutedType, deadline)
			if err != nil {
				return nil, err
			}
			roomUserMuted = &types.RoomUserMuted{
				Id:       id,
				RoomId:   roomId,
				UserId:   userid,
				ListType: mutedType,
				Deadline: deadline,
			}
			res = append(res, roomUserMuted)
		} else {
			if Mutedinfo.ListType == mutedType {

				res = append(res, Mutedinfo)
			}
		}
	}
	return res, err
}

//查找recv群消息
func FindReceiveLogById(logId, userId string) (*types.RoomMsgReceive, error) {
	res := &types.RoomMsgReceive{}
	receiveInfo, err := cache.Cache.GetReceiveLogAll()
	if err != nil {
		l.Warn("find ReceiveLog from redis failed", "err", err)
	}
	if receiveInfo == nil {
		return mysql.FindReceiveLogById(logId, userId)
	}
	for _, v := range receiveInfo {
		if v.RoomMsgId == logId && v.ReceiveId == userId {
			res.Id = v.Id
			res.ReceiveId = userId
			res.RoomMsgId = logId
			res.State = v.State
		}
	}

	return res, nil
}

//添加群聊接收日志
func AppendMemberRevLog(id int64, userId, logId string, state int) error {
	info := &types.RoomMsgReceive{
		Id:        utility.ToString(id),
		RoomMsgId: logId,
		ReceiveId: userId,
		State:     state,
	}
	err := cache.Cache.SaveReceiveLog(info)
	return err
}

//获取群中所有管理员
func GetRoomManagerAndMaster(roomId string) ([]*types.MemberJoinUser, error) {
	res := make([]*types.MemberJoinUser, 0)
	info, err := cache.Cache.GetRoomUser(roomId)
	if err != nil {
		l.Warn("find Room_User from redis failed", "err", err)
	}
	if info == nil {
		//通过roomid找到user和room_user表的所有数据 找出有关的数据
		RoomMemberInfo, err := mysql.FindRoomMembers(roomId, types.SearchAll)
		if err != nil {
			l.Warn("Room Member info can not find", err)
		}
		//MemberInfo:=&types.RoomMember{}
		for _, v := range RoomMemberInfo {
			//循环存储到redis
			err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{v.RoomMember})
			if err != nil {
				return nil, err
			}
		}
		return mysql.GetRoomManagerAndMaster(roomId)
	}
	//key是userid ,val是level
	for key, val := range info {
		//val>1是群主或者管理员
		if val > 1 {
			MemberJoinUser, err := FindRoomMemberById(roomId, key, types.RoomUserNotDeleted)
			if err != nil {
				l.Warn("redis can not Update Muted member", "err", err)
			}
			res = append(res, MemberJoinUser)
		}
	}
	return res, err
}

//根据id查询群成员信息
func FindRoomMemberById(roomId, userId string, isDel int) (*types.MemberJoinUser, error) {
	memberInfo, err := cache.Cache.FindRoomMemberInfo(roomId, userId)
	if err != nil {
		l.Warn("find member info from redis failed", "err", err)
	}

	if memberInfo == nil {
		//从mysql 获取并更新缓存
		member, err := mysql.FindRoomMemberById(roomId, userId, isDel) // db.GetRoomMemberInfo(roomId, userId, isDel)
		if member == nil {
			return nil, err
		}
		m := &types.RoomMember{
			Id:           member.Id,
			RoomId:       member.RoomId,
			UserId:       member.RoomMember.UserId,
			UserNickname: member.UserNickname,
			Level:        member.Level,
			NoDisturbing: member.NoDisturbing,
			CommonUse:    member.CommonUse,
			RoomTop:      member.RoomTop,
			CreateTime:   member.RoomMember.CreateTime,
			IsDelete:     member.IsDelete,
			Source:       member.Source,
		}
		u := &types.User{
			UserId:      member.RoomMember.UserId,
			MarkId:      member.MarkId,
			Uid:         member.Uid,
			AppId:       member.AppId,
			Username:    member.Username,
			Account:     member.Account,
			UserLevel:   member.UserLevel,
			Verified:    member.Verified, //是否实名制
			Avatar:      member.Avatar,
			Phone:       member.Phone,
			Email:       member.Email,
			InviteCode:  member.InviteCode,
			DeviceToken: member.DeviceToken,
		}
		err = cache.Cache.SaveRoomUser(roomId, &[]*types.RoomMember{m})
		if err != nil {
			l.Warn("redis can not save room member")
		}
		return &types.MemberJoinUser{
			RoomMember: m,
			User:       u,
		}, nil
	}
	userInfo, err := GetUserInfoById(userId)
	if err != nil {
		return nil, err
	}
	ret := &types.MemberJoinUser{
		RoomMember: memberInfo,
		User:       userInfo,
	}
	return ret, nil
}

//查找群中消息免打扰的成员
func FindSetNoDisturbingMembers(roomId string) ([]*types.RoomMember, error) {
	res := make([]*types.RoomMember, 0)
	//先获取关系 room 对应的userid
	info, err := cache.Cache.GetRoomUser(roomId)
	if err != nil {
		l.Warn("find Room_User from redis failed", "err", err)
	}
	if info == nil {
		//通过roomid找到user和room_user表的所有数据 找出有关的数据
		RoomMemberInfo, err := mysql.FindRoomMembers(roomId, types.SearchAll)
		if err != nil {
			l.Warn("Room Member info can not find", err)
		}
		//MemberInfo:=&types.RoomMember{}
		for _, v := range RoomMemberInfo {
			//循环存储到redis
			err = cache.Cache.AddRoomUser(roomId, &[]*types.RoomMember{v.RoomMember})
			if err != nil {
				return nil, err
			}
		}
		return mysql.FindSetNoDisturbingMembers(roomId)
	}
	//key是userid ,val是level
	for key := range info {
		MemberJoinUser, err := FindRoomMemberById(roomId, key, types.RoomUserNotDeleted)
		if err != nil {
			l.Warn("redis can not Update Muted member", err)
		}
		if MemberJoinUser.NoDisturbing == types.NoDisturbingOn {
			m := &types.RoomMember{
				Id:           MemberJoinUser.Id,
				RoomId:       MemberJoinUser.RoomId,
				UserId:       MemberJoinUser.RoomMember.UserId,
				UserNickname: MemberJoinUser.UserNickname,
				Level:        MemberJoinUser.Level,
				NoDisturbing: MemberJoinUser.NoDisturbing,
				CommonUse:    MemberJoinUser.CommonUse,
				RoomTop:      MemberJoinUser.RoomTop,
				CreateTime:   MemberJoinUser.RoomMember.CreateTime,
				IsDelete:     MemberJoinUser.IsDelete,
				Source:       MemberJoinUser.Source,
			}
			res = append(res, m)
		}
	}
	return res, err
}

func CreateNewRoom(tx types.Tx, roomId int64, creater, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, randomRoomId string, createTime int64) error {
	info := &types.Room{
		Id:              utility.ToString(roomId),
		MarkId:          randomRoomId,
		Name:            roomName,
		Avatar:          roomAvatar,
		MasterId:        creater,
		CreateTime:      createTime,
		CanAddFriend:    canAddFriend,
		JoinPermission:  joinPermission,
		RecordPermision: recordPermission,
		AdminMuted:      adminMuted,
		MasterMuted:     masterMuted,
		Encrypt:         encrypt,
		IsDelete:        types.RoomNotDeleted,
	}
	err := cache.Cache.SaveRoomInfo(info)
	if err != nil {
		l.Warn("redis can not save room")
	}
	return err
}

//根据room id获取群信息
func SearchRoomInfo(roomid string, isDel int) (*types.Room, error) {
	roomInfo, err := cache.Cache.FindRoomInfo(roomid)
	if err != nil {
		l.Warn("find room info from redis failed", "err", err)
	}
	if roomInfo == nil {
		//从mysql 获取并更新缓存
		room, err := mysql.SearchRoomInfo(roomid, isDel)
		if room == nil {
			l.Warn("Room info can not find")
			return nil, err
		}
		if room.IsDelete == types.RoomNotDeleted {
			err = cache.Cache.SaveRoomInfo(room)
			if err != nil {
				l.Warn("redis can not save room member")
			}
		}
		return room, nil
	}
	if isDel != types.RoomDeletedOrNot && roomInfo.IsDelete != isDel {
		return nil, err
	}
	return roomInfo, nil
}

//删除群
func DeleteRoomInfoById(roomId string) error {
	info, err := cache.Cache.GetRoomUser(roomId)
	if err != nil {
		l.Warn("find Room_User from redis failed", "err", err)
	}
	for userId := range info {
		//val>1是群主或者管理员
		err = cache.Cache.DeleteRoomUserInfo(roomId, userId)
		if err != nil {
			l.Warn("redis can not Delete RoomUserInfo", "err", err)
		}
	}
	err = cache.Cache.DeleteRoomInfo(roomId)
	if err != nil {
		l.Warn("Room info can not delete", "err", err)
	}
	err = cache.Cache.DelRoomUserAll(roomId)
	if err != nil {
		return err
	}
	return cache.Cache.DeleteRoomUserMuted(roomId)
}

//删除群成员
func DeleteRoomMemberById(userId, roomId string) (bool, error) {
	b := true
	//删除群成员消息
	err := cache.Cache.DeleteRoomUserInfo(roomId, userId)
	if err != nil {
		l.Warn("RoomUserInfo can not delete", "err", err)
		b = false
		return b, err
	}
	//删除关系
	err = cache.Cache.DelRoomUser(roomId, userId)
	if err != nil {
		l.Warn("Room_User can not delete", "err", err)
		b = false
		return b, err
	}
	return b, nil
}

//根据logid查群聊天信息
func FindRoomChatLogByContentId(logId string) (*types.RoomLogJoinUser, error) {
	info, err := cache.Cache.GetRoomMsgContent(logId)
	if err != nil {
		l.Warn("find Room_Msg_Content from redis failed", "err", err)
	}
	if info == nil {
		contxt, err := mysql.FindRoomChatLogByContentId(logId)
		roomlog := &types.RoomLog{
			Id:       contxt.Id,
			MsgId:    contxt.MsgId,
			RoomId:   contxt.RoomId,
			SenderId: contxt.SenderId,
			IsSnap:   contxt.IsSnap,
			MsgType:  contxt.MsgType,
			Content:  contxt.Content,
			Datetime: contxt.Datetime,
			Ext:      contxt.Ext,
			IsDelete: contxt.IsDelete,
		}
		err = cache.Cache.SaveRoomMsgContent(roomlog)
		if err != nil {
			return nil, err
		}
		return contxt, err
	}
	userInfo, err := GetUserInfoById(info.SenderId)
	if err != nil {
		l.Warn("user info can not find", "err", err)
	}

	return &types.RoomLogJoinUser{
		RoomLog: info,
		User:    userInfo,
	}, nil
}

//根据senderid和msgid查群聊天信息
func FindRoomChatLogByMsgId(senderId, msgId string) (*types.RoomLogJoinUser, error) {
	res := &types.RoomLogJoinUser{}
	info, err := cache.Cache.GetRoomMsgContentAll()
	if err != nil {
		l.Warn("user info can not find", "err", err)
	}
	userInfo, err := cache.Cache.GetUserInfoById(senderId)

	if err != nil {
		l.Warn("find Room_Msg_Content from redis failed", "err", err)
	}
	if info == nil {
		return mysql.FindRoomChatLogByMsgId(senderId, msgId)
	}
	//v为RoomLog
	for _, v := range info {
		content, _ := FindRoomChatLogByContentId(v.Id)
		if content.SenderId == senderId && content.MsgId == msgId {
			roomlog := &types.RoomLog{
				Id:       content.Id,
				MsgId:    content.MsgId,
				RoomId:   content.RoomId,
				SenderId: content.SenderId,
				IsSnap:   content.IsSnap,
				MsgType:  content.MsgType,
				Content:  content.Content,
				Datetime: content.Datetime,
				Ext:      content.Ext,
				IsDelete: content.IsDelete,
			}
			res = &types.RoomLogJoinUser{
				RoomLog: roomlog,
				User:    userInfo,
			}
		}
	}
	return res, nil
}

//根据id删除群聊天信息
func DelRoomChatLogById(logId string) error {
	info, err := cache.Cache.GetRoomMsgContent(logId)
	if err != nil {
		l.Warn("find Room_Msg_Content from redis failed", "err", err)
	}
	if info == nil {
		return nil
	}
	return cache.Cache.DeleteRoomMsgContent(logId)

}

func AlertRoomRevStateByRevId(Id string, state int) error {

	return cache.Cache.UpdateReceiveLog(Id, state)
}

// 添加群聊聊天日志
func AppendRoomChatLog(logId int64, userId, roomId, msgId string, msgType, isSnap int, content, ext string, time int64) error {

	info := &types.RoomLog{
		Id:       utility.ToString(logId),
		MsgId:    msgId,
		RoomId:   roomId,
		SenderId: userId,
		IsSnap:   isSnap,
		MsgType:  msgType,
		Content:  content,
		Datetime: time,
		Ext:      ext,
		IsDelete: types.RoomMsgNotDelete,
	}
	return cache.Cache.SaveRoomMsgContent(info)
}

//获取用户创建群个数上限
func GetCreateRoomsLimit(appId string, level int) (int, error) {
	info, err := cache.Cache.FindUserConfig(appId)
	if err != nil {
		l.Warn("Find UserConfig from redis failed", "err", err)
	}
	if info == nil {
		lm, err := mysql.GetCreateRoomsLimit(appId, level)
		if err != nil {
			l.Warn("Find UserConfig from mysql failed", "err", err)
		}
		err = cache.Cache.SaveUserConfig(appId, level, lm)
		if err != nil {
			l.Warn("SaveUserConfig to redis failed", "err", err)
		}
		return lm, err
	}
	return info.Limit, nil
}

//获取群的成员数上限
func GetRoomMembersLimit(appId string, level int) (int, error) {
	info, err := cache.Cache.FindRoomConfig(appId)
	if err != nil {
		l.Warn("Find RoomConfig from redis failed", "err", err)
	}
	if info == nil {
		lm, err := mysql.GetRoomMembersLimit(appId, level)
		if err != nil {
			l.Warn("Find RoomConfig from mysql failed", "err", err)
		}
		err = cache.Cache.SaveRoomConfig(appId, level, lm)
		if err != nil {
			l.Warn("SaveRoomConfig to redis failed", "err", err)
		}
		return lm, err
	}
	return info.Limit, nil
}

//设置用户创建群个数上限
func SetCreateRoomsLimit(appId string, level, limit int) error {
	info, err := cache.Cache.FindUserConfig(appId)
	if err != nil {
		l.Warn("Find UserConfig from redis failed", "err", err)
	}
	if info == nil {
		lm, err := mysql.GetCreateRoomsLimit(appId, level)
		if err != nil {
			l.Warn("Find UserConfig from mysql failed", "err", err)
		}
		return cache.Cache.SaveRoomConfig(appId, level, lm)
	}
	return cache.Cache.SaveRoomConfig(appId, level, limit)
}

//设置群的成员数上限
func SetRoomMembersLimit(appId string, level, limit int) error {
	info, err := cache.Cache.FindRoomConfig(appId)
	if err != nil {
		l.Warn("Find RoomConfig from redis failed", "err", err)
	}
	if info == nil {
		lm, err := mysql.GetRoomMembersLimit(appId, level)
		if err != nil {
			l.Warn("Find RoomConfig from mysql failed", "err", err)
		}
		return cache.Cache.SaveUserConfig(appId, level, lm)
	}
	return cache.Cache.SaveUserConfig(appId, level, limit)
}

//设置为加v认证群
func SetRoomVerifyed(roomId, vInfo string) error {
	//设置更新操作，不管redis存不存在，都要取mysql
	info, err := mysql.SearchRoomInfo(roomId, types.RoomNotDeleted)
	roominfo, err := cache.Cache.FindRoomInfo(roomId)
	if err != nil {
		l.Warn("find room info from redis failed", "err", err)
	}
	//redis里不存在,直接存，存在，更新
	if roominfo == nil {
		err = cache.Cache.SaveRoomInfo(info)
		if err != nil {
			l.Warn("redis can not save room Info")
			return err
		}
	}
	err = cache.Cache.UpdateRoomInfo(roomId, "identification", utility.ToString(types.Verified))
	err = cache.Cache.UpdateRoomInfo(roomId, "identificationInfo", vInfo)
	if err != nil {
		l.Warn("redis can not update room info")
		return err
	}
	return nil
}
