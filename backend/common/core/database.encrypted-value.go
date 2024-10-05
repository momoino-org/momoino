package core

import (
	"context"
	"fmt"
	"reflect"

	"gorm.io/gorm/schema"
)

type EncryptionSerializer struct {
	Encryptor Encryptor
}

var _ schema.SerializerInterface = (*EncryptionSerializer)(nil)

func (es EncryptionSerializer) Scan(
	ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	dbValue interface{},
) error {
	if dbValue != nil {
		var encryptedValue string

		switch v := dbValue.(type) {
		case []byte:
			encryptedValue = string(v)
		case string:
			encryptedValue = v
		default:
			return fmt.Errorf("failed to convert value type to string: %#v", dbValue)
		}

		decryptedValue, err := es.Encryptor.Decrypt(encryptedValue)
		if err != nil {
			return err
		}

		field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(*decryptedValue))
	}

	return nil
}

func (es EncryptionSerializer) Value(
	ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	fieldValue interface{},
) (interface{}, error) {
	var plainTextInBytes []byte

	switch v := fieldValue.(type) {
	case []byte:
		plainTextInBytes = v
	case string:
		plainTextInBytes = []byte(v)
	default:
		return nil, fmt.Errorf("failed to convert value type to string: %#v", fieldValue)
	}

	return es.Encryptor.Encrypt(plainTextInBytes)
}
