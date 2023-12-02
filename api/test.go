package api

import (
	"context"

	testv1 "github.com/koo04/gateway-test/internal/gen/proto/go/api/v1"
)

func (s *server) GetTest(ctx context.Context, req *testv1.GetTestRequest) (*testv1.TestResponse, error) {
	testString := ctx.Value(ContextTestString{})

	if testString == nil {
		return &testv1.TestResponse{Data: "no test string in the context"}, nil
	}

	return &testv1.TestResponse{Data: "from the context: " + testString.(string)}, nil
}
