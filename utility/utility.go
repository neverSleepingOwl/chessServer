package utility
//some useful staff


func Delete(a[]interface{},at int){
	a = append(a[:at], a[at+1:]...)
}