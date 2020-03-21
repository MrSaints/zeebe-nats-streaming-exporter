package main

import (
	"os"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"

	"github.com/mrsaints/zeebe-nats-streaming-exporter/indexers/importer"
)

func main() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Setting up importer ...")

	log.Info().Msg("Connecting to database ...")
	conn, err := pop.Connect("development")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to connect to database.")
	}

	popImporter := &PopImporter{conn}
	importer := importer.NewImporter().WithDeploymentRecordHandler(
		popImporter.deploymentRecordHandler,
	).WithWorkflowInstanceRecordHandler(
		popImporter.workflowInstanceRecordHandler,
	)

	log.Info().Msg("Connecting to NATS ...")
	nc, err := nats.Connect("nats://localhost")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to connect to NATS.")
	}
	defer nc.Close()

	log.Info().Msg("Setting up STAN ...")
	sc, err := stan.Connect("zeebe", "zeebe-nats-streaming-sql-indexer", stan.NatsConn(nc))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to set-up STAN.")
	}
	defer sc.Close()

	log.Info().Msg("Subscribing ...")
	aw, err := time.ParseDuration("120s")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to parse duration for ACK.")
	}
	qsub, err := sc.QueueSubscribe("zeebe", "zeebe-nats-streaming-sql-indexer", func(m *stan.Msg) {
		if err := importer.Import(m.Data); err != nil {
			log.Error().
				Uint64("Sequence", m.Sequence).
				Str("Subject", m.Subject).
				Int64("Timestamp", m.Timestamp).
				Err(err).
				Msg("Failed to handle message.")
		} else {
			m.Ack()
			log.Debug().
				Uint64("Sequence", m.Sequence).
				Str("Subject", m.Subject).
				Int64("Timestamp", m.Timestamp).
				Msg("Message handled, and acknowledged.")
		}
	}, stan.DurableName("zeebe-nats-streaming-sql-indexer"), stan.SetManualAckMode(), stan.AckWait(aw))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to subscribe to durable queue.")
	}
	defer qsub.Close()

	log.Info().Msg("Importer set-up completed. Waiting to receive messages ...")
	select {}
}
