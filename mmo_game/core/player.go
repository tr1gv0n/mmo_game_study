package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"mmo_game/pb"
	"sync"
	"zinx/ziface"
)

type Player struct {
	Pid int32	//玩家ID
	Conn ziface.IConnection	//当前玩家的链接(与对应客户端通信)
	X float32	//平面的x轴坐标
	Y float32	//高度
	Z float32	//平面的y轴坐标
	V float32	//玩家脸朝向的方向
}
// playerID 生成器
var PidGen int32 = 1	//用于生产玩家ID计数器
var IdLock sync.Mutex	//保护PidGen生成器的互斥锁
//初始化玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	//分配一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	//创建一个玩家对象
	p:=&Player{
		Pid:id,
		Conn:conn,
		X:float32(160 + rand.Intn(10)),	//随机生成玩家上线所在的x轴坐标
		Y:0,
		Z:float32(140 + rand.Intn(10)),	//随机在140坐标点附近 y轴坐标上线
		V:0,	//角度为0
	}
	return  p
}
//玩家可以和对端客户端发送消息的方法
func (p *Player)SendMsg(msgID uint32,proto_struct proto.Message) error {
	//要将proto结构体 转换成 二进制的数据
	binary_proto_data,err := proto.Marshal(proto_struct)
	if err != nil {
		fmt.Println("marshal proto error",err)
		return err
	}
	//再调用zinx原生的connecton.Send（msgID, 二进制数据）
	if err := p.Conn.Send(msgID,binary_proto_data);err != nil{
		fmt.Println("Player send error!",err)
		return err
	}
	return nil
}

/*
 服务器给客户端发送玩家初始ID
*/
func (p *Player)ReturnPid()  {
	proto_msg := &pb.SyncPid{
		Pid:p.Pid,
	}
	//将这个消息 发送给客户端
	p.SendMsg(1,proto_msg)
}
//服务器给客户端发送一个玩家的初始化位置信息
func (p *Player)ReturnPlayerPosition()  {
	//组建MsgID:200消息
	proto_msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:2,	//2 -坐标信息
		Data:&pb.BroadCast_P{
			P:&pb.Position{
				X:p.X,
				Y:p.Y,
				Z:p.Z,
				V:p.V,
			},
		},
	}
	//将这个消息 发送给客户端
	p.SendMsg(200,proto_msg)
}

//将聊天数据广播给全部的在线玩家
func (p *Player)SendTalkMsgToAll(content string)  {
	/*
		message BroadCast{
		int32 Pid=1;
		int32 Tp=2;
		oneof Data {
			string Content=3;
			Position P=4;
			int32 ActionData=5;
		}
	}
	*/
	//定义一个广播的proto消息数据类型
	proto_msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:1,
		Data:&pb.BroadCast_Content{
			Content:content,
		},
	}
	//获取全部的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//想全部的玩家进行广播proto_msg数据
	for _,player := range players{
		player.SendMsg(200,proto_msg)
	}
}

//得到当前玩家周围玩家
func (p *Player)GetSurroundingPlayers() []*Player  {
	pids := WorldMgrObj.AoiMgr.GetSurroundPIDsByPos(p.X,p.Z)
	fmt.Println("Surrounding players = ",pids)
	players := make([]*Player,0,len(pids))
	for _,pid := range pids{
		players = append(players,WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}
//将自己的消息同步给周围玩家
func (p *Player)SyncSurrounding()  {
	//获取当前玩家的周围九宫格的玩家
	players := p.GetSurroundingPlayers()
	//构建一个广播信息200，循环发送g
	proto_msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:2,
		Data:&pb.BroadCast_P{
			P:&pb.Position{
				X:p.X,
				Y:p.Y,
				Z:p.Z,
				V:p.V,
			},
		},
	}
	//将当前玩家id和位置信息发送给周边玩家
	for _,player := range players{
		player.SendMsg(200,proto_msg)
	}

	//将其他玩家告诉当前玩家  （让当前玩家能够看见周边玩家的坐标）
	//构建一个202消息  players的信息 告知当前玩家 p.send(202, ... )
	//得到全部周边玩家的player集合message Player
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//制作一个message Player 消息
		p_1 := &pb.Player{
			Pid:player.Pid,
			P:&pb.Position{
				X:player.X,
				Y:player.Y,
				Z:player.Z,
				V:player.V,
			},
		}

		fmt.Println("pid = ", player.Pid,"position=",player)

		players_proto_msg = append(players_proto_msg, p_1)
	}
	//创建一个 Message SyncPlayers
	syncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}
	//将当前的周边的全部的玩家信息 发送给当前的客户端
	p.SendMsg(202,syncPlayers_proto_msg)
}

func (p *Player)UpdatePosition(x,y,z,v float32)  {
	p.X=x
	p.Y=y
	p.Z=z
	p.V=v

	proto_msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:4,
		Data:&pb.BroadCast_P{
			P:&pb.Position{
				X:p.X,
				Y:p.Y,
				Z:p.Z,
				V:p.V,
			},
		},
	}

	players := p.GetSurroundingPlayers()

	for _,player := range players{
		player.SendMsg(200,proto_msg)
	}
}

func (p *Player) OffLine()  {
	//得到当前玩家的周边的玩家
	players := p.GetSurroundingPlayers()

	//制作一个消息MsgId:201
	proto_msg := &pb.SyncPid{
		Pid:p.Pid,
	}
	//给周边的玩家广播一个消息
	for _,player := range players {
		player.SendMsg(201,proto_msg)
	}
	//将该下线的玩家 从世界管理体验
	WorldMgrObj.RemovePlayerByPid(p.Pid)
	//将该下线玩家从地图AOIManager中移出
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid),p.X,p.Z)
}