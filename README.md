# Note
## Database schema: 
- Points required for a particular product is stored in two different tables: "reward_perks" and "store_product". Make sure it is managed in business workflow. Allows 
flexibility but must be careful in managing this.

- Review on delete rules
- maybe add time as well to the dates

# TODO
## CI/CD 
- Include database verification during CI/CD

# Docker 
- Since there are two entry points now (server and worker) Dockerfile needs to be updated and the main.go for worker needs to be fulfilled