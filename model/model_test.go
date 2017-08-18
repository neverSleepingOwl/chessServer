package model

import (
	"testing"
	"fmt"
	"chessServer/utility"
)



func Render(p []utility.Point, fig utility.Point)(bool){
	emptyField:=true
	for i:=0;i < 8;i++{
		for j:=0;j<8;j++{
			emptyField = true
			comp := utility.ConstructPoint(i,j)
			for _,element:=range p{
				if !element.CheckFieldBoundaries(){
					return false
				}
				if element.Equal(fig){
					return false
				}
				if comp.Equal(element){
					fmt.Print("*")
					emptyField = false
					break
				}
			}
			if comp.Equal(fig){
				fmt.Print("H")
			}else if emptyField{
				fmt.Print("X")
			}

		}
		fmt.Println()
	}
	fmt.Println("________")
	return true
}

func Check(p []utility.Point, fig utility.Point)(bool){
	for _,element:=range p{
		if element.Equal(fig){
			return false
		}
		if !element.CheckFieldBoundaries(){
			return false
		}
	}
	return true
}
func TestLinearFigure_AttacksAvailable(t *testing.T) {
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructBishop(i,j,0)
			p:=Check(temp.AttacksAvailable(),temp.Point)
			if !p{
				t.Error("Wrong Bishop fields")
			}
			fmt.Println()
		}
	}
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructRook(i,j,0)
			p:=Check(temp.AttacksAvailable(),temp.Point)
			if !p{
				t.Error("Wrong Rook fields")
			}
			fmt.Println()
		}
	}
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructQueen(i,j,0)
			p:=Check(temp.AttacksAvailable(), temp.Point)
			if !p{
				t.Error("Wrong Queen fields")
			}
			fmt.Print()
		}
	}
}

func TestNonLinearFigure_StepsAvailable(t *testing.T) {
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructKnight(i,j,0)
			p:=Check(temp.AttacksAvailable(), temp.Point)
			if !p{
				t.Error("Wrong Knight fields")
			}
			fmt.Print()
		}
	}
	for i:=0;i < 8;i++{
		for j:=0;j < 8;j++{
			temp:=ConstructKing(i,j,0)
			p:=Check(temp.AttacksAvailable(), temp.Point)
			if !p{
				t.Error("Wrong King fields")
			}
			fmt.Print()
		}
	}
}

func TestPawn_AttacksAvailable(t *testing.T) {

}

func TestPawn_StepsAvailable(t *testing.T) {

}
