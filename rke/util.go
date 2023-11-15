package rke

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/blang/semver"
	ghodssyaml "github.com/ghodss/yaml"
	gover "github.com/hashicorp/go-version"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

var parsedRangeAtLeast124 = semver.MustParseRange(">= 1.24.0-rancher0")

func splitImportID(s string) ([]string, error) {
	sep := ":"
	if len(s) == 0 {
		return nil, fmt.Errorf("Import ID is nil")
	}

	result := strings.Split(s, sep)
	if len(result) != 2 && len(result) != 3 {
		return nil, fmt.Errorf("Import ID bad format")
	}
	return result, nil
}

func base64Encode(s string) string {
	if len(s) == 0 {
		return ""
	}
	data := []byte(s)

	return base64.StdEncoding.EncodeToString(data)
}

func base64Decode(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(s)

	return string(data), err
}

func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func toArrayString(in []interface{}) []string {
	out := make([]string, len(in))
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toArrayInterface(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func toMapString(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toMapInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for i, v := range in {
		out[i] = v
	}
	return out
}

func jsonToMapInterface(in string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := json.Unmarshal([]byte(in), &out)
	if err != nil {
		return nil, err
	}
	return out, err
}

func mapInterfaceToJSON(in map[string]interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	out, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func jsonToInterface(in string, out interface{}) error {
	if out == nil {
		return nil
	}
	err := json.Unmarshal([]byte(in), out)
	if err != nil {
		return err
	}
	return err
}

func interfaceToMap(in interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := make(map[string]interface{})

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func interfaceToJSON(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	out, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func yamlToMapInterface(in string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(in), &out)
	if err != nil {
		return nil, err
	}
	return out, err
}

func yamlToInterface(in string, out interface{}) error {
	if out == nil {
		return nil
	}
	err := yaml.Unmarshal([]byte(in), out)
	if err != nil {
		return err
	}
	return err
}

func ghodssyamlToMapInterface(in string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := ghodssyaml.Unmarshal([]byte(in), &out)
	if err != nil {
		return nil, err
	}
	return out, err
}

func ghodssyamlToInterface(in string, out interface{}) error {
	if out == nil {
		return nil
	}
	err := ghodssyaml.Unmarshal([]byte(in), out)
	if err != nil {
		return err
	}
	return err
}

func interfaceToYaml(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	out, err := yaml.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func interfaceToGhodssyaml(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	out, err := ghodssyaml.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func fileExist(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}

func sortVersions(list map[string]string) ([]*gover.Version, error) {
	var versions []*gover.Version
	for key := range list {
		v, err := gover.NewVersion(key)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	sort.Sort(gover.Collection(versions))
	return versions, nil
}

func getClusterVersion(version string) (semver.Version, error) {
	var parsedVersion semver.Version
	if len(version) <= 1 || !strings.HasPrefix(version, "v") {
		return parsedVersion, fmt.Errorf("%s is not valid version", version)
	}
	parsedVersion, err := semver.Parse(version[1:])
	if err != nil {
		return parsedVersion, fmt.Errorf("%s is not valid semver", version)
	}
	return parsedVersion, nil
}

func getLatestVersion(list map[string]string) (string, error) {
	sorted, err := sortVersions(list)
	if err != nil {
		return "", err
	}

	return sorted[len(sorted)-1].String(), nil
}

func getNewUUID() string {
	newuid, _ := uuid.NewV4()
	return newuid.String()
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func privateKeyToPEM(key *rsa.PrivateKey) string {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	return string(pemdata)
}

func certificateToPEM(cert *x509.Certificate) string {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
	return string(pemdata)
}
