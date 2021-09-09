/**
 @author: robert
 @date: 2021/9/9
**/
package lib

import (
	"fmt"
	"testing"
)

func TestJwtAuth(t *testing.T) {

	st, err := JwtAuth("robert", "develop")
	if err != nil {
		t.Logf("err:%s", err)
	}

	fmt.Println(st)
}
