package user

import "gorm.io/gorm"

const tableName = "user"

const bufSize = 256

type Repository interface {
	GetAll() ([]*User, error)
	Range(func([]*User) (stop bool)) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db: db,
	}
}

type repositoryImpl struct {
	db *gorm.DB
}

func (r repositoryImpl) GetAll() ([]*User, error) {
	var users []*User
	err := r.db.
		Table(tableName).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r repositoryImpl) Range(f func([]*User) (stop bool)) error {
	currentMaxID := int64(0)

	for {
		var users []*User

		err := r.db.
			Table(tableName).
			Where("id > ?", currentMaxID).
			Order("id").
			Limit(bufSize).
			Scan(&users).Error
		if err != nil {
			return err
		}

		if f(users) {
			break
		}

		if len(users) < bufSize { // 最后一页
			break
		}

		currentMaxID = users[bufSize-1].ID
	}

	return nil
}
