package db

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/semi-technologies/weaviate/adapters/repos/db/helpers"
	"github.com/semi-technologies/weaviate/adapters/repos/db/inverted"
	"github.com/semi-technologies/weaviate/adapters/repos/db/lsmkv"
)

func (s *Shard) deleteFromInvertedIndicesLSM(props []inverted.Property,
	docID uint64) error {
	for _, prop := range props {
		b := s.store.Bucket(helpers.BucketFromPropNameLSM(prop.Name))
		if b == nil {
			return fmt.Errorf("no bucket for prop '%s' found", prop.Name)
		}

		hashBucket := s.store.Bucket(helpers.HashBucketFromPropNameLSM(prop.Name))
		if b == nil {
			return fmt.Errorf("no hash bucket for prop '%s' found", prop.Name)
		}

		if prop.HasFrequency {
			for _, item := range prop.Items {
				if err := s.deleteInvertedIndexItemWithFrequencyLSM(b, hashBucket, item,
					docID); err != nil {
					return errors.Wrapf(err, "extend index with item '%s'",
						string(item.Data))
				}
			}
		} else {
			for _, item := range prop.Items {
				if err := s.deleteInvertedIndexItemLSM(b, hashBucket, item, docID); err != nil {
					return errors.Wrapf(err, "extend index with item '%s'",
						string(item.Data))
				}
			}
		}
	}

	return nil
}

func (s *Shard) deleteInvertedIndexItemWithFrequencyLSM(b, hashBucket *lsmkv.Bucket,
	item inverted.Countable, docID uint64) error {
	if b.Strategy() != lsmkv.StrategyMapCollection {
		panic("prop has frequency, but bucket does not have 'Map' strategy")
	}

	hash, err := generateRowHash()
	if err != nil {
		return err
	}

	if err := hashBucket.Put(item.Data, hash); err != nil {
		return err
	}

	docIDBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(docIDBytes, docID)

	return b.MapDeleteKey(item.Data, docIDBytes)
}

func (s *Shard) deleteInvertedIndexItemLSM(b, hashBucket *lsmkv.Bucket,
	item inverted.Countable, docID uint64) error {
	if b.Strategy() != lsmkv.StrategySetCollection {
		panic("prop has no frequency, but bucket does not have 'Set' strategy")
	}

	hash, err := generateRowHash()
	if err != nil {
		return err
	}

	if err := hashBucket.Put(item.Data, hash); err != nil {
		return err
	}

	docIDBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(docIDBytes, docID)

	return b.SetDeleteSingle(item.Data, docIDBytes)
}