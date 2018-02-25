package main

import (
	"os"

	"github.com/mzohreva/GoGraphviz/graphviz"
)

func main() {
	G := graphviz.Graph{}
	G.MakeDirected()
	n1 := G.AddNode("Hello")
	n2 := G.AddNode("World")
	n3 := G.AddNode("Hi")
	n4 := G.AddNode("NULL")
	G.AddEdge(n1, n2, "next")
	G.AddEdge(n2, n3, "next")
	G.AddEdge(n3, n4, "next")
	G.MakeSameRank(n1, n2, n3, n4)

	G.GraphAttribute(graphviz.NodeSep, "0.5")

	G.DefaultNodeAttribute(graphviz.Shape, graphviz.ShapeBox)
	G.DefaultNodeAttribute(graphviz.FontName, "Courier")
	G.DefaultNodeAttribute(graphviz.FontSize, "14")
	G.DefaultNodeAttribute(graphviz.Style, graphviz.StyleFilled+","+graphviz.StyleRounded)
	G.DefaultNodeAttribute(graphviz.FillColor, "yellow")

	G.NodeAttribute(n4, graphviz.Shape, graphviz.ShapeCircle)
	G.NodeAttribute(n4, graphviz.Style, graphviz.StyleDashed)

	G.DefaultEdgeAttribute(graphviz.FontName, "Courier")
	G.DefaultEdgeAttribute(graphviz.FontSize, "12")

	G.GenerateDOT(os.Stdout)
}
