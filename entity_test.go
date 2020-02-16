// Copyright © 2020 - present. liyongfei <liyongfei@walktotop.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package entity

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestEntity_Get(t *testing.T) {
	e := new(Entity)
	e.Set("name", "jack")
	log.Println()
	if e.Get("name") != "jack" {
		t.Error("Get name value is not jack")
	}
}

func TestEntity_GetData(t *testing.T) {
	e := new(Entity)
	e.Set("name", "jack")
	if e.GetData()["name"] != "jack" {
		t.Error("GetData func is error")
	}
}

func Test_NewByJson(t *testing.T) {
	f, err := ioutil.ReadFile("test_data.json")
	if err != nil {
		t.Fatal("read fail", err)
	}

	entity := NewByJSON(f)

	if !entity.GetBool("event:simulator") {
		t.Errorf("GetBool 'event:simulator' val is not true")
	}

	clientContext := New(entity.GetStringMapSlice("clientContext")[1])
	offsetInMilliseconds := clientContext.GetInt("payload:offsetInMilliseconds")

	if offsetInMilliseconds != 1023785 {
		t.Errorf("GetInt 'payload:offsetInMilliseconds' val is not 1023785")
	}
}
