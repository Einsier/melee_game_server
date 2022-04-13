package aoi

import "testing"

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:
 */

func TestNil(t *testing.T) {
	type S struct {
		int
	}
	var ss [10]*S
	if ss[0] != nil{
		t.Fatalf("not nil")
	}
}
