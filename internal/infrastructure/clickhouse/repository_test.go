package clickhouse_test

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/mahsandr/arman-challenge/config"
	"github.com/mahsandr/arman-challenge/internal/domain/models"
	"github.com/mahsandr/arman-challenge/internal/infrastructure/clickhouse"
	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v2/chpool"
)

func TestAddUserSegments(t *testing.T) {
	tests := []struct {
		name     string
		segments []*models.UserSegment
		wantErr  bool
	}{
		{
			name: "valid segments",
			segments: []*models.UserSegment{
				{UserID: "user1", Segment: "segment1", RegistredAt: uint32(time.Now().Unix())},
				{UserID: "user2", Segment: "segment2", RegistredAt: uint32(time.Now().Unix())},
			},
			wantErr: false,
		},
	}

	// Setup ClickHouse connection
	cfg, _ := config.GetConfig("../../.env")

	conn, err := chpool.New(cfg.ClickHouseAddr)
	if err != nil {
		t.Fatalf("failed to create ClickHouse connection: %v", err)
	}
	defer conn.Close()
	setupTestTable(t, conn)
	defer downTestTable(t, conn)
	repo, err := clickhouse.NewClickHouseRepository(conn, "testsegments", "")
	if err != nil {
		t.Fatalf("failed to create ClickHouse repository: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SaveUserSegments(context.Background(), tt.segments)
			if tt.wantErr {
				assert.Error(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestGetSegmentCount(t *testing.T) {
	tests := []struct {
		name      string
		segment   string
		setupData []*models.UserSegment
		wantCount uint32
		wantErr   bool
	}{
		{
			name:    "segment exists",
			segment: "segment1",
			setupData: []*models.UserSegment{
				{UserID: "user1", Segment: "segment1", RegistredAt: uint32(time.Now().Unix())},
				{UserID: "user2", Segment: "segment1", RegistredAt: uint32(time.Now().Unix())},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:    "segment does not exist",
			segment: "segment2",
			setupData: []*models.UserSegment{
				{UserID: "user1", Segment: "segment1", RegistredAt: uint32(time.Now().Unix())},
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	// Setup ClickHouse connection
	cfg, _ := config.GetConfig("../../.env")

	conn, err := chpool.New(cfg.ClickHouseAddr)
	if err != nil {
		t.Fatalf("failed to create ClickHouse connection: %v", err)
	}
	defer conn.Close()
	setupTestTable(t, conn)
	defer downTestTable(t, conn)

	repo, err := clickhouse.NewClickHouseRepository(conn, "testsegments", "testsegmentsview")
	if err != nil {
		t.Fatalf("failed to create ClickHouse repository: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Insert setup data
			err := repo.SaveUserSegments(context.Background(), tt.setupData)
			if err != nil {
				t.Fatalf("failed to insert setup data: %v", err)
			}

			// Test GetSegmentCount
			count, err := repo.GetSegmentUsersCount(context.Background(), tt.segment)
			if tt.wantErr {
				assert.Error(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
			}
		})
	}
}

//go:embed test_queries/create_segment_table.sql
var createTestSegmentQuery string

//go:embed test_queries/create_segmenet_view.sql
var createSegmentViewQuery string

//go:embed test_queries/create_temp_table.sql
var createTempTableQuery string

func setupTestTable(t *testing.T, conn chpool.Pool) {
	err := conn.Exec(context.Background(), createTestSegmentQuery)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	err = conn.Exec(context.Background(), createTempTableQuery)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	err = conn.Exec(context.Background(), createSegmentViewQuery)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
}
func downTestTable(t *testing.T, conn chpool.Pool) {
	query := `DROP TABLE IF EXISTS testsegments`
	err := conn.Exec(context.Background(), query)
	if err != nil {
		t.Fatalf("failed to drop test table: %v", err)
	}
	query = `DROP TABLE IF EXISTS testsegmentsview`
	err = conn.Exec(context.Background(), query)
	if err != nil {
		t.Fatalf("failed to drop test table: %v", err)
	}
	query = `DROP TABLE IF EXISTS temptable`
	err = conn.Exec(context.Background(), query)
	if err != nil {
		t.Fatalf("failed to drop test table: %v", err)
	}
}
