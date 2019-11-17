package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ben-han-cn/gorest"
	"github.com/ben-han-cn/gorest/adaptor"
	goresterr "github.com/ben-han-cn/gorest/error"
	"github.com/ben-han-cn/gorest/resource"
	"github.com/ben-han-cn/gorest/resource/schema"
)

var (
	version = resource.APIVersion{
		Group:   "zdns.cloud.example",
		Version: "example/v1",
	}
	clusterKind = resource.DefaultKindName(Cluster{})
	nodeKind    = resource.DefaultKindName(Node{})
)

type Cluster struct {
	resource.ResourceBase `json:",inline"`
	Name                  string         `json:"name" rest:"required=true,minLen=1,maxLen=10"`
	NodeCount             int            `json:"nodeCount" rest:"required=true,min=1,max=1000"`
	MapData               map[string]int `json:"mapData" rest:"required=true"`

	nodes []*Node `json:"-"`
}

type Node struct {
	resource.ResourceBase `json:",inline"`
	Address               string `json:"address,omitempty" rest:"required=true,minLen=7,maxLen=13"`
	IsWorker              bool   `json:"isWorker"`
}

func (c Cluster) CreateActions(name string) *resource.Action {
	switch name {
	case "encode":
		return &resource.Action{
			Name:  "encode",
			Input: &Input{},
		}
	case "decode":
		return &resource.Action{
			Name:  "decode",
			Input: &Input{},
		}
	default:
		return nil
	}
}

type Input struct {
	Data string `json:"data,omitempty"`
}

func (n Node) GetParents() []resource.ResourceKind {
	return []resource.ResourceKind{Cluster{}}
}

func (n Node) CreateDefaultResource() resource.Resource {
	return &Node{
		IsWorker: true,
	}
}

type State struct {
	clusters []*Cluster
	lock     sync.Mutex
}

func newState() *State {
	return &State{}
}

func (s *State) AddCluster(cluster *Cluster) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if c := s.getCluster(cluster.Name); c != nil {
		return fmt.Errorf("cluster %s already exist", cluster.Name)
	}
	s.clusters = append(s.clusters, cluster)
	return nil
}

func (s *State) GetCluster(name string) *Cluster {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.getCluster(name)
}

func (s *State) getCluster(name string) *Cluster {
	for _, c := range s.clusters {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (s *State) GetClusters() []*Cluster {
	s.lock.Lock()
	defer s.lock.Unlock()

	cl := make([]*Cluster, len(s.clusters))
	copy(cl, s.clusters)
	return cl
}

func (s *State) AddNode(clustreName string, node *Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	c := s.getCluster(clustreName)
	if c == nil {
		return fmt.Errorf("unknown cluster with name %s", clustreName)
	}

	if n := s.getNode(c, node.Address); n != nil {
		return fmt.Errorf("add duplicate node with address %s", node.Address)
	}
	c.nodes = append(c.nodes, node)
	return nil
}

func (s *State) GetNode(clusterName string, address string) *Node {
	s.lock.Lock()
	defer s.lock.Unlock()

	c := s.getCluster(clusterName)
	if c == nil {
		return nil
	}

	return s.getNode(c, address)
}

func (s *State) getNode(c *Cluster, address string) *Node {
	for _, n := range c.nodes {
		if n.Address == address {
			return n
		}
	}
	return nil
}

func (s *State) GetNodes(clusterName string) []*Node {
	s.lock.Lock()
	defer s.lock.Unlock()

	c := s.getCluster(clusterName)
	return s.getNodes(c)
}

func (s *State) getNodes(c *Cluster) []*Node {
	nodes := make([]*Node, len(c.nodes))
	copy(nodes, c.nodes)
	return nodes
}

func (s *State) DeleteNode(clusterName string, address string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	c := s.getCluster(clusterName)
	if c == nil {
		return fmt.Errorf("unknwon cluster with name %s", clusterName)
	}

	return s.deleteNode(c, address)
}

func (s *State) deleteNode(c *Cluster, address string) error {
	for i, n := range c.nodes {
		if n.Address == address {
			c.nodes = append(c.nodes[:i], c.nodes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("unknown node with address %s", address)
}

type clusterHandler struct {
	clusters *State
}

func newClusterHandler(s *State) *clusterHandler {
	return &clusterHandler{
		clusters: s,
	}
}

func (h *clusterHandler) Create(ctx *resource.Context) (resource.Resource, *goresterr.APIError) {
	cluster := ctx.Resource.(*Cluster)
	cluster.SetID(cluster.Name)
	cluster.SetCreationTimestamp(time.Now())
	if err := h.clusters.AddCluster(cluster); err != nil {
		return nil, goresterr.NewAPIError(goresterr.DuplicateResource, err.Error())
	} else {
		return cluster, nil
	}
}

func (h *clusterHandler) List(ctx *resource.Context) interface{} {
	//return []int{1, 2, 3}
	return 1
	//return h.clusters.GetClusters()
}

func (h *clusterHandler) Get(ctx *resource.Context) resource.Resource {
	return h.clusters.GetCluster(ctx.Resource.GetID())
}

func (h *clusterHandler) Action(ctx *resource.Context) (interface{}, *goresterr.APIError) {
	r := ctx.Resource
	input, _ := r.GetAction().Input.(*Input)
	switch r.GetAction().Name {
	case "encode":
		return base64.StdEncoding.EncodeToString([]byte(input.Data)), nil
	case "decode":
		if data, e := base64.StdEncoding.DecodeString(input.Data); e != nil {
			return nil, goresterr.NewAPIError(goresterr.InvalidFormat, e.Error())
		} else {
			return string(data), nil
		}
	default:
		panic("it should never come here")
	}
}

type nodeHandler struct {
	clusters *State
}

func newNodeHandler(s *State) *nodeHandler {
	return &nodeHandler{
		clusters: s,
	}
}

func (h *nodeHandler) Create(ctx *resource.Context) (resource.Resource, *goresterr.APIError) {
	node := ctx.Resource.(*Node)
	if ip := net.ParseIP(node.Address); ip == nil {
		return nil, goresterr.NewAPIError(goresterr.InvalidFormat, "address isn't valid ipv4 address")
	}

	node.SetID(node.Address)
	if err := h.clusters.AddNode(node.GetParent().GetID(), node); err != nil {
		return nil, goresterr.NewAPIError(goresterr.NotFound, err.Error())
	}
	return node, nil
}

func (h *nodeHandler) Delete(ctx *resource.Context) *goresterr.APIError {
	node := ctx.Resource.(*Node)
	if err := h.clusters.DeleteNode(node.GetParent().GetID(), node.GetID()); err != nil {
		return goresterr.NewAPIError(goresterr.NotFound, err.Error())
	} else {
		return nil
	}
}

func (h *nodeHandler) List(ctx *resource.Context) interface{} {
	node := ctx.Resource.(*Node)
	return h.clusters.GetNodes(node.GetParent().GetID())
}

func (h *nodeHandler) Get(ctx *resource.Context) interface{} {
	node := ctx.Resource.(*Node)
	return h.clusters.GetNode(node.GetParent().GetID(), node.GetID())
}

func main() {
	schemas := schema.NewSchemaManager()
	state := newState()
	schemas.Import(&version, Cluster{}, newClusterHandler(state))
	schemas.Import(&version, Node{}, newNodeHandler(state))
	router := gin.Default()
	adaptor.RegisterHandler(router, gorest.NewAPIServer(schemas), schemas.GenerateResourceRoute())
	router.Run("0.0.0.0:1234")
}
