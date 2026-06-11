package responses

import (
	"fmt"
	"time"
)

func ResponseID() string {
	return fmt.Sprintf("resp_%d", time.Now().UnixNano())
}

func MessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

func EventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
