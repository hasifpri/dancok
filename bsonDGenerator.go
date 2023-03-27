package dancok

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BsonDGenerator struct {
	DefaultFieldForSort string
}

func NewBsonDGenerator(defaultFieldForSort string) *BsonDGenerator {
	return &BsonDGenerator{defaultFieldForSort}
}

func (g *BsonDGenerator) ParseFilter(param SelectParameter) primitive.D {
	bsonFilter := bson.D{}
	if len(param.FilterDescriptors) > 0 {
		isFirstFilter := true
		for _, filter := range param.FilterDescriptors {
			if isFirstFilter {
				bsonFilter = append(bsonFilter, bson.E{Key: filter.FieldName, Value: bson.D{bson.E{Key: g.GetOperator(filter.Operator), Value: filter.Value}}})
				isFirstFilter = false
			} else {
				if filter.Condition == And {
					bsonFilter = append(bsonFilter, bson.E{Key: filter.FieldName, Value: bson.D{bson.E{Key: g.GetOperator(filter.Operator), Value: filter.Value}}})
				} else { // OR
					//Todo
					//bsonFilter = append(bsonFilter, bson.E{Key: "$or", Value: bson.A{bson.E{Key: filter.FieldName, Value: bson.D{bson.E{Key: g.GetOperator(filter.Operator), Value: filter.Value}}}}})
				}
			}
		}
	}

	if len(param.CompositeFilterDescriptors) > 0 {
		//isFirstCompositeFilter := true
		bsonA := bson.A{}
		for _, filter := range param.CompositeFilterDescriptors {
			// if isFirstCompositeFilter {
			// 	isFirstCompositeFilter = false
			// } else {

			// }
			isFirstItem := true
			for _, item := range filter.GroupFilterDescriptor.Items {
				if isFirstItem {
					isFirstItem = false
				}
				bsonA = append(bsonA, bson.D{bson.E{Key: item.FieldName, Value: bson.D{bson.E{Key: g.GetOperator(item.Operator), Value: item.Value}}}})
			}
		}

		bsonE := bson.E{Key: "$or", Value: bsonA}
		bsonFilter = append(bsonFilter, bsonE)
	}

	return bsonFilter
}

func (g *BsonDGenerator) ParseSort(param SelectParameter) primitive.D {
	bsonSort := bson.D{}

	if len(param.SortDescriptors) > 0 {
		for _, sort := range param.SortDescriptors {
			if sort.SortDirection == Ascending {
				bsonSort = append(bsonSort, primitive.E{Key: sort.FieldName, Value: 1})
			} else {
				bsonSort = append(bsonSort, primitive.E{Key: sort.FieldName, Value: -1})
			}
		}
	} else {
		bsonSort = append(bsonSort, primitive.E{Key: g.DefaultFieldForSort, Value: 1})
	}

	return bsonSort
}

func (g *BsonDGenerator) GetOperator(operator Operator) string {
	filterText := ""
	switch opt := operator; opt {
	case IsEqual:
		filterText = "$eq"
	case IsNotEqual:
		filterText = "$ne"
	case IsLessThan:
		filterText = "$lt"
	case IsLessThanOrEqual:
		filterText = "$lte"
	case IsMoreThan:
		filterText = "$gt"
	case IsMoreThanOrEqual:
		filterText = "$gte"
	case IsContain:
		filterText = "$regex"
	case IsBeginWith:
		filterText = "$regex"
	case IsEndWith:
		filterText = "$regex"
	case IsBetween:
		filterText = "$eq" //TODO
	case IsIn:
		filterText = "$eq" //TODO
	case IsNotIn:
		filterText = "$eq" //TODO
	}

	return filterText
}
