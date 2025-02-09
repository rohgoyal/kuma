package common_test

import (
	mesh_proto "github.com/Kong/kuma/api/mesh/v1alpha1"
	core_mesh "github.com/Kong/kuma/pkg/core/resources/apis/mesh"
	sds_auth "github.com/Kong/kuma/pkg/sds/auth"
	auth_common "github.com/Kong/kuma/pkg/sds/auth/common"
	"github.com/Kong/kuma/pkg/test/resources/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDataplaneIdentity()", func() {
	It("should get identity from Dataplane with inbounds", func() {
		// given
		dpRes := core_mesh.DataplaneResource{
			Meta: &model.ResourceMeta{
				Mesh: "demo",
				Name: "dp-1",
			},
			Spec: mesh_proto.Dataplane{
				Networking: &mesh_proto.Dataplane_Networking{
					Inbound: []*mesh_proto.Dataplane_Networking_Inbound{
						{
							Interface: "127.0.0.1:8080:80",
							Tags: map[string]string{
								"service": "backend",
							},
						},
					},
				},
			},
		}

		// when
		id, err := auth_common.GetDataplaneIdentity(&dpRes)

		// then
		Expect(err).ToNot(HaveOccurred())

		// and
		Expect(id).To(Equal(sds_auth.Identity{
			Mesh:    "demo",
			Service: "backend",
		}))
	})

	It("should get identity from Dataplane with gateway", func() {
		// given
		dpRes := core_mesh.DataplaneResource{
			Meta: &model.ResourceMeta{
				Mesh: "demo",
				Name: "dp-1",
			},
			Spec: mesh_proto.Dataplane{
				Networking: &mesh_proto.Dataplane_Networking{
					Gateway: &mesh_proto.Dataplane_Networking_Gateway{
						Tags: map[string]string{
							"service": "edge",
						},
					},
				},
			},
		}

		// when
		id, err := auth_common.GetDataplaneIdentity(&dpRes)

		// then
		Expect(err).ToNot(HaveOccurred())

		// and
		Expect(id).To(Equal(sds_auth.Identity{
			Mesh:    "demo",
			Service: "edge",
		}))
	})

	It("should throw an error on dataplane without services", func() {
		// given
		dpRes := core_mesh.DataplaneResource{
			Meta: &model.ResourceMeta{
				Mesh: "demo",
				Name: "dp-1",
			},
			Spec: mesh_proto.Dataplane{},
		}

		// when
		_, err := auth_common.GetDataplaneIdentity(&dpRes)

		// then
		Expect(err).To(MatchError("Dataplane has no services associated with it"))
	})
})
