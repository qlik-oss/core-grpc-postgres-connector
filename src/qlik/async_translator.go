package qlik

import (
	"strconv"
	"github.com/jackc/pgx"
)

var emptyString = "";
func translateStringToQix(srcValue interface{}) (*string, float64, bool) {
	if srcValue == nil {
		return &emptyString, 0, true
	}
	var val = srcValue.(string)
	return &val, 0, false
}
func translateInt32ToQix(srcValue interface{}) (*string, float64, bool) {
	var val = strconv.FormatInt(int64(srcValue.(int32)), 10)
	return &val, float64(int64(srcValue.(int32))), srcValue == nil
}

type TranslatorFunction func(srcValue interface{}) (*string, float64, bool);

func buildRowBundle(tempQixRowList [][]interface{}, translators []TranslatorFunction) *BundledRows2 {

	var columnCount, rowCount = int64(len(translators)), int64(len(tempQixRowList))
	var sliceSize = columnCount * rowCount;
	var rowBundle = BundledRows2{ColumnCount: columnCount, RowCount: rowCount, Numerics: make([]float64, 0, sliceSize), Strings: make([]string, 0, sliceSize), StrIsNulls: make([]bool, 0, sliceSize)}

	for r := 0; r<len(tempQixRowList); r++ {
		var postgresCols = tempQixRowList[r]
		for i := 0; i < len(postgresCols); i++ {
			var translatorFunc = translators[i]
			var str, val, strIsNull = translatorFunc(postgresCols[i]);
			rowBundle.Strings = append(rowBundle.Strings, *str)
			rowBundle.Numerics = append(rowBundle.Numerics, val)
			rowBundle.StrIsNulls = append(rowBundle.StrIsNulls, strIsNull)
		}
	}
	return &rowBundle
}

func ( this *AsyncTranslator) buildRowBundle2(tempQixRowList [][]interface{}) *BundledRows2 {

	var typeConsts = getTypeConstants(this.fieldDescriptors);
	var columnCount, rowCount = int64(len(this.fieldDescriptors)), int64(len(tempQixRowList))
	var sliceSize = columnCount * rowCount;
	var rowBundle = BundledRows2{ColumnCount: columnCount, RowCount: rowCount, Numerics: make([]float64, sliceSize), Strings: make([]string, sliceSize), StrIsNulls: make([]bool, sliceSize)}

	var loc = 0;
	for r := 0; r<len(tempQixRowList); r++ {
		var postgresCols = tempQixRowList[r]
		for i := 0; i < len(postgresCols); i++ {
			var srcValue = postgresCols[i];
			switch typeConsts[i] {
			case STRING:
				if srcValue != nil {
					rowBundle.Strings[loc] = srcValue.(string)
					rowBundle.Numerics[loc] = 0
					rowBundle.StrIsNulls[loc] = false
				} else {
					rowBundle.Strings[loc] = emptyString
					rowBundle.Numerics[loc] = 0
					rowBundle.StrIsNulls[loc] = true
				}
			case INT4:
				rowBundle.Strings[loc] = strconv.FormatInt(int64(srcValue.(int32)), 10)
				rowBundle.Numerics[loc] = float64(int64(srcValue.(int32)))
				rowBundle.StrIsNulls[loc] = false
			}
			loc++
		}
	}
	return &rowBundle
}

func ( this *AsyncTranslator) buildRowBundle3(tempQixRowList [][]interface{}) *BundledRows3 {

	var typeConsts = getTypeConstants(this.fieldDescriptors);
	var columnCount, rowCount = len(this.fieldDescriptors), int64(len(tempQixRowList))
	var rowBundle = BundledRows3{RowCount: rowCount, Column: make([]*Column, columnCount)}

	for i := 0; i < columnCount; i++ {
		var column = &Column{}
		switch typeConsts[i] {
		case STRING:
			column.StrIsNulls=make([]bool, rowCount);
			column.Strings=make([]string, rowCount);
			column.Numerics=nil;
			for r := 0; r < len(tempQixRowList); r++ {
				var srcValue = tempQixRowList[r][i]
				if srcValue != nil {
					column.Strings[r] = srcValue.(string)
					column.StrIsNulls[r] = false
				} else {
					column.Strings[r] = emptyString
					column.StrIsNulls[r] = true
				}
			}
		case INT4:
			column.StrIsNulls=nil;
			column.Strings=make([]string, rowCount);
			column.Numerics=make([]float64, rowCount);
			for r := 0; r < len(tempQixRowList); r++ {
				var srcValue = tempQixRowList[r][i]
				column.Strings[r] = strconv.FormatInt(int64(srcValue.(int32)), 10)
				column.Numerics[r] = float64(int64(srcValue.(int32)))
			}
		}
		rowBundle.Column[i] = column;
	}
	return &rowBundle
}

func getTranslators(fieldDescriptors []pgx.FieldDescription) []TranslatorFunction {
	var translators = make([]TranslatorFunction, len(fieldDescriptors));
	for i, fieldDescr := range fieldDescriptors {
		switch fieldDescr.DataTypeName {
		case "int4":
			translators[i] = translateInt32ToQix
		default:
			translators[i] = translateStringToQix
		}
	}
	return translators
}

const STRING, INT4 = 0,1;

func getTypeConstants(fieldDescriptors []pgx.FieldDescription) []int {
	var translators = make([]int, len(fieldDescriptors));
	for i, fieldDescr := range fieldDescriptors {
		switch fieldDescr.DataTypeName {
		case "int4":
			translators[i] = INT4
		default:
			translators[i] = STRING
		}
	}
	return translators
}

/**
 *	Class AsyncStreamWriter
 */

type AsyncTranslator struct {
	writer *AsyncStreamWriter
	fieldDescriptors []pgx.FieldDescription
	channel chan [][]interface{}
}

func NewAsyncTranslator(writer *AsyncStreamWriter, fieldDescriptors []pgx.FieldDescription) *AsyncTranslator {
	var this = &AsyncTranslator{writer, fieldDescriptors, make(chan [][]interface{}, 100)}
	go this.run()
	return this
}
func (this *AsyncTranslator) Write(values [][]interface{}) {
	this.channel <- values
}

func (this *AsyncTranslator) Close() {
	close(this.channel)
}

func (this *AsyncTranslator) run() {
	//var translators = getTranslators(this.fieldDescriptors);
	for tempQixRowList := range this.channel {
		//this.writer.Write(buildRowBundle(tempQixRowList, translators))
		//var resultChunk = &ResultChunk{ResultSpec: nil, CellsByRow:this.buildRowBundle2(tempQixRowList)}
		var resultChunk = &ResultChunk{ResultSpec: nil, CellsByColumn:this.buildRowBundle3(tempQixRowList)}
		this.writer.Write(resultChunk)
	}
	this.writer.Close();
}