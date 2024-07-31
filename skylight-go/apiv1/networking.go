package apiv1

import "github.com/gogf/gf/v2/frame/g"

type NetworkingGetReq struct {
	g.Meta `path:"*" tags:"Networkings" method:"get"`
}
type NetworkingGetRes struct {
	g.Meta `mime:"application/json" example:"{\"networkings\":[]}"`
}
type NetworkingPostReq struct {
	g.Meta `path:"*" tags:"Networkings" method:"post"`
}
type NetworkingPostRes struct {
	g.Meta `mime:"application/json" example:"{\"networking\":[]}"`
}
type NetworkingPutReq struct {
	g.Meta `path:"*" tags:"Networkings" method:"post"`
}
type NetworkingPutRes struct {
	g.Meta `mime:"application/json" example:"{\"networking\":[]}"`
}
type NetworkingDeleteReq struct {
	g.Meta `path:"*" tags:"Networkings" method:"delete"`
}
type NetworkingDeleteRes struct {
	g.Meta `mime:"application/json"`
}
