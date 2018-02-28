package models

type RefundInfo struct {
	ID        int64 `xorm:"pk autoincr" json:"id"`
	CreatedAt int64 `xorm:"not null default 0 int index" json:"createdAt"`
	UpdatedAt int64 `xorm:"not null default 0 int index" json:"-"`
}
