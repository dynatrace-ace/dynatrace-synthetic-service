package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// Test dynatrace credential behavior: values should be read from dynatrace secret unless secret name has been overridden in dynatrace config file.
// If neither is available, an error should be produced.
func TestDynatraceK8CredentialsReader_GetDynatraceCredentials(t *testing.T) {

	wantDynatraceCredentials, err := NewDynatraceCredentials("https://mySampleEnv.live.dynatrace.com", testDynatraceAPIToken)
	assert.NoError(t, err)
	wantDynatraceHTTPCredentials, err := NewDynatraceCredentials("http://mySampleEnv.live.dynatrace.com", testDynatraceAPIToken)
	assert.NoError(t, err)

	type args struct {
		secretName string
	}
	tests := []struct {
		name    string
		secret  *v1.Secret
		args    args
		want    *DynatraceCredentials
		wantErr bool
	}{
		{
			name:   "with no secret",
			secret: &v1.Secret{},
			args: args{
				secretName: "dynatrace",
			},
			wantErr: true,
		},
		{
			name: "with dynatrace secret",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT":    "https://mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			args: args{
				secretName: "dynatrace",
			},
			want:    wantDynatraceCredentials,
			wantErr: false,
		},
		{
			name: "with dynatrace secret - want HTTPS",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT":    "mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			args: args{
				secretName: "dynatrace",
			},
			want:    wantDynatraceCredentials,
			wantErr: false,
		},
		{
			name: "with dynatrace secret - want HTTP",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT":    "http://mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			args: args{
				secretName: "dynatrace",
			},
			want:    wantDynatraceHTTPCredentials,
			wantErr: false,
		},
		{
			name: "with dynatrace secret - invalid URL",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT":    "//mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			args: args{
				secretName: "dynatrace",
			},

			wantErr: true,
		},
		{
			name: "with dynatrace secret - invalid scheme",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT":    "ftp://mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			args: args{
				secretName: "dynatrace",
			},

			wantErr: true,
		},
		{
			name: "with dynatrace secret - invalid token",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT":    "ftp://mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": "dcO.public.private",
				}),
			args: args{
				secretName: "dynatrace",
			},

			wantErr: true,
		},
		{
			name: "with dynatrace secret - no token",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT": "ftp://mySampleEnv.live.dynatrace.com",
				}),
			args: args{
				secretName: "dynatrace",
			},

			wantErr: true,
		},
		{
			name: "with dynatrace_other secret, with other secret name",
			secret: createTestSecret(
				"dynatrace_other",
				map[string]string{
					"DT_TENANT":    "https://mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			args: args{
				secretName: "dynatrace_other",
			},
			want:    wantDynatraceCredentials,
			wantErr: false,
		},
		{
			name: "with dynatrace_other secret, with invalid secret name",
			secret: createTestSecret(
				"dynatrace_other",
				map[string]string{
					"DT_TENANT":    "https://mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			args: args{
				secretName: "dynatrace_other2",
			},
			wantErr: true,
		},
		{
			name: "with dynatrace secret, empty secret name",
			secret: createTestSecret(
				"dynatrace",
				map[string]string{
					"DT_TENANT":    "https://mySampleEnv.live.dynatrace.com",
					"DT_API_TOKEN": testDynatraceAPIToken,
				}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			secretReader := NewK8sSecretReader(fake.NewSimpleClientset(tt.secret))
			credentialsProvider := NewDynatraceK8sSecretReader(secretReader)
			got, err := credentialsProvider.GetDynatraceCredentials(context.Background(), tt.args.secretName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
