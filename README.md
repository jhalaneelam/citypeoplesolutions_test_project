## Getting Started

Follow these steps to get your project up and running.

### 1. Clone the Repository

https://github.com/jhalaneelam/citypeoplesolutions_test_project.git

# Steps to run the PostgreSQL and Redis
docker network create go_docker_compose
docker-compose up -d

### 2. To connect to the postgres choose either of the below commands:

docker exec -it go_docker_compose-postgres-1 psql -U root -d sensors_db -W
If prompted add password as root@123

                OR

docker exec -it -e PGPASSWORD=root@123 go_docker_compose-postgres-1 psql -U root -d sensors_db
