package main

import (
	"fmt"
	"mmo_game/apis"
	"mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

func OnConnectionLost(conn ziface.IConnection)  {
	//得到当前下线玩家
	pid,_ := conn.GetProperty("pid")

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//玩家下线业务
	player.OffLine()
}
//当前客户端建立链接之后触发Hook函数
func OnConnectionAdd(conn ziface.IConnection)  {
	fmt.Println("conn Add...")
	//创建一个玩家 将链接和玩家模块绑定
	p := core.NewPlayer(conn)
	//给客户端发送一个msgID:1
	p.ReturnPid()
	//给客户端发送一个msgID:200
	p.ReturnPlayerPosition()
	//上线成功了
	//将玩家对象添加到世界管理器中
	core.WorldMgrObj.AddPlayer(p)
	//给conn添加一个属性 pid属性
	conn.SetProperty("pid",p.Pid)
	//同步周边玩家
	p.SyncSurrounding()

	fmt.Println("---->player ID = ",p.Pid,"Online ...",",Player num=",len(core.WorldMgrObj.Players))
}

func main()  {
	s := znet.NewServer("MMO Game Server")
	//注册一些 链接创建/销毁的 Hook钩子函数
	s.AddOnConnStart(OnConnectionAdd)
	s.AddOnConnStop(OnConnectionLost)
	//针对MsgID2 建立路由业务
	s.AddRouter(2,&apis.WorldChat{})
	s.AddRouter(3,&apis.Move{})
	//注册一些路由业务
	s.Serve()
}