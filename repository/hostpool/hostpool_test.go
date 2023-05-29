package hostpool

import (
	"context"
	"reflect"
	"testing"
)

func TestRepository_AppendHost(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		mock     func()
		wantPool []string
		wantMap  map[string]bool
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				host: "host_4",
			},
			wantErr: false,
			mock: func() {
				pool = append(pool, "host_1", "host_2", "host_3")
				poolMap = map[string]bool{
					"host_1": true,
					"host_2": true,
					"host_3": true,
				}
			},
			wantPool: []string{"host_1", "host_2", "host_3", "host_4"},
			wantMap: map[string]bool{
				"host_1": true,
				"host_2": true,
				"host_3": true,
				"host_4": true,
			},
		},
		{
			name: "success duplicate",
			args: args{
				ctx:  context.Background(),
				host: "host_3",
			},
			wantErr: false,
			mock: func() {
				pool = append(pool, "host_1", "host_2", "host_3")
				poolMap = map[string]bool{
					"host_1": true,
					"host_2": true,
					"host_3": true,
				}
			},
			wantPool: []string{"host_1", "host_2", "host_3"},
			wantMap: map[string]bool{
				"host_1": true,
				"host_2": true,
				"host_3": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				pool = []string{}
				poolMap = map[string]bool{}
			}()

			if err := r.AppendHost(tt.args.ctx, tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("AppendHost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(pool, tt.wantPool) {
				t.Errorf("AppendHost() pool = %v, wantPool %v", pool, tt.wantPool)
			}
			if !reflect.DeepEqual(poolMap, tt.wantMap) {
				t.Errorf("AppendHost() poolMap = %v, wantMap %v", poolMap, tt.wantMap)
			}
		})
	}
}

func TestRepository_GetCurrentIndex(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantRes int
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			wantRes: 11,
			wantErr: false,
			mock: func() {
				index = 11
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				index = 0
			}()

			gotRes, err := r.GetCurrentIndex(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GetCurrentIndex() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestRepository_GetHostListFromPool(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			wantRes: []string{"host1"},
			wantErr: false,
			mock: func() {
				pool = []string{"host1"}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				pool = []string{}
			}()

			gotRes, err := r.GetHostListFromPool(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostListFromPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetHostListFromPool() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestRepository_GetHostListLength(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantRes int
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			wantRes: 1,
			wantErr: false,
			mock: func() {
				pool = []string{"host1"}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				pool = []string{}
			}()

			gotRes, err := r.GetHostListLength(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostListLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GetHostListLength() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestRepository_IncrementIndex(t *testing.T) {
	type args struct {
		ctx       context.Context
		increment int
	}
	tests := []struct {
		name    string
		args    args
		wantRes int
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				ctx:       context.Background(),
				increment: 2,
			},
			wantRes: 1,
			wantErr: false,
			mock: func() {
				pool = []string{"host1", "host2", "host3", "host4"}
				index = 3
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				pool = []string{}
				index = 0
			}()

			gotRes, err := r.IncrementIndex(tt.args.ctx, tt.args.increment)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrementIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("IncrementIndex() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestRepository_RemoveHostByHostAddress(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		mock     func()
		wantPool []string
		wantMap  map[string]bool
		wantIdx  int
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				host: "host3",
			},
			wantErr: false,
			mock: func() {
				pool = []string{"host1", "host2", "host3", "host4"}
				index = 3
				poolMap = map[string]bool{
					"host1": true,
					"host2": true,
					"host3": true,
					"host4": true,
				}
			},
			wantPool: []string{"host1", "host2", "host4"},
			wantMap: map[string]bool{
				"host1": true,
				"host2": true,
				"host3": false,
				"host4": true,
			},
			wantIdx: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				pool = []string{}
				poolMap = map[string]bool{}
				index = 0
			}()

			if err := r.RemoveHostByHostAddress(tt.args.ctx, tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("RemoveHostByHostAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(pool, tt.wantPool) {
				t.Errorf("RemoveHostByHostAddress() pool = %v, wantPool %v", pool, tt.wantPool)
			}
			if !reflect.DeepEqual(poolMap, tt.wantMap) {
				t.Errorf("RemoveHostByHostAddress() poolMap = %v, wantMap %v", poolMap, tt.wantMap)
			}
			if !reflect.DeepEqual(index, tt.wantIdx) {
				t.Errorf("RemoveHostByHostAddress() idx = %v, wantIdx %v", index, tt.wantIdx)
			}
		})
	}
}

func TestRepository_RequeueFirstHostToLast(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		mock     func()
		wantPool []string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			mock: func() {
				pool = []string{"host1", "host2", "host4"}
			},
			wantPool: []string{"host2", "host4", "host1"},
		},
		{
			name: "success with one host",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			mock: func() {
				pool = []string{"host1"}
			},
			wantPool: []string{"host1"},
		},
		{
			name: "success empty host",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			mock: func() {
				pool = []string{}
			},
			wantPool: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				pool = []string{}
				poolMap = map[string]bool{}
				index = 0
			}()

			if err := r.RequeueFirstHostToLast(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RequeueFirstHostToLast() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(pool, tt.wantPool) {
				t.Errorf("RequeueFirstHostToLast() pool = %v, wantPool %v", pool, tt.wantPool)
			}
		})
	}
}

func TestRepository_SetIndex(t *testing.T) {
	type args struct {
		ctx      context.Context
		newIndex int
	}
	tests := []struct {
		name    string
		args    args
		wantRes int
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				newIndex: 11,
			},
			wantRes: 2,
			wantErr: false,
			mock: func() {
				pool = append(pool, "host_1", "host_2", "host_3")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{}
			tt.mock()
			defer func() {
				pool = []string{}
				index = 0
			}()

			gotRes, err := r.SetIndex(tt.args.ctx, tt.args.newIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("SetIndex() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
