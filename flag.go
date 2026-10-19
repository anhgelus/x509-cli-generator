package main

import "fmt"

type Slice[T any] struct {
	Content []T
}

func (sl *Slice[T]) String() string {
	if len(sl.Content) == 0 {
		return "<nil>"
	}
	return fmt.Sprint(sl.Content)
}

func (sl *Slice[T]) Set(v string) error {
	var val T
	_, err := fmt.Sscan(v, &val)
	if err != nil {
		return err
	}
	sl.Content = append(sl.Content, val)
	return nil
}
