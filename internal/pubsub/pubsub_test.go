package pubsub

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestDefaultEvent_ID(t *testing.T) {
	type fields struct {
		ResourceID     uuid.UUID
		ActionType     ActionType
		AdditionalData map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "should return the resource id",
			fields: fields{
				ResourceID: uuid.MustParse("cfd2e31e-8a0b-4fd2-8af7-38cbaf2e05f7"),
			},
			want: uuid.MustParse("cfd2e31e-8a0b-4fd2-8af7-38cbaf2e05f7"),
		},
		{
			name: "should return the empty id",
			fields: fields{
				ResourceID: uuid.UUID{},
			},
			want: uuid.UUID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultEvent{
				ResourceID:     tt.fields.ResourceID,
				ActionType:     tt.fields.ActionType,
				AdditionalData: tt.fields.AdditionalData,
			}
			if got := d.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultEvent.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultEvent_Action(t *testing.T) {
	type fields struct {
		ResourceID     uuid.UUID
		ActionType     ActionType
		AdditionalData map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		want   ActionType
	}{
		{
			name: "should return the action type create",
			fields: fields{
				ActionType: ActionTypeCreate,
			},
			want: ActionTypeCreate,
		},
		{
			name: "should return the action type update",
			fields: fields{
				ActionType: ActionTypeUpdate,
			},
			want: ActionTypeUpdate,
		},
		{
			name: "should return the action type delete",
			fields: fields{
				ActionType: ActionTypeDelete,
			},
			want: ActionTypeDelete,
		},
		{
			name: "empty",
			fields: fields{
				ActionType: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultEvent{
				ResourceID:     tt.fields.ResourceID,
				ActionType:     tt.fields.ActionType,
				AdditionalData: tt.fields.AdditionalData,
			}
			if got := d.Action(); got != tt.want {
				t.Errorf("DefaultEvent.Action() = %v, want %v", got, tt.want)
			}
		})
	}
}
