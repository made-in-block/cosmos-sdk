package ormtable_test

import (
	"bytes"
	"testing"

	"github.com/cosmos/cosmos-sdk/orm/model/ormtable"

	"google.golang.org/protobuf/testing/protocmp"

	"gotest.tools/v3/assert"

	"github.com/cosmos/cosmos-sdk/orm/internal/testkv"
	"github.com/cosmos/cosmos-sdk/orm/internal/testpb"
)

func TestSingleton(t *testing.T) {
	val := &testpb.ExampleSingleton{}
	singleton, err := ormtable.Build(ormtable.Options{
		MessageType: val.ProtoReflect().Type(),
	})
	assert.NilError(t, err)
	store := testkv.NewSplitMemBackend()

	found, err := singleton.Has(store, nil)
	assert.NilError(t, err)
	assert.Assert(t, !found)
	assert.NilError(t, singleton.Save(store, val, ormtable.SAVE_MODE_DEFAULT))
	found, err = singleton.Has(store, nil)
	assert.NilError(t, err)
	assert.Assert(t, found)

	val.Foo = "abc"
	val.Bar = 3
	assert.NilError(t, singleton.Save(store, val, ormtable.SAVE_MODE_DEFAULT))

	var val2 testpb.ExampleSingleton
	found, err = singleton.Get(store, nil, &val2)
	assert.NilError(t, err)
	assert.DeepEqual(t, val, &val2, protocmp.Transform())

	buf := &bytes.Buffer{}
	assert.NilError(t, singleton.ExportJSON(store, buf))
	assert.NilError(t, singleton.ValidateJSON(bytes.NewReader(buf.Bytes())))
	store2 := testkv.NewSplitMemBackend()
	assert.NilError(t, singleton.ImportJSON(store2, bytes.NewReader(buf.Bytes())))

	var val3 testpb.ExampleSingleton
	found, err = singleton.Get(store, nil, &val3)
	assert.NilError(t, err)
	assert.DeepEqual(t, val, &val3, protocmp.Transform())
}
