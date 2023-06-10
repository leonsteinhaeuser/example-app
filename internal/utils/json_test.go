package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteJSON(t *testing.T) {
	type args struct {
		statusCode int
		v          interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write json",
			args: args{
				statusCode: 200,
				v:          nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			if err := WriteJSON(rec, tt.args.statusCode, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("WriteJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadJSON(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "read json",
			args: args{
				data: map[string]any{
					"foo": "bar",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := json.Marshal(tt.args.data)
			if err != nil {
				t.Errorf("ReadJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			rslt := make(map[string]any)
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(dt))
			if err := ReadJSON(req, &rslt); (err != nil) != tt.wantErr {
				t.Errorf("ReadJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			diff := cmp.Diff(tt.args.data, rslt)
			if diff != "" {
				t.Errorf("ReadJSON() = %v", diff)
			}
		})
	}
}
