package user

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

//go:embed user.sql
var f embed.FS

func mockDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	file, err := f.ReadFile("user.sql")
	if err != nil {
		return nil, err
	}

	db = db.Exec(string(file))

	return db, nil
}

func mockRepo() (Repository, error) {
	db, err := mockDB()
	if err != nil {
		return nil, err
	}
	return NewRepository(db), nil
}

func Test_repositoryImpl_GetAll(t *testing.T) {
	a := assert.New(t)

	repo, err := mockRepo()

	users, err := repo.GetAll()
	a.NoError(err)

	var userIDs []int64
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	a.Equal(1000, len(userIDs))
}

func Test_repositoryImpl_Range(t *testing.T) {
	a := assert.New(t)

	repo, err := mockRepo()

	var userIDs []int64
	err = repo.Range(func(users []*User) (stop bool) {
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
		return false
	})
	a.NoError(err)

	a.Equal(1000, len(userIDs))
}
