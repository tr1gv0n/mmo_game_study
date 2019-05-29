package main

import (
	"fmt"
	"mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)
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

	fmt.Println("---->player ID = ",p.Pid,"Online ...",",Player num=",len(core.WorldMgrObj.Players))
}

func main()  {
	s := znet.NewServer("MMO Game Server")
	//注册一些 链接创建/销毁的 Hook钩子函数
	s.AddOnConnStart(OnConnectionAdd)
	//注册一些路由业务
	s.Serve()
}