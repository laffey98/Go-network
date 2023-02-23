package main

import (
	"fmt"
	. "network/data"
	pr "network/print"
)

//待完成：
//assemble函数
//Start_socket_server()
//

func main() {

	//初始化实验管道
	exch := make(chan Msg_ex)

	//接收Net变量
	//test
	var net1 Net
	net1.Node_num = 4
	net1.Graph = make(map[string][]string, 4)
	net1.Graph["A"] = []string{"B"}
	net1.Graph["B"] = []string{"A", "C", "D"}
	net1.Graph["C"] = []string{"B"}
	net1.Graph["D"] = []string{"B"}
	//test

	//初始化并完善部分拓扑管道映射表
	var chmap Chmap
	chmap.Ch_tp = make(map[string][4]chan []byte, net1.Node_num)
	chmap.Port_map = make(map[string][]Edge, net1.Node_num)
	for key := range net1.Graph {

		//fmt.Println(key, v)

		var charr [4]chan []byte
		for i := 0; i < 4; i++ {
			charr[i] = make(chan []byte)
		}
		chmap.Ch_tp[key] = charr

		var edgearr = make([]Edge, len(net1.Graph[key]))

		//fmt.Println("len(net1.Graph[key]):", len(net1.Graph[key]))

		for i, s := range net1.Graph[key] {
			edgearr[i].My_port = i
			edgearr[i].S_Dst = s
			//edgearr[i].Ch_tp = make(chan []byte)

			//fmt.Println(i, s)

		}
		chmap.Port_map[key] = edgearr
	}
	pr.Printvar("CreateMsg", "chmap has created and completed partly. chmap:", chmap)
	//fmt.Println(chmap)

	//初始化并完善节点信息
	node := make([]Node, net1.Node_num)
	i := 0
	for key := range net1.Graph {
		node[i].S = key
		node[i].Ch_lk_dn = make(chan Msg_phy)
		node[i].Ch_lk_up = make(chan Msg_net)
		node[i].Ch_net_dn = make(chan Msg_lk)
		node[i].Ch_net_up = make(chan Msg_app)
		node[i].Ch_ph_up = make(chan Msg_lk)
		node[i].Ch_app_dn = make(chan Msg_net)
		node[i].Ch_app_in = make(chan Msg_app)
		i++
	}
	pr.Printvar("CreateMsg", "node array has created and completed totally. node array:", node)

	//fmt.Println(node)
	/* 	chc := make(chan []byte)
	   	ccc := chc
	   	fmt.Println(chc, ccc) */

	//完善所有拓扑管道映射表
	for source, v := range chmap.Port_map {
		for i, src_edge := range v {
			for _, dst_edge := range chmap.Port_map[src_edge.S_Dst] {
				if dst_edge.S_Dst == source {
					//fmt.Println(1)
					//fmt.Println(src_edge.Ch_tp)
					chmap.Port_map[source][i].Ch_tp = chmap.Ch_tp[src_edge.S_Dst][dst_edge.My_port]
					//src_edge.Ch_tp = chmap.Ch_tp[src_edge.S_Dst][dst_edge.My_port]
					//fmt.Println(src_edge.Ch_tp)
				}
			}
		}
	}
	pr.Printvar("CreateMsg", "chmap array has completed totally. chmap:", chmap)

	//初始化并完善部分so管道映射表
	var list_so_ch List_so_ch
	list_so_ch.Ch_so = make(chan Msg_so_send)
	s := make([]string, net1.Node_num)
	i = 0
	for key := range net1.Graph {
		s[i] = key
		i++
	}
	list_so_ch.S = s

	//fmt.Println(list_so_ch)
	pr.Printvar("CreateMsg", "list_so_ch has created and completed partly. list_so_ch:", list_so_ch)

	//启动服务器协程
	go Start_socket_server(list_so_ch, exch)

	//启动websocket
	if ok := StartWS(); !ok {
		pr.Printvar("WS", "StartWS false")
		var msg Msg_ex
		p := &msg
		p.String("StartWS false")
		exch <- msg
	}

	//启动网络创建，大量协程展开！！！
	if ok := StartNet(chmap, node, exch, list_so_ch.Ch_so); !ok {
		pr.Printvar("CreateMsg", "StartPhy false")
		var msg Msg_ex
		p := &msg
		p.String("StartPhy false")
		exch <- msg
	}

	//websocket循环处理
	for {
		msg := <-exch
		fmt.Println(msg)
	}

}

func StartWS() bool {
	fmt.Println("StartWS sth")
	return true
}

//创建net，拓扑ch，节点ch，与socket通信的协程
func StartNet(chmap Chmap, node []Node, exch chan Msg_ex, ch_so chan Msg_so_send) bool {
	for _, n := range node {
		go Start_phy_layer(chmap, n, exch)
		go Start_lk_layer(n, exch, ch_so)
		go Start_net_layer(n, exch, ch_so)
		go Start_app_layer(n, exch, ch_so)
	}
	fmt.Println("StartNet sth")
	return true
}
func Start_socket_server(list_so_ch List_so_ch, exch chan Msg_ex) {
	ok := true
	//开启socket服务器，连接节点数量的客户端，写socket管道映射表，开协程 并等待消息，筛选后发送至对应管道
	//主协程用socket管道映射表，select等待管道，并对应发送
	if !ok {
		pr.Printvar("Start_socket_server", "Start_socket_server false")
		var msg Msg_ex
		p := &msg
		p.String("Start_socket_server false")
		exch <- msg
	} else {
		pr.Printvar("Start_socket_server", "Start_socket_server success")
	}
}

//以下四个函数是抽象的层实体（select监听所有拓扑通道，所有节点通道，并做相应处理）
func Start_phy_layer(chmap Chmap, n Node, exch chan Msg_ex) {
	pr.Printvar("Start_phy_layer", "Start phy_layer of", n.S)
	for {
		select {
		case msg := <-chmap.Ch_tp[n.S][0]:
			phy_operate(msg)
			var next Msg_lk
			next.Msg = msg
			n.Ch_ph_up <- next
		case msg := <-chmap.Ch_tp[n.S][1]:
			phy_operate(msg)
			var next Msg_lk
			next.Msg = msg
			n.Ch_ph_up <- next
		case msg := <-chmap.Ch_tp[n.S][2]:
			phy_operate(msg)
			var next Msg_lk
			next.Msg = msg
			n.Ch_ph_up <- next
		case msg := <-chmap.Ch_tp[n.S][3]:
			phy_operate(msg)
			var next Msg_lk
			next.Msg = msg
			n.Ch_ph_up <- next
		case msg := <-n.Ch_lk_dn:
			for _, edge := range chmap.Port_map[n.S] {
				if msg.S_Dst == edge.S_Dst {
					edge.Ch_tp <- msg.Msg
				}
			}
		}
	}
}
func Start_lk_layer(n Node, exch chan Msg_ex, ch_so chan Msg_so_send) {
	pr.Printvar("Start_lk_layer", "Start lk_layer of", n.S)
	select {
	case msg := <-n.Ch_ph_up:
		msg_so_send := Msg_so_send{S: n.S, LayerSrc: Ph_up, Len: len(msg.Msg), Msg: msg.Msg}
		ch_so <- msg_so_send
	case msg := <-n.Ch_net_dn:
		msg_so_send := Msg_so_send{S: n.S, LayerSrc: Net_dn, Len: len(msg.Msg), Msg: msg.Msg}
		ch_so <- msg_so_send
	}
}
func Start_net_layer(n Node, exch chan Msg_ex, ch_so chan Msg_so_send) {
	pr.Printvar("Start_net_layer", "Start net_layer of", n.S)
	select {
	case msg := <-n.Ch_lk_up:
		msg_so_send := Msg_so_send{S: n.S, LayerSrc: Lk_up, Len: len(msg.Msg), Msg: msg.Msg}
		ch_so <- msg_so_send
	case msg := <-n.Ch_app_dn:
		msg_so_send := Msg_so_send{S: n.S, LayerSrc: App_dn, Len: len(msg.Msg), Msg: msg.Msg}
		ch_so <- msg_so_send
	}
}
func Start_app_layer(n Node, exch chan Msg_ex, ch_so chan Msg_so_send) {
	pr.Printvar("Start_app_layer", "Start app_layer of", n.S)
	select {
	case msg := <-n.Ch_net_up:
		assemble(msg)
	case msg := <-n.Ch_app_in:
		msg_so_send := Msg_so_send{S: n.S, LayerSrc: App_in, Len: len(msg.Msg), Msg: msg.Msg}
		ch_so <- msg_so_send
	}
}

//对物理层的byte数组做物理层仿真处理
func phy_operate(msg []byte) {

	//调用phy_operate包内的函数
	fmt.Println("phy_operate sth")
}

//到站组装
func assemble(msg Msg_app) {}
