/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package keepalive defines configurable parameters for point-to-point
// healthcheck.
package keepalive

import (
	"time"
)

// ClientParameters is used to set keepalive parameters on the client-side.
// These configure how the client will actively probe to notice when a
// connection is broken and send pings so intermediaries will be aware of the
// liveness of the connection. Make sure these parameters are set in
// coordination with the keepalive policy on the server, as incompatible
// settings can result in closing of connection.
//
// ClientParameters是客户端keepalive相关的参数
// 这配置可以让client主动探测发现一个连接是否损坏。请确保这些参数的设置与服务器上的keepalive策略一致，
// 因为不兼容的设置可能会导致连接关闭。注：这些可以在下面EnforcementPolicy的参数中上体现
type ClientParameters struct {
	// After a duration of this time if the client doesn't see any activity it
	// pings the server to see if the transport is still alive.
	// If set below 10s, a minimum value of 10s will be used instead.
	//
	// 在Time时候内，连接没有活跃（没有收到一个包），就会开始发送ping。
	// 如果设置的值小于10s，就取10s
	Time time.Duration // The current default value is infinity. // 默认无穷大
	// After having pinged for keepalive check, the client waits for a duration
	// of Timeout and if no activity is seen even after that the connection is
	// closed.
	//
	// 在发送了ping之后，client等待Timeout的时间后，连接还是没有活跃，就会关闭连接
	Timeout time.Duration // The current default value is 20 seconds.
	// If true, client sends keepalive pings even with no active RPCs. If false,
	// when there are no active RPCs, Time and Timeout will be ignored and no
	// keepalive pings will be sent.
	//
	// 如果设为true：没有活跃的stream，也会走发送ping的逻辑
	// 如果设为false：没有活跃的stream，就会阻塞在那里，只有有新的stream被创建出来，才会继续走发送ping的逻辑
	PermitWithoutStream bool // false by default.
}

// ServerParameters is used to set keepalive and max-age parameters on the
// server-side.
type ServerParameters struct {
	// MaxConnectionIdle is a duration for the amount of time after which an
	// idle connection would be closed by sending a GoAway. Idleness duration is
	// defined since the most recent time the number of outstanding RPCs became
	// zero or the connection establishment.
	//
	// 如果客户端空闲了 t 时间， 就发送 GoAway，踢掉连接
	MaxConnectionIdle time.Duration // The current default value is infinity.
	// MaxConnectionAge is a duration for the maximum amount of time a
	// connection may exist before it will be closed by sending a GoAway. A
	// random jitter of +/-10% will be added to MaxConnectionAge to spread out
	// connection storms.
	//
	// 如果一个连接生命时间超过 t，就发送 GoAway，踢掉连接
	MaxConnectionAge time.Duration // The current default value is infinity.
	// MaxConnectionAgeGrace is an additive period after MaxConnectionAge after
	// which the connection will be forcibly closed.
	//
	// 因为MaxConnectionAge而给client发送了GoAway，可以设置在等待时间MaxConnectionAgeGrace后强制关闭连接
	MaxConnectionAgeGrace time.Duration // The current default value is infinity.
	// After a duration of this time if the server doesn't see any activity it
	// pings the client to see if the transport is still alive.
	// If set below 1s, a minimum value of 1s will be used instead.
	//
	// 在Time时候内，连接没有活跃（没有收到一个包），就会开始发送ping。
	Time time.Duration // The current default value is 2 hours.
	// After having pinged for keepalive check, the server waits for a duration
	// of Timeout and if no activity is seen even after that the connection is
	// closed.
	//
	// 在发送了ping之后，server等待Timeout的时间后，连接还是没有活跃，就会关闭连接
	Timeout time.Duration // The current default value is 20 seconds.
}

// EnforcementPolicy is used to set keepalive enforcement policy on the
// server-side. Server will close connection with a client that violates this
// policy.
type EnforcementPolicy struct {
	// MinTime is the minimum amount of time a client should wait before sending
	// a keepalive ping.
	//
	// client可以发送ping的时间间隔，如果小于这个值，server认为ping太多，会发送GoAway给client，踢掉连接
	MinTime time.Duration // The current default value is 5 minutes.
	// If true, server allows keepalive pings even when there are no active
	// streams(RPCs). If false, and client sends ping when there are no active
	// streams, server will send GOAWAY and close the connection.
	//
	// true: 没有活跃的流（active streams），也接受pings。不然将会发送GoAway给client，踢掉连接
	PermitWithoutStream bool // false by default.
}
