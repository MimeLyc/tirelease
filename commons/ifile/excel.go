package ifile

import (
	"fmt"
	"os"
	"reflect"
	"time"

	excelize "github.com/xuri/excelize/v2"
)

// Create or cover excel file
// Using `tag` of variable to fetch head and values
// Will automatically create the parent dirs if not exists
// If the file is already exists, append sheet to it
func CreateExcelSheetByTag[T interface{}](s []T, dir, filename, sheetName string) error {
	if len(s) == 0 {
		return nil
	}
	CreateFileRecursively(dir, "")

	qualifiedName := fmt.Sprintf("%s/%s", dir, filename)

	f := excelize.NewFile()
	if _, err := os.Stat(qualifiedName); err == nil {
		f, err = excelize.OpenFile(qualifiedName)
		if err != nil {
			return err
		}
	}
	index := f.NewSheet(sheetName)
	// delete default sheet which named "Sheet1"
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)
	f = setSheetHead(s, sheetName, f)
	f = setSheetValue(s, sheetName, f)

	if err := f.SaveAs(qualifiedName); err != nil {
		return err
	}

	f.Close()
	return nil
}

// set sheet name using tag of variables
func setSheetHead[T interface{}](s []T, sheetName string, f *excelize.File) *excelize.File {
	sType := reflect.TypeOf(s[0])
	fields := reflect.VisibleFields(sType)
	styleID, _ := f.NewStyle(
		&excelize.Style{
			Font: &excelize.Font{
				Color: "#777777",
				Bold:  true,
			},
			Protection: &excelize.Protection{
				Locked: true,
			},
		},
	)

	i := 0
	for _, field := range fields {
		fieldName := field.Tag.Get("excel")
		if fieldName == "" {
			continue
		}
		column, _ := excelize.ColumnNumberToName(i)
		cellName, _ := excelize.JoinCellName(column, 1)
		f.SetCellValue(sheetName, cellName, fieldName)
		i++
	}
	f.SetRowStyle(sheetName, 1, 1, styleID)

	return f
}

// set sheet value using reflect to fetch value
func setSheetValue[T interface{}](s []T, sheetName string, f *excelize.File) *excelize.File {
	sType := reflect.TypeOf(s[0])
	fields := reflect.VisibleFields(sType)

	for i, row := range s {
		j := 0
		for _, field := range fields {
			fieldTag := field.Tag.Get("excel")
			if fieldTag == "" {
				continue
			}
			value := reflect.ValueOf(row)
			fieldValue := reflect.Indirect(value).FieldByName(field.Name)
			valueString := convertValueToString(fieldValue)
			columnNum, _ := excelize.ColumnNumberToName(j)
			cellName, _ := excelize.JoinCellName(columnNum, i+2)
			f.SetCellValue(sheetName, cellName, valueString)
			j++
		}
	}
	return f
}

func convertValueToString(value reflect.Value) string {
	rawValue := value
	if value.Kind() == reflect.Pointer {
		rawValue = value.Elem()
	}
	if !rawValue.IsValid() {
		return ""
	}

	if rawValue.CanInt() {
		return fmt.Sprintf("%d", rawValue.Int())
	} else if rawValue.CanFloat() {
		return fmt.Sprintf("%v", rawValue.Float())
	} else if rawValue.Kind() == reflect.String {
		return string(rawValue.String())
	} else if rawValue.Type().String() == "time.Time" {
		return rawValue.Interface().(time.Time).String()
	} else if rawValue.Kind() == reflect.Bool {
		return fmt.Sprintf("%t", rawValue.Bool())
	} else {
		return "Type not supported"
	}
}
