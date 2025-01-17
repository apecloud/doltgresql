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

import (
	"fmt"

	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/doltgresql/server/functions/framework"
	pgtypes "github.com/dolthub/doltgresql/server/types"
)

// initToHex registers the functions to the catalog.
func initToHex() {
	framework.RegisterFunction(to_hex_int32)
	framework.RegisterFunction(to_hex_int64)
}

// to_hex_int32 represents the PostgreSQL function of the same name, taking the same parameters.
var to_hex_int32 = framework.Function1{
	Name:       "to_hex",
	Return:     pgtypes.Text,
	Parameters: [1]pgtypes.DoltgresType{pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]pgtypes.DoltgresType, val1 any) (any, error) {
		return fmt.Sprintf("%x", uint64(val1.(int32))), nil
	},
}

// to_hex_int64 represents the PostgreSQL function of the same name, taking the same parameters.
var to_hex_int64 = framework.Function1{
	Name:       "to_hex",
	Return:     pgtypes.Text,
	Parameters: [1]pgtypes.DoltgresType{pgtypes.Int64},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]pgtypes.DoltgresType, val1 any) (any, error) {
		return fmt.Sprintf("%x", uint64(val1.(int64))), nil
	},
}
