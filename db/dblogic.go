package db

func InsertPacket(packetID, userID, toID, tType, size, amount, remark string, cType, coin int, time int64) error {
	const sqlStr = "insert into red_packet_log(packet_id,ctype,user_id,to_id,coin,size,amount,remark,type, created_at) values(?,?,?,?,?,?,?,?,?,?)"
	_, _, err := conn.Exec(sqlStr, packetID, cType, userID, toID, coin, size, amount, remark, tType, time)
	return err
}

func GetRedPacket(packetId string) ([]map[string]string, error) {
	return conn.Query(`SELECT r.packet_id, r.user_id, r.type, r.coin, r.size, r.amount, r.remark, r.created_at,r.to_id, r.ctype, u.username, 
		u.avatar, u.uid FROM red_packet_log as r, user as u WHERE r.packet_id=? AND r.user_id=u.user_id`, packetId)
}
