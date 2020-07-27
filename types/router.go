package types

const groupPrefix = "Group-"
const roomPrefix = "Room-"

func GetGroupRouteById(str string) string {
	return groupPrefix + str
}

func GetRoomRouteById(str string) string {
	return roomPrefix + str
}
