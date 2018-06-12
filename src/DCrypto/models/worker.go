package models

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
)

func (w *Worker) Unique() string {
	return fmt.Sprintf("%v:%v", w.IP, w.Port)
}

func (w *Worker) Decrypt(f *Frame, wc *workerManager) ([]byte, error) {

	if tcpConn, err := net.Dial("tcp", w.Unique()); err != nil {
		wc.Delete(w)
		return nil, err
	} else {
		defer tcpConn.Close()

		enc := gob.NewEncoder(tcpConn)
		if err := enc.Encode(f); err != nil {
			log.Printf("[ERROR]\tGob Encode fail, %v", err)
			return nil, err
		}

		// start read response

		dec := gob.NewDecoder(tcpConn)
		ret := FrameRet{}
		if err := dec.Decode(&ret); err != nil {
			log.Printf("[ERROR]\tGob Decode fail, %v", err)
			return nil, err
		}

		if ret.Error != "" {
			return nil, errors.New(ret.Error)
		}
		return ret.Data, nil
	}
}
