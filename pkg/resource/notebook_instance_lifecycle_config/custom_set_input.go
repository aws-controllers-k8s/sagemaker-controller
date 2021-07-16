package notebook_instance_lifecycle_config

import (
	"encoding/base64"

	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

func (rm *resourceManager) fixNotebookLFinput(input *svcsdk.CreateNotebookInstanceLifecycleConfigInput) *svcsdk.CreateNotebookInstanceLifecycleConfigInput {
	if input.OnCreate != nil {
		f1 := []*svcsdk.NotebookInstanceLifecycleHook{}
		for _, f1iter := range input.OnCreate {
			f1elem := &svcsdk.NotebookInstanceLifecycleHook{}
			if f1iter.Content != nil {
				b64 := base64.StdEncoding.EncodeToString([]byte(*f1iter.Content))
				f1elem.SetContent(b64)
			}
			f1 = append(f1, f1elem)
		}
		input.SetOnCreate(f1)
	}
	if input.OnStart != nil {
		if input.OnStart != nil {
			f2 := []*svcsdk.NotebookInstanceLifecycleHook{}
			for _, f2iter := range input.OnStart {
				f2elem := &svcsdk.NotebookInstanceLifecycleHook{}
				if f2iter.Content != nil {
					b64 := base64.StdEncoding.EncodeToString([]byte(*f2iter.Content))
					f2elem.SetContent(b64)
				}
				f2 = append(f2, f2elem)
			}
			input.SetOnStart(f2)
		}
	}
	return input
}

func (rm *resourceManager) fixNotebookLFinput_update(input *svcsdk.UpdateNotebookInstanceLifecycleConfigInput) *svcsdk.UpdateNotebookInstanceLifecycleConfigInput {
	if input.OnCreate != nil {
		f1 := []*svcsdk.NotebookInstanceLifecycleHook{}
		for _, f1iter := range input.OnCreate {
			f1elem := &svcsdk.NotebookInstanceLifecycleHook{}
			if f1iter.Content != nil {
				b64 := base64.StdEncoding.EncodeToString([]byte(*f1iter.Content))
				f1elem.SetContent(b64)
			}
			f1 = append(f1, f1elem)
		}
		input.SetOnCreate(f1)
	}
	if input.OnStart != nil {
		if input.OnStart != nil {
			f2 := []*svcsdk.NotebookInstanceLifecycleHook{}
			for _, f2iter := range input.OnStart {
				f2elem := &svcsdk.NotebookInstanceLifecycleHook{}
				if f2iter.Content != nil {
					b64 := base64.StdEncoding.EncodeToString([]byte(*f2iter.Content))
					f2elem.SetContent(b64)
				}
				f2 = append(f2, f2elem)
			}
			input.SetOnStart(f2)
		}
	}
	return input
}
