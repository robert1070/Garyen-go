package mysql

import (
	"fmt"
	"testing"
)

func TestAcquireST(t *testing.T) {
	selectStmt := `select * from t1.t1 a join t1.t2 b on a.id = b.id where a.user_id = 1;
	select * from t2.t2;`
	extract, err := AcquireST(selectStmt)
	if err != nil {
		t.Logf("parse err, err:%s", err)
	}

	for _, v := range extract {
		fmt.Printf("schema:%s, table:%s, SQLtype:%s\n", v.Schema, v.Table, v.SQLType)
	}
}
