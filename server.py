import http.server
import socketserver
import os

# Set the port number
PORT = 8000

# Specify the path to your HTML file
FILE_NAME = "source.html"

# Custom request handler to serve the specific file
class CustomHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/':
            self.path = FILE_NAME # Serve the specified file by default
        return super().do_GET()

# Set up the HTTP server with the custom handler
with socketserver.TCPServer(("", PORT), CustomHandler) as httpd:
    print(f"Serving {FILE_NAME} at http://localhost:{PORT}")
    httpd.serve_forever()
