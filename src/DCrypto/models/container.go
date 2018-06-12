package models

import (
	"math/rand"
	"time"
)

func NewWorkerManager() *workerManager {
	return &workerManager{
		workerList: make([]*Worker, 0),
		workerHash: make(map[string]int, 0),
	}
}

func (manager *workerManager) Add(worker *Worker) bool {
	manager.Lock()
	defer manager.Unlock()

	for {
		workerId := worker.Unique()

		if index, ok := manager.workerHash[workerId]; ok {
			manager.workerList[index].UpdateTime = worker.UpdateTime
			break
		}
		manager.workerList = append(manager.workerList, worker)
		manager.size = len(manager.workerList)
		manager.workerHash[workerId] = manager.size - 1
		manager.capacity = cap(manager.workerList)
		break
	}
	return true
}

func (manager *workerManager) Delete(worker *Worker) bool {
	manager.Lock()
	defer manager.Unlock()

	workerId := worker.Unique()

	if index, ok := manager.workerHash[workerId]; ok {
		lastWorker := manager.workerList[manager.size-1]

		manager.workerList = manager.workerList[0 : manager.size-1]
		delete(manager.workerHash, workerId)

		manager.size = len(manager.workerList)
		manager.capacity = cap(manager.workerList)

		if index < manager.size {
			manager.workerList[index] = lastWorker
			manager.workerHash[lastWorker.Unique()] = index
		}

		return true
	}
	return false
}

func (manager *workerManager) Choose(limit int) (res []*Worker) {

	manager.RLock()
	defer manager.RUnlock()
	if limit < 1 || manager.size < 1 {
		return
	}

	rand.Seed(time.Now().UnixNano())

	res = make([]*Worker, limit)
	for i := 0; i < limit; i++ {
		idx := rand.Int()
		res[i] = manager.workerList[idx%manager.size]
	}
	return res
}

func (manager *workerManager) GetAllWorkers() []*Worker {
	manager.RLock()
	defer manager.RUnlock()
	return manager.workerList
}
