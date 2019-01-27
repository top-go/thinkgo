package memory

import (
	"errors"
	"testing"
	"time"
)

type CacheUser struct {
	Name string
	Age  int
}


func getStore() *Store {
	return NewStore("cache")
}

func TestStore(t *testing.T) {
	s := getStore()
	var a int
	var b string
	var c CacheUser

	err := s.Get("a", &a)
	if err == nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	err = s.Get("b", &b)
	if err == nil{
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	s.Put("a", 1, 10*time.Minute)
	s.Put("b", "thinkgo", 10*time.Minute)

	err = s.Get("a", &a)
	if err != nil {
		t.Error(err)
	}

	if a != 1{
		t.Error("Expect: ", 1)
	}

	err = s.Get("b", &b)
	if err != nil {
		t.Error(err)
	}

	if b != "thinkgo"{
		t.Error("Expect: ", "thinkgo")
	}

	err = s.Put(
		"user", CacheUser{
			Name: "alice",
			Age:  16,
		},
		10*time.Minute,
	)
	if err != nil {
		t.Error(err)
	}

	err = s.Get("user", &c)
	if err != nil {
		t.Error(err)
	}

	t.Logf("user:name=%s,age=%d", c.Name, c.Age)
}

func TestStoreDuration(t *testing.T)  {
	s := getStore()
	var a int

	s.Put("a", 3, 20*time.Millisecond)

	<-time.After(21 * time.Millisecond)
	err := s.Get("a", &a)
	if err == nil {
		t.Error("Found a when it should have been automatically deleted")
	}
}


func TestStoreForgetAndExist(t *testing.T) {
	s := getStore()
	err := s.Put("forget", "Forget me", 10*time.Minute)
	if err != nil {
		t.Error(err)
	}

	exist := s.Exist("forget")
	if exist != true {
		t.Error(errors.New("Expect true"))
	}

	err = s.Forget("forget")
	if err != nil {
		t.Error(err)
	}

	exist = s.Exist("forget")
	if exist == true {
		t.Error(errors.New("Expect false"))
	}
}

func TestStoreFlush(t *testing.T) {
	s := getStore()
	err := s.Put("Flush", "Flush all", 10*time.Minute)
	if err != nil {
		t.Error(err)
	}

	exist := s.Exist("Flush")
	if exist != true {
		t.Error(errors.New("Expect true"))
	}

	err = s.Flush()
	if err != nil {
		t.Error(err)
	}

	exist = s.Exist("Flush")
	if exist == true {
		t.Error(errors.New("Expect false"))
	}
}