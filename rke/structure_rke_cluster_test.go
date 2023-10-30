package rke

import "testing"

func Test_k8sVersionRequiresCri(t *testing.T) {
	type args struct {
		kuberenetsVersion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "v1.26.4-rancher2-1",
			args: args{
				kuberenetsVersion: "v1.26.4-rancher2-1",
			},
			want: true,
		},
		{
			name: "v1.25.9-rancher2-2",
			args: args{
				kuberenetsVersion: "v1.25.9-rancher2-2",
			},
			want: true,
		},
		{
			name: "v1.25.9-rancher2-1",
			args: args{
				kuberenetsVersion: "v1.25.9-rancher2-1",
			},
			want: true,
		},
		{
			name: "v1.25.6-rancher4-1",
			args: args{
				kuberenetsVersion: "v1.25.6-rancher4-1",
			},
			want: true,
		},
		{
			name: "v1.25.6-rancher2-1",
			args: args{
				kuberenetsVersion: "v1.25.6-rancher2-1",
			},
			want: true,
		},
		{
			name: "v1.24.13-rancher2-2",
			args: args{
				kuberenetsVersion: "v1.24.13-rancher2-2",
			},
			want: true,
		},
		{
			name: "v1.24.13-rancher2-1",
			args: args{
				kuberenetsVersion: "v1.24.13-rancher2-1",
			},
			want: true,
		},
		{
			name: "v1.24.10-rancher4-1",
			args: args{
				kuberenetsVersion: "v1.24.10-rancher4-1",
			},
			want: true,
		},
		{
			name: "v1.24.10-rancher2-1",
			args: args{
				kuberenetsVersion: "v1.24.10-rancher2-1",
			},
			want: true,
		},
		{
			name: "v1.24.9-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.24.9-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.24.8-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.24.8-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.24.6-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.24.6-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.24.4-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.24.4-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.24.2-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.24.2-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.16-rancher2-3",
			args: args{
				kuberenetsVersion: "v1.23.16-rancher2-3",
			},
			want: true,
		},
		{
			name: "v1.23.16-rancher2-1",
			args: args{
				kuberenetsVersion: "v1.23.16-rancher2-1",
			},
			want: true,
		},
		{
			name: "v1.23.15-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.15-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.14-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.14-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.12-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.12-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.10-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.10-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.8-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.8-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.7-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.7-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.6-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.6-rancher1-1",
			},
			want: true,
		},
		{
			name: "v1.23.4-rancher1-2",
			args: args{
				kuberenetsVersion: "v1.23.4-rancher1-2",
			},
			want: true,
		},
		{
			name: "v1.23.4-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.23.4-rancher1-2",
			},
			want: true,
		},
		{
			name: "v1.22.17-rancher1-2",
			args: args{
				kuberenetsVersion: "v1.22.17-rancher1-2",
			},
			want: false,
		},
		{
			name: "v1.22.17-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.17-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.16-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.16-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.15-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.15-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.13-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.13-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.11-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.11-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.10-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.10-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.9-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.9-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.7-rancher1-2",
			args: args{
				kuberenetsVersion: "v1.22.7-rancher1-2",
			},
			want: false,
		},
		{
			name: "v1.22.7-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.7-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.6-rancher1-1",
			args: args{
				kuberenetsVersion: "v1.22.6-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.22.5-rancher2-1",
			args: args{
				kuberenetsVersion: "v1.22.5-rancher2-1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := k8sVersionRequiresCri(tt.args.kuberenetsVersion); got != tt.want {
				t.Errorf("k8sVersionRequiresCri() = %v, want %v", got, tt.want)
			}
		})
	}
}