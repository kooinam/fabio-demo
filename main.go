package main

import (
	fab "github.com/kooinam/fabio"
)

func main() {
	fab.ConfigureAndServe(&FabInitializer{})
}
