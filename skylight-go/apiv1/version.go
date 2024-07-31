package apiv1

import "github.com/gogf/gf/v2/frame/g"

type VersionGetReq struct {
	g.Meta `path:"/version" tags:"configmaps" method:"get"`
}
type VersionGetRes struct {
	g.Meta `mime:"application/json" example:"{\"version\": "xx"}"`
}
