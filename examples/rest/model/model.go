package model

type UserModel struct {
    ID         int    `gorm:"column:id;type:int(11) auto_increment;primary_key;comment:'主键'" json:"id"`
    UserName   string `gorm:"column:name;type:varchar(100);not null;comment:'用户名'" json:"user_name"`
    Password   string `gorm:"column:name;type:varchar(100);not null;comment:'密码'" json:"password"`
    Sex        int    `gorm:"column:name;type:int(11);not null;default 0;comment:'性别'" json:"sex"`
    CreateTime int    `gorm:"column:name;type:timestamp;not null;comment:'创建时间'" json:"create_time"`
}

func (m *UserModel) TableName() string {
    return "user"
}
