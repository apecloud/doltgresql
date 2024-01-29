// Copyright 2024 Dolthub, Inc.
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

package functions

import "strings"

// lower represents the PostgreSQL function of the same name.
var lower = Function{
	Name:      "lower",
	Overloads: []interface{}{lower_string},
}

// lower_string is one of the overloads of lower.
func lower_string(text StringType) (StringType, error) {
	if text.IsNull {
		return StringType{IsNull: true}, nil
	}
	//TODO: this doesn't respect collations
	return StringType{Value: strings.ToLower(text.Value)}, nil
}