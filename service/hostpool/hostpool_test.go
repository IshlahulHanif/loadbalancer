package hostpool

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/loadbalancer/entity/host"
	hostEntity "github.com/loadbalancer/entity/host"
	"github.com/loadbalancer/pkg/config"
	"testing"
)

func TestService_HealthCheckAllHost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		config  config.Config
		usecase func() usecase
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp host.HealthCheckAllHostResp
		wantErr  bool
	}{
		{
			name: "success with failed to return unhealthy host to pool and failed to remove host from pool",
			fields: fields{
				config: config.Config{
					HostList: []string{"host1", "host2"},
				},
				usecase: func() usecase {
					mockUsecaseHostpool := NewMockUsecaseHostpool(ctrl)
					mockUsecasePoolClient := NewMockUsecasePoolClient(ctrl)

					mockUsecaseHostpool.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host3", "host2"}, nil)

					mockUsecasePoolClient.EXPECT().PingHost(gomock.Any(), "host1").Return(true)
					mockUsecasePoolClient.EXPECT().PingHost(gomock.Any(), "host2").Return(false)
					mockUsecasePoolClient.EXPECT().PingHost(gomock.Any(), "host3").Return(true)

					mockUsecaseHostpool.EXPECT().AddHost(gomock.Any(), "host1").Return(errors.New("err"))

					mockUsecaseHostpool.EXPECT().RemoveHost(gomock.Any(), "host2").Return(errors.New("err"))

					return usecase{
						hostpool:   mockUsecaseHostpool,
						poolclient: mockUsecasePoolClient,
					}
				},
			},
			wantResp: host.HealthCheckAllHostResp{
				HealthyHosts: []string{"host1", "host3"},
				DownHosts:    []string{"host2"},
			},
			wantErr: false,
		},
		{
			name: "err get host list from pool",
			fields: fields{
				config: config.Config{
					HostList: []string{"host1", "host2"},
				},
				usecase: func() usecase {
					mockUsecaseHostpool := NewMockUsecaseHostpool(ctrl)
					mockUsecasePoolClient := NewMockUsecasePoolClient(ctrl)

					mockUsecaseHostpool.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{}, errors.New("err"))

					return usecase{
						hostpool:   mockUsecaseHostpool,
						poolclient: mockUsecasePoolClient,
					}
				},
			},
			wantResp: host.HealthCheckAllHostResp{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				config:  tt.fields.config,
				usecase: tt.fields.usecase(),
			}
			gotResp, err := s.HealthCheckAllHost(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheckAllHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var gotRespHealthyMap, gotRespDownMap = make(map[string]bool), make(map[string]bool)
			for _, downHost := range gotResp.DownHosts {
				gotRespDownMap[downHost] = true
			}
			for _, healthyHost := range gotResp.HealthyHosts {
				gotRespHealthyMap[healthyHost] = true
			}

			for _, downHost := range tt.wantResp.DownHosts {
				if !gotRespDownMap[downHost] {
					t.Errorf("HealthCheckAllHost() gotResp = %v, want %v", gotResp, tt.wantResp)
				}
			}
			for _, healthyHost := range tt.wantResp.HealthyHosts {
				if !gotRespHealthyMap[healthyHost] {
					t.Errorf("HealthCheckAllHost() gotResp = %v, want %v", gotResp, tt.wantResp)
				}
			}
		})
	}
}

func TestService_ManageHost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		config  config.Config
		usecase func() usecase
	}
	type args struct {
		ctx context.Context
		req host.ManageHostReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success add",
			fields: fields{
				usecase: func() usecase {
					mockUsecaseHostpool := NewMockUsecaseHostpool(ctrl)
					mockUsecaseHostpool.EXPECT().AddHost(gomock.Any(), "host1").Return(nil)
					mockUsecaseHostpool.EXPECT().AddHost(gomock.Any(), "host2").Return(nil)

					return usecase{
						hostpool: mockUsecaseHostpool,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				req: host.ManageHostReq{
					Operation: hostEntity.Add,
					Data: []string{
						"host1",
						"host2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "err add",
			fields: fields{
				usecase: func() usecase {
					mockUsecaseHostpool := NewMockUsecaseHostpool(ctrl)
					mockUsecaseHostpool.EXPECT().AddHost(gomock.Any(), "host1").Return(nil)
					mockUsecaseHostpool.EXPECT().AddHost(gomock.Any(), "host2").Return(errors.New("err"))

					return usecase{
						hostpool: mockUsecaseHostpool,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				req: host.ManageHostReq{
					Operation: hostEntity.Add,
					Data: []string{
						"host1",
						"host2",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "success remove",
			fields: fields{
				usecase: func() usecase {
					mockUsecaseHostpool := NewMockUsecaseHostpool(ctrl)
					mockUsecaseHostpool.EXPECT().RemoveHost(gomock.Any(), "host1").Return(nil)
					mockUsecaseHostpool.EXPECT().RemoveHost(gomock.Any(), "host2").Return(nil)

					return usecase{
						hostpool: mockUsecaseHostpool,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				req: host.ManageHostReq{
					Operation: hostEntity.Remove,
					Data: []string{
						"host1",
						"host2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "err remove",
			fields: fields{
				usecase: func() usecase {
					mockUsecaseHostpool := NewMockUsecaseHostpool(ctrl)
					mockUsecaseHostpool.EXPECT().RemoveHost(gomock.Any(), "host1").Return(nil)
					mockUsecaseHostpool.EXPECT().RemoveHost(gomock.Any(), "host2").Return(errors.New("err"))

					return usecase{
						hostpool: mockUsecaseHostpool,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				req: host.ManageHostReq{
					Operation: hostEntity.Remove,
					Data: []string{
						"host1",
						"host2",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "err default",
			fields: fields{
				usecase: func() usecase {
					return usecase{}
				},
			},
			args: args{
				ctx: context.Background(),
				req: host.ManageHostReq{
					Operation: hostEntity.Unknown,
					Data:      []string{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				config:  tt.fields.config,
				usecase: tt.fields.usecase(),
			}
			if err := s.ManageHost(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ManageHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
