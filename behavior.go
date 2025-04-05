package main

type Trans struct{}

type None struct{}

type Lt struct {
	Layer LayerIndex
	Tap   KeyCode
}

type To struct {
	Layer LayerIndex
}

type Mt struct {
	Hold KeyCode
	Tap  KeyCode
}

type Kp struct {
	Tap KeyCode
}

type KpKp struct {
	Hold KeyCode
	Tap  KeyCode
}

type Rmt struct {
	Hold KeyCode
	Tap  KeyCode
}
