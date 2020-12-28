# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Enable decomposition of numerical event outcomes into digits signed separately using different nonces.

### Fixed
- Issue with concurrent requests for an event that is not yet in the DB.

## [0.0.4] - 2020-26-10

### Changed
- Updated cryptographic module to use latest bip340.
