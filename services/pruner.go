package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"github.com/apertoire/mlog"
	"os"
)

type Pruner struct {
	Bus    *bus.Bus
	Config *helper.Config
}

func (self *Pruner) Start() {
	mlog.Info("starting Pruner service ...")

	go self.react()

	mlog.Info("Pruner service started")
}

func (self *Pruner) Stop() {
	// nothing right now
	mlog.Info("Pruner service stopped")
}

func (self *Pruner) react() {
	for {
		select {
		case msg := <-self.Bus.PruneMovies:
			go self.doPruneMovies(msg.Reply)
		}
	}
}

func (self *Pruner) doPruneMovies(reply chan string) {
	mlog.Info("Looking for something to prune")

	msg := message.ListMovies{make(chan []*message.Movie)}
	self.Bus.ListMovies <- &msg
	items := <-msg.Reply

	for _, item := range items {

		if _, err := os.Stat(item.Location); err != nil {
			if os.IsNotExist(err) {
				mlog.Info("UP FOR DELETION: [%d] %s (%s))", item.Id, item.Title, item.Location)
				self.Bus.DeleteMovie <- item
			}
		}

	}
}
