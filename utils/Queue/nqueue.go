// // 定义无锁队列结构
package Queue

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	//"time"
)

const forCount = 4

type NQueue[T any] struct {
	head      *node[T]
	tail      *node[T]
	status    atomic.Bool
	lock      atomic.Bool
	count     atomic.Int64
	recvCond  *sync.Cond
	recvLock  sync.RWMutex
	nodePool  sync.Pool
	zeroValue T ///零值
	closeOnce sync.Once
}

type node[T any] struct {
	value T
	next  *node[T]
}

// 新建队列，返回一个空队列
func NewNQueue[T any]() *NQueue[T] {
	qq := &NQueue[T]{}
	qq.status.Swap(true)
	qq.lock.Swap(false)
	qq.recvCond = sync.NewCond(qq.recvLock.RLocker())
	qq.nodePool = sync.Pool{
		New: func() any {
			return &node[T]{
				value: qq.zeroValue,
				next:  nil,
			}
		},
	}

	return qq
}

func (q *NQueue[T]) Close() {
	defer func() {
		q.recvLock.Lock()
		defer q.recvLock.Unlock()
		q.recvCond.Broadcast()
	}()
	q.status.CompareAndSwap(true, false)

	q.closeOnce.Do(func() {})
}

// 插入，将给定的值v放在队列的尾部
func (q *NQueue[T]) Enqueue(v T) error {
	defer func() {
		q.recvLock.Lock()
		defer q.recvLock.Unlock()
		q.recvCond.Broadcast()
	}()

	if !q.status.Load() {
		return fmt.Errorf("is close")
	}
	// n := &node[T]{value: v}

	n := q.nodePool.Get().(*node[T])
	n.value = v
	n.next = nil

	for {
		if !q.status.Load() {
			n.value = q.zeroValue
			q.nodePool.Put(n)
			return fmt.Errorf("is close")
		}
		if q.lock.CompareAndSwap(false, true) {

			if q.head == nil {
				q.head = n
			} else {
				if q.tail == nil {
					q.tail = n
					q.head.next = q.tail
				} else {
					oldTail := q.tail
					oldTail.next = n
					q.tail = n
				}
			}
			q.lock.Swap(false)

			q.count.Add(1)

			return nil

		}
		gosched()
	}
}

// /不阻塞
func (q *NQueue[T]) Dequeue() (t T, isClose bool) {
	defer func() { isClose = !q.status.Load() }()
	for {
		tt, ok, headIsNil := q.dequeue()
		if ok {
			t = tt
			return
		}
		if headIsNil {
			return
		}
		gosched()

	}
}

// 移除，删除并返回队列头部的值,如果队列为空，则返回nil
func (q *NQueue[T]) dequeue() (t T, ok bool, headIsNil bool) {
	ok = false
	headIsNil = false
	if q.lock.CompareAndSwap(false, true) {
		defer q.lock.Swap(false)
		if q.head == nil {
			headIsNil = true
			return
			// return nil, err
		} else {

			oldHead := q.head
			if oldHead.next == nil {
				q.head = nil
			} else {
				q.head = oldHead.next
				if q.head == q.tail {
					q.tail = nil
				}
			}

			ok = true
			q.count.Add(-1)
			t = oldHead.value

			oldHead.value = q.zeroValue
			oldHead.next = nil
			q.nodePool.Put(oldHead)
			return
		}
	}
	return
}

func (q *NQueue[T]) dequeueWait() (t T, ok bool) {
	q.recvLock.RLock()
	if q.status.Load() && q.count.Load() == 0 {
		q.recvCond.Wait()
	}
	q.recvLock.RUnlock()

	t, ok, _ = q.dequeue()
	return
}

// /阻塞    返回值t
func (q *NQueue[T]) DequeueWait() (t T, isClose bool) {
	defer func() { isClose = !q.status.Load() }()
	for {
		t = q.zeroValue
		tt, ok := q.dequeueWait()

		if ok {
			t = tt
			return
		}
		if !q.status.Load() && q.count.Load() == 0 {
			return
		}
		gosched()
	}
}

type DequeueFunc[T any] func(t T, isClose bool) bool

func (q *NQueue[T]) DequeueFunc(fn DequeueFunc[T]) (err error) {
	for {
		t, ok := q.dequeueWait()
		if ok {
			if !fn(t, !q.status.Load()) {
				return
			}
		}
		if !q.status.Load() && q.count.Load() == 0 {
			return fmt.Errorf("queue is close and empty")
		}
		gosched()
	}
}

func (q *NQueue[T]) Count() int64 {
	return q.count.Load()
}

func (q *NQueue[T]) Status() bool {
	return q.status.Load()
}

func gosched() {
	// for i := 0; i < forCount; i++ {
	runtime.Gosched()
	// }
}
