package cmd

type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	r := &Registry{commands: make(map[string]Command)}
	return r
}

func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

func (r *Registry) List() []Command {
	list := []Command{}
	for _, cmd := range r.commands {
		list = append(list, cmd)
	}
	return list
}
