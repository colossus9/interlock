package service

import (
	"encoding/csv"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/opts"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	"github.com/docker/engine-api/types/swarm"
	"github.com/docker/go-connections/nat"
	units "github.com/docker/go-units"
	"github.com/spf13/cobra"
)

var (
	// DefaultReplicas is the default replicas to use for a replicated service
	DefaultReplicas uint64 = 1
)

type int64Value interface {
	Value() int64
}

type memBytes int64

func (m *memBytes) String() string {
	return units.BytesSize(float64(m.Value()))
}

func (m *memBytes) Set(value string) error {
	val, err := units.RAMInBytes(value)
	*m = memBytes(val)
	return err
}

func (m *memBytes) Type() string {
	return "MemoryBytes"
}

func (m *memBytes) Value() int64 {
	return int64(*m)
}

type nanoCPUs int64

func (c *nanoCPUs) String() string {
	return big.NewRat(c.Value(), 1e9).FloatString(3)
}

func (c *nanoCPUs) Set(value string) error {
	cpu, ok := new(big.Rat).SetString(value)
	if !ok {
		return fmt.Errorf("Failed to parse %v as a rational number", value)
	}
	nano := cpu.Mul(cpu, big.NewRat(1e9, 1))
	if !nano.IsInt() {
		return fmt.Errorf("value is too precise")
	}
	*c = nanoCPUs(nano.Num().Int64())
	return nil
}

func (c *nanoCPUs) Type() string {
	return "NanoCPUs"
}

func (c *nanoCPUs) Value() int64 {
	return int64(*c)
}

// DurationOpt is an option type for time.Duration that uses a pointer. This
// allows us to get nil values outside, instead of defaulting to 0
type DurationOpt struct {
	value *time.Duration
}

// Set a new value on the option
func (d *DurationOpt) Set(s string) error {
	v, err := time.ParseDuration(s)
	d.value = &v
	return err
}

// Type returns the type of this option
func (d *DurationOpt) Type() string {
	return "duration-ptr"
}

// String returns a string repr of this option
func (d *DurationOpt) String() string {
	if d.value != nil {
		return d.value.String()
	}
	return "none"
}

// Value returns the time.Duration
func (d *DurationOpt) Value() *time.Duration {
	return d.value
}

// Uint64Opt represents a uint64.
type Uint64Opt struct {
	value *uint64
}

// Set a new value on the option
func (i *Uint64Opt) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	i.value = &v
	return err
}

// Type returns the type of this option
func (i *Uint64Opt) Type() string {
	return "uint64-ptr"
}

// String returns a string repr of this option
func (i *Uint64Opt) String() string {
	if i.value != nil {
		return fmt.Sprintf("%v", *i.value)
	}
	return "none"
}

// Value returns the uint64
func (i *Uint64Opt) Value() *uint64 {
	return i.value
}

// MountOpt is a Value type for parsing mounts
type MountOpt struct {
	values []swarm.Mount
}

// Set a new mount value
func (m *MountOpt) Set(value string) error {
	csvReader := csv.NewReader(strings.NewReader(value))
	fields, err := csvReader.Read()
	if err != nil {
		return err
	}

	mount := swarm.Mount{}

	volumeOptions := func() *swarm.VolumeOptions {
		if mount.VolumeOptions == nil {
			mount.VolumeOptions = &swarm.VolumeOptions{
				Labels: make(map[string]string),
			}
		}
		if mount.VolumeOptions.DriverConfig == nil {
			mount.VolumeOptions.DriverConfig = &swarm.Driver{}
		}
		return mount.VolumeOptions
	}

	bindOptions := func() *swarm.BindOptions {
		if mount.BindOptions == nil {
			mount.BindOptions = new(swarm.BindOptions)
		}
		return mount.BindOptions
	}

	setValueOnMap := func(target map[string]string, value string) {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) == 1 {
			target[value] = ""
		} else {
			target[parts[0]] = parts[1]
		}
	}

<<<<<<< HEAD
	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) == 1 && strings.ToLower(parts[0]) == "writable" {
			mount.Writable = true
=======
	// Set writable as the default
	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) == 1 && strings.ToLower(parts[0]) == "readonly" {
			mount.ReadOnly = true
			continue
		}

		if len(parts) == 1 && strings.ToLower(parts[0]) == "volume-nocopy" {
			volumeOptions().NoCopy = true
>>>>>>> 12a5469... start on swarm services; move to glade
			continue
		}

		if len(parts) != 2 {
			return fmt.Errorf("invalid field '%s' must be a key=value pair", field)
		}

		key, value := parts[0], parts[1]
		switch strings.ToLower(key) {
		case "type":
			mount.Type = swarm.MountType(strings.ToUpper(value))
		case "source":
			mount.Source = value
		case "target":
			mount.Target = value
<<<<<<< HEAD
		case "writable":
			mount.Writable, err = strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for writable: %s", value)
			}
		case "bind-propagation":
			bindOptions().Propagation = swarm.MountPropagation(strings.ToUpper(value))
		case "volume-populate":
			volumeOptions().Populate, err = strconv.ParseBool(value)
=======
		case "readonly":
			ro, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for readonly: %s", value)
			}
			mount.ReadOnly = ro
		case "bind-propagation":
			bindOptions().Propagation = swarm.MountPropagation(strings.ToUpper(value))
		case "volume-nocopy":
			volumeOptions().NoCopy, err = strconv.ParseBool(value)
>>>>>>> 12a5469... start on swarm services; move to glade
			if err != nil {
				return fmt.Errorf("invalid value for populate: %s", value)
			}
		case "volume-label":
			setValueOnMap(volumeOptions().Labels, value)
		case "volume-driver":
			volumeOptions().DriverConfig.Name = value
		case "volume-driver-opt":
			if volumeOptions().DriverConfig.Options == nil {
				volumeOptions().DriverConfig.Options = make(map[string]string)
			}
			setValueOnMap(volumeOptions().DriverConfig.Options, value)
		default:
<<<<<<< HEAD
			return fmt.Errorf("unexpected key '%s' in '%s'", key, value)
=======
			return fmt.Errorf("unexpected key '%s' in '%s'", key, field)
>>>>>>> 12a5469... start on swarm services; move to glade
		}
	}

	if mount.Type == "" {
		return fmt.Errorf("type is required")
	}

	if mount.Target == "" {
		return fmt.Errorf("target is required")
	}

<<<<<<< HEAD
=======
	if mount.VolumeOptions != nil && mount.Source == "" {
		return fmt.Errorf("source is required when specifying volume-* options")
	}

	if mount.Type == swarm.MountType("BIND") && mount.VolumeOptions != nil {
		return fmt.Errorf("cannot mix 'volume-*' options with mount type '%s'", swarm.MountTypeBind)
	}
	if mount.Type == swarm.MountType("VOLUME") && mount.BindOptions != nil {
		return fmt.Errorf("cannot mix 'bind-*' options with mount type '%s'", swarm.MountTypeVolume)
	}

>>>>>>> 12a5469... start on swarm services; move to glade
	m.values = append(m.values, mount)
	return nil
}

// Type returns the type of this option
func (m *MountOpt) Type() string {
	return "mount"
}

// String returns a string repr of this option
func (m *MountOpt) String() string {
	mounts := []string{}
	for _, mount := range m.values {
		repr := fmt.Sprintf("%s %s %s", mount.Type, mount.Source, mount.Target)
		mounts = append(mounts, repr)
	}
	return strings.Join(mounts, ", ")
}

// Value returns the mounts
func (m *MountOpt) Value() []swarm.Mount {
	return m.values
}

type updateOptions struct {
	parallelism uint64
	delay       time.Duration
}

type resourceOptions struct {
	limitCPU      nanoCPUs
	limitMemBytes memBytes
	resCPU        nanoCPUs
	resMemBytes   memBytes
}

func (r *resourceOptions) ToResourceRequirements() *swarm.ResourceRequirements {
	return &swarm.ResourceRequirements{
		Limits: &swarm.Resources{
			NanoCPUs:    r.limitCPU.Value(),
			MemoryBytes: r.limitMemBytes.Value(),
		},
		Reservations: &swarm.Resources{
			NanoCPUs:    r.resCPU.Value(),
			MemoryBytes: r.resMemBytes.Value(),
		},
	}
}

type restartPolicyOptions struct {
	condition   string
	delay       DurationOpt
	maxAttempts Uint64Opt
	window      DurationOpt
}

func (r *restartPolicyOptions) ToRestartPolicy() *swarm.RestartPolicy {
	return &swarm.RestartPolicy{
		Condition:   swarm.RestartPolicyCondition(r.condition),
		Delay:       r.delay.Value(),
		MaxAttempts: r.maxAttempts.Value(),
		Window:      r.window.Value(),
	}
}

func convertNetworks(networks []string) []swarm.NetworkAttachmentConfig {
	nets := []swarm.NetworkAttachmentConfig{}
	for _, network := range networks {
		nets = append(nets, swarm.NetworkAttachmentConfig{Target: network})
	}
	return nets
}

type endpointOptions struct {
	mode  string
	ports opts.ListOpts
}

func (e *endpointOptions) ToEndpointSpec() *swarm.EndpointSpec {
	portConfigs := []swarm.PortConfig{}
	// We can ignore errors because the format was already validated by ValidatePort
	ports, portBindings, _ := nat.ParsePortSpecs(e.ports.GetAll())

	for port := range ports {
		portConfigs = append(portConfigs, convertPortToPortConfig(port, portBindings)...)
	}

	return &swarm.EndpointSpec{
		Mode:  swarm.ResolutionMode(strings.ToLower(e.mode)),
		Ports: portConfigs,
	}
}

func convertPortToPortConfig(
	port nat.Port,
	portBindings map[nat.Port][]nat.PortBinding,
) []swarm.PortConfig {
	ports := []swarm.PortConfig{}

	for _, binding := range portBindings[port] {
		hostPort, _ := strconv.ParseUint(binding.HostPort, 10, 16)
		ports = append(ports, swarm.PortConfig{
			//TODO Name: ?
			Protocol:      swarm.PortConfigProtocol(strings.ToLower(port.Proto())),
			TargetPort:    uint32(port.Int()),
			PublishedPort: uint32(hostPort),
		})
	}
	return ports
}

// ValidatePort validates a string is in the expected format for a port definition
func ValidatePort(value string) (string, error) {
	portMappings, err := nat.ParsePortSpec(value)
	for _, portMapping := range portMappings {
		if portMapping.Binding.HostIP != "" {
			return "", fmt.Errorf("HostIP is not supported by a service.")
		}
	}
	return value, err
}

type serviceOptions struct {
	name    string
	labels  opts.ListOpts
	image   string
	command []string
	args    []string
	env     opts.ListOpts
	workdir string
	user    string
	mounts  MountOpt

	resources resourceOptions
	stopGrace DurationOpt

	replicas Uint64Opt
	mode     string

	restartPolicy restartPolicyOptions
	constraints   []string
	update        updateOptions
	networks      []string
	endpoint      endpointOptions
<<<<<<< HEAD
=======

	registryAuth bool
>>>>>>> 12a5469... start on swarm services; move to glade
}

func newServiceOptions() *serviceOptions {
	return &serviceOptions{
		labels: opts.NewListOpts(runconfigopts.ValidateEnv),
		env:    opts.NewListOpts(runconfigopts.ValidateEnv),
		endpoint: endpointOptions{
			ports: opts.NewListOpts(ValidatePort),
		},
	}
}

func (opts *serviceOptions) ToService() (swarm.ServiceSpec, error) {
	var service swarm.ServiceSpec

	service = swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   opts.name,
			Labels: runconfigopts.ConvertKVStringsToMap(opts.labels.GetAll()),
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image:           opts.image,
				Command:         opts.command,
				Args:            opts.args,
				Env:             opts.env.GetAll(),
				Dir:             opts.workdir,
				User:            opts.user,
				Mounts:          opts.mounts.Value(),
				StopGracePeriod: opts.stopGrace.Value(),
			},
			Resources:     opts.resources.ToResourceRequirements(),
			RestartPolicy: opts.restartPolicy.ToRestartPolicy(),
			Placement: &swarm.Placement{
				Constraints: opts.constraints,
			},
		},
		Mode: swarm.ServiceMode{},
		UpdateConfig: &swarm.UpdateConfig{
			Parallelism: opts.update.parallelism,
			Delay:       opts.update.delay,
		},
		Networks:     convertNetworks(opts.networks),
		EndpointSpec: opts.endpoint.ToEndpointSpec(),
	}

	switch opts.mode {
	case "global":
		if opts.replicas.Value() != nil {
			return service, fmt.Errorf("replicas can only be used with replicated mode")
		}

		service.Mode.Global = &swarm.GlobalService{}
	case "replicated":
		service.Mode.Replicated = &swarm.ReplicatedService{
			Replicas: opts.replicas.Value(),
		}
	default:
		return service, fmt.Errorf("Unknown mode: %s", opts.mode)
	}
	return service, nil
}

<<<<<<< HEAD
// addServiceFlags adds all flags that are common to both `create` and `update.
=======
// addServiceFlags adds all flags that are common to both `create` and `update`.
>>>>>>> 12a5469... start on swarm services; move to glade
// Any flags that are not common are added separately in the individual command
func addServiceFlags(cmd *cobra.Command, opts *serviceOptions) {
	flags := cmd.Flags()
	flags.StringVar(&opts.name, flagName, "", "Service name")
	flags.VarP(&opts.labels, flagLabel, "l", "Service labels")

	flags.VarP(&opts.env, "env", "e", "Set environment variables")
	flags.StringVarP(&opts.workdir, "workdir", "w", "", "Working directory inside the container")
	flags.StringVarP(&opts.user, flagUser, "u", "", "Username or UID")
	flags.VarP(&opts.mounts, flagMount, "m", "Attach a mount to the service")

	flags.Var(&opts.resources.limitCPU, flagLimitCPU, "Limit CPUs")
	flags.Var(&opts.resources.limitMemBytes, flagLimitMemory, "Limit Memory")
	flags.Var(&opts.resources.resCPU, flagReserveCPU, "Reserve CPUs")
	flags.Var(&opts.resources.resMemBytes, flagReserveMemory, "Reserve Memory")
	flags.Var(&opts.stopGrace, flagStopGracePeriod, "Time to wait before force killing a container")

	flags.Var(&opts.replicas, flagReplicas, "Number of tasks")

<<<<<<< HEAD
	flags.StringVar(&opts.restartPolicy.condition, flagRestartCondition, "", "Restart when condition is met (none, on_failure, or any)")
=======
	flags.StringVar(&opts.restartPolicy.condition, flagRestartCondition, "", "Restart when condition is met (none, on-failure, or any)")
>>>>>>> 12a5469... start on swarm services; move to glade
	flags.Var(&opts.restartPolicy.delay, flagRestartDelay, "Delay between restart attempts")
	flags.Var(&opts.restartPolicy.maxAttempts, flagRestartMaxAttempts, "Maximum number of restarts before giving up")
	flags.Var(&opts.restartPolicy.window, flagRestartWindow, "Window used to evaluate the restart policy")

	flags.StringSliceVar(&opts.constraints, flagConstraint, []string{}, "Placement constraints")

	flags.Uint64Var(&opts.update.parallelism, flagUpdateParallelism, 0, "Maximum number of tasks updated simultaneously")
	flags.DurationVar(&opts.update.delay, flagUpdateDelay, time.Duration(0), "Delay between updates")

	flags.StringSliceVar(&opts.networks, flagNetwork, []string{}, "Network attachments")
<<<<<<< HEAD
	flags.StringVar(&opts.endpoint.mode, flagEndpointMode, "", "Endpoint mode(Valid values: vip, dnsrr)")
	flags.VarP(&opts.endpoint.ports, flagPublish, "p", "Publish a port as a node port")
=======
	flags.StringVar(&opts.endpoint.mode, flagEndpointMode, "", "Endpoint mode (vip or dnsrr)")
	flags.VarP(&opts.endpoint.ports, flagPublish, "p", "Publish a port as a node port")

	flags.BoolVar(&opts.registryAuth, flagRegistryAuth, false, "Send registry authentication details to Swarm agents")
>>>>>>> 12a5469... start on swarm services; move to glade
}

const (
	flagConstraint         = "constraint"
	flagEndpointMode       = "endpoint-mode"
	flagLabel              = "label"
	flagLimitCPU           = "limit-cpu"
	flagLimitMemory        = "limit-memory"
	flagMode               = "mode"
	flagMount              = "mount"
	flagName               = "name"
	flagNetwork            = "network"
	flagPublish            = "publish"
	flagReplicas           = "replicas"
	flagReserveCPU         = "reserve-cpu"
	flagReserveMemory      = "reserve-memory"
	flagRestartCondition   = "restart-condition"
	flagRestartDelay       = "restart-delay"
	flagRestartMaxAttempts = "restart-max-attempts"
	flagRestartWindow      = "restart-window"
	flagStopGracePeriod    = "stop-grace-period"
	flagUpdateDelay        = "update-delay"
	flagUpdateParallelism  = "update-parallelism"
	flagUser               = "user"
<<<<<<< HEAD
=======
	flagRegistryAuth       = "registry-auth"
>>>>>>> 12a5469... start on swarm services; move to glade
)