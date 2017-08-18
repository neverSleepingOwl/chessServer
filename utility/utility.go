package utility

import "math"

//package containing utility functions
type Point struct{
	X int
	Y int
}

func (p Point) CheckFieldBoundaries() bool{
	return p.X >=0 && p.X < 8 && p.Y >=0 && p.Y <8
}

func (p Point)Equal(compared Point)bool{
	return p.X == compared.X &&p.Y == compared.Y
}

func (p Point)Add(value Point)(Point){
	return Point{value.X+p.X, value.Y+p.Y}
}

func (p Point)Substract(value Point)(Point){
	return Point{p.X-value.X, p.Y-value.Y}
}

type Vector struct{
	Point
}

type Line struct{
	Begin Point
	End   Point
}

func (l Line)abs()float64{
	p:=Point{l.End.X - l.Begin.X, l.End.Y - l.Begin.Y}
	return math.Sqrt(float64(p.X*p.X+p.Y*p.Y))
}

func (l Line)tan()(float64, bool){
	p:=Point{l.End.X - l.Begin.X, l.End.Y - l.Begin.Y}
	if p.X !=0{
		return  float64(p.Y)/float64(p.X),true
	}else{
		return 0,false
	}
}

func (l Line) Intersect(point Point)bool{
	if l.End.X == l.Begin.X{
		if l.End.X == point.X{
			vec:=Line{l.Begin, point}
			if sameSign(l.End.Y - l.Begin.Y,point.Y-l.Begin.Y) && l.abs() >= vec.abs(){
				return true
			}else{
				return false
			}
		}else{
			return false
		}
	}else{
		vec:=Line{l.Begin, point}
		return l.abs() > vec.abs() && l.tan() == vec.tan()
	}
}

func sameSign(a,b int)bool{
	return (a > 0) == (b > 0)
}



