package utils

import (
	"go.mongodb.org/mongo-driver/bson"
)

func MarshallStructtoBson(v interface{}) ([]byte, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	return data, err
}

func UnMarshallBsonArray() {

}
