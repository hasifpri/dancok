package dancok

import "strconv"

type SqlGenerator struct {
	TableName           string
	DefaultFieldForSort string
}

func NewSqlGenerator(tableName string, defaultFieldForSort string) *SqlGenerator {
	return &SqlGenerator{tableName, defaultFieldForSort}
}

func (g *SqlGenerator) Generate(param SelectParameter) string {
	result := ""
	result = "select * from (select ROW_NUMBER() OVER(" + g.ParseSort(param) + ") as RowNumber,* from " + g.TableName + " " + g.ParseFilter(param) + ") T " + g.ParsePaging(param)
	return result
}

func (g *SqlGenerator) Parse(param SelectParameter) string {
	result := g.ParseFilter(param) + g.ParseSort(param)

	return result
}

func (g *SqlGenerator) ParseFilter(param SelectParameter) string {
	filterText := ""
	if len(param.FilterDescriptors) > 0 {
		filterText = " WHERE "
		isFirstFilter := true
		for _, filter := range param.FilterDescriptors {
			if isFirstFilter {
				filterText = filterText + filter.FieldName
				isFirstFilter = false
			} else {
				if filter.Condition == And {
					filterText = filterText + " AND " + filter.FieldName
				} else {
					filterText = filterText + " OR " + filter.FieldName
				}
			}

			switch opt := filter.Operator; opt {
			case IsEqual:
				filterText = filterText + " = '" + filter.Value.(string) + "'"
			case IsNotEqual:
				filterText = filterText + " != '" + filter.Value.(string) + "'"
			case IsLessThan:
				filterText = filterText + " < " + filter.Value.(string)
			case IsLessThanOrEqual:
				filterText = filterText + " <= " + filter.Value.(string)
			case IsMoreThan:
				filterText = filterText + " > " + filter.Value.(string)
			case IsMoreThanOrEqual:
				filterText = filterText + " >= " + filter.Value.(string)
			case IsContain:
				filterText = filterText + " LIKE '%" + filter.Value.(string) + "%'"
			case IsBeginWith:
				filterText = filterText + " LIKE '" + filter.Value.(string) + "%'"
			case IsEndWith:
				filterText = filterText + " LIKE '%" + filter.Value.(string) + "'"
			case IsBetween:
				filterText = filterText + " BETWEEN '" + filter.Value.(string) + "' AND '" + filter.Value2.(string) + "'"
			case IsIn:
				filterText = filterText + " IN (" + ParseRangeValues(filter.RangeValues) + ")"
			case IsNotIn:
				filterText = filterText + " NOT IN (" + ParseRangeValues(filter.RangeValues) + ")"
			}
		}
	}

	if len(param.CompositeFilterDescriptors) > 0 {
		isFirstCompositeFilter := true
		for _, filter := range param.CompositeFilterDescriptors {
			if isFirstCompositeFilter {
				if filterText == "" {
					filterText = " WHERE ("
				} else {
					filterText = filterText + " " + string(filter.Condition) + " ("
				}
				isFirstCompositeFilter = false
			} else {
				if filter.Condition == And {
					filterText = filterText + " AND ("
				} else {
					filterText = filterText + " OR ("
				}
			}

			isFirstItem := true
			for _, item := range filter.GroupFilterDescriptor.Items {
				if isFirstItem {
					switch opt := item.Operator; opt {
					case IsEqual:
						filterText = filterText + item.FieldName + " = '" + item.Value.(string) + "'"
					case IsNotEqual:
						filterText = filterText + item.FieldName + " != '" + item.Value.(string) + "'"
					case IsLessThan:
						filterText = filterText + item.FieldName + " < " + item.Value.(string)
					case IsLessThanOrEqual:
						filterText = filterText + item.FieldName + " <= " + item.Value.(string)
					case IsMoreThan:
						filterText = filterText + item.FieldName + " > " + item.Value.(string)
					case IsMoreThanOrEqual:
						filterText = filterText + item.FieldName + " >= " + item.Value.(string)
					}

					isFirstItem = false
				} else {
					switch opt := item.Operator; opt {
					case IsEqual:
						filterText = filterText + " " + string(filter.GroupFilterDescriptor.Condition) + " " + item.FieldName + " = '" + item.Value.(string) + "'"
					case IsNotEqual:
						filterText = filterText + " " + string(filter.GroupFilterDescriptor.Condition) + " " + item.FieldName + " != '" + item.Value.(string) + "'"
					case IsLessThan:
						filterText = filterText + " " + string(filter.GroupFilterDescriptor.Condition) + " " + item.FieldName + " < " + item.Value.(string)
					case IsLessThanOrEqual:
						filterText = filterText + " " + string(filter.GroupFilterDescriptor.Condition) + " " + item.FieldName + " <= " + item.Value.(string)
					case IsMoreThan:
						filterText = filterText + " " + string(filter.GroupFilterDescriptor.Condition) + " " + item.FieldName + " > " + item.Value.(string)
					case IsMoreThanOrEqual:
						filterText = filterText + " " + string(filter.GroupFilterDescriptor.Condition) + " " + item.FieldName + " >= " + item.Value.(string)
					}
				}
			}

			filterText = filterText + ")"
		}
	}

	return filterText
}

func (g *SqlGenerator) ParsePaging(param SelectParameter) string {
	pagingText := ""

	if param.PageDescriptor.PageIndex == 0 && param.PageDescriptor.PageSize == 0 {
		pagingText = ""
	} else {
		startRowIndex := (param.PageDescriptor.PageIndex * param.PageDescriptor.PageSize) + 1
		endRowIndex := (param.PageDescriptor.PageIndex + 1) * param.PageDescriptor.PageSize
		pagingText = " where RowNumber between " + strconv.FormatInt(int64(startRowIndex), 10) + " and " + strconv.FormatInt(int64(endRowIndex), 10)
	}
	return pagingText
}

func (g *SqlGenerator) ParseSort(param SelectParameter) string {
	sortText := " "

	if len(param.SortDescriptors) > 0 {
		isFirstSort := true
		sortText = sortText + "order by"
		for _, sort := range param.SortDescriptors {
			if isFirstSort {
				sortText = sortText + " " + sort.FieldName
				isFirstSort = false
			} else {
				sortText = sortText + "," + sort.FieldName
			}

			if sort.SortDirection == Ascending {
				sortText = sortText + " asc"
			} else {
				sortText = sortText + " desc"
			}
		}
	} else {
		sortText = sortText + " order by " + g.DefaultFieldForSort + " desc"
	}

	return sortText
}

func ParseRangeValues(values []any) string {
	valueText := ""
	if len(values) > 0 {
		isFirstValue := true
		_, isStringType := values[0].(string)
		if isStringType {
			for _, v := range values {
				if isFirstValue {
					valueText = "'" + v.(string) + "'"
					isFirstValue = false
				} else {
					valueText = valueText + ",'" + v.(string) + "'"
				}
			}
		} else {
			for _, v := range values {
				if isFirstValue {
					valueText = string(v.(int32))
					isFirstValue = false
				} else {
					valueText = valueText + "," + string(v.(int32))
				}
			}
		}
	}
	return valueText
}
