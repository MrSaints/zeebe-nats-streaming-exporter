package models

import (
	schema "github.com/mrsaints/zeebe-nats-streaming-exporter/indexers/exporter_protocol"
)

func IsWorkflowInstanceElementRecord(r *schema.WorkflowInstanceRecord) bool {
	return r.GetMetadata().GetKey() != r.GetWorkflowInstanceKey()
}
