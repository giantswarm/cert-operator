# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Adjust helm chart to be used with `config-controller`.
- Replace `jwt-go` with `golang-jwt/jwt`.
- Manage Secrets in the same namespace in which CertConfigs are found.
- Make `expirationThreshold` configurable.

## [1.0.1] - 2021-02-23

### Fixed

- Add `list` permission for `cluster.x-k8s.io`.

## [1.0.0] - 2021-02-23

### Changed

- Update Kubernetes dependencies to 1.18 versions.
- Reconcile `CertConfig`s based on their `cert-operator.giantswarm.io/version` label.

### Removed

- Stop using the `VersionBundle` version.

### Added

- Add network policy resource.
- Added lookup for nodepool clusters in other namespaces than `default`.

## [0.1.0-2] - 2020-08-11

### Fixed

- Skip validation of reference versions like `0.1.0-1`.
- Continue to export vault token expiration time as 0 when lookup fails.

### Changed

- Update `apiextensions` to `0.4.1` version.
- Set version `0.1.0` in `project.go`.
- Use `architect` `2.1.2` in github actions.

## [0.1.0-1] - 2020-08-07

### Added

- Add `k8s-jwt-to-vault-token` init container to ensure *vault* token secret exists.
- Add Github automation workflows.

## [0.1.0] 2020-05-15

### Changed

- No longer ensure CertConfig CRD.
- Use architect-orb to release cert-operator.

### Added

- First release.

[Unreleased]: https://github.com/giantswarm/cert-operator/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/giantswarm/cert-operator/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/cert-operator/compare/v0.1.0-2...v1.0.0
[0.1.0-2]: https://github.com/giantswarm/cert-operator/compare/v0.1.0-1...v0.1.0-2
[0.1.0-1]: https://github.com/giantswarm/cert-operator/compare/v0.1.0...v0.1.0-1
[0.1.0]: https://github.com/giantswarm/cert-operator/releases/tag/v0.1.0
