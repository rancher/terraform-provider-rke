package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterServicesETCDBackupS3Conf      *rancher.S3BackupConfig
	testRKEClusterServicesETCDBackupS3Interface []interface{}
	testRKEClusterServicesETCDBackupConf        *rancher.BackupConfig
	testRKEClusterServicesETCDBackupInterface   []interface{}
	testRKEClusterServicesETCDConf              rancher.ETCDService
	testRKEClusterServicesETCDInterface         []interface{}
)

func init() {
	testRKEClusterServicesETCDBackupS3Conf = &rancher.S3BackupConfig{
		AccessKey:  "access_key",
		BucketName: "bucket_name",
		CustomCA:   "custom_ca",
		Endpoint:   "endpoint",
		Folder:     "folder",
		Region:     "region",
	}
	testRKEClusterServicesETCDBackupS3Interface = []interface{}{
		map[string]interface{}{
			"access_key":  "access_key",
			"bucket_name": "bucket_name",
			"custom_ca":   base64Encode("custom_ca"),
			"endpoint":    "endpoint",
			"folder":      "folder",
			"region":      "region",
		},
	}
	testRKEClusterServicesETCDBackupConf = &rancher.BackupConfig{
		Enabled:        newTrue(),
		IntervalHours:  20,
		Retention:      10,
		S3BackupConfig: testRKEClusterServicesETCDBackupS3Conf,
		SafeTimestamp:  true,
		Timeout:        500,
	}
	testRKEClusterServicesETCDBackupInterface = []interface{}{
		map[string]interface{}{
			"enabled":          true,
			"interval_hours":   20,
			"retention":        10,
			"s3_backup_config": testRKEClusterServicesETCDBackupS3Interface,
			"safe_timestamp":   true,
			"timeout":          500,
		},
	}
	testRKEClusterServicesETCDConf = rancher.ETCDService{
		BackupConfig: testRKEClusterServicesETCDBackupConf,
		CACert:       "XXXXXXXX",
		Cert:         "YYYYYYYY",
		Creation:     "creation",
		ExternalURLs: []string{"url_one", "url_two"},
		GID:          1001,
		Key:          "ZZZZZZZZ",
		Path:         "/etcd",
		Retention:    "6h",
		Snapshot:     newTrue(),
		UID:          1001,
	}
	testRKEClusterServicesETCDConf.ExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesETCDConf.WindowsExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesETCDConf.ExtraArgsArray = map[string][]string{
		"arg1": {"v1"},
		"arg2": {"v2"},
	}
	testRKEClusterServicesETCDConf.WindowsExtraArgsArray = map[string][]string{
		"arg1": {"v1"},
		"arg2": {"v2"},
	}
	testRKEClusterServicesETCDConf.ExtraBinds = []string{"bind_one", "bind_two"}
	testRKEClusterServicesETCDConf.ExtraEnv = []string{"env_one", "env_two"}
	testRKEClusterServicesETCDConf.Image = "image"
	testRKEClusterServicesETCDInterface = []interface{}{
		map[string]interface{}{
			"backup_config": testRKEClusterServicesETCDBackupInterface,
			"ca_cert":       "XXXXXXXX",
			"cert":          "YYYYYYYY",
			"creation":      "creation",
			"external_urls": []interface{}{"url_one", "url_two"},
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"win_extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_args_array": []interface{}{
				map[string]interface{}{
					"extra_arg": []interface{}{
						map[string]interface{}{
							"argument": "arg1",
							"values":   []interface{}{"v1"},
						},
						map[string]interface{}{
							"argument": "arg2",
							"values":   []interface{}{"v2"},
						},
					},
				},
			},
			"win_extra_args_array": []interface{}{
				map[string]interface{}{
					"extra_arg": []interface{}{
						map[string]interface{}{
							"argument": "arg1",
							"values":   []interface{}{"v1"},
						},
						map[string]interface{}{
							"argument": "arg2",
							"values":   []interface{}{"v2"},
						},
					},
				},
			},
			"extra_binds": []interface{}{"bind_one", "bind_two"},
			"extra_env":   []interface{}{"env_one", "env_two"},
			"gid":         1001,
			"image":       "image",
			"key":         "ZZZZZZZZ",
			"path":        "/etcd",
			"retention":   "6h",
			"snapshot":    true,
			"uid":         1001,
		},
	}
}

func TestFlattenRKEClusterServicesEtcdBackupConfigS3(t *testing.T) {

	cases := []struct {
		Input          *rancher.S3BackupConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesETCDBackupS3Conf,
			testRKEClusterServicesETCDBackupS3Interface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterServicesEtcdBackupConfigS3(tc.Input, testRKEClusterServicesETCDBackupS3Interface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterServicesEtcdBackupConfig(t *testing.T) {

	cases := []struct {
		Input          *rancher.BackupConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesETCDBackupConf,
			testRKEClusterServicesETCDBackupInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterServicesEtcdBackupConfig(tc.Input, testRKEClusterServicesETCDBackupInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterServicesEtcd(t *testing.T) {

	cases := []struct {
		Input          rancher.ETCDService
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesETCDConf,
			testRKEClusterServicesETCDInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterServicesEtcd(tc.Input, testRKEClusterServicesETCDInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesEtcdBackupConfigS3(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.S3BackupConfig
	}{
		{
			testRKEClusterServicesETCDBackupS3Interface,
			testRKEClusterServicesETCDBackupS3Conf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesEtcdBackupConfigS3(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesEtcdBackupConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.BackupConfig
	}{
		{
			testRKEClusterServicesETCDBackupInterface,
			testRKEClusterServicesETCDBackupConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesEtcdBackupConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesEtcd(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.ETCDService
	}{
		{
			testRKEClusterServicesETCDInterface,
			testRKEClusterServicesETCDConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesEtcd(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
