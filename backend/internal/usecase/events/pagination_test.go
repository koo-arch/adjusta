package events

import (
	"testing"
)

func TestBuildPaginationOutput(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		page       int
		perPage    int
		totalItems int
		want       PaginationOutput
	}{
		{
			name:       "割り切れる件数",
			page:       1,
			perPage:    20,
			totalItems: 40,
			want:       PaginationOutput{Page: 1, PerPage: 20, TotalItems: 40, TotalPages: 2},
		},
		{
			name:       "余りがある件数は切り上げ",
			page:       2,
			perPage:    20,
			totalItems: 41,
			want:       PaginationOutput{Page: 2, PerPage: 20, TotalItems: 41, TotalPages: 3},
		},
		{
			name:       "0件は totalPages 0",
			page:       1,
			perPage:    20,
			totalItems: 0,
			want:       PaginationOutput{Page: 1, PerPage: 20, TotalItems: 0, TotalPages: 0},
		},
		{
			name:       "perPage 未指定は既定値 20 に補正",
			page:       0,
			perPage:    0,
			totalItems: 5,
			want:       PaginationOutput{Page: 1, PerPage: 20, TotalItems: 5, TotalPages: 1},
		},
		{
			name:       "perPage が totalItems 未満",
			page:       1,
			perPage:    3,
			totalItems: 10,
			want:       PaginationOutput{Page: 1, PerPage: 3, TotalItems: 10, TotalPages: 4},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := buildPaginationOutput(tc.page, tc.perPage, tc.totalItems)
			if got != tc.want {
				t.Fatalf("buildPaginationOutput(%d, %d, %d) = %+v, want %+v", tc.page, tc.perPage, tc.totalItems, got, tc.want)
			}
		})
	}
}

func TestToEventSearchOptionsPagination(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		page       int
		perPage    int
		wantOffset int
		wantLimit  int
	}{
		{
			name:       "1ページ目は offset 0",
			page:       1,
			perPage:    20,
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			name:       "3ページ目は (page-1)*perPage",
			page:       3,
			perPage:    10,
			wantOffset: 20,
			wantLimit:  10,
		},
		{
			name:       "page 未指定は 1 ページ目扱い",
			page:       0,
			perPage:    20,
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			name:       "perPage 0 は無制限(offset/limit なし)",
			page:       5,
			perPage:    0,
			wantOffset: 0,
			wantLimit:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := toEventSearchOptions(EventSearchOptions{Page: tc.page, PerPage: tc.perPage})
			if got.EventOffset != tc.wantOffset || got.EventLimit != tc.wantLimit {
				t.Fatalf("toEventSearchOptions(page=%d, perPage=%d) = offset %d / limit %d, want offset %d / limit %d",
					tc.page, tc.perPage, got.EventOffset, got.EventLimit, tc.wantOffset, tc.wantLimit)
			}
		})
	}
}
