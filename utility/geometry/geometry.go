package geometry

import (
	"math"
	"github.com/chessServer/utility/logger"
)

//package containing utility functions
type Point struct{
	X int	`json:"x"`
	Y int	`json:"y"`
}

func ConstructPoint(x,y int)(Point){
	return Point{X: x, Y: y}
}

func (p Point) CheckFieldBoundaries() bool{
	return p.X >=0 && p.X < 8 && p.Y >=0 && p.Y <8
}

func (p Point)Equal(compared Point)bool{
	return p.X == compared.X && p.Y == compared.Y
}

func (p Point)Add(value Point)(Point){
	return Point{value.X +p.X, value.Y +p.Y}
}

func (p Point) Subtract(value Point)(Point){
	return Point{p.X -value.X, p.Y -value.Y}
}

type Vector struct{
	Point
}

type Line struct{
	begin Point
	end   Point
}

func ConstructLine(begin,end Point)(Line){
	return Line{begin: begin, end: end}
}

func (l Line)abs()float64{
	p:=Point{l.end.X - l.begin.X, l.end.Y - l.begin.Y}
	return math.Sqrt(float64(p.X*p.X +p.Y*p.Y))
}

func (l Line)tan()(float64, bool){
	p:=Point{l.end.X - l.begin.X, l.end.Y - l.begin.Y}
	if p.X !=0{
		return  float64(p.Y)/float64(p.X),true
	}else{
		return 0,false
	}
}

//5 5    1 5     5 1
func (l Line) Intersect(point Point)bool{
	if l.end.X == l.begin.X {
		if l.end.X == point.X {
			vec:=Line{l.begin, point}
			if sameSign(l.end.Y- l.begin.Y,point.Y-l.begin.Y) && l.abs() >= vec.abs(){
				logger.WriteLog(6,"From line 67:Intersects: ", "Line: ", l, " Point: ",point)
				return true
			}else{
				logger.WriteLog(6,"From line 69:Doesn't intersect: ", "Line: ", l, " Point: ",point)
				return false
			}
		}else{
			logger.WriteLog(6,"From line 73: Doesn't intersect: ", "Line: ", l, " Point: ",point)
			return false
		}
	}else{
		vec:=Line{l.begin, point}
		wayTan,_:=l.tan()
		wayToObstacleTan,ok:=vec.tan()
		if !ok{
			return false
		}
		output := l.abs() >= vec.abs() && wayTan == wayToObstacleTan
		output =output && sameSign(l.end.X- l.begin.X,vec.end.X - vec.begin.X)
		output = output && sameSign(l.end.Y- l.begin.Y,vec.end.Y - vec.begin.Y)
		if output{
			logger.WriteLog(6,"From line 84:Intersects: ", "Line: ", l, " Point: ",point)
		}else{
			logger.WriteLog(6,"From line 86: Doesn't intersect: ", "Line: ", l, " Point: ",point)
		}
		return output
	}
}

func sameSign(a,b int)bool{
	return (a > 0) == (b > 0)
}



