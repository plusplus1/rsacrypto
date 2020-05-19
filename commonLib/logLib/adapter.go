package logLib

import (
	"github.com/sirupsen/logrus"
)

const (
	keyWorker = "worker"
)

type Adapter struct {
	workerName string
}

func NewLogAdapter(workerName string) Adapter {
	return Adapter{
		workerName: workerName,
	}
}

func (a Adapter) WithFields(fds logrus.Fields) *logrus.Entry {
	if fds == nil {
		fds = logrus.Fields{"worker": a.workerName}
	} else if _, ok := fds[keyWorker]; !ok {
		fds[keyWorker] = a.workerName
	}
	return logrus.WithFields(fds)
}

func (a Adapter) WithField(k string, v interface{}) *logrus.Entry {
	return a.WithFields(logrus.Fields{k: v})
}

func (a Adapter) Infof(format string, args ...interface{}) {
	logrus.WithField(keyWorker, a.workerName).Infof(format, args...)
}

func (a Adapter) Debugf(format string, args ...interface{}) {
	logrus.WithField(keyWorker, a.workerName).Debugf(format, args...)
}

func (a Adapter) Errorf(format string, args ...interface{}) {
	logrus.WithField(keyWorker, a.workerName).Errorf(format, args...)
}

func (a Adapter) Warnf(format string, args ...interface{}) {
	logrus.WithField(keyWorker, a.workerName).Warnf(format, args...)
}
