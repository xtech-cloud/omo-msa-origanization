package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
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
	tmp.Remark = info.Remark
	tmp.Quotes = info.Quotes
	tmp.Devices = info.Devices()
	return tmp
}

func (mine *RoomService)AddOne(ctx context.Context, in *pb.ReqRoomAdd, out *pb.ReplyRoomInfo) error {
	path := "Room.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pb.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Owner)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
		return nil
	}

	if scene.HadRoomByName(in.Name) {
		out.Status = outError(path,"the name repeated ", pb.ResultStatus_Repeated)
		return nil
	}

	Room, err := scene.CreateRoom(in)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRoom(Room)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRoomInfo) error {
	path := "Room.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	var info *cache.RoomInfo
	if len(in.Parent) > 0 {
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
			return nil
		}
		info = scene.GetRoom(in.Uid)
	}else{
		info = cache.Context().GetRoom(in.Uid)
	}

	if info == nil {
		out.Status = outError(path,"the room not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchRoom(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "Room.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path,"the user is empty ", pb.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "Room.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveRoom(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService)GetList(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyRoomList) error {
	path := "Room.getList"
	inLog(path, in)
	scene := cache.Context().GetScene(in.Parent)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
		return nil
	}
	var list []*cache.RoomInfo
	if in.Key == "" {
		list = scene.GetRooms()
	}else if in.Key == "product" {
		tp,er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path,er.Error(), pb.ResultStatus_DBException)
			return nil
		}
		list = scene.GetRoomsByType(uint8(tp))
	}else if in.Key == "quote" {
		list = scene.GetRoomsByQuote(in.Value)
	}else{
		list = make([]*cache.RoomInfo, 0, 1)
	}

	out.List = make([]*pb.RoomInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchRoom(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *RoomService) UpdateBase (ctx context.Context, in *pb.ReqRoomUpdate, out *pb.ReplyInfo) error {
	path := "Room.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetRoom(in.Uid)
	if info == nil {
		out.Status = outError(path,"the room not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if len(in.Name) > 0 || len(in.Remark) > 0 {
		scene := cache.Context().GetScene(info.Scene)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
			return nil
		}
		if in.Name != info.Name && scene.HadRoomByName(in.Name) {
			out.Status = outError(path,"the room name repeated ", pb.ResultStatus_Repeated)
			return nil
		}
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *RoomService) UpdateQuotes (ctx context.Context, in *pb.ReqRoomQuotes, out *pb.ReplyRoomInfo) error {
	path := "Room.updateQuotes"
	inLog(path, in)
	if len(in.Scene) < 1 || len(in.Room) < 1 {
		out.Status = outError(path,"the scene or room is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetRoomBy(in.Scene, in.Room)
	if info == nil {
		out.Status = outError(path,"the room not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateQuotes(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *RoomService) AppendDevice (ctx context.Context, in *pb.ReqRoomDevice, out *pb.ReplyRoomDevices) error {
	path := "scene.appendDevice"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetRoom(in.Uid)
	if info == nil {
		out.Status = outError(path,"the room not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendDevice(in.Device, in.Remark, in.Type)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.Devices()
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoomService) SubtractDevice (ctx context.Context, in *pb.ReqRoomDevice, out *pb.ReplyRoomDevices) error {
	path := "scene.subtractDevice"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetRoom(in.Uid)
	if info == nil {
		out.Status = outError(path,"the room not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractDevice(in.Device)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.Devices()
	out.Status = outLog(path, out)
	return nil
}


