/**
 @author: robert
 @date: 2021/9/9
**/
package lib

import (
	"fmt"
	"testing"
)

func TestDjangoEncrypt(t *testing.T) {
	s1 := "test01"
	s2 := "test02"

	newStr := DjangoEncrypt(s1, s2)
	fmt.Println(newStr)
}
