package com.fyianlai.zeebe.exporter.nats

import io.nats.streaming.NatsStreaming
import io.nats.streaming.Options
import io.nats.streaming.StreamingConnection
import io.zeebe.exporter.api.Exporter
import io.zeebe.exporter.api.context.Context
import io.zeebe.exporter.api.context.Controller
import io.zeebe.exporter.proto.RecordTransformer
import io.zeebe.protocol.record.Record
import java.util.*
import org.slf4j.Logger

class NatsStreamingExporter : Exporter {
    private lateinit var logger: Logger
    private lateinit var configuration: NatsStreamingExporterConfiguration

    private lateinit var streamingConnection: StreamingConnection
    private lateinit var controller: Controller

    private lateinit var transformer: (Record<*>) -> ByteArray

    override fun configure(context: Context) {
        this.logger = context.logger
        this.configuration = context.configuration.instantiate(NatsStreamingExporterConfiguration::class.java)

        if (this.configuration.format.equals("proto", ignoreCase = true)) {
            this.transformer = ::recordToProto
        } else if (this.configuration.format.equals("json", ignoreCase = true)) {
            this.transformer = ::recordToJson
        } else {
            throw IllegalArgumentException("Expected `format` to be either configured to 'proto' or 'json', got: ${this.configuration.format}")
        }

        this.logger.debug("Exporter configured with ${this.configuration}")
    }

    override fun open(controller: Controller) {
        this.controller = controller

        val streamingOptions = Options.Builder().natsUrl(this.configuration.serverUrl)
            .maxPubAcksInFlight(this.configuration.maxPubAcksInFlight).build()

        val uuid = UUID.randomUUID()
        val clientId = "${this.configuration.clientIdPrefix}$uuid"

        this.logger.debug("Connecting to ${this.configuration.serverUrl} with $clientId")

        this.streamingConnection = NatsStreaming.connect(this.configuration.clusterId, clientId, streamingOptions)

        this.logger.info("Exporter opened")
    }

    override fun close() {
        this.streamingConnection.close()

        this.logger.debug("Exporter closed")
    }

    override fun export(record: Record<*>) {
        val message = this.transformer(record)
        this.streamingConnection.publish("zeebe", message)
        this.controller.updateLastExportedRecordPosition(record.position)

        this.logger.debug("Exported record $record")
    }

    companion object {
        fun recordToProto(record: Record<*>): ByteArray {
            return RecordTransformer.toGenericRecord(record).toByteArray()
        }

        fun recordToJson(record: Record<*>): ByteArray {
            return record.toJson().toByteArray()
        }
    }
}
