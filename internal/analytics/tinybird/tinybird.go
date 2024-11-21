package tinybird

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/jonashiltl/openchangelog/internal/analytics"
	"github.com/jonashiltl/openchangelog/internal/xlog"
	"github.com/olivere/ndjson"
)

const (
	data_source            = "analytics_events"
	events_api             = "https://api.tinybird.co/v0/events"
	default_flush_interval = 10 * time.Second
	default_batch_size     = 50
)

type TinybirdOptions struct {
	AccessToken string
}

func New(opts TinybirdOptions) analytics.Emitter {
	b := &bird{
		buffer:        make([]analytics.Event, 0),
		flushInterval: default_flush_interval,
		batchSize:     default_batch_size,
		client:        &http.Client{Timeout: time.Second * 10},
		opts:          opts,
	}
	go b.startFlusher()
	return b
}

type bird struct {
	buffer        []analytics.Event
	mutex         sync.Mutex
	flushInterval time.Duration
	batchSize     int
	client        *http.Client
	opts          TinybirdOptions
}

func (b *bird) Emit(e analytics.Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.buffer = append(b.buffer, e)
	if len(b.buffer) >= b.batchSize {
		b.flush()
	}
}

func (b *bird) startFlusher() {
	ticker := time.NewTicker(b.flushInterval)
	defer ticker.Stop()

	for range ticker.C {
		b.mutex.Lock()
		b.flush()
		b.mutex.Unlock()
	}
}

func (b *bird) flush() {
	if len(b.buffer) == 0 {
		return
	}

	batch := b.buffer
	b.buffer = make([]analytics.Event, 0)

	go b.sendBatch(batch)
}

func (b *bird) sendBatch(events []analytics.Event) error {
	url := fmt.Sprintf("%s?name=%s", events_api, data_source)

	var buf bytes.Buffer
	writer := ndjson.NewWriter(&buf)
	for _, event := range events {
		if err := writer.Encode(event); err != nil {
			return err
		}
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		slog.Error("failed create new analytics request to tinybird", xlog.ErrAttr(err))
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.opts.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		slog.Error("failed to send events to tinybird", xlog.ErrAttr(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > http.StatusAccepted {
		slog.Error("received error status from tinybird", slog.String("status", resp.Status))
		return fmt.Errorf("received error status from tinybird: %s", resp.Status)
	}

	return nil
}
