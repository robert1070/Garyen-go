package dms

type OpType int8

const (
	DataChange   OpType = iota + 1 // 数据变更
	StructChange                   // 结构变更
	DataExport                     // 数据导出
	PermApply                      // 权限申请
	Query                          // 数据查询
)

type ProType int8

const (
	Submitted ProType = iota + 1
	Reviewed
	Refused
	Executed
	Executing
	Finished
	Defeated
	Revoked
	Schedule
)

type EType int8

const (
	PROD EType = iota + 1
	GRAY
	TEST
	DEV
)

type CoreSQLOrder struct {
	ID           int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	OrderId      string  `gorm:"type:varchar(32);not null;unique_index:uk_order_id;" json:"order_id"`
	Title        string  `gorm:"type:varchar(64);not null;default '';" json:"title"`
	Demand       string  `gorm:"type:varchar(128);not null;default '';" json:"demand"`
	SQLType      OpType  `gorm:"type:tinyint(4) unsigned;not null;default 1;" json:"sql_type"`
	Sponsor      string  `gorm:"type:varchar(12);not null;" json:"sponsor"`
	Reviewer     string  `gorm:"type:varchar(12);not null;default '';" json:"reviewer"`
	Executor     string  `gorm:"type:varchar(12);not null;default '';" json:"executor"`
	Env          EType   `gorm:"type:tinyint(4) unsigned;not null;" json:"env"`
	Delay        string  `gorm:"varchar(24);not null;default '';" json:"delay"`
	IsBackup     uint    `gorm:"type:type:tinyint(4);unsigned;not null;default 0;" json:"is_backup"`
	DBId         uint    `gorm:"type:int(11);not null;default '0';" json:"db_id"`
	DataSourceId uint    `gorm:"type:int(11);not null;" json:"data_source_id"`
	RemoteHost   string  `gorm:"type:varchar(24);not null;" json:"remote_host"`
	RemotePort   uint    `gorm:"type:int(11);not null;" json:"remote_port"`
	DBName       string  `gorm:"type:varchar(24);not null;" json:"db_name"`
	Progress     ProType `gorm:"type:tinyint(4);not null default 1;" json:"progress"`
	Step         uint    `gorm:"type:tinyint(4);not null default 1;" json:"step"`
	Remark       string  `gorm:"type:varchar(24);not null;default '';" json:"remark"`
	PermDetail   string  `gorm:"type:text;;" json:"perm_detail"`
	Contents     string  `gorm:"type:text;" json:"contents"`
	FileFormat   string  `gorm:"type:char(4);not null;default 'xlsx';" json:"file_format"`
	GmtCreate    int64   `gorm:"type:int(11);not null;" json:"gmt_create"`
	GmtModified  int64   `gorm:"type:int(11);not null;" json:"gmt_modified"`
}

func (c *CoreSQLOrder) TableName() string {
	return "core_sql_order"
}
