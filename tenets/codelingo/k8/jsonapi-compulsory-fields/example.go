package main

import "fmt"

type KindWithoutName struct {
	UID         types.UID                   `json:"uid" protobuf:"bytes,1,opt,name=uid"`
	Kind        metav1.GroupVersionKind     `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	Resource    metav1.GroupVersionResource `json:"resource" protobuf:"bytes,3,opt,name=resource"`
	SubResource string                      `json:"subResource,omitempty" protobuf:"bytes,4,opt,name=subResource"`
	Namespace   string                      `json:"namespace,omitempty" protobuf:"bytes,6,opt,name=namespace"`
	Operation   Operation                   `json:"operation" protobuf:"bytes,7,opt,name=operation"`
	UserInfo    authenticationv1.UserInfo   `json:"userInfo" protobuf:"bytes,8,opt,name=userInfo"`
	Object      runtime.RawExtension        `json:"object,omitempty" protobuf:"bytes,9,opt,name=object"`
	OldObject   runtime.RawExtension        `json:"oldObject,omitempty" protobuf:"bytes,10,opt,name=oldObject"`
	DryRun      *bool                       `json:"dryRun,omitempty" protobuf:"varint,11,opt,name=dryRun"`
}

type KindWithoutNamespace struct {
	UID         types.UID                   `json:"uid" protobuf:"bytes,1,opt,name=uid"`
	Kind        metav1.GroupVersionKind     `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	Resource    metav1.GroupVersionResource `json:"resource" protobuf:"bytes,3,opt,name=resource"`
	SubResource string                      `json:"subResource,omitempty" protobuf:"bytes,4,opt,name=subResource"`
	Name        string                      `json:"name,omitempty" protobuf:"bytes,5,opt,name=name"`
	Operation   Operation                   `json:"operation" protobuf:"bytes,7,opt,name=operation"`
	UserInfo    authenticationv1.UserInfo   `json:"userInfo" protobuf:"bytes,8,opt,name=userInfo"`
	Object      runtime.RawExtension        `json:"object,omitempty" protobuf:"bytes,9,opt,name=object"`
	OldObject   runtime.RawExtension        `json:"oldObject,omitempty" protobuf:"bytes,10,opt,name=oldObject"`
	DryRun      *bool                       `json:"dryRun,omitempty" protobuf:"varint,11,opt,name=dryRun"`
}

type AdmissionRequest struct {
	UID         types.UID                   `json:"uid" protobuf:"bytes,1,opt,name=uid"`
	Kind        metav1.GroupVersionKind     `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	Resource    metav1.GroupVersionResource `json:"resource" protobuf:"bytes,3,opt,name=resource"`
	SubResource string                      `json:"subResource,omitempty" protobuf:"bytes,4,opt,name=subResource"`
	Name        string                      `json:"name,omitempty" protobuf:"bytes,5,opt,name=name"`
	Namespace   string                      `json:"namespace,omitempty" protobuf:"bytes,6,opt,name=namespace"`
	Operation   Operation                   `json:"operation" protobuf:"bytes,7,opt,name=operation"`
	UserInfo    authenticationv1.UserInfo   `json:"userInfo" protobuf:"bytes,8,opt,name=userInfo"`
	Object      runtime.RawExtension        `json:"object,omitempty" protobuf:"bytes,9,opt,name=object"`
	OldObject   runtime.RawExtension        `json:"oldObject,omitempty" protobuf:"bytes,10,opt,name=oldObject"`
	DryRun      *bool                       `json:"dryRun,omitempty" protobuf:"varint,11,opt,name=dryRun"`
}

// So, look for ~name=kind~ which does not have ~name=namespace~, nor ~name=name~.

// Need to look for 'Objects'.
// we want ~name=apiVersion~ and ~name=kind~ this time, but for apis
// that return
// I should deal with jsonapi-obj-compulsory-metadata first to learn how to detect objects
// vim +/"Every object SHOULD have the following metadata in a nested object field" "$HOME/notes2018/ws/codelingo/issues/kubernetes/tenets-found.org"

// look for this
// APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,3,opt,name=apiVersion"`
