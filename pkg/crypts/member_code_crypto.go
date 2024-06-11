package crypts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

//  Contains Env Var for activating Family Plan Member Code Crypto.
var (
	FPMemberCodeCrypto = os.Getenv("FP_MEMBER_CODE_CRYPTO")
	reUUID             = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	errorLayout        = "Member Code Crypto: %v"
)

const (
	familyPlanMemberCodePrefix = `FPMC__`
)

// MemberCode contains value needed to encrypt or decrypt Member Code.
type MemberCode struct {
	group   string
	parent  string
	child   string
	subsid  string
	msisdn  string
	groupID int
}

// GetChild get child msisdn inside this `crypts#MemberCode`.
func (mc MemberCode) GetGroupID() int {
	return mc.groupID
}

// GetChild get child msisdn inside this `crypts#MemberCode`.
func (mc MemberCode) GetChild() string {
	return mc.child
}

// GetGroup get group of group inside this `crypts#MemberCode`.
func (mc MemberCode) GetGroup() string {
	return mc.group
}

// GetParent get parent msisdn inside this `crypts#MemberCode`.
func (mc MemberCode) GetParent() string {
	return mc.parent
}

// ToString output joined value found inside this `crypts#MemberCode`.
func (mc *MemberCode) ToString() string {
	return mc.group + "|" + mc.parent + "|" + mc.child + "|" + strconv.Itoa(mc.groupID)
}

// SetChild set child msisdn to this `crypts#MemberCode`,
// this is alternative to `crypts#MemberCode.Pack` method.
func (mc *MemberCode) SetChild(msisdn string) {
	if msisdn == "" {
		mc.child = uuid.New().String()
		return
	}
	mc.child = msisdn
}

// IsChildEmpty check if this Member Code generated from child perspective.
func (mc *MemberCode) IsChildEmpty() bool {
	return reUUID.MatchString(mc.child)
}

// IsParent check if this Member Code generated from parent perspective.
func (mc *MemberCode) IsParent() bool {
	return mc.parent == mc.child
}

// Pack encrypt the member code crypto, honor env var `FP_MEMBER_CODE_CRYPTO`,
// if env var set to `1`, will activate encryption, if not will call `crypts#MemberCode.ToString` method.
func (mc MemberCode) Pack(child string) string {
	mc.SetChild(child)

	if FPMemberCodeCrypto != "1" {
		return mc.ToString()
	}

	derivedKey := saltKey(mc.msisdn, mc.subsid)
	block, cipherErr := aes.NewCipher(derivedKey)
	if cipherErr != nil {
		panic(fmt.Sprintf(errorLayout, cipherErr))
	}

	gcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		panic(fmt.Sprintf(errorLayout, gcmErr))
	}

	nonce := make([]byte, gcm.NonceSize())
	dst := []byte(familyPlanMemberCodePrefix)

	if _, randErr := rand.Read(nonce); randErr != nil {
		panic(fmt.Sprintf(errorLayout, randErr))
	}

	plainText := mc.ToString()
	res := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	dst = append(dst, res...)
	return base64.RawURLEncoding.EncodeToString(dst)
}

// MemberCodeCryptoUnpack like `crypts#MemberCode.Pack` but do the reverse.
func MemberCodeCryptoUnpack(memberCode string, msisdn, subsid string) MemberCode {
	gID := 0
	if FPMemberCodeCrypto != "1" {
		splitted := strings.Split(memberCode, "|")
		if len(splitted) > 3 {
			gID, _ = strconv.Atoi(splitted[3])
		}
		return MemberCode{group: splitted[0], parent: splitted[1], child: splitted[2], groupID: gID}
	}

	buff, base64Err := base64.RawURLEncoding.DecodeString(memberCode)
	if base64Err != nil {
		panic(fmt.Sprintf(errorLayout, base64Err))
	}

	buff = buff[len(familyPlanMemberCodePrefix):]

	derivedKey := saltKey(msisdn, subsid)
	block, cipherErr := aes.NewCipher(derivedKey)
	if cipherErr != nil {
		panic(fmt.Sprintf(errorLayout, cipherErr))
	}

	gcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		panic(fmt.Sprintf(errorLayout, gcmErr))
	}
	nonceSize := gcm.NonceSize()
	nonce, data := buff[:nonceSize], buff[nonceSize:]
	dst := make([]byte, 0)

	dec, decErr := gcm.Open(dst, nonce, data, nil)
	if decErr != nil {
		panic(fmt.Sprintf(errorLayout, decErr))
	}

	splitted := strings.Split(string(dec), "|")
	if len(splitted) > 3 {
		gID, _ = strconv.Atoi(splitted[3])
	}
	return MemberCode{group: splitted[0], parent: splitted[1], child: splitted[2], msisdn: msisdn, subsid: subsid, groupID: gID}
}

// NewMemberCode create a new `*crypts#MemberCode`
func NewMemberCode(group, parent, msisdn, subsId string, groupID int) MemberCode {
	return MemberCode{
		group:   group,
		parent:  parent,
		msisdn:  msisdn,
		subsid:  subsId,
		groupID: groupID,
	}
}
