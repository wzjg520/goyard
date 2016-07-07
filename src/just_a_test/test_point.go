package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	a := &ListNode{1, nil}
	b := &ListNode{8, nil}
	a.Next = b

	for {
		if a.Next != nil {
			fmt.Println(a.Val)
			fmt.Println(a.Next)
			a = a.Next
			fmt.Println(a)
			fmt.Println(a.Next)
		} else {
			break
		}

	}
}
