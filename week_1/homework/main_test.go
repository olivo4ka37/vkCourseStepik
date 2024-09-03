package main

import (
	"bytes"
	"fmt"
	"testing"
)

const testFullResult = `├───project
│	├───file.txt (19b)
│	└───gopher.png (70372b)
├───static
│	├───a_lorem
│	│	├───dolor.txt (empty)
│	│	├───gopher.png (70372b)
│	│	└───ipsum
│	│		└───gopher.png (70372b)
│	├───css
│	│	└───body.css (28b)
│	├───empty.txt (empty)
│	├───html
│	│	└───index.html (57b)
│	├───js
│	│	└───site.js (10b)
│	└───z_lorem
│		├───dolor.txt (empty)
│		├───gopher.png (70372b)
│		└───ipsum
│			└───gopher.png (70372b)
├───zline
│	├───empty.txt (empty)
│	└───lorem
│		├───dolor.txt (empty)
│		├───gopher.png (70372b)
│		└───ipsum
│			└───gopher.png (70372b)
└───zzfile.txt (empty)
`

func TestTreeFull(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, "testdata", true)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	fmt.Println(compareStrings(result, testFullResult))
	if result != testFullResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testFullResult)
	}
}

const testDirResult = `├───project
├───static
│	├───a_lorem
│	│	└───ipsum
│	├───css
│	├───html
│	├───js
│	└───z_lorem
│		└───ipsum
└───zline
	└───lorem
		└───ipsum
`

func TestTreeDir(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, "testdata", false)
	if err != nil {
		t.Errorf("test for OK Failed - error")
	}
	result := out.String()
	fmt.Println(compareStrings(result, testDirResult))
	if result != testDirResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testDirResult)
	}
}

func compareStrings(str1, str2 string) string {
	length := len(str1)
	if len(str2) < length {
		length = len(str2)
	}

	for i := 0; i < length; i++ {
		if str1[i] != str2[i] {
			return fmt.Sprintf("Difference at position %d: '%c' vs '%c'", i, str1[i], str2[i])
		}
	}

	if len(str1) != len(str2) {
		return fmt.Sprintf("Strings differ in length: %d vs %d", len(str1), len(str2))
	}

	return "Strings are identical"
}
