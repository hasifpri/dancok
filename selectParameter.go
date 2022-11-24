package dancok

type SelectParameter struct {
	FilterDescriptors          []FilterDescriptor
	CompositeFilterDescriptors []CompositeFilterDescriptor
	SortDescriptors            []SortDescriptor
	PageDescriptor             PageDescriptor
}
