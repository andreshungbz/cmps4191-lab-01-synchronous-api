# CMPS4191 Laboratory 1

## Measuring a Synchronous API Contract

| Key               | Value                                                                                              |
| ----------------- | -------------------------------------------------------------------------------------------------- |
| **Student Name**  | [Andres Hung](https://github.com/andreshungbz) & [Jennessa Sierra](https://github.com/jennxsierra) |
| **Student Email** | 2018118240@ub.edu.bz & 2021153908@ub.edu.bz                                                        |
| **Course**        | CMPS4191 - Advanced Web Technologies                                                               |
| **Due Date**      | August 19, 2026                                                                                    |

## Running the Application

### Docker Compose

```
docker compose up
```

### Manual Method

#### Prerequisites

- curl
- go
- golang-migrate
- make
- PostgreSQL

#### Database Setup

```
CREATE ROLE cmps4191_lab_01_user WITH LOGIN PASSWORD 'cmps4191_lab_01_password';
CREATE DATABASE cmps4191_lab_01;
ALTER DATABASE cmps4191_lab_01 OWNER TO cmps4191_lab_01_user;
```

#### Application Setup

```
cp .envrc.example .envrc
make db/migrations/up
make run
```
