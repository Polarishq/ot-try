package backend

func init() {
	Constants = constants{
		ComponentName: "ot-try",
		Version:       "v1",
	}
}

type constants struct {
	ComponentName string
	Version       string
}

// Constants contains static information regarding this service
var Constants constants

func (c *constants) ServiceName() string {
	return c.ComponentName + "-" + c.Version
}
