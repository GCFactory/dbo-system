package main

import (
	"os"
	"os/exec"
	"fmt"
)

const (
	GraphTypeUsual string = "Usual"
	GraphTypeDig string = "Gigraph"
)

const (
	GraphOrientationTopBottom string = "TB"
	GraphOrientationBottomTop string = "BT"
	GraphOrientationLeftRight string = "LR"
	GraphOrientationRightLeft string = "RL"
)

type Node struct {
	Name string
	prop map[string]string
}

func createNode(name string) *Node {

	return &Node{
		Name: name,
		prop: make(map[string]string),
	}

}

func printProperties(properties map[string]string) string {

	result := "["

	for key, value := range properties {
		result += "\"" + key + "\"=\"" + value + "\";"
	}

	result += "]"

	return result

}

func (node* Node) print() string {

	result := ""

	if node != nil && node.Name != "" {

		result = node.Name + " " + printProperties(node.prop) + ";"

	}

	return result

}

func (node *Node) setProperty(key string, value string) {

	node.prop[key] = value

}

type Edge struct {
	graphType string
	prop map[string]string
	from *Node
	to *Node
}

func createEdge(from *Node, to *Node) *Edge{

	return &Edge{
		graphType: GraphTypeUsual,
		prop: make(map[string]string),
		from: from,
		to: to,
	}

}

func (edge *Edge) setProperty(key string, value string) {

	edge.prop[key] = value

}

func (edge *Edge) print() string {

	result := ""

	if edge != nil &&
		edge.from != nil &&
		edge.to != nil {

		connect := "--"

		if edge.graphType == GraphTypeDig {
			connect = "->"
		}

		result = edge.from.Name + " " + connect + " " + edge.to.Name + " " + printProperties(edge.prop) + ";"

	}

	return result

}

func (edge *Edge) setGraphType(graph_type string) {

	edge.graphType = graph_type

}

func (edge *Edge) setStartCluster(name string) {
	edge.prop["ltail"] = name
}

func (edge *Edge) setEndCluster(name string) {
	edge.prop["lhead"] = name
}

type SubGraph struct {
	name string
	label string
	nodes []*Node
	sub_graphs []*SubGraph
	prop map[string]string
}

func createSubGraph(name string) *SubGraph{

	return &SubGraph{
		name: name,
		label: "",
		nodes: make([]*Node, 0),
		sub_graphs: make([]*SubGraph, 0),
		prop: make(map[string]string),
	}

}

func (graph *SubGraph) getFullName() string {
	return "cluster_" + graph.name
}

func (graph *SubGraph) setLabel(label string) {
	graph.label = label
}

func (graph *SubGraph) addNodes(nodes ...*Node) {
	graph.nodes = append(graph.nodes, nodes...)
}

func (graph *SubGraph) addSubGraphs(graphs ...*SubGraph) {
	graph.sub_graphs = append(graph.sub_graphs, graphs...)
}

func (graph *SubGraph) setProperty(key string, value string) {
	graph.prop[key] = value
}

func (graph *SubGraph) print() string {

	result := "subgraph \"cluster_" + graph.name + "\" {\n"

	if graph.label != "" {
		result += "\tlabel=\"" + graph.label + "\";\n"
	}

	if len(graph.nodes) != 0 {
		result += "\t"
		for _, node := range graph.nodes {
			result += node.Name + "; "
		}
		result += "\n"
	}

	result += "\tgraph["
	for key, value := range graph.prop {
		result += key + "=\"" + value + "\";"
	}
	result += "];\n"

	if len(graph.sub_graphs) != 0 {

		for _, sub_graph := range graph.sub_graphs {

			result += sub_graph.print() + "\n"

		}

	}

	result += "\t}"

	return result

}

const (
	RankTypeSame string = "same"
	RankTypeMin string = "min"
	RankTypeSource string = "source"
)

type Rank struct {
	rank_type string
	nodes []*Node
}

func createRank() *Rank {

	return &Rank{
		rank_type: RankTypeSame,
		nodes: make([]*Node, 0),
	}

}

func (rank *Rank) addNodes(nodes ...*Node) {
	rank.nodes = append(rank.nodes, nodes...)
}

func (rank *Rank) setType(rank_type string) {
	rank.rank_type = rank_type
}

func (rank *Rank) print() string {

	result := "{ "

	result += "rank=\"" + rank.rank_type + "\";"

	for _, node := range rank.nodes {
		result += " " + node.Name + ";"
	}

	result += " }"

	return result

}

type Graph struct{
	orientation string
	nodes []*Node
	edges []*Edge
	graph_type string
	sub_graphs []*SubGraph
	rank_is_on bool
	ranks []*Rank
}

func createGraph(orientation string, graph_type string) *Graph {

	return &Graph{
		orientation: orientation,
		nodes: make([]*Node, 0),
		edges: make([]*Edge, 0),
		graph_type: graph_type,
		sub_graphs: make([]*SubGraph, 0),
		rank_is_on: false,
		ranks: make([]*Rank, 0),
	}

}

func (graph *Graph) addEdges(edges ...*Edge){

	graph.edges = append(graph.edges, edges...)

}

func (graph *Graph) addNodes(nodes ...*Node){

	graph.nodes = append(graph.nodes, nodes...)

}

func (graph *Graph) addSubGraphs(graphs ...*SubGraph) {

	graph.sub_graphs = append(graph.sub_graphs, graphs...)

}

func (graph *Graph) print() string {

	result := ""

	if graph != nil {
		switch graph.graph_type {
			case GraphTypeDig: {
				result += "digraph {\n"
			}
			case GraphTypeUsual: {
				result += "graph {\n"
			}
			default: {
				result += "graph {\n"
			}
		}

		switch graph.orientation {
			case 	GraphOrientationTopBottom,
				GraphOrientationBottomTop,
				GraphOrientationLeftRight,
				GraphOrientationRightLeft:
			{
				result += "\trankdir=" + graph.orientation + ";\n"
			}
		}

		if graph.rank_is_on {
			result += "\tnewrank=\"true\"\n"
		}

		if graph.nodes != nil {

			for _, node := range graph.nodes {

				result += "\t" + node.print() + "\n"

			}

		}

		if graph.edges != nil {

			for _, edge := range graph.edges {

				edge.setGraphType(graph.graph_type)

				result += "\t" + edge.print() + "\n"

			}

		}

		if graph.sub_graphs != nil {

			for _, sub_graph := range graph.sub_graphs {

				result += "\t" + sub_graph.print() + "\n"
			}

		}

		if graph.ranks != nil {

			for _, rank := range graph.ranks {
				result += "\t" + rank.print() + "\n"
			}

		}

		result += "}\n"

	}

	return result

}

func (graph *Graph) turnOnRank() {
	graph.rank_is_on = true
}

func (graph *Graph) turnOffRank() {
	graph.rank_is_on = false
}

func (graph *Graph) addRanks(ranks ...*Rank) {
	graph.ranks = append(graph.ranks, ranks...)
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func createGraphImage(config_file_name string, img_file_name string, graph_config string) {

	local_config_file_name := "./" + config_file_name + ".dot"
	local_img_file_name := "./" + img_file_name + ".svg"

	ex, err := exists(local_config_file_name)
	if err != nil {
		return
	}
	if ex {
		err = os.Remove(local_config_file_name)
		if err != nil {
			return
		}
	}

	ex, err = exists(local_img_file_name)
	if err != nil {
		return
	}
	if ex {
		err = os.Remove(local_img_file_name)
		if err != nil {
			return
		}
	}

	cfg_file, err := os.OpenFile(local_config_file_name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	defer cfg_file.Close()

	if err != nil {

		return
	}

	nbt, err := cfg_file.WriteString( graph_config )
	if err != nil || nbt == 0 {
		return
	}

	_, err = os.Create(local_img_file_name)
	if err != nil {
		return
	}

	cmd := exec.Command("dot", "-Tsvg", local_config_file_name,"-o", local_img_file_name)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
}

func main(){
	fmt.Println("begin\n")

	g := createGraph(GraphOrientationTopBottom, GraphTypeDig)
	g.turnOnRank()

	sub_graph_operation := createSubGraph("sub_graph_operation")
	sub_graph_operation.setLabel("Операция открытия счёта")

	sub_graph_check := createSubGraph("sub_graph_check")

	node_saga_check := createNode("node_saga_check")
	node_saga_check.setProperty("label", "SAGA получения данных пользователя")

	sub_graph_check_events := createSubGraph("sub_graph_check_events")

	node_event_check := createNode("node_event_check")
	node_event_check.setProperty("label", "Операция резервирования счёта")

	sub_graph_check_events.addNodes(
		node_event_check,
	)

	sub_graph_check.addNodes(
		node_saga_check,
	)

	sub_graph_check.addSubGraphs(
		sub_graph_check_events,
	)

	sub_graph_reserve := createSubGraph("sub_graph_reserve")

	node_saga_reserve := createNode("node_saga_reserve")
	node_saga_reserve.setProperty("label", "SAGA резервирования счёта")

	sub_graph_reserve.addNodes(
		node_saga_reserve,
	)

	sub_graph_create := createSubGraph("sub_graph_create")

	node_saga_create := createNode("node_saga_create")
	node_saga_create.setProperty("label", "SAGA создания счёта")

	sub_graph_create.addNodes(
		node_saga_create,
	)

	sub_graph_open_and_add := createSubGraph("sub_graph_open_and_add")

	node_saga_open_and_add := createNode("node_saga_open_and_add")
	node_saga_open_and_add.setProperty("label", "SAGA открытия счёта")

	sub_graph_open_and_add.addNodes(
		node_saga_open_and_add,
	)

	sub_graph_operation.addSubGraphs(
		sub_graph_check,
		sub_graph_reserve,
		sub_graph_create,
		sub_graph_open_and_add,
	)

	g.addSubGraphs(
		sub_graph_operation,
	)

	g.addNodes(
		node_saga_check,
		node_saga_reserve,
		node_saga_create,
		node_saga_open_and_add,
		node_event_check,
	)

	edge_saga1_saga2 := createEdge(node_saga_check, node_saga_reserve)
	// edge_saga1_saga2.setStartCluster(sub_graph_check.getFullName())
	// edge_saga1_saga2.setEndCluster(sub_graph_reserve.getFullName())

	edge_saga1_event1 := createEdge(node_saga_check, node_event_check)
	edge_saga1_event1.setProperty("arrowsize", "0.0")
	// edge_saga1_event1.setEndCluster(sub_graph_check_events.getFullName())

	edge_saga2_saga3 := createEdge(node_saga_reserve, node_saga_create)
	edge_saga3_saga4 := createEdge(node_saga_create, node_saga_open_and_add)

	g.addEdges(
		edge_saga1_saga2,
		edge_saga2_saga3,
		edge_saga3_saga4,
		edge_saga1_event1,
	)

	rank_saga := createRank()
	rank_saga.addNodes(
		node_saga_check,
		node_saga_reserve,
		node_saga_create,
		node_saga_open_and_add,
	)

	g.addRanks(
		rank_saga,
	)

	g_config := g.print()

	fmt.Println(g_config)

	createGraphImage("config", "test", g_config)

	fmt.Println("end")
}
