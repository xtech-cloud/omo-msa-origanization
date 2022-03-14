package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
	"strconv"
)

type RoomService struct {}

func switchRoom(info *cache.RoomInfo) *pb.RoomInfo {
	tmp := new(pb.RoomInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Owner = info.Scene
	tmp.Remark = info.Remark
	tmp.Quotes = info.Quotes
	tmp.Devices = info.Devices()
	return tmp
}

func (mine *RoomService)AddOne(ctx context.Context, in *pb.ReqRoomAdd, out *pb.ReplyRoomInfo) error {
	path := "room.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Owner)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	if scene.HadRoomByName(in.Name) {
		out.Status = outError(path,"the name repeated ", pbstatus.ResultStatus_Repeated)
		return nil
	}

	Room, err := scene.CreateRoom(in)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRoom(Room)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRoomInfo) error {
	path := "room.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.RoomInfo
	if len(in.Parent) > 0 {
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info = scene.GetRoom(in.Uid)
	}else{
		info = cache.Context().GetRoom(in.Uid)
	}

	if info == nil {
		out.Status = outError(path,"the room not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchRoom(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "room.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path,"the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "room.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveRoom(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService)GetList(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyRoomList) error {
	path := "room.getList"
	inLog(path, in)
	var list []*cache.RoomInfo
	if in.Parent == "" {
		if in.Key == "device" {
			list = cache.Context().GetRoomsByDevice(in.Value)
		}else if in.Key == "quote" {
			list = cache.Context().GetRoomsByQuote(in.Value)
		}else{
			list = make([]*cache.RoomInfo, 0, 1)
		}
	}else{
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Key == "" {
			list = scene.GetRooms()
		}else if in.Key == "product" {
			tp,er := strconv.ParseUint(in.Value, 10, 32)
			if er != nil {
				out.Status = outError(path,er.Error(), pbstatus.ResultStatus_DBException)
				return nil
			}
			list = scene.GetRoomsByType(uint8(tp))
		}else if in.Key == "quote" {
			list = scene.GetRoomsByQuote(in.Value)
		}else if in.Key == "device" {
			list = scene.GetRoomsByDevice(in.Value)
		}else{
			list = make([]*cache.RoomInfo, 0, 1)
		}
	}


	out.List = make([]*pb.RoomInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchRoom(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *RoomService) UpdateBase (ctx context.Context, in *pb.ReqRoomUpdate, out *pb.ReplyInfo) error {
	path := "room.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetRoom(in.Uid)
	if info == nil {
		out.Status = outError(path,"the room not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if len(in.Name) > 0 || len(in.Remark) > 0 {
		scene := cache.Context().GetScene(info.Scene)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Name != info.Name && scene.HadRoomByName(in.Name) {
			out.Status = outError(path,"the room name repeated ", pbstatus.ResultStatus_Repeated)
			return nil
		}
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *RoomService) UpdateQuotes (ctx context.Context, in *pb.ReqRoomQuotes, out *pb.ReplyRoomInfo) error {
	path := "room.updateQuotes"
	inLog(path, in)
	if len(in.Scene) < 1 || len(in.Room) < 1 {
		out.Status = outError(path,"the scene or room is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Scene)
	if scene == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	info := scene.GetRoom(in.Room)
	if info == nil {
		out.Status = outError(path,"the room not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	scene.ClearQuotes(in.Operator, in.List)
	err := info.UpdateQuotes(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *RoomService) AppendDevice (ctx context.Context, in *pb.ReqRoomDevice, out *pb.ReplyRoomDevices) error {
	path := "room.appendDevice"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if cache.Context().HadBindDeviceInRoom(in.Device) {
		out.Status = outError(path,"the device had bind by other room", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := cache.Context().GetRoom(in.Uid)
	if info == nil {
		out.Status = outError(path,"the room not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendDevice(in.Device, in.Remark, in.Type)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Devices()
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService) SubtractDevice (ctx context.Context, in *pb.ReqRoomDevice, out *pb.ReplyRoomDevices) error {
	path := "room.subtractDevice"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetRoom(in.Uid)
	if info == nil {
		out.Status = outError(path,"the room not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractDevice(in.Device)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Devices()
	out.Status = outLog(path, out)
	return nil
}


