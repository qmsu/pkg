package utils

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestCopyWithFormatTime(t *testing.T) {
	type User struct {
		Name    string
		Role    string
		Age     int
		AddTime time.Time
	}

	type ModelUser struct {
		Name    string
		Role    string
		Age     int
		AddTime string
	}

	user := User{
		Name:    "lily",
		Role:    "admin",
		Age:     10,
		AddTime: time.Now(),
	}
	var mUser ModelUser
	err := CopyWithFormatTime(&mUser, &user)
	if err != nil {
		t.Fatal(err.Error())
	}
	b, _ := json.Marshal(mUser)
	// mUser: {"Name":"lily","Role":"admin","Age":10,"AddTime":"2023-07-10 10:49:17"}
	fmt.Println("mUser:", string(b))
}

func BenchmarkCopyWithFormatTime(t *testing.B) {
	type User struct {
		Name    string
		Role    string
		Age     int
		AddTime time.Time
	}

	type ModelUser struct {
		Name    string
		Role    string
		Age     int
		AddTime string
	}

	user := User{
		Name:    "lily",
		Role:    "admin",
		Age:     10,
		AddTime: time.Now(),
	}
	var mUser ModelUser
	err := CopyWithFormatTime(&mUser, &user)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("mUser:", mUser)
	// BenchmarkCopyWithFormatTime-8   	1000000000	         0.0000426 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkCopy(t *testing.B) {
	type User struct {
		Name    string
		Role    string
		Age     int
		AddTime time.Time
	}

	type ModelUser struct {
		Name    string
		Role    string
		Age     int
		AddTime string
	}

	user := User{
		Name:    "lily",
		Role:    "admin",
		Age:     10,
		AddTime: time.Now(),
	}
	var mUser ModelUser
	mUser = ModelUser{
		Name:    user.Name,
		Role:    user.Role,
		Age:     user.Age,
		AddTime: user.AddTime.Format("2006-01-02 15:04:05"),
	}
	fmt.Println("mUser:", mUser)
	// BenchmarkCopy-8   	1000000000	         0.0000246 ns/op	       0 B/op	       0 allocs/op
}
