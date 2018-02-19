package postgres

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	qlik "github.com/qlik-ea/postgres-grpc-connector/qlik"
)

// GetTypes returns a list of fields and their types.
func (a *AsyncTranslator) GetTypes() []*qlik.FieldInfo {

	var array = make([]*qlik.FieldInfo, len(a.fieldDescriptors))

	for i, fieldDescr := range a.fieldDescriptors {
		var semanticType = qlik.SemanticType_DEFAULT
		var fieldAttrType = qlik.FieldAttrType_TEXT
		switch fieldDescr.DataTypeName {
		case "varchart", "text":
			fieldAttrType = qlik.FieldAttrType_TEXT
		case "int8", "int4", "char", "int2", "oid":
			fieldAttrType = qlik.FieldAttrType_INTEGER
		case "float4", "float8":
			fieldAttrType = qlik.FieldAttrType_REAL
		case "timestamp", "timestamptz":
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
		array[i] = &qlik.FieldInfo{
			Name:            a.fieldDescriptors[i].Name,
			SemanticType:    semanticType,
			FieldAttributes: &qlik.FieldAttributes{Type: fieldAttrType},
		}
	}
	return array
}

// AsyncTranslator defines the translator interface.
type AsyncTranslator struct {
	writer           *AsyncStreamWriter
	fieldDescriptors []pgx.FieldDescription
	channel          chan [][]interface{}
}

// NewAsyncTranslator constructs a new translator.
func NewAsyncTranslator(writer *AsyncStreamWriter, fieldDescriptors []pgx.FieldDescription) *AsyncTranslator {
	var this = &AsyncTranslator{writer, fieldDescriptors, make(chan [][]interface{}, 10)}
	go this.run()
	return this
}

// GetDataResponseMetadata returns the metadata for a specific dataset.
func (a *AsyncTranslator) GetDataResponseMetadata() *qlik.GetDataResponse {
	var array = a.GetTypes()
	return &qlik.GetDataResponse{FieldInfo: array, TableName: ""}
}

func (a *AsyncTranslator) buildDataChunk(tempQixRowList [][]interface{}) *qlik.DataChunk {
	var types = a.GetTypes()
	var columnCount = len(a.fieldDescriptors)
	var maxSize = len(tempQixRowList) * len(a.fieldDescriptors)
	var dataChunk = qlik.DataChunk{StringBucket: make([]string, 0, maxSize),
		DoubleBucket: make([]float64, 0, maxSize),
		StringCodes:  make([]int32, 0, maxSize),
		NumberCodes:  make([]int64, 0, 2*maxSize)}

	if len(tempQixRowList) > 0 {
		for r := 0; r < len(tempQixRowList); r++ {
			for c := 0; c < columnCount; c++ {
				if types[c].SemanticType == qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC {
					var srcValue = tempQixRowList[r][c]
					switch tempQixRowList[r][c].(type) {
					case time.Time:
						dataChunk.DoubleBucket = append(dataChunk.DoubleBucket, float64(srcValue.(time.Time).Unix()))
						dataChunk.StringCodes = append(dataChunk.StringCodes, -1)
						dataChunk.NumberCodes = append(dataChunk.NumberCodes, int64(len(dataChunk.DoubleBucket)-1))
					default:
						dataChunk = addNull(dataChunk)
						fmt.Println("Unknown format", reflect.TypeOf(srcValue))
					}
				} else {
					switch types[c].FieldAttributes.Type {
					case qlik.FieldAttrType_TEXT:
						var srcValue = tempQixRowList[r][c]
						if srcValue != nil {
							switch tempQixRowList[r][c].(type) {
							case string:
								dataChunk = addString(dataChunk, srcValue.(string))
							default:
								dataChunk = addNull(dataChunk)
								fmt.Println("Unknown format", reflect.TypeOf(srcValue))
							}
						} else {
							dataChunk = addString(dataChunk, "")
						}
					case qlik.FieldAttrType_REAL:
						var srcValue = tempQixRowList[r][c]
						switch tempQixRowList[r][c].(type) {
						case float64:
							dataChunk = addNumber(dataChunk, srcValue.(float64))
						case float32:
							dataChunk = addNumber(dataChunk, float64(srcValue.(float32)))
						case pgtype.Numeric:
							var value = srcValue.(pgtype.Numeric)
							var tempValue float64
							value.AssignTo(&tempValue)
							dataChunk = addNumber(dataChunk, tempValue)
						case pgtype.Decimal:
							var value = srcValue.(pgtype.Decimal)
							var tempValue float64
							value.AssignTo(&tempValue)
							dataChunk = addNumber(dataChunk, tempValue)
						default:
							dataChunk = addNull(dataChunk)
							fmt.Println("Unknown format", reflect.TypeOf(srcValue))
						}
					case qlik.FieldAttrType_INTEGER:
						var srcValue = tempQixRowList[r][c]
						switch tempQixRowList[r][c].(type) {
						case int:
							dataChunk = addInteger(dataChunk, int64(srcValue.(int)))
						case int64:
							dataChunk = addInteger(dataChunk, int64(srcValue.(int64)))
						case int32:
							dataChunk = addInteger(dataChunk, int64(srcValue.(int32)))
						case int16:
							dataChunk = addInteger(dataChunk, int64(srcValue.(int16)))
						case int8:
							dataChunk = addInteger(dataChunk, int64(srcValue.(int8)))
						case bool:
							if srcValue.(bool) {
								dataChunk = addInteger(dataChunk, 1)
							} else {
								dataChunk = addInteger(dataChunk, 0)
							}
						default:
							dataChunk = addNull(dataChunk)
							fmt.Println("Unknown format", reflect.TypeOf(srcValue))
						}
					}
				}
			}
		}
	}
	return &dataChunk
}

func addString(dataChunk qlik.DataChunk, stringValue string) qlik.DataChunk {
	if stringValue != "" {
		dataChunk.StringBucket = append(dataChunk.StringBucket, stringValue)
		dataChunk.StringCodes = append(dataChunk.StringCodes, int32(len(dataChunk.StringBucket)-1))
		dataChunk.NumberCodes = append(dataChunk.NumberCodes, -1)
	} else {
		dataChunk.StringCodes = append(dataChunk.StringCodes, -2)
		dataChunk.NumberCodes = append(dataChunk.NumberCodes, -1)
	}
	return dataChunk
}

func addNumber(dataChunk qlik.DataChunk, numberValue float64) qlik.DataChunk {
	dataChunk.DoubleBucket = append(dataChunk.DoubleBucket, numberValue)
	dataChunk.StringCodes = append(dataChunk.StringCodes, -1)
	dataChunk.NumberCodes = append(dataChunk.NumberCodes, int64(len(dataChunk.DoubleBucket)-1))
	return dataChunk
}

func addInteger(dataChunk qlik.DataChunk, intValue int64) qlik.DataChunk {
	dataChunk.StringCodes = append(dataChunk.StringCodes, -1)
	dataChunk.NumberCodes = append(dataChunk.NumberCodes, -2)
	dataChunk.NumberCodes = append(dataChunk.NumberCodes, intValue)
	return dataChunk
}

func addNull(dataChunk qlik.DataChunk) qlik.DataChunk {
	dataChunk.StringCodes = append(dataChunk.StringCodes, -1)
	dataChunk.NumberCodes = append(dataChunk.NumberCodes, -1)
	return dataChunk
}

func (a *AsyncTranslator) Write(values [][]interface{}) {
	a.channel <- values
}

// Close will close the underlying channel.
func (a *AsyncTranslator) Close() {
	close(a.channel)
}

func (a *AsyncTranslator) run() {
	for tempQixRowList := range a.channel {
		var resultChunk = a.buildDataChunk(tempQixRowList)
		a.writer.Write(resultChunk)
	}
	a.writer.Close()
}
