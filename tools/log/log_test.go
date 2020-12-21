package log

import (
	"testing"
)

func TestLogPath(t *testing.T) {
	//fmt.Print(getCurrentDirectory())
	//fmt.Println(os.Args[1])
	n4 := ListNode{4, nil}
	n3 := ListNode{3, &n4}
	n2 := ListNode{2, &n3}
	n1 := ListNode{1, &n2}
	swapPairs(&n1)
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	pre := head
	post := head.Next
	nt := post
	ln := pre
	for post != nil {
		mid := post.Next
		post.Next = pre
		pre.Next = nil
		if mid == nil || mid.Next == nil {
			break
		}
		pre = mid
		post = mid.Next
		ln.Next = post
		ln = pre
	}
	return nt
}
