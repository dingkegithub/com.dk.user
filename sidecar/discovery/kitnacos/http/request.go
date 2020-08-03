package http

type RegisterInstanceRequest struct {
	Ip          string  `json:"ip"`
	Port        uint16  `json:"port"`
	NamespaceId string  `json:"namespaceId"`
	Weight      float64 `json:"weight"`
	Enabled     bool    `json:"enabled"`
	Healthy     bool    `json:"healthy"`
	Metadata    string  `json:"metadata"`
	ClusterName string  `json:"clusterName"`
	ServiceName string  `json:"serviceName"`
	GroupName   string  `json:"groupName"`
	Ephemeral   bool    `json:"ephemeral"`
}

type DeregisterInstanceRequest struct {
	Ip          string `json:"ip"`
	Port        uint16 `json:"port"`
	NamespaceId string `json:"namespaceId"`
	ClusterName string `json:"clusterName"`
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	Ephemeral   bool   `json:"ephemeral"`
}

type ModifyInstanceRequest struct {
	Ip          string  `json:"ip"`
	Port        int     `json:"port"`
	NamespaceId string  `json:"namespaceId"`
	Weight      float64 `json:"weight"`
	Enabled     bool    `json:"enabled"`
	Metadata    string  `json:"metadata"`
	ClusterName string  `json:"clusterName"`
	ServiceName string  `json:"serviceName"`
	GroupName   string  `json:"groupName"`
	Ephemeral   bool    `json:"ephemeral"`
}

type ListInstanceRequest struct {
	NamespaceId string `json:"namespaceId"`
	ClusterName string `json:"clusterName"`
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	HealthyOnly bool   `json:"healthyOnly"`
}

type HostInfo struct {
	Valid      bool                   `json:"valid"`
	Marked     bool                   `json:"marked"`
	InstanceId string                 `json:"instanceId"`
	Port       uint16                 `json:"port"`
	Ip         string                 `json:"ip"`
	Weight     float64                `json:"weight"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type ListInstanceResponse struct {
	Dom             string      `json:"dom"`
	CacheMillis     uint64      `json:"cacheMillis"`
	UseSpecifiedURL bool        `json:"useSpecifiedURL"`
	Hosts           []*HostInfo `json:"hosts"`
	Checksum        string      `json:"checksum"`
	LastRefTime     uint64      `json:"lastRefTime"`
	Env             string      `json:"env"`
	Clusters        string      `json:"clusters"`
}

type DetailInstanceRequest struct {
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	NamespaceId string `json:"namespaceId"`
	ClusterName string `json:"clusterName"`
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	Ephemeral   bool   `json:"ephemeral"`
	HealthyOnly bool   `json:"healthyOnly"`
}

type DetailInstanceResponse struct {
	Metadata    map[string]interface{} `json:"metadata"`
	InstanceId  string                 `json:"instanceId"`
	Port        int                    `json:"port"`
	Service     string                 `json:"service"`
	Healthy     bool                   `json:"healthy"`
	Ip          string                 `json:"ip"`
	ClusterName string                 `json:"clusterName"`
	Weight      float64                `json:"weight"`
}

type Beat struct {
	Cluster     string                 `json:"cluster"`
	Ip          string                 `json:"ip"`
	Metadata    map[string]interface{} `json:"metadata"`
	Port        uint16                 `json:"port"`
	Scheduled   bool                   `json:"scheduled"`
	ServiceName string                 `json:"serviceName"`
	Weight      float64                `json:"weight"`
}

type HeartbeatRequest struct {
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	Ephemeral   bool   `json:"ephemeral"`
	Beat        *Beat   `json:"beat"`
}

type ServiceCreateRequest struct {
	ServiceName      string  `json:"serviceName"`
	GroupName        string  `json:"groupName"`
	NamespaceId      string  `json:"namespaceId"`
	ProtectThreshold float64 `json:"protectThreshold"`
	Metadata         string  `json:"metadata"`
	Selector         string  `json:"selector"`
}

type ServiceDeleteRequest struct {
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	NamespaceId string `json:"namespaceId"`
}

type ServiceModifyRequest struct {
	ServiceName      string  `json:"serviceName"`
	GroupName        string  `json:"groupName"`
	NamespaceId      string  `json:"namespaceId"`
	ProtectThreshold float64 `json:"protectThreshold"`
	Metadata         string  `json:"metadata"`
	Selector         string  `json:"selector"`
}

type ServiceRetrieve struct {
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	NamespaceId string `json:"namespaceId"`
}

type Selector struct {
	Type string `json:"type"`
}

type HealthChecker struct {
	Type string `json:"type"`
}

type Cluster struct {
	HealthChecker HealthChecker          `json:"healthChecker"`
	Metadata      map[string]interface{} `json:"metadata"`
	Name          string                 `json:"name"`
}

type ServiceRetrieveResponse struct {
	Metadata         map[string]interface{} `json:"metadata"`
	GroupName        string                 `json:"groupName"`
	NamespaceId      string                 `json:"namespaceId"`
	Name             string                 `json:"name"`
	Selector         Selector               `json:"selector"`
	ProtectThreshold float64                `json:"protectThreshold"`
	Clusters         Cluster                `json:"clusters"`
}

type ServiceQueryListRequest struct {
	PageNo      int    `json:"pageNo"`
	PageSize    int    `json:"pageSize"`
	GroupName   string `json:"groupName"`
	NamespaceId string `json:"namespaceId"`
}

type ServiceQueryListResponse struct {
	Count int      `json:"count"`
	Doms  []string `json:"doms"`
}

type ClusterServer struct {
	Ip             string  `json:"ip"`
	ServePort      uint16  `json:"servePort"`
	Site           string  `json:"site"`
	Weight         float64 `json:"weight"`
	AdWeight       float64 `json:"adWeight"`
	Alive          bool    `json:"alive"`
	LastRefTime    uint64  `json:"lastRefTime"`
	LastRefTimeStr string  `json:"lastRefTimeStr"`
	Key            string  `json:"key"`
}

type ClusterServersQueryRequest struct {
	Healthy bool `json:"healthy"`
}

type ClusterServersQueryResponse struct {
	Servers []*ClusterServer `json:"servers"`
}

type Leader struct {
	HeartbeatDueMs uint64 `json:"heartbeatDueMs"`
	Ip             string `json:"ip"`
	LeaderDueMs    uint64 `json:"leaderDueMs"`
	State          string `json:"state"`
	Term           uint64 `json:"term"`
	VoteFor        string `json:"voteFor"`
}

type ClusterLeaderResponse struct {
	Leader Leader `json:"leader"`
}
