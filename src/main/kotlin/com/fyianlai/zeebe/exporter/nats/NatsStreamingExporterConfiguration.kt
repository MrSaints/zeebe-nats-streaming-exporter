package com.fyianlai.zeebe.exporter.nats

data class NatsStreamingExporterConfiguration(
    val serverUrl: String = "nats://localhost",
    val clusterId: String = "zeebe",
    val clientIdPrefix: String = "zeebe-exporter-",
    val channel: String = "zeebe",
    val maxPubAcksInFlight: Int = 10000,
    val format: String = "proto"
)