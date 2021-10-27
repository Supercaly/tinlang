package tin

import "fmt"

type fileLocation struct {
	fileName string
	col      int
	row      int
}

func (fl fileLocation) String() string {
	return fmt.Sprintf("%s:%d:%d", fl.fileName, fl.row+1, fl.col+1)
}
