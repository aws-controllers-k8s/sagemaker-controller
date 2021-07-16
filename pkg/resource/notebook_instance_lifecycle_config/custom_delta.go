package notebook_instance_lifecycle_config

import (
	"encoding/base64"
	"fmt"

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
	//check length
	if a.ko.Spec.OnCreate != nil && b.ko.Spec.OnCreate != nil {
		if len(a.ko.Spec.OnCreate) != len(b.ko.Spec.OnCreate) {
			delta.Add("Spec.OnCreate", a.ko.Spec.OnCreate, b.ko.Spec.OnCreate)
			return delta
		}
	} else {
		return delta
	}
	//Check variables, a and b have to be of equal length
	onCreateLen := len(a.ko.Spec.OnCreate)
	if a.ko.Spec.OnCreate != nil && b.ko.Spec.OnCreate != nil {
		for i := 0; i < onCreateLen; i++ {
			abb := *a.ko.Spec.OnCreate[i].Content
			bbb := *b.ko.Spec.OnCreate[i].Content
			if err, re := IsBase64(*a.ko.Spec.OnCreate[i].Content); err == false {
				abb = re
			}
			if err, re := IsBase64(*b.ko.Spec.OnCreate[i].Content); err == false {
				bbb = re
			}
			if abb != bbb {
				fmt.Println(abb, bbb)
				delta.Add("Spec.OnCreate", a.ko.Spec.OnCreate, b.ko.Spec.OnCreate)
				return delta
			}
		}
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
	//check length
	if a.ko.Spec.OnStart != nil && b.ko.Spec.OnStart != nil {
		if len(a.ko.Spec.OnStart) != len(b.ko.Spec.OnStart) {
			delta.Add("Spec.OnStart", a.ko.Spec.OnStart, b.ko.Spec.OnStart)
			return delta
		}
	} else {
		return delta
	}
	//Check variables, a and b have to be of equal length
	onStartLen := len(a.ko.Spec.OnStart)
	if a.ko.Spec.OnStart != nil && b.ko.Spec.OnStart != nil {
		for i := 0; i < onStartLen; i++ {
			abb := *a.ko.Spec.OnStart[i].Content
			bbb := *b.ko.Spec.OnStart[i].Content
			if err, re := IsBase64(*a.ko.Spec.OnStart[i].Content); err == false {
				abb = re
			}
			if err, re := IsBase64(*b.ko.Spec.OnStart[i].Content); err == false {
				bbb = re
			}
			if abb != bbb {
				fmt.Println(abb, bbb)
				delta.Add("Spec.OnStart", a.ko.Spec.OnStart, b.ko.Spec.OnStart)
				return delta
			}
		}

	}
	return delta
}

func IsBase64(s string) (bool, string) {
	res, err := base64.StdEncoding.DecodeString(s)
	return err != nil, string(res)
}
