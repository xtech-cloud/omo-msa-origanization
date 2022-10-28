package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"omo.msa.organization/proxy/nosql"
	"omo.msa.organization/tool"
)

//房间大厅
type RoomInfo struct {
	baseInfo
	Remark string
	Scene  string
	Quotes []string
}

func (mine *cacheContext) GetRoom(uid string) *RoomInfo {
	for _, scene := range mine.scenes {
		info := scene.GetRoom(uid)
		if info != nil {
			return info
		}
	}
	return nil
}

func (mine *cacheContext) GetRoomsByDevice(sn string) []*RoomInfo {
	list := make([]*RoomInfo, 0, 10)
	for _, scene := range mine.scenes {
		arr := scene.GetRoomsByDevice(sn)
		if arr != nil && len(arr) > 0 {
			list = append(list, arr...)
		}
	}
	return list
}

func (mine *cacheContext) GetRoomsByQuote(quote string) []*RoomInfo {
	list := make([]*RoomInfo, 0, 10)
	for _, scene := range mine.scenes {
		arr := scene.GetRoomsByQuote(quote)
		if arr != nil && len(arr) > 0 {
			list = append(list, arr...)
		}
	}
	return list
}

func (mine *cacheContext) HadBindDeviceInRoom(sn string) bool {
	for _, scene := range mine.scenes {
		arr := scene.GetRoomsByDevice(sn)
		if arr != nil && len(arr) > 0 {
			return true
		}
	}
	return false
}

func (mine *cacheContext) GetRoomBy(scene, uid string) *RoomInfo {
	for _, item := range mine.scenes {
		if item.UID == scene {
			return item.GetRoom(uid)
		}
	}
	return nil
}

func (mine *cacheContext) RemoveRoom(uid, operator string) error {
	for _, scene := range mine.scenes {
		if scene.HadRoom(uid) {
			return scene.RemoveRoom(uid, operator)
		}
	}
	return nil
}

func (mine *RoomInfo) initInfo(db *nosql.Room) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Scene = db.Scene
	mine.Quotes = db.Quotes
	if mine.Quotes == nil {
		mine.Quotes = make([]string, 0, 1)
	}
}

func (mine *RoomInfo) UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateRoomBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *RoomInfo) Devices() []*DeviceInfo {
	dbs, _ := nosql.GetDevicesByRoom(mine.Scene, mine.UID)
	devices := make([]*DeviceInfo, 0, 5)
	for _, device := range dbs {
		tmp := new(DeviceInfo)
		tmp.initInfo(device)
		devices = append(devices, tmp)
	}
	return devices
}

func (mine *RoomInfo) UpdateQuotes(operator string, list []string) error {
	if list == nil {
		list = make([]string, 0, 1)
	}

	err := nosql.UpdateRoomQuotes(mine.UID, operator, list)
	if err == nil {
		mine.Quotes = list
		mine.Operator = operator
	}
	return err
}

func (mine *RoomInfo) HadQuote(quote string) bool {
	for i := 0; i < len(mine.Quotes); i += 1 {
		if mine.Quotes[i] == quote {
			return true
		}
	}
	return false
}

func (mine *RoomInfo) HadQuotes(quotes []string) bool {
	for i := 0; i < len(mine.Quotes); i += 1 {
		if tool.HasItem(quotes, mine.Quotes[i]) {
			return true
		}
	}
	return false
}

func (mine *RoomInfo) UpdateDisplays(sn, group, operator string, showing bool, displays []string) error {
	device := mine.GetDevice(sn)
	if device == nil {
		return errors.New("the device had not found by sn")
	}
	if displays == nil {
		displays = make([]string, 0, 1)
	}
	return device.UpdateShowings(operator, displays)
}

func (mine *RoomInfo) HadDevice(sn string) bool {
	devices := mine.Devices()
	for _, device := range devices {
		if device.SN == sn {
			return true
		}
	}
	return false
}

func (mine *RoomInfo) HadDeviceByType(tp uint8) bool {
	devices := mine.Devices()
	for _, device := range devices {
		if device.Type == uint32(tp) {
			return true
		}
	}
	return false
}

func (mine *RoomInfo) Products() []*pb.ProductInfo {
	devices := mine.Devices()
	list := make([]*pb.ProductInfo, 0, len(devices))
	for _, device := range devices {
		tmp := &pb.ProductInfo{Uid: device.SN, Room: device.Room, Area: device.Area, Type: device.Type, Remark: device.Remark, Question: device.Question}
		tmp.Displays = cacheCtx.SwitchDisplays(device.Type, device.Displays)
		list = append(list, tmp)
	}
	return list
}

func (mine *cacheContext) SwitchDisplays(tp uint32, arr []string) []*pb.DisplayInfo {
	list := make([]*pb.DisplayInfo, 0, 10)
	tmp := new(pb.DisplayInfo)
	tmp.Group = ""
	tmp.Showings = arr
	list = append(list, tmp)
	//prepares := mine.GetPrepareDisplays(tp)
	//for _, prepare := range prepares {
	//	tmp := new(pb.DisplayInfo)
	//	tmp.Group = prepare.Group
	//	tmp.Prepares = prepare.Showings
	//	tmp.Showings = arr
	//	list = append(list, tmp)
	//}
	return list
}

func (mine *RoomInfo) GetDevice(sn string) *DeviceInfo {
	devices := mine.Devices()
	for _, device := range devices {
		if device.SN == sn {
			return device
		}
	}
	return nil
}

//func (mine *RoomInfo)GetPrepareDisplays(tp uint32) []*proxy.DisplayInfo {
//	list := make([]*proxy.DisplayInfo, 0, 3)
//	for _, item := range mine.displays {
//		if item.Type == tp {
//			list = append(list, item)
//		}
//	}
//	return list
//}

func (mine *RoomInfo) AppendDevice(area, device, remark, operator string, tp uint32) error {
	if mine.HadDevice(device) {
		return nil
	}
	info, err := cacheCtx.checkDevice(mine.Scene, mine.UID, area, device, remark, operator, tp)
	if err == nil {
		return info.UpdateRoom(mine.UID, area, operator)
	}
	return err
}
