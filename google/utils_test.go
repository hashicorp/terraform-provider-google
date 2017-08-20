package google

import (
	"reflect"
	"strings"
	"testing"
)

func TestConvertStringArr(t *testing.T) {

	input := make([]interface{}, 3)
	input[0] = "aaa"
	input[1] = "bbb"
	input[2] = "aaa"

	expected := []string{"aaa", "bbb", "ccc"}
	actual := convertStringArr(input)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("(%s) did not match expected value: %s", actual, expected)
	}
}

func TestConvertAndMapStringArr(t *testing.T) {

	input := make([]interface{}, 3)
	input[0] = "aaa"
	input[1] = "bbb"
	input[2] = "aaa"

	expected := []string{"AAA", "BBB", "CCC"}
	actual := convertAndMapStringArr(input, strings.ToUpper)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("(%s) did not match expected value: %s", actual, expected)
	}
}

func TestExtractLastResourceFromUri_withUrl(t *testing.T) {
	actual := extractLastResourceFromUri("http://something.com/one/two/three")
	expected := "three"
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

func TestExtractLastResourceFromUri_WithStaticValue(t *testing.T) {
	actual := extractLastResourceFromUri("three")
	expected := "three"
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}
