package space

import "context"

var globalAppSpace = NewSpace(context.Background(), GlobalNameKey)

func GlobalSpace() context.Context {
	return globalAppSpace
}
