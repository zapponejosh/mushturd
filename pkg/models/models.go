package models

type Checkpoint struct {
	Name    string
	In      bool
	DogsIn  int
	Out     bool
	DogsOut int
}

type Musher struct {
	CurrentPos       int
	Name             string
	Rookie           bool
	Bib              int
	LatestCheckpoint Checkpoint
	Speed            float32
	InCheckpoint     bool
	Status           string
}
