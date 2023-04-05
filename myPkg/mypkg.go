package mypkg

var owner = struct {
	name string
}{
	name: "kimkyungsub",
}

func Owner() struct{ name string } {
	return owner
}
