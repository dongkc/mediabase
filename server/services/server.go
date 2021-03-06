package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/helper"
	"apertoire.net/mediabase/server/message"
	"apertoire.net/mediabase/server/model"
	"apertoire.net/mediabase/server/static"
	"fmt"
	"github.com/apertoire/mlog"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

const apiVersion string = "/api/v1"
const docPath string = ""

type Server struct {
	Bus    *bus.Bus
	Config *model.Config
	r, s   *gin.Engine
}

func (self *Server) Start() {
	mlog.Info("starting server service")

	self.r = gin.New()

	self.r.Use(gin.Recovery())
	// self.r.Use(helper.Logging())

	path := filepath.Join(".", "web")

	var abs string
	var err error
	if abs, err = filepath.Abs(path); err != nil {
		mlog.Info("unable to get absolute path: %s, ", err)
		return
	}

	mlog.Info("server root path is: %s", abs)

	// self.r.Use(static.Serve("./web/"))
	self.r.Use(static.Serve(path))
	self.r.NoRoute(self.redirect)

	api := self.r.Group(apiVersion)
	{
		api.GET("/movies/cover", self.getCover)
		api.POST("/movies", self.getMovies)
		api.POST("/movies/search", self.searchMovies)

		api.GET("/movies/duplicates", self.getDuplicates)

		api.PUT("/movies/watched", self.watchedMovie)
		api.POST("/movies/fix", self.fixMovie)
		api.POST("/movies/prune", self.pruneMovies)

		api.GET("/import", self.importMovies)
		api.GET("/import/status", self.importMoviesStatus)

		api.GET("/config", self.getConfig)
		api.PUT("/config", self.saveConfig)
	}

	mlog.Info("service started listening on %s:%s", self.Config.Host, self.Config.Port)

	go self.r.Run(fmt.Sprintf("%s:%s", self.Config.Host, self.Config.Port))
}

func (self *Server) Stop() {
	mlog.Info("server service stopped")
	// nothing here
}

func (self *Server) redirect(c *gin.Context) {
	c.Redirect(301, "/index.html")
}

func (self *Server) getConfig(c *gin.Context) {
	msg := message.GetConfig{Reply: make(chan *model.Config)}
	self.Bus.GetConfig <- &msg

	reply := <-msg.Reply
	c.JSON(200, &reply)
}

func (self *Server) saveConfig(c *gin.Context) {
	var conf model.Config

	c.Bind(&conf)
	msg := message.SaveConfig{Config: &conf, Reply: make(chan *model.Config)}
	self.Bus.SaveConfig <- &msg

	self.Config = <-msg.Reply
	c.JSON(200, &self.Config)
}

func (self *Server) getCover(c *gin.Context) {
	msg := message.Movies{Reply: make(chan *message.MoviesDTO)}
	self.Bus.GetCover <- &msg
	reply := <-msg.Reply

	// mlog.Info("response is: %s", reply)

	// helper.WriteJson(w, 200, &reply)
	c.JSON(200, &reply)
}

func (self *Server) getMovies(c *gin.Context) {
	var options message.Options

	c.Bind(&options)

	mlog.Info("bocelli: %+v", options)

	msg := message.Movies{Options: options, Reply: make(chan *message.MoviesDTO)}
	self.Bus.GetMovies <- &msg
	reply := <-msg.Reply

	// mlog.Info("response is: %s", reply)

	c.JSON(200, &reply)
}

func (self *Server) importMovies(c *gin.Context) {
	mlog.Info("importMovies: you know .. i got here")

	msg := message.Status{Reply: make(chan *message.Context)}
	self.Bus.ImportMovies <- &msg
	reply := <-msg.Reply

	c.JSON(200, &reply)
}

func (self *Server) importMoviesStatus(c *gin.Context) {
	msg := message.Status{Reply: make(chan *message.Context)}
	self.Bus.ImportMoviesStatus <- &msg
	reply := <-msg.Reply

	c.JSON(200, &reply)
}

func (self *Server) searchMovies(c *gin.Context) {
	mlog.Info("searchMovies: are you a head honcho ?")

	var options message.Options

	c.Bind(&options)

	mlog.Info("anyway the wind blows: %+v", options)

	msg := message.Movies{Options: options, Reply: make(chan *message.MoviesDTO)}
	self.Bus.SearchMovies <- &msg
	reply := <-msg.Reply

	// mlog.Info("%s", reply)

	c.JSON(200, &reply)
}

func (self *Server) pruneMovies(c *gin.Context) {
	mlog.Info("pruning .. i got here")

	msg := message.PruneMovies{Reply: make(chan string)}
	self.Bus.PruneMovies <- &msg
	reply := <-msg.Reply

	data := struct {
		Description string
	}{Description: reply}
	c.JSON(200, &data)
}

func (self *Server) getDuplicates(c *gin.Context) {
	msg := message.Movies{Reply: make(chan *message.MoviesDTO)}
	self.Bus.ShowDuplicates <- &msg
	reply := <-msg.Reply

	// mlog.Info("response is: %s", reply)
	c.JSON(200, &reply)
}

func (self *Server) watchedMovie(c *gin.Context) {
	var movie message.Movie

	c.Bind(&movie)
	// mlog.Info("%+v", movie)

	msg := message.SingleMovie{Movie: &movie, Reply: make(chan *message.Movie)}
	self.Bus.WatchedMovie <- &msg
	reply := <-msg.Reply

	c.JSON(200, &reply)
}

func (self *Server) fixMovies(w http.ResponseWriter, req *http.Request) {
	self.Bus.FixMovies <- 1
	helper.WriteJson(w, 200, "ok")
}

func (self *Server) fixMovie(c *gin.Context) {
	var movie message.Movie

	c.Bind(&movie)
	mlog.Info("%+v", movie)

	msg := message.SingleMovie{Movie: &movie, Reply: make(chan *message.Movie)}
	self.Bus.FixMovie <- &msg

	data := struct {
		Status bool `json:"status"`
	}{Status: true}

	c.JSON(200, &data)
}
