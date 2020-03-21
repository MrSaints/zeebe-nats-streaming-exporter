package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"gopkg.in/lucsky/cuid.v1"

	schema "github.com/mrsaints/zeebe-nats-streaming-exporter/indexers/exporter_protocol"
)

type Workflow struct {
	ID        string `json:"id" db:"id"`
	Key       int64  `json:"key" db:"key"`
	ProcessID string `json:"process_id" db:"process_id"`
	Version   uint32 `json:"version" db:"version"`
	Resource  string `json:"resource" db:"resource"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IndexedAt time.Time `json:"indexed_at" db:"indexed_at"`
}

func (w Workflow) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

type Workflows []Workflow

func (w Workflows) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

func (w *Workflow) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (w *Workflow) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (w *Workflow) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func NewWorkflowsFromProto(r *schema.DeploymentRecord) ([]Workflow, error) {
	deployedResources := map[string]*schema.DeploymentRecord_Resource{}
	for _, deployedResource := range r.GetResources() {
		deployedResources[deployedResource.GetResourceName()] = deployedResource
	}

	var workflows []Workflow
	for _, deployedWorkflow := range r.GetDeployedWorkflows() {
		relatedResource := deployedResources[deployedWorkflow.GetResourceName()]
		newWorkflow := Workflow{
			ID:        cuid.New(),
			Key:       deployedWorkflow.GetWorkflowKey(),
			ProcessID: deployedWorkflow.GetBpmnProcessId(),
			Version:   uint32(deployedWorkflow.GetVersion()),
			Resource:  string(relatedResource.GetResource()),
			CreatedAt: time.Unix(0, r.GetMetadata().GetTimestamp()*int64(time.Millisecond)),
			IndexedAt: time.Now().UTC(),
		}
		workflows = append(workflows, newWorkflow)
	}

	return workflows, nil
}
