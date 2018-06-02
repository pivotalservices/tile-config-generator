package generator

type NetworkProperties struct {
	Network                   *Name  `yaml:"network,omitempty"`
	ServiceNetwork            *Name  `yaml:"service_network,omitempty"`
	OtherAvailabilityZones    []Name `yaml:"other_availability_zones"`
	SingletonAvailabilityZone *Name  `yaml:"singleton_availability_zone"`
}

type Name struct {
	Name string `yaml:"name"`
}

//go:generate counterfeiter -o ./fakes/metadata.go --fake-name FakeMetadata . metadata
type metadata interface {
	UsesServiceNetwork() bool
}

func NewNetworkProperties(metadata metadata) *NetworkProperties {
	props := &NetworkProperties{}
	props.Network = &Name{
		Name: "((network_name))",
	}
	if metadata.UsesServiceNetwork() {
		props.ServiceNetwork = &Name{
			Name: "((service_network_name))",
		}
	}
	props.SingletonAvailabilityZone = &Name{
		Name: "((singleton_availability_zone))",
	}
	props.OtherAvailabilityZones = append(props.OtherAvailabilityZones, Name{
		Name: "((singleton_availability_zone))",
	})
	return props
}
