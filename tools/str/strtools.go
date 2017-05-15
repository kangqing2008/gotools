package str

import (
	"fmt"
	"strconv"
)

//对float64进行修约
func RF64(value float64,precision int)float64{
	s := fmt.Sprintf("%0." + strconv.Itoa(precision) + "f", value)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

//对float64进行修约,并返回字符串
func RF64Str(value float64,precision int)string{
	s := fmt.Sprintf("%0." + strconv.Itoa(precision) + "f", value)
	return s
}