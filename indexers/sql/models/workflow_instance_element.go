package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"gopkg.in/lucsky/cuid.v1"

	schema "github.com/mrsaints/zeebe-nats-streaming-exporter/indexers/exporter_protocol"
)

type WorkflowInstanceElement struct {
	ID                  string `json:"id" db:"id"`
	Key                 int64  `json:"key" db:"key"`
	PartitionID         uint32 `json:"partition_id" db:"partition_id"`
	Position            int64  `json:"position" db:"position"`
	WorkflowKey         int64  `json:"workflow_key" db:"workflow_key"`
	WorkflowInstanceKey int64  `json:"workflow_instance_key" db:"workflow_instance_key"`

	Intent       string `json:"intent" db:"intent"`
	ElementID    string `json:"element_id" db:"element_id"`
	ElementType  string `json:"element_type" db:"element_type"`
	FlowScopeKey int64  `json:"flow_scope_key" db:"flow_scope_key"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IndexedAt time.Time `json:"indexed_at" db:"indexed_at"`
}

func (w WorkflowInstanceElement) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

type WorkflowInstanceElements []WorkflowInstanceElement

func (w WorkflowInstanceElements) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

func (w *WorkflowInstanceElement) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (w *WorkflowInstanceElement) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (w *WorkflowInstanceElement) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func NewWorkflowInstanceElementFromProto(r *schema.WorkflowInstanceRecord) (WorkflowInstanceElement, error) {
	if !IsWorkflowInstanceElementRecord(r) {
		return WorkflowInstanceElement{}, fmt.Errorf("workflow instance record is for the instance, not an element: %s == %s", r.GetMetadata().GetKey(), r.GetWorkflowInstanceKey())
	}

	return WorkflowInstanceElement{
		ID:                  cuid.New(),
		Key:                 r.GetWorkflowInstanceKey(),
		PartitionID:         uint32(r.GetMetadata().GetPartitionId()),
		Position:            r.GetMetadata().GetPosition(),
		WorkflowKey:         r.GetWorkflowKey(),
		WorkflowInstanceKey: r.GetWorkflowInstanceKey(),
		Intent:              r.GetMetadata().GetIntent(),
		ElementID:           r.GetElementId(),
		ElementType:         r.GetBpmnElementType().String(),
		FlowScopeKey:        r.GetFlowScopeKey(),
		CreatedAt:           time.Unix(0, r.GetMetadata().GetTimestamp()*int64(time.Millisecond)),
		IndexedAt:           time.Now().UTC(),
	}, nil
}
