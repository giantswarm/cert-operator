package ssh

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"

	"golang.org/x/crypto/ssh"
)

// Creates a new RSA key pair with the given key length. The private key will be
// of pem format and the public key will be of OpenSSH format.
func generateRSAKeys(keyBits int) (publicKeyRsa string, privateKeyRsa string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %s", err)
	}

	privateKeyRsa = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))

	sshPublicKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %s", err)
	}
	publicKeyRsa = "ssh-rsa " + base64.StdEncoding.EncodeToString(sshPublicKey.Marshal())
	return
}

// Public key and the script to install the key are uploaded to remote machine.
// Public key is either added or removed from authorized_keys file using the
// script. Default script is for a Linux machine and hence the path of the
// authorized_keys file is hard coded to resemble Linux.
//
// The last param 'install' if false, uninstalls the key.
func (b *backend) installPublicKeyInTarget(adminUser, username, ip string, port int, hostkey, dynamicPublicKey, installScript string, install bool) error {
	// Transfer the newly generated public key to remote host under a random
	// file name. This is to avoid name collisions from other requests.
	_, publicKeyFileName, err := b.GenerateSaltedOTP()
	if err != nil {
		return err
	}

	comm, err := createSSHComm(adminUser, ip, port, hostkey)
	if err != nil {
		return err
	}
	defer comm.Close()

	err = comm.Upload(publicKeyFileName, bytes.NewBufferString(dynamicPublicKey), nil)
	if err != nil {
		return fmt.Errorf("error uploading public key: %s", err)
	}

	// Transfer the script required to install or uninstall the key to the remote
	// host under a random file name as well. This is to avoid name collisions
	// from other requests.
	scriptFileName := fmt.Sprintf("%s.sh", publicKeyFileName)
	err = comm.Upload(scriptFileName, bytes.NewBufferString(installScript), nil)
	if err != nil {
		return fmt.Errorf("error uploading install script: %s", err)
	}

	// Create a session to run remote command that triggers the script to install
	// or uninstall the key.
	session, err := comm.NewSession()
	if err != nil {
		return fmt.Errorf("unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return fmt.Errorf("invalid session object")
	}
	defer session.Close()

	authKeysFileName := fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)

	var installOption string
	if install {
		installOption = "install"
	} else {
		installOption = "uninstall"
	}

	// Give execute permissions to install script, run and delete it.
	chmodCmd := fmt.Sprintf("chmod +x %s", scriptFileName)
	scriptCmd := fmt.Sprintf("./%s %s %s %s", scriptFileName, installOption, publicKeyFileName, authKeysFileName)
	rmCmd := fmt.Sprintf("rm -f %s", scriptFileName)
	targetCmd := fmt.Sprintf("%s;%s;%s", chmodCmd, scriptCmd, rmCmd)

	session.Run(targetCmd)
	return nil
}

// Takes an IP address and role name and checks if the IP is part
// of CIDR blocks belonging to the role.
func roleContainsIP(s logical.Storage, roleName string, ip string) (bool, error) {
	if roleName == "" {
		return false, fmt.Errorf("missing role name")
	}

	if ip == "" {
		return false, fmt.Errorf("missing ip")
	}

	roleEntry, err := s.Get(fmt.Sprintf("roles/%s", roleName))
	if err != nil {
		return false, fmt.Errorf("error retrieving role '%s'", err)
	}
	if roleEntry == nil {
		return false, fmt.Errorf("role '%s' not found", roleName)
	}

	var role sshRole
	if err := roleEntry.DecodeJSON(&role); err != nil {
		return false, fmt.Errorf("error decoding role '%s'", roleName)
	}

	if matched, err := cidrListContainsIP(ip, role.CIDRList); err != nil {
		return false, err
	} else {
		return matched, nil
	}
}

// Checks if the comma separated list of CIDR blocks are all valid and they
// dont conflict with each other.
func validateCIDRList(cidrList string) (string, error) {
	// Check if the blocks are parsable
	c := strings.Split(cidrList, ",")
	for _, item := range c {
		_, _, err := net.ParseCIDR(item)
		if err != nil {
			return "", err
		}
	}

	var overlaps string
	for i := 0; i < len(c)-1; i++ {
		for j := i + 1; j < len(c); j++ {
			overlap, err := cidrOverlap(c[i], c[j])
			if err != nil {
				return "", err
			}
			if overlap {
				overlaps = fmt.Sprintf("%s [%s,%s]", overlaps, c[i], c[j])
			}
		}
	}

	return overlaps, nil
}

// Tells if the CIDR blocks overlap with eath other. Applying the mask of bigger
// block to both addresses and checking for its equality to detect an overlap.
func cidrOverlap(c1, c2 string) (bool, error) {
	ip1, net1, err := net.ParseCIDR(c1)
	if err != nil {
		return false, err
	}
	maskLen1, _ := net1.Mask.Size()

	ip2, net2, err := net.ParseCIDR(c2)
	if err != nil {
		return false, err
	}
	maskLen2, _ := net2.Mask.Size()

	// Choose the mask of bigger block
	mask := net1.Mask
	if maskLen2 < maskLen1 {
		mask = net2.Mask
	}

	return bytes.Equal(ip1.Mask(mask), ip2.Mask(mask)), nil
}

// Returns true if the IP supplied by the user is part of the comma
// separated CIDR blocks
func cidrListContainsIP(ip, cidrList string) (bool, error) {
	if len(cidrList) == 0 {
		return false, fmt.Errorf("IP does not belong to role")
	}
	for _, item := range strings.Split(cidrList, ",") {
		_, cidrIPNet, err := net.ParseCIDR(item)
		if err != nil {
			return false, fmt.Errorf("invalid CIDR entry '%s'", item)
		}
		if cidrIPNet.Contains(net.ParseIP(ip)) {
			return true, nil
		}
	}
	return false, nil
}

func createSSHComm(username, ip string, port int, hostkey string) (*comm, error) {
	signer, err := ssh.ParsePrivateKey([]byte(hostkey))
	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	connfunc := func() (net.Conn, error) {
		c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 15*time.Second)
		if err != nil {
			return nil, err
		}

		if tcpConn, ok := c.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(5 * time.Second)
		}

		return c, nil
	}
	config := &SSHCommConfig{
		SSHConfig:    clientConfig,
		Connection:   connfunc,
		Pty:          false,
		DisableAgent: true,
	}

	return SSHCommNew(fmt.Sprintf("%s:%d", ip, port), config)
}
