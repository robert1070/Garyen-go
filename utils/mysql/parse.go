/**
 @author: robert
 @date: 2021/8/20
**/
package mysql

import (
	"fmt"
	"strings"

	"Garyen-go/tidb/parser"
	"Garyen-go/tidb/parser/ast"
	"Garyen-go/tidb/parser/format"
	_ "Garyen-go/tidb/parser/tidb-types/parser_driver"
)

type colX struct {
	colNames   []string
	tableNames []string
	schema     string
	inNum      int
}

func (v *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	if name, ok := in.(*ast.TableName); ok {
		v.tableNames = append(v.tableNames, name.Name.O)
		v.schema = name.Schema.O
	}

	if patterIn, ok := in.(*ast.PatternInExpr); ok {
		v.inNum = len(patterIn.List)
		if v.inNum > 10 {
			patterIn.List = patterIn.List[:10]

			patterIn.List = append(patterIn.List, ast.NewValueExpr(fmt.Sprintf("/* %d */", v.inNum)))
		}
	}
	return in, false
}

func (v *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extractToColX1(node *ast.StmtNode) colX {
	v := new(colX)
	(*node).Accept(v)
	return *v
}

// return sql type
type Type string

const (
	Truncate     Type = "truncate"
	DropSchema   Type = "dropSchema"
	DropTable    Type = "dropTable"
	AlterTable   Type = "alterTable"
	RenameTable  Type = "renameTable"
	CreateTable  Type = "createTable"
	CreateIndex  Type = "createIndex"
	CreateSchema Type = "createSchema"
	Insert       Type = "insert"
	Update       Type = "update"
	Delete       Type = "delete"
	Select       Type = "select"
)

func (t Type) String() string {
	return string(t)
}

func GetTypeFromStmt(stmt ast.StmtNode) Type {
	switch stmt.(type) {
	case *ast.TruncateTableStmt:
		return Truncate
	case *ast.DropDatabaseStmt:
		return DropSchema
	case *ast.DropTableStmt:
		return DropTable
	case *ast.AlterTableStmt:
		return AlterTable
	case *ast.RenameTableStmt:
		return RenameTable
	case *ast.CreateTableStmt:
		return CreateTable
	case *ast.CreateIndexStmt:
		return CreateIndex
	case *ast.CreateDatabaseStmt:
		return CreateSchema
	case *ast.InsertStmt:
		return Insert
	case *ast.UpdateStmt:
		return Update
	case *ast.DeleteStmt:
		return Delete
	case *ast.SelectStmt:
		return Select
	default:
		return "unknown type"
	}
}

func parseStmtNode(SQL string) (stmts []ast.StmtNode, err error) {
	stmts, _, err = parser.New().Parse(SQL, "", "")
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

// generate sql
func GenerateSQL(stmt ast.StmtNode) (string, error) {
	builder := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, builder)
	if err := stmt.Restore(ctx); err != nil {
		return "", err
	}

	return builder.String(), nil
}

type SQLExtract struct {
	Schema  string   `json:"schema"`
	Table   []string `json:"table"`
	SQLType string   `json:"sql_type"`
	SQLStmt string   `json:"sql_stmt"`
}

// AcquireST return schema and table
func AcquireST(SQL string) ([]*SQLExtract, error) {
	var (
		err   error
		stmts []ast.StmtNode
	)

	extract := make([]*SQLExtract, 0)

	stmts, err = parseStmtNode(SQL)
	if err != nil {
		return nil, err
	}

	for _, rst := range stmts {
		SQLType := GetTypeFromStmt(rst)
		prop := extractToColX1(&rst)

		extract = append(extract, &SQLExtract{
			Schema:  prop.schema,
			Table:   prop.tableNames,
			SQLType: SQLType.String(),
			SQLStmt: SQL,
		})
	}
	return extract, nil
}
