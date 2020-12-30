package domain

import (
	"context"
	"time"
)

// RecordTimestamp ...
type RecordTimestamp struct {
	CreatedAt      time.Time `json:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	LastModifiedBy int       `json:"lastModifiedBy"`
	DeletedAt      time.Time `pg:",soft_delete" json:"deletedAt"`
}

var recordTimestampGroupMetadata = GetMetadataMap(&RecordTimestamp{})

// GetRecordTimestampSQLField ...
func GetRecordTimestampSQLField(jsonField string) string {
	if metadata, ok := recordTimestampGroupMetadata[jsonField]; ok {
		return metadata.SQLTag
	}

	return ""
}

// BeforeUpdate ...
func (r *RecordTimestamp) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if r.LastModifiedAt.IsZero() {
		r.LastModifiedAt = time.Now()
	}
	return ctx, nil
}

// BeforeInsert ...
func (r *RecordTimestamp) BeforeInsert(ctx context.Context) (context.Context, error) {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	return ctx, nil
}
