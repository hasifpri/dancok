package dancok

import "strconv"

type RediSearchGenerator struct {
	DefaultFieldForSort string
}

func NewRediSearchGenerator(defaultFieldForSort string) *RediSearchGenerator {
	return &RediSearchGenerator{defaultFieldForSort}
}

func (g *RediSearchGenerator) Generate(param SelectParameter) string {
	result := "'" + g.ParseFilter(param) + "'" + g.ParseSort(param) + g.ParsePaging(param)

	return result
}

func (g *RediSearchGenerator) Parse(param SelectParameter) string {
	result := g.ParseFilter(param) + g.ParseSort(param)

	return result
}

func (g *RediSearchGenerator) ParseFilter(param SelectParameter) string {
	filterText := "*"
	if len(param.FilterDescriptors) > 0 {
		filterText = ""
		isFirstFilter := true
		for _, filter := range param.FilterDescriptors {
			currentFilter := ""
			currentFilter = "@" + filter.FieldName + ":"

			switch opt := filter.Operator; opt {
			case IsEqual:
				currentFilter = currentFilter + filter.Value.(string)
			case IsNotEqual:
				currentFilter = "-" + currentFilter + filter.Value.(string)
			case IsLessThan:
				currentFilter = currentFilter + "[-inf (" + filter.Value.(string) + "]"
			case IsLessThanOrEqual:
				currentFilter = currentFilter + "[-inf " + filter.Value.(string) + "]"
			case IsMoreThan:
				currentFilter = currentFilter + "[(" + filter.Value.(string) + " +inf]"
			case IsMoreThanOrEqual:
				currentFilter = currentFilter + "[" + filter.Value.(string) + " +inf]"
			case IsContain:
				currentFilter = currentFilter + "*" + filter.Value.(string) + "*"
			case IsBeginWith:
				currentFilter = currentFilter + filter.Value.(string) + "*"
			case IsEndWith:
				currentFilter = currentFilter + "*" + filter.Value.(string)
			case IsBetween:
				currentFilter = currentFilter + "[" + filter.Value.(string) + " " + filter.Value2.(string) + "]"
			case IsIn:
				currentFilter = currentFilter + ParseRangeValuesRediSearch(filter.RangeValues)
			case IsNotIn:
				currentFilter = "-" + currentFilter + ParseRangeValuesRediSearch(filter.RangeValues)
			}

			if isFirstFilter {
				filterText = filterText + "(" + currentFilter + ")"
				isFirstFilter = false
			} else {
				if filter.Condition == And {
					filterText = filterText + " (" + currentFilter + ")"
				} else {
					filterText = filterText + "|(" + currentFilter + ")"
				}
			}
		}
	}

	if len(param.CompositeFilterDescriptors) > 0 {
		isFirstCompositeFilter := true
		for _, filter := range param.CompositeFilterDescriptors {
			if isFirstCompositeFilter {
				if filterText == "*" {
					filterText = ""
				} else {
					if filter.Condition == And {
						filterText = filterText + " ("
					} else {
						filterText = filterText + "|("
					}
				}
				isFirstCompositeFilter = false
			} else {
				if filter.Condition == And {
					filterText = filterText + " "
				} else {
					filterText = filterText + "|"
				}
			}

			isFirstItem := true
			currentGroupFilter := ""
			separatorCondition := " "
			if filter.GroupFilterDescriptor.Condition == And {
				separatorCondition = " "
			} else {
				separatorCondition = "|"
			}
			for _, item := range filter.GroupFilterDescriptor.Items {
				currentItemGroupFilter := ""
				currentItemGroupFilter = "@" + item.FieldName + ":"

				switch opt := item.Operator; opt {
				case IsEqual:
					currentItemGroupFilter = currentItemGroupFilter + item.Value.(string)
				case IsNotEqual:
					currentItemGroupFilter = "-" + currentItemGroupFilter + item.Value.(string)
				case IsLessThan:
					currentItemGroupFilter = currentItemGroupFilter + "[-inf (" + item.Value.(string) + "]"
				case IsLessThanOrEqual:
					currentItemGroupFilter = currentItemGroupFilter + "[-inf " + item.Value.(string) + "]"
				case IsMoreThan:
					currentItemGroupFilter = currentItemGroupFilter + "[(" + item.Value.(string) + " +inf]"
				case IsMoreThanOrEqual:
					currentItemGroupFilter = currentItemGroupFilter + "[" + item.Value.(string) + " +inf]"
				case IsContain:
					currentItemGroupFilter = currentItemGroupFilter + "*" + item.Value.(string) + "*"
				case IsBeginWith:
					currentItemGroupFilter = currentItemGroupFilter + item.Value.(string) + "*"
				case IsEndWith:
					currentItemGroupFilter = currentItemGroupFilter + "*" + item.Value.(string)
				}

				if isFirstItem {
					currentGroupFilter = currentGroupFilter + "(" + currentItemGroupFilter + ")"
					isFirstItem = false
				} else {
					currentGroupFilter = currentGroupFilter + separatorCondition + "(" + currentItemGroupFilter + ")"
				}
			}

			filterText = filterText + currentGroupFilter + ")"
		}
	}

	return filterText
}

func (g *RediSearchGenerator) ParsePaging(param SelectParameter) string {
	pagingText := ""

	if param.PageDescriptor.PageIndex == 0 && param.PageDescriptor.PageSize == 0 {
		pagingText = ""
	} else {
		startRowIndex := (param.PageDescriptor.PageIndex * param.PageDescriptor.PageSize)

		pagingText = " LIMIT " + strconv.FormatInt(int64(startRowIndex), 10) + " " + strconv.FormatInt(int64(param.PageDescriptor.PageSize), 10)
	}
	return pagingText
}

func (g *RediSearchGenerator) ParseSort(param SelectParameter) string {
	sortText := " "

	if len(param.SortDescriptors) > 0 {
		sortText = " SORTBY"
		sortText = sortText + " " + param.SortDescriptors[0].FieldName
		if param.SortDescriptors[0].SortDirection == Ascending {
			sortText = sortText + " ASC"
		} else {
			sortText = sortText + " DESC"
		}

	} else {
		sortText = sortText + " SORTBY " + g.DefaultFieldForSort + " DESC"
	}

	return sortText
}

func ParseRangeValuesRediSearch(values []any) string {
	valueText := ""
	if len(values) > 0 {
		valueText = "("
		isFirstValue := true
		_, isStringType := values[0].(string)
		if isStringType {
			for _, v := range values {
				if isFirstValue {
					valueText = valueText + v.(string)
					isFirstValue = false
				} else {
					valueText = valueText + "|" + v.(string)
				}
			}
		} else {
			for _, v := range values {
				if isFirstValue {
					valueText = string(v.(int32))
					isFirstValue = false
				} else {
					valueText = valueText + "|" + string(v.(int32))
				}
			}
		}

		valueText = valueText + ")"
	}
	return valueText
}
