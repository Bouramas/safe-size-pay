# Safe-size Payments Demo App

Welcome to the Safe-size Payments Demo App, an app which consists of:
- an API made with Golang 
- a UI made with plain HTML,CSS,JS for the back and front end respectively.

### Overview

The application is a prototype of the payment flow integration with [Viva's Smart Checkout](https://developer.viva.com/smart-checkout/). Users can:
- signup/login
- randomly generate an amount to be paid
- create a Viva order
- see a QR code where they can execute the payment
- see the outcome of the payment
- see the transaction history

### Key Technologies

- **Golang 1.22:** Leveraging the power and efficiency of Golang for robust backend development.
- **MySQL:** A reliable relational database management system for persistent data storage.
- **HTML,CSS,JS:** Used for a simple UI
- **🐳 Docker:** Simplifies deployment and ensures consistent environments across different platforms.

## Getting Started

### Running Locally

You need to bring up the API and the DB following the instructions above and then open 
the User Interface by opening the [index.html](ui/index.html) file in a browser.

#### Using Docker Compose 🐳

Start both MySQL and the API service with a single command:

```bash
docker-compose up -d --build mysql api
```

#### Running Separately

Alternatively, build and run the API service and MySQL individually:

```bash
# Build the image
docker build -t safe-size-pay-image .

# Run the image with environment variables and network setup
docker run --env-file .env --network safe-size-pay-stack -p 8080:8080 safe-size-pay-image
```

### Setting Up the API

1. **Configure MySQL Connection**

   Set the `MYSQL_DSN` environment variable to connect the API with MySQL. Example:

   ```bash
   export MYSQL_DSN="root:password@tcp(mysql:3306)/safe_size_db?parseTime=true&sql_mode=NO_ZERO_DATE"
   ```

2. **Adjust Build Configuration**

   Modify `GOOS` and `GOARCH` in the `Makefile` according to your local machine architecture.

3. **Build and Run**

   Use the Makefile to build the API and execute the compiled executable:

   ```bash
   make build
   ./safe-size-pay
   ```

### Additional Resources

Explore the API functionalities using the Postman collection available in the [docs](docs) folder.
   
   - Use the Signup call to create your user
   - Use the Login call first in order to acquire a token before calling the rest of the API

### TODO:

- Add Unit Tests


## ☎️ Get in Touch

I'm always open to discussions, collaborations, and feedback. If you have any questions or just want to connect, feel free to reach out!

- **Email:** gbouramas@gmail.com
- **LinkedIn:** [Giannis Bouramas](https://www.linkedin.com/in/bouramas)
