package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"
	"strconv"
)

func QryFriends(ctx *httputils.HttpContext) {
	offset := ctx.Query("offset")
	count := 20
	var err error
	countStr := ctx.Query("count")
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			count = 20
		}
	}
	code, friends := services.QryFriends(ctx.ToRpcCtx(), &pbobjs.FriendListReq{
		Limit:  int64(count),
		Offset: offset,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &models.Friends{
		Items:  []*pbobjs.UserObj{},
		Offset: friends.Offset,
	}
	for _, friend := range friends.Items {
		ret.Items = append(ret.Items, &pbobjs.UserObj{
			UserId:   friend.UserId,
			Nickname: friend.Nickname,
			Avatar:   friend.Avatar,
			Pinyin:   friend.Pinyin,
			IsFriend: true,
		})
	}
	ctx.ResponseSucc(ret)
}

func QryFriendsWithPage(ctx *httputils.HttpContext) {
	var err error
	page := 1
	pageStr := ctx.Query("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
	}

	size := 20
	sizeStr := ctx.Query("size")
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			size = 20
		}
	}
	orderTag := ctx.Query("order_tag")
	code, friends := services.QryFriendsWithPage(ctx.ToRpcCtx(), &pbobjs.FriendListWithPageReq{
		Page:     int64(page),
		Size:     int64(size),
		OrderTag: orderTag,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &models.Friends{
		Items: []*pbobjs.UserObj{},
	}
	for _, friend := range friends.Items {
		ret.Items = append(ret.Items, &pbobjs.UserObj{
			UserId:   friend.UserId,
			Nickname: friend.Nickname,
			Avatar:   friend.Avatar,
			Pinyin:   friend.Pinyin,
			IsFriend: true,
		})
	}
	ctx.ResponseSucc(ret)
}

func AddFriend(ctx *httputils.HttpContext) {
	req := models.Friend{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddFriends(ctx.ToRpcCtx(), &pbobjs.FriendIdsReq{
		FriendIds: []string{req.FriendId},
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func DelFriend(ctx *httputils.HttpContext) {
	req := models.FriendIds{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelFriends(ctx.ToRpcCtx(), &pbobjs.FriendIdsReq{
		FriendIds: req.FriendIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func ApplyFriend(ctx *httputils.HttpContext) {
	req := pbobjs.ApplyFriend{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.ApplyFriend(ctx.ToRpcCtx(), &pbobjs.ApplyFriend{
		FriendId: req.FriendId,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func ConfirmFriend(ctx *httputils.HttpContext) {
	req := pbobjs.ConfirmFriend{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.ConfirmFriend(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func MyFriendApplications(ctx *httputils.HttpContext) {
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyFriendApplications(ctx.ToRpcCtx(), &pbobjs.QryFriendApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func MyPendingFriendApplications(ctx *httputils.HttpContext) {
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyPendingFriendApplications(ctx.ToRpcCtx(), &pbobjs.QryFriendApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func FriendApplications(ctx *httputils.HttpContext) {
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryFriendApplications(ctx.ToRpcCtx(), &pbobjs.QryFriendApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}
