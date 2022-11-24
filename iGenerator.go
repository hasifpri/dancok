package dancok

type IGenerator interface {
	Generate(SelectParameter) string
	Parse(SelectParameter) string
	ParseFilter(SelectParameter) string
	ParsePaging(SelectParameter) string
	ParseSort(SelectParameter) string
}
