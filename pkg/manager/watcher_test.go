package manager

import (
	"reflect"
	"testing"

	"github.com/kictto/kube-vip/pkg/bgp"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestParseBgpAnnotations(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Annotations: map[string]string{}},
	}

	_, _, err := parseBgpAnnotations(node, "bgp")
	if err == nil {
		t.Fatal("Parsing BGP annotations should return an error when no annotations exist")
	}

	node.Annotations = map[string]string{
		"bgp/node-asn": "65000",
		"bgp/peer-asn": "64000",
		"bgp/src-ip":   "10.0.0.254",
	}

	bgpConfig, bgpPeer, err := parseBgpAnnotations(node, "bgp")
	if err != nil {
		t.Fatal("Parsing BGP annotations should return nil when minimum config is met")
	}

	assert.Equal(t, uint32(65000), bgpConfig.AS, "bgpConfig.AS parsed incorrectly")
	assert.Equal(t, uint32(64000), bgpPeer.AS, "bgpPeer.AS parsed incorrectly")
	assert.Equal(t, "10.0.0.254", bgpConfig.RouterID, "bgpConfig.RouterID parsed incorrectly")

	node.Annotations = map[string]string{
		"bgp/node-asn": "65000",
		"bgp/peer-asn": "64000",
		"bgp/src-ip":   "10.0.0.254",
		"bgp/peer-ip":  "10.0.0.1,10.0.0.2,10.0.0.3",
		"bgp/bgp-pass": "cGFzc3dvcmQ=", // password
	}

	bgpConfig, bgpPeer, err = parseBgpAnnotations(node, "bgp")
	if err != nil {
		t.Fatal("Parsing BGP annotations should return nil when minimum config is met")
	}

	bgpPeers := []bgp.Peer{
		{Address: "10.0.0.1", AS: uint32(64000), Password: "password"},
		{Address: "10.0.0.2", AS: uint32(64000), Password: "password"},
		{Address: "10.0.0.3", AS: uint32(64000), Password: "password"},
	}
	assert.Equal(t, bgpPeers, bgpConfig.Peers, "bgpConfig.Peers parsed incorrectly")
	assert.Equal(t, "10.0.0.3", bgpPeer.Address, "bgpPeer.Address parsed incorrectly")
	assert.Equal(t, "password", bgpPeer.Password, "bgpPeer.Password parsed incorrectly")
}

func Test_parseBgpAnnotations(t *testing.T) {
	type args struct {
		node   *v1.Node
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    bgp.Config
		want1   bgp.Peer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseBgpAnnotations(tt.args.node, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBgpAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBgpAnnotations() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseBgpAnnotations() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
