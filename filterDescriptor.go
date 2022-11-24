package dancok

type FilterDescriptor struct {
	FieldName   string
	Operator    Operator
	Condition   Condition
	Value       any
	Value2      any
	RangeValues []any
}
