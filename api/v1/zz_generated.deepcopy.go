//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nodered) DeepCopyInto(out *Nodered) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nodered.
func (in *Nodered) DeepCopy() *Nodered {
	if in == nil {
		return nil
	}
	out := new(Nodered)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Nodered) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NoderedList) DeepCopyInto(out *NoderedList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Nodered, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NoderedList.
func (in *NoderedList) DeepCopy() *NoderedList {
	if in == nil {
		return nil
	}
	out := new(NoderedList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NoderedList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NoderedSettings) DeepCopyInto(out *NoderedSettings) {
	*out = *in
	out.EditorTheme = in.EditorTheme
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NoderedSettings.
func (in *NoderedSettings) DeepCopy() *NoderedSettings {
	if in == nil {
		return nil
	}
	out := new(NoderedSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NoderedSettingsEditorTheme) DeepCopyInto(out *NoderedSettingsEditorTheme) {
	*out = *in
	out.Page = in.Page
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NoderedSettingsEditorTheme.
func (in *NoderedSettingsEditorTheme) DeepCopy() *NoderedSettingsEditorTheme {
	if in == nil {
		return nil
	}
	out := new(NoderedSettingsEditorTheme)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NoderedSettingsEditorThemePage) DeepCopyInto(out *NoderedSettingsEditorThemePage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NoderedSettingsEditorThemePage.
func (in *NoderedSettingsEditorThemePage) DeepCopy() *NoderedSettingsEditorThemePage {
	if in == nil {
		return nil
	}
	out := new(NoderedSettingsEditorThemePage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NoderedSpec) DeepCopyInto(out *NoderedSpec) {
	*out = *in
	out.Settings = in.Settings
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NoderedSpec.
func (in *NoderedSpec) DeepCopy() *NoderedSpec {
	if in == nil {
		return nil
	}
	out := new(NoderedSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NoderedStatus) DeepCopyInto(out *NoderedStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NoderedStatus.
func (in *NoderedStatus) DeepCopy() *NoderedStatus {
	if in == nil {
		return nil
	}
	out := new(NoderedStatus)
	in.DeepCopyInto(out)
	return out
}
