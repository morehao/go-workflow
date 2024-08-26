package objProcDef

type ProcDefHistoryBaseInfo struct {
	Company    string `json:"company" form:"company"`       //
	DeployTime string `json:"deployTime" form:"deployTime"` //
	Name       string `json:"name" form:"name"`             //
	Resource   string `json:"resource" form:"resource"`     //
	UserID     string `json:"userid" form:"userid"`         //
	Username   string `json:"username" form:"username"`     //
	Version    int32  `json:"version" form:"version"`       //
}
