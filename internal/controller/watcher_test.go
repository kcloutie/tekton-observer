package controller

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kcloutie/tekton-observer/internal/tektonobserver"
	"github.com/kcloutie/tekton-observer/pkg/events"
	"github.com/kcloutie/tekton-observer/test/utils"
	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestTektonObservationReconciler_findConfigsFromPipelineRun(t *testing.T) {
	testLogger := zaptest.NewLogger(t)
	log := zapr.NewLogger(testLogger)
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx         context.Context
		pipelineRun client.Object
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		want                  []reconcile.Request
		wantProcessAnnotation string
	}{
		{
			name: "Test with pipelineRun not found",
			fields: fields{
				Client: utils.NewFakeClient(),
				Scheme: scheme.Scheme,
			},
			args: args{
				ctx:         context.Background(),
				pipelineRun: utils.NewPipelineRun("test-namespace", "test-name", map[string]string{}, false),
			},
			want: []reconcile.Request{},
		},

		{
			name: "Test with pipelineRun found and observation not found",
			fields: fields{
				Client: utils.NewFakeClient(utils.NewPipelineRun("test-namespace", "test-name", map[string]string{}, false)),
				Scheme: scheme.Scheme,
			},
			args: args{
				ctx:         context.Background(),
				pipelineRun: utils.NewPipelineRun("test-namespace", "test-name", map[string]string{}, false),
			},
			want: []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: "test-namespace",
						Name:      tektonobserver.ObservationCrdName,
					},
				},
			},
			wantProcessAnnotation: tektonobserver.ProcessingState,
		},
		{
			name: "Test with pipelineRun already started not done",
			fields: fields{
				Client: utils.NewFakeClient(utils.NewPipelineRun("test-namespace", "test-name", map[string]string{
					tektonobserver.PipelineProcessingStateAnnotation: tektonobserver.ProcessingStartState,
				}, false)),
				Scheme: scheme.Scheme,
			},
			args: args{
				ctx: context.Background(),
				pipelineRun: utils.NewPipelineRun("test-namespace", "test-name", map[string]string{
					tektonobserver.PipelineProcessingStateAnnotation: tektonobserver.ProcessingStartState,
				}, false),
			},
			want:                  []reconcile.Request{},
			wantProcessAnnotation: tektonobserver.ProcessingStartState,
		},
		{
			name: "Test with pipelineRun already started done",
			fields: fields{
				Client: utils.NewFakeClient(utils.NewPipelineRun("test-namespace", "test-name", map[string]string{
					tektonobserver.PipelineProcessingStateAnnotation: tektonobserver.ProcessingStartState,
				}, true)),
				Scheme: scheme.Scheme,
			},
			args: args{
				ctx: context.Background(),
				pipelineRun: utils.NewPipelineRun("test-namespace", "test-name", map[string]string{
					tektonobserver.PipelineProcessingStateAnnotation: tektonobserver.ProcessingStartState,
				}, true),
			},
			want: []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: "test-namespace",
						Name:      tektonobserver.ObservationCrdName,
					},
				},
			},
			wantProcessAnnotation: tektonobserver.ProcessingStartState,
		},
		{
			name: "Test with pipelineRun already processed",
			fields: fields{
				Client: utils.NewFakeClient(utils.NewPipelineRun("test-namespace", "test-name", map[string]string{
					tektonobserver.PipelineProcessingStateAnnotation: tektonobserver.ProcessingCompleteState,
				}, false)),
				Scheme: scheme.Scheme,
			},
			args: args{
				ctx: context.Background(),
				pipelineRun: utils.NewPipelineRun("test-namespace", "test-name", map[string]string{
					tektonobserver.PipelineProcessingStateAnnotation: tektonobserver.ProcessingCompleteState,
				}, false),
			},
			want:                  []reconcile.Request{},
			wantProcessAnnotation: tektonobserver.ProcessingCompleteState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TektonObservationReconciler{
				Client:       tt.fields.Client,
				Scheme:       tt.fields.Scheme,
				EventEmitter: events.NewEventEmitter(tt.fields.Client, &log, ""),
			}

			got := r.findConfigsFromPipelineRun(tt.args.ctx, tt.args.pipelineRun)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TektonObservationReconciler.findConfigsFromPipelineRun() = %v, want %v", got, tt.want)
			}
			if tt.wantProcessAnnotation != "" {
				pr := &tknv1.PipelineRun{}
				err := r.Get(tt.args.ctx, types.NamespacedName{
					Namespace: tt.args.pipelineRun.GetNamespace(),
					Name:      tt.args.pipelineRun.GetName(),
				}, pr)
				if err != nil {
					t.Errorf("TektonObservationReconciler.findConfigsFromPipelineRun() = %v, want %v", got, tt.want)
				}
				if pr.Annotations[tektonobserver.PipelineProcessingStateAnnotation] != tt.wantProcessAnnotation {
					t.Errorf("TektonObservationReconciler.findConfigsFromPipelineRun() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
