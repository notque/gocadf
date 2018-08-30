package cadf

import (
	"fmt"
	"net"
	"time"
	"net/http"
	"strconv"
	"github.com/sapcc/go-bits/gopherpolicy"
	"github.com/gofrs/uuid"
)

// Event contains the CADF event according to CADF spec, section 6.6.1 Event (data)
// Extensions: requestPath (OpenStack, IBM), initiator.project_id/domain_id
// Omissions: everything that we do not use or not expose to API users
//  The JSON annotations are for parsing the result from ElasticSearch AND for generating the Hermes API response
type Event struct {
	TypeURI   string `json:"typeURI"`
	ID        string `json:"id"`
	EventTime string `json:"eventTime"`
	Action    string `json:"action"`
	EventType string `json:"eventType"`
	Outcome   string `json:"outcome"`
	Reason    Reason `json:"reason,omitempty"`
	Initiator   Resource     `json:"initiator"`
	Target      Resource     `json:"target"`
	Observer    Resource     `json:"observer"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// requestPath is an extension of OpenStack's pycadf which is supported by IBM as well
	RequestPath string `json:"requestPath,omitempty"`
}


// Resource contains attributes describing a (OpenStack-) Resource
type Resource struct {
	TypeURI   string `json:"typeURI"`
	Name      string `json:"name,omitempty"`
	Domain    string `json:"domain,omitempty"`
	ID        string `json:"id"`
	Addresses []struct {
		URL  string `json:"url"`
		Name string `json:"name,omitempty"`
	} `json:"addresses,omitempty"`
	Host Host `json:"host,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// project_id and domain_id are OpenStack extensions (introduced by Keystone and keystone(audit)middleware)
	ProjectID string `json:"project_id,omitempty"`
	DomainID  string `json:"domain_id,omitempty"`
}

// Reason contains HTTP Code and Type, and is optional in the CADF spec
type Reason struct {
	ReasonCode string `json:"reasonCode,omitempty"`
	ReasonType string `json:"reasonType,omitempty"`
}

// Host contains optional Information about the Host
type Host struct {
	ID      string `json:"id,omitempty"`
	Address string `json:"address,omitempty"`
	Agent   string `json:"agent,omitempty"`
	Platform string `json:"platform,omitempty"`
}

// Attachment contains self-describing extensions to the event
type Attachment struct {
	// Note: name is optional in CADF spec. to permit unnamed attachments
	Name string `json:"name,omitempty"`
	// this is messed-up in the spec.: the schema and examples says contentType. But the text often refers to typeURI.
	// Using typeURI would surely be more consistent. OpenStack uses typeURI, IBM supports both
	// (but forgot the name property)
	TypeURI string `json:"typeURI"`
	// Content contains the payload of the attachment. In theory this means any type.
	// In practise we have to decide because otherwise ES does based one first value
	// An interface allows arrays of json content. This should be json in the content.
	Content interface{} `json:"content"`
}

// Timestamp for proper CADF format
type Timestamp struct {
	time.Time
}

// MarshalJSON for cadf format time
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"2006-01-02T15:04:05.999Z"`)), nil
}

// UnmarshalJSON for cadf format time
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.Parse(`"2006-01-02T15:04:05.999Z"`, string(data))
	return
}

//EventParams contains parameters for creating an audit event.
type EventParams struct {
	Token        *gopherpolicy.Token
	Request      *http.Request
	ReasonCode   int
	Time         string
	ObserverUUID string
	DomainID     string
	ProjectID    string
	ServiceType  string
	ResourceName string
	RejectReason string
}

// NewEvent takes the necessary parameters and returns a new audit event.
func (p EventParams) NewEvent() Event {
	targetID := p.ProjectID
	if p.ProjectID == "" {
		targetID = p.DomainID
	}

	outcome := "failure"
	if p.ReasonCode == http.StatusOK {
		outcome = "success"
	}

	return Event{
		TypeURI:   "http://schemas.dmtf.org/cloud/audit/1.0/event", 
		ID:        generateUUID(), 
		EventTime: p.Time,
		EventType: "activity", // Activity is all we use for auditing. Activity/Monitor/Control
		Action:    "update", // Create/Update/Delete 
		Outcome:   outcome, // Success/Failure/Pending
		Reason: Reason{
			ReasonType: "HTTP",
			ReasonCode: strconv.Itoa(p.ReasonCode),
		},
		Initiator: Resource{
			TypeURI:   "service/security/account/user",
			Name:      p.Token.Context.Auth["user_name"],
			ID:        p.Token.Context.Auth["user_id"],
			Domain:    p.Token.Context.Auth["domain_name"],
			DomainID:  p.Token.Context.Auth["domain_id"],
			ProjectID: p.Token.Context.Auth["project_id"],
			Host: Host{
				Address: TryStripPort(p.Request.RemoteAddr),
				Agent:   p.Request.Header.Get("User-Agent"),
			},
		},
		Target: Resource{
			TypeURI:   fmt.Sprintf("fix/%s/%s/me", p.ServiceType, p.ResourceName),
			ID:        targetID,
			DomainID:  p.DomainID,
			ProjectID: p.ProjectID,
		},
		Observer: Resource{
			TypeURI: "service/typeURI",
			Name:    "nameofservice",
			ID:      p.ObserverUUID,
		},
		RequestPath: p.Request.URL.String(),
	}
}

//Generate an UUID based on random numbers (RFC 4122).
func generateUUID() string {
	u := uuid.Must(uuid.NewV4())
	return u.String()
}

//TryStripPort returns a host without the port number
func TryStripPort(hostPort string) string {
	host, _, err := net.SplitHostPort(hostPort)
	if err == nil {
		return host
	}
	return hostPort
}