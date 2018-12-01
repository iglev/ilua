package ilua

import (
	"github.com/iglev/ilua/log"
)

func init() {
	// set std logger
	log.SetLogger(nil)
}
