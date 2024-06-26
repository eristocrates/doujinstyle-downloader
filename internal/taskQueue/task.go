package taskQueue

import (
	"github.com/relepega/doujinstyle-downloader/internal/downloader"
	"github.com/relepega/doujinstyle-downloader/internal/playwrightWrapper"
)

type Task struct {
	Active           bool
	Done             bool
	ServiceNumber    int
	DownloadProgress int8
	Error            error
	UrlSlug          string
	ch               chan int8
}

func NewTask(s string, serviceNumber int) *Task {
	return &Task{
		Active:           false,
		Done:             false,
		ServiceNumber:    serviceNumber,
		DownloadProgress: -1, // -1: The downloader cannot calculate the download progress
		Error:            nil,
		UrlSlug:          s,
	}
}

func (t *Task) Activate() {
	t.Active = true
}

func (t *Task) Deactivate() {
	t.Active = false
}

func (t *Task) MarkAsDone(e error) {
	t.Active = false
	t.Done = true
	t.Error = e
}

func (t *Task) Reset() {
	if t.Active {
		return
	}

	t.Active = false
	t.Done = false
	t.Error = nil
}

func (t *Task) SetDownloadProgress(p int8) {
	t.DownloadProgress = p
	t.ch <- p
}

func (t *Task) Run(pwc *playwrightWrapper.PwContainer) error {
	err := downloader.Download(
		t.UrlSlug,
		&pwc.Browser,
		&t.DownloadProgress,
		t.ServiceNumber,
	)

	return err
}
