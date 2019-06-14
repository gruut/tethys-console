package utils

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func plusOne(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range in {
			out <- v + 1
		}
	}()

	return out
}
func TestPlusOne(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"PlusOne",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan int)
			values := []int{1, 2, 3}

			go func(values []int) {
				for _, v := range values {
					in <- v
				}
				close(in)
			}(values)

			for v := range plusOne(in) {
				fmt.Println(v)
			}
		})
	}
}

func TestChain(t *testing.T) {
	type args struct {
		ps           []IntPipe
		initialValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"PlusTwo",
			args{[]IntPipe{plusOne, plusOne}, 1},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plusTwo := Chain(tt.args.ps...)

			ch := make(chan int)
			go func() {
				defer close(ch)
				ch <- 1
			}()
			result := <-plusTwo(ch)

			if result != tt.want {
				log.Fatalf("Got a unexpected value, expected: %d, actual: %d", tt.want, result)
			}
		})
	}
}

func TestFanIn(t *testing.T) {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)

	go func() {
		defer close(c1)
		c1<-1
	}()
	go func() {
		defer close(c2)
		c2<-2
	}()
	go func() {
		defer close(c3)
		c3<-3
	}()

	channels := []<-chan int{c1, c2, c3}

	type args struct {
		ins []<-chan int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"FanIn",
			args{channels},
			[]int{1,2,3,4,5,6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FanIn(tt.args.ins...)
			
			arr := []int{}
			for v := range got {
				arr = append(arr, v)
			}

			if !reflect.DeepEqual(arr, tt.want) {
				t.Errorf("FanIn() = %v, want %v", arr, tt.want)
			}
		})
	}
}
