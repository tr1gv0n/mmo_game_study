package core

import (
	"fmt"
	"testing"
)

func TestGrid(t *testing.T)  {
	player1 := "玩家1"
	player2 := "玩家2"

	//单元测试Grid模块
	g := NewGrid(1,2,3,4,5)

	g.Add(1,player1)
	g.Add(2,player2)
	//打印格子信息
	fmt.Println(g)
}
