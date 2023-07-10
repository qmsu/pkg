package utils

import (
	"errors"
	"time"

	"github.com/jinzhu/copier"
)

// 复制对象
// 深度拷贝
// 将时间格式化为 string 类型
// toValue 目标对象指针
// fromValue 原对象指针
func CopyWithFormatTime(toValue interface{}, fromValue interface{}) (err error) {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type not matching")
					}
					return s.Format("2006-01-02 15:04:05"), nil
				},
			},
		},
	})
}
