package usecase

import (
	"fmt"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"strconv"
)

const (
	GraphTypeUsual string = "Usual"
	GraphTypeDig   string = "Gigraph"
)

const (
	GraphNoOrientation        string = ""
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

func (node *Node) print() string {

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
	prop      map[string]string
	from      *Node
	to        *Node
}

func createEdge(from *Node, to *Node) *Edge {

	return &Edge{
		graphType: GraphTypeUsual,
		prop:      make(map[string]string),
		from:      from,
		to:        to,
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
	name       string
	label      string
	nodes      []*Node
	sub_graphs []*SubGraph
	prop       map[string]string
}

func createSubGraph(name string) *SubGraph {

	return &SubGraph{
		name:       name,
		label:      "",
		nodes:      make([]*Node, 0),
		sub_graphs: make([]*SubGraph, 0),
		prop:       make(map[string]string),
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
	RankTypeSame   string = "same"
	RankTypeMin    string = "min"
	RankTypeMax    string = "max"
	RankTypeSource string = "source"
)

type Rank struct {
	rank_type string
	number    int
	nodes     []*Node
}

func createRank() *Rank {

	return &Rank{
		rank_type: RankTypeSame,
		number:    0,
		nodes:     make([]*Node, 0),
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

	if rank.number != 0 {
		result += " " + strconv.Itoa(rank.number) + ";"
	}

	for _, node := range rank.nodes {
		result += " " + node.Name + ";"
	}

	result += " }"

	return result

}

func (rank *Rank) setRankNumber(number int) {
	rank.number = number
}

type Graph struct {
	orientation          string
	nodes                []*Node
	edges                []*Edge
	edgesBetweenClusters bool
	graph_type           string
	sub_graphs           []*SubGraph
	rank_is_on           bool
	ranks                []*Rank
}

func createGraph(orientation string, graph_type string) *Graph {

	return &Graph{
		orientation:          orientation,
		nodes:                make([]*Node, 0),
		edges:                make([]*Edge, 0),
		edgesBetweenClusters: false,
		graph_type:           graph_type,
		sub_graphs:           make([]*SubGraph, 0),
		rank_is_on:           false,
		ranks:                make([]*Rank, 0),
	}

}

func (graph *Graph) addEdges(edges ...*Edge) {

	graph.edges = append(graph.edges, edges...)

}

func (graph *Graph) addNodes(nodes ...*Node) {

	graph.nodes = append(graph.nodes, nodes...)

}

func (graph *Graph) addSubGraphs(graphs ...*SubGraph) {

	graph.sub_graphs = append(graph.sub_graphs, graphs...)

}

func (graph *Graph) print() string {

	result := ""

	if graph != nil {
		switch graph.graph_type {
		case GraphTypeDig:
			{
				result += "digraph {\n"
			}
		case GraphTypeUsual:
			{
				result += "graph {\n"
			}
		default:
			{
				result += "graph {\n"
			}
		}

		switch graph.orientation {
		case GraphOrientationTopBottom,
			GraphOrientationBottomTop,
			GraphOrientationLeftRight,
			GraphOrientationRightLeft:
			{
				result += "\trankdir=" + graph.orientation + ";\n"
			}
		case GraphNoOrientation:
			{
			}
		}

		if graph.edgesBetweenClusters {
			result += "\tcompound=\"true\"\n"
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

func (graph *Graph) turnOnEdgesBetweenClusters() {
	graph.edgesBetweenClusters = true
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

func createGraphImage(folder_path string, config_file_name string, img_file_name string, graph_config string) (string, error) {

	local_config_file_name := folder_path + "/" + config_file_name + ".dot"
	local_img_file_name := folder_path + "/" + img_file_name + ".svg"

	ex, err := exists(local_config_file_name)
	if err != nil {
		return "", err
	}
	if ex {
		err = os.Remove(local_config_file_name)
		if err != nil {
			return "", err
		}
	}

	ex, err = exists(local_img_file_name)
	if err != nil {
		return "", err
	}
	if ex {
		err = os.Remove(local_img_file_name)
		if err != nil {
			return "", err
		}
	}

	cfg_file, err := os.OpenFile(local_config_file_name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	defer cfg_file.Close()

	if err != nil {
		return "", err
	}

	nbt, err := cfg_file.WriteString(graph_config)
	if err != nil || nbt == 0 {
		return "", err
	}

	_, err = os.Create(local_img_file_name)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("dot", "-Tsvg", local_config_file_name, "-o", local_img_file_name)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды dot:", err)
		return "", err
	}

	return img_file_name + ".svg", nil
}

func CreateGraph(operation_tree *models.OperationTree, folder_path string, file_name string) (string, error) {

	g := createGraph(GraphNoOrientation, GraphTypeDig)
	//g.turnOnRank()
	g.turnOnEdgesBetweenClusters()

	main_graph := createSubGraph("main_graph")
	main_graph.setLabel("Operation:" + operation_tree.OperationName)

	saga_sub_graph_list := make(map[uuid.UUID]*SubGraph, 0)
	event_nodes_list := make(map[uuid.UUID]*Node, 0)

	ranks_list := make([]*Rank, 0)
	for saga_id, saga := range operation_tree.SagaList {
		saga_sub_graph := createSubGraph(saga.Name + "_sub_graph")
		saga_sub_graph.setLabel("SAGA:" + saga.Name)

		saga_sub_graph_list[saga_id] = saga_sub_graph
		switch saga.Status {
		case 30:
			{
				//saga_sub_graph.setProperty("fillcolor", "seagreen")
				//saga_sub_graph.setProperty("style", "filled")
			}
		}

		main_graph.addSubGraphs(saga_sub_graph)

		//var rank *Rank = nil
		//if len(saga.Events) > 1 {
		//	rank = createRank()
		//}
		for _, event_id := range saga.Events {
			if event, ok := operation_tree.EventList[event_id]; ok {
				event_node := createNode(event.Name + "_event")
				//if rank != nil {
				//	rank.addNodes(event_node)
				//}

				event_node.setProperty("label", event.Name)
				switch event.Status {
				case 0:
					{
					}
				case 10:
					{
						event_node.setProperty("color", "gray")
						event_node.setProperty("style", "filled")
					}
				case 20:
					{
						event_node.setProperty("color", "cornflowerblue")
						event_node.setProperty("style", "filled")
					}
				case 30:
					{
						event_node.setProperty("color", "seagreen")
						event_node.setProperty("style", "filled")
					}
				case 40:
					{
						event_node.setProperty("color", "orangered")
						event_node.setProperty("style", "filled")
					}
				case 50:
					{
						event_node.setProperty("color", "orange")
						event_node.setProperty("style", "filled")
					}
				case 250:
					{
						event_node.setProperty("color", "darkorchid4")
						event_node.setProperty("style", "filled")
					}
				case 255:
					{
						event_node.setProperty("color", "darkred")
						event_node.setProperty("style", "filled")
					}
				}

				event_nodes_list[event.Id] = event_node

				saga_sub_graph.addNodes(event_node)
				g.addNodes(event_node)
			}
		}

		//if rank != nil {
		//	g.addRanks(rank)
		//}

	}
	for _, saga_depend := range operation_tree.SagaDependList {
		saga, ok := operation_tree.SagaList[saga_depend.ParentId]
		if ok {
			event_parent, ok := operation_tree.EventList[saga.Events[0]]
			if ok {
				sub_grapth_parent, ok := saga_sub_graph_list[saga.Id]
				if ok {
					saga, ok := operation_tree.SagaList[saga_depend.ChildId]
					if ok {
						event_child, ok := operation_tree.EventList[saga.Events[0]]
						if ok {
							node_event_child, ok := event_nodes_list[event_child.Id]
							if ok {
								node_event_parent, ok := event_nodes_list[event_parent.Id]
								if ok {
									sub_graph_child, ok := saga_sub_graph_list[saga.Id]
									if ok {
										edge := createEdge(node_event_parent, node_event_child)
										edge.setStartCluster(sub_grapth_parent.getFullName())
										edge.setEndCluster(sub_graph_child.getFullName())

										g.addEdges(edge)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	g.addSubGraphs(main_graph)

	// Создаём легенду
	sub_graph_legend := createSubGraph("legend")
	sub_graph_legend.setLabel("Legend")

	sub_graph_legend_nodes := createSubGraph("legend_nodes")
	sub_graph_legend_nodes.setLabel("Nodes")

	node_unknown := createNode("node_unknown")
	node_unknown.setProperty("label", "Unknown node")

	node_created := createNode("node_created")
	node_created.setProperty("label", "Created")
	node_created.setProperty("color", "gray")
	node_created.setProperty("style", "filled")

	node_processing := createNode("node_processing")
	node_processing.setProperty("label", "Processing")
	node_processing.setProperty("color", "cornflowerblue")
	node_processing.setProperty("style", "filled")

	node_success := createNode("node_success")
	node_success.setProperty("label", "Success")
	node_success.setProperty("color", "seagreen")
	node_success.setProperty("style", "filled")

	node_error := createNode("node_error")
	node_error.setProperty("label", "Error")
	node_error.setProperty("color", "darkred")
	node_error.setProperty("style", "filled")

	node_fall_back := createNode("node_fall_back")
	node_fall_back.setProperty("label", "Fall back")
	node_fall_back.setProperty("color", "orangered")
	node_fall_back.setProperty("style", "filled")

	node_fall_back_success := createNode("node_fall_back_success")
	node_fall_back_success.setProperty("label", "Fall back success")
	node_fall_back_success.setProperty("color", "orange")
	node_fall_back_success.setProperty("style", "filled")

	node_fall_back_error := createNode("node_fall_back_error")
	node_fall_back_error.setProperty("label", "Fall back error")
	node_fall_back_error.setProperty("color", "darkorchid4")
	node_fall_back_error.setProperty("style", "filled")

	sub_graph_legend_nodes.addNodes(
		node_unknown,
		node_created,
		node_processing,
		node_success,
		node_error,
		node_fall_back,
		node_fall_back_success,
		node_fall_back_error,
	)
	g.addNodes(
		node_unknown,
		node_created,
		node_processing,
		node_success,
		node_error,
		node_fall_back,
		node_fall_back_success,
		node_fall_back_error,
	)

	//rank_legend_nodes := createRank()
	//rank_legend_nodes.addNodes(
	//	node_unknown,
	//	node_created,
	//	node_processing,
	//	node_success,
	//	node_error,
	//	node_fall_back,
	//	node_fall_back_success,
	//	node_fall_back_error,
	//)

	//g.addRanks(rank_legend_nodes)

	sub_graph_legend.addSubGraphs(sub_graph_legend_nodes)

	g.addSubGraphs(sub_graph_legend)

	g.addRanks(ranks_list...)

	return createGraphImage(folder_path, "config_"+file_name, file_name, g.print())

	//fmt.Println("begin\n")
	//

	//
	//sub_graph_operation := createSubGraph("sub_graph_operation")
	//sub_graph_operation.setLabel("Операция открытия счёта")
	//
	//sub_graph_check := createSubGraph("sub_graph_check")
	//
	//node_saga_check := createNode("node_saga_check")
	//node_saga_check.setProperty("label", "SAGA получения данных пользователя")
	//
	//sub_graph_check_events := createSubGraph("sub_graph_check_events")
	//
	//node_event_check := createNode("node_event_check")
	//node_event_check.setProperty("label", "Операция резервирования счёта")
	//
	//sub_graph_check_events.addNodes(
	//	node_event_check,
	//)
	//
	//sub_graph_check.addNodes(
	//	node_saga_check,
	//)
	//
	//sub_graph_check.addSubGraphs(
	//	sub_graph_check_events,
	//)
	//
	//sub_graph_reserve := createSubGraph("sub_graph_reserve")
	//
	//node_saga_reserve := createNode("node_saga_reserve")
	//node_saga_reserve.setProperty("label", "SAGA резервирования счёта")
	//
	//sub_graph_reserve.addNodes(
	//	node_saga_reserve,
	//)
	//
	//sub_graph_create := createSubGraph("sub_graph_create")
	//
	//node_saga_create := createNode("node_saga_create")
	//node_saga_create.setProperty("label", "SAGA создания счёта")
	//
	//sub_graph_create.addNodes(
	//	node_saga_create,
	//)
	//
	//sub_graph_open_and_add := createSubGraph("sub_graph_open_and_add")
	//
	//node_saga_open_and_add := createNode("node_saga_open_and_add")
	//node_saga_open_and_add.setProperty("label", "SAGA открытия счёта")
	//
	//sub_graph_open_and_add.addNodes(
	//	node_saga_open_and_add,
	//)
	//
	//sub_graph_operation.addSubGraphs(
	//	sub_graph_check,
	//	sub_graph_reserve,
	//	sub_graph_create,
	//	sub_graph_open_and_add,
	//)
	//
	//g.addSubGraphs(
	//	sub_graph_operation,
	//)
	//
	//g.addNodes(
	//	node_saga_check,
	//	node_saga_reserve,
	//	node_saga_create,
	//	node_saga_open_and_add,
	//	node_event_check,
	//)
	//
	//edge_saga1_saga2 := createEdge(node_saga_check, node_saga_reserve)
	//// edge_saga1_saga2.setStartCluster(sub_graph_check.getFullName())
	//// edge_saga1_saga2.setEndCluster(sub_graph_reserve.getFullName())
	//
	//edge_saga1_event1 := createEdge(node_saga_check, node_event_check)
	//edge_saga1_event1.setProperty("arrowsize", "0.0")
	//// edge_saga1_event1.setEndCluster(sub_graph_check_events.getFullName())
	//
	//edge_saga2_saga3 := createEdge(node_saga_reserve, node_saga_create)
	//edge_saga3_saga4 := createEdge(node_saga_create, node_saga_open_and_add)
	//
	//g.addEdges(
	//	edge_saga1_saga2,
	//	edge_saga2_saga3,
	//	edge_saga3_saga4,
	//	edge_saga1_event1,
	//)
	//
	//rank_saga := createRank()
	//rank_saga.addNodes(
	//	node_saga_check,
	//	node_saga_reserve,
	//	node_saga_create,
	//	node_saga_open_and_add,
	//)
	//
	//g.addRanks(
	//	rank_saga,
	//)
	//
	//g_config := g.print()
	//
	//fmt.Println(g_config)
	//
	//createGraphImage("config", "test", g_config)
	//
	//fmt.Println("end")
}
