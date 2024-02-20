package tekton

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kcloutie/tekton-observer/internal/tektonobserver"
	"github.com/kcloutie/tekton-observer/pkg/events"
	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"go.uber.org/zap/zapcore"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PacLabelPrefix = "pipelinesascode.tekton.dev"
)

type PipelineRunData struct {
	RawPipelineRun  *tknv1.PipelineRun `json:"rawPipelineRun,omitempty" yaml:"rawPipelineRun,omitempty"`
	VariableValues  map[string]string  `json:"variables,omitempty" yaml:"variables,omitempty"`
	PacLabels       map[string]string  `json:"pacLabels,omitempty" yaml:"pacLabels,omitempty"`
	Namespace       string             `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	PipelineRunName string             `json:"pipelineRunName,omitempty" yaml:"pipelineRunName,omitempty"`
	PipelineName    string             `json:"pipelineName,omitempty" yaml:"pipelineName,omitempty"`
	StartTime       *metav1.Time       `json:"startTime,omitempty" yaml:"startTime,omitempty"`
	CompletionTime  *metav1.Time       `json:"completionTime,omitempty" yaml:"completionTime,omitempty"`
	TotalTime       *string            `json:"totalTime,omitempty" yaml:"totalTime,omitempty"`
	Attributes      map[string]string  `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	// PipelineStatus     string             `json:"pipelineStatus,omitempty" yaml:"pipelineStatus,omitempty"`
}

func GetPipelineRunData(ctx context.Context, pipelineRun *tknv1.PipelineRun, eventEmitter *events.EventEmitter) (*PipelineRunData, error) {
	variables := GetPipelineVariables(ctx, pipelineRun)
	pacLabels := GetLabelsWithPrefix(pipelineRun, PacLabelPrefix)

	return &PipelineRunData{
		RawPipelineRun:  pipelineRun,
		VariableValues:  variables,
		PacLabels:       pacLabels,
		PipelineRunName: pipelineRun.Name,
		PipelineName:    GetPipelineName(pipelineRun, pacLabels),
		StartTime:       pipelineRun.Status.StartTime,
		CompletionTime:  pipelineRun.Status.CompletionTime,
	}, nil

}

func GetPipelineVariables(ctx context.Context, pipelineRun *tknv1.PipelineRun) map[string]string {
	variables := make(map[string]string)
	for _, param := range pipelineRun.Spec.Params {
		variables[param.Name] = param.Value.StringVal
	}

	if pipelineRun.Status.PipelineSpec != nil {
		for _, item := range pipelineRun.Status.PipelineSpec.Params {
			_, ok := variables[item.Name]
			if !ok {
				if item.Default != nil {
					variables[item.Name] = item.Default.StringVal
				}
			}
		}
	}

	return variables
}

func GetLabelsWithPrefix(pipelineRun *tknv1.PipelineRun, prefix string) map[string]string {
	labels := make(map[string]string)
	for key, value := range pipelineRun.Labels {
		if strings.HasPrefix(key, prefix) {
			labels[strings.ReplaceAll(key, prefix+"/", "")] = value
		}
	}
	return labels
}

func GetTotalTime(startTime *metav1.Time, completionTime *metav1.Time) string {
	totalTime := "Unknown"
	if startTime != nil && completionTime != nil && !startTime.IsZero() && !completionTime.IsZero() {
		totalTime = completionTime.Sub(startTime.Time).String()
	}
	return totalTime
}

func GetPipelineName(pipelineRun *tknv1.PipelineRun, variables map[string]string) string {
	pipelineRefName := "cannot-determine-pipeline-name"

	//PAC pipeline name
	val, ok := variables["{original-prname}"]
	if ok {
		return val
	}

	if pipelineRun.Spec.PipelineRef != nil {
		if pipelineRun.Spec.PipelineRef.Name != "" {
			return pipelineRun.Spec.PipelineRef.Name
		}
		if pipelineRun.Spec.PipelineRef.ResolverRef.Params != nil {
			for _, item := range pipelineRun.Spec.PipelineRef.ResolverRef.Params {
				if item.Name == "name" {
					return item.Value.StringVal
				}
			}
		}

	}
	if pipelineRun.Spec.PipelineSpec != nil {
		pipelineRefName := pipelineRun.GetName()
		lastChar := strings.LastIndex(pipelineRefName, "-")
		if lastChar > -1 {
			return pipelineRefName[:lastChar]
		}
		return pipelineRefName
	}

	return pipelineRefName
}

func GetAttributes(ctx context.Context, pipelineRun *tknv1.PipelineRun, eventEmitter *events.EventEmitter) map[string]string {
	attributes := make(map[string]string)
	rawAttr, exists := pipelineRun.Annotations[tektonobserver.AttributesAnnotation]
	if !exists {
		return attributes
	}
	err := json.Unmarshal([]byte(rawAttr), &attributes)
	if err != nil {
		eventEmitter.EmitMessagePipelineRun(ctx, pipelineRun, zapcore.ErrorLevel, "Attributes", fmt.Sprintf("Error unmarshalling the contents of the '%s' annotation. Ensure the contents are valid JSON", tektonobserver.AttributesAnnotation))
	}
	return attributes
}
