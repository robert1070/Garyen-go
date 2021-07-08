// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"errors"
	"fmt"
)

// SQLError records an error information, from executing SQL.
type SQLError struct {
	Code    uint16
	Message string
	State   string
}

// Error prints errors, with a formatted string.
func (e *SQLError) Error() string {
	return fmt.Sprintf("ERROR %d (%s): %s", e.Code, e.State, e.Message)
}

// NewErr generates a SQL error, with an error code and default format specifier defined in ErrName.
func NewErr(errCode uint16, args ...interface{}) *SQLError {
	e := &SQLError{Code: errCode}

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DefaultMySQLState
	}

	if format, ok := ErrName[errCode]; ok {
		e.Message = fmt.Sprintf(format, args...)
	} else {
		e.Message = fmt.Sprint(args...)
	}

	return e
}

// NewErrf creates a SQL error, with an error code and a format specifier.
func NewErrf(errCode uint16, format string, args ...interface{}) *SQLError {
	e := &SQLError{Code: errCode}

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DefaultMySQLState
	}

	e.Message = fmt.Sprintf(format, args...)

	return e
}

var (
	ErrNoMasterConn     = errors.New("no master connection")
	ErrNoSlaveConn      = errors.New("no slave connection")
	ErrNoDefaultNode    = errors.New("no default node")
	ErrNoMasterDB       = errors.New("no master database")
	ErrNoSlaveDB        = errors.New("no slave database")
	ErrNoDatabase       = errors.New("no database")
	ErrMasterDown       = errors.New("master is down")
	ErrSlaveDown        = errors.New("slave is down")
	ErrDatabaseClose    = errors.New("database is close")
	ErrConnIsNil        = errors.New("connection is nil")
	ErrBadConn          = errors.New("connection was bad")
	ErrMalformPacket    = errors.New("malform packet error")
	ErrIgnoreSQL        = errors.New("ignore this sql")
	ErrAddressNull      = errors.New("address is nil")
	ErrInvalidArgument  = errors.New("argument is invalid")
	ErrInvalidCharset   = errors.New("charset is invalid")
	ErrCmdUnsupported   = errors.New("command unsupported")
	ErrLocationsCount   = errors.New("locations count is not equal")
	ErrNoCriteria       = errors.New("plan have no criteria")
	ErrNoRouteNode      = errors.New("no route node")
	ErrResultNil        = errors.New("result is nil")
	ErrSumColumnType    = errors.New("sum column type error")
	ErrSelectInInsert   = errors.New("select in insert not allowed")
	ErrInsertInMulti    = errors.New("insert in multi node")
	ErrUpdateInMulti    = errors.New("update in multi node")
	ErrDeleteInMulti    = errors.New("delete in multi node")
	ErrReplaceInMulti   = errors.New("replace in multi node")
	ErrExecInMulti      = errors.New("exec in multi node")
	ErrTransInMulti     = errors.New("transaction in multi node")
	ErrNoPlan           = errors.New("statement have no plan")
	ErrNoPlanRule       = errors.New("statement have no plan rule")
	ErrUpdateKey        = errors.New("routing key in update expression")
	ErrStmtConvert      = errors.New("statement fail to convert")
	ErrExprConvert      = errors.New("expr fail to convert")
	ErrConnNotEqual     = errors.New("the length of conns not equal sqls")
	ErrKeyOutOfRange    = errors.New("shard key not in key range")
	ErrMultiShard       = errors.New("insert or replace has multiple shard targets")
	ErrIRNoColumns      = errors.New("insert or replace must specify columns")
	ErrIRNoShardingKey  = errors.New("insert or replace not contain sharding key")
	ErrColsLenNotMatch  = errors.New("insert or replace cols and values length not match")
	ErrDateIllegal      = errors.New("date format illegal")
	ErrDateRangeIllegal = errors.New("date range format illegal")
	ErrDateRangeCount   = errors.New("date range count is not equal")
	ErrSlaveExist       = errors.New("slave has exist")
	ErrSlaveNotExist    = errors.New("slave has not exist")
	ErrBlackSQLExist    = errors.New("black sql has exist")
	ErrBlackSQLNotExist = errors.New("black sql has not exist")
	ErrInsertTooComplex = errors.New("insert is too complex")
	ErrSQLNULL          = errors.New("sql is null")
	ErrNoWhere          = errors.New("no where condition")
	ErrInvalidSQLMode   = errors.New("invalid sql mode")
	ErrInternalServer   = errors.New("internal server error")
	ErrAccessDenied     = errors.New("access denied")
	ErrNoShardKey       = errors.New("no shared key or shared key in 'or' block")
	ErrWhereLimitEmpty  = errors.New("where limit is empty")
	ErrShardOrderLimit  = errors.New("unsupport query order or limit in multi table")
	ErrNoDBSelected     = errors.New("no database selected")
)

var (
	ErrStrMaxInsertValues = "max insert values %d, large than %d"
	ErrStrMaxWhereIn      = "max where in %d, large than %d"
	ErrStrMaxWhereLimit   = "max where limit %d, large than %d"
	ErrStrMaxPacketSize   = "max packet size %d, large than %d"
)
