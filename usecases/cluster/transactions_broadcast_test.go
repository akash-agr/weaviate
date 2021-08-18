package cluster

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcastOpenTransaction(t *testing.T) {
	client := &fakeClient{}
	state := &fakeState{[]string{"host1", "host2", "host3"}}

	bc := NewTxBroadcaster(state, client)

	tx := &Transaction{ID: "foo"}

	err := bc.BroadcastTransaction(context.Background(), tx)
	require.Nil(t, err)

	assert.ElementsMatch(t, []string{"host1", "host2", "host3"}, client.openCalled)
}

func TestBroadcastAbortTransaction(t *testing.T) {
	client := &fakeClient{}
	state := &fakeState{[]string{"host1", "host2", "host3"}}

	bc := NewTxBroadcaster(state, client)

	tx := &Transaction{ID: "foo"}

	err := bc.BroadcastAbortTransaction(context.Background(), tx)
	require.Nil(t, err)

	assert.ElementsMatch(t, []string{"host1", "host2", "host3"}, client.abortCalled)
}

func TestBroadcastCommitTransaction(t *testing.T) {
	client := &fakeClient{}
	state := &fakeState{[]string{"host1", "host2", "host3"}}

	bc := NewTxBroadcaster(state, client)

	tx := &Transaction{ID: "foo"}

	err := bc.BroadcastCommitTransaction(context.Background(), tx)
	require.Nil(t, err)

	assert.ElementsMatch(t, []string{"host1", "host2", "host3"}, client.commitCalled)
}

type fakeState struct {
	hosts []string
}

func (f *fakeState) Hostnames() []string {
	return f.hosts
}

type fakeClient struct {
	sync.Mutex
	openCalled   []string
	abortCalled  []string
	commitCalled []string
}

func (f *fakeClient) OpenTransaction(ctx context.Context, host string, tx *Transaction) error {
	f.Lock()
	defer f.Unlock()

	f.openCalled = append(f.openCalled, host)
	return nil
}

func (f *fakeClient) AbortTransaction(ctx context.Context, host string, tx *Transaction) error {
	f.Lock()
	defer f.Unlock()

	f.abortCalled = append(f.abortCalled, host)
	return nil
}

func (f *fakeClient) CommitTransaction(ctx context.Context, host string, tx *Transaction) error {
	f.Lock()
	defer f.Unlock()

	f.commitCalled = append(f.commitCalled, host)
	return nil
}