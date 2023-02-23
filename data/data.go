package data

const (
	//0:app_dn 1:net_dn 2:net_up 3:lk_up 4:lk_dn 5:ph_up 6:app_in
	App_dn = iota
	Net_dn
	Net_up
	Lk_up
	Lk_dn
	Ph_up
	App_in
)

type Net struct {
	Node_num int                 //节点数
	Graph    map[string][]string //图
}

//这里特别注意端口和位置的一一对应，0-左，1-上，2-右，3-下
//一个节点最多四个端口
//拓扑管道映射表
type Chmap struct {
	//node_num int
	Ch_tp    map[string][4]chan []byte
	Port_map map[string][]Edge
}

//一条网络边，端口对应信息
type Edge struct {
	My_port   int
	Your_port int
	Ch_tp     chan []byte //注意！！！是出节点端口通道
	S_Dst     string
	//Port_Dst int
}

//节点信息
type Node struct {
	//节点字母
	S string
	//socket 通道
	//ch_so chan []byte
	//拓扑通道
	//num_ch_tp int
	//ch_tp     []chan []byte
	//节点通道(进节点)
	Ch_app_in chan Msg_app
	Ch_ph_up  chan Msg_lk
	Ch_lk_dn  chan Msg_phy
	Ch_lk_up  chan Msg_net
	Ch_net_up chan Msg_app
	Ch_net_dn chan Msg_lk
	Ch_app_dn chan Msg_net
	//Ch_source chan []byte
}

//socket端口管道映射表
type List_so_ch struct {
	So_port []int
	S       []string
	Ch_so   chan Msg_so_send
}

//应用层报文
type Msg_app struct {
	Msg []byte
}

//网络层报文
type Msg_net struct {
	Msg []byte
}

//链路层报文

type Msg_lk struct {
	Msg []byte
}

//物理层报文

type Msg_phy struct {
	S_Dst string
	Msg   []byte
}

//socket发送报文

type Msg_so_send struct {
	S        string
	LayerSrc int
	Len      int
	Msg      []byte
}

//socket接收报文
type Msg_so_rev struct {
	S        string
	LayerDst int
	S_Dst    string
	Len      int
	Msg      []byte
}

//实验信息结构
type Msg_ex struct {
	Type_num int //1:string
	Msg      []byte
}

func (msg *Msg_ex) String(s string) {
	msg.Type_num = 1
	msg.Msg = []byte(s)
}
