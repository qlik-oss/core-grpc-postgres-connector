package qlik

import (
	"github.com/jackc/pgx"
)

var emptyString = ""

func RowsToGetDataResponse(rows *pgx.Rows) *GetDataResponse {
	var descrs = rows.FieldDescriptions()
	var types = getTypeConstants(descrs)
	var array = make([]*FieldInfo, len(descrs))

	for i := range descrs {
		array[i] = &FieldInfo{descrs[i].Name, types[i]}
	}
	return &GetDataResponse{array, "dummy"}
}


func getTypeConstants(fieldDescriptors []pgx.FieldDescription) []FieldType {
	var translators = make([]FieldType, len(fieldDescriptors))
	for i, fieldDescr := range fieldDescriptors {
		switch fieldDescr.DataTypeName {
		case "int4":
			translators[i] = FieldType_INTEGER
		default:
			translators[i] = FieldType_ASCII
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
	var this = &AsyncTranslator{writer, fieldDescriptors, make(chan [][]interface{}, 10)}
	go this.run()
	return this
}

func ( this *AsyncTranslator) GetDataResponseMetadata() *GetDataResponse {
	var types = getTypeConstants(this.fieldDescriptors)
	var array = make([]*FieldInfo, len(this.fieldDescriptors))

	for i := range this.fieldDescriptors {
		array[i] = &FieldInfo{this.fieldDescriptors[i].Name, types[i]}
	}
	return &GetDataResponse{array, ""}
}

func ( this *AsyncTranslator) buildRowBundle(tempQixRowList [][]interface{}) *BundledRows {
	var typeConsts = getTypeConstants(this.fieldDescriptors)
	var columnCount, rowCount = len(this.fieldDescriptors), int64(len(tempQixRowList))
	var rowBundle = BundledRows{Cols: make([]*Column, columnCount)}

	for i := 0; i < columnCount; i++ {
		var column = &Column{}
		switch typeConsts[i] {
		case FieldType_ASCII:
			column.StrIsNulls=make([]bool, rowCount)
			column.Strings=make([]string, rowCount)
			column.Numbers=nil
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
		case FieldType_INTEGER:
			column.StrIsNulls=nil
			column.Strings=nil
			column.Numbers=make([]float64, rowCount)
			for r := 0; r < len(tempQixRowList); r++ {
				var srcValue = tempQixRowList[r][i]
				column.Numbers[r] = float64(int64(srcValue.(int32)))
			}
		}
		rowBundle.Cols[i] = column
	}
	return &rowBundle
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
		var resultChunk = this.buildRowBundle(tempQixRowList)
		this.writer.Write(resultChunk)
	}
	this.writer.Close()
}