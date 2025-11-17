package authorize

import (
	"context"
	"fmt"
	"time"

	"github.com/neee333ko/IAM/internal/authzserver/analytics"
	"github.com/neee333ko/component-base/pkg/json"
	"github.com/neee333ko/log"
	"github.com/ory/ladon"
)

type PolicyGetter interface {
	GetPolicy(username string) ([]*ladon.DefaultPolicy, error)
}

type Authorize struct {
	pg PolicyGetter
}

func NewAuthorize(pg PolicyGetter) *Authorize {
	return &Authorize{pg: pg}
}

func (a *Authorize) FindRequestCandidates(ctx context.Context, r *ladon.Request) (ladon.Policies, error) {
	value := r.Context["username"]
	username := value.(string)

	policies, err := a.pg.GetPolicy(username)
	if err != nil {
		return nil, err
	}

	var ppolicies []ladon.Policy = make([]ladon.Policy, 0, len(policies))
	for _, item := range policies {
		ppolicies = append(ppolicies, item)
	}

	return ppolicies, nil
}

func (a *Authorize) LogRejectedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	var record analytics.AnalyticRecord

	if len(deciders) == 0 {
		record.Conclusion = "rejected because of no matched policy."
	} else {
		record.Deciders = marshal(pool)
		record.Conclusion = fmt.Sprintf("rejected by matched policy: %s\n", marshal(deciders[len(deciders)-1]))
	}

	record.Timestamp = time.Now().Format(time.RFC3339)
	record.Username = request.Context["username"].(string)
	record.Request = marshal(request)
	record.Pool = marshal(pool)
	record.Effect = ladon.DenyAccess

	record.SetExpire(0)
	if err := analytics.GetAnalytics().SendRecord(record); err != nil {
		log.Warnf("failed to write record: %s because analytics is down.", marshal(record))
	}
}

func (a *Authorize) LogGrantedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	record := analytics.AnalyticRecord{
		Timestamp:  time.Now().Format(time.RFC3339),
		Username:   request.Context["username"].(string),
		Request:    marshal(request),
		Pool:       marshal(pool),
		Deciders:   marshal(deciders),
		Conclusion: fmt.Sprintf("granted by %d policies", len(deciders)),
		Effect:     ladon.AllowAccess,
	}

	record.SetExpire(0)
	if err := analytics.GetAnalytics().SendRecord(record); err != nil {
		log.Warnf("failed to write record: %s because analytics is down.", marshal(record))
	}
}

func marshal(data any) string {
	jsonData, _ := json.Marshal(data)

	return string(jsonData)
}
