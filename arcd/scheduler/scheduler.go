/*
 * Arc - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package scheduler

import (
	"github.com/evilsocket/arc/arcd/db"
	"github.com/evilsocket/arc/arcd/events"
	"github.com/evilsocket/arc/arcd/log"
	"time"
)

func worker(secs int) {
	period := time.Duration(secs) * time.Second

	log.Debugf("Scheduler started with a %v period.", period)

	for {
		time.Sleep(period)

		db.Lock()

		for _, store := range db.GetStores() {
			for _, r := range store.Children() {
				meta := r.Meta()
				if r.Expired() {
					if r.WasNotified() == false {
						events.Add(events.RecordExpired(r))
						r.SetNotified(true)
					}

					if meta.Prune {
						log.Infof("Pruning record %d ( %s ) ...", meta.Id, meta.Title)
						if _, err := store.Del(meta.Id); err != nil {
							log.Errorf("Error while deleting record %d: %s.", meta.Id, err)
						}
					}
				}
			}
		}

		db.Unlock()
	}
}

func Start(period int) {
	go worker(period)
}
