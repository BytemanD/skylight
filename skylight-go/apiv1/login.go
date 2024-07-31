package apiv1

import "github.com/gogf/gf/v2/frame/g"

type LoginPostReq struct {
	g.Meta `path:"/login" tags:"configmaps" method:"post"`
}
type LoginPostRes struct {
	g.Meta `mime:"application/json" example:"{}"`
}
type LoginGetReq struct {
	g.Meta `path:"/login" tags:"configmaps" method:"get"`
}
type LoginGetRes struct {
	g.Meta `mime:"application/json" example:"{}"`
}
