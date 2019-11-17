# Resource以及ResourceKind
- 父资源是不是一个资源
- 资源支持的操作是由handler决定的
- resource type和url
- 资源类型和资源区分开

```go
func SetPodSchema(schema *resttypes.Schema, handler resttypes.Handler) {
    schema.Handler = handler
    schema.CollectionMethods = []string{"GET"}
    schema.ResourceMethods = []string{"GET", "DELETE"}
    schema.Parents = []string{DeploymentType, DaemonSetType, StatefulSetType, JobType, CronJobType}
}
```
# Handler接口
- Crate/Update应该返回resource
- List接口保证返回的是resource的list

# Required和正确性验证
- 要在json数据unmarshal到go对象之前做required的检查
- default可以用源编程但是通过kind返回对象更好（map是问题）
- 正确性验证支持子结构体

# route的性能
- 资源是从bottom到top
- route的资源查找是从top到bottom

```go
type Node struct {
    resource.ResourceBase `json:",inline"`
    Address               string `json:"address,omitempty" rest:"required=true,minLen=7,maxLen=13"`
    IsWorker              bool   `json:"isWorker"`
}

func (n Node) GetParents() []resource.ResourceKind {
    return []resource.ResourceKind{Cluster{}}
}

func (n Node) CreateDefaultResource() resource.Resource {
    return &Node{
        IsWorker: true,
    }   
}
schemas.Import(&version, Cluster{}, newClusterHandler(state))
```
