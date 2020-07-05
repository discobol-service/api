package io

import (
	"testing"
	"context"
)

func TestGetPg (t *testing.T) {
	pg := GetPg()
	rows, err := pg.Connect.Query(context.Background(), "select 1")

	if err != nil {
		t.Error(err.Error())
		return
	}

	for rows.Next() {
		var one int64
		err = rows.Scan(&one)

		if err != nil {
			t.Error(err.Error())
			return
		}

		if one != 1 {
			t.Error("Unexpected returns")
			return
		}
	}
}
