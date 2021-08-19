package notebook_instance_lifecycle_config

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// customDeltaOnCreate compares the difference between
// two onCreate fields. It makes use of the ackcompare.Splice.
func customDeltaOnCreate(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) *ackcompare.Delta {
	if ackcompare.HasNilDifference(a.ko.Spec.OnCreate, b.ko.Spec.OnCreate) {
		delta.Add("Spec.OnCreate", a.ko.Spec.OnCreate, b.ko.Spec.OnCreate)
		return delta
	}
	var lista []*string
	var listb []*string

	for _, s := range a.ko.Spec.OnCreate {
		lista = append(lista, s.Content)
	}
	for _, s := range b.ko.Spec.OnCreate {
		listb = append(listb, s.Content)
	}
	if !ackcompare.SliceStringPEqual(lista, listb) {
		delta.Add("Spec.OnCreate", a.ko.Spec.OnCreate, b.ko.Spec.OnCreate)
	}

	return delta
}

// customDeltaOnStart compares the difference between
// two onStart fields. It makes use of the ackcompare.Splice.
func customDeltaOnStart(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) *ackcompare.Delta {
	if ackcompare.HasNilDifference(a.ko.Spec.OnStart, b.ko.Spec.OnStart) {
		delta.Add("Spec.OnStart", a.ko.Spec.OnStart, b.ko.Spec.OnStart)
		return delta
	}

	var lista []*string
	var listb []*string

	for _, s := range a.ko.Spec.OnStart {
		lista = append(lista, s.Content)
	}
	for _, s := range b.ko.Spec.OnStart {
		listb = append(listb, s.Content)
	}
	if !ackcompare.SliceStringPEqual(lista, listb) {
		delta.Add("Spec.OnStart", a.ko.Spec.OnStart, b.ko.Spec.OnStart)
	}

	return delta
}
