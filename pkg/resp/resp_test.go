package resp

import (
	"bufio"
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestRespConn_readLine(t *testing.T) {
	type fields struct {
		Reader *bufio.Reader
		Writer *bufio.Writer
		conn   net.Conn
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "valid int",
			fields: fields{
				Reader: bufio.NewReader(strings.NewReader("6\r\n")),
			},
			want:    []byte("6"),
			wantErr: false,
		},
		{
			name: "valid string",
			fields: fields{
				Reader: bufio.NewReader(strings.NewReader("test\r\n")),
			},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name: "invalid - missing CR",
			fields: fields{
				Reader: bufio.NewReader(strings.NewReader("test\n")),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RespConn{
				Reader: tt.fields.Reader,
				Writer: tt.fields.Writer,
				conn:   tt.fields.conn,
			}
			got, err := r.readLine()
			if (err != nil) != tt.wantErr {
				t.Errorf("RespConn.readLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RespConn.readLine() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestRespConn_readInt(t *testing.T) {
	type fields struct {
		Reader *bufio.Reader
		Writer *bufio.Writer
		conn   net.Conn
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "valid int",
			fields: fields{
				Reader: bufio.NewReader(strings.NewReader("1234\r\n")),
			},
			want:    1234,
			wantErr: false,
		},
		//TODO: test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RespConn{
				Reader: tt.fields.Reader,
				Writer: tt.fields.Writer,
				conn:   tt.fields.conn,
			}
			got, err := r.readInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("RespConn.readInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RespConn.readInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
