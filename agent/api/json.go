// Copyright 2014-2015 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package api

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/amazon-ecs-agent/agent/utils"
)

func (ts *TaskStatus) UnmarshalJSON(b []byte) error {
	if strings.ToLower(string(b)) == "null" {
		*ts = TaskStatusNone
		return nil
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		*ts = TaskStatusNone
		return errors.New("TaskStatus must be a string or null")
	}
	strStatus := string(b[1 : len(b)-1])
	// 'UNKNOWN' and 'DEAD' for Compatibility with v1.0.0 state files
	if strStatus == "UNKNOWN" {
		*ts = TaskStatusNone
		return nil
	}
	if strStatus == "DEAD" {
		*ts = TaskStopped
		return nil
	}

	stat, ok := taskStatusMap[strStatus]
	if !ok {
		*ts = TaskStatusNone
		return errors.New("Unrecognized TaskStatus")
	}
	*ts = stat
	return nil
}

func (ts *TaskStatus) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return nil, nil
	}
	return []byte(`"` + ts.String() + `"`), nil
}

func (cs *ContainerStatus) UnmarshalJSON(b []byte) error {
	if strings.ToLower(string(b)) == "null" {
		*cs = ContainerStatusNone
		return nil
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		*cs = ContainerStatusNone
		return errors.New("ContainerStatus must be a string or null; Got " + string(b))
	}
	strStatus := string(b[1 : len(b)-1])
	// 'UNKNOWN' and 'DEAD' for Compatibility with v1.0.0 state files
	if strStatus == "UNKNOWN" {
		*cs = ContainerStatusNone
		return nil
	}
	if strStatus == "DEAD" {
		*cs = ContainerStopped
		return nil
	}

	stat, ok := containerStatusMap[strStatus]
	if !ok {
		*cs = ContainerStatusNone
		return errors.New("Unrecognized ContainerStatus")
	}
	*cs = stat
	return nil
}

func (cs *ContainerStatus) MarshalJSON() ([]byte, error) {
	if cs == nil {
		return nil, nil
	}
	return []byte(`"` + cs.String() + `"`), nil
}

// A type alias that doesn't have a custom unmarshaller so we can unmarshal into
// something without recursing
type ContainerOverridesCopy ContainerOverrides

// This custom unmarshaller is needed because the json sent to us as a string
// rather than a fully typed object. We support both formats in the hopes that
// one day everything will be fully typed
// Note: the `json:",string"` tag DOES NOT apply here; it DOES NOT work with
// struct types, only ints/floats/etc. We're basically doing that though
// We also intentionally fail if there are any keys we were unable to unmarshal
// into our struct
func (overrides *ContainerOverrides) UnmarshalJSON(b []byte) error {
	regular := ContainerOverridesCopy{}

	// Try to do it the strongly typed way first
	err := json.Unmarshal(b, &regular)
	if err == nil {
		err = utils.CompleteJsonUnmarshal(b, regular)
		if err == nil {
			*overrides = ContainerOverrides(regular)
			return nil
		}
		err = utils.NewMultiError(errors.New("Error unmarshalling ContainerOverrides"), err)
	}

	// Now the stringly typed way
	var str string
	err2 := json.Unmarshal(b, &str)
	if err2 != nil {
		return utils.NewMultiError(errors.New("Could not unmarshal ContainerOverrides into either an object or string respectively"), err, err2)
	}

	// We have a string, let's try to unmarshal that into a typed object
	err3 := json.Unmarshal([]byte(str), &regular)
	if err3 == nil {
		err3 = utils.CompleteJsonUnmarshal([]byte(str), regular)
		if err3 == nil {
			*overrides = ContainerOverrides(regular)
			return nil
		} else {
			err3 = utils.NewMultiError(errors.New("Error unmarshalling ContainerOverrides"), err3)
		}
	}

	return utils.NewMultiError(errors.New("Could not unmarshal ContainerOverrides in any supported way"), err, err2, err3)
}

// UnmarshalJSON for TaskVolume determines the name and volume type, and
// unmarshals it into the appropriate HostVolume fulfilling interfaces
func (tv *TaskVolume) UnmarshalJSON(b []byte) error {
	// Format: {name: volumeName, host: emptyVolumeOrHostVolume}
	intermediate := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &intermediate); err != nil {
		return err
	}
	name, ok := intermediate["name"]
	if !ok {
		return errors.New("invalid Volume; must include a name")
	}
	if err := json.Unmarshal(name, &tv.Name); err != nil {
		return err
	}

	if rawhostdata, ok := intermediate["host"]; ok {
		// Default to trying to unmarshal it as a FSHostVolume
		var hostvolume FSHostVolume
		err := json.Unmarshal(rawhostdata, &hostvolume)
		if err != nil {
			return err
		}
		if hostvolume.FSSourcePath == "" {
			// If the FSSourcePath is empty, that must mean it was not an
			// FSHostVolume (empty path is invalid for that type). The only other
			// type is an empty volume, so unmarshal it as such.
			emptyVolume := &EmptyHostVolume{}
			json.Unmarshal(rawhostdata, emptyVolume)
			tv.Volume = emptyVolume
		} else {
			tv.Volume = &hostvolume
		}
		return nil
	}

	return errors.New("unrecognized volume type; try updating me")
}

func (tv *TaskVolume) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})

	result["name"] = tv.Name

	switch v := tv.Volume.(type) {
	case *FSHostVolume:
		result["host"] = v
	case *EmptyHostVolume:
		result["host"] = v
	default:
		log.Crit("Unknown task volume type in marshal")
	}
	return json.Marshal(result)
}

// UnmarshalJSON for TransportProtocol determines whether to use TCP or UDP,
// setting TCP as the zero-value but treating other unrecognized values as
// errors
func (tp *TransportProtocol) UnmarshalJSON(b []byte) error {
	if strings.ToLower(string(b)) == "null" {
		*tp = TransportProtocolTCP
		log.Warn("Unmarshalled nil TransportProtocol as TCP")
		return nil
	}
	switch string(b) {
	case `"tcp"`:
		*tp = TransportProtocolTCP
	case `"udp"`:
		*tp = TransportProtocolUDP
	default:
		*tp = TransportProtocolTCP
		return errors.New("TransportProtocol must be \"tcp\" or \"udp\"; Got " + string(b))
	}
	return nil
}

func (tp *TransportProtocol) MarshalJSON() ([]byte, error) {
	if tp == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + tp.String() + `"`), nil
}
