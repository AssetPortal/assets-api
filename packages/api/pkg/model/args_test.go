package model

import (
	"testing"
)

func TestPagination_Validate(t *testing.T) {
	tests := []struct {
		name     string
		limit    *int
		offset   *int
		expected Pagination
	}{
		{
			name:   "Limit and Offset are nil",
			limit:  nil,
			offset: nil,
			expected: Pagination{
				Limit:  intPtr(MAX_LIMIT),
				Offset: intPtr(0),
			},
		},
		{
			name:   "Limit is nil, Offset is set",
			limit:  nil,
			offset: intPtr(10),
			expected: Pagination{
				Limit:  intPtr(MAX_LIMIT),
				Offset: intPtr(10),
			},
		},
		{
			name:   "Limit is set, Offset is nil",
			limit:  intPtr(50),
			offset: nil,
			expected: Pagination{
				Limit:  intPtr(50),
				Offset: intPtr(0),
			},
		},
		{
			name:   "Limit exceeds MAX_LIMIT",
			limit:  intPtr(200),
			offset: intPtr(5),
			expected: Pagination{
				Limit:  intPtr(MAX_LIMIT),
				Offset: intPtr(5),
			},
		},
		{
			name:   "Limit and Offset are set within bounds",
			limit:  intPtr(50),
			offset: intPtr(5),
			expected: Pagination{
				Limit:  intPtr(50),
				Offset: intPtr(5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Limit:  tt.limit,
				Offset: tt.offset,
			}
			p.Validate()
			if *p.Limit != *tt.expected.Limit {
				t.Errorf("expected Limit %d, got %d", *tt.expected.Limit, *p.Limit)
			}
			if *p.Offset != *tt.expected.Offset {
				t.Errorf("expected Offset %d, got %d", *tt.expected.Offset, *p.Offset)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
