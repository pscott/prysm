package kv

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"go.etcd.io/bbolt"
)

func Test_migrateArchivedIndex(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, db *bbolt.DB)
		eval  func(t *testing.T, db *bbolt.DB)
	}{
		{
			name: "only runs once",
			setup: func(t *testing.T, db *bbolt.DB) {
				err := db.Update(func(tx *bbolt.Tx) error {
					_, err := tx.CreateBucketIfNotExists(archivedRootBucket)
					assert.NoError(t, err)
					if err := tx.Bucket(archivedRootBucket).Put(bytesutil.Uint64ToBytesLittleEndian(2048), []byte("foo")); err != nil {
						return err
					}
					return tx.Bucket(migrationsBucket).Put(migrationArchivedIndex0Key, migrationCompleted)
				})
				assert.NoError(t, err)
			},
			eval: func(t *testing.T, db *bbolt.DB) {
				err := db.View(func(tx *bbolt.Tx) error {
					v := tx.Bucket(archivedRootBucket).Get(bytesutil.Uint64ToBytesLittleEndian(2048))
					if !bytes.Equal(v, []byte("foo")) {
						return fmt.Errorf("did not receive correct data for key 2048, wanted 'foo' got %s", v)
					}
					return nil
				})
				assert.NoError(t, err)
			},
		},
		{
			name: "migrates and deletes entries",
			setup: func(t *testing.T, db *bbolt.DB) {
				err := db.Update(func(tx *bbolt.Tx) error {
					_, err := tx.CreateBucketIfNotExists(archivedRootBucket)
					assert.NoError(t, err)
					_, err = tx.CreateBucketIfNotExists(slotsHasObjectBucket)
					assert.NoError(t, err)
					if err := tx.Bucket(archivedRootBucket).Put(bytesutil.Uint64ToBytesLittleEndian(2048), []byte("foo")); err != nil {
						return err
					}

					b, err := encode(context.Background(), &ethpb.SignedBeaconBlock{Block: &ethpb.BeaconBlock{Slot: 2048}})
					if err != nil {
						return err
					}
					return tx.Bucket(blocksBucket).Put([]byte("foo"), b)
				})
				assert.NoError(t, err)
			},
			eval: func(t *testing.T, db *bbolt.DB) {
				err := db.View(func(tx *bbolt.Tx) error {
					k := uint64(2048)
					if v := tx.Bucket(stateSlotIndicesBucket).Get(bytesutil.Uint64ToBytesBigEndian(k)); !bytes.Equal(v, []byte("foo")) {
						return fmt.Errorf("did not receive correct data for key %d, wanted 'foo' got %v", k, v)
					}
					return nil
				})
				assert.NoError(t, err)
			},
		},
		{
			name: "deletes old buckets",
			setup: func(t *testing.T, db *bbolt.DB) {
				err := db.Update(func(tx *bbolt.Tx) error {
					_, err := tx.CreateBucketIfNotExists(archivedRootBucket)
					assert.NoError(t, err)
					_, err = tx.CreateBucketIfNotExists(slotsHasObjectBucket)
					assert.NoError(t, err)
					return tx.Bucket(slotsHasObjectBucket).Put(savedStateSlotsKey, []byte("foo"))
				})
				assert.NoError(t, err)
			},
			eval: func(t *testing.T, db *bbolt.DB) {
				err := db.View(func(tx *bbolt.Tx) error {
					assert.Equal(t, (*bbolt.Bucket)(nil), tx.Bucket(slotsHasObjectBucket), "Expected %v to be deleted", savedStateSlotsKey)
					assert.Equal(t, (*bbolt.Bucket)(nil), tx.Bucket(archivedRootBucket), "Expected %v to be deleted", savedStateSlotsKey)
					return nil
				})
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t).db
			tt.setup(t, db)
			assert.NoError(t, db.Update(migrateArchivedIndex), "migrateArchivedIndex(tx) error")
			tt.eval(t, db)
		})
	}
}
