package model

//Meta defines the struct for meta info that is
//replicated across each Raft Group member
type Meta struct {
	ClusterID    string   `json:"cluster_id"`
	Members      []string `json:"members"`
	DatabaseName string   `json:"db_name"`
	Collection   []string `json:"collection"`
	Partition    string   `json:"partition"`
	Replication  string   `json:"replication"`
	Namespace    []string `json:"namespace"`
}
