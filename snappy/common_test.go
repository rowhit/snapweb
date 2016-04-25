/*
 * Copyright (C) 2014-2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package snappy

import (
	"errors"
	"testing"

	"github.com/ubuntu-core/snappy/client"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type fakeSnapdClient struct {
	snaps  []*client.Snap
	err    error
	filter client.SnapFilter
}

func newDefaultSnap() *client.Snap {
	snap := &client.Snap{
		Description:   "WebRTC Video chat server for Snappy",
		DownloadSize:  6930947,
		Icon:          "/1.0/icons/chatroom.ogra/icon",
		InstalledSize: 18976651,
		Name:          "chatroom",
		Developer:     "ogra",
		Status:        client.StatusActive,
		Type:          client.TypeApp,
		Version:       "0.1-8",
	}
	return snap
}

func newSnap(name string) *client.Snap {
	snap := newDefaultSnap()
	snap.Name = name
	return snap
}

func (f *fakeSnapdClient) Icon(pkgID string) (*client.Icon, error) {
	icon := &client.Icon{
		Filename: "icon.png",
		Content:  []byte("png"),
	}
	return icon, nil
}

func (f *fakeSnapdClient) Services(pkg string) (map[string]*client.Service, error) {
	return nil, errors.New("the package has no services")
}

func (f *fakeSnapdClient) Snap(name string) (*client.Snap, *client.ResultInfo, error) {
	if len(f.snaps) > 0 {
		return f.snaps[0], nil, f.err
	}
	return nil, nil, f.err
}

func (f *fakeSnapdClient) FilterSnaps(filter client.SnapFilter) ([]*client.Snap, *client.ResultInfo, error) {
	f.filter = filter // record the filter used

	return f.snaps, nil, f.err
}

func (f *fakeSnapdClient) AddSnap(name string) (string, error) {
	return "", nil
}

func (f *fakeSnapdClient) RemoveSnap(name string) (string, error) {
	return "", nil
}

var _ snapdClient = (*fakeSnapdClient)(nil)

type fakeSnapdClientServicesNoExternalUI struct {
	fakeSnapdClient
}

func (f *fakeSnapdClientServicesNoExternalUI) Services(pkg string) (map[string]*client.Service, error) {
	internal := map[string]client.ServicePort{"ui": client.ServicePort{Port: "200/tcp"}}
	external := map[string]client.ServicePort{"web": client.ServicePort{Port: "1024/tcp"}}
	s1 := &client.Service{
		Spec: client.ServiceSpec{
			Ports: client.ServicePorts{
				Internal: internal,
				External: external,
			},
		},
	}

	s2 := &client.Service{}

	services := map[string]*client.Service{
		"s1": s1,
		"s2": s2,
	}

	return services, nil
}

type fakeSnapdClientServicesExternalUI struct {
	fakeSnapdClient
}

func (f *fakeSnapdClientServicesExternalUI) Services(pkg string) (map[string]*client.Service, error) {
	s1 := &client.Service{}

	internal := map[string]client.ServicePort{"ui": client.ServicePort{Port: "200/tcp"}}
	external := map[string]client.ServicePort{"ui": client.ServicePort{Port: "1024/tcp"}}
	s2 := &client.Service{
		Spec: client.ServiceSpec{
			Ports: client.ServicePorts{
				Internal: internal,
				External: external,
			},
		},
	}

	services := map[string]*client.Service{
		"s1": s1,
		"s2": s2,
	}

	return services, nil
}
