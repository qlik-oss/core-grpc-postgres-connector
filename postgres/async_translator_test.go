package postgres

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	qlik "github.com/qlik-ea/core-grpc-postgres-connector/qlik"
)

var fieldDescriptions = []pgx.FieldDescription{
	{DataTypeName: "varchar"},
	{DataTypeName: "text"},
	{DataTypeName: "int8"},
	{DataTypeName: "int4"},
	{DataTypeName: "char"},
	{DataTypeName: "int2"},
	{DataTypeName: "oid"},
	{DataTypeName: "float4"},
	{DataTypeName: "float8"},
	{DataTypeName: "timestamp"},
	{DataTypeName: "timestamptz"},
	{DataTypeName: "date"},
	{DataTypeName: "numeric"},
	{DataTypeName: "decimal"},
	{DataTypeName: "bool"},
	{DataTypeName: "text"},
}

var time1, _ = time.Parse(time.RFC3339, "20150717T00:00:00+00:00")

var postgresRowData = [][]interface{}{
	{
		"varchar",
		"text",
		8,
		4,
		2,
		1,
		9,
		2.4,
		4.8,
		time1,
		time1,
		time1,
		pgtype.Numeric{Int: big.NewInt(10), Exp: 0, Status: pgtype.Present},
		pgtype.Decimal{},
		true,
		nil,
	},
}

func TestAsyncTranslator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AsyncTranslator")
}

var _ = Describe("AsyncTranslator", func() {
	Context("Translate fieldDescriptions", func() {
		var asyncTranslator = AsyncTranslator{fieldDescriptors: fieldDescriptions}
		var typeConstants = asyncTranslator.GetDataResponseMetadata().FieldInfo
		It("should match the expected type constants", func() {
			Expect(typeConstants[0].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_TEXT))
			Expect(typeConstants[1].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_TEXT))
			Expect(typeConstants[2].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[3].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[4].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[5].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[6].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[7].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_REAL))
			Expect(typeConstants[8].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_REAL))
			Expect(typeConstants[9].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_TIMESTAMP))
			Expect(typeConstants[9].SemanticType).To(Equal(qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC))
			Expect(typeConstants[10].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_TIMESTAMP))
			Expect(typeConstants[10].SemanticType).To(Equal(qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC))
			Expect(typeConstants[11].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_DATE))
			Expect(typeConstants[11].SemanticType).To(Equal(qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC))
			Expect(typeConstants[12].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_REAL))
			Expect(typeConstants[13].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_REAL))
			Expect(typeConstants[14].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[15].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_TEXT))
		})
	})

	Context("Translate a row bundle", func() {
		var asyncTranslator = AsyncTranslator{fieldDescriptors: fieldDescriptions}
		var bundle = asyncTranslator.buildDataChunk(postgresRowData)
		var expectedOutCome = qlik.DataChunk{
			StringBucket: []string{"varchar", "text"},
			DoubleBucket: []float64{2.4, 4.8, float64(time1.Unix()), float64(time1.Unix()), float64(time1.Unix()), 10, 0},
			StringCodes:  []int32{0, 1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -2},
			NumberCodes:  []int64{-1, -1, -2, 8, -2, 4, -2, 2, -2, 1, -2, 9, 0, 1, 2, 3, 4, 5, 6, -2, 1, -1}}

		It("should match the expected type contants", func() {
			Expect(len(bundle.StringBucket)).To(BeIdenticalTo(2))
			for i := 0; i < len(bundle.StringBucket); i++ {
				Expect(bundle.StringBucket[i]).To(BeEquivalentTo(expectedOutCome.StringBucket[i]))
			}
			Expect(len(bundle.DoubleBucket)).To(BeIdenticalTo(7))
			for i := 0; i < len(bundle.DoubleBucket); i++ {
				Expect(bundle.DoubleBucket[i]).To(BeEquivalentTo(expectedOutCome.DoubleBucket[i]))
			}
			Expect(len(bundle.StringCodes)).To(BeIdenticalTo(16))
			for i := 0; i < len(bundle.StringCodes); i++ {
				Expect(bundle.StringCodes[i]).To(BeEquivalentTo(expectedOutCome.StringCodes[i]))
			}
			Expect(len(bundle.NumberCodes)).To(BeIdenticalTo(22))
			for i := 0; i < len(bundle.NumberCodes); i++ {
				Expect(bundle.NumberCodes[i]).To(BeEquivalentTo(expectedOutCome.NumberCodes[i]))
			}
		})
	})

})
