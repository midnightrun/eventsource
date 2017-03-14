package sqlstore_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/savaki/eventsource"
	"github.com/savaki/eventsource/provider/sqlstore"
	"github.com/stretchr/testify/assert"
)

type EntitySetFirst struct {
	eventsource.Model
	First string
}

type EntitySetLast struct {
	eventsource.Model
	Last string
}

func TestStore_Save(t *testing.T) {
	ctx := context.Background()
	tableName := "entity_events"

	// Ensure table exists

	db := MustOpen()
	err := sqlstore.CreateMySQL(ctx, db, tableName)
	assert.Nil(t, err)
	db.Close()

	// Return

	aggregateID := strconv.FormatInt(time.Now().UnixNano(), 10)
	first := EntitySetFirst{
		Model: eventsource.Model{
			ID:      aggregateID,
			Version: 1,
		},
		First: "first",
	}
	second := EntitySetLast{
		Model: eventsource.Model{
			ID:      aggregateID,
			Version: 2,
		},
		Last: "last",
	}

	serializer := eventsource.JSONSerializer()
	serializer.Bind(first, second)

	store := sqlstore.New(tableName, Open, sqlstore.WithDebug(os.Stderr))

	err = store.Save(context.Background(), serializer, first, second)
	assert.Nil(t, err)

	events, version, err := store.Fetch(context.Background(), serializer, aggregateID, 0)
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{&first, &second}, events)
	assert.Equal(t, second.Model.Version, version)
}
