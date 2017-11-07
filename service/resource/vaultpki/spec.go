package vaultpki

type VaultPKIState struct {
	BackendExists       bool
	CACertificateExists bool
}

type ChangeType string

const (
	BackendChange       ChangeType = "BackendChange"
	CACertificateChange ChangeType = "CACertificateChange"
)
