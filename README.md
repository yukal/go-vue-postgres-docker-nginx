# go-vue-postgres-docker-nginx
Golang webserver starter kit

Depending on your requirements, you can customize and use the server configuration for different types of development, such as a monolith, a microservice, or a hybrid approach.

Each microservice or lambda function you can assign to a specific route as a [proxy path](data/.web/nginx/templates/default.conf.template#L55) using nginx configuration.

### Run in a development mode
1. Prepare docker images
   ```bash
   ./run up:dev
   ```

2. Start serving
   ```bash
   ./run server
   ```

   ```bash
   ./run client
   ```

3. Check if routes are available
  - Web: http://localhost:8080
  - API: http://localhost:8080/api/ping
  - Img: http://localhost:8080/img/index.jpg

4. Cleanup
   ```bash
   ./run dn:dev
   ```

### Run in a production mode
1. Prepare docker images
   ```bash
   ./run up:prod
   ```

2. Check if routes are available
  - Web: http://localhost:8080
  - API: http://localhost:8080/api/ping
  - Img: http://localhost:8080/img/index.jpg

3. Cleanup
   ```bash
   ./run dn:prod
   ```
