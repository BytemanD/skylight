package apiv1

import "github.com/gogf/gf/v2/frame/g"

type ComputingGetReq struct {
	g.Meta `path:"*" tags:"Computings" method:"get"`
}
type ComputingGetRes struct {
	g.Meta `mime:"application/json" example:"{\"networkings\":[]}"`
}
type ComputingPostReq struct {
	g.Meta `path:"*" tags:"Computings" method:"post"`
}
type ComputingPostRes struct {
	g.Meta `mime:"application/json" example:"{\"networking\":[]}"`
}
type ComputingPutReq struct {
	g.Meta `path:"*" tags:"Computings" method:"post"`
}
type ComputingPutRes struct {
	g.Meta `mime:"application/json" example:"{\"networking\":[]}"`
}
type ComputingDeleteReq struct {
	g.Meta `path:"*" tags:"Computings" method:"delete"`
}
type ComputingDeleteRes struct {
	g.Meta `mime:"application/json"`
}
