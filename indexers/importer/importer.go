package importer

import (
	"fmt"

	"github.com/rs/zerolog/log"

	proto "github.com/golang/protobuf/proto"
	ptypes "github.com/golang/protobuf/ptypes"

	schema "github.com/mrsaints/zeebe-nats-streaming-exporter/indexers/exporter_protocol"
)

type DeploymentRecordHandler func(d schema.DeploymentRecord) error
type WorkflowInstanceRecordHandler func(w schema.WorkflowInstanceRecord) error

type Importer struct {
	deploymentRecordHandler       DeploymentRecordHandler
	workflowInstanceRecordHandler WorkflowInstanceRecordHandler
}

func NewImporter() *Importer {
	return &Importer{
		deploymentRecordHandler:       defaultDeploymentRecordHandler,
		workflowInstanceRecordHandler: defaultWorkflowInstanceRecordHandler,
	}
}

func (i *Importer) WithDeploymentRecordHandler(deploymentRecordHandler DeploymentRecordHandler) *Importer {
	i.deploymentRecordHandler = deploymentRecordHandler
	return i
}

func (i *Importer) WithWorkflowInstanceRecordHandler(workflowInstanceRecordHandler WorkflowInstanceRecordHandler) *Importer {
	i.workflowInstanceRecordHandler = workflowInstanceRecordHandler
	return i
}

func (i *Importer) Import(d []byte) error {
	log.Debug().Msg("Importing generic record.")

	var r schema.Record
	if err := proto.Unmarshal(d, &r); err != nil {
		return fmt.Errorf("failed to unmarshal Record: %w", err)
	}

	log.Debug().
		Str("TypeUrl", r.GetRecord().GetTypeUrl()).
		Msg("Determining generic record type.")

	if ptypes.Is(r.GetRecord(), &schema.DeploymentRecord{}) {
		var deploymentRecord schema.DeploymentRecord
		if err := ptypes.UnmarshalAny(r.GetRecord(), &deploymentRecord); err != nil {
			return fmt.Errorf("failed to unmarshal DeploymentRecord: %w", err)
		}

		log.Debug().
			Int32("PartitionId", deploymentRecord.GetMetadata().GetPartitionId()).
			Int64("Position", deploymentRecord.GetMetadata().GetPosition()).
			Int64("Key", deploymentRecord.GetMetadata().GetKey()).
			Int64("Timestamp", deploymentRecord.GetMetadata().GetTimestamp()).
			Str("RecordType", deploymentRecord.GetMetadata().GetRecordType().String()).
			Str("Intent", deploymentRecord.GetMetadata().GetIntent()).
			Str("ValueType", deploymentRecord.GetMetadata().GetValueType().String()).
			Msg("Importing deployment record.")

		return i.deploymentRecordHandler(deploymentRecord)
	}

	if ptypes.Is(r.GetRecord(), &schema.WorkflowInstanceRecord{}) {
		var instanceRecord schema.WorkflowInstanceRecord
		if err := ptypes.UnmarshalAny(r.GetRecord(), &instanceRecord); err != nil {
			return fmt.Errorf("failed to unmarshal WorkflowInstanceRecord: %w", err)
		}

		log.Debug().
			Int32("PartitionId", instanceRecord.GetMetadata().GetPartitionId()).
			Int64("Position", instanceRecord.GetMetadata().GetPosition()).
			Int64("Key", instanceRecord.GetMetadata().GetKey()).
			Int64("Timestamp", instanceRecord.GetMetadata().GetTimestamp()).
			Str("RecordType", instanceRecord.GetMetadata().GetRecordType().String()).
			Str("Intent", instanceRecord.GetMetadata().GetIntent()).
			Str("ValueType", instanceRecord.GetMetadata().GetValueType().String()).
			Msg("Importing workflow instance record.")

		return i.workflowInstanceRecordHandler(instanceRecord)
	}

	log.Info().
		Str("TypeUrl", r.GetRecord().GetTypeUrl()).
		Msg("Generic record was not be handled in importer.")

	return nil
}

func defaultDeploymentRecordHandler(d schema.DeploymentRecord) error {
	log.Warn().
		Int32("PartitionId", d.GetMetadata().GetPartitionId()).
		Int64("Position", d.GetMetadata().GetPosition()).
		Int64("Key", d.GetMetadata().GetKey()).
		Int64("Timestamp", d.GetMetadata().GetTimestamp()).
		Str("RecordType", d.GetMetadata().GetRecordType().String()).
		Str("Intent", d.GetMetadata().GetIntent()).
		Str("ValueType", d.GetMetadata().GetValueType().String()).
		Msg("No deployment record handler attached to importer.")
	return nil
}

func defaultWorkflowInstanceRecordHandler(w schema.WorkflowInstanceRecord) error {
	log.Warn().
		Int32("PartitionId", w.GetMetadata().GetPartitionId()).
		Int64("Position", w.GetMetadata().GetPosition()).
		Int64("Key", w.GetMetadata().GetKey()).
		Int64("Timestamp", w.GetMetadata().GetTimestamp()).
		Str("RecordType", w.GetMetadata().GetRecordType().String()).
		Str("Intent", w.GetMetadata().GetIntent()).
		Str("ValueType", w.GetMetadata().GetValueType().String()).
		Msg("No workflow instance record handler attached to importer.")
	return nil
}
