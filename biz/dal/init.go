package dal

import "github.com/2451965602/LMS/biz/dal/db"

func Init() error {
	err := db.Init()
	if err != nil {
		return err
	}

	return nil
}
