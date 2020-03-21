package main

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/rs/zerolog/log"

	schema "github.com/mrsaints/zeebe-nats-streaming-exporter/indexers/exporter_protocol"
	"github.com/mrsaints/zeebe-nats-streaming-exporter/indexers/sql/models"
)

const (
	DEPLOYMENT_PARTITION = 1
	CREATED_INTENT       = "CREATED"
)

type PopImporter struct {
	Connection *pop.Connection
}

func (p *PopImporter) deploymentRecordHandler(d schema.DeploymentRecord) error {
	intent := d.GetMetadata().GetIntent()
	partitionId := d.GetMetadata().GetPartitionId()

	if intent != CREATED_INTENT || partitionId != DEPLOYMENT_PARTITION {
		log.Debug().
			Int32("PartitionId", d.GetMetadata().GetPartitionId()).
			Str("Intent", d.GetMetadata().GetIntent()).
			Msg("Skipping import of deployment record.")
		return nil
	}

	workflows, err := models.NewWorkflowsFromProto(&d)
	if err != nil {
		return fmt.Errorf("failed to convert proto to workflows: %w", err)
	}

	for _, workflow := range workflows {
		log.Debug().
			Int64("Key", workflow.Key).
			Str("ProcessID", workflow.ProcessID).
			Uint32("Version", workflow.Version).
			Msg("Storing workflow.")

		var wf models.Workflow
		totalExistingWorkflows, err := p.Connection.Where("process_id = ?", workflow.ProcessID).Where("version = ?", workflow.Version).Count(&wf)
		if err != nil {
			return fmt.Errorf("failed to find existing workflows (%s@%s): %w", workflow.ProcessID, workflow.Version, err)
		}
		if totalExistingWorkflows > 0 {
			continue
		}

		if err := p.Connection.Create(&workflow); err != nil {
			return fmt.Errorf("failed to store workflows (%s@%s): %w", workflow.ProcessID, workflow.Version, err)
		}
	}

	return nil
}

func (p *PopImporter) workflowInstanceRecordHandler(w schema.WorkflowInstanceRecord) error {
	if models.IsWorkflowInstanceElementRecord(&w) {
		log.Debug().
			Int64("Key", w.GetMetadata().GetKey()).
			Int64("WorkflowInstanceKey", w.GetWorkflowInstanceKey()).
			Str("ElementId", w.GetElementId()).
			Msg("Workflow instance record is for an element.")

		workflowInstanceElement, err := models.NewWorkflowInstanceElementFromProto(&w)
		if err != nil {
			return fmt.Errorf("failed to convert proto to workflow instance element (%s@%s): %w", w.GetElementId(), w.GetWorkflowInstanceKey(), err)
		}

		log.Debug().
			Int64("Key", workflowInstanceElement.Key).
			Uint32("PartitionID", workflowInstanceElement.PartitionID).
			Int64("Position", workflowInstanceElement.Position).
			Int64("WorkflowKey", workflowInstanceElement.WorkflowKey).
			Int64("WorkflowInstanceKey", workflowInstanceElement.WorkflowInstanceKey).
			Str("Intent", workflowInstanceElement.Intent).
			Str("ElementID", workflowInstanceElement.ElementID).
			Str("ElementType", workflowInstanceElement.ElementType).
			Int64("FlowScopeKey", workflowInstanceElement.FlowScopeKey).
			Msg("Storing workflow instance element.")

		if err := p.Connection.Create(&workflowInstanceElement); err != nil {
			return fmt.Errorf("failed to store workflow instance element (%s@%s): %w", w.GetElementId(), w.GetWorkflowInstanceKey(), err)
		}
	} else {
		log.Debug().
			Int64("Key", w.GetMetadata().GetKey()).
			Int64("WorkflowInstanceKey", w.GetWorkflowInstanceKey()).
			Msg("Workflow instance record is for the instance.")

		workflowInstance, err := models.NewWorkflowInstanceFromProto(&w)
		if err != nil {
			return fmt.Errorf("failed to convert proto to workflow instance (%s): %w", w.GetWorkflowInstanceKey(), err)
		}

		log.Debug().
			Int64("Key", workflowInstance.Key).
			Uint32("PartitionID", workflowInstance.PartitionID).
			Int64("WorkflowKey", workflowInstance.WorkflowKey).
			Str("ProcessID", workflowInstance.ProcessID).
			Uint32("Version", workflowInstance.Version).
			Str("Intent", workflowInstance.Intent).
			Msg("Storing workflow instance.")

		if err := p.Connection.Create(&workflowInstance); err != nil {
			return fmt.Errorf("failed to store workflow instance (%s): %w", w.GetWorkflowInstanceKey(), err)
		}
	}

	return nil
}
