package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/tlarsendataguy/goalteryx/sdk"
)

type CopyData func(sdk.Record, map[string]interface{}) error

func FindFieldAndGenerateCopier(field string, incomingInfo sdk.IncomingRecordInfo) (CopyData, error) {
	var copier CopyData
	for _, incomingField := range incomingInfo.Fields() {
		if field == incomingField.Name {
			switch incomingField.Type {
			case `Byte`, `Int16`, `Int32`, `Int64`:
				intField, _ := incomingInfo.GetIntField(field)
				getInt := intField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) error {
					value, isNull := getInt(copyFrom)
					if isNull {
						copyTo[field] = nil
						return nil
					}
					copyTo[field] = value
					return nil
				}
				return copier, nil
			case `String`, `WString`, `V_String`, `V_WString`:
				stringField, _ := incomingInfo.GetStringField(field)
				getString := stringField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) error {
					value, isNull := getString(copyFrom)
					if isNull {
						copyTo[field] = nil
						return nil
					}
					copyTo[field] = value
					return nil
				}
				return copier, nil
			case `Bool`:
				boolField, _ := incomingInfo.GetBoolField(field)
				getBool := boolField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) error {
					value, isNull := getBool(copyFrom)
					if isNull {
						copyTo[field] = nil
						return nil
					}
					copyTo[field] = value
					return nil
				}
				return copier, nil
			case `Date`:
				timeField, _ := incomingInfo.GetTimeField(field)
				getTime := timeField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) error {
					value, isNull := getTime(copyFrom)
					if isNull {
						copyTo[field] = nil
						return nil
					}
					copyTo[field] = dbtype.Date(value)
					return nil
				}
				return copier, nil
			case `DateTime`:
				timeField, _ := incomingInfo.GetTimeField(field)
				getTime := timeField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) error {
					value, isNull := getTime(copyFrom)
					if isNull {
						copyTo[field] = nil
						return nil
					}
					copyTo[field] = value
					return nil
				}
				return copier, nil
			case `Float`, `Double`, `FixedDecimal`:
				floatField, _ := incomingInfo.GetFloatField(field)
				getFloat := floatField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) error {
					value, isNull := getFloat(copyFrom)
					if isNull {
						copyTo[field] = nil
						return nil
					}
					copyTo[field] = value
					return nil
				}
				return copier, nil
			case `Blob`:
				blobField, _ := incomingInfo.GetBlobField(field)
				getBlob := blobField.GetValue
				copier = func(copyFrom sdk.Record, copyTo map[string]interface{}) error {
					bytes := getBlob(copyFrom)
					if bytes == nil {
						copyTo[field] = nil
						return nil
					}
					var listData []interface{}
					err := json.Unmarshal(bytes, &listData)
					if err == nil {
						copyTo[field] = listData
						return nil
					}
					copyTo[field] = nil
					return errors.New(fmt.Sprintf(`error: field does not contain a valid JSON list: %v`, err.Error()))
				}
				return copier, nil
			}
		}
	}
	return nil, fmt.Errorf(`could not find field '%v' in the record`, field)
}
