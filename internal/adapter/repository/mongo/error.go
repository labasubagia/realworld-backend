package mongo

import (
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
	"go.mongodb.org/mongo-driver/mongo"
)

var mapException = map[int]string{
	11000: exception.TypeValidation,
}

func getWriteConcernCode(err error) string {
	if err == nil {
		return ""
	}
	fail, ok := err.(mongo.WriteException)
	if !ok {
		return ""
	}
	if len(fail.WriteErrors) == 0 {
		return ""
	}
	return mapException[fail.WriteErrors[0].Code]
}

func intoException(err error) *exception.Exception {
	if err == nil {
		return nil
	}
	typeWriteConcern := getWriteConcernCode(err)
	if typeWriteConcern != "" {
		return exception.New(typeWriteConcern, err.Error(), err)
	}
	return exception.Into(err)
}
