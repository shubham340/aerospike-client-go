// Copyright 2013-2017 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import xornd "github.com/aerospike/aerospike-client-go/types/rand"

// Statement encapsulates query statement parameters.
type Statement struct {
	// Namespace determines query Namespace
	Namespace string

	// SetName determines query Set name (Optional)
	SetName string

	// IndexName determines query index name (Optional)
	// If not set, the server will determine the index from the filter's bin name.
	IndexName string

	// BinNames detemines bin names (optional)
	BinNames []string

	// Filters determine query filters (Optional)
	// Currently, only one filter is allowed by the server on a secondary index lookup.
	// If multiple filters are necessary, see QueryFilter example for a workaround.
	// QueryFilter demonstrates how to add additional filters in an user-defined
	// aggregation function.
	Filters []*Filter

	packageName  string
	functionName string
	functionArgs []Value

	// Ordered list of predicate expressions
	predExps  []predExp
	predicate *boolExpression

	// TaskId determines query task id. (Optional)
	TaskId uint64

	// determines if the query should return data
	returnData bool
}

// NewStatement initializes a new Statement instance.
func NewStatement(ns string, set string, binNames ...string) *Statement {
	return &Statement{
		Namespace:  ns,
		SetName:    set,
		BinNames:   binNames,
		returnData: true,
		TaskId:     uint64(xornd.Int64()),
	}
}

// Addfilter adds a filter to the statement.
// Aerospike Server currently only supports using a single filter per statement/query.
func (stmt *Statement) Addfilter(filter *Filter) error {
	stmt.Filters = append(stmt.Filters, filter)

	return nil
}

// Adds a low-level predicate to the statement. Used fr internal testing only.
func (stmt *Statement) AddPredExp(predexp predExp) error {
	stmt.predExps = append(stmt.predExps, predexp)
	return nil
}

// SetPredicate accepts a predicate to evaluate on the server-side for non-indexed bins.
func (stmt *Statement) SetPredicate(exp *boolExpression) error {
	stmt.predicate = exp
	return nil
}

// SetAggregateFunction sets aggregation function parameters.
// This function will be called on both the server
// and client for each selected item.
func (stmt *Statement) SetAggregateFunction(packageName string, functionName string, functionArgs []Value, returnData bool) {
	stmt.packageName = packageName
	stmt.functionName = functionName
	stmt.functionArgs = functionArgs
	stmt.returnData = returnData
}

// IsScan determines is the Statement is a full namespace/set scan or a selective Query.
func (stmt *Statement) IsScan() bool {
	return len(stmt.Filters) == 0
}

// Always set the taskId client-side to a non-zero random value
func (stmt *Statement) setTaskId() {
	for stmt.TaskId == 0 {
		stmt.TaskId = uint64(xornd.Int64())
	}
}
