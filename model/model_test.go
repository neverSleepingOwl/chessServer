package model

import (
	"testing"
	"fmt"
	"chessServer/utility"
)

func TestConstructBishop(t *testing.T) {
	fmt.Println("echo")
}

func Render(p []utility.Point, fig utility.Point)(bool){	//	function, iterating through array
	for i:=0;i < 8;i++{
		for j:=0;j<8;j++{
			for _,element:=range p{
				if !element.CheckFieldBoundaries(){
					return false
				}
				if element.Equal(fig){
					return false
				}
				if comp := utility.ConstructPoint(i,j);comp.Equal(element){
					fmt.Print("*")
				}else if comp.Equal(fig){
					fmt.Print("X")
				}else{
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
	return true
}
func TestLinearFigure_AttacksAvailable(t *testing.T) {
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructBishop(i,j,0)
			p:=Render(temp.AttacksAvailable(),temp.Point)
			if !p{
				t.Error("Wrong Bishop fields")
			}
		}
	}
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructRook(i,j,0)
			p:=Render(temp.AttacksAvailable(),temp.Point)
			if !p{
				t.Error("Wrong Rook fields")
			}
		}
	}
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructQueen(i,j,0)
			p:=Render(temp.AttacksAvailable(), temp.Point)
			if !p{
				t.Error("Wrong Queen fields")
			}
		}
	}
}
