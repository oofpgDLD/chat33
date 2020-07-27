package result

import (
	"fmt"
	"testing"
)

func Test_ParseError(t *testing.T) {
	fmt.Println(ParseError(-1000, ""))
}
