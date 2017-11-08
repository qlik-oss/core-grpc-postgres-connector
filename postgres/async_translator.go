package postgres

import (
	"github.com/jackc/pgx"
	"github.com/qlik-ea/postgres-grpc-connector/qlik"
	"time"
	"github.com/jackc/pgx/pgtype"
	"fmt"
	"reflect"
)

func (this *AsyncTranslator) GetTypes() []*qlik.FieldInfo {

	var array = make([]*qlik.FieldInfo, len(this.fieldDescriptors))

	for i, fieldDescr := range this.fieldDescriptors {
		var semanticType = qlik.SemanticType_DEFAULT;
		var fieldAttrType = qlik.FieldAttrType_TEXT;
		switch fieldDescr.DataTypeName {
		case "varchart", "text":
			fieldAttrType = qlik.FieldAttrType_TEXT
		case "int8", "int4", "char", "int2", "oid":
			fieldAttrType = qlik.FieldAttrType_INTEGER
		case "float4", "float8":
			fieldAttrType = qlik.FieldAttrType_REAL
		case "timestamp", "timestampz":
			fieldAttrType = qlik.FieldAttrType_TIMESTAMP
			semanticType = qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC
		case "date":
			fieldAttrType = qlik.FieldAttrType_DATE
			semanticType = qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC
		case "numeric", "decimal":
			fieldAttrType = qlik.FieldAttrType_REAL
		case "bool":
			fieldAttrType = qlik.FieldAttrType_INTEGER

		default:
			fieldAttrType = qlik.FieldAttrType_TEXT
		}
		array[i] = &qlik.FieldInfo{this.fieldDescriptors[i].Name, semanticType, &qlik.FieldAttributes{Type: fieldAttrType}}

	}
	return array
}

/**
 *	Class AsyncStreamWriter
 */

type AsyncTranslator struct {
	writer           *qlik.AsyncStreamWriter
	fieldDescriptors []pgx.FieldDescription
	channel chan [][]interface{}
}

func NewAsyncTranslator(writer *qlik.AsyncStreamWriter, fieldDescriptors []pgx.FieldDescription) *AsyncTranslator {
	var this = &AsyncTranslator{writer, fieldDescriptors, make(chan [][]interface{}, 10)}
	go this.run()
	return this
}

func (this *AsyncTranslator) GetDataResponseMetadata() *qlik.GetDataResponse {
	var array = this.GetTypes();
	return &qlik.GetDataResponse{array, ""}
}

func (this *AsyncTranslator) buildRowBundle(tempQixRowList [][]interface{}) *qlik.DataChunk {
	var types = this.GetTypes();
	var columnCount, rowCount = len(this.fieldDescriptors), int64(len(tempQixRowList))
	var rowBundle = qlik.DataChunk{Cols: make([]*qlik.Column, columnCount)}

	if len(tempQixRowList) > 0 {
		for c := 0; c < columnCount; c++ {
			var column = &qlik.Column{}
			if types[c].SemanticType == qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC {
				column.Integers = make([]int64, rowCount)
				for r := 0; r < len(tempQixRowList); r++ {
					var srcValue = tempQixRowList[r][c]
					switch tempQixRowList[r][c].(type) {
					case time.Time:
						column.Integers[r] = srcValue.(time.Time).Unix()
					default:
						fmt.Println(srcValue)
					}
				}
			} else {
				switch types[c].FieldAttributes.Type {
				case qlik.FieldAttrType_TEXT:
					column.Flags = make([]qlik.ValueFlag, rowCount)
					column.Strings = make([]string, rowCount)
					for r := 0; r < len(tempQixRowList); r++ {
						var srcValue = tempQixRowList[r][c]
						if srcValue != nil {
							switch tempQixRowList[0][c].(type) {
							case string:
								column.Strings[r] = srcValue.(string)
								column.Flags[r] = qlik.ValueFlag_Normal
							default:
								column.Strings[r] = "<Unsupported format>"
								column.Flags[r] = qlik.ValueFlag_Normal
							}
						} else {
							column.Strings[r] = ""
							column.Flags[r] = qlik.ValueFlag_Null
						}
					}
				case qlik.FieldAttrType_REAL:
					column.Doubles = make([]float64, rowCount)
					for r := 0; r < len(tempQixRowList); r++ {
						var srcValue = tempQixRowList[r][c]
						switch tempQixRowList[r][c].(type) {
						case float64:
							column.Doubles[r] = float64(srcValue.(float64))
						case float32:
							column.Doubles[r] = float64(srcValue.(float32))
						case pgtype.Numeric:
							var value = srcValue.(pgtype.Numeric)
							value.AssignTo(&column.Doubles[r])
						case pgtype.Decimal:
							var value = srcValue.(pgtype.Decimal)
							value.AssignTo(&column.Doubles[r])
						default:
							fmt.Println(srcValue)
							fmt.Println("Unknown format", srcValue);
						}
					}
				case qlik.FieldAttrType_INTEGER:
					column.Integers = make([]int64, rowCount)
					for r := 0; r < len(tempQixRowList); r++ {
						var srcValue = tempQixRowList[r][c]
						switch tempQixRowList[r][c].(type) {
						case int:
							column.Integers[r] = int64(srcValue.(int))
						case int64:
							column.Integers[r] = srcValue.(int64)
						case int32:
							column.Integers[r] = int64(srcValue.(int32))
						case int16:
							column.Integers[r] = int64(srcValue.(int16))
						case int8:
							column.Integers[r] = int64(srcValue.(int8))
						case bool:
							if srcValue.(bool) {
								column.Integers[r] = -1
							} else {
								column.Integers[r] = 0
							}
						default:
							fmt.Println(reflect.TypeOf(srcValue))
						}
					}
				}
			}
			rowBundle.Cols[c] = column
		}
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
	for tempQixRowList := range this.channel {
		var resultChunk = this.buildRowBundle(tempQixRowList)
		this.writer.Write(resultChunk)
	}
	this.writer.Close()
}
