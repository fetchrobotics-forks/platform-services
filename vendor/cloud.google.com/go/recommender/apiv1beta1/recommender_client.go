// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gapic-generator. DO NOT EDIT.

package recommender

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/golang/protobuf/proto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	recommenderpb "google.golang.org/genproto/googleapis/cloud/recommender/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// CallOptions contains the retry settings for each method of Client.
type CallOptions struct {
	ListRecommendations         []gax.CallOption
	GetRecommendation           []gax.CallOption
	MarkRecommendationClaimed   []gax.CallOption
	MarkRecommendationSucceeded []gax.CallOption
	MarkRecommendationFailed    []gax.CallOption
}

func defaultClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("recommender.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultCallOptions() *CallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.3,
				})
			}),
		},
	}
	return &CallOptions{
		ListRecommendations:         retry[[2]string{"default", "idempotent"}],
		GetRecommendation:           retry[[2]string{"default", "idempotent"}],
		MarkRecommendationClaimed:   retry[[2]string{"default", "non_idempotent"}],
		MarkRecommendationSucceeded: retry[[2]string{"default", "non_idempotent"}],
		MarkRecommendationFailed:    retry[[2]string{"default", "non_idempotent"}],
	}
}

// Client is a client for interacting with Recommender API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type Client struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	client recommenderpb.RecommenderClient

	// The call options for this service.
	CallOptions *CallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewClient creates a new recommender client.
//
// Provides recommendations for cloud customers for various categories like
// performance optimization, cost savings, reliability, feature discovery, etc.
// These recommendations are generated automatically based on analysis of user
// resources, configuration and monitoring metrics.
func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn:        conn,
		CallOptions: defaultCallOptions(),

		client: recommenderpb.NewRecommenderClient(conn),
	}
	c.setGoogleClientInfo()
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *Client) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *Client) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *Client) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// ListRecommendations lists recommendations for a Cloud project. Requires the recommender.*.list
// IAM permission for the specified recommender.
func (c *Client) ListRecommendations(ctx context.Context, req *recommenderpb.ListRecommendationsRequest, opts ...gax.CallOption) *RecommendationIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ListRecommendations[0:len(c.CallOptions.ListRecommendations):len(c.CallOptions.ListRecommendations)], opts...)
	it := &RecommendationIterator{}
	req = proto.Clone(req).(*recommenderpb.ListRecommendationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*recommenderpb.Recommendation, string, error) {
		var resp *recommenderpb.ListRecommendationsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.client.ListRecommendations(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Recommendations, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.PageSize)
	it.pageInfo.Token = req.PageToken
	return it
}

// GetRecommendation gets the requested recommendation. Requires the recommender.*.get
// IAM permission for the specified recommender.
func (c *Client) GetRecommendation(ctx context.Context, req *recommenderpb.GetRecommendationRequest, opts ...gax.CallOption) (*recommenderpb.Recommendation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetRecommendation[0:len(c.CallOptions.GetRecommendation):len(c.CallOptions.GetRecommendation)], opts...)
	var resp *recommenderpb.Recommendation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.GetRecommendation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MarkRecommendationClaimed mark the Recommendation State as Claimed. Users can use this method to
// indicate to the Recommender API that they are starting to apply the
// recommendation themselves. This stops the recommendation content from being
// updated.
//
// MarkRecommendationClaimed can be applied to recommendations in CLAIMED,
// SUCCEEDED, FAILED, or ACTIVE state.
//
// Requires the recommender.*.update IAM permission for the specified
// recommender.
func (c *Client) MarkRecommendationClaimed(ctx context.Context, req *recommenderpb.MarkRecommendationClaimedRequest, opts ...gax.CallOption) (*recommenderpb.Recommendation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.MarkRecommendationClaimed[0:len(c.CallOptions.MarkRecommendationClaimed):len(c.CallOptions.MarkRecommendationClaimed)], opts...)
	var resp *recommenderpb.Recommendation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.MarkRecommendationClaimed(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MarkRecommendationSucceeded mark the Recommendation State as Succeeded. Users can use this method to
// indicate to the Recommender API that they have applied the recommendation
// themselves, and the operation was successful. This stops the recommendation
// content from being updated.
//
// MarkRecommendationSucceeded can be applied to recommendations in ACTIVE,
// CLAIMED, SUCCEEDED, or FAILED state.
//
// Requires the recommender.*.update IAM permission for the specified
// recommender.
func (c *Client) MarkRecommendationSucceeded(ctx context.Context, req *recommenderpb.MarkRecommendationSucceededRequest, opts ...gax.CallOption) (*recommenderpb.Recommendation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.MarkRecommendationSucceeded[0:len(c.CallOptions.MarkRecommendationSucceeded):len(c.CallOptions.MarkRecommendationSucceeded)], opts...)
	var resp *recommenderpb.Recommendation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.MarkRecommendationSucceeded(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MarkRecommendationFailed mark the Recommendation State as Failed. Users can use this method to
// indicate to the Recommender API that they have applied the recommendation
// themselves, and the operation failed. This stops the recommendation content
// from being updated.
//
// MarkRecommendationFailed can be applied to recommendations in ACTIVE,
// CLAIMED, SUCCEEDED, or FAILED state.
//
// Requires the recommender.*.update IAM permission for the specified
// recommender.
func (c *Client) MarkRecommendationFailed(ctx context.Context, req *recommenderpb.MarkRecommendationFailedRequest, opts ...gax.CallOption) (*recommenderpb.Recommendation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.MarkRecommendationFailed[0:len(c.CallOptions.MarkRecommendationFailed):len(c.CallOptions.MarkRecommendationFailed)], opts...)
	var resp *recommenderpb.Recommendation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.MarkRecommendationFailed(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RecommendationIterator manages a stream of *recommenderpb.Recommendation.
type RecommendationIterator struct {
	items    []*recommenderpb.Recommendation
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*recommenderpb.Recommendation, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *RecommendationIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *RecommendationIterator) Next() (*recommenderpb.Recommendation, error) {
	var item *recommenderpb.Recommendation
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *RecommendationIterator) bufLen() int {
	return len(it.items)
}

func (it *RecommendationIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
