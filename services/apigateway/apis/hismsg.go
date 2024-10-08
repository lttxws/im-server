package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"sort"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func RecallHisMsgs(ctx *gin.Context) {
	var req models.RecallHisMsgsReq
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}

	converId := commonservices.GetConversationId(req.FromId, req.TargetId, pbobjs.ChannelType(req.ChannelType))
	code, _, err := services.SyncApiCall(ctx, "recall_msg", req.FromId, converId, &pbobjs.RecallMsgReq{
		TargetId:    req.TargetId,
		ChannelType: pbobjs.ChannelType(req.ChannelType),
		MsgId:       req.MsgId,
		MsgTime:     req.MsgTime,
		Exts:        commonservices.Map2KvItems(req.Exts),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func CleanHisMsgs(ctx *gin.Context) {
	var req models.CleanHisMsgsReq
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "clean_hismsg", req.FromId, req.TargetId, &pbobjs.CleanHisMsgReq{
		TargetId:        req.TargetId,
		ChannelType:     pbobjs.ChannelType(req.ChannelType),
		CleanMsgTime:    req.CleanTime,
		CleanTimeOffset: req.CleanTimeOffset,
		CleanScope:      int32(req.CleanScope),
		SenderId:        req.SenderId,
	}, nil)

	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	tools.SuccessHttpResp(ctx, nil)
}

func DelHisMsgs(ctx *gin.Context) {
	var req models.DelHisMsgsReq
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	//del_hismsg
	msgs := []*pbobjs.SimpleMsg{}
	for _, m := range req.Msgs {
		msgs = append(msgs, &pbobjs.SimpleMsg{
			MsgId: m.MsgId,
		})
	}
	code, _, err := services.SyncApiCall(ctx, "del_msg", req.FromId, req.TargetId, &pbobjs.DelHisMsgsReq{
		TargetId:    req.TargetId,
		ChannelType: pbobjs.ChannelType(req.ChannelType),
		Msgs:        msgs,
		DelScope:    int32(req.DelScope),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryHisMsgs(ctx *gin.Context) {
	channelTypeStr := ctx.Query("channel_type")
	fromIdStr := ctx.Query("from_id")
	targetIdStr := ctx.Query("target_id")

	channelTypeInt, err := tools.String2Int64(channelTypeStr)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_MSG_PARAM_ILLEGAL)
		return
	}
	channelType := pbobjs.ChannelType_Private
	if channelTypeInt == 1 {
		channelType = pbobjs.ChannelType_Private
	} else if channelTypeInt == 2 {
		channelType = pbobjs.ChannelType_Group
	}

	startTimeStr := ctx.Query("start")
	startTimeInt, err := tools.String2Int64(startTimeStr)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_MSG_PARAM_ILLEGAL)
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
	converId := commonservices.GetConversationId(fromIdStr, targetIdStr, channelType)
	code, resp, err := services.SyncApiCall(ctx, "qry_hismsgs", fromIdStr, converId, &pbobjs.QryHisMsgsReq{
		TargetId:    targetIdStr,
		ChannelType: channelType,
		StartTime:   startTimeInt,
		Count:       int32(count),
		Order:       int32(order),
	}, func() proto.Message {
		return &pbobjs.DownMsgSet{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	msgs := resp.(*pbobjs.DownMsgSet)
	respHisMsgs := &models.HisMsgs{
		Msgs:       []*models.HisMsg{},
		IsFinished: msgs.IsFinished,
	}
	if order == 0 {
		sort.Slice(msgs.Msgs, func(i, j int) bool {
			return msgs.Msgs[i].MsgTime > msgs.Msgs[j].MsgTime
		})
	}
	for _, msg := range msgs.Msgs {
		respHisMsgs.Msgs = append(respHisMsgs.Msgs, &models.HisMsg{
			SenderId:    msg.SenderId,
			ReceiverId:  msg.TargetId,
			ChannelType: int32(msg.ChannelType),
			MsgId:       msg.MsgId,
			MsgTime:     msg.MsgTime,
			MsgType:     msg.MsgType,
			MsgContent:  string(msg.MsgContent),
		})
	}
	tools.SuccessHttpResp(ctx, respHisMsgs)
}

func MarkRead(ctx *gin.Context) {
	var req MarkReadReq
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	pbMsgs := []*pbobjs.SimpleMsg{}
	for _, msgId := range req.MsgIds {
		pbMsgs = append(pbMsgs, &pbobjs.SimpleMsg{
			MsgId: msgId,
		})
	}
	converId := commonservices.GetConversationId(req.UserId, req.TargetId, pbobjs.ChannelType(req.ChannelType))
	services.AsyncApiCall(ctx, "mark_read", req.UserId, converId, &pbobjs.MarkReadReq{
		TargetId:    req.TargetId,
		ChannelType: pbobjs.ChannelType(req.ChannelType),
		Msgs:        pbMsgs,
	})
	tools.SuccessHttpResp(ctx, nil)
}

type MarkReadReq struct {
	UserId      string   `json:"user_id"`
	TargetId    string   `json:"target_id"`
	ChannelType int32    `json:"channel_type"`
	MsgIds      []string `json:"msg_ids"`
}

func ImportHisMsg(ctx *gin.Context) {
	var msg models.HisMsg
	if err := ctx.BindJSON(&msg); err != nil || msg.SenderId == "" || msg.MsgType == "" || msg.MsgContent == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if msg.ReceiverId == "" {
		msg.ReceiverId = msg.TargetId
	}
	if msg.ReceiverId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if msg.ChannelType == int32(pbobjs.ChannelType_Private) {
		services.AsyncSendMsg(ctx, "imp_pri_msg", msg.SenderId, msg.ReceiverId, &pbobjs.UpMsg{
			MsgType:    msg.MsgType,
			MsgContent: []byte(msg.MsgContent),
			Flags:      handleHisMsgFlag(msg),
			MsgTime:    msg.MsgTime,
		}, false, "")
	} else if msg.ChannelType == int32(pbobjs.ChannelType_Group) {
		services.AsyncSendMsg(ctx, "imp_grp_msg", msg.SenderId, msg.ReceiverId, &pbobjs.UpMsg{
			MsgType:    msg.MsgType,
			MsgContent: []byte(msg.MsgContent),
			Flags:      handleHisMsgFlag(msg),
			MsgTime:    msg.MsgTime,
		}, false, "")
	} else {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_ILLEGAL)
		return
	}
	services.AsyncApiCall(ctx, "", msg.SenderId, msg.ReceiverId, &pbobjs.DownMsg{})
	tools.SuccessHttpResp(ctx, nil)
}

func handleHisMsgFlag(msg models.HisMsg) int32 {
	var flag int32 = 0
	if msg.IsStorage == nil || *msg.IsStorage {
		flag = commonservices.SetStoreMsg(flag)
	}
	if msg.IsCount == nil || *msg.IsCount {
		flag = commonservices.SetCountMsg(flag)
	}
	return flag
}
