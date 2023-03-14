package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
	"strings"
)

type AreaService struct{}

func switchArea(info *cache.AreaInfo) *pb.AreaInfo {
	tmp := new(pb.AreaInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Parent = info.Parent
	tmp.Owner = info.Owner
	tmp.Type = info.Type
	tmp.Template = info.Template
	tmp.Question = info.Question
	tmp.Device = info.Device
	tmp.Displays = info.Displays
	tmp.Catalog = info.Catalog
	return tmp
}

func (mine *AreaService) AddOne(ctx context.Context, in *pb.ReqAreaAdd, out *pb.ReplyAreaInfo) error {
	path := "area.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateArea(in.Name, in.Remark, in.Owner, in.Parent, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchArea(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AreaService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAreaInfo) error {
	path := "area.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetArea(in.Uid)
	if er != nil {
		out.Status = outError(path, "the area not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchArea(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AreaService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "area.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AreaService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "area.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetArea(in.Uid)
	if er != nil {
		out.Status = outError(path, "the area not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *AreaService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAreaList) error {
	path := "area.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AreaService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyAreaList) error {
	path := "area.getListByFilter"
	inLog(path, in)
	var list []*cache.AreaInfo
	var err error
	if in.Key == "" {
		list = cache.Context().GetAreasByOwner(in.Scene)
	} else if in.Key == "template" {
		list = cache.Context().GetAreasByTemplate(in.Scene, in.Value)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.AreaInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchArea(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AreaService) UpdateBase(ctx context.Context, in *pb.ReqAreaBase, out *pb.ReplyInfo) error {
	path := "area.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetArea(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, "")
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AreaService) UpdateByFilter(ctx context.Context, in *pb.ReqUpdateFilter, out *pb.ReplyInfo) error {
	path := "area.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetArea(in.Uid)
	if er != nil {
		out.Status = outError(path, "the area not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "question" {
		err = info.UpdateQuestion(in.Value, in.Operator)
	} else if in.Key == "sn" {
		err = info.UpdateDevice2(in.Value, in.Operator)
	} else if in.Key == "template" {
		err = info.UpdateTemplate(in.Value, in.Operator)
	} else if in.Key == "device" {
		arr := strings.Split(in.Value, ";")
		if len(arr) != 2 {
			err = errors.New("the device format is error")
		} else {
			tp := parseInt(arr[0])
			sn := arr[1]
			err = info.UpdateDevice(sn, in.Operator, uint32(tp))
		}
	} else if in.Key == "type" {
		tp := parseInt(in.Value)
		err = info.UpdateType(uint32(tp), in.Operator)
	} else if in.Key == "catalog" {
		err = info.UpdateCatalog(in.Value, in.Operator)
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
