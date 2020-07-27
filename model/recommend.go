package model

import (
	"sync"
	"time"

	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/types"
)

const (
	//120m ping
	timeInterval = 120

	DesignationRoom = 1
	AutoRoom        = 0
)

var recommendRooms = make(map[string]*RecommendRoom)
var rTimer *time.Timer

func init() {
	rTimer = time.NewTimer(timeInterval * time.Minute)
}

type RecommendRoom struct {
	sync.RWMutex

	//指定群
	data []*types.Room
	//根据member排序
	dataMem []*types.Room
	//根据message number排序
	dataMsg []*types.Room
}

func NewRecommendRoom() *RecommendRoom {
	return &RecommendRoom{
		data:    make([]*types.Room, 0),
		dataMem: make([]*types.Room, 0),
		dataMsg: make([]*types.Room, 0),
	}
}

func GetRecommendRoom(appId string) *RecommendRoom {
	return recommendRooms[appId]
}

func (t *RecommendRoom) SetData(data []*types.Room) {
	t.Lock()
	defer t.Unlock()
	t.data = data
}

func (t *RecommendRoom) AddData(item *types.Room) {
	t.Lock()
	defer t.Unlock()
	t.data = append(t.data, item)
}

func (t *RecommendRoom) DelData(item *types.Room) {
	t.Lock()
	defer t.Unlock()
	for i := 0; i < len(t.data); i++ {
		if t.data[i].Id == item.Id {
			t.data = append(t.data[:i], t.data[i+1:]...)
			return
		}
	}
}

func (t *RecommendRoom) SetDataMem(data []*types.Room) {
	t.Lock()
	defer t.Unlock()
	t.dataMem = data
}

func (t *RecommendRoom) SetDataMsg(data []*types.Room) {
	t.Lock()
	defer t.Unlock()
	t.dataMsg = data
}

func (t *RecommendRoom) GetCopyData() []*types.Room {
	t.RLock()
	defer t.RUnlock()
	r := make([]*types.Room, 0)
	r = append(r, t.data...)
	return r
}

func (t *RecommendRoom) GetCopyDataMem() []*types.Room {
	t.RLock()
	defer t.RUnlock()
	r := make([]*types.Room, 0)
	r = append(r, t.dataMem...)
	return r
}

func (t *RecommendRoom) GetCopyDataMsg() []*types.Room {
	t.RLock()
	defer t.RUnlock()
	r := make([]*types.Room, 0)
	r = append(r, t.dataMsg...)
	return r
}

func (t *RecommendRoom) Clear() {
	t.Lock()
	defer t.Unlock()

	t.data = t.data[0:0]
	t.dataMem = t.dataMem[0:0]
	t.dataMsg = t.dataMsg[0:0]
}

func StartServe() {
	go func() {
		defer func() {
			rTimer.Stop()
		}()
		RefreshRooms()
		for {
			<-rTimer.C
			rTimer.Reset(timeInterval * time.Minute)
			RefreshRooms()
		}
	}()
}

func RefreshRooms() {
	queryTime := utility.MillionSecondAddDate(utility.NowMillionSecond(), 0, 0, -7)
	for _, app := range app.GetApps() {
		if r, ok := recommendRooms[app.AppId]; !ok || r == nil {
			recommendRooms[app.AppId] = NewRecommendRoom()
		}
		reRoom := recommendRooms[app.AppId]
		reRoom.Clear()
		//获取指定推荐群
		rs1, err := orm.FindAllRecommendRooms(app.AppId)
		if rs1 != nil && err == nil {
			reRoom.SetData(rs1)
		}

		//获取一周内发言人数排前的群
		rs2, err := orm.RoomsOrderActiveMember(app.AppId, queryTime)
		if rs2 != nil && err == nil {
			reRoom.SetDataMem(rs2)
		}

		//获取一周内发言条数排前的群
		rs3, err := orm.RoomsOrderActiveMsg(app.AppId, queryTime)
		if rs3 != nil && err == nil {
			reRoom.SetDataMsg(rs3)
		}
	}
}
