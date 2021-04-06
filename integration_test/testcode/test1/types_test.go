// Copyright (c) 2021 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package test1_test

import (
	"testing"
)

func TestRefreshableTypes(t *testing.T) {
	//s1 := test1.SuperStruct{
	//	String: "string",
	//	NestedStruct: test1.NestedStruct{
	//		FieldA: "fieldA",
	//	},
	//}
	//r1 := refreshable.NewDefaultRefreshable(s1)
	//ssR1 := test1.NewRefreshingSuperStruct(r1)
	//rString := ssR1.String()
	//assert.Equal(t, "string", rString.CurrentString())
	//rFieldA := ssR1.NestedStruct().FieldA()
	//assert.Equal(t, "fieldA", rFieldA.CurrentString())
	//
	//err := r1.Update(test1.SuperStruct{
	//	String: "new string",
	//	NestedStruct: test1.NestedStruct{
	//		FieldA: "new fieldA",
	//	},
	//})
	//require.NoError(t, err)
	//assert.Equal(t, "new string", rString.CurrentString())
	//assert.Equal(t, "new fieldA", rFieldA.CurrentString())
}
