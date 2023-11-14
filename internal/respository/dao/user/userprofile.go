package user

type Profile struct {
	ID              int64  `gorm:"primaryKey,autoIncrement"`
	Email           string `grom:"foreignkey"`
	Username        string
	Birthday        int64
	Personalprofile string `gorm:"type:text" json:"personalProfile"`
}

func (p Profile) GetTableName() string {
	return "profile"
}
