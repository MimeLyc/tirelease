package ifile

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCreateExcelSheetByTag(t *testing.T) {
	type test struct {
		Var1 string `excel:"name1"`
		Var2 string
	}

	testSlice := []test{
		{
			Var1: "111",
			Var2: "222",
		},
	}
	CreateExcelSheetByTag(testSlice, "tmp", "test.xlsx", "tmpsheet1")
}

func TestTimeKindString(t *testing.T) {
	test := time.Now()
	testValue := reflect.ValueOf(test)
	testKindString := testValue.Type().String()
	fmt.Printf("%s", testKindString)
}
