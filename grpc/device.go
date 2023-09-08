package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
)

type DeviceService struct{}

func switchDevice(info *cache.DeviceInfo) *pb.DeviceInfo {
	tmp := new(pb.DeviceInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Sn = info.SN
	tmp.Status = uint32(info.Status)
	tmp.Owner = info.Scene
	tmp.Quote = info.Quote
	tmp.Activated = info.ActiveTime
	tmp.Expiry = info.Expired
	tmp.Os = info.OS
	tmp.Certificate = info.Certificate
	tmp.Type = uint32(info.Type)
	return tmp
}

func (mine *DeviceService) AddOne(ctx context.Context, in *pb.ReqDeviceAdd, out *pb.ReplyDeviceInfo) error {
	path := "device.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	tmp, _ := cache.Context().GetDeviceBySN(in.Sn)
	if tmp != nil {
		out.Status = outError(path, "the sn is repeated ", pbstatus.ResultStatus_Repeated)
		return nil
	}

	info, err := cache.Context().CreateDevice(in.Owner, in.Name, in.Sn, in.Remark, in.Operator, uint8(in.Type))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchDevice(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DeviceService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyDeviceInfo) error {
	path := "device.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.DeviceInfo
	var er error
	if in.Operator == "sn" {
		info, er = cache.Context().GetDeviceBySN(in.Uid)
	} else {
		info, er = cache.Context().GetDevice(in.Uid)
	}

	if er != nil {
		out.Status = outError(path, "the device not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchDevice(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DeviceService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "device.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.Key == "count" {
		out.Count = uint32(cache.Context().GetDeviceCount())
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *DeviceService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "device.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetDevice(in.Uid)
	if er != nil {
		out.Status = outError(path, "the device not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *DeviceService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyDeviceList) error {
	path := "device.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *DeviceService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyDeviceList) error {
	path := "device.getListByFilter"
	inLog(path, in)
	var list []*cache.DeviceInfo
	var err error
	if in.Scene == "" {
		if in.Key == "status" {
			st := parseInt(in.Value)
			list, err = cache.Context().GetDevicesByStatus(int32(st))
		} else {
			err = errors.New("the key not defined")
		}
	} else {
		if in.Key == "" {
			list, err = cache.Context().GetDevicesByOwner(in.Scene)
		} else if in.Key == "type" {

		} else {
			err = errors.New("the key not defined")
		}
	}

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.DeviceInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchDevice(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *DeviceService) UpdateBase(ctx context.Context, in *pb.ReqDeviceBase, out *pb.ReplyInfo) error {
	path := "device.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetDevice(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *DeviceService) UpdateByFilter(ctx context.Context, in *pb.ReqUpdateFilter, out *pb.ReplyInfo) error {
	path := "device.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetDevice(in.Uid)
	if er != nil {
		out.Status = outError(path, "the device not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "certificate" {
		err = info.UpdateCertificate(in.Value, in.Operator)
	} else if in.Key == "scene" {
		err = info.UpdateScene(in.Value, in.Operator)
	} else if in.Key == "aspect" {
		err = info.UpdateAspect(in.Value, in.Operator)
	} else if in.Key == "type" {
		tp := parseInt(in.Value)
		err = info.UpdateType(in.Operator, uint8(tp))
	} else {
		err = errors.New("the field not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *DeviceService) Bind(ctx context.Context, in *pb.ReqDeviceBind, out *pb.ReplyInfo) error {
	path := "device.bind"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetDeviceBySN(in.Uid)
	if er != nil {
		out.Status = outError(path, "the device not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	//if info.Status != cache.DeviceIdle {
	//	out.Status = outError(path, "the device had bind ", pbstatus.ResultStatus_Prohibition)
	//	return nil
	//}
	err := info.Bind(in.Quote, in.Os, in.Operator, in.Activated, uint64(in.Expiry))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
