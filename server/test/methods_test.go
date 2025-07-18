package test

import (
	"fmt"
	"testing"

	"github.com/CristianVega28/goserver/core/middleware"
)

func TestNameMiddleware(t *testing.T) {
	result := middleware.CreateMapMiddleware()
	fmt.Println(result)
}
