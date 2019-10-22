package server

import (
	"context"
	"encoding/json"
	"fmt"
	fluxevent "github.com/fluxcd/flux/pkg/event"
	"github.com/gin-gonic/gin"
	"github.com/mintel/fluxcloud-filebeat/pkg/handler"
	"log"
	"net/http"
	"reflect"
)

type Server struct {
	handler handler.Handler
	router  *gin.Engine
	server  *http.Server
}

func NewServer(port int, handler handler.Handler) *Server {
	router := gin.Default()
	s := &Server{
		handler: handler,
		router:  router,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
	}
	router.POST("/v1/event", s.handleFluxEvent)
	router.GET("/healthz", s.healthz)
	return s
}

func (s *Server) healthz(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *Server) handleFluxEvent(c *gin.Context) {
	var msg struct {
		Event fluxevent.Event
	}
	c.Status(http.StatusOK)
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filebeatMsg := s.handler.BuildMessage(msg.Event)
	if err := s.handler.Handle(filebeatMsg); err != nil {
		dumpedMsg, _ := json.Marshal(filebeatMsg)
		log.Println("Failed to send filebeat message: {}", err, dumpedMsg)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) Start() error {
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(reflect.TypeOf(err).String())
			log.Println("failed to run server", err)
		}
	}()
	return nil
}

func (s *Server) Close() error {
	return s.server.Shutdown(context.Background())
}
