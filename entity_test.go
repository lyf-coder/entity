// Copyright © 2020 - present. liyongfei <liyongfei@walktotop.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package entity

import (
	"io/ioutil"
	"testing"
)

func Test_NewByJson(t *testing.T) {
	f, err := ioutil.ReadFile("test_data.json")
	if err != nil {
		t.Fatal("read fail", err)
	}

	entity := NewByJson(f)

	if !entity.GetBool("event:simulator"){
		t.Errorf("GetBool 'event:simulator' val is not true")
	}


	clientContext := New(entity.GetStringMapSlice("clientContext")[1])
	offsetInMilliseconds := clientContext.GetInt("payload:offsetInMilliseconds")

	if offsetInMilliseconds != 1023785 {
		t.Errorf("GetInt 'payload:offsetInMilliseconds' val is not 1023785")
	}
}
