package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mmo_game/core"
	"mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

type Move struct {
	znet.BaseRouter
}

func ( m *Move)Handle(request ziface.IRequest)  {
	proto_msg := &pb.Position{}
	proto.Unmarshal(request.GetMsg().GetMsgData(),proto_msg)

	pid,_ := request.GetConnection().GetProperty("pid")
	fmt.Println("player id = ",pid.(int32),"move-->",proto_msg.X,",",proto_msg.Z,",",proto_msg.V)
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePosition(proto_msg.X,proto_msg.Y,proto_msg.Z,proto_msg.V)
}
