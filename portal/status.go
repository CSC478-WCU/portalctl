package portal

import "encoding/json"

// StatusNode is a single node's published info (subset).
type StatusNode struct {
	Hostname string `json:"hostname"`
	IPv4     string `json:"ipv4"`
}

// StatusAggregate collects per-aggregate status and nodes.
type StatusAggregate struct {
	Status string                `json:"status"`
	Nodes  map[string]StatusNode `json:"nodes"`
}

// StatusPayload matches experimentStatus -j JSON payload.
type StatusPayload struct {
	UUID               string                       `json:"uuid"`
	URL                string                       `json:"url"`
	Expires            string                       `json:"expires"`
	Status             string                       `json:"status"`
	AggregateStatus    map[string]StatusAggregate   `json:"aggregate_status"`
	InstanceCertificate string                      `json:"instance_certificate,omitempty"`
	InstancePrivateKey  string                      `json:"instance_private_key,omitempty"`
}

// ParseStatusJSON decodes the server's JSON string.
func ParseStatusJSON(s string) (*StatusPayload, error) {
	var p StatusPayload
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// FlattenNodes flattens all aggregates to one map of client_id->node.
func FlattenNodes(p *StatusPayload) map[string]StatusNode {
	out := make(map[string]StatusNode)
	for _, agg := range p.AggregateStatus {
		for cid, n := range agg.Nodes {
			out[cid] = n
		}
	}
	return out
}
