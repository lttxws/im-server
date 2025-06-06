package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"
	"time"

	"google.golang.org/protobuf/proto"
)

type GroupMsgActor struct {
	bases.BaseActor
}

func (actor *GroupMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if upMsg, ok := input.(*pbobjs.UpMsg); ok {
		logs.WithContext(ctx).Infof("group_id:%s\tmsg_type:%s\tflag:%d", bases.GetTargetIdFromCtx(ctx), upMsg.MsgType, upMsg.Flags)
		code, msgId, sendTime, msgSeq, clientMsgId, memberCount, modifiedMsg := services.SendGroupMsg(ctx, upMsg)
		userPubAck := bases.CreateGrpPubAckWraper(ctx, code, msgId, sendTime, msgSeq, clientMsgId, memberCount, modifiedMsg)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d\tmsg_id:%s", code, msgId)
	} else {
		userPubAck := bases.CreateGrpPubAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, "", time.Now().UnixMilli(), 0, "", 0, nil)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("upMsg is illigal. upMsg:%v", upMsg)
	}
}

func (actor *GroupMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
