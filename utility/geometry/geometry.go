package geometry

import "math"

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
	return p.X == compared.X &&p.Y == compared.Y
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

func (l Line) Intersect(point Point)bool{
	if l.end.X == l.begin.X {
		if l.end.X == point.X {
			vec:=Line{l.begin, point}
			if sameSign(l.end.Y- l.begin.Y,point.Y-l.begin.Y) && l.abs() >= vec.abs(){
				return true
			}else{
				return false
			}
		}else{
			return false
		}
	}else{
		vec:=Line{l.begin, point}
		wayTan,_:=l.tan()
		wayToObstacleTan,_:=vec.tan()
		return l.abs() > vec.abs() && wayTan == wayToObstacleTan
	}
}

func sameSign(a,b int)bool{
	return (a > 0) == (b > 0)
}



