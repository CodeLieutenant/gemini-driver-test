package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strings"
	"time"

	_ "embed"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

func must[T any](data T, err error) T {
	if err != nil {
		log.Fatalf("An Error has occurred: %+v", err)
	}

	return data
}

var (
	pk0Value1 = `[b573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe90934618df4bb22a8ec7deb573537bbe9093461]`
	pk0Value2 = `[5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe2568226e19f5b163b293932ba3c2c4fe]`

	pk1Value1 = must(gocql.ParseUUID("82aa8b9e-5c3b-1d6b-b5e6-fc5cee71c66c"))
	pk1Value2 = must(gocql.ParseUUID("4a5194b1-4bd4-1b14-9a6a-fc5cee71c66c"))
)

// pk0  timeuuid,
// pk1  date,
// ck0  float,
// ck1  ascii,
// col0 inet,
// col1 frozen<list<timeuuid>>,
// col2 text,
// col3 date,
// col4 frozen<set<time>>,
// col5 double,
// col6 blob,
// col7 text,

type Data struct {
	pk0  gocql.UUID
	pk1  time.Time
	ck0  float32
	ck1  string
	col0 net.IP
	col1 []gocql.UUID
	col2 string
	col3 time.Time
	col4 []time.Time
	col5 float64
	col6 []byte
	col7 string
}

//go:embed schema.cql
var Schema string

//go:embed query.cql
var Query string

var Keyspace string = "CREATE KEYSPACE IF NOT EXISTS test WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1}"

var (
	hosts       string
	cqlVersion  string
	consistency string

	drop bool
)

func scyllaHosts(hosts string) []string {
	h := strings.Split(hosts, ",")

	for i := range h {
		h[i] = strings.TrimSpace(h[i])
	}

	return h
}

func withoutPreparedStatement(ctx context.Context, session *gocql.Session, stmt string) error {
	query := session.Query(stmt)
	defer query.Release()

	log.Printf("Query: %s", stmt)
	if err := query.WithContext(ctx).Exec(); err != nil {
		return errors.Wrap(err, "Failed to Execute Query table")
	}

	var data Data

	err := query.Scan(
		&data.pk0,
		&data.pk1,
		&data.ck0,
		&data.ck1,
		&data.col0,
		&data.col1,
		&data.col2,
		&data.col3,
		&data.col4,
		&data.col5,
		&data.col6,
		&data.col7,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to scan data")
	}

	log.Printf("Data: %+v", data)

	return nil
}

func withPreparedStatement(ctx context.Context, session *gocql.Session, stmt string) error {
	query := session.Query(stmt, pk0Value1, pk0Value2, pk1Value1, pk1Value2)
	defer query.Release()

	log.Printf("Query: %s", stmt)
	if err := query.WithContext(ctx).Exec(); err != nil {
		return errors.Wrap(err, "Failed to Execute Query table")
	}

	var data Data

	err := query.Scan(
		&data.pk0,
		&data.pk1,
		&data.ck0,
		&data.ck1,
		&data.col0,
		&data.col1,
		&data.col2,
		&data.col3,
		&data.col4,
		&data.col5,
		&data.col6,
		&data.col7,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to scan data")
	}

	log.Printf("Data: %+v", data)

	return nil
}

func main() {
	flag.StringVar(&hosts, "hosts", "localhost", "ScyllaDB hosts")
	flag.StringVar(&cqlVersion, "cql-version", "3.0.0", "CQL version")
	flag.StringVar(&consistency, "consistency", "LOCAL_QUORUM", "Consistency level")
	flag.BoolVar(&drop, "drop-keyspace", true, "Drop Keyspace at the and")

	flag.Parse()

	ctx := context.Background()

	cluster := gocql.NewCluster(scyllaHosts(hosts)...)

	cluster.Timeout = 1 * time.Second
	cluster.ConnectTimeout = 1 * time.Second
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        60 * time.Second,
		NumRetries: 5,
	}
	cluster.CQLVersion = cqlVersion
	cluster.Consistency = gocql.ParseConsistency(consistency)
	cluster.ProtoVersion = 4
	cluster.DefaultIdempotence = true
	cluster.NumConns = 100
	cluster.MaxPreparedStmts = 100
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("datacenter1"))

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to create ScyllDB Session: %v", err)
	}

	defer session.Close()

	log.Printf("Connected to ScyllaDB: %s", hosts)
	log.Printf("Creating Keyspace: %s", Keyspace)

	if err = session.Query(Keyspace).WithContext(ctx).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	query := session.Query(Schema)

	log.Printf("Creating Schema: %s", query.String())
	if err = query.WithContext(ctx).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	if drop {
		defer func() {
			if err = session.Query("DROP KEYSPACE test;").WithContext(ctx).Exec(); err != nil {
				log.Fatalf("Failed to drop table: %v", err)
			}
		}()
	}

	stmts := strings.Split(Query, ";")

	if err = withPreparedStatement(ctx, session, stmts[1]); err != nil {
		log.Printf("Failed to execute query with prepared statement: %+v", err)
	}

	if err = withoutPreparedStatement(ctx, session, stmts[0]); err != nil {
		log.Printf("Failed to execute query without prepared statement: %+v", err)
	}
}
