package aggregate_test

import (
	"github.com/matiux/memo/domain/aggregate"
	"time"
)

var memoId = aggregate.NewUUIDv4()
var body = "Vegetables are good"
var creationDate = time.Now()

func createMemo() *aggregate.Memo {

	return aggregate.NewMemo(memoId, body, creationDate)
}
