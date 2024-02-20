package tekton

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPipelineVariables(t *testing.T) {
	type args struct {
		ctx         context.Context
		pipelineRun *tknv1.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Test with no params",
			args: args{
				ctx: context.Background(),
				pipelineRun: &tknv1.PipelineRun{
					Spec: tknv1.PipelineRunSpec{
						Params: []tknv1.Param{},
					},
					Status: tknv1.PipelineRunStatus{
						PipelineRunStatusFields: tknv1.PipelineRunStatusFields{
							PipelineSpec: &tknv1.PipelineSpec{
								Params: []tknv1.ParamSpec{},
							},
						},
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "Test with params",
			args: args{
				ctx: context.Background(),
				pipelineRun: &tknv1.PipelineRun{
					Spec: tknv1.PipelineRunSpec{
						Params: []tknv1.Param{
							{
								Name:  "param1",
								Value: tknv1.ParamValue{Type: tknv1.ParamTypeString, StringVal: "value1"},
							},
						},
					},
					Status: tknv1.PipelineRunStatus{
						PipelineRunStatusFields: tknv1.PipelineRunStatusFields{
							PipelineSpec: &tknv1.PipelineSpec{
								Params: []tknv1.ParamSpec{
									{
										Name:    "param2",
										Default: &tknv1.ParamValue{Type: tknv1.ParamTypeString, StringVal: "default2"},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]string{
				"param1": "value1",
				"param2": "default2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPipelineVariables(tt.args.ctx, tt.args.pipelineRun); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPipelineVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLabelsWithPrefix(t *testing.T) {
	type args struct {
		pipelineRun *tknv1.PipelineRun
		prefix      string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Test with no labels",
			args: args{
				pipelineRun: &tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{},
					},
				},
				prefix: "prefix",
			},
			want: map[string]string{},
		},
		{
			name: "Test with labels without prefix",
			args: args{
				pipelineRun: &tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"label1": "value1",
							"label2": "value2",
						},
					},
				},
				prefix: "prefix",
			},
			want: map[string]string{},
		},
		{
			name: "Test with labels with prefix",
			args: args{
				pipelineRun: &tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"prefix/label1": "value1",
							"label2":        "value2",
						},
					},
				},
				prefix: "prefix",
			},
			want: map[string]string{
				"label1": "value1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLabelsWithPrefix(tt.args.pipelineRun, tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLabelsWithPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPipelineName(t *testing.T) {
	type args struct {
		pipelineRun         tknv1.PipelineRun
		replacementVariable map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with pipeline name in replacement variable",
			args: args{
				pipelineRun:         tknv1.PipelineRun{},
				replacementVariable: map[string]string{"{original-prname}": "pipeline1"},
			},
			want: "pipeline1",
		},
		{
			name: "Test with pipeline name in pipeline reference",
			args: args{
				pipelineRun: tknv1.PipelineRun{
					Spec: tknv1.PipelineRunSpec{
						PipelineRef: &tknv1.PipelineRef{
							Name: "pipeline2",
						},
					},
				},
				replacementVariable: map[string]string{},
			},
			want: "pipeline2",
		},
		{
			name: "Test with pipeline name in resolver reference",
			args: args{

				pipelineRun: tknv1.PipelineRun{
					Spec: tknv1.PipelineRunSpec{
						PipelineRef: &tknv1.PipelineRef{
							ResolverRef: tknv1.ResolverRef{
								Params: []tknv1.Param{
									{
										Name: "name",
										Value: tknv1.ParamValue{
											Type:      tknv1.ParamTypeString,
											StringVal: "pipeline3",
										},
									},
								},
							},
						},
					},
				},
				replacementVariable: map[string]string{},
			},
			want: "pipeline3",
		},
		{
			name: "Test with no built-in pipeline name",
			args: args{
				pipelineRun: tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Name: "base-rhel-9.3-z76kp",
					},
					Spec: tknv1.PipelineRunSpec{
						PipelineSpec: &tknv1.PipelineSpec{
							DisplayName: "base-rhel-9.3",
						},
					},
				},
				replacementVariable: map[string]string{},
			},
			want: "base-rhel-9.3",
		},
		{
			name: "Test with no built-in no dash in name",
			args: args{
				pipelineRun: tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1234566",
					},
					Spec: tknv1.PipelineRunSpec{
						PipelineSpec: &tknv1.PipelineSpec{
							DisplayName: "test1234566",
						},
					},
				},
				replacementVariable: map[string]string{},
			},
			want: "test1234566",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPipelineName(&tt.args.pipelineRun, tt.args.replacementVariable); got != tt.want {
				t.Errorf("GetPipelineName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAttributes(t *testing.T) {
	testLogger := zaptest.NewLogger(t)
	log := zapr.NewLogger(testLogger)
	type args struct {
		ctx         context.Context
		pipelineRun *tknv1.PipelineRun
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Test with no attributes annotation",
			args: args{
				ctx: context.Background(),
				pipelineRun: &tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "Test with invalid attributes annotation",
			args: args{
				ctx: context.Background(),
				pipelineRun: &tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							tektonobserver.AttributesAnnotation: "{",
						},
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "Test with valid attributes annotation",
			args: args{
				ctx: context.Background(),
				pipelineRun: &tknv1.PipelineRun{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							tektonobserver.AttributesAnnotation: `{"key1":"value1","key2":"value2"}`,
						},
					},
				},
			},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := utils.NewFakeClient()
			evtEmit := events.NewEventEmitter(client, &log, "")

			if got := GetAttributes(tt.args.ctx, tt.args.pipelineRun, evtEmit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
