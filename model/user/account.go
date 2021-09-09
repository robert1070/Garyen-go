/**
 @author: robert
 @date: 2021/8/19
**/
package user

type CoreAccount struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Username    string `gorm:"type:varchar(50);not null;index:idx_username" json:"username"`
	Password    string `gorm:"type:varchar(150);not null" json:"password"`
	Rule        string `gorm:"type:varchar(10);not null" json:"rule"`
	RealName    string `gorm:"type:varchar(50);" json:"real_name"`
	Email       string `gorm:"type:varchar(50);" json:"email"`
	State       int8   `gorm:"type:tinyint;not null;default '1';" json:"email"` // 0-正常，1-启用
	GmtCreate   int64  `gorm:"type:int(11);not null;" json:"gmt_create"`
	GmtModified int64  `gorm:"type:int(11);not null;" json:"gmt_modified"`
}

func (c *CoreAccount) TableName() string {
	return "core_account"
}
