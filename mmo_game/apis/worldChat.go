package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mmo_game/core"
	"mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

//世界聊天 路由业务
type WorldChat struct {
	znet.BaseRouter
}

func (wc *WorldChat)Handle(request ziface.IRequest)  {
	//1 解析客户端传递过来的protobuf数据
	proto_msg := &pb.Talk{}
	if err := proto.Unmarshal(request.GetMsg().GetMsgData(),proto_msg);err != nil {
		fmt.Println("Talk message unmarshal error",err)
		return
	}
	//通过获取连接属性 得到当前的玩家ID
	pid,err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("get pid err",err)
		return
	}
	//通过pid来得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//当前的聊天数据广播给全部的在线玩家
	//当前玩家的window客户端发送过来的消息
	player.SendTalkMsgToAll(proto_msg.GetContent())

}