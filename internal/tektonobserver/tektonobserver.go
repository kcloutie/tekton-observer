package tektonobserver

const (
	ObservationCrdName                = "tekton-observer"
	GroupName                         = "observer.tkn.dev"
	PipelineProcessingStateAnnotation = GroupName + "/processing-state"
	ProcessingState                   = "processing"
	ProcessingCompleteState           = "complete"
	ProcessingStartState              = "started"
	// PipelineProcessedStartAnnotation    = GroupName + "/processed-start"
	// PipelineProcessedCompleteAnnotation = GroupName + "/processed-complete"
	AttributesAnnotation      = GroupName + "/attributes"
	ObserverNameAnnotation    = GroupName + "/name"
	PubSubProjectIDAnnotation = GroupName + "/pubsub-project-id"
	PubSubTopicIDAnnotation   = GroupName + "/pubsub-topic"

	V1Version             = "v1"
	TektonObservationKind = "TektonObservation"
)
