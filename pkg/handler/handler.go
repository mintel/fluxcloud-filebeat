package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	fluxevent "github.com/fluxcd/flux/pkg/event"
	"github.com/mintel/fluxcloud-filebeat/pkg/config"
	"github.com/mintel/fluxcloud-filebeat/pkg/types"
	"html/template"
	"log"
	"net"
	"strings"

	"github.com/mintel/fluxcloud-filebeat/pkg/utils"
)

type Handler interface {
	BuildMessage(event fluxevent.Event) types.Message
	Handle(msg types.Message) error
}

type handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) (Handler, error) {
	return &handler{
		config: config,
	}, nil
}

func (h *handler) BuildMessage(event fluxevent.Event) types.Message {

	// Default title for the event
	title := event.String()

	// Commits are a simplified version of fluxevent.Commits
	var commits []types.Commit

	// Tags are just a list of strings holding useful information that may be
	// filtered on later (in ES/Grafana etc
	var tags []string
	tags = append(tags, event.Type)

	// Parse out events - this is mostly to grab the commits and tags.
	switch event.Type {
	case fluxevent.EventRelease:
		metadata := event.Metadata.(*fluxevent.ReleaseEventMetadata)
		commit := types.Commit{Message: metadata.Cause.Message, User: metadata.Cause.User, Revision: metadata.Revision, ShortRevision: shortRevision(metadata.Revision)}
		commits = append(commits, commit)
		tags = append(tags, shortRevision(metadata.Revision))

	case fluxevent.EventCommit:
		metadata := event.Metadata.(*fluxevent.CommitEventMetadata)

		message := "See affected workloads"
		var user string

		if len(metadata.Spec.Cause.Message) > 0 {
			message = metadata.Spec.Cause.Message
		}
		if len(metadata.Spec.Cause.User) > 0 {
			user = metadata.Spec.Cause.User
		}

		commit := types.Commit{Message: message, User: user, Revision: metadata.Revision, ShortRevision: shortRevision(metadata.Revision)}
		commits = append(commits, commit)
		tags = append(tags, shortRevision(metadata.Revision))

	case fluxevent.EventAutoRelease:
		metadata := event.Metadata.(*fluxevent.AutoReleaseEventMetadata)
		imageIDs := metadata.Result.ChangedImages()
		if len(imageIDs) == 0 {
			imageIDs = []string{"<no image>"}
		}

		title = fmt.Sprintf(
			"Automatically released %s",
			strings.Join(imageIDs, ", "),
		)

		commit := types.Commit{Message: "See affected workloads", Revision: metadata.Revision, ShortRevision: shortRevision(metadata.Revision)}
		commits = append(commits, commit)
		tags = append(tags, shortRevision(metadata.Revision))

	case fluxevent.EventSync:
		metadata := event.Metadata.(*fluxevent.SyncEventMetadata)
		commitCount := len(metadata.Commits)

		if commitCount > 0 {
			for _, c := range metadata.Commits {
				commit := types.Commit{Message: c.Message, Revision: c.Revision, ShortRevision: shortRevision(c.Revision)}
				commits = append(commits, commit)
				tags = append(tags, shortRevision(c.Revision))
			}
		}

		title = fmt.Sprintf(
			"Synced %d commits", len(metadata.Commits),
		)
	}

	// Common across all events
	affectedWorkloads := make([]string, len(event.ServiceIDs)-1)

	// Parse out namespaces and services into separate fields.
	affectedNamespaces := []string{}

	for _, serviceID := range event.ServiceIDs {
		namespace, kind, name := serviceID.Components()
		txt := fmt.Sprintf("%s/%s", kind, name)
		if !utils.StringInSlice(txt, affectedWorkloads) {
			affectedWorkloads = append(affectedWorkloads, txt)
		}

		if !utils.StringInSlice(namespace, affectedNamespaces) {
			affectedNamespaces = append(affectedNamespaces, namespace)
		}

		if !utils.StringInSlice(name, tags) {
			tags = append(tags, name)
		}
	}

	// Create message with everything we need - this all gets passed on to FileBeat.
	message := types.Message{
		Title:      title,
		EventType:  event.Type,
		Namespaces: affectedNamespaces,
		Workloads:  affectedWorkloads,
		Commits:    commits,
		VCSRootURL: h.config.VCSRootURL,
		Tags:       tags,
	}

	if len(h.config.KeepFluxEvents) > 0 {
		// Event is rather verbose, so it's optional
		message.Event = event
	}

	return message
}

func (h *handler) Handle(message types.Message) error {

	// Format Message as HTML via template.
	// This is because Grafana annotations can use an HTML formatted field.
	paths := []string{
		"templates/event.html",
	}

	// Render it and put it back into the Body field.
	t := template.Must(template.New("event.html").ParseFiles(paths...))
	var renderedMessageBuffer bytes.Buffer
	err := t.Execute(&renderedMessageBuffer, message)
	renderedMessage := renderedMessageBuffer.String()
	message.Body = renderedMessage

	if err != nil {
		log.Println("Could not format event in template", err)
		return nil
	}

	// Now dump the message to a tcp socket.
	if len(h.config.FileBeatAddress) > 0 {
		tcpAddr, err := net.ResolveTCPAddr("tcp", h.config.FileBeatAddress)
		if err != nil {
			log.Println("ResolveTCPAddr failed:", err.Error())
			return nil
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Println("Dial failed:", err.Error())
			return nil
		}

		err = json.NewEncoder(conn).Encode(message)
		if err != nil {
			log.Println("Could encode message:", err)
			return nil
		}
	}

	return nil
}

// Helper to create a short git-commit rev.
func shortRevision(rev string) string {
	if len(rev) <= 7 {
		return rev
	}
	return rev[:7]
}
