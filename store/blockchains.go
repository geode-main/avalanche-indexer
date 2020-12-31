package store

import "gorm.io/gorm"

type BlockchainsStore struct {
	*gorm.DB
}
