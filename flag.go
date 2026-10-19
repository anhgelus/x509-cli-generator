package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

type Duration struct {
	dur time.Duration
	has bool
}

func (d *Duration) String() string {
	if d == nil {
		return "<empty>"
	}
	return d.dur.String()
}

func (d *Duration) Set(v string) error {
	d.has = true
	v = strings.TrimSpace(v)
	parts := strings.Split(v, " ")
	for _, p := range parts {
		if len(p) < 2 {
			return fmt.Errorf("malformed duration string: %s", p)
		}
		last := len(p) - 1
		n, err := strconv.Atoi(p[:last])
		if err != nil {
			return err
		}
		switch p[last] {
		case 'd':
			d.dur += time.Duration(n) * time.Hour * 24
		case 'm':
			d.dur += time.Duration(n) * time.Hour * 24 * 30
		case 'y':
			d.dur += time.Duration(n) * time.Hour * 24 * 365
		default:
			return fmt.Errorf("malformed duration string: %s", p)
		}
	}
	return nil
}
