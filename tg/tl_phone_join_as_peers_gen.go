// Code generated by gotdgen, DO NOT EDIT.

package tg

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tdp"
	"github.com/gotd/td/tgerr"
)

// No-op definition for keeping imports.
var (
	_ = bin.Buffer{}
	_ = context.Background()
	_ = fmt.Stringer(nil)
	_ = strings.Builder{}
	_ = errors.Is
	_ = sort.Ints
	_ = tdp.Format
	_ = tgerr.Error{}
)

// PhoneJoinAsPeers represents TL type `phone.joinAsPeers#afe5623f`.
//
// See https://core.telegram.org/constructor/phone.joinAsPeers for reference.
type PhoneJoinAsPeers struct {
	// Peers field of PhoneJoinAsPeers.
	Peers []PeerClass
	// Chats field of PhoneJoinAsPeers.
	Chats []ChatClass
	// Users field of PhoneJoinAsPeers.
	Users []UserClass
}

// PhoneJoinAsPeersTypeID is TL type id of PhoneJoinAsPeers.
const PhoneJoinAsPeersTypeID = 0xafe5623f

func (j *PhoneJoinAsPeers) Zero() bool {
	if j == nil {
		return true
	}
	if !(j.Peers == nil) {
		return false
	}
	if !(j.Chats == nil) {
		return false
	}
	if !(j.Users == nil) {
		return false
	}

	return true
}

// String implements fmt.Stringer.
func (j *PhoneJoinAsPeers) String() string {
	if j == nil {
		return "PhoneJoinAsPeers(nil)"
	}
	type Alias PhoneJoinAsPeers
	return fmt.Sprintf("PhoneJoinAsPeers%+v", Alias(*j))
}

// FillFrom fills PhoneJoinAsPeers from given interface.
func (j *PhoneJoinAsPeers) FillFrom(from interface {
	GetPeers() (value []PeerClass)
	GetChats() (value []ChatClass)
	GetUsers() (value []UserClass)
}) {
	j.Peers = from.GetPeers()
	j.Chats = from.GetChats()
	j.Users = from.GetUsers()
}

// TypeID returns type id in TL schema.
//
// See https://core.telegram.org/mtproto/TL-tl#remarks.
func (*PhoneJoinAsPeers) TypeID() uint32 {
	return PhoneJoinAsPeersTypeID
}

// TypeName returns name of type in TL schema.
func (*PhoneJoinAsPeers) TypeName() string {
	return "phone.joinAsPeers"
}

// TypeInfo returns info about TL type.
func (j *PhoneJoinAsPeers) TypeInfo() tdp.Type {
	typ := tdp.Type{
		Name: "phone.joinAsPeers",
		ID:   PhoneJoinAsPeersTypeID,
	}
	if j == nil {
		typ.Null = true
		return typ
	}
	typ.Fields = []tdp.Field{
		{
			Name:       "Peers",
			SchemaName: "peers",
		},
		{
			Name:       "Chats",
			SchemaName: "chats",
		},
		{
			Name:       "Users",
			SchemaName: "users",
		},
	}
	return typ
}

// Encode implements bin.Encoder.
func (j *PhoneJoinAsPeers) Encode(b *bin.Buffer) error {
	if j == nil {
		return fmt.Errorf("can't encode phone.joinAsPeers#afe5623f as nil")
	}
	b.PutID(PhoneJoinAsPeersTypeID)
	b.PutVectorHeader(len(j.Peers))
	for idx, v := range j.Peers {
		if v == nil {
			return fmt.Errorf("unable to encode phone.joinAsPeers#afe5623f: field peers element with index %d is nil", idx)
		}
		if err := v.Encode(b); err != nil {
			return fmt.Errorf("unable to encode phone.joinAsPeers#afe5623f: field peers element with index %d: %w", idx, err)
		}
	}
	b.PutVectorHeader(len(j.Chats))
	for idx, v := range j.Chats {
		if v == nil {
			return fmt.Errorf("unable to encode phone.joinAsPeers#afe5623f: field chats element with index %d is nil", idx)
		}
		if err := v.Encode(b); err != nil {
			return fmt.Errorf("unable to encode phone.joinAsPeers#afe5623f: field chats element with index %d: %w", idx, err)
		}
	}
	b.PutVectorHeader(len(j.Users))
	for idx, v := range j.Users {
		if v == nil {
			return fmt.Errorf("unable to encode phone.joinAsPeers#afe5623f: field users element with index %d is nil", idx)
		}
		if err := v.Encode(b); err != nil {
			return fmt.Errorf("unable to encode phone.joinAsPeers#afe5623f: field users element with index %d: %w", idx, err)
		}
	}
	return nil
}

// GetPeers returns value of Peers field.
func (j *PhoneJoinAsPeers) GetPeers() (value []PeerClass) {
	return j.Peers
}

// MapPeers returns field Peers wrapped in PeerClassArray helper.
func (j *PhoneJoinAsPeers) MapPeers() (value PeerClassArray) {
	return PeerClassArray(j.Peers)
}

// GetChats returns value of Chats field.
func (j *PhoneJoinAsPeers) GetChats() (value []ChatClass) {
	return j.Chats
}

// MapChats returns field Chats wrapped in ChatClassArray helper.
func (j *PhoneJoinAsPeers) MapChats() (value ChatClassArray) {
	return ChatClassArray(j.Chats)
}

// GetUsers returns value of Users field.
func (j *PhoneJoinAsPeers) GetUsers() (value []UserClass) {
	return j.Users
}

// MapUsers returns field Users wrapped in UserClassArray helper.
func (j *PhoneJoinAsPeers) MapUsers() (value UserClassArray) {
	return UserClassArray(j.Users)
}

// Decode implements bin.Decoder.
func (j *PhoneJoinAsPeers) Decode(b *bin.Buffer) error {
	if j == nil {
		return fmt.Errorf("can't decode phone.joinAsPeers#afe5623f to nil")
	}
	if err := b.ConsumeID(PhoneJoinAsPeersTypeID); err != nil {
		return fmt.Errorf("unable to decode phone.joinAsPeers#afe5623f: %w", err)
	}
	{
		headerLen, err := b.VectorHeader()
		if err != nil {
			return fmt.Errorf("unable to decode phone.joinAsPeers#afe5623f: field peers: %w", err)
		}
		for idx := 0; idx < headerLen; idx++ {
			value, err := DecodePeer(b)
			if err != nil {
				return fmt.Errorf("unable to decode phone.joinAsPeers#afe5623f: field peers: %w", err)
			}
			j.Peers = append(j.Peers, value)
		}
	}
	{
		headerLen, err := b.VectorHeader()
		if err != nil {
			return fmt.Errorf("unable to decode phone.joinAsPeers#afe5623f: field chats: %w", err)
		}
		for idx := 0; idx < headerLen; idx++ {
			value, err := DecodeChat(b)
			if err != nil {
				return fmt.Errorf("unable to decode phone.joinAsPeers#afe5623f: field chats: %w", err)
			}
			j.Chats = append(j.Chats, value)
		}
	}
	{
		headerLen, err := b.VectorHeader()
		if err != nil {
			return fmt.Errorf("unable to decode phone.joinAsPeers#afe5623f: field users: %w", err)
		}
		for idx := 0; idx < headerLen; idx++ {
			value, err := DecodeUser(b)
			if err != nil {
				return fmt.Errorf("unable to decode phone.joinAsPeers#afe5623f: field users: %w", err)
			}
			j.Users = append(j.Users, value)
		}
	}
	return nil
}

// Ensuring interfaces in compile-time for PhoneJoinAsPeers.
var (
	_ bin.Encoder = &PhoneJoinAsPeers{}
	_ bin.Decoder = &PhoneJoinAsPeers{}
)