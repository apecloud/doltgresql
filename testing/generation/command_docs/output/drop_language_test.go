// Copyright 2023 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import "testing"

func TestDropLanguage(t *testing.T) {
	tests := []QueryParses{
		Unimplemented("DROP LANGUAGE name"),
		Unimplemented("DROP PROCEDURAL LANGUAGE name"),
		Unimplemented("DROP LANGUAGE IF EXISTS name"),
		Unimplemented("DROP PROCEDURAL LANGUAGE IF EXISTS name"),
		Unimplemented("DROP LANGUAGE name CASCADE"),
		Unimplemented("DROP PROCEDURAL LANGUAGE name CASCADE"),
		Unimplemented("DROP LANGUAGE IF EXISTS name CASCADE"),
		Unimplemented("DROP PROCEDURAL LANGUAGE IF EXISTS name CASCADE"),
		Unimplemented("DROP LANGUAGE name RESTRICT"),
		Unimplemented("DROP PROCEDURAL LANGUAGE name RESTRICT"),
		Unimplemented("DROP LANGUAGE IF EXISTS name RESTRICT"),
		Unimplemented("DROP PROCEDURAL LANGUAGE IF EXISTS name RESTRICT"),
	}
	RunTests(t, tests)
}