package cmd

import (
	"reflect"
	"testing"
)

func TestFilterAndSortPodNames_FilterAndSort(t *testing.T) {
	in := []string{
		"kube-scheduler",
		"coredns-abc",
		"coredns-def",
		"etcd-node",
	}

	got := filterAndSortPodNames(in, "core")

	want := []string{
		"coredns-abc",
		"coredns-def",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestFilterAndSortPodNames_NoPrefixReturnsAllSorted(t *testing.T) {
	in := []string{
		"kube-scheduler",
		"coredns-def",
		"coredns-abc",
	}

	got := filterAndSortPodNames(in, "")

	want := []string{
		"coredns-abc",
		"coredns-def",
		"kube-scheduler",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestFilterAndSortPodNames_NoMatch(t *testing.T) {
	in := []string{
		"kube-scheduler",
		"coredns-def",
	}

	got := filterAndSortPodNames(in, "abc")

	if len(got) != 0 {
		t.Fatalf("expected empty result, got %#v", got)
	}
}
