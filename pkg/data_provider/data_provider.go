/*
 * Copyright (C) 2020-2022, IrineSistiana
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package data_provider

import (
	"fmt"
	"github.com/IrineSistiana/mosdns/v4/pkg/notifier"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type DataManager struct {
	pm sync.RWMutex
	ps map[string]*DataProvider
}

type DataListener interface {
	Update(newData []byte) error
}

func NewDataManager() *DataManager {
	return &DataManager{
		ps: make(map[string]*DataProvider),
	}
}

func (m *DataManager) AddDataProvider(name string, p *DataProvider) {
	m.pm.Lock()
	defer m.pm.Unlock()
	m.ps[name] = p
}

func (m *DataManager) GetDataProvider(name string) *DataProvider {
	m.pm.RLock()
	defer m.pm.RUnlock()
	return m.ps[name]
}

type DataProviderConfig struct {
	Tag        string `yaml:"tag"`
	File       string `yaml:"file"`
	AutoReload bool   `yaml:"auto_reload"`
}

type DataProvider struct {
	logger     *zap.Logger
	file       string
	autoReload bool

	v         atomic.Value
	lm        sync.Mutex
	listeners []DataListener

	notifier *notifier.Notifier
	sc       *notifier.SafeClose
}

func NewDataProvider(lg *zap.Logger, cfg *DataProviderConfig) (*DataProvider, error) {
	dp := new(DataProvider)
	dp.logger = lg
	dp.file = cfg.File
	dp.autoReload = cfg.AutoReload

	dp.notifier = notifier.NewNotifier()
	dp.sc = notifier.NewSafeClose()

	if err := dp.init(); err != nil {
		return nil, err
	}
	return dp, nil
}

func (ds *DataProvider) init() error {
	v, err := ds.loadFromDisk()
	if err != nil {
		return err
	}
	ds.v.Store(v)

	if ds.autoReload {
		if err := ds.startFsWatcher(); err != nil {
			return fmt.Errorf("failed to start fs watcher, %w", err)
		}
	}
	return nil
}

func (ds *DataProvider) Close() {
	ds.sc.Done()
	ds.sc.CloseWait()
}

func (ds *DataProvider) GetNotifier() *notifier.Notifier {
	return ds.notifier
}

// LoadAndAddListener loads the DataListener, returns any error that occurs, and
// add this DataListener to this DataProvider.
func (ds *DataProvider) LoadAndAddListener(l DataListener) error {
	v := ds.GetData()
	if err := l.Update(v); err != nil {
		return err
	}
	ds.lm.Lock()
	defer ds.lm.Unlock()
	ds.listeners = append(ds.listeners, l)
	return nil
}

func (ds *DataProvider) GetData() []byte {
	return ds.v.Load().([]byte)
}

// updateData updates ds' v, notify the notifier and trigger all listeners.
func (ds *DataProvider) updateData(newData []byte) {
	ds.v.Store(newData)
	ds.notifier.Notify()
	ds.lm.Lock()
	defer ds.lm.Unlock()
	for _, l := range ds.listeners {
		if err := l.Update(newData); err != nil {
			ds.logger.Error(
				"failed to update data listener",
				zap.Error(err),
			)
		}
	}
}

func (ds *DataProvider) loadFromDisk() ([]byte, error) {
	return os.ReadFile(ds.file)
}

func (ds *DataProvider) startFsWatcher() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := w.Add(ds.file); err != nil {
		return err
	}

	go func() {
		var delayTimer *time.Timer
		for {
			select {
			case e, ok := <-w.Events:
				if !ok {
					return
				}
				ds.logger.Info(
					"fs event",
					zap.Stringer("event", e),
				)
				if delayTimer == nil {
					delayTimer = time.AfterFunc(time.Second, func() {
						ds.logger.Info(
							"reloading file",
							zap.String("file", ds.file),
						)
						if v, err := ds.loadFromDisk(); err != nil {
							ds.logger.Error(
								"failed to reload file",
								zap.String("file", ds.file),
								zap.Error(err),
							)
						} else {
							ds.logger.Info(
								"file reloaded",
								zap.String("file", ds.file),
							)
							ds.updateData(v)
						}
					})
				} else {
					delayTimer.Reset(time.Second)
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				ds.logger.Error("fs notify error", zap.Error(err))
			case <-ds.sc.ReceiveCloseSignal():
				return
			}
		}
	}()
	return nil
}