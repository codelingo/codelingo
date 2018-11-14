package mypack

// good
func NewGoodWorkerGood() (*GoodWorker) {
}

// bad
func NewBadWorkerBad() (worker.Worker) {
}
