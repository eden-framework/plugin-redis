package redis

func Command(name string, args ...interface{}) *CMD {
	return &CMD{
		name: name,
		args: args,
	}
}

type CMD struct {
	name string
	args []interface{}
}
