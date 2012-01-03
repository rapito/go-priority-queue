package prio

func (q *Queue) Get(i int) Interface {
	return q.get(i)
}
