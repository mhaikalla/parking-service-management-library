package crypts

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTrip(t *testing.T) {

	msisdn1 := "099909989897"
	msisdn2 := "776768655565"
	subsID1 := "989879"
	subsID2 := "989880"
	group1 := "9898"
	// group2 := "009900"
	groupID1 := 123

	// parent encryption
	memberCode1 := NewMemberCode(group1, msisdn1, msisdn1, subsID1, groupID1)
	packedMemberCode := memberCode1.Pack(msisdn2)
	assert.Equal(t, strings.Join([]string{group1, msisdn1, msisdn2, strconv.Itoa(groupID1)}, "|"), packedMemberCode)
	assert.NotEmpty(t, packedMemberCode)
	memberCode2 := MemberCodeCryptoUnpack(packedMemberCode, msisdn1, subsID1)
	assert.NotEmpty(t, memberCode2)
	assert.Equal(t, group1, memberCode2.GetGroup())
	assert.Equal(t, msisdn2, memberCode2.GetChild())
	assert.Equal(t, msisdn1, memberCode2.GetParent())
	assert.Equal(t, groupID1, memberCode2.GetGroupID())

	FPMemberCodeCrypto = "1"

	packedMemberCode = memberCode1.Pack(msisdn2)
	assert.NotEmpty(t, packedMemberCode)
	memberCode3 := MemberCodeCryptoUnpack(packedMemberCode, msisdn1, subsID1)
	assert.NotEmpty(t, memberCode3)
	assert.Equal(t, group1, memberCode3.GetGroup())
	assert.Equal(t, msisdn2, memberCode3.GetChild())
	assert.Equal(t, msisdn1, memberCode3.GetParent())
	assert.Equal(t, groupID1, memberCode3.GetGroupID())
	assert.False(t, memberCode3.IsParent())
	assert.False(t, memberCode3.IsChildEmpty())

	memberCode1.SetChild(msisdn1)
	memberCode1.SetChild("")
	packedMemberCode = memberCode1.Pack(msisdn1)
	assert.NotEmpty(t, packedMemberCode)
	memberCode3 = MemberCodeCryptoUnpack(packedMemberCode, msisdn1, subsID1)
	assert.NotEmpty(t, memberCode3)
	assert.Equal(t, group1, memberCode3.GetGroup())
	assert.Equal(t, msisdn1, memberCode3.GetChild())
	assert.Equal(t, msisdn1, memberCode3.GetParent())
	assert.True(t, memberCode3.IsParent())
	assert.False(t, memberCode3.IsChildEmpty())

	assert.Panics(t, func() { MemberCodeCryptoUnpack(packedMemberCode, msisdn2, subsID1) })
	assert.Panics(t, func() { MemberCodeCryptoUnpack(packedMemberCode, msisdn1, subsID2) })
	assert.Panics(t, func() { MemberCodeCryptoUnpack(packedMemberCode, msisdn2, subsID2) })

}
