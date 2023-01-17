package user

import (
	"embed"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

//go:embed user.sql
var f embed.FS

func readSQL() (string, error) {
	file, err := f.ReadFile("user.sql")
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func mockDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sql, err := readSQL()
	if err != nil {
		return nil, err
	}

	db = db.Exec(sql)

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
	Convey("获取全量数据", t, func() {
		repo, err := mockRepo()
		So(err, ShouldBeNil)

		users, err := repo.GetAll()
		So(err, ShouldBeNil)

		var userIDs []int64
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}

		So(len(userIDs), ShouldEqual, 1_000)
	})
}

func Test_repositoryImpl_Range(t *testing.T) {
	Convey("使用 Range 获取部分数据", t, func() {
		repo, err := mockRepo()
		So(err, ShouldBeNil)

		Convey("获取全部的 ID", func() {
			var userIDs []int64
			err = repo.Range(func(users []*User) (stop bool) {
				for _, user := range users {
					userIDs = append(userIDs, user.ID)
				}
				return false
			})

			So(err, ShouldBeNil)
			So(len(userIDs), ShouldEqual, 1_000)
		})

		Convey("获取 100 条 ID", func() {
			var userIDs []int64
			count := 0
			err = repo.Range(func(users []*User) (stop bool) {
				for _, user := range users {
					userIDs = append(userIDs, user.ID)
					count++
					if count == 100 {
						return true
					}
				}
				return false
			})

			So(err, ShouldBeNil)
			So(len(userIDs), ShouldEqual, 100)
		})
	})
}
