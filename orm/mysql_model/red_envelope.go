package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func InsertPacketInfo(packetID, userID, toID, tType, size, amount, remark string, cType, coin int, time int64) error {
	return db.InsertPacket(packetID, userID, toID, tType, size, amount, remark, cType, coin, time)
}

func GetRedPacketInfo(packetId string) (*types.RedPacketLog, error) {
	maps, err := db.GetRedPacket(packetId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil //errors.New("未找到红包记录")
	}
	info := maps[0]
	return &types.RedPacketLog{
		PacketId:  info["packet_id"],
		CType:     utility.ToInt(info["ctype"]),
		UserId:    info["user_id"],
		ToId:      info["to_id"],
		Coin:      utility.ToInt(info["coin"]),
		Size:      utility.ToInt(info["size"]),
		Amount:    utility.ToFloat64(info["amount"]),
		Remark:    info["remark"],
		Type:      utility.ToInt(info["type"]),
		CreatedAt: utility.ToInt64(info["created_at"]),
	}, nil
}
