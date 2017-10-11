package geometry

import(
	"testing"
)

func TestPoint_CheckFieldBoundaries(t *testing.T) {
		out:=[]Point{{0,-1},{-1,0}, {0,8}, {8,0}}
		for _,elem:=range out{
			if elem.CheckFieldBoundaries(){
				t.Error("Incorrect Bounds",elem)
			}
		}
		for i:=0;i<8;i++{
			for j:=0;j<8;j++{
				if ok:=ConstructPoint(i,j);!ok.CheckFieldBoundaries(){
					t.Error("Incorrect Bounds",ok)
				}

			}
		}
}

/*func TestPoint_Add(t *testing.T) {
	a,b:=make([]Point,1000000), make([]Point,1000000)
	for i:=0;i<1000;i++{
		for j:=0;j<1000;j++{
			a[i*1000+j]=ConstructPoint(i,j)
			b[i*1000+j]=ConstructPoint(i,j)
		}
	}
	for i,_:=range a{
		if a[i].Add(b[i]) != ConstructPoint(a[i].X+b[i].X, a[i].Y+b[i].Y){
			t.Error("Failed summ", a[i])
		}
	}
}*/

func TestLine_Intersect(t *testing.T) {
	l := Line{Point{7,7},Point{0,0}}
	p := Point{2,3}
	if l.Intersect(p){
		t.Error("Incorrect intersect")
	}
}

func TestPoint_Equal(t *testing.T) {

}

func TestPoint_Substract(t *testing.T) {

}