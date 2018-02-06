package key

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

// AllowedDomains computes a comma separated list of alternative names where the
// first item is the common name. This has to be considered in ToAltNames when
// reverse computing the list of allowed domains.
func AllowedDomains(ID, commonNameFormat string, altNames []string) string {
	commonName := fmt.Sprintf(commonNameFormat, ID)
	domains := append([]string{commonName}, altNames...)
	return strings.Join(domains, ",")
}

func ListRolesPath(ID string) string {
	return fmt.Sprintf("pki-%s/roles/", ID)
}

func ReadRolePath(ID string, organizations []string) string {
	return fmt.Sprintf("pki-%s/roles/%s", ID, RoleName(ID, organizations))
}

func RoleName(ID string, organizations []string) string {
	if len(organizations) == 0 {
		// If organizations isn't set, use the role that was created when the PKI
		// for this cluster was first setup.
		return fmt.Sprintf("role-%s", ID)
	}

	// Compute a url-safe hash of the organizations that stays the same regardless
	// of the order of the organizations supplied.
	return fmt.Sprintf("role-org-%s", computeOrgHash(organizations))
}

// ToAltNames takes a string as provided by AllowedDomains and returns the list
// of alternative names as taken by AllowedDomains. Note this implies dropping
// the first item of the parsed list.
func ToAltNames(a string) []string {
	if a == "" {
		return nil
	}

	altNames := strings.Split(a, ",")
	altNames = altNames[1:]

	return altNames
}

func ToOrganizations(o string) []string {
	if o == "" {
		return nil
	}

	organizations := strings.Split(o, ",")

	return organizations
}

func WriteRolePath(ID string, organizations []string) string {
	return fmt.Sprintf("pki-%s/roles/%s", ID, RoleName(ID, organizations))
}

// computeOrgHash computes a hash for the role that can issue these
// organizations. Since we want to reuse roles when possible, we should try to
// make sure that the same list of organizations returns the same hash
// (regardless of the order). The reason we don't use just the organizations
// that the user provided is because that could potentially be a very long list,
// or otherwise contain characters that are not allowed in URLs.
func computeOrgHash(organizations []string) string {
	sort.Strings(organizations)
	s := strings.Join(organizations, ",")

	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
