package controller

import (
	"context"

	"github.com/go-logr/logr"
	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// func (r *TektonObservationReconciler) pipelineRunIsDone(log logr.Logger, pipelineRun *tknv1.PipelineRun) bool {

// 	if !pipelineRun.IsDone() {
// 		// We do not care about pipelines that have not finished running
// 		log.V(4).Info("PipelineRun has the status of unknown so it is likely still running, skipping")
// 		return false
// 	}

// 	log.V(2).Info("PipelineRun is done...lets process it!")
// 	return true
// }

func (r *TektonObservationReconciler) updatePipelineRunAnnotation(ctx context.Context, key, value string, pipelineRun tknv1.PipelineRun, log logr.Logger) error {
	updated := pipelineRun.DeepCopy()
	updated.Annotations[key] = value

	patch := client.MergeFrom(&pipelineRun)
	err := r.Patch(ctx, updated, patch)

	if err != nil {
		return err
	}
	return nil
}

// func (r *TektonObservationReconciler) updatePipelineRunLabel(ctx context.Context, label string, pipelineRun tknv1.PipelineRun, log logr.Logger) error {
// 	updated := pipelineRun.DeepCopy()
// 	updated.Labels[label] = fmt.Sprintf("%d", time.Now().UTC().Unix())

// 	patch := client.MergeFrom(&pipelineRun)
// 	err := r.Patch(ctx, updated, patch)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
