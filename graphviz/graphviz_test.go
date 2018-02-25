package graphviz_test

import (
	"os"

	"github.com/mzohreva/GoGraphviz/graphviz"
)

// This example shows how Graph can be used to display a simple linked list.
// The output can be piped to the dot tool to generate an image.
func Example_linkedList() {
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
	// output:
	// strict digraph {
	//   nodesep = "0.5";
	//   node [ shape = "box" ]
	//   node [ fontname = "Courier" ]
	//   node [ fontsize = "14" ]
	//   node [ style = "filled,rounded" ]
	//   node [ fillcolor = "yellow" ]
	//   edge [ fontname = "Courier" ]
	//   edge [ fontsize = "12" ]
	//   n0 [label="Hello"]
	//   n1 [label="World"]
	//   n2 [label="Hi"]
	//   n3 [label="NULL", shape="circle", style="dashed"]
	//   {rank=same; n0; n1; n2; n3; }
	//   n0 -> n1 [label="next"]
	//   n1 -> n2 [label="next"]
	//   n2 -> n3 [label="next"]
	// }
}
