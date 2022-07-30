package errors

import (
	"fmt"
	v1 "github.com/letscrum/letscrum/apis/general/v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func HandleErr(err error) v1.ErrorResponse {
	s := status.Convert(err)
	var details []*anypb.Any
	if len(s.Details()) > 0 {
		for _, d := range s.Details() {
			strWrapper := wrapperspb.String(fmt.Sprintf("%v", d))
			strAny, _ := anypb.New(strWrapper)
			details = append(details, strAny)
		}
	}
	return v1.ErrorResponse{
		Code:    int32(s.Code()),
		Message: s.Message(),
		Details: details,
	}
}
