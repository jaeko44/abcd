// +build !ignore_autogenerated_openshift

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package api

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	reflect "reflect"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

// RegisterDeepCopies adds deep-copy functions to the given scheme. Public
// to allow building arbitrary schemes.
func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_Group, InType: reflect.TypeOf(&Group{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_GroupList, InType: reflect.TypeOf(&GroupList{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_Identity, InType: reflect.TypeOf(&Identity{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_IdentityList, InType: reflect.TypeOf(&IdentityList{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_User, InType: reflect.TypeOf(&User{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_UserIdentityMapping, InType: reflect.TypeOf(&UserIdentityMapping{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_UserList, InType: reflect.TypeOf(&UserList{})},
	)
}

func DeepCopy_api_Group(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*Group)
		out := out.(*Group)
		*out = *in
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*v1.ObjectMeta)
		}
		if in.Users != nil {
			in, out := &in.Users, &out.Users
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
		return nil
	}
}

func DeepCopy_api_GroupList(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*GroupList)
		out := out.(*GroupList)
		*out = *in
		if in.Items != nil {
			in, out := &in.Items, &out.Items
			*out = make([]Group, len(*in))
			for i := range *in {
				if err := DeepCopy_api_Group(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func DeepCopy_api_Identity(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*Identity)
		out := out.(*Identity)
		*out = *in
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*v1.ObjectMeta)
		}
		if in.Extra != nil {
			in, out := &in.Extra, &out.Extra
			*out = make(map[string]string)
			for key, val := range *in {
				(*out)[key] = val
			}
		}
		return nil
	}
}

func DeepCopy_api_IdentityList(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*IdentityList)
		out := out.(*IdentityList)
		*out = *in
		if in.Items != nil {
			in, out := &in.Items, &out.Items
			*out = make([]Identity, len(*in))
			for i := range *in {
				if err := DeepCopy_api_Identity(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func DeepCopy_api_User(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*User)
		out := out.(*User)
		*out = *in
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*v1.ObjectMeta)
		}
		if in.Identities != nil {
			in, out := &in.Identities, &out.Identities
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
		if in.Groups != nil {
			in, out := &in.Groups, &out.Groups
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
		return nil
	}
}

func DeepCopy_api_UserIdentityMapping(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*UserIdentityMapping)
		out := out.(*UserIdentityMapping)
		*out = *in
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*v1.ObjectMeta)
		}
		return nil
	}
}

func DeepCopy_api_UserList(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*UserList)
		out := out.(*UserList)
		*out = *in
		if in.Items != nil {
			in, out := &in.Items, &out.Items
			*out = make([]User, len(*in))
			for i := range *in {
				if err := DeepCopy_api_User(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
