package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
)

type MaintainService struct{}

func switchMaintain(info *cache.MaintainInfo) *pb.MaintainInfo {
	tmp := new(pb.MaintainInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator

	tmp.Scene = info.Scene
	tmp.Type = uint32(info.Type)
	tmp.Remark = info.Remark
	tmp.Date = info.Date
	tmp.Area = info.Area
	tmp.Submitter = info.Submitter
	return tmp
}

func (mine *MaintainService) AddOne(ctx context.Context, in *pb.ReqMaintainAdd, out *pb.ReplyMaintainInfo) error {
	path := "maintain.addOne"
	inLog(path, in)

	area, er := cache.Context().GetArea(in.Area)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	info, err := cache.Context().CreateMaintain(in, area.Device)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchMaintain(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MaintainService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMaintainInfo) error {
	path := "maintain.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the maintain uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetMaintain(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchMaintain(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MaintainService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "maintain.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the maintain uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *MaintainService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyMaintainList) error {
	path := "maintain.getByFilter"
	inLog(path, in)
	var array []*cache.MaintainInfo
	if in.Key == "" {
		array, _ = cache.Context().GetMaintainByScene(in.Scene)
	} else if in.Key == "area" {
		array, _ = cache.Context().GetMaintainByArea(in.Scene, in.Value)
	}
	out.List = make([]*pb.MaintainInfo, 0, len(array))
	for _, info := range array {
		item := switchMaintain(info)
		out.List = append(out.List, item)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *MaintainService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "maintain.getStatistic"
	inLog(path, in)
	if in.Key == "size" {

	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *MaintainService) UpdateByFilter(ctx context.Context, in *pb.ReqUpdateFilter, out *pb.ReplyInfo) error {
	path := "maintain.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the motion uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}
