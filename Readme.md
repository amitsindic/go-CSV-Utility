# GoLang REST API for CSV File Management

This project is a GoLang REST API for managing CSV files. It allows users to upload a CSV file, insert its data into a MySQL database, and cache the data in Redis. Additionally, it provides endpoints for retrieving, updating, and deleting records stored in both the database and Redis.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Project Structure](#project-structure)


## Introduction

This project provides a RESTful API written in GoLang for managing CSV files. It leverages Gin framework for routing and MySQL for database operations. The API allows users to upload a CSV file, insert its data into a MySQL database, cache the data in Redis for faster retrieval, and perform CRUD operations on the stored records.

## Features

- Upload a CSV file and insert its data into MySQL database
- Cache data in Redis for faster retrieval
- Retrieve records from both database and Redis
- Update records by first name in both database and Redis
- Delete records by first name from both database and Redis

## Getting Started

### Prerequisites

Before running the project, you need to have the following installed:

- Go (version 1.21.3 or higher)
- MySQL
- Redis

### Installation

- Redis setup (https://kasunprageethdissanayake.medium.com/installing-redis-x64-3-2-100-on-windows-and-running-redis-server-94db3a98ae3d)
- Mysql setup (https://www.dataquest.io/blog/install-mysql-windows/)



# Usage

- cd choiceTechLabs
- go mod tidy
- go run cmd/main.go

- Once the server is running, you can interact with the API using any HTTP client (e.g., curl, Postman). Refer to the API endpoints     section for available endpoints and their usage.


# Endpoints
- POST /upload: (localhost:8080/upload) -> Upload a CSV file and insert its data into the database and cache.
- GET /data: (localhost:8080/data) -> Retrieve all records.
- GET /data?first_name={first_name}: (localhost:8080/data?first_name={first_name}) -> Retrives data for particular name
- PATCH /data?first_name={first_name}: (localhost:8080/data?first_name={first_name}) -> Update record by firstName in both database and Redis.
- DELETE /data?first_name={first_name}: (localhost:8080/data?first_name={first_name}) -> Delete a record by first name from both database and Redis.


# Project Structure

project_name/
  |- cmd/
  |    |- server/
  |    |    server.go
  |    |- main.go
  |- internal/
       |- config/
       |    |- config.go
       |- sql/
       |    |- script.sql
       |- handlers/
       |    |- csv_handler.go
       |    |- data_handler.go
       |    |- delete_handler.go
       |    |- update_handler.go
       |- services/
       |    |- csv_service.go
       |    |- db_service.go
       |    |- redis_service.go
       |- models/
            |- record.go