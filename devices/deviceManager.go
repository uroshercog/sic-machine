package devices

// DeviceManager ..
type DeviceManager struct {
	devices map[byte]Device
}

// Get ...
func (dm *DeviceManager) Get(fd byte) Device {
	if dev, ok := dm.devices[fd]; ok {
		return dev
	}

	// TODO: Device does not exist, open it
	// This is a file device
	d := NewFileDevice(fd)
	dm.devices[fd] = d
	return d
}

// Set ...
func (dm *DeviceManager) Set(fd byte, device Device) {
	dm.devices[fd] = device
}

// New ...
func New() *DeviceManager {
	return &DeviceManager{
		devices: map[byte]Device{
			0: NewStdinDevice(),
			1: NewStdoutDevice(),
			2: NewStderrDevice(),
		},
	}
}
