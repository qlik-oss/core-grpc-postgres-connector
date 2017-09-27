package postgres

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jackc/pgx"
	"github.com/qlik-ea/postgres-grpc-connector/qlik"
	"github.com/jackc/pgx/pgtype"
	"time"
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
	{DataTypeName: "timestampz"},
	{DataTypeName: "date"},
	{DataTypeName: "numeric"},
	{DataTypeName: "decimal"},
	{DataTypeName: "bool"},
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
		pgtype.Numeric{},
		pgtype.Decimal{},
		true,
	},
}

func TestAsyncTranslator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AsyncTranslator")
}

var _ = Describe("AsyncTranslator", func() {
	Context("Translate fieldDscriptions", func() {
		var asyncTranslator = AsyncTranslator{ fieldDescriptors: fieldDescriptions}
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
			Expect(typeConstants[9].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[9].SemanticType).To(Equal(qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC))
			Expect(typeConstants[10].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[10].SemanticType).To(Equal(qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC))
			Expect(typeConstants[11].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
			Expect(typeConstants[11].SemanticType).To(Equal(qlik.SemanticType_UNIX_SECONDS_SINCE_1970_UTC))
			Expect(typeConstants[12].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_REAL))
			Expect(typeConstants[13].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_REAL))
			Expect(typeConstants[14].FieldAttributes.Type).To(Equal(qlik.FieldAttrType_INTEGER))
		})
	})

	Context("Translate a row bundle", func() {
		var asyncTranslator = AsyncTranslator{fieldDescriptors: fieldDescriptions}
		var bundle = asyncTranslator.buildRowBundle(postgresRowData)


		var expectedBundle = qlik.DataChunk{Cols: []*qlik.Column {
			{Strings: []string {"varchar"}, Flags:[]qlik.ValueFlag{qlik.ValueFlag_Normal}},
			{Strings: []string {"text"}, Flags:[]qlik.ValueFlag{qlik.ValueFlag_Normal}},
			{Integers: []int64 {8}},
			{Integers: []int64 {4}},
			{Integers: []int64 {2}},
			{Integers: []int64 {1}},
			{Integers: []int64 {9}},
			{Doubles: []float64 {2.4}},
			{Doubles: []float64 {4.8}},
			{Integers: []int64 {int64(time1.Unix())}},
			{Integers: []int64 {int64(time1.Unix())}},
			{Integers: []int64 {int64(time1.Unix())}},
			{Doubles: []float64 {0}},
			{Doubles: []float64 {0}},
			{Integers: []int64 {-1}},
		}}

		It("should match the expected type constants", func() {
			Expect(bundle.Cols[0]).To(BeEquivalentTo(expectedBundle.Cols[0]))
			Expect(bundle.Cols[1]).To(BeEquivalentTo(expectedBundle.Cols[1]))
			Expect(bundle.Cols[2]).To(BeEquivalentTo(expectedBundle.Cols[2]))
			Expect(bundle.Cols[3]).To(BeEquivalentTo(expectedBundle.Cols[3]))
			Expect(bundle.Cols[4]).To(BeEquivalentTo(expectedBundle.Cols[4]))
			Expect(bundle.Cols[5]).To(BeEquivalentTo(expectedBundle.Cols[5]))
			Expect(bundle.Cols[6]).To(BeEquivalentTo(expectedBundle.Cols[6]))
			Expect(bundle.Cols[7]).To(BeEquivalentTo(expectedBundle.Cols[7]))
			Expect(bundle.Cols[8]).To(BeEquivalentTo(expectedBundle.Cols[8]))
			Expect(bundle.Cols[9]).To(BeEquivalentTo(expectedBundle.Cols[9]))
			Expect(bundle.Cols[10]).To(BeEquivalentTo(expectedBundle.Cols[10]))
			Expect(bundle.Cols[11]).To(BeEquivalentTo(expectedBundle.Cols[11]))
			Expect(bundle.Cols[12]).To(BeEquivalentTo(expectedBundle.Cols[12]))
			Expect(bundle.Cols[13]).To(BeEquivalentTo(expectedBundle.Cols[13]))
			Expect(bundle.Cols[14]).To(BeEquivalentTo(expectedBundle.Cols[14]))
		})
	})

})
