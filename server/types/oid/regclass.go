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

package oid

import (
	"fmt"
	"strconv"

	"github.com/dolthub/dolt/go/libraries/doltcore/sqle/resolve"
	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/doltgresql/server/settings"
)

// regclass_IoInput is the implementation for IoInput that avoids circular dependencies by being declared in a separate
// package.
func regclass_IoInput(ctx *sql.Context, input string) (uint32, error) {
	// If the string just represents a number, then we return it.
	if parsedOid, err := strconv.ParseUint(input, 10, 32); err == nil {
		return uint32(parsedOid), nil
	}
	sections, err := ioInputSections(input)
	if err != nil {
		return 0, err
	}
	if err = regclass_IoInputValidation(ctx, input, sections); err != nil {
		return 0, err
	}

	var database string
	var searchSchemas []string
	var relationName string
	switch len(sections) {
	case 1:
		database = ctx.GetCurrentDatabase()
		searchSchemas, err = resolve.SearchPath(ctx)
		if err != nil {
			return 0, err
		}
		relationName = sections[0]
	case 3:
		database = ctx.GetCurrentDatabase()
		searchSchemas = []string{sections[0]}
		relationName = sections[2]
	case 5:
		database = sections[0]
		searchSchemas = []string{sections[2]}
		relationName = sections[4]
	default:
		return 0, fmt.Errorf("regclass failed validation")
	}

	// Iterate over all of the items to find which relation matches.
	// Postgres does not need to worry about name conflicts since everything is created in the same naming space, but
	// GMS and Dolt use different naming spaces, so for now we just ignore potential name conflicts and return the first
	// match found.
	resultOid := uint32(0)
	err = IterateDatabase(ctx, database, Callbacks{
		Index: func(ctx *sql.Context, schema ItemSchema, table ItemTable, index ItemIndex) (cont bool, err error) {
			idxName := index.Item.ID()
			if idxName == "PRIMARY" {
				idxName = fmt.Sprintf("%s_pkey", index.Item.Table())
			}
			if relationName == idxName {
				resultOid = index.OID
				return false, nil
			}
			return true, nil
		},
		Sequence: func(ctx *sql.Context, schema ItemSchema, sequence ItemSequence) (cont bool, err error) {
			if sequence.Item.Name == relationName {
				resultOid = sequence.OID
				return false, nil
			}
			return true, nil
		},
		Table: func(ctx *sql.Context, schema ItemSchema, table ItemTable) (cont bool, err error) {
			if table.Item.Name() == relationName {
				resultOid = table.OID
				return false, nil
			}
			return true, nil
		},
		View: func(ctx *sql.Context, schema ItemSchema, view ItemView) (cont bool, err error) {
			if view.Item.Name == relationName {
				resultOid = view.OID
				return false, nil
			}
			return true, nil
		},
		SearchSchemas: searchSchemas,
	})
	if err != nil || resultOid != 0 {
		return resultOid, err
	}
	return 0, fmt.Errorf(`relation "%s" does not exist`, input)
}

// regclass_IoOutput is the implementation for IoOutput that avoids circular dependencies by being declared in a separate
// package.
func regclass_IoOutput(ctx *sql.Context, oid uint32) (string, error) {
	// Find all the schemas on the search path. If a schema is on the search path, then it is not included in the
	// output of relation name. If the relation's schema is not on the search path, then it is explicitly included.
	schemasMap, err := settings.GetCurrentSchemasAsMap(ctx)
	if err != nil {
		return "", err
	}

	// The pg_catalog schema is always implicitly part of the search path
	// https://www.postgresql.org/docs/current/ddl-schemas.html#DDL-SCHEMAS-CATALOG
	schemasMap["pg_catalog"] = struct{}{}

	output := strconv.FormatUint(uint64(oid), 10)
	err = RunCallback(ctx, oid, Callbacks{
		Index: func(ctx *sql.Context, schema ItemSchema, table ItemTable, index ItemIndex) (cont bool, err error) {
			output = index.Item.ID()
			if output == "PRIMARY" {
				schemaName := schema.Item.SchemaName()
				if _, ok := schemasMap[schemaName]; ok {
					output = fmt.Sprintf("%s_pkey", index.Item.Table())
				} else {
					output = fmt.Sprintf("%s.%s_pkey", schemaName, index.Item.Table())
				}
			}
			return false, nil
		},
		Sequence: func(ctx *sql.Context, schema ItemSchema, sequence ItemSequence) (cont bool, err error) {
			schemaName := schema.Item.SchemaName()
			if _, ok := schemasMap[schemaName]; ok {
				output = sequence.Item.Name
			} else {
				output = fmt.Sprintf("%s.%s", schemaName, sequence.Item.Name)
			}
			return false, nil
		},
		Table: func(ctx *sql.Context, schema ItemSchema, table ItemTable) (cont bool, err error) {
			schemaName := schema.Item.SchemaName()
			if _, ok := schemasMap[schemaName]; ok {
				output = table.Item.Name()
			} else {
				output = fmt.Sprintf("%s.%s", schemaName, table.Item.Name())
			}
			return false, nil
		},
		View: func(ctx *sql.Context, schema ItemSchema, view ItemView) (cont bool, err error) {
			schemaName := schema.Item.SchemaName()
			if _, ok := schemasMap[schemaName]; ok {
				output = view.Item.Name
			} else {
				output = fmt.Sprintf("%s.%s", schemaName, view.Item.Name)
			}
			return false, nil
		},
	})
	return output, err
}

// regclass_IoInputValidation handles the validation for the parsed sections in regclass_IoInput.
func regclass_IoInputValidation(ctx *sql.Context, input string, sections []string) error {
	switch len(sections) {
	case 1:
		return nil
	case 3:
		if sections[1] != "." {
			return fmt.Errorf("invalid name syntax")
		}
		return nil
	case 5:
		if sections[1] != "." || sections[3] != "." {
			return fmt.Errorf("invalid name syntax")
		}
		return nil
	case 7:
		if sections[1] != "." || sections[3] != "." || sections[5] != "." {
			return fmt.Errorf("invalid name syntax")
		}
		return fmt.Errorf("improper qualified name (too many dotted names): %s", input)
	default:
		return fmt.Errorf("invalid name syntax")
	}
}
