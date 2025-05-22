package test

import (
	"fmt"
	"testing"

	"github.com/CristianVega28/goserver/server"
)

func TestNameMiddleware(t *testing.T) {
	result := server.CreateMapMiddleware()
	fmt.Println(result)
}
