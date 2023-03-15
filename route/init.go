package route

var Routes = map[string]*RouterData{}

func Add(name string, route *RouterData) {
	Routes[name] = route
}

func Get(name string) *RouterData {
	return Routes[name]
}
