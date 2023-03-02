package service

import (
	"MyLNPU/internal/cache"
	"MyLNPU/internal/db"
	"MyLNPU/internal/model"
	"errors"
	"github.com/redis/go-redis/v9"
)

func GetSystemNotice() (string, error) {
	notice, err := cache.Get("lnpu:notice")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			//缓存为空则从mysql获取
			n, err := db.GetSystemNotice()
			if err != nil {
				return "", err
			}
			cache.Set("lnpu:notice", n.Notice, 0)
			return n.Notice, nil
		}
		return "", err
	}
	return notice, nil
}

func UpdateSystemNotice(notice *model.Notice) error {
	err := cache.Del("lnpu:notice")
	if err != nil {
		return err
	}
	err = db.UpdateSystemNotice(notice)
	if err != nil {
		return err
	}
	return nil
}
