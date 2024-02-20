package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	requestTimesBuckets = []float64{.5, 1, 2.5, 5, 10}

	SendEmailRequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tknobs_send_email_request_duration_seconds",
		Help:    "Histogram of send email request time in seconds",
		Buckets: requestTimesBuckets,
	}, []string{"success"})

	EmitEventRequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tknobs_emit_event_request_duration_seconds",
		Help:    "Histogram of emit event request time in seconds",
		Buckets: requestTimesBuckets,
	}, []string{"success"})

	GithubRequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tknobs_github_request_duration_seconds",
		Help:    "Histogram of github request time in seconds",
		Buckets: requestTimesBuckets,
	}, []string{"route", "method", "status_code"})

	GoogleRequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tknobs_google_request_duration_seconds",
		Help:    "Histogram of google API request time in seconds",
		Buckets: requestTimesBuckets,
	}, []string{"route", "method", "success"})

	KubernetesRequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tknobs_kubernetes_request_duration_seconds",
		Help:    "Histogram of kubernetes API request time in seconds",
		Buckets: requestTimesBuckets,
	}, []string{"route", "method", "status_code"})

	WebexRequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tknobs_webex_request_duration_seconds",
		Help:    "Histogram of webex API request time in seconds",
		Buckets: requestTimesBuckets,
	}, []string{"route", "method", "status_code"})

	ProcessPipelineTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tknobs_process_pipeline_duration_seconds",
		Help:    "Histogram of the time it takes to process a pipeline in seconds",
		Buckets: requestTimesBuckets,
	}, []string{"executeStatus"})

	PipelineRunsProcessedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_processed_pipeline_runs_total",
			Help: "Number of pipeline runs processed",
		},
	)
	PipelineRunsStartedProcessingTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_started_processing_pipeline_runs_total",
			Help: "Number of pipeline runs that started processing",
		},
	)

	LogsSavedToGcsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_logs_saved_total",
			Help: "Number of pipeline run logs saved to a GCS bucket",
		},
	)

	LogsSavedToGcsSkippedDisabledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_logs_save_skipped_disabled_total",
			Help: "Number of pipeline run logs skipped being saved because it was disabled",
		},
	)

	LogsSavedToGcsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_logs_save_failed_total",
			Help: "Number of pipeline run logs that failed to be saved to a GCS bucket",
		},
	)

	WebexMessagesSentTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_webex_messages_sent_total",
			Help: "Number of webex notification messages sent",
		},
	)

	WebexMessagesSkippedDisabledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_webex_messages_skipped_disabled_total",
			Help: "Number of webex notification messages that were not sent because it was disabled",
		},
	)

	WebexMessagesSkippedSuccessTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_webex_messages_skipped_success_total",
			Help: "Number of webex notification messages that were not sent because the pipeline run was a success",
		},
	)

	WebexMessagesFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_webex_messages_failed_total",
			Help: "Number of webex notification messages that failed to be sent",
		},
	)

	GithubCommentSkippedDisabledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_comment_skipped_disabled_total",
			Help: "Number of github comments that were not created because it was disabled",
		},
	)

	GithubCommitCommentCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_commit_comment_created_total",
			Help: "Number of github commit comments created",
		},
	)

	GithubCommitCommentFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_commit_comment_failed_total",
			Help: "Number of github commit comments failed to create",
		},
	)

	GithubPrCommentCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_pr_comment_created_total",
			Help: "Number of github pr comments created",
		},
	)

	GithubPrCommentFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_pr_comment_failed_total",
			Help: "Number of github pr comments failed to create",
		},
	)

	GithubPrCommentSkippedNotPrTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_pr_comment_skipped_not_pr_total",
			Help: "Number of github pr comments skipped because it was not a pr",
		},
	)

	GithubStatusSkippedDisabledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_status_skipped_disabled_total",
			Help: "Number of github status was not set because it was disabled",
		},
	)

	GithubStatusFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_status_failed_total",
			Help: "Number of github status failed to be set",
		},
	)

	GithubStatusCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_status_created_total",
			Help: "Number of github status created",
		},
	)

	GithubStatusSkippedFailedExistedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_status_skipped_failed_existed_total",
			Help: "Number of github status was not set because an existing status check existed that was in a failed state",
		},
	)

	GithubDeploymentSkippedDisabledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_deployment_skipped_disabled_total",
			Help: "Number of github deployment was not set because it was disabled",
		},
	)

	GithubDeploymentFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_deployment_failed_total",
			Help: "Number of github deployment failed to be set",
		},
	)

	GithubDeploymentCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_deployment_created_total",
			Help: "Number of github deployment created",
		},
	)

	GithubDeploymentSkippedEventNotMatchTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_github_deployment_skipped_event_not_match_total",
			Help: "Number of github deployments that were not created because the event did not match",
		},
	)

	EmailSkippedDisabledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_email_skipped_disabled_total",
			Help: "Number of emails that were not sent because it was disabled",
		},
	)

	EmailFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_email_failed_total",
			Help: "Number of emails that failed to send",
		},
	)

	EmailCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_email_created_total",
			Help: "Number of emails that were sent",
		},
	)

	EmailSkippedSuccessTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_email_skipped_success_total",
			Help: "Number of emails that were not sent because the pipeline run was a success",
		},
	)

	PubSubSkippedDisabledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_pubsub_skipped_disabled_total",
			Help: "Number of pubsub messages that were not sent because it was disabled",
		},
	)

	PubSubSentTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_pubsub_sent_total",
			Help: "Number of pubsub messages sent",
		},
	)

	PubSubFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_pubsub_failed_total",
			Help: "Number of pubsub messages that failed to send",
		},
	)

	PubSubGlobalSentTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_pubsub_global_sent_total",
			Help: "Number of global pubsub messages sent",
		},
	)

	PubSubGlobalFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_pubsub_global_failed_total",
			Help: "Number of global pubsub messages that failed to send",
		},
	)

	SendEmailFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tknobs_email_send_failed_total",
			Help: "Number of times sending an email to the smtp server failed",
		},
	)
)

func InitMetrics() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		PipelineRunsProcessedTotal,
		PipelineRunsStartedProcessingTotal,
		LogsSavedToGcsTotal,
		LogsSavedToGcsSkippedDisabledTotal,
		LogsSavedToGcsFailedTotal,
		WebexMessagesSentTotal,
		WebexMessagesSkippedDisabledTotal,
		WebexMessagesSkippedSuccessTotal,
		WebexMessagesFailedTotal,
		GithubCommitCommentCreatedTotal,
		GithubCommentSkippedDisabledTotal,
		GithubCommitCommentFailedTotal,
		GithubPrCommentCreatedTotal,
		GithubPrCommentFailedTotal,
		GithubPrCommentSkippedNotPrTotal,
		GithubStatusSkippedDisabledTotal,
		GithubStatusFailedTotal,
		GithubStatusCreatedTotal,
		GithubStatusSkippedFailedExistedTotal,
		GithubDeploymentSkippedDisabledTotal,
		GithubDeploymentFailedTotal,
		GithubDeploymentCreatedTotal,
		GithubDeploymentSkippedEventNotMatchTotal,
		EmailSkippedDisabledTotal,
		EmailFailedTotal,
		EmailCreatedTotal,
		EmailSkippedSuccessTotal,
		PubSubSkippedDisabledTotal,
		PubSubSentTotal,
		PubSubFailedTotal,
		PubSubGlobalSentTotal,
		PubSubGlobalFailedTotal,
		SendEmailFailedTotal,
		SendEmailRequestTimeHistogram,
		EmitEventRequestTimeHistogram,
		GithubRequestTimeHistogram,
		GoogleRequestTimeHistogram,
		KubernetesRequestTimeHistogram,
		WebexRequestTimeHistogram,
		ProcessPipelineTimeHistogram,
	)

}

func GetStatusCode(err error) string {

	if err == nil {
		return "200"
	}
	return "400"
}
