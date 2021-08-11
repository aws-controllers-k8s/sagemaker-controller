package notebook_instance_lifecycle_config

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func modifyDeltaCreate(
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

	onCreateLenA := len(a.ko.Spec.OnCreate)
	onCreateLenB := len(b.ko.Spec.OnCreate)
	for i := 0; i < onCreateLenA; i++ {
		val := *a.ko.Spec.OnCreate[i].Content
		lista = append(lista, &val)
	}
	for i := 0; i < onCreateLenB; i++ {
		val := *b.ko.Spec.OnCreate[i].Content
		listb = append(listb, &val)
	}
	if !ackcompare.SliceStringPEqual(lista, listb) {
		delta.Add("Spec.OnCreate", a.ko.Spec.OnCreate, b.ko.Spec.OnCreate)
	}

	return delta
}

func modifyDeltaStart(
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

	onStartLenA := len(a.ko.Spec.OnStart)
	onStartLenB := len(b.ko.Spec.OnStart)

	for i := 0; i < onStartLenA; i++ {
		val := *a.ko.Spec.OnStart[i].Content
		lista = append(lista, &val)
	}
	for i := 0; i < onStartLenB; i++ {
		val := *b.ko.Spec.OnStart[i].Content
		listb = append(listb, &val)
	}
	if !ackcompare.SliceStringPEqual(lista, listb) {
		delta.Add("Spec.OnStart", a.ko.Spec.OnStart, b.ko.Spec.OnStart)
	}

	return delta
}
