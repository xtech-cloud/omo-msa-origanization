package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"omo.msa.organization/cache"
)

type GroupService struct {}

func switchGroup(info *cache.GroupInfo) *pb.GroupInfo {
	tmp := new(pb.GroupInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Location = info.Location
	tmp.Master = info.Master
	tmp.Assistant = info.Assistant
	tmp.Contact = info.Contact
	tmp.Members = info.AllMembers()
	tmp.Scene = info.Scene
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Address = new(pb.AddressInfo)
	tmp.Address.Country = info.Address.Country
	tmp.Address.Province = info.Address.Province
	tmp.Address.City = info.Address.City
	tmp.Address.Zone = info.Address.Zone
	return tmp
}

func (mine *GroupService)AddOne(ctx context.Context, in *pb.ReqGroupAdd, out *pb.ReplyGroupOne) error {
	path := "group.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pb.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Scene)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
		return nil
	}

	if scene.HadGroupByName(in.Name) {
		out.Status = outError(path,"not found the scene ", pb.ResultStatus_Repeated)
		return nil
	}

	group, err := scene.CreateGroup(in)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchGroup(group)
	out.Status = outLog(path, out)
	return nil
}

func (mine *GroupService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyGroupOne) error {
	path := "group.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	var info *cache.GroupInfo
	if len(in.Parent) > 0 {
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
			return nil
		}
		info = scene.GetGroup(in.Uid)
	}else{
		info = cache.Context().GetGroup(in.Uid)
	}

	if info == nil {
		out.Status = outError(path,"the group not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchGroup(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *GroupService)GetByContact(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyGroupList) error {
	path := "group.getByContact"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the phone is empty ", pb.ResultStatus_Empty)
		return nil
	}
	list := cache.Context().GetGroupByContact(in.Uid)
	out.List = make([]*pb.GroupInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchGroup(info))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *GroupService)GetByUser(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyGroupList) error {
	path := "group.getByUser"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the user is empty ", pb.ResultStatus_Empty)
		return nil
	}
	list := cache.Context().GetGroupByMember(in.Uid)
	out.List = make([]*pb.GroupInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchGroup(info))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *GroupService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "group.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveGroup(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *GroupService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyGroupList) error {
	path := "group.getList"
	inLog(path, in)
	scene := cache.Context().GetScene(in.Parent)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
		return nil
	}
	total,max,list := scene.GetGroups(in.Number, in.Page)
	out.List = make([]*pb.GroupInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchGroup(value))
	}
	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *GroupService) UpdateBase (ctx context.Context, in *pb.ReqGroupUpdate, out *pb.ReplyInfo) error {
	path := "group.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Uid)
	if info == nil {
		out.Status = outError(path,"the group not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}

	if len(in.Name) > 0 || len(in.Remark) > 0 {
		scene := cache.Context().GetScene(info.Scene)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pb.ResultStatus_NotExisted)
			return nil
		}
		if in.Name != info.Name && scene.HadGroupByName(in.Name) {
			out.Status = outError(path,"the department name repeated ", pb.ResultStatus_Repeated)
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

func (mine *GroupService) UpdateAddress (ctx context.Context, in *pb.RequestAddress, out *pb.ReplyGroupOne) error {
	path := "group.updateAddress"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Uid)
	if info == nil {
		out.Status = outError(path,"the group not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Location != info.Location {
		err = info.UpdateLocation(in.Location, in.Operator)
		if err != nil {
			out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
			return nil
		}
	}
	err = info.UpdateAddress(in.Country, in.Province, in.City, in.Zone, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchGroup(info)
	out.Status = outLog(path, out)
	return err
}

func (mine *GroupService) UpdateLocation (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "group.updateLocation"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Group not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateLocation(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *GroupService) UpdateContact (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "group.updateContact"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Group not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateContact(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *GroupService) UpdateMaster (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "group.updateMaster"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Uid)
	if info == nil {
		out.Status = outError(path,"the group not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateMaster(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *GroupService) UpdateAssistant (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "group.updateAssistant"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Uid)
	if info == nil {
		out.Status = outError(path,"the group not found ", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateMaster(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *GroupService) AppendMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "group.appendMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Parent)
	if info == nil {
		out.Status = outError(path,"the group not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}

func (mine *GroupService) SubtractMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "group.subtractMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the uid is empty ", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetGroup(in.Parent)
	if info == nil {
		out.Status = outError(path,"the group not found ", pb.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}


