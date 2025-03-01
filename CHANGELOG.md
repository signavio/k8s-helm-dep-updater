## [2.1.5](https://github.com/signavio/k8s-helm-dep-updater/compare/v2.1.4...v2.1.5) (2025-01-20)


### Bug Fixes

* **deps:** update golang dependencies ([#70](https://github.com/signavio/k8s-helm-dep-updater/issues/70)) ([90f427b](https://github.com/signavio/k8s-helm-dep-updater/commit/90f427bcc7b73931d07dd6889877b38c00f2a0a0))

## [2.1.4](https://github.com/signavio/k8s-helm-dep-updater/compare/v2.1.3...v2.1.4) (2025-01-08)


### Bug Fixes

* make oci login work ([#68](https://github.com/signavio/k8s-helm-dep-updater/issues/68)) ([389a882](https://github.com/signavio/k8s-helm-dep-updater/commit/389a8829ffb6da809f50617b80ff0221b52eb3da))

## [2.1.3](https://github.com/signavio/k8s-helm-dep-updater/compare/v2.1.2...v2.1.3) (2024-12-23)


### Bug Fixes

* **deps:** update module helm.sh/helm/v3 to v3.16.4 ([#67](https://github.com/signavio/k8s-helm-dep-updater/issues/67)) ([32c6f4b](https://github.com/signavio/k8s-helm-dep-updater/commit/32c6f4b75caa76cc54b8c20df40c76b2ef5f806b))

## [2.1.2](https://github.com/signavio/k8s-helm-dep-updater/compare/v2.1.1...v2.1.2) (2024-12-16)


### Bug Fixes

* **deps:** update golang dependencies to v0.32.0 ([#66](https://github.com/signavio/k8s-helm-dep-updater/issues/66)) ([82b37b4](https://github.com/signavio/k8s-helm-dep-updater/commit/82b37b4eb51e6a181a80fd1fa37f045ac14e678c))

## [2.1.1](https://github.com/signavio/k8s-helm-dep-updater/compare/v2.1.0...v2.1.1) (2024-12-15)


### Bug Fixes

* not update helm repo if there are no repos ([#64](https://github.com/signavio/k8s-helm-dep-updater/issues/64)) ([abbe2ff](https://github.com/signavio/k8s-helm-dep-updater/commit/abbe2ff9ebc470be35c0837a90ec82dbcb027e92))

# [2.1.0](https://github.com/signavio/k8s-helm-dep-updater/compare/v2.0.1...v2.1.0) (2024-12-14)


### Features

* use random helm cache dir ([#63](https://github.com/signavio/k8s-helm-dep-updater/issues/63)) ([c725450](https://github.com/signavio/k8s-helm-dep-updater/commit/c725450b42f2cf49c404881d9eaf8fc3cf318fc8))

## [2.0.1](https://github.com/signavio/k8s-helm-dep-updater/compare/v2.0.0...v2.0.1) (2024-12-02)


### Bug Fixes

* **deps:** update module github.com/stretchr/testify to v1.10.0 ([#58](https://github.com/signavio/k8s-helm-dep-updater/issues/58)) ([b5bbb19](https://github.com/signavio/k8s-helm-dep-updater/commit/b5bbb192f311722cbdae95e5440e258f4c19571a))

# [2.0.0](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.4.6...v2.0.0) (2024-12-01)


### Features

* hierarchical parallel dependency updates ([#55](https://github.com/signavio/k8s-helm-dep-updater/issues/55)) ([5e71804](https://github.com/signavio/k8s-helm-dep-updater/commit/5e71804d3ee2229427fb070a1b1d435aedf69b0d))


### BREAKING CHANGES

* Dependency updates are now executed in parallel, which may affect existing workflows relying on sequential updates.

## [1.4.6](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.4.5...v1.4.6) (2024-11-25)


### Bug Fixes

* **deps:** update golang dependencies to v0.31.3 ([54416a4](https://github.com/signavio/k8s-helm-dep-updater/commit/54416a4b580681b36e90775a97d74f57da74dc49))

## [1.4.5](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.4.4...v1.4.5) (2024-11-18)


### Bug Fixes

* **deps:** update module helm.sh/helm/v3 to v3.16.3 ([c742b4c](https://github.com/signavio/k8s-helm-dep-updater/commit/c742b4c0590890b071c810e1925a9a025ebcbc36))

## [1.4.4](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.4.3...v1.4.4) (2024-10-28)


### Bug Fixes

* **deps:** update golang dependencies to v0.31.2 ([a8caa08](https://github.com/signavio/k8s-helm-dep-updater/commit/a8caa084a8597c015a3ab5a2c98f71fe0e228266))

## [1.4.3](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.4.2...v1.4.3) (2024-10-14)


### Bug Fixes

* **deps:** update module helm.sh/helm/v3 to v3.16.2 ([e2cf290](https://github.com/signavio/k8s-helm-dep-updater/commit/e2cf290037152605f251aed42fa0ea9db9436c11))

## [1.4.2](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.4.1...v1.4.2) (2024-09-16)


### Bug Fixes

* **deps:** update golang dependencies ([7e77b9d](https://github.com/signavio/k8s-helm-dep-updater/commit/7e77b9da8ed0fe3e0642442b7e64f0a050b5937e))

## [1.4.1](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.4.0...v1.4.1) (2024-08-19)


### Bug Fixes

* **deps:** update golang dependencies ([0b79402](https://github.com/signavio/k8s-helm-dep-updater/commit/0b79402ecc5728ce06036af234c3af6995d49b5f))

# [1.4.0](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.11...v1.4.0) (2024-07-29)


### Features

* add helm repo add functionality ([5635455](https://github.com/signavio/k8s-helm-dep-updater/commit/5635455eb940bc97b6b5df3b1f3353c71bddb553))

## [1.3.11](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.10...v1.3.11) (2024-07-22)


### Bug Fixes

* **deps:** update golang dependencies to v0.30.3 ([86680b1](https://github.com/signavio/k8s-helm-dep-updater/commit/86680b1b89694b9ee28f180d068febf8b67c398b))

## [1.3.10](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.9...v1.3.10) (2024-07-15)


### Bug Fixes

* **deps:** update module helm.sh/helm/v3 to v3.15.3 ([b84e884](https://github.com/signavio/k8s-helm-dep-updater/commit/b84e8841849989f85e7e1768dfe7668cf6a673de))

## [1.3.9](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.8...v1.3.9) (2024-06-17)


### Bug Fixes

* **deps:** update golang dependencies ([ea55d90](https://github.com/signavio/k8s-helm-dep-updater/commit/ea55d90ac4bb4261f9427e9490f7a3ad0cddc310))

## [1.3.8](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.7...v1.3.8) (2024-05-27)


### Bug Fixes

* **deps:** update module helm.sh/helm/v3 to v3.15.1 ([0fad59c](https://github.com/signavio/k8s-helm-dep-updater/commit/0fad59cefab63545456bef8ff32d4753343c66d3))

## [1.3.7](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.6...v1.3.7) (2024-05-20)


### Bug Fixes

* **deps:** update golang dependencies ([aeae774](https://github.com/signavio/k8s-helm-dep-updater/commit/aeae774729269e82b5381df2cb8c44c1c9afc874))

## [1.3.6](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.5...v1.3.6) (2024-04-22)


### Bug Fixes

* **deps:** update golang dependencies to v0.30.0 ([f666b0c](https://github.com/signavio/k8s-helm-dep-updater/commit/f666b0cf2617cd369459ba68719419d9695232b8))

## [1.3.5](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.4...v1.3.5) (2024-04-15)


### Bug Fixes

* **deps:** update module helm.sh/helm/v3 to v3.14.4 ([4ce981c](https://github.com/signavio/k8s-helm-dep-updater/commit/4ce981c539d556ae9e37cee921970a5fa02b1edf))

## [1.3.4](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.3...v1.3.4) (2024-03-18)


### Bug Fixes

* **deps:** update golang dependencies ([65266fb](https://github.com/signavio/k8s-helm-dep-updater/commit/65266fbaceb8f5818c42a30bdc8dc9679ec5a0eb))

## [1.3.3](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.2...v1.3.3) (2024-03-05)


### Bug Fixes

* add darwin relase ([a040916](https://github.com/signavio/k8s-helm-dep-updater/commit/a04091627744b69e490a632c98d551730bdfd1ec))

## [1.3.2](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.1...v1.3.2) (2024-02-27)


### Bug Fixes

* trigger release of new version ([912b9f6](https://github.com/signavio/k8s-helm-dep-updater/commit/912b9f6a650f9b180b84ad92e6263a5c965a1102))

## [1.3.1](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.3.0...v1.3.1) (2024-02-27)


### Bug Fixes

* **deps:** update module helm.sh/helm/v3 to v3.14.1 [security] ([b721ae9](https://github.com/signavio/k8s-helm-dep-updater/commit/b721ae938d20e9937b0d4495f5825f97fa2bfe0f))
* **deps:** update module helm.sh/helm/v3 to v3.14.2 [security] ([5725ac0](https://github.com/signavio/k8s-helm-dep-updater/commit/5725ac01aa8e55c4f65c3be948986ddac40960fa))
* updating plugin-version and script only requires sh now ([5a01359](https://github.com/signavio/k8s-helm-dep-updater/commit/5a013592727063c34fb902c2f85f35ab4eb9a4d3))

# [1.3.0](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.2.1...v1.3.0) (2024-02-02)


### Bug Fixes

* use k8s-helm-dep-updater v1.2.1 ([857e989](https://github.com/signavio/k8s-helm-dep-updater/commit/857e98946b9c2e1fd79ce48f871d93bff648f9f6))


### Features

* update k8s-helm-dep-updater with latest url ([501475e](https://github.com/signavio/k8s-helm-dep-updater/commit/501475e2e4b02c3b3b65d1b1e59c687d829d7c9f))

## [1.2.1](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.2.0...v1.2.1) (2024-01-31)


### Bug Fixes

* update  goreleaser.yaml ([3688306](https://github.com/signavio/k8s-helm-dep-updater/commit/36883061548d643077cad99537fa8c09d33cd75a))

# [1.2.0](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.1.1...v1.2.0) (2024-01-31)


### Features

* **SIGCORE-756:** modifying uname in name template ([65e2625](https://github.com/signavio/k8s-helm-dep-updater/commit/65e2625eeffd5e76f9d28aa4723fd973b5f74789))
* **SIGCORE-756:** updating readme ([09bdc28](https://github.com/signavio/k8s-helm-dep-updater/commit/09bdc28f65c9b3b806b04733f257b770ca478da6))

## [1.1.1](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.1.0...v1.1.1) (2023-06-17)


### Bug Fixes

* configurable secret names for registries ([dcf080a](https://github.com/signavio/k8s-helm-dep-updater/commit/dcf080acd0f7ab41789591af3003c62cedd70793))

# [1.1.0](https://github.com/signavio/k8s-helm-dep-updater/compare/v1.0.0...v1.1.0) (2023-05-23)


### Features

* add helm plugin ([c2fa847](https://github.com/signavio/k8s-helm-dep-updater/commit/c2fa8475700b95b79b2f719178f11c4593990e9e))

# 1.0.0 (2023-05-23)


### Bug Fixes

* empty registries by default ([46225d2](https://github.com/signavio/k8s-helm-dep-updater/commit/46225d228b7b5c459866a0631f9bc6e8d8af3e81))


### Features

* add go releaser and setup ci ([487beee](https://github.com/signavio/k8s-helm-dep-updater/commit/487beee048c9b03d3ee20e846a66b10538b4feda))
* init helm dep updater ([585ddea](https://github.com/signavio/k8s-helm-dep-updater/commit/585ddeab088979adebe73fa6bc5a46c0e1ccfcae))
