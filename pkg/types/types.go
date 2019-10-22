package types

import (
	fluxevent "github.com/fluxcd/flux/pkg/event"
)

// Represents a Flux event that will get sent to an exporter
type Message struct {
	Title      string
	TitleLink  string
	EventType  string
	Namespaces []string
	Workloads  []string
	Commits    []Commit
	Body       string
	VCSRootURL string
	Tags       []string
	Event      fluxevent.Event
}

type Commit struct {
	User          string
	Revision      string
	ShortRevision string
	Message       string
}
