package utils

import (
	"log"
	"testing"
)

/**
*@Author Sly
*@Date 2022/1/15
*@Version 1.0
*@Description:
 */

func TestFloatEqual(t *testing.T) {
	if FloatEqual(0, -5.1259995e-06) == false {
		log.Fatalln("wrong!")
	}
}
