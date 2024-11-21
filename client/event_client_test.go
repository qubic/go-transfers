package client

import "testing"

func Test_logEvents(t *testing.T) {

	err := logEvents()
	if err != nil {
		t.Error(err)
	}

}
