package discovery

import (
	"reflect"
	"strings"
	"testing"
)

type staticHostEntry struct {
	in  string
	out []string
}

func TestLookupStaticHost(t *testing.T) {
	tests := []struct {
		name string
		ents []staticHostEntry
	}{
		{
			"testdata/hosts",
			[]staticHostEntry{
				{"odin", []string{"127.0.0.2", "127.0.0.3", "::2"}},
				{"thor", []string{"127.1.1.1"}},
				{"ullr", []string{"127.1.1.2"}},
				{"ullrhost", []string{"127.1.1.2"}},
				{"localhost", []string{"fe80::1%lo0"}},
			},
		},
		{
			"testdata/singleline-hosts", // see golang.org/issue/6646
			[]staticHostEntry{
				{"odin", []string{"127.0.0.2"}},
			},
		},
		{
			"testdata/ipv4-hosts", // see golang.org/issue/8996
			[]staticHostEntry{
				{"localhost", []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}},
				{"localhost.localdomain", []string{"127.0.0.3"}},
			},
		},
		{
			"testdata/ipv6-hosts", // see golang.org/issue/8996
			[]staticHostEntry{
				{"localhost", []string{"::1", "fe80::1", "fe80::2%lo0", "fe80::3%lo0"}},
				{"localhost.localdomain", []string{"fe80::3%lo0"}},
			},
		},
		{
			"testdata/case-hosts", // see golang.org/issue/12806
			[]staticHostEntry{
				{"PreserveMe", []string{"127.0.0.1", "::1"}},
				{"PreserveMe.local", []string{"127.0.0.1", "::1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			names, _, err := readHostsFile(tt.name)
			if err != nil {
				t.Fatalf("readHostsFile() err = %v", err)
			}
			for _, ent := range tt.ents {
				ins := []string{ent.in, absDomainName([]byte(ent.in)), strings.ToLower(ent.in), strings.ToUpper(ent.in)}
				for _, in := range ins {
					addrs := names[prepareHostLookup(in)]
					if !reflect.DeepEqual(addrs, ent.out) {
						t.Errorf("lookupStaticHost(%s)\ngot:  %v\nwant: %v", in, addrs, ent.out)
					}
				}
			}
		})
	}
}

func TestLookupStaticAddr(t *testing.T) {
	tests := []struct {
		name string
		ents []staticHostEntry
	}{
		{
			"testdata/hosts",
			[]staticHostEntry{
				{"255.255.255.255", []string{"broadcasthost"}},
				{"127.0.0.2", []string{"odin"}},
				{"127.0.0.3", []string{"odin"}},
				{"::2", []string{"odin"}},
				{"127.1.1.1", []string{"thor"}},
				{"127.1.1.2", []string{"ullr", "ullrhost"}},
				{"fe80::1%lo0", []string{"localhost"}},
			},
		},
		{
			"testdata/singleline-hosts", // see golang.org/issue/6646
			[]staticHostEntry{
				{"127.0.0.2", []string{"odin"}},
			},
		},
		{
			"testdata/ipv4-hosts", // see golang.org/issue/8996
			[]staticHostEntry{
				{"127.0.0.1", []string{"localhost"}},
				{"127.0.0.2", []string{"localhost"}},
				{"127.0.0.3", []string{"localhost", "localhost.localdomain"}},
			},
		},
		{
			"testdata/ipv6-hosts", // see golang.org/issue/8996
			[]staticHostEntry{
				{"::1", []string{"localhost"}},
				{"fe80::1", []string{"localhost"}},
				{"fe80::2%lo0", []string{"localhost"}},
				{"fe80::3%lo0", []string{"localhost", "localhost.localdomain"}},
			},
		},
		{
			"testdata/case-hosts", // see golang.org/issue/12806
			[]staticHostEntry{
				{"127.0.0.1", []string{"PreserveMe", "PreserveMe.local"}},
				{"::1", []string{"PreserveMe", "PreserveMe.local"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, addrs, err := readHostsFile(tt.name)
			if err != nil {
				t.Fatalf("readHostsFile() err = %v", err)
			}
			for _, ent := range tt.ents {
				hosts := addrs[ent.in]
				for i := range ent.out {
					ent.out[i] = absDomainName([]byte(ent.out[i]))
				}
				if !reflect.DeepEqual(hosts, ent.out) {
					t.Errorf("lookupStaticAddr(%s)\ngot:  %v\nwant: %v", ent.in, hosts, ent.out)
				}
			}
		})
	}
}
