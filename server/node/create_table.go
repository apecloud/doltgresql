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

package node

import (
	"fmt"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/plan"
	"github.com/dolthub/go-mysql-server/sql/rowexec"

	"github.com/dolthub/doltgresql/core"
	"github.com/dolthub/doltgresql/server/auth"
)

// CreateTable is a node that implements functionality specifically relevant to Doltgres' table creation needs.
type CreateTable struct {
	gmsCreateTable *plan.CreateTable
	sequences      []*CreateSequence
}

var _ sql.ExecSourceRel = (*CreateTable)(nil)

// NewCreateTable returns a new *CreateTable.
func NewCreateTable(createTable *plan.CreateTable, sequences []*CreateSequence) *CreateTable {
	return &CreateTable{
		gmsCreateTable: createTable,
		sequences:      sequences,
	}
}

// Children implements the interface sql.ExecSourceRel.
func (c *CreateTable) Children() []sql.Node {
	return c.gmsCreateTable.Children()
}

// IsReadOnly implements the interface sql.ExecSourceRel.
func (c *CreateTable) IsReadOnly() bool {
	return false
}

// Resolved implements the interface sql.ExecSourceRel.
func (c *CreateTable) Resolved() bool {
	return c.gmsCreateTable != nil && c.gmsCreateTable.Resolved()
}

// RowIter implements the interface sql.ExecSourceRel.
func (c *CreateTable) RowIter(ctx *sql.Context, r sql.Row) (sql.RowIter, error) {
	var userRole auth.Role
	auth.LockRead(func() {
		userRole = auth.GetRole(ctx.Client().User)
	})
	if !userRole.IsValid() {
		return nil, fmt.Errorf(`role "%s" does not exist`, ctx.Client().User)
	}

	createTableIter, err := rowexec.DefaultBuilder.Build(ctx, c.gmsCreateTable, r)
	if err != nil {
		return nil, err
	}

	schemaName, err := core.GetSchemaName(ctx, c.gmsCreateTable.Db, "")
	if err != nil {
		return nil, err
	}
	for _, sequence := range c.sequences {
		sequence.schema = schemaName
		_, err = sequence.RowIter(ctx, r)
		if err != nil {
			_ = createTableIter.Close(ctx)
			return nil, err
		}
	}

	auth.LockWrite(func() {
		auth.AddOwner(auth.OwnershipKey{
			PrivilegeObject: auth.PrivilegeObject_TABLE,
			Schema:          schemaName,
			Name:            c.gmsCreateTable.Name(),
		}, userRole.ID())
		err = auth.PersistChanges()
	})
	if err != nil {
		return nil, err
	}
	return createTableIter, err
}

// Schema implements the interface sql.ExecSourceRel.
func (c *CreateTable) Schema() sql.Schema {
	return c.gmsCreateTable.Schema()
}

// String implements the interface sql.ExecSourceRel.
func (c *CreateTable) String() string {
	return c.gmsCreateTable.String()
}

// WithChildren implements the interface sql.ExecSourceRel.
func (c *CreateTable) WithChildren(children ...sql.Node) (sql.Node, error) {
	gmsCreateTable, err := c.gmsCreateTable.WithChildren(children...)
	if err != nil {
		return nil, err
	}
	return &CreateTable{
		gmsCreateTable: gmsCreateTable.(*plan.CreateTable),
		sequences:      c.sequences,
	}, nil
}
