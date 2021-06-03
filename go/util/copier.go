package util

import (
	"fmt"
	"github.com/tlarsen7572/goalteryx/sdk"
)

type CopyData func(sdk.Record, map[string]interface{})

func FindFieldAndGenerateCopier(field string, incomingInfo sdk.IncomingRecordInfo) (CopyData, error) {
	var copier CopyData
	for _, incomingField := range incomingInfo.Fields() {
		if field == incomingField.Name {
			switch incomingField.Type {
			case `Byte`, `Int16`, `Int32`, `Int64`:
				intField, _ := incomingInfo.GetIntField(field)
				getInt := intField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getInt(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				return copier, nil
			case `String`, `WString`, `V_String`, `V_WString`:
				stringField, _ := incomingInfo.GetStringField(field)
				getString := stringField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getString(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				return copier, nil
			case `Bool`:
				boolField, _ := incomingInfo.GetBoolField(field)
				getBool := boolField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getBool(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				return copier, nil
			case `Date`, `DateTime`:
				timeField, _ := incomingInfo.GetTimeField(field)
				getTime := timeField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getTime(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				return copier, nil
			case `Float`, `Double`, `FixedDecimal`:
				floatField, _ := incomingInfo.GetFloatField(field)
				getFloat := floatField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) {
					value, isNull := getFloat(copyFrom)
					if isNull {
						copyTo[field] = nil
						return
					}
					copyTo[field] = value
				}
				return copier, nil
			}
		}
	}
	return nil, fmt.Errorf(`could not find field '%v' in the record`, field)
}
