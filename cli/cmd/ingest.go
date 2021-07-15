package cmd

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/ipc"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
)

var (
	ingestSocketLookupDelay = time.Millisecond * 500
	ingestReconnectDelay    = time.Millisecond * 250
)

type IngestCommand struct {
	db     *store.DB
	logger *logrus.Logger

	ipcRoot   string
	ipcChains []string
}

func NewIngestCommand(db *store.DB, logger *logrus.Logger, ipcRoot string, ipcChains []string) IngestCommand {
	return IngestCommand{
		db:     db,
		logger: logger,

		ipcRoot:   ipcRoot,
		ipcChains: ipcChains,
	}
}

func (cmd IngestCommand) Run() error {
	topics := make([]*model.RawMessageTopic, len(cmd.ipcChains))

	for idx, chain := range cmd.ipcChains {
		cmd.logger.Info("creating raw message topic:", chain)

		chunks := strings.SplitN(chain, "-", 3)

		topic, err := cmd.db.RawMessages.CreateOrFindTopic(chunks[2], chunks[0], chunks[1])
		if err != nil {
			return err
		}
		topics[idx] = topic
	}

	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(len(topics))

	for _, topic := range topics {
		go func(topic *model.RawMessageTopic) {
			cmd.startReadLoop(ctx, topic, cmd.messageWriter(topic))
		}(topic)
	}

	wg.Wait()
	return nil
}

func (cmd IngestCommand) messageWriter(topic *model.RawMessageTopic) ipc.WriterFn {
	return func(msg ipc.Message) error {
		err := cmd.db.RawMessages.CreateMessage(topic.ID, msg.Data, msg.Time)
		if err != nil {
			cmd.logger.WithError(err).WithField("topic", topic.Chain).Error("cant write message")
		}
		return err
	}
}

func (cmd IngestCommand) startReadLoop(ctx context.Context, topic *model.RawMessageTopic, writer ipc.WriterFn) {
	cmd.logger.
		WithField("chain", topic.Chain).
		WithField("vm", topic.VM).
		Info("starting IPC socket reader")

	for {
		socketPath := ipc.FindDesicionsSocketPath(topic.Chain, cmd.ipcRoot)
		if socketPath == "" {
			cmd.logger.WithField("chain", topic.Chain).Error("cant find socket file")
			time.Sleep(ingestSocketLookupDelay)
			continue
		}

		cmd.logger.WithField("file", socketPath).Info("connecting to socket file")
		socket, err := ipc.Dial(socketPath)
		if err != nil {
			cmd.logger.
				WithError(err).
				WithField("file", socketPath).
				Error("cant connect to socket file")

			time.Sleep(ingestReconnectDelay)
			continue
		}

		if err := socket.Start(writer); err != nil {
			log.Println("socket write error:", err)
		}

		socket.Close()
		time.Sleep(ingestReconnectDelay)
	}
}
