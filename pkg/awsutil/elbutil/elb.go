package elbutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
)

type LoadBalancer struct {
	// The name of the load balancer.
	Name string

	// DNSName is the DNS name for the load balancer. CNAME records can be
	// created that point to this location.
	DNSName string

	// Load balancer scheme
	Scheme string

	// The SSL Certificate to associate with the load balancer.
	SSLCert string

	// InstancePort is the port that this load balancer forwards requests to
	// on the host.
	InstancePort int64

	InstanceStates []InstanceState

	// Tags contain the tags attached to the LoadBalancer
	Tags map[string]string
}

type InstanceState struct {
	InstanceID  string
	ReasonCode  string
	State       string
	Description string
}

type Client struct {
	client *elb.ELB
}

func New(c *elb.ELB) *Client {
	return &Client{client: c}
}

func (c *Client) ListByTags(tags map[string]string) ([]*LoadBalancer, error) {
	var (
		nextMarker *string
		lbs        []*LoadBalancer
	)

	for {
		out, err := c.client.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
			Marker:   nextMarker,
			PageSize: aws.Long(20), // Set this to 20, because DescribeTags has a limit of 20 on the LoadBalancerNames attribute.
		})
		if err != nil {
			return nil, err
		}

		if len(out.LoadBalancerDescriptions) == 0 {
			continue
		}

		// Create a names slice and descriptions map.
		names := make([]*string, len(out.LoadBalancerDescriptions))
		descs := map[string]*elb.LoadBalancerDescription{}

		for i, d := range out.LoadBalancerDescriptions {
			names[i] = d.LoadBalancerName
			descs[*d.LoadBalancerName] = d
		}

		// Find all the tags for this batch of load balancers.
		out2, err := c.client.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: names})
		if err != nil {
			return lbs, err
		}

		// Append matching load balancers to our result set.
		for _, d := range out2.TagDescriptions {
			if containsTags(tags, d.Tags) {
				lb := descs[*d.LoadBalancerName]
				var instancePort int64
				var sslCert string

				if len(lb.ListenerDescriptions) > 0 {
					instancePort = *lb.ListenerDescriptions[0].Listener.InstancePort
					for _, ld := range lb.ListenerDescriptions {
						if ld.Listener.SSLCertificateID != nil {
							sslCert = *ld.Listener.SSLCertificateID
						}
					}
				}

				// Get the instance health
				hOut, err := c.client.DescribeInstanceHealth(&elb.DescribeInstanceHealthInput{
					LoadBalancerName: d.LoadBalancerName,
				})
				if err != nil {
					return lbs, err
				}

				iStates := make([]InstanceState, len(hOut.InstanceStates))
				for i, s := range hOut.InstanceStates {
					iStates[i] = InstanceState{
						InstanceID:  *s.InstanceID,
						ReasonCode:  *s.ReasonCode,
						State:       *s.State,
						Description: *s.Description,
					}
				}

				lbs = append(lbs, &LoadBalancer{
					Name:           *lb.LoadBalancerName,
					DNSName:        *lb.DNSName,
					Scheme:         *lb.Scheme,
					SSLCert:        sslCert,
					InstancePort:   instancePort,
					Tags:           mapTags(d.Tags),
					InstanceStates: iStates,
				})
			}
		}

		nextMarker = out.NextMarker
		if nextMarker == nil || *nextMarker == "" {
			// No more items
			break
		}
	}

	return lbs, nil
}

// mapTags takes a list of []*elb.Tag's and converts them into a map[string]string
func mapTags(tags []*elb.Tag) map[string]string {
	tagMap := make(map[string]string)
	for _, t := range tags {
		tagMap[*t.Key] = *t.Value
	}

	return tagMap
}

// elbTags takes a map[string]string and converts it to the elb.Tag format.
func elbTags(tags map[string]string) []*elb.Tag {
	var e []*elb.Tag

	for k, v := range tags {
		e = append(e, elbTag(k, v))
	}

	return e
}

func elbTag(k, v string) *elb.Tag {
	return &elb.Tag{
		Key:   aws.String(k),
		Value: aws.String(v),
	}
}

// containsTags ensures that b contains all of the tags in a.
func containsTags(a map[string]string, b []*elb.Tag) bool {
	for k, v := range a {
		t := elbTag(k, v)
		if !containsTag(t, b) {
			return false
		}
	}
	return true
}

func containsTag(t *elb.Tag, tags []*elb.Tag) bool {
	for _, t2 := range tags {
		if *t.Key == *t2.Key && *t.Value == *t2.Value {
			return true
		}
	}
	return false
}
