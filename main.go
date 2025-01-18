package main

func main() {

	e := Engine{}
	if err := e.Init(); err != nil {
		e.Destroy()
		panic(err)
	}
	defer e.Destroy()

	game := Game{}

	err := game.Init(&e)
	if err != nil {
		panic(err)
	}

	e.Run()
}


