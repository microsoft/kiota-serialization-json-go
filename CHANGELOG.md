# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

## [0.5.3] - 2022-06-09

### Changed

- Fixed a bug where new lines in string values would not be escaped generating invalid JSON.

## [0.5.2] - 2022-06-07

### Changed

- Upgrades abstractions and yaml dependencies.

## [0.5.1] - 2022-05-30

### Changed

 - Updated supported types for Additional Data, unsupported types now throwing an error instead of ignoring.
 - Changed logic that trims excessive commas to be called only once on serialization.

## [0.5.0] - 2022-05-26

### Changed

 - Updated reference to abstractions to support enum responses.

## [0.4.0] - 2022-05-19

### Changed

- Upgraded abstractions version.

## [0.3.2] - 2022-05-11

### Changed

- Serialization writer close method now clears the internal array and can be used to reset the writer.

## [0.3.1] - 2022-05-03

### Changed

- Fixed an issue where quotes in string values would not be escaped. #11
- Fixed an issue where int64 and byte values would get a double key. #12, #13

## [0.3.0] - 2022-04-19

### Changed

- Upgraded abstractions to 0.4.0.
- Upgraded to go 18.

## [0.2.1] - 2022-04-14

### Changed

- Fixed a bug where dates, date only, time only and duration would not serialize properly.

## [0.2.0] - 2022-04-04

### Changed

- Breaking: simplifies the field deserializers.

## [0.1.0] - 2022-03-30

### Added

- Initial tagged release of the library.
