package events

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	observerv1 "github.com/kcloutie/tekton-observer/api/tektonobserver/v1"
	"github.com/kcloutie/tekton-observer/internal/tektonobserver"
	"github.com/kcloutie/tekton-observer/pkg/formatting"
	"github.com/kcloutie/tekton-observer/pkg/metrics"
	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewEventEmitter(client client.Client, logger *logr.Logger, controllerInstance string) *EventEmitter {
	return &EventEmitter{
		client:             client,
		logger:             logger,
		controllerInstance: controllerInstance,
	}
}

type EventEmitter struct {
	client             client.Client
	logger             *logr.Logger
	controllerInstance string
}

func (e *EventEmitter) SetLogger(logger *logr.Logger) {
	e.logger = logger
}

func (e *EventEmitter) EmitMessage(ctx context.Context, tektonObserver *observerv1.TektonObservation, loggerLevel zapcore.Level, reason, message string) {
	if tektonObserver != nil {
		start := time.Now()
		event := makeEvent(tektonObserver, loggerLevel, reason, message, e.controllerInstance)
		err := e.client.Create(ctx, event, &client.CreateOptions{})

		if err != nil {
			e.logger.Info(fmt.Sprintf("Cannot create event: %s", err.Error()), "event", fmt.Sprintf("%+v", event))
		}
		duration := time.Since(start)
		metrics.EmitEventRequestTimeHistogram.WithLabelValues(fmt.Sprintf("%v", metrics.GetStatusCode(err))).Observe(duration.Seconds())
	}
}

func (e *EventEmitter) EmitMessagePipelineRun(ctx context.Context, pipelineRun *tknv1.PipelineRun, loggerLevel zapcore.Level, reason, message string) {
	if pipelineRun != nil {
		start := time.Now()
		event := makeEventPipelineRun(pipelineRun, "PipelineRun", loggerLevel, reason, message, e.controllerInstance)
		err := e.client.Create(ctx, event, &client.CreateOptions{})
		if err != nil {
			e.logger.Info(fmt.Sprintf("Cannot create event: %s", err.Error()), "event", fmt.Sprintf("%+v", event))
		}
		duration := time.Since(start)
		metrics.EmitEventRequestTimeHistogram.WithLabelValues(fmt.Sprintf("%v", metrics.GetStatusCode(err))).Observe(duration.Seconds())
	}
}

func makeEventPipelineRun(pipelineRun client.Object, kind string, loggerLevel zapcore.Level, reason, message string, controllerInstance string) *v1.Event {
	event := &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pipelineRun.GetName() + "-",
			Namespace:    pipelineRun.GetNamespace(),
			Labels: map[string]string{
				tektonobserver.ObserverNameAnnotation: formatting.CleanValueKubernetes(pipelineRun.GetName()),
			},
			Annotations: map[string]string{
				tektonobserver.ObserverNameAnnotation: pipelineRun.GetName(),
			},
		},
		Action:  reason,
		Message: message,
		Reason:  reason,
		Type:    v1.EventTypeWarning,
		InvolvedObject: v1.ObjectReference{
			APIVersion:      fmt.Sprintf("%v/%v", tektonobserver.GroupName, tektonobserver.V1Version),
			Kind:            kind,
			Namespace:       pipelineRun.GetNamespace(),
			Name:            pipelineRun.GetName(),
			UID:             pipelineRun.GetUID(),
			ResourceVersion: pipelineRun.GetResourceVersion(),
		},
		ReportingController: "kcloutie/tekton-observer",
		ReportingInstance:   controllerInstance,
		EventTime:           metav1.MicroTime{Time: time.Now()},
		Source: v1.EventSource{
			Component: "Tekton Observer",
		},
	}
	if loggerLevel == zap.InfoLevel {
		event.Type = v1.EventTypeNormal
	}

	return event
}

func makeEvent(tektonObserver *observerv1.TektonObservation, loggerLevel zapcore.Level, reason, message string, controllerInstance string) *v1.Event {
	event := &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: tektonObserver.Name + "-",
			Namespace:    tektonObserver.Namespace,
			Labels: map[string]string{
				tektonobserver.ObserverNameAnnotation: formatting.CleanValueKubernetes(tektonObserver.Name),
			},
			Annotations: map[string]string{
				tektonobserver.ObserverNameAnnotation: tektonObserver.Name,
			},
		},
		Action:  reason,
		Message: message,
		Reason:  reason,
		Type:    v1.EventTypeWarning,
		InvolvedObject: v1.ObjectReference{
			APIVersion:      fmt.Sprintf("%v/%v", tektonobserver.GroupName, tektonobserver.V1Version),
			Kind:            tektonobserver.TektonObservationKind,
			Namespace:       tektonObserver.Namespace,
			Name:            tektonObserver.Name,
			UID:             tektonObserver.UID,
			ResourceVersion: tektonObserver.ResourceVersion,
		},
		ReportingController: "kcloutie/tekton-observer",
		ReportingInstance:   controllerInstance,
		EventTime:           metav1.MicroTime{Time: time.Now()},
		Source: v1.EventSource{
			Component: "Tekton Observer",
		},
	}
	if loggerLevel == zap.InfoLevel {
		event.Type = v1.EventTypeNormal
	}

	return event
}

type EventMin struct {
	// The object that this event is about.
	InvolvedObject v1.ObjectReference `json:"involvedObject" protobuf:"bytes,2,opt,name=involvedObject"`

	// This should be a short, machine understandable string that gives the reason
	// for the transition into the object's current status.
	// TODO: provide exact specification for format.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`

	// A human-readable description of the status of this operation.
	// TODO: decide on maximum length.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`

	// Type of this event (Normal, Warning), new types could be added in the future
	// +optional
	Type string `json:"type,omitempty" protobuf:"bytes,9,opt,name=type"`

	// Time when this Event was first observed.
	// +optional
	EventTime metav1.MicroTime `json:"eventTime,omitempty" protobuf:"bytes,10,opt,name=eventTime"`

	Component string `json:"component,omitempty" protobuf:"bytes,1,opt,name=component"`
}

func NewEventMin(event v1.Event) EventMin {
	evt := EventMin{}
	evt.InvolvedObject = *event.InvolvedObject.DeepCopy()
	evt.Reason = event.Reason
	evt.Message = event.Message
	evt.Type = event.Type
	evt.EventTime = metav1.MicroTime(event.CreationTimestamp)
	evt.Component = event.Source.Component
	return evt
}
