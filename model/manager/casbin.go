package manager

type CasbinInfo struct {
	RoleID string `json:"roleid,omitempty" bson:"v0"` //角色id，标识为角色，要设置忽略在 "22"
	Path   string `json:"path" bson:"v1"`
	Method string `json:"method" bson:"v2"`
}

type CasbinReq struct {
	RoleID      string       `json:"roleid"` // 权限id
	CasbinInfos []CasbinInfo `json:"casbinInfos"`
}

func DefaultCasbin() CasbinReq {
	return CasbinReq{
		RoleID: "222", //"222"是超级管理员 角色标识
		CasbinInfos: []CasbinInfo{
			{Path: "/pub/healthz", Method: "GET"},
			{Path: "/pub/login", Method: "POST"},
		}}
}
