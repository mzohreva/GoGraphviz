package main

import (
	"log"

	"github.com/mzohreva/GoGraphviz/graphviz"
)

func main() {
	list := linkedList{}
	visualizeLinkedList(list, "list-1.png", "Empty linked list")

	list.append("hello")
	list.append("world")
	visualizeLinkedList(list, "list-2.png", "After appending 'hello' and 'world'")

	list.prepend("first!")
	visualizeLinkedList(list, "list-3.png", "After prepending 'first!'")
}

func visualizeLinkedList(list linkedList, filename, title string) {
	G := &graphviz.Graph{}
	nodes := make([]int, 1)
	nodes[0] = G.AddNode("head")
	head := list.head
	for head != nil {
		nodes = append(nodes, G.AddNode(head.value))
		head = head.next
	}
	nilNode := G.AddNode("nil")
	nodes = append(nodes, nilNode)

	G.AddEdge(nodes[0], nodes[1], "")
	for i := 1; i < len(nodes)-1; i++ {
		G.AddEdge(nodes[i], nodes[i+1], "next")
	}
	G.NodeAttribute(nodes[0], graphviz.Shape, graphviz.ShapeNone)
	G.NodeAttribute(nilNode, graphviz.Shape, graphviz.ShapeNone)
	G.MakeSameRank(nodes[0], nodes[1], nodes[2:]...)

	G.DefaultNodeAttribute(graphviz.Shape, graphviz.ShapeBox)
	G.DefaultNodeAttribute(graphviz.FontName, "Courier")
	G.DefaultEdgeAttribute(graphviz.FontName, "Courier")
	G.GraphAttribute(graphviz.NodeSep, "0.5")
	G.SetTitle("\n\n" + title)
	G.MakeDirected()

	err := G.GenerateImage("dot", filename, "png")
	if err != nil {
		log.Fatal(err)
	}
}

type linkedList struct {
	head, tail *listNode
}

func (list *linkedList) append(value string) {
	node := &listNode{value, nil}
	if list.head == nil {
		list.head, list.tail = node, node
		return
	}
	list.tail.next = node
	list.tail = node
}

func (list *linkedList) prepend(value string) {
	node := &listNode{value, nil}
	if list.head == nil {
		list.head, list.tail = node, node
		return
	}
	node.next = list.head
	list.head = node
}

type listNode struct {
	value string
	next  *listNode
}
