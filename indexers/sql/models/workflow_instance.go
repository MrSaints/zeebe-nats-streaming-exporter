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

type WorkflowInstance struct {
	ID          string `json:"id" db:"id"`
	Key         int64  `json:"key" db:"key"`
	PartitionID uint32 `json:"partition_id" db:"partition_id"`
	WorkflowKey int64  `json:"workflow_key" db:"workflow_key"`
	ProcessID   string `json:"process_id" db:"process_id"`
	Version     uint32 `json:"version" db:"version"`

	Intent string `json:"intent" db:"intent"`

	ParentWorkflowInstanceKey int64 `json:"parent_workflow_instance_key" db:"parent_workflow_instance_key"`
	ParentElementInstanceKey  int64 `json:"parent_element_instance_key" db:"parent_element_instance_key"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IndexedAt time.Time `json:"indexed_at" db:"indexed_at"`
}

func (w WorkflowInstance) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

type WorkflowInstances []WorkflowInstance

func (w WorkflowInstances) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

func (w *WorkflowInstance) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (w *WorkflowInstance) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (w *WorkflowInstance) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func NewWorkflowInstanceFromProto(r *schema.WorkflowInstanceRecord) (WorkflowInstance, error) {
	if IsWorkflowInstanceElementRecord(r) {
		return WorkflowInstance{}, fmt.Errorf("workflow instance record is for an element, not the instance: %s != %s", r.GetMetadata().GetKey(), r.GetWorkflowInstanceKey())
	}

	intent := r.GetMetadata().GetIntent()

	return WorkflowInstance{
		ID:                        cuid.New(),
		Key:                       r.GetWorkflowInstanceKey(),
		PartitionID:               uint32(r.GetMetadata().GetPartitionId()),
		WorkflowKey:               r.GetWorkflowKey(),
		ProcessID:                 r.GetBpmnProcessId(),
		Version:                   uint32(r.GetVersion()),
		Intent:                    intent,
		ParentWorkflowInstanceKey: r.GetParentWorkflowInstanceKey(),
		ParentElementInstanceKey:  r.GetParentElementInstanceKey(),
		CreatedAt:                 time.Unix(0, r.GetMetadata().GetTimestamp()*int64(time.Millisecond)),
		IndexedAt:                 time.Now().UTC(),
	}, nil
}
