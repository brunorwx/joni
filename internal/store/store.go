package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/brunorwx/joni/internal/model"
	"go.etcd.io/bbolt"
)

var snippetsBucket = []byte("snippets")

type Store struct {
	db *bbolt.DB
}

func OpenStore(dir string) (*Store, error) {
	path := filepath.Join(dir, "joni.db")
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(snippetsBucket)
		return e
	})
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Add(sn *model.Snippet) (int64, error) {
	if sn == nil {
		return 0, errors.New("nil snippet")
	}
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(snippetsBucket)
		id, _ := b.NextSequence()
		sn.ID = int64(id)
		sn.CreatedAt = time.Now().UTC()
		data, e := json.Marshal(sn)
		if e != nil {
			return e
		}
		key := itob(sn.ID)
		return b.Put(key, data)
	})
	if err != nil {
		return 0, err
	}
	return sn.ID, nil
}

func (s *Store) Get(id int64) (*model.Snippet, error) {
	var sn model.Snippet
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(snippetsBucket)
		v := b.Get(itob(id))
		if v == nil {
			return fmt.Errorf("not found")
		}
		return json.Unmarshal(v, &sn)
	})
	if err != nil {
		return nil, err
	}
	return &sn, nil
}

func (s *Store) List() ([]*model.Snippet, error) {
	var out []*model.Snippet
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(snippetsBucket)
		return b.ForEach(func(k, v []byte) error {
			var sn model.Snippet
			if err := json.Unmarshal(v, &sn); err != nil {
				return err
			}
			out = append(out, &sn)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) Delete(id int64) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(snippetsBucket)
		return b.Delete(itob(id))
	})
}

func (s *Store) Search(q string) ([]*model.Snippet, error) {
	if q == "" {
		return nil, errors.New("empty query")
	}
	all, err := s.List()
	if err != nil {
		return nil, err
	}
	var out []*model.Snippet
	for _, sn := range all {
		if contains(sn.Content, q) || contains(sn.Description, q) || tagsContain(sn.Tags, q) {
			out = append(out, sn)
		}
	}
	return out, nil
}

func tagsContain(tags []string, q string) bool {
	for _, t := range tags {
		if contains(t, q) {
			return true
		}
	}
	return false
}

func contains(hay, needle string) bool {
	return strings.Index(hay, needle) >= 0
}

func itob(v int64) []byte {
	b := make([]byte, 8)
	for i := uint(0); i < 8; i++ {
		b[7-i] = byte(v >> (i * 8))
	}
	return b
}
