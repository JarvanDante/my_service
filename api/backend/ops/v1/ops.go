// Package v1 后台运营配置契约(公告/跳转位/敏感词管理)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------- 公告 ----------

type AnnItem struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	TextNode  string `json:"text_node"`
	Cover     string `json:"cover"`
	JumpUrl   string `json:"jump_url"`
	SysType   string `json:"sys_type"`
	StartAt   string `json:"start_at"`
	EndAt     string `json:"end_at"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type AnnListReq struct {
	g.Meta `path:"/announcement" method:"get" tags:"Backend/Ops" summary:"公告列表"`
	Status string `json:"status"` // 空=全部  0=关闭  1=开启
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type AnnListRes struct {
	List  []AnnItem `json:"list"`
	Total int       `json:"total"`
}

type AnnCreateReq struct {
	g.Meta   `path:"/announcement" method:"post" tags:"Backend/Ops" summary:"新增公告"`
	Title    string `json:"title" v:"required#标题必填"`
	Content  string `json:"content"`
	TextNode string `json:"text_node"`
	Cover    string `json:"cover"`
	JumpUrl  string `json:"jump_url"`
	SysType  string `json:"sys_type"`                   // 默认 app
	StartAt  string `json:"start_at"`                   // 默认 now
	EndAt    string `json:"end_at" v:"required#结束时间必填"` // 如 2027-12-31 23:59:59
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type AnnCreateRes struct {
	Id int64 `json:"id"`
}

type AnnUpdateReq struct {
	g.Meta   `path:"/announcement/{id}" method:"put" tags:"Backend/Ops" summary:"更新公告"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	TextNode string `json:"text_node"`
	Cover    string `json:"cover"`
	JumpUrl  string `json:"jump_url"`
	SysType  string `json:"sys_type"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type AnnUpdateRes struct{}

type AnnDeleteReq struct {
	g.Meta `path:"/announcement/{id}" method:"delete" tags:"Backend/Ops" summary:"删除公告"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type AnnDeleteRes struct{}

// ---------- 跳转位 ----------

type JtItem struct {
	Id          int64  `json:"id"`
	CnName      string `json:"cn_name"`
	EnName      string `json:"en_name"`
	Avatar      string `json:"avatar"`
	Link        string `json:"link"`
	PicJumpLink string `json:"pic_jump_link"`
	Location    int    `json:"location"`
	Rank        int    `json:"rank"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type JtListReq struct {
	g.Meta   `path:"/jumptab" method:"get" tags:"Backend/Ops" summary:"跳转位列表"`
	Location int    `json:"location"` // 0=全部
	Status   string `json:"status"`   // 空=全部  0=禁用  1=启用
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type JtListRes struct {
	List  []JtItem `json:"list"`
	Total int      `json:"total"`
}

type JtCreateReq struct {
	g.Meta      `path:"/jumptab" method:"post" tags:"Backend/Ops" summary:"新增跳转位"`
	CnName      string `json:"cn_name" v:"required#名称必填"`
	EnName      string `json:"en_name"`
	Avatar      string `json:"avatar"`
	Link        string `json:"link"`
	PicJumpLink string `json:"pic_jump_link"`
	Location    int    `json:"location" v:"required|min:1#位置必填|位置非法"`
	Rank        int    `json:"rank"`
	Status      int    `json:"status" v:"in:0,1#状态非法"`
}
type JtCreateRes struct {
	Id int64 `json:"id"`
}

type JtUpdateReq struct {
	g.Meta      `path:"/jumptab/{id}" method:"put" tags:"Backend/Ops" summary:"更新跳转位"`
	Id          int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	CnName      string `json:"cn_name"`
	EnName      string `json:"en_name"`
	Avatar      string `json:"avatar"`
	Link        string `json:"link"`
	PicJumpLink string `json:"pic_jump_link"`
	Location    int    `json:"location"`
	Rank        int    `json:"rank"`
	Status      int    `json:"status" v:"in:0,1#状态非法"`
}
type JtUpdateRes struct{}

type JtDeleteReq struct {
	g.Meta `path:"/jumptab/{id}" method:"delete" tags:"Backend/Ops" summary:"删除跳转位"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type JtDeleteRes struct{}

// ---------- 敏感词 ----------

type FwItem struct {
	Id        int64  `json:"id"`
	Word      string `json:"word"`
	CreatedAt string `json:"created_at"`
}

type FwListReq struct {
	g.Meta  `path:"/filterword" method:"get" tags:"Backend/Ops" summary:"敏感词列表"`
	Keyword string `json:"keyword"` // 模糊搜索
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type FwListRes struct {
	List  []FwItem `json:"list"`
	Total int      `json:"total"`
}

type FwAddReq struct {
	g.Meta `path:"/filterword" method:"post" tags:"Backend/Ops" summary:"批量添加敏感词(重复自动跳过)"`
	Words  []string `json:"words" v:"required#词列表必填"`
}
type FwAddRes struct {
	Added int `json:"added"` // 实际新增数(重复跳过)
}

type FwDeleteReq struct {
	g.Meta `path:"/filterword/{id}" method:"delete" tags:"Backend/Ops" summary:"删除敏感词"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type FwDeleteRes struct{}
