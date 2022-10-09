//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	"unsafe"

	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	configv1alpha1 "k8s.io/component-base/config/v1alpha1"
	"k8s.io/controller-manager/config"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*ControllerLeaderConfiguration)(nil), (*config.ControllerLeaderConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ControllerLeaderConfiguration_To_config_ControllerLeaderConfiguration(a.(*ControllerLeaderConfiguration), b.(*config.ControllerLeaderConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ControllerLeaderConfiguration)(nil), (*ControllerLeaderConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ControllerLeaderConfiguration_To_v1alpha1_ControllerLeaderConfiguration(a.(*config.ControllerLeaderConfiguration), b.(*ControllerLeaderConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LeaderMigrationConfiguration)(nil), (*config.LeaderMigrationConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LeaderMigrationConfiguration_To_config_LeaderMigrationConfiguration(a.(*LeaderMigrationConfiguration), b.(*config.LeaderMigrationConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LeaderMigrationConfiguration)(nil), (*LeaderMigrationConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LeaderMigrationConfiguration_To_v1alpha1_LeaderMigrationConfiguration(a.(*config.LeaderMigrationConfiguration), b.(*LeaderMigrationConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.GenericControllerManagerConfiguration)(nil), (*GenericControllerManagerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GenericControllerManagerConfiguration_To_v1alpha1_GenericControllerManagerConfiguration(a.(*config.GenericControllerManagerConfiguration), b.(*GenericControllerManagerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*GenericControllerManagerConfiguration)(nil), (*config.GenericControllerManagerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration(a.(*GenericControllerManagerConfiguration), b.(*config.GenericControllerManagerConfiguration), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_ControllerLeaderConfiguration_To_config_ControllerLeaderConfiguration(in *ControllerLeaderConfiguration, out *config.ControllerLeaderConfiguration, s conversion.Scope) error {
	out.Name = in.Name
	out.Component = in.Component
	return nil
}

// Convert_v1alpha1_ControllerLeaderConfiguration_To_config_ControllerLeaderConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_ControllerLeaderConfiguration_To_config_ControllerLeaderConfiguration(in *ControllerLeaderConfiguration, out *config.ControllerLeaderConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_ControllerLeaderConfiguration_To_config_ControllerLeaderConfiguration(in, out, s)
}

func autoConvert_config_ControllerLeaderConfiguration_To_v1alpha1_ControllerLeaderConfiguration(in *config.ControllerLeaderConfiguration, out *ControllerLeaderConfiguration, s conversion.Scope) error {
	out.Name = in.Name
	out.Component = in.Component
	return nil
}

// Convert_config_ControllerLeaderConfiguration_To_v1alpha1_ControllerLeaderConfiguration is an autogenerated conversion function.
func Convert_config_ControllerLeaderConfiguration_To_v1alpha1_ControllerLeaderConfiguration(in *config.ControllerLeaderConfiguration, out *ControllerLeaderConfiguration, s conversion.Scope) error {
	return autoConvert_config_ControllerLeaderConfiguration_To_v1alpha1_ControllerLeaderConfiguration(in, out, s)
}

func autoConvert_v1alpha1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration(in *GenericControllerManagerConfiguration, out *config.GenericControllerManagerConfiguration, s conversion.Scope) error {
	out.Port = in.Port
	out.Address = in.Address
	out.MinResyncPeriod = in.MinResyncPeriod
	if err := configv1alpha1.Convert_v1alpha1_ClientConnectionConfiguration_To_config_ClientConnectionConfiguration(&in.ClientConnection, &out.ClientConnection, s); err != nil {
		return err
	}
	out.ControllerStartInterval = in.ControllerStartInterval
	if err := configv1alpha1.Convert_v1alpha1_LeaderElectionConfiguration_To_config_LeaderElectionConfiguration(&in.LeaderElection, &out.LeaderElection, s); err != nil {
		return err
	}
	out.Controllers = *(*[]string)(unsafe.Pointer(&in.Controllers))
	if err := configv1alpha1.Convert_v1alpha1_DebuggingConfiguration_To_config_DebuggingConfiguration(&in.Debugging, &out.Debugging, s); err != nil {
		return err
	}
	out.LeaderMigrationEnabled = in.LeaderMigrationEnabled
	if err := Convert_v1alpha1_LeaderMigrationConfiguration_To_config_LeaderMigrationConfiguration(&in.LeaderMigration, &out.LeaderMigration, s); err != nil {
		return err
	}
	return nil
}

func autoConvert_config_GenericControllerManagerConfiguration_To_v1alpha1_GenericControllerManagerConfiguration(in *config.GenericControllerManagerConfiguration, out *GenericControllerManagerConfiguration, s conversion.Scope) error {
	out.Port = in.Port
	out.Address = in.Address
	out.MinResyncPeriod = in.MinResyncPeriod
	if err := configv1alpha1.Convert_config_ClientConnectionConfiguration_To_v1alpha1_ClientConnectionConfiguration(&in.ClientConnection, &out.ClientConnection, s); err != nil {
		return err
	}
	out.ControllerStartInterval = in.ControllerStartInterval
	if err := configv1alpha1.Convert_config_LeaderElectionConfiguration_To_v1alpha1_LeaderElectionConfiguration(&in.LeaderElection, &out.LeaderElection, s); err != nil {
		return err
	}
	out.Controllers = *(*[]string)(unsafe.Pointer(&in.Controllers))
	if err := configv1alpha1.Convert_config_DebuggingConfiguration_To_v1alpha1_DebuggingConfiguration(&in.Debugging, &out.Debugging, s); err != nil {
		return err
	}
	out.LeaderMigrationEnabled = in.LeaderMigrationEnabled
	if err := Convert_config_LeaderMigrationConfiguration_To_v1alpha1_LeaderMigrationConfiguration(&in.LeaderMigration, &out.LeaderMigration, s); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_LeaderMigrationConfiguration_To_config_LeaderMigrationConfiguration(in *LeaderMigrationConfiguration, out *config.LeaderMigrationConfiguration, s conversion.Scope) error {
	out.LeaderName = in.LeaderName
	out.ResourceLock = in.ResourceLock
	out.ControllerLeaders = *(*[]config.ControllerLeaderConfiguration)(unsafe.Pointer(&in.ControllerLeaders))
	return nil
}

// Convert_v1alpha1_LeaderMigrationConfiguration_To_config_LeaderMigrationConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_LeaderMigrationConfiguration_To_config_LeaderMigrationConfiguration(in *LeaderMigrationConfiguration, out *config.LeaderMigrationConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_LeaderMigrationConfiguration_To_config_LeaderMigrationConfiguration(in, out, s)
}

func autoConvert_config_LeaderMigrationConfiguration_To_v1alpha1_LeaderMigrationConfiguration(in *config.LeaderMigrationConfiguration, out *LeaderMigrationConfiguration, s conversion.Scope) error {
	out.LeaderName = in.LeaderName
	out.ResourceLock = in.ResourceLock
	out.ControllerLeaders = *(*[]ControllerLeaderConfiguration)(unsafe.Pointer(&in.ControllerLeaders))
	return nil
}

// Convert_config_LeaderMigrationConfiguration_To_v1alpha1_LeaderMigrationConfiguration is an autogenerated conversion function.
func Convert_config_LeaderMigrationConfiguration_To_v1alpha1_LeaderMigrationConfiguration(in *config.LeaderMigrationConfiguration, out *LeaderMigrationConfiguration, s conversion.Scope) error {
	return autoConvert_config_LeaderMigrationConfiguration_To_v1alpha1_LeaderMigrationConfiguration(in, out, s)
}
