package core

import (
	"fmt"
	"testing"
)

func TestAOIManager_init(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0,250,5,0,250,5)
	fmt.Println(aoiMgr)
}

func TestAOIManagerSurround(t *testing.T)  {
	aoiMgr := NewAOIManager(0,250,5,0,250,5)
	//求出每个格子周边的九宫格信息
	for gid,_ :=range aoiMgr.grids{
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid:",gid,"grids num =",len(grids))
		//当前九宫格的ID集合
		gIDs := make([]int,0,len(grids))
		for _,grid := range grids{
			gIDs = append(gIDs,grid.GID)
		}
		fmt.Println("grids IDs are",gIDs)
	}
	fmt.Println("========>")

	playerIDs := aoiMgr.GetSurroundPIDsByPos(175,88)
	fmt.Println("palyerIDs :",playerIDs)
}
