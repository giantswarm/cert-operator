# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add possibility to run cert-operator as a unique app, reconciling special version '0.0.0'.

## [2.0.1] - 2022-04-04

### Fixed

- Bump go module major version.

## [2.0.0] - 2022-03-31

### Changed

- Use v1beta1 CAPI CRDs.
- Bump `giantswarm/apiextensions` to `v6.0.0`.
- Bump `giantswarm/exporterkit` to `v1.0.0`.
- Bump `giantswarm/microendpoint` to `v1.0.0`.
- Bump `giantswarm/microerror` to `v0.4.0`.
- Bump `giantswarm/microkit` to `v1.0.0`.
- Bump `giantswarm/micrologger` to `v0.6.0`.
- Bump `giantswarm/k8sclient` to `v7.0.1`.
- Bump `giantswarm/operatorkit` to `v7.0.1`.
- Bump k8s dependencies to `v0.22.2`.
- Bump `controller-runtime` to `v0.10.3`.

## [1.3.0] - 2022-01-03

### Changed

- Use `RenewSelf` instead of `LookupSelf` to prevent expiration of `Vault token`.

## [1.2.0] - 2021-10-15

### Changed

- Introducing `v1alpha3` CR's.

### Added

- Add check to ensure that the `Cluster` resource is in the same namespace as the `certConfig` before creating the secret there.

## [1.1.0] - 2021-09-28

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

[Unreleased]: https://github.com/giantswarm/cert-operator/compare/v2.0.1...HEAD
[2.0.1]: https://github.com/giantswarm/cert-operator/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/giantswarm/cert-operator/compare/v1.3.0...v2.0.0
[1.3.0]: https://github.com/giantswarm/cert-operator/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/giantswarm/cert-operator/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/cert-operator/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/giantswarm/cert-operator/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/cert-operator/compare/v0.1.0-2...v1.0.0
[0.1.0-2]: https://github.com/giantswarm/cert-operator/compare/v0.1.0-1...v0.1.0-2
[0.1.0-1]: https://github.com/giantswarm/cert-operator/compare/v0.1.0...v0.1.0-1
[0.1.0]: https://github.com/giantswarm/cert-operator/releases/tag/v0.1.0
