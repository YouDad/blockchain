package global

import (
	"sync"
)

var syncMutex sync.Mutex
var updateMutex sync.Mutex

func SyncLock() {
	// log.SetCallerLevel(1)
	// log.Debugln("sync lock before")
	// log.SetCallerLevel(0)
	syncMutex.Lock()
	// log.SetCallerLevel(1)
	// log.Debugln("sync lock after")
	// log.SetCallerLevel(0)
}

func SyncUnlock() {
	syncMutex.Unlock()
	// log.SetCallerLevel(1)
	// log.Debugln("sync unlock")
	// log.SetCallerLevel(0)
}

func UpdateLock() {
	// log.SetCallerLevel(1)
	// log.Debugln("update lock before")
	// log.SetCallerLevel(0)
	updateMutex.Lock()
	// log.SetCallerLevel(1)
	// log.Debugln("update lock after")
	// log.SetCallerLevel(0)
}

func UpdateUnlock() {
	// log.SetCallerLevel(1)
	// log.Debugln("update unlock")
	// log.SetCallerLevel(0)
	updateMutex.Unlock()
}
