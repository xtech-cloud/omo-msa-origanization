package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
)

type RegionService struct {}

func switchRegion(info *cache.RegionInfo) *pb.RegionInfo {
	tmp := new(pb.RegionInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Location = info.Location
	tmp.Master = info.Master
	tmp.Members = info.Members
	tmp.Scene = info.Scene
	tmp.Parent = info.Parent
	tmp.Entity = info.Entity
	tmp.Address = new(pb.AddressInfo)
	tmp.Address.Country = info.Address.Country
	tmp.Address.Province = info.Address.Province
	tmp.Address.City = info.Address.City
	tmp.Address.Zone = info.Address.Zone
	return tmp
}

func (mine *RegionService)AddOne(ctx context.Context, in *pb.ReqRegionAdd, out *pb.ReplyRegionOne) error {
	path := "region.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Scene)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	if scene.HadRegionByName(in.Name) {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_Repeated)
		return nil
	}

	region, err := scene.CreateRegion(in)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRegion(region)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RegionService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRegionOne) error {
	path := "region.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.RegionInfo
	var err error
	if len(in.Parent) > 0 {
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info,err = scene.GetRegion(in.Uid)
	}else{
		info,err = cache.Context().GetRegion(in.Uid)
	}

	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchRegion(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RegionService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "region.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path,"the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *RegionService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "region.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetRegion(in.Uid)
	if er != nil {
		out.Status = outError(path,er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *RegionService)GetList(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyRegionList) error {
	path := "region.getList"
	inLog(path, in)
	scene := cache.Context().GetScene(in.Scene)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var total uint32
	var max uint32
	var list []*cache.RegionInfo
	if in.Key == "" {
		total,max,list = scene.GetRegions(in.Number, in.Page)
		out.List = make([]*pb.RegionInfo, 0, len(list))
		for _, value := range list {
			out.List = append(out.List, switchRegion(value))
		}
	}

	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *RegionService) UpdateBase (ctx context.Context, in *pb.ReqRegionUpdate, out *pb.ReplyInfo) error {
	path := "region.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetRegion(in.Uid)
	if er != nil {
		out.Status = outError(path,er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if len(in.Name) > 0 || len(in.Remark) > 0 {
		scene := cache.Context().GetScene(info.Scene)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Name != info.Name && scene.HadRegionByName(in.Name) {
			out.Status = outError(path,"the department name repeated ", pbstatus.ResultStatus_Repeated)
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

func (mine *RegionService) UpdateByFilter(ctx context.Context, in *pb.ReqUpdateFilter, out *pb.ReplyInfo) error {
	path := "region.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	_,er := cache.Context().GetRegion(in.Uid)
	if er != nil {
		out.Status = outError(path,er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *RegionService) UpdateAddress (ctx context.Context, in *pb.RequestAddress, out *pb.ReplyRegionOne) error {
	path := "region.updateAddress"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetRegion(in.Uid)
	if er != nil {
		out.Status = outError(path,er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Location != info.Location {
		err = info.UpdateLocation(in.Location, in.Operator)
		if err != nil {
			out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}
	err = info.UpdateAddress(in.Country, in.Province, in.City, in.Zone, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRegion(info)
	out.Status = outLog(path, out)
	return err
}

func (mine *RegionService) AppendMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "region.appendMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetRegion(in.Parent)
	if er != nil {
		out.Status = outError(path,er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Members
	out.Status = outLog(path, out)
	return nil
}

func (mine *RegionService) SubtractMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "region.subtractMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetRegion(in.Parent)
	if er != nil {
		out.Status = outError(path,er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Members
	out.Status = outLog(path, out)
	return nil
}


