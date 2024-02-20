package controller

import (
	"context"
	"fmt"

	obsv1 "github.com/kcloutie/tekton-observer/api/tektonobserver/v1"
	"github.com/kcloutie/tekton-observer/internal/tektonobserver"
	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *TektonObservationReconciler) findConfigsFromPipelineRun(ctx context.Context, pipelineRun client.Object) []reconcile.Request {

	pipelineRunObject := &tknv1.PipelineRun{}
	pipelineRunName := types.NamespacedName{
		Namespace: pipelineRun.GetNamespace(),
		Name:      pipelineRun.GetName(),
	}

	log := log.FromContext(ctx).WithValues("namespace", pipelineRunName.Namespace, "PipelineRun", pipelineRunName.Name, "PipelineUid", pipelineRun.GetUID(), "clusterName", tektonobserver.ControllerConfiguration.GetClusterName())

	if err := r.Get(ctx, pipelineRunName, pipelineRunObject); err != nil {
		// Its likely the pipelineRun was deleted if it does not exist...so just return
		log.V(2).Info("Cannot find PipelineRun...it was likely deleted")
		return []reconcile.Request{}
	}

	observationNamespacedName := types.NamespacedName{
		Namespace: pipelineRun.GetNamespace(),
		Name:      tektonobserver.ObservationCrdName,
	}

	observation := obsv1.TektonObservation{}

	err := r.Get(ctx, observationNamespacedName, &observation)
	if err != nil {
		err = r.Create(ctx, &observation, &client.CreateOptions{})
		if err != nil {
			mess := "Failed to create TektonObservation CR"
			log.V(0).Error(err, mess, "observationName", observationNamespacedName)
			r.EventEmitter.EmitMessagePipelineRun(ctx, pipelineRunObject, zapcore.ErrorLevel, "Create TektonObservation CR", fmt.Sprintf("%v. %v", mess, err))
		}
	}

	state, processedStateLabelFound := pipelineRunObject.Annotations[tektonobserver.PipelineProcessingStateAnnotation]

	if processedStateLabelFound && state == tektonobserver.ProcessingCompleteState {
		log.V(3).Info("PipelineRun has already been processed...skipping")
		return []reconcile.Request{}
	}

	if !processedStateLabelFound {
		err = r.updatePipelineRunAnnotation(ctx, tektonobserver.PipelineProcessingStateAnnotation, tektonobserver.ProcessingState, *pipelineRunObject, log)
		if err != nil {
			log.Error(err, "Failed to update PipelineRun label")
		}

	} else if state == tektonobserver.ProcessingStartState {
		if !pipelineRunObject.IsDone() {
			log.V(2).Info("PipelineRun is still running...skipping")
			return []reconcile.Request{}
		}
	}
	return []reconcile.Request{{NamespacedName: observationNamespacedName}}
}
