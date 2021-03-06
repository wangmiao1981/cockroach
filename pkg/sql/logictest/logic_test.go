// Copyright 2017 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package logictest

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

// TestLogic runs logic tests that were written by hand to test various
// CockroachDB features. The tests use a similar methodology to the SQLLite
// Sqllogictests.
//
// See the comments in logic.go for more details.
func TestLogic(t *testing.T) {
	defer leaktest.AfterTest(t)()
	RunLogicTest(t, "testdata/logic_test/[^.]*")
}

// TestSqlLiteLogic runs all logic tests from CockroachDB's fork of
// Sqllogictest:
//
//   https://www.sqlite.org/sqllogictest/doc/trunk/about.wiki
//
// This fork contains many generated tests created by the SqlLite project that
// ensure the tested SQL database returns correct statement and query output.
// The logic tests are reasonably independent of the specific dialect of each
// database so that they can be easily retargeted. In fact, the expected output
// for each test can be generated by one database and then used to verify the
// output of another database.
//
// By default, this test is skipped, unless the `bigtest` flag is specified.
// The reason for this is that these tests are contained in another repo that
// must be present on the machine, and because they take a long time to run.
//
// See the comments in logic.go for more details.
func TestSqlLiteLogic(t *testing.T) {
	defer leaktest.AfterTest(t)()

	if !*bigtest {
		t.Skip("-bigtest flag must be specified to run this test")
	}

	logicTestPath := build.Default.GOPATH + "/src/github.com/cockroachdb/sqllogictest"
	if _, err := os.Stat(logicTestPath); os.IsNotExist(err) {
		fullPath, err := filepath.Abs(logicTestPath)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("unable to find sqllogictest repo: %s\n"+
			"git clone https://github.com/cockroachdb/sqllogictest %s",
			logicTestPath, fullPath)
		return
	}
	globs := []string{
		logicTestPath + "/test/index/between/*/*.test",
		logicTestPath + "/test/index/commute/*/*.test",
		logicTestPath + "/test/index/delete/*/*.test",
		logicTestPath + "/test/index/in/*/*.test",
		logicTestPath + "/test/index/orderby/*/*.test",
		logicTestPath + "/test/index/orderby_nosort/*/*.test",
		logicTestPath + "/test/index/view/*/*.test",

		// TODO(pmattis): Incompatibilities in numeric types.
		// For instance, we type SUM(int) as a decimal since all of our ints are
		// int64.
		// logicTestPath + "/test/random/expr/*.test",

		// TODO(pmattis): We don't support correlated subqueries.
		// logicTestPath + "/test/select*.test",

		// TODO(pmattis): We don't support unary + on strings.
		// logicTestPath + "/test/index/random/*/*.test",
		// [uses joins] logicTestPath + "/test/random/aggregates/*.test",
		// [uses joins] logicTestPath + "/test/random/groupby/*.test",
		// [uses joins] logicTestPath + "/test/random/select/*.test",
	}

	RunLogicTest(t, globs...)
}
