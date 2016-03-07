// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Source: lock.vdl

// Package lock defines the interface and implementation
// for managing a physical lock.
//
// Each lock device runs a Vanadium RPC service that offers
// methods for locking and unlocking it. The Vanadium RPC protocol
// allows clients to securely communicate with this service in a
// peer to peer manner.
//
// The key for a lock is a blessing obtained during lock initialization.
// Only clients that present this blessing or extensions of it have access
// to the lock. Clients can delegate access to the lock to other principals
// by blessing them using this 'key' blessing.
package lock

import (
	// VDL system imports
	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/rpc"
	"v.io/v23/vdl"

	// VDL user imports
	"v.io/v23/security"
)

// LockStatus  indicates the status (locked or unlocked) of a lock.
type LockStatus int32

func (LockStatus) __VDLReflect(struct {
	Name string `vdl:"v.io/x/lock.LockStatus"`
}) {
}

func (m LockStatus) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if err := t.FromInt(int64(m), __VDLType_lock_v_io_x_lock_LockStatus); err != nil {
		return err
	}
	return nil
}

func (m LockStatus) MakeVDLTarget() vdl.Target {
	return nil
}

func init() {
	vdl.Register((*LockStatus)(nil))
}

var __VDLType_lock_v_io_x_lock_LockStatus *vdl.Type = vdl.TypeOf(LockStatus(0))

func __VDLEnsureNativeBuilt_lock() {
}

const Locked = LockStatus(0)

const Unlocked = LockStatus(1)

// UnclaimedLockClientMethods is the client interface
// containing UnclaimedLock methods.
//
// UnclaimedLock represents an unclaimed lock device. It is the state
// in which the lock would be after a "factory reset".
//
// Claim is used to initialize the lock and create a blessing for the caller.
//
// Once initialized, this interface will be disabled and the device will instead
// export the 'Lock' interface. Only principals that present a blessing obtained
// by a call to UnclaimedLock.Claim, or an extension of it, will be authorized.
type UnclaimedLockClientMethods interface {
	// Claim makes the device export the "Lock" interface and returns a blessing
	// bound to the caller that can be used to invoke methods on the "Lock"
	// interface.
	//
	// The 'name' is the blessing name that the device will subsequently use to
	// authenticate to its callers.
	Claim(_ *context.T, name string, _ ...rpc.CallOpt) (security.Blessings, error)
}

// UnclaimedLockClientStub adds universal methods to UnclaimedLockClientMethods.
type UnclaimedLockClientStub interface {
	UnclaimedLockClientMethods
	rpc.UniversalServiceMethods
}

// UnclaimedLockClient returns a client stub for UnclaimedLock.
func UnclaimedLockClient(name string) UnclaimedLockClientStub {
	return implUnclaimedLockClientStub{name}
}

type implUnclaimedLockClientStub struct {
	name string
}

func (c implUnclaimedLockClientStub) Claim(ctx *context.T, i0 string, opts ...rpc.CallOpt) (o0 security.Blessings, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Claim", []interface{}{i0}, []interface{}{&o0}, opts...)
	return
}

// UnclaimedLockServerMethods is the interface a server writer
// implements for UnclaimedLock.
//
// UnclaimedLock represents an unclaimed lock device. It is the state
// in which the lock would be after a "factory reset".
//
// Claim is used to initialize the lock and create a blessing for the caller.
//
// Once initialized, this interface will be disabled and the device will instead
// export the 'Lock' interface. Only principals that present a blessing obtained
// by a call to UnclaimedLock.Claim, or an extension of it, will be authorized.
type UnclaimedLockServerMethods interface {
	// Claim makes the device export the "Lock" interface and returns a blessing
	// bound to the caller that can be used to invoke methods on the "Lock"
	// interface.
	//
	// The 'name' is the blessing name that the device will subsequently use to
	// authenticate to its callers.
	Claim(_ *context.T, _ rpc.ServerCall, name string) (security.Blessings, error)
}

// UnclaimedLockServerStubMethods is the server interface containing
// UnclaimedLock methods, as expected by rpc.Server.
// There is no difference between this interface and UnclaimedLockServerMethods
// since there are no streaming methods.
type UnclaimedLockServerStubMethods UnclaimedLockServerMethods

// UnclaimedLockServerStub adds universal methods to UnclaimedLockServerStubMethods.
type UnclaimedLockServerStub interface {
	UnclaimedLockServerStubMethods
	// Describe the UnclaimedLock interfaces.
	Describe__() []rpc.InterfaceDesc
}

// UnclaimedLockServer returns a server stub for UnclaimedLock.
// It converts an implementation of UnclaimedLockServerMethods into
// an object that may be used by rpc.Server.
func UnclaimedLockServer(impl UnclaimedLockServerMethods) UnclaimedLockServerStub {
	stub := implUnclaimedLockServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implUnclaimedLockServerStub struct {
	impl UnclaimedLockServerMethods
	gs   *rpc.GlobState
}

func (s implUnclaimedLockServerStub) Claim(ctx *context.T, call rpc.ServerCall, i0 string) (security.Blessings, error) {
	return s.impl.Claim(ctx, call, i0)
}

func (s implUnclaimedLockServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implUnclaimedLockServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{UnclaimedLockDesc}
}

// UnclaimedLockDesc describes the UnclaimedLock interface.
var UnclaimedLockDesc rpc.InterfaceDesc = descUnclaimedLock

// descUnclaimedLock hides the desc to keep godoc clean.
var descUnclaimedLock = rpc.InterfaceDesc{
	Name:    "UnclaimedLock",
	PkgPath: "v.io/x/lock",
	Doc:     "// UnclaimedLock represents an unclaimed lock device. It is the state\n// in which the lock would be after a \"factory reset\".\n//\n// Claim is used to initialize the lock and create a blessing for the caller.\n//\n// Once initialized, this interface will be disabled and the device will instead\n// export the 'Lock' interface. Only principals that present a blessing obtained\n// by a call to UnclaimedLock.Claim, or an extension of it, will be authorized.",
	Methods: []rpc.MethodDesc{
		{
			Name: "Claim",
			Doc:  "// Claim makes the device export the \"Lock\" interface and returns a blessing\n// bound to the caller that can be used to invoke methods on the \"Lock\"\n// interface.\n//\n// The 'name' is the blessing name that the device will subsequently use to\n// authenticate to its callers.",
			InArgs: []rpc.ArgDesc{
				{"name", ``}, // string
			},
			OutArgs: []rpc.ArgDesc{
				{"", ``}, // security.Blessings
			},
		},
	},
}

// LockClientMethods is the client interface
// containing Lock methods.
//
// Lock is the interface for managing a physical lock.
//
// Only principals that present a blessing obtained by a call to UnclaimedLock.Claim,
// or an extension of it, will be authorized.
type LockClientMethods interface {
	// Lock locks the lock.
	Lock(*context.T, ...rpc.CallOpt) error
	// Unlock unlocks the lock.
	Unlock(*context.T, ...rpc.CallOpt) error
	// Status returns the current status (locked or unlocked) of the
	// lock.
	Status(*context.T, ...rpc.CallOpt) (LockStatus, error)
}

// LockClientStub adds universal methods to LockClientMethods.
type LockClientStub interface {
	LockClientMethods
	rpc.UniversalServiceMethods
}

// LockClient returns a client stub for Lock.
func LockClient(name string) LockClientStub {
	return implLockClientStub{name}
}

type implLockClientStub struct {
	name string
}

func (c implLockClientStub) Lock(ctx *context.T, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Lock", nil, nil, opts...)
	return
}

func (c implLockClientStub) Unlock(ctx *context.T, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Unlock", nil, nil, opts...)
	return
}

func (c implLockClientStub) Status(ctx *context.T, opts ...rpc.CallOpt) (o0 LockStatus, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Status", nil, []interface{}{&o0}, opts...)
	return
}

// LockServerMethods is the interface a server writer
// implements for Lock.
//
// Lock is the interface for managing a physical lock.
//
// Only principals that present a blessing obtained by a call to UnclaimedLock.Claim,
// or an extension of it, will be authorized.
type LockServerMethods interface {
	// Lock locks the lock.
	Lock(*context.T, rpc.ServerCall) error
	// Unlock unlocks the lock.
	Unlock(*context.T, rpc.ServerCall) error
	// Status returns the current status (locked or unlocked) of the
	// lock.
	Status(*context.T, rpc.ServerCall) (LockStatus, error)
}

// LockServerStubMethods is the server interface containing
// Lock methods, as expected by rpc.Server.
// There is no difference between this interface and LockServerMethods
// since there are no streaming methods.
type LockServerStubMethods LockServerMethods

// LockServerStub adds universal methods to LockServerStubMethods.
type LockServerStub interface {
	LockServerStubMethods
	// Describe the Lock interfaces.
	Describe__() []rpc.InterfaceDesc
}

// LockServer returns a server stub for Lock.
// It converts an implementation of LockServerMethods into
// an object that may be used by rpc.Server.
func LockServer(impl LockServerMethods) LockServerStub {
	stub := implLockServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implLockServerStub struct {
	impl LockServerMethods
	gs   *rpc.GlobState
}

func (s implLockServerStub) Lock(ctx *context.T, call rpc.ServerCall) error {
	return s.impl.Lock(ctx, call)
}

func (s implLockServerStub) Unlock(ctx *context.T, call rpc.ServerCall) error {
	return s.impl.Unlock(ctx, call)
}

func (s implLockServerStub) Status(ctx *context.T, call rpc.ServerCall) (LockStatus, error) {
	return s.impl.Status(ctx, call)
}

func (s implLockServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implLockServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{LockDesc}
}

// LockDesc describes the Lock interface.
var LockDesc rpc.InterfaceDesc = descLock

// descLock hides the desc to keep godoc clean.
var descLock = rpc.InterfaceDesc{
	Name:    "Lock",
	PkgPath: "v.io/x/lock",
	Doc:     "// Lock is the interface for managing a physical lock.\n//\n// Only principals that present a blessing obtained by a call to UnclaimedLock.Claim,\n// or an extension of it, will be authorized.",
	Methods: []rpc.MethodDesc{
		{
			Name: "Lock",
			Doc:  "// Lock locks the lock.",
		},
		{
			Name: "Unlock",
			Doc:  "// Unlock unlocks the lock.",
		},
		{
			Name: "Status",
			Doc:  "// Status returns the current status (locked or unlocked) of the\n// lock.",
			OutArgs: []rpc.ArgDesc{
				{"", ``}, // LockStatus
			},
		},
	},
}
