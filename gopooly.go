package gopooly

import (
	"context"
	"sync"
	"sync/atomic"
)

//Type untuk fungsi yang nantinya akan di overwrite ketika package di instance
type emptyFunc func(ctx context.Context, args interface{}) (interface{}, error)

//Struct untuk membungkus data pada task queue yang dikirimkan melalui channel
type taskFunc struct {
	ctx      context.Context
	args     interface{}
	taskFunc emptyFunc
}

//Main struct
type threadPool struct {
	wg        *sync.WaitGroup
	taskTotal int64
	ExecFunc  emptyFunc
	taskQueue chan taskFunc
}

//New fungsi ini digunakan untuk membuat instance struct threadpool
//dengan menyetel buffer channel dan jumlah pekerja yang digunakan.
//Ketika menggunakan fungsi ini maka user (pengguna package) harus melakukan penyesuaian secara
//terpisah business logic yang akan digunakan oleh para pekerja.
func New(workers, bufferSize int) *threadPool {
	t := &threadPool{
		wg:        new(sync.WaitGroup),
		taskQueue: make(chan taskFunc, bufferSize),
	}

	//Membuat pekerja dan menyetel proses yang akan dijalankan
	for i := 0; i < workers; i++ {
		t.wg.Add(1)

		go func() {
			defer t.wg.Done()

			for task := range t.taskQueue {
				atomic.AddInt64(&t.taskTotal, 1)

				task.taskFunc(task.ctx, task.args)
			}
		}()
	}

	return t
}

//NewFunc fungsi ini sama dengan New() yang membedakan adalah business logic untuk workers dikirimkan ketika melakukan instance fungsi
func NewFunc(workers, bufferSize int, execFunc emptyFunc) *threadPool {
	t := New(workers, bufferSize)
	t.ExecFunc = execFunc

	return t
}

//Len fungsi ini mengembalikan jumlah task queue yang belum diproses workers
func (t *threadPool) Len() int { return len(t.taskQueue) }

//Cap fungsi ini mengembalikan kapasitas yang dapat ditampung pada task queue
func (t *threadPool) Cap() int { return cap(t.taskQueue) }

//TaskTotal fungsi ini mengembalikan total task yang sudah di eksekusi oleh workers
func (t *threadPool) TaskTotal() int64 { return t.taskTotal }

//Process fungsi ini digunakan untuk mengirimkan request (eksekusi) kedalam antrian eksekusi (task queue)
func (t *threadPool) Process(ctx context.Context, args interface{}) {
	t.taskQueue <- taskFunc{ctx: ctx, args: args, taskFunc: t.ExecFunc}
}

//Close fungsi ini digunakan untuk menutup task queue. Fungsi ini wajib di-eksekusi ketika main app di tutup untuk menghindari
//workers sedang bekerja dan app diakhiri. (Gracefully Shutdown)
func (t *threadPool) Close() {
	close(t.taskQueue)

	t.wg.Wait()
}
