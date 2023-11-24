package rke

import "testing"

func Test_k8sVersionRequiresCri(t *testing.T) {
	type args struct {
		kubernetesVersion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
        {
            name: "v1.26.9-rancher1-1",
            args: args{
                kubernetesVersion: "v1.26.9-rancher1-1",
            },
            want: true,
        },
		{
			name: "v1.26.4-rancher2-1",
			args: args{
				kubernetesVersion: "v1.26.4-rancher2-1",
			},
			want: true,
		},
		{
			name: "v1.25.9-rancher2-2",
			args: args{
				kubernetesVersion: "v1.25.9-rancher2-2",
			},
			want: true,
		},
		{
			name: "v1.24.13-rancher2-2",
			args: args{
				kubernetesVersion: "v1.24.13-rancher2-2",
			},
			want: true,
		},
		{
			name: "v1.23.16-rancher2-3",
			args: args{
				kubernetesVersion: "v1.23.16-rancher2-3",
			},
			want: false,
		},
		{
			name: "v1.22.17-rancher1-2",
			args: args{
				kubernetesVersion: "v1.22.17-rancher1-2",
			},
			want: false,
		},
		{
			name: "v1.21.14-rancher1-1",
			args: args{
				kubernetesVersion: "v1.21.14-rancher1-1",
			},
			want: false,
		},
		{
			name: "v1.20.15-rancher2-2",
			args: args{
				kubernetesVersion: "v1.20.15-rancher2-2",
			},
			want: false,
		},
		{
			name: "invalid",
			args: args{
				kubernetesVersion: "invalid",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := k8sVersionRequiresCri(tt.args.kubernetesVersion); got != tt.want {
				t.Errorf("k8sVersionRequiresCri() = %v, want %v", got, tt.want)
			}
		})
	}
}
