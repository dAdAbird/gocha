package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCommandsParse(t *testing.T) {
	cmds := make(ResponsesSet)
	cmdTest := func(_ *Session, _ []byte) error {
		return nil
	}
	cmds.Register("TEST", cmdTest)

	cases := []struct {
		data       []byte
		shouldCmd  RespHandler
		shouldArgs []byte
		shouldErr  error
	}{
		{
			data:       []byte("TEST test msg\n"),
			shouldCmd:  cmdTest,
			shouldArgs: []byte("test msg"),
			shouldErr:  nil,
		},
		{
			data:       []byte("TEST\n"),
			shouldCmd:  cmdTest,
			shouldArgs: nil,
			shouldErr:  nil,
		},
		{
			data:       []byte("TEST"),
			shouldCmd:  cmdTest,
			shouldArgs: nil,
			shouldErr:  nil,
		},
		{
			data:       []byte("NOCMD\n"),
			shouldCmd:  nil,
			shouldArgs: nil,
			shouldErr:  ErrNoCommand,
		},
		{
			data:       []byte("nothingtoshow"),
			shouldCmd:  nil,
			shouldArgs: nil,
			shouldErr:  ErrNoCommand,
		},
		{
			data:       []byte("\n"),
			shouldCmd:  nil,
			shouldArgs: nil,
			shouldErr:  ErrNoCommand,
		},
		{
			data:       []byte(""),
			shouldCmd:  nil,
			shouldArgs: nil,
			shouldErr:  ErrNoCommand,
		},
		{
			data:       []byte("TEST test msg\r\n"),
			shouldCmd:  cmdTest,
			shouldArgs: []byte("test msg"),
			shouldErr:  nil,
		},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%s", testCase.data), func(t *testing.T) {
			cmd, args, err := cmds.Parse(testCase.data)
			if &cmd != &testCase.shouldCmd &&
				!bytes.Equal(args, testCase.shouldArgs) ||
				err != testCase.shouldErr {
				t.Errorf("Wrong result. Got <%v>, <%s>, <%v>, should be %v", cmd, args, err, testCase)
			}
		})
	}
}
