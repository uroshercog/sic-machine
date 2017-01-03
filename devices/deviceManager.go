package devices

// DeviceManager ..
type DeviceManager struct {
	devices map[uint8]Device
}

// Get ...
func (dm *DeviceManager) Get(fd byte) Device {
	dev, ok := dm.devices[fd]

	if ok {
		return dev
	}

	// TODO: Device does not exist, open it
	return nil
}

// Set ...
func (dm *DeviceManager) Set(fd byte, device Device) {
	dm.devices[fd] = device
}

// New ...
func New() *DeviceManager {
	return &DeviceManager{
		devices: map[uint8]Device{
			0: NewStdinDevice(),
			1: NewStdoutDevice(),
			2: NewStderrDevice(),
		},
	}
}
