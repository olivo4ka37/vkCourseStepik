package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var test1 = `1
2
3
4
5
`

var test1Result = `1
2
3
4
5
`

var test2Fail = `
1
3
3
7
1
`

func TestUniq(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(test1))
	out := new(bytes.Buffer)

	err := uniq(in, out)
	if err != nil {
		t.Errorf("test func uniq is failed")
	}

	result := out.String()
	if result != test1Result {
		t.Errorf("result of test is not matched\n Result of test:%v \n Correct result:%v", result, test1Result)
	}
}

func TestForError(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(test2Fail))
	out := new(bytes.Buffer)

	err := uniq(in, out)
	if err == nil {
		result := out.String()
		t.Errorf("Test do not return error, when it should test is: \n%v, result is: \n%v", test2Fail, result)
	}
}
