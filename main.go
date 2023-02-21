package main

import (
	"fmt"
	. "network/data"
	pr "network/print"
)

//待完成：
//使用pr打印信息

func main() {

	//初始化实验管道
	exch := make(chan string)

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

	//初始化并完善拓扑管道映射表
	var chmap Chmap
	chmap.Ch_tp = make(map[string][4]chan []byte, net1.Node_num)
	chmap.Port_map = make(map[string][]Edge, net1.Node_num)
	for key, _ := range net1.Graph {

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

			//fmt.Println(i, s)

		}
		chmap.Port_map[key] = edgearr

	}
	//fmt.Println(chmap)

	//初始化并完善节点信息
	node := make([]Node, net1.Node_num)
	i := 0
	for key, _ := range net1.Graph {
		node[i].S = key
		node[i].Ch_lk_dn = make(chan []byte)
		node[i].Ch_lk_up = make(chan []byte)
		node[i].Ch_net_dn = make(chan []byte)
		node[i].Ch_net_up = make(chan []byte)
		node[i].Ch_ph_up = make(chan []byte)
		node[i].Ch_app_dn = make(chan []byte)
		i++
	}

	fmt.Println(node)

	//初始化并完善部分so管道映射表
	var list_so_ch List_so_ch
	list_so_ch.Ch_so = make(chan Msg_so_send)
	s := make([]string, net1.Node_num)
	i = 0
	for key, _ := range net1.Graph {
		s[i] = key
		i++
	}
	list_so_ch.S = s

	//fmt.Println(list_so_ch)

	//启动服务器协程
	go Start_socket_server(list_so_ch, exch)

	//启动websocket
	if ok := StartWS(); !ok {
		pr.Printvar("WS", "StartWS false")
		exch <- "StartWS false"
	}

	//启动网络创建，大量协程展开！！！
	if ok := StartNet(chmap, node, exch, list_so_ch.Ch_so); !ok {
		pr.Printvar("main", "StartPhy false")
		exch <- "StartPhy false"
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
func StartNet(chmap Chmap, node []Node, exch chan string, ch_so chan Msg_so_send) bool {
	fmt.Println("StartNet sth")
	return true
}
func Start_socket_server(list_so_ch List_so_ch, exch chan string) {
	ok := true
	//开启socket服务器，连接节点数量的客户端，写socket管道映射表，开协程 并等待消息，筛选后发送至对应管道
	//主协程用socket管道映射表，select等待管道，并对应发送
	if !ok {
		pr.Printvar("Start_socket_server", "Start_socket_server false")
		exch <- "Start_socket_server false"
	} else {
		pr.Printvar("Start_socket_server", "Start_socket_server success")
	}
}
