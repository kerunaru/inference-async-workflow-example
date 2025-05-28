# INFERENCE ASYNC WORKFLOW EXAMPLE

This project demonstrates a complete inference workflow in the context of a web application, showcasing best practices for interacting with third party inference APIs.

## USAGE

To run the application, execute the following command:

```bash
make build # Build the application
docker compose up # Bring up the application and its dependencies
```

Now you can access the application at `http://localhost:8080`.

The client will connect to a WebSocket at `http://localhost:8081/ws`. It will receive a response from the WebSocket server when the inference is complete.

## LICENSE

This project is licensed under the MIT License - see the LICENSE file for details.
