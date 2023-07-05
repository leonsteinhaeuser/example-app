package thread

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func TestThreadFromThreadModel(t *testing.T) {
	type args struct {
		m *ThreadModel
	}
	tests := []struct {
		name string
		args args
		want *Thread
	}{
		{
			name: "Converts a ThreadModel to a Thread",
			args: args{
				m: &ThreadModel{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: gorm.DeletedAt{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					Title:     "Test Title",
					Body:      "Test Body",
					AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					KeywordIDs: pq.StringArray{
						"00000000-0000-0000-0000-000000000000",
						"00000000-0000-0000-0000-000000000001",
					},
				},
			},
			want: &Thread{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				DeletedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Title:     "Test Title",
				Body:      "Test Body",
				AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				KeywordIDs: []uuid.UUID{
					uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ThreadFromThreadModel(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ThreadFromThreadModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThreadModelFromThread(t *testing.T) {
	type args struct {
		t *Thread
	}
	tests := []struct {
		name string
		args args
		want *ThreadModel
	}{
		{
			name: "Converts a Thread to a ThreadModel",
			args: args{
				t: &Thread{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Title:     "Test Title",
					Body:      "Test Body",
					AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					KeywordIDs: []uuid.UUID{
						uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					},
				},
			},
			want: &ThreadModel{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				DeletedAt: gorm.DeletedAt{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				Title:     "Test Title",
				Body:      "Test Body",
				AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				KeywordIDs: pq.StringArray{
					"00000000-0000-0000-0000-000000000000",
					"00000000-0000-0000-0000-000000000001",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ThreadModelFromThread(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ThreadModelFromThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThreadsFromThreadModels(t *testing.T) {
	type args struct {
		models []*ThreadModel
	}
	tests := []struct {
		name string
		args args
		want []*Thread
	}{
		{
			name: "Converts a slice of ThreadModels to Threads",
			args: args{
				models: []*ThreadModel{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						DeletedAt: gorm.DeletedAt{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
						Title:     "Test Title",
						Body:      "Test Body",
						AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						KeywordIDs: pq.StringArray{
							"00000000-0000-0000-0000-000000000000",
							"00000000-0000-0000-0000-000000000001",
						},
					},
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						DeletedAt: gorm.DeletedAt{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
						Title:     "Test Title 2",
						Body:      "Test Body 2",
						AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						KeywordIDs: pq.StringArray{
							"00000000-0000-0000-0000-000000000000",
							"00000000-0000-0000-0000-000000000001",
						},
					},
				},
			},
			want: []*Thread{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Title:     "Test Title",
					Body:      "Test Body",
					AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					KeywordIDs: []uuid.UUID{
						uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					},
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Title:     "Test Title 2",
					Body:      "Test Body 2",
					AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					KeywordIDs: []uuid.UUID{
						uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ThreadsFromThreadModels(tt.args.models)

			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf("ThreadsFromThreadModels() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestThreadModelsFromThreads(t *testing.T) {
	type args struct {
		threads []*Thread
	}
	tests := []struct {
		name string
		args args
		want []*ThreadModel
	}{
		{
			name: "Converts a slice of Threads to ThreadModels",
			args: args{
				threads: []*Thread{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						DeletedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Title:     "Test Title",
						Body:      "Test Body",
						AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						KeywordIDs: []uuid.UUID{
							uuid.MustParse("00000000-0000-0000-0000-000000000000"),
							uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						},
					},
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						DeletedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Title:     "Test Title 2",
						Body:      "Test Body 2",
						AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						KeywordIDs: []uuid.UUID{
							uuid.MustParse("00000000-0000-0000-0000-000000000000"),
							uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						},
					},
				},
			},
			want: []*ThreadModel{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: gorm.DeletedAt{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					Title:     "Test Title",
					Body:      "Test Body",
					AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					KeywordIDs: pq.StringArray{
						"00000000-0000-0000-0000-000000000000",
						"00000000-0000-0000-0000-000000000001",
					},
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					DeletedAt: gorm.DeletedAt{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					Title:     "Test Title 2",
					Body:      "Test Body 2",
					AuthorID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					KeywordIDs: pq.StringArray{
						"00000000-0000-0000-0000-000000000000",
						"00000000-0000-0000-0000-000000000001",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ThreadModelsFromThreads(tt.args.threads)

			if len(got) != len(tt.want) {
				t.Errorf("ThreadModelsFromThreads() expected length = %d, got %d", len(tt.want), len(got))
			}

			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ThreadModelsFromThreads() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
