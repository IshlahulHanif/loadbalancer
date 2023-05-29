package hostpool

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestUsecase_AddHost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repo func() repository
	}
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success normal",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().AppendHost(gomock.Any(), "host_x").Return(nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx:  context.Background(),
				host: "host_x",
			},
			wantErr: false,
		},
		{
			name: "err append hots",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().AppendHost(gomock.Any(), "host_x").Return(errors.New("test err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx:  context.Background(),
				host: "host_x",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Usecase{
				repo: tt.fields.repo(),
			}
			if err := u.AddHost(tt.args.ctx, tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("AddHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_GetHostDequeueRoundRobin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repo func() repository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes string
		wantErr bool
	}{
		{
			name: "success normal",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host_1", "host_2", "host_3"}, nil)
					hostPoolRepoMock.EXPECT().GetCurrentIndex(gomock.Any()).Return(1, nil)
					hostPoolRepoMock.EXPECT().IncrementIndex(gomock.Any(), 1).Return(2, nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "host_2",
			wantErr: false,
		},
		{
			name: "success index overload",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host_1", "host_2", "host_3"}, nil)
					hostPoolRepoMock.EXPECT().GetCurrentIndex(gomock.Any()).Return(3, nil)
					hostPoolRepoMock.EXPECT().SetIndex(gomock.Any(), 0).Return(0, nil)
					hostPoolRepoMock.EXPECT().RequeueFirstHostToLast(gomock.Any()).Return(nil)
					hostPoolRepoMock.EXPECT().IncrementIndex(gomock.Any(), 1).Return(2, nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "host_1",
			wantErr: false,
		},
		{
			name: "err incr idx",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host_1", "host_2", "host_3"}, nil)
					hostPoolRepoMock.EXPECT().GetCurrentIndex(gomock.Any()).Return(3, nil)
					hostPoolRepoMock.EXPECT().SetIndex(gomock.Any(), 0).Return(0, nil)
					hostPoolRepoMock.EXPECT().RequeueFirstHostToLast(gomock.Any()).Return(nil)
					hostPoolRepoMock.EXPECT().IncrementIndex(gomock.Any(), 1).Return(2, errors.New("test err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "",
			wantErr: true,
		},
		{
			name: "err idx out of bound",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host_1", "host_2", "host_3"}, nil)
					hostPoolRepoMock.EXPECT().GetCurrentIndex(gomock.Any()).Return(3, nil)
					hostPoolRepoMock.EXPECT().SetIndex(gomock.Any(), 0).Return(10, nil)
					hostPoolRepoMock.EXPECT().RequeueFirstHostToLast(gomock.Any()).Return(nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "",
			wantErr: true,
		},
		{
			name: "err rearrange queue",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host_1", "host_2", "host_3"}, nil)
					hostPoolRepoMock.EXPECT().GetCurrentIndex(gomock.Any()).Return(3, nil)
					hostPoolRepoMock.EXPECT().SetIndex(gomock.Any(), 0).Return(10, nil)
					hostPoolRepoMock.EXPECT().RequeueFirstHostToLast(gomock.Any()).Return(errors.New("test err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "",
			wantErr: true,
		},
		{
			name: "err set idx",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host_1", "host_2", "host_3"}, nil)
					hostPoolRepoMock.EXPECT().GetCurrentIndex(gomock.Any()).Return(3, nil)
					hostPoolRepoMock.EXPECT().SetIndex(gomock.Any(), 0).Return(0, errors.New("test err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "",
			wantErr: true,
		},
		{
			name: "err get curr idx",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host_1", "host_2", "host_3"}, nil)
					hostPoolRepoMock.EXPECT().GetCurrentIndex(gomock.Any()).Return(3, errors.New("test err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "",
			wantErr: true,
		},
		{
			name: "err empty host pool",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{}, nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "",
			wantErr: true,
		},
		{
			name: "err get host list",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{}, errors.New("test err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Usecase{
				repo: tt.fields.repo(),
			}
			gotRes, err := u.GetHostDequeueRoundRobin(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostDequeueRoundRobin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GetHostDequeueRoundRobin() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestUsecase_GetHostListFromPool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repo func() repository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []string
		wantErr bool
	}{
		{
			name: "success normal",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{"host"}, nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: []string{"host"},
			wantErr: false,
		},
		{
			name: "success empty",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{}, nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: []string{},
			wantErr: false,
		},
		{
			name: "err get hots",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().GetHostListFromPool(gomock.Any()).Return([]string{}, errors.New("err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantRes: []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Usecase{
				repo: tt.fields.repo(),
			}
			gotRes, err := u.GetHostListFromPool(tt.args.ctx)
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

func TestUsecase_RemoveHost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repo func() repository
	}
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success normal",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().RemoveHostByHostAddress(gomock.Any(), "host_x").Return(nil)

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx:  context.Background(),
				host: "host_x",
			},
			wantErr: false,
		},
		{
			name: "err RemoveHost host",
			fields: fields{
				repo: func() repository {
					hostPoolRepoMock := NewMockRepoHostpool(ctrl)
					hostPoolRepoMock.EXPECT().RemoveHostByHostAddress(gomock.Any(), "host_x").Return(errors.New("test err"))

					return repository{
						hostpool: hostPoolRepoMock,
					}
				},
			},
			args: args{
				ctx:  context.Background(),
				host: "host_x",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Usecase{
				repo: tt.fields.repo(),
			}
			if err := u.RemoveHost(tt.args.ctx, tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("RemoveHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
