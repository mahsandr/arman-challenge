package clickhouse

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/mahsandr/arman-challenge/internal/domain/models"
	"github.com/mahsandr/arman-challenge/internal/domain/repository"
	"github.com/vahid-sohrabloo/chconn/v2"
	"github.com/vahid-sohrabloo/chconn/v2/chpool"
	"github.com/vahid-sohrabloo/chconn/v2/column"
)

//go:embed queries/insert_segments.sql
var insertSegmentsQuery string

var _ repository.SegmentRepository = &ClickHouseRepository{}

type ClickHouseRepository struct {
	conn         chpool.Pool
	segmentsTbl  string
	segmentsView string
}

// ClickHouseRepository return and initialize clickhouse object
func NewClickHouseRepository(conn chpool.Pool, segmentsTbl, segmentsView string) (*ClickHouseRepository, error) {
	return &ClickHouseRepository{
		conn:         conn,
		segmentsTbl:  segmentsTbl,
		segmentsView: segmentsView,
	}, nil
}

func (c *ClickHouseRepository) SaveUserSegments(ctx context.Context, segments []*models.UserSegment) error {
	userIDCol := column.NewString().LowCardinality()
	segmentCol := column.NewString().LowCardinality()
	registeredAtCol := column.New[uint32]()

	for _, segment := range segments {
		userIDCol.Append(segment.UserID)
		segmentCol.Append(segment.Segment)
		registeredAtCol.Append(segment.RegistredAt)
	}
	err := c.conn.InsertWithOption(ctx, insertSegmentsQuery,
		&chconn.QueryOptions{
			Parameters: chconn.NewParameters(chconn.StringParameter("tableName", c.segmentsTbl)),
		},
		userIDCol, segmentCol, registeredAtCol)
	if err != nil {
		return fmt.Errorf("error inserting user segments: %w", err)
	}
	return nil
}

//go:embed queries/count_users.sql
var selectSegmentCountQuery string

func (c *ClickHouseRepository) GetSegmentUsersCount(ctx context.Context, segment string) (uint32, error) {
	countCol := column.New[uint64]()
	stmt, err := c.conn.SelectWithOption(
		context.Background(),
		selectSegmentCountQuery,
		&chconn.QueryOptions{
			Parameters: chconn.NewParameters(
				chconn.StringParameter("viewName", c.segmentsView),
				chconn.StringParameter("segment", segment)),
		},
		countCol,
	)
	if err != nil {
		return 0, fmt.Errorf("error on select:%w", err)
	}
	defer stmt.Close()
	if !stmt.Next() {
		return 0, nil
	}
	count := countCol.Row(0)

	if stmt.Err() != nil {
		return 0, stmt.Err()
	}
	return uint32(count), nil
}
