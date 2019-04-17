/*
Copyright 2017 Gravitational, Inc.

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

package utils

// CopyByteSlice returns a copy of the byte slice.
func CopyByteSlice(in []byte) []byte {
	if in == nil {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

// CopyByteSlices returns a copy of the byte slices.
func CopyByteSlices(in [][]byte) [][]byte {
	if in == nil {
		return nil
	}
	out := make([][]byte, len(in))
	for i := range in {
		out[i] = CopyByteSlice(in[i])
	}
	return out
}

// JoinStringSlices joins two string slices and returns a resulting slice
func JoinStringSlices(a []string, b []string) []string {
	if len(a)+len(b) == 0 {
		return nil
	}
	out := make([]string, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	return out
}

// CopyStrings makes a deep copy of the passed in string slice and returns
// the copy.
func CopyStrings(in []string) []string {
	if in == nil {
		return nil
	}

	out := make([]string, len(in))
	copy(out, in)

	return out
}

// CopyStringMapSlices makes a deep copy of the passed in map[string][]string
// and returns the copy.
func CopyStringMapSlices(a map[string][]string) map[string][]string {
	if a == nil {
		return nil
	}

	out := make(map[string][]string)
	for key, values := range a {
		vout := make([]string, len(values))
		copy(vout, values)
		out[key] = vout
	}

	return out
}

// CopyStringMap makes a deep copy of a map[string]string and returns the copy.
func CopyStringMap(a map[string]string) map[string]string {
	if a == nil {
		return nil
	}

	out := make(map[string]string)
	for key, value := range a {
		out[key] = value
	}

	return out
}

// CopyStringMapInterface makes a deep copy of the passed in map[string]interface{}
// and returns the copy.
func CopyStringMapInterface(a map[string]interface{}) map[string]interface{} {
	if a == nil {
		return nil
	}

	out := make(map[string]interface{})
	for key, value := range a {
		out[key] = value
	}

	return out
}
