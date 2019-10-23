package main

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestIsDeletable(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2019-07-18T11:05:06+09:00")
	duration, _ := time.ParseDuration("10m")
	m := jobManager{now: now, deletedDuration: duration}

	ct, _ := time.Parse(time.RFC3339, "2019-07-18T10:05:06+09:00")
	j := &batchv1.Job{Status: batchv1.JobStatus{CompletionTime: &metav1.Time{Time: ct}}}
	if !m.isDeletable(j) {
		t.Error("this job is deletable")
	}

	ct, _ = time.Parse(time.RFC3339, "2019-07-18T11:00:00+09:00")
	j = &batchv1.Job{Status: batchv1.JobStatus{CompletionTime: &metav1.Time{Time: ct}}}
	if m.isDeletable(j) {
		t.Error("this job is not deletable")
	}
}
