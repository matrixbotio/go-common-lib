package nano

import (
	"fmt"

	"github.com/jaevor/go-nanoid"
)

func ID() string {
	id, err := nanoid.Standard(11)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return id()
}
