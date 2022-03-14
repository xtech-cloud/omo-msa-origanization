package cache

import (
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"omo.msa.organization/proxy"
	"omo.msa.organization/proxy/nosql"
	"omo.msa.organization/tool"
	"time"
)

type RoomInfo struct {
	baseInfo
	Remark string
	Scene string
	Quotes []string
	devices []*proxy.DeviceInfo
}

func (mine *cacheContext)GetRoom(uid string) *RoomInfo {
	for _, scene := range mine.scenes {
		info := scene.GetRoom(uid)
		if  info != nil {
			return info
		}
	}
	return nil
}

func (mine *cacheContext)GetRoomsByDevice(sn string) []*RoomInfo {
	list := make([]*RoomInfo, 0, 10)
	for _, scene := range mine.scenes {
		arr := scene.GetRoomsByDevice(sn)
		if  arr != nil && len(arr) > 0 {
			list = append(list, arr...)
		}
	}
	return list
}

func (mine *cacheContext)GetRoomsByQuote(quote string) []*RoomInfo {
	list := make([]*RoomInfo, 0, 10)
	for _, scene := range mine.scenes {
		arr := scene.GetRoomsByQuote(quote)
		if  arr != nil && len(arr) > 0 {
			list = append(list, arr...)
		}
	}
	return list
}

func (mine *cacheContext)HadBindDeviceInRoom(sn string) bool {
	for _, scene := range mine.scenes {
		arr := scene.GetRoomsByDevice(sn)
		if  arr != nil && len(arr) > 0 {
			return true
		}
	}
	return false
}

func (mine *cacheContext)GetRoomBy(scene, uid string) *RoomInfo {
	for _, item := range mine.scenes {
		if item.UID == scene {
			return item.GetRoom(uid)
		}
	}
	return nil
}

func (mine *cacheContext)RemoveRoom(uid, operator string) error {
	for _, scene := range mine.scenes {
		if scene.HadRoom(uid) {
			return scene.RemoveRoom(uid, operator)
		}
	}
	return nil
}

func (mine *RoomInfo)initInfo(db *nosql.Room)  {
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
	mine.devices = db.Devices
	if mine.devices == nil {
		mine.devices = make([]*proxy.DeviceInfo, 0, 1)
	}

}

func (mine *RoomInfo)UpdateBase(name, remark, operator string) error {
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

func (mine *RoomInfo)UpdateQuotes(operator string, list []string) error {
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

func (mine *RoomInfo)HadQuote(quote string) bool {
	for i := 0;i < len(mine.Quotes);i += 1 {
		if mine.Quotes[i] == quote {
			return true
		}
	}
	return false
}

func (mine *RoomInfo)HadQuotes(quotes []string) bool {
	for i := 0;i < len(mine.Quotes);i += 1 {
		if tool.HasItem(quotes, mine.Quotes[i]) {
			return true
		}
	}
	return false
}

func (mine *RoomInfo)HadDevice(sn string) bool {
	for _, device := range mine.devices {
		if device.SN == sn {
			return true
		}
	}
	return false
}

func (mine *RoomInfo)HadDeviceByType(tp uint8) bool {
	for _, device := range mine.devices {
		if device.Type == tp {
			return true
		}
	}
	return false
}

func (mine *RoomInfo)Devices() []*pb.ProductInfo {
	devices := make([]*pb.ProductInfo, 0, len(mine.devices))
	for _, device := range mine.devices {
		devices = append(devices, &pb.ProductInfo{Uid: device.SN, Type: uint32(device.Type), Remark: device.Remark})
	}
	return devices
}

func (mine *RoomInfo)AppendDevice(device, remark string, tp uint32) error {
	if mine.HadDevice(device){
		return nil
	}
	info := &proxy.DeviceInfo{SN: device, Remark: remark, Type: uint8(tp), Updated: time.Now()}
	err := nosql.AppendRoomDevice(mine.UID, info)
	if err == nil {
		mine.devices = append(mine.devices, info)
	}
	return err
}

func (mine *RoomInfo)SubtractDevice(uid string) error {
	if !mine.HadDevice(uid){
		return nil
	}
	err := nosql.SubtractRoomDevice(mine.UID, uid)
	if err == nil {
		for i := 0;i < len(mine.devices);i += 1 {
			if mine.devices[i].SN == uid {
				if i == len(mine.devices) - 1 {
					mine.devices = append(mine.devices[:i])
				}else{
					mine.devices = append(mine.devices[:i], mine.devices[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

