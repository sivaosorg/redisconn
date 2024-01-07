# redisconn

![GitHub contributors](https://img.shields.io/github/contributors/sivaosorg/gocell)
![GitHub followers](https://img.shields.io/github/followers/sivaosorg)
![GitHub User's stars](https://img.shields.io/github/stars/pnguyen215)

A Golang Redis connector library with features for Redis Pub/Sub, key-value operations, and distributed locking.

## Table of Contents

- [redisconn](#redisconn)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Modules](#modules)
    - [Running Tests](#running-tests)
    - [Tidying up Modules](#tidying-up-modules)
    - [Upgrading Dependencies](#upgrading-dependencies)
    - [Cleaning Dependency Cache](#cleaning-dependency-cache)

## Introduction

Welcome to the Redis Connector for Go repository! This library provides a set of tools for seamless interaction with Redis in your Go applications. It includes features for Redis Pub/Sub, key-value operations, and distributed locking.

## Prerequisites

Golang version v1.20

## Installation

- Latest version

```bash
go get -u github.com/sivaosorg/redisconn@latest
```

- Use a specific version (tag)

```bash
go get github.com/sivaosorg/redisconn@v0.0.1
```

## Modules

Explain how users can interact with the various modules.

### Running Tests

To run tests for all modules, use the following command:

```bash
make test
```

### Tidying up Modules

To tidy up the project's Go modules, use the following command:

```bash
make tidy
```

### Upgrading Dependencies

To upgrade project dependencies, use the following command:

```bash
make deps-upgrade
```

### Cleaning Dependency Cache

To clean the Go module cache, use the following command:

```bash
make deps-clean-cache
```
